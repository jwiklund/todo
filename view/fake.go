package view

import (
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo/fake"
)

func newFake() (*fake.Fake, Todo) {
	r := fake.New()
	e, _ := ext.New(nil)
	s := State{}
	v, _ := New(r, e, s)
	return r, v
}
