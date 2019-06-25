package wcg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var word = WcgWords{}

func TestWcgWords_Length(t *testing.T) {
	length := word.Length()
	assert.False(t, length < 0)
}

func TestWcgWords_FindIndex(t *testing.T) {
	assert := assert.New(t)

	w := word.FindIndex(-1)
	assert.Equal(w, "")

	w = word.FindIndex(10)
	assert.NotEqual(w, "")

	w = word.FindIndex(100000)
	assert.Equal(w, "")
}
