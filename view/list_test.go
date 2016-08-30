package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	r, v := newFake()

	r.Add("message", nil)

	if ts, e := v.List(); assert.Nil(t, e) {
		if !assert.Equal(t, 1, len(ts)) {
			return
		}

		assert.Equal(t, "message", ts[0].Message)
	}
}
