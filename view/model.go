package view

import (
	"sort"

	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
)

// Todo a task repository view model
type Todo interface {
	List(filter func(todo.Task) bool) ([]todo.Task, error)
	Add(string, map[string]string) (todo.Task, error)
	Get(string) (todo.Task, error)
	Update(todo.Task) error

	SyncAll(dryRun bool) error
	Sync(name string, dryRun bool) error

	State() State

	Close() error
}

type view struct {
	repo  todo.RepoBegin
	ext   ext.Externals
	state State
}

// New create a new Todo
func New(repo todo.RepoBegin, ext ext.Externals, state State) (Todo, error) {
	return &view{repo, ext, state}, nil
}

// List todo tasks, with relative ids, resets relative ids
func (t *view) List(filter func(todo.Task) bool) ([]todo.Task, error) {
	ts, err := t.repo.List()
	if err != nil {
		return ts, err
	}
	filtered := make([]todo.Task, len(ts))[0:0]
	for _, t := range ts {
		if filter(t) {
			filtered = append(filtered, t)
		}
	}
	sort.Sort(sorter(filtered))
	return t.state.Remapp(filtered)
}

// Add return task with absolute id
func (t *view) Add(message string, attr map[string]string) (todo.Task, error) {
	task, err := t.repo.Add(message, attr)
	if err != nil {
		return task, err
	}
	upd, err := t.ext.Handle(task)
	if err != nil {
		return upd, err
	}
	if !task.Equal(upd) {
		return upd, t.repo.Update(upd)
	}
	return upd, nil
}

// Get from relative ID, return with relative id
func (t *view) Get(id string) (todo.Task, error) {
	aid, err := t.state.ToDB(id)
	if err != nil {
		return todo.Task{}, err
	}
	raw, err := t.repo.Get(aid)
	if err != nil {
		return raw, err
	}
	raw.ID = id
	return raw, nil
}

// Update uses task with relative ID
func (t *view) Update(task todo.Task) error {
	id, err := t.state.ToDB(task.ID)
	if err != nil {
		return err
	}
	task.ID = id
	mod, err := t.ext.Handle(task)
	if err != nil {
		return err
	}
	return t.repo.Update(mod)
}

func (t *view) Close() error {
	err1 := t.ext.Close()
	err2 := t.repo.Close()

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

func (t *view) Sync(name string, dryRun bool) error {
	return t.ext.Sync(t.repo, name, dryRun)
}

func (t *view) SyncAll(dryRun bool) error {
	return t.ext.SyncAll(t.repo, dryRun)
}

func (t *view) State() State {
	return t.state
}
