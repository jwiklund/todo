package jira

import (
	"io/ioutil"

	jira "github.com/andygrunwald/go-jira"
	"github.com/jwiklund/todo/ext/internal"
	"github.com/jwiklund/todo/todo"
	"github.com/pkg/errors"
)

func (t *extJira) Sync(r todo.RepoBegin, dryRun bool) error {
	localTasks, err := r.List()
	if err != nil {
		return err
	}

	query := "status != Done AND project = " + t.project
	if t.label != "" {
		query = query + " AND labels = " + t.label
	}
	issues, res, err := t.client.Issue.Search(query, nil)
	if err != nil {
		if res != nil {
			body, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			jiraLog.Debug("Jira response ", string(body))
		}
		return errors.Wrap(err, "Could not list issues")
	}

	externalTasks := tasksFor(t.id, issues)

	return internal.SyncHelper(r, t.id, dryRun, externalTasks, localTasks)
}

func tasksFor(extID string, issues []jira.Issue) []todo.Task {
	var res []todo.Task
	for _, issue := range issues {
		res = append(res, todo.Task{
			Message: issue.Fields.Summary,
			State:   jiraState(issue.Fields.Status.Name),
			Attr: map[string]string{
				"external":    extID,
				extID + ".id": issue.Key,
			},
		})
	}
	return res
}
