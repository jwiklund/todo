package todo

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var (
	todoLog = logrus.WithField("comp", "todo")
	// ErrorNotFound error returned when task was not found.
	ErrorNotFound = errors.New("Task not found")
)

// Repo a todo repository
type Repo interface {
	List() ([]Task, error)
	Add(string, map[string]string) (Task, error)
	Get(string) (Task, error)
	GetByExternal(remoteID, externalID string) (Task, error)
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
