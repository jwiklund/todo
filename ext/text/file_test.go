package text

import (
	"bytes"
	"testing"

	"github.com/jwiklund/todo/todo/fake"
	"github.com/stretchr/testify/assert"
)

func TestReadSource(t *testing.T) {
	source := func(bytes string) [][]byte {
		return newText("text", "path", []byte(bytes)).source
	}

	assert.Equal(t, 0, len(newText("text", "path", nil).source))
	assert.Equal(t, 1, len(source("")))
	assert.Equal(t, 1, len(source("1")))
	assert.Equal(t, 2, len(source("\n")))
	assert.Equal(t, 2, len(source("1\n2")))
}

func TestAddSync(t *testing.T) {
	target := text{
		id:      "text",
		path:    "path",
		updated: false,
		source:  [][]byte{},
	}
	r := fake.New()
	task := r.MustAdd("message", map[string]string{"external": "text"})
	task, err := target.Handle(task)
	if !assert.Nil(t, err) {
		return
	}
	r.MustUpdate(task)
	if e := target.Sync(r, false); !assert.Nil(t, e) {
		return
	}
	assert.Equal(t, "message", str(target.source))
	assert.Equal(t, "0", r.MustGet("0").Attr["text.id"])
}

func TestExistingAddSync(t *testing.T) {
	target := text{
		id:      "text",
		path:    "path",
		updated: false,
		source: [][]byte{
			[]byte("message1"),
		},
	}
	r := fake.New()
	task := r.MustAdd("message", map[string]string{"external": "text"})
	task, err := target.Handle(task)
	if !assert.Nil(t, err) {
		return
	}
	r.MustUpdate(task)
	if e := target.Sync(r, false); !assert.Nil(t, e) {
		return
	}
	assert.Equal(t, "message1\nmessage", str(target.source))
	assert.Equal(t, "1", r.MustGet("0").Attr["text.id"])
}

func str(lines [][]byte) string {
	return string(bytes.Join(lines, []byte("\n")))
}
