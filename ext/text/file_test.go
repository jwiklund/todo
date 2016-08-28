package text

import (
	"testing"

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
