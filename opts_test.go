package main

import (
	"testing"

	opt "github.com/docopt/docopt-go"
	"github.com/stretchr/testify/assert"
)

func parse(t *testing.T, args ...string) map[string]interface{} {
	if args == nil {
		args = make([]string, 0)
	}
	opts, err := opt.Parse(usage, args, false, "1.0", false, false)
	if err != nil {
		assert.Fail(t, "Failed parse "+err.Error())
	}
	return opts
}

func expectParseFailure(t *testing.T, message string, args ...string) {
	_, err := opt.Parse(usage, args, false, "1.0", true, false)
	if err == nil {
		assert.Fail(t, message)
	}
}

func TestUpdate(t *testing.T) {
	expectParseFailure(t, "update requires id", "update")
}

func TestUpdateAttribute(t *testing.T) {
	expectParseFailure(t, "update -a requires key", "update", "1", "-a")
	expectParseFailure(t, "update -a requires value", "update", "1", "-a", "key")
	opts := parse(t, "update", "1", "-a", "key", "value")
	assert.Equal(t, true, opts["-a"])
	assert.Equal(t, "key", opts["<key>"])
	assert.Equal(t, "value", opts["<value>"])
}

func TestUpdateState(t *testing.T) {
	expectParseFailure(t, "update -s requires state", "update", "1", "-s")
	opts := parse(t, "update", "1", "-s", "state")
	assert.Equal(t, true, opts["-s"])
	assert.Equal(t, "state", opts["<state>"])
}

func TestUpdateMessage(t *testing.T) {
	expectParseFailure(t, "update -m requires message", "update", "1", "-m")
	opts := parse(t, "update", "1", "-m", "a", "message")
	assert.Equal(t, true, opts["-m"])
	assert.Equal(t, []string{"a", "message"}, opts["<message>"])
	opts = parse(t, "update", "1", "-m", "message")
	assert.Equal(t, []string{"message"}, opts["<message>"])
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
