package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitUri(t *testing.T) {
	uri, user, pass, err := splitURI("http://user:pass@host/uri")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "http://host/uri", uri)
	assert.Equal(t, "user", user)
	assert.Equal(t, "pass", pass)
}
