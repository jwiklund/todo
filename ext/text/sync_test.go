package text

import (
	"testing"

	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/todo/fake"
	"github.com/stretchr/testify/assert"
)

func TestSyncEmpty(t *testing.T) {
	r := fake.New()
	target := &text{"text", "/tmp", false, [][]byte{}}

	if err := target.Sync(r, false); !assert.Nil(t, err) {
		return
	}
}

func TestSyncAddSingle(t *testing.T) {
	r := fake.New()
	target := &text{"text", "/tmp", false, [][]byte{[]byte("line")}}

	if err := target.Sync(r, false); !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, []todo.Task{
		todo.Task{
			ID:      "0",
			Message: "line",
			State:   todo.StateTodo,
			Attr: map[string]string{
				"external": "text",
				"text.id":  "0",
			},
		},
	}, r.MustList())
}

func TestSyncSingle(t *testing.T) {
	r := fake.New()
	r.AddWithAttr("line", map[string]string{
		"external": "text",
		"text.id":  "0",
	})
	target := &text{"text", "/tmp", false, [][]byte{[]byte("line")}}

	if err := target.Sync(r, false); !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 1, len(r.MustList()))
}

func TestSyncUpdate(t *testing.T) {
	r := fake.New()
	r.AddWithAttr("original", map[string]string{
		"external": "text",
		"text.id":  "0",
	})
	target := &text{"text", "/tmp", false, [][]byte{[]byte("update")}}

	if err := target.Sync(r, false); !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, []todo.Task{
		todo.Task{
			ID:      "0",
			Message: "update",
			State:   todo.StateTodo,
			Attr: map[string]string{
				"external": "text",
				"text.id":  "0",
			},
		},
	}, r.MustList())
}

func TestSyncDoubleLine(t *testing.T) {
	r := fake.New()
	r.AddWithAttr("original", map[string]string{
		"external": "text",
		"text.id":  "0",
	})
	target := &text{"text", "/tmp", false, [][]byte{
		[]byte("update"),
		[]byte("new"),
	}}

	if err := target.Sync(r, false); !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, "update", r.MustGet("0").Message)
	assert.Equal(t, "new", r.MustGet("1").Message)
}
