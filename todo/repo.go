package todo

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var (
	log = logrus.WithField("comp", "todo")
)

// State valid state, one of "todo", "doing", "waiting", done"
type State string

func (s State) String() string {
	return string(s)
}

// StateValid check if state is a valid state
func StateValid(state string) bool {
	return state == StateFrom(state).String()
}

// StateFrom returns a valid state from a given string (or todo)
func StateFrom(state string) State {
	switch state {
	case "todo":
		return State(state)
	case "doing":
		return State(state)
	case "waiting":
		return State(state)
	case "done":
		return State(state)
	}
	return "todo"
}

// Task a todo task
type Task struct {
	ID      string
	State   State
	Message string
}

func (t Task) String() string {
	return "(" + t.ID + ")\t" + t.Message
}

// IsCurrent return true if task is not waiting or archived
func (t Task) IsCurrent() bool {
	return t.State == "todo" || t.State == "doing"
}

// Repo a todo repository
type Repo interface {
	Close() error
	List() ([]Task, error)
	Add(string) (Task, error)
	Get(string) (Task, error)
	Update(Task) error
}

// RepoFromPath Return a repository stored in path
func RepoFromPath(path string) (Repo, error) {
	t, p := splitPath(path)
	log.Debugf("Open %s at %s", t, p)
	if t != "sqlite" {
		return nil, errors.Errorf("Invalid repo type %s", t)
	}

	return newSQL("sqlite3", p)
}
