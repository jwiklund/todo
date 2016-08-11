package main

import (
	"testing"

	opt "github.com/docopt/docopt-go"
	"github.com/stretchr/testify/assert"
)

func parse(t *testing.T, args ...string) map[string]interface{} {
	opts, err := opt.Parse(usage, args, false, "1.0", false, false)
	if err != nil {
		assert.Fail(t, "Failed parse %v", err)
	}
	return opts
}

func expectParseFailure(t *testing.T, message string, args ...string) {
	_, err := opt.Parse(usage, args, false, "1.0", true, false)
	if err == nil {
		assert.Fail(t, message)
	}
}

func TestUpdateState(t *testing.T) {
	opts := parse(t, "update", "1", "state")
	assert.Equal(t, "state", opts["<state>"])
}

func TestUpdateAttributes(t *testing.T) {
	opts := parse(t, "update", "1", "-a", "key", "value")
	assert.Equal(t, "key", opts["<key>"])
	assert.Equal(t, "value", opts["<value>"])
}

func TestUpdateAttributesState(t *testing.T) {
	opts := parse(t, "update", "1", "-a", "key", "value", "state")
	assert.Equal(t, "key", opts["<key>"])
	assert.Equal(t, "value", opts["<value>"])
	assert.Equal(t, "state", opts["<state>"])
}

func TestAddMessage(t *testing.T) {
	opts := parse(t, "add", "a", "message")
	assert.Equal(t, []string{"a", "message"}, opts["<message>"])
}

func TestAddAttributeMessage(t *testing.T) {
	opts := parse(t, "add", "-a", "key", "value", "a", "message")
	assert.Equal(t, []string{"a", "message"}, opts["<message>"])
	assert.Equal(t, "key", opts["<key>"])
	assert.Equal(t, "value", opts["<value>"])
}

func TestRepo(t *testing.T) {
	expectParseFailure(t, "-r requires argument", "-r")

	opts := parse(t, "-r", "repo")
	assert.Equal(t, "repo", opts["-r"])
}

func TestConfig(t *testing.T) {
	expectParseFailure(t, "-c requires argument", "-c")

	opts := parse(t, "-c", "cfg")
	assert.Equal(t, "cfg", opts["-c"])
}

func TestOpts(t *testing.T) {
	opts := parse(t, "-a", "-v")
	assert.Equal(t, true, opts["-a"])
	assert.Equal(t, true, opts["-v"])
	opts = parse(t)
	assert.Equal(t, false, opts["-a"])
	assert.Equal(t, false, opts["-v"])
}
