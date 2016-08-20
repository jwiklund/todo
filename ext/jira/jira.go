package jira

import (
	"net/http"
	"time"

	"regexp"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/andygrunwald/go-jira"
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
	"github.com/pkg/errors"
)

var jiraLog = logrus.WithField("comp", "ext.jira")

func init() {
	ext.Register("jira", func(cfg ext.ExternalConfig) (ext.External, error) {
		url, user, pass, err := splitURI(cfg.URI)
		if err != nil {
			return nil, err
		}
		jiraLog.Debugf("%v", cfg)
		return New(cfg.ID, url, user, pass, cfg.Extra)
	})
}

func splitURI(uri string) (string, string, string, error) {
	r := regexp.MustCompile("^(https?://)([^:/@]+):([^@/]+)@(.*)$")
	match := r.FindStringSubmatch(uri)
	if match == nil {
		return "", "", "", errors.New("invalid uri")
	}
	return match[1] + match[4], match[2], match[3], nil
}

// New create a new jira mirror
func New(id, url, user, pass string, extra map[string]string) (ext.External, error) {

	project, ok := extra["project"]
	if !ok {
		return nil, errors.New("project is required for jira")
	}
	label, _ := extra["label"]
	stateTransitions := map[string]string{}
	for _, state := range todo.States {
		if transition, ok := extra[state.String()+"_transition"]; ok {
			stateTransitions[state.String()] = transition
		}
	}

	return &extJira{
		id:          id,
		project:     project,
		label:       label,
		transitions: stateTransitions,
		cli:         nil,
		clientFn: func() (*jira.Client, error) {
			httpClient := http.Client{
				Timeout: time.Duration(10 * time.Second),
			}
			jiraClient, err := jira.NewClient(&httpClient, url)
			if err != nil {
				return nil, errors.Wrap(err, "Could not create jira client")
			}
			_, err = jiraClient.Authentication.AcquireSessionCookie(user, pass)
			if err != nil {
				return nil, errors.Wrap(err, "Could not login to jira")
			}
			return jiraClient, nil
		}}, nil
}

type extJira struct {
	id          string
	project     string
	label       string
	transitions map[string]string
	cli         *jira.Client
	clientFn    func() (*jira.Client, error)
}

func (t *extJira) client() (*jira.Client, error) {
	if t.cli != nil {
		return t.cli, nil
	}
	cli, err := t.clientFn()
	t.cli = cli
	return t.cli, err
}

func (t *extJira) Handle(task todo.Task) (todo.Task, error) {
	if a := task.Attr["external"]; a != t.id {
		return task, nil
	}
	if task.Message == "" {
		return task, errors.New("message is required for jira tasks")
	}
	client, err := t.client()
	if err != nil {
		return task, err
	}
	if a := task.Attr[t.id+".id"]; a != "" {
		issue, res, err := client.Issue.Get(a)
		if err != nil {
			if res != nil {
				res.Body.Close()
			}
			return task, errors.Wrap(err, "Could not get jira issue")
		}
		if issue.Fields.Summary != task.Message {
			if err := t.updateJiraFields(a, task.Message, task.Attr); err != nil {
				return task, err
			}
		}
		if issue.Fields.Status.Name != jiraStatus(task.State) {
			if err := t.updateJiraStatus(a, task.State); err != nil {
				return task, err
			}
		}
		return task, nil
	}
	var labels []string
	if t.label != "" {
		labels = []string{t.label}
	}
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Summary: task.Message,
			Type: jira.IssueType{
				Name: "Story",
			},
			Project: jira.Project{
				Key: "EX",
			},
			Labels: labels,
		},
	}
	i, res, err := client.Issue.Create(&issue)
	defer res.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		jiraLog.Debug("jira response ", string(body))
		return task, errors.Wrap(err, "Could not create jira issue")
	}
	task.Attr[t.id+".id"] = i.Key
	return task, nil
}

func (t *extJira) Close() error {
	return nil
}

func (t *extJira) updateJiraFields(extID string, message string, attr map[string]string) error {
	updated := struct {
		Fields struct {
			Summary string `json:"summary,omitempty"`
		} `json:"fields"`
	}{}
	client, err := t.client()
	if err != nil {
		return err
	}
	req, err := client.NewRequest("PUT", "/rest/api/2/issue/"+extID, updated)
	if err != nil {
		return errors.Wrap(err, "Could not create put request")
	}
	res, err := client.Do(req, nil)
	defer res.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		jiraLog.Debug("Jira response ", string(body))
		return errors.Wrap(err, "Could not update jira issue")
	}
	jiraLog.Debugf("%+v", res)
	return nil
}

func jiraStatus(state todo.State) string {
	switch state {
	case todo.StateTodo:
		return "To Do"
	case todo.StateDoing:
		return "In Progress"
	case todo.StateWaiting:
		return "To Do"
	case todo.StateDone:
		return "Done"
	default:
		return "To Do"
	}
}

func (t *extJira) updateJiraStatus(extID string, state todo.State) error {
	transition := struct {
		Transition struct {
			ID string `json:"id"`
		} `json:"transition"`
	}{}
	transitionID, ok := t.transitions[state.String()]
	if !ok {
		return errors.New("No transition for state " + state.String())
	}
	transition.Transition.ID = transitionID
	client, err := t.client()
	if err != nil {
		return err
	}
	path := "/rest/api/2/issue/" + extID + "/transitions"
	req, err := client.NewRequest("POST", path, &transition)
	if err != nil {
		return errors.Wrap(err, "Could not create update request")
	}
	res, err := client.Do(req, nil)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "Could not update state")
	}
	return nil
}
