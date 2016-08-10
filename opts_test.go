package main

import (
	"testing"

	opt "github.com/docopt/docopt-go"
	"github.com/stretchr/testify/assert"
)

func parse(t *testing.T, args ...string) map[string]interface{} {
    opts, err := opt.Parse(usage, args, true, "1.0", false, false)
    if err != nil {
        assert.Fail(t, "Failed parse %v", err)
    }
    return opts
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
