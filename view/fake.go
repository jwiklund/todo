package view

import (
	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/todo/fake"
)

type fakExt struct {
}

func (e *fakExt) Handle(t todo.Task) (todo.Task, error) {
	return t, nil
}

func (e *fakExt) Sync(todo.RepoBegin, string, bool) error {
	return nil
}

func (e *fakExt) SyncAll(todo.RepoBegin, bool) error {
	return nil
}

func (e *fakExt) Close() error {
	return nil
}

func newFake() (*fake.Fake, Todo) {
	r := fake.New()
	e := fakExt{}
	s := State{}
	v, _ := New(r, &e, s)
	return r, v
}
