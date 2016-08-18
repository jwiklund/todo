package main

import (
	"testing"

	"strings"

	"github.com/jwiklund/todo/ext"
	"github.com/stretchr/testify/assert"
)

func TestToml(t *testing.T) {
	c, e := readConfigToml(strings.NewReader(`
	repo = "repo"
	
	[id1]
	uri = "uri1"
	type = "type"
    key="value"
    
	[id2]
	uri = "uri2"
	type = "type"
    key="value"
    `))

	if !assert.Nil(t, e) {
		return
	}
	assert.Equal(t, config{
		External: []ext.ExternalConfig{
			ext.ExternalConfig{
				ID:   "id1",
				Type: "type",
				URI:  "uri1",
				Extra: map[string]string{
					"key": "value",
				},
			},
			ext.ExternalConfig{
				ID:   "id2",
				Type: "type",
				URI:  "uri2",
				Extra: map[string]string{
					"key": "value",
				},
			},
		},
		Repo: "repo",
	}, c)
}

func TestInvalidExternal(t *testing.T) {
	_, e := readConfigToml(strings.NewReader(`
	repo = "repo"
	
	[id1]
	`))
	assert.NotNil(t, e)
}
