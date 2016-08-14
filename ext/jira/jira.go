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
		return New(cfg.ID, url, user, pass)
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
func New(id, url, user, pass string) (ext.External, error) {
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

	return &extJira{id, jiraClient}, nil
}

type extJira struct {
	id     string
	client *jira.Client
}

func (t *extJira) Handle(task todo.Task) (todo.Task, error) {
	if a := task.Attr["external"]; a != t.id {
		return task, nil
	}
	if task.Message == "" {
		return task, errors.New("message is required for jira tasks")
	}
	if a := task.Attr[t.id+".id"]; a != "" {
		issue, _, err := t.client.Issue.Get(a)
		if err != nil {
			return task, errors.Wrap(err, "Could not get jira issue")
		}
		if issue.Fields.Summary != task.Message {
			updated := struct {
				Fields struct {
					Summary string `json:"summary"`
				} `json:"fields"`
			}{}
			updated.Fields.Summary = task.Message
			req, _ := t.client.NewRequest("PUT", "/rest/api/2/issue/"+a, updated)
			res, err := t.client.Do(req, nil)
			defer res.Body.Close()
			if err != nil {
				body, _ := ioutil.ReadAll(res.Body)
				jiraLog.Debug("Jira response ", string(body))
				return task, errors.Wrap(err, "Could not update jira issue")
			}
			jiraLog.Debugf("%+v", res)
			return task, nil
		}
	} else {
		issue := jira.Issue{
			Fields: &jira.IssueFields{
				Summary: task.Message,
				Type: jira.IssueType{
					Name: "Story",
				},
				Project: jira.Project{
					Key: "EX",
				},
			},
		}
		i, res, err := t.client.Issue.Create(&issue)
		defer res.Body.Close()
		if err != nil {
			body, _ := ioutil.ReadAll(res.Body)
			jiraLog.Debug("jira response ", string(body))
			return task, errors.Wrap(err, "Could not create jira issue")
		}
		task.Attr[t.id+".id"] = i.Key
		return task, nil
	}
	return task, errors.New("Invalid state")
}

func (t *extJira) Sync(r todo.RepoBegin) error {
	return errors.New("Not implemented")
}

func (t *extJira) Close() error {
	return nil
}
