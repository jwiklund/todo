package internal

import (
	"testing"

	"github.com/jwiklund/todo/todo"
	"github.com/stretchr/testify/assert"
)

func newChange() Change {
	return Change{
		map[string]string{},
		map[string]string{},
		[]string{},
	}
}

func TestDiffEmpty(t *testing.T) {
	t1 := todo.Task{Message: "message"}
	t2 := todo.Task{Attr: map[string]string{"key": "value"}}

	assert.Equal(t, newChange(), Compare(t1, t1))
	assert.Equal(t, newChange(), Compare(t2, t2))
}

func TestDiffMod(t *testing.T) {
	t1 := todo.Task{Message: "message"}
	t2 := todo.Task{Attr: map[string]string{"key": "value"}}
	c1 := newChange()
	c1.Added["message"] = "message"
	c2 := newChange()
	c2.Added["key"] = "value"

	assert.Equal(t, c1, Compare(todo.Task{}, t1))
	assert.Equal(t, c2, Compare(todo.Task{}, t2))
}

func TestDiffRem(t *testing.T) {
	t1 := todo.Task{Message: "message"}
	t2 := todo.Task{Attr: map[string]string{"key": "value"}}
	c1 := newChange()
	c1.Removed = []string{"message"}
	c2 := newChange()
	c2.Removed = []string{"key"}

	assert.Equal(t, c1, Compare(t1, todo.Task{}))
	assert.Equal(t, c2, Compare(t2, todo.Task{}))
}

func TestApplyAdded(t *testing.T) {
	t1 := todo.Task{Attr: map[string]string{}}
	c1 := newChange()
	c1.Added["message"] = "message"
	c1.Added["key"] = "value"

	c1.Apply(&t1)

	assert.Equal(t, todo.Task{
		Message: "message",
		Attr:    map[string]string{"key": "value"},
	}, t1)
}

func TestApplyModify(t *testing.T) {
	t1 := todo.Task{Message: "message", Attr: map[string]string{"key": "value"}}
	c1 := newChange()
	c1.Modified["message"] = "message1"
	c1.Removed = []string{"key"}

	c1.Apply(&t1)

	assert.Equal(t, todo.Task{
		Message: "message1",
		Attr:    map[string]string{},
	}, t1)
}
