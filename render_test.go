package main

import (
	"testing"

	"bytes"

	"github.com/jwiklund/todo/todo"
	"github.com/stretchr/testify/assert"
)

func TestRenderOne(t *testing.T) {
	bs := bytes.Buffer{}
	renderOne(todo.Task{
		ID:      "0",
		State:   todo.StateTodo,
		Message: "message",
	}, &bs)
	assert.Equal(t, "(0)   none  todo  message\n", bs.String())
}

func TestRenderOneAttr(t *testing.T) {
	bs := bytes.Buffer{}
	renderOne(todo.Task{
		ID:      "0",
		State:   todo.StateDone,
		Message: "message",
		Attr: map[string]string{
			"key": "value",
		},
	}, &bs)
	assert.Equal(t, "(0)   none  done  message\n      key   value\n", bs.String())
}

func TestRenderList(t *testing.T) {
	bs := bytes.Buffer{}
	renderList([]todo.Task{todo.Task{
		ID:      "0",
		State:   todo.StateTodo,
		Message: "message",
	}}, &bs)
	assert.Equal(t, "(0)   none  todo  message\n", bs.String())
}
