package view

import (
	"testing"

	"github.com/jwiklund/todo/todo"
	"github.com/stretchr/testify/assert"
)

func TestSortById(t *testing.T) {
	r, v := newFake()

	r.Add("message1", nil)
	r.Add("message2", nil)

	ts, _ := v.List()
	assert.Equal(t, []string{"message1", "message2"}, messages(ts))
}

func TestSortWithPriority(t *testing.T) {
	r, v := newFake()

	r.Add("message1", map[string]string{"prio": "2"})
	r.Add("message2", map[string]string{"prio": "1"})

	ts, _ := v.List()
	assert.Equal(t, []string{"message2", "message1"}, messages(ts))
}

func TestSortWithMixedPriority(t *testing.T) {
	r, v := newFake()

	r.Add("message1", map[string]string{"prio": "3"})
	r.Add("message2", nil)
	r.Add("message3", map[string]string{"prio": "1"})

	ts, _ := v.List()
	assert.Equal(t, []string{"message3", "message1", "message2"}, messages(ts))
}

func messages(ts []todo.Task) []string {
	res := []string{}
	for _, t := range ts {
		res = append(res, t.Message)
	}
	return res
}
