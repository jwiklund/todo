package todo

import (
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var (
	todoLog = logrus.WithField("comp", "todo")
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
	Attr    map[string]string
}

func (t Task) String() string {
	return "(" + t.ID + ")\t" + t.State.String() + "\t" + t.Message
}

// Equal check if tasks are equal
func (t Task) Equal(t2 Task) bool {
	return reflect.DeepEqual(t, t2)
}

// IsCurrent return true if task is not waiting or archived
func (t Task) IsCurrent() bool {
	return t.State == "todo" || t.State == "doing"
}

// Repo a todo repository
type Repo interface {
	List() ([]Task, error)
	Add(string) (Task, error)
	AddWithAttr(string, map[string]string) (Task, error)
	Get(string) (Task, error)
	Update(Task) error

	Close() error
}

// RepoBegin with begin
type RepoBegin interface {
	Repo
	Begin() (RepoCommit, error)
}

// RepoCommit with commit
type RepoCommit interface {
	Repo
	Commit() error
}

// RepoFromPath Return a repository stored in path
func RepoFromPath(path string) (RepoBegin, error) {
	t, p := splitPath(path)
	todoLog.Debugf("Open %s at %s", t, p)
	if t != "sqlite" {
		return nil, errors.Errorf("Invalid repo type %s", t)
	}

	return newSQL("sqlite3", p)
}
