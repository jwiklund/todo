package main

import (
	"testing"

	"strings"

	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/view"
	"github.com/stretchr/testify/assert"
)

func TestToml(t *testing.T) {
	c, e := readConfigToml(strings.NewReader(`
	repo = "repo"
	
	[id1]
	uri = "uri1"
	type = "type"
    key1 = "value1"
    
	[id2]
	uri = "uri2"
	type = "type"
    key2 = "value2"
    `))

	if !assert.Nil(t, e) {
		return
	}
	assert.Contains(t, c.External, ext.ExternalConfig{
		ID:   "id1",
		Type: "type",
		URI:  "uri1",
		Extra: map[string]string{
			"key1": "value1",
		},
	})
	assert.Contains(t, c.External, ext.ExternalConfig{
		ID:   "id2",
		Type: "type",
		URI:  "uri2",
		Extra: map[string]string{
			"key2": "value2",
		},
	})
	assert.Equal(t, "repo", c.Repo)
}

func TestInvalidExternal(t *testing.T) {
	_, e := readConfigToml(strings.NewReader(`
	repo = "repo"
	
	[id1]
	`))
	assert.NotNil(t, e)
}

func TestWriteState(t *testing.T) {
	state := view.State{
		Mapping: map[string]string{
			"0": "1",
		},
	}
	bs := bytes.Buffer{}
	e := toml.NewEncoder(&bs)
	if !assert.Nil(t, e.Encode(state)) {
		return
	}
	assert.Equal(t, "[Mapping]\n  0 = \"1\"\n", bs.String())
}
