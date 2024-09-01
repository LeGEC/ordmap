package ordmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetDelete(t *testing.T) {
	// test basic actions on a Map
	var m Map[string, int]
	m.Set("a", 1)
	m.Set("c", 2)
	m.Set("b", 3)

	assert.Equal(t, 1, m.Get("a"))
	assert.Equal(t, 2, m.Get("c"))
	assert.Equal(t, 3, m.Get("b"))
	assert.Equal(t, []string{"a", "c", "b"}, m.Keys())

	m.Set("c", 4)
	assert.Equal(t, 4, m.Get("c"))
	assert.Equal(t, []string{"a", "c", "b"}, m.Keys())

	m.Delete("c")
	assert.Equal(t, 0, m.Get("c"))
	assert.Equal(t, []string{"a", "b"}, m.Keys())

	m.Set("c", 5)
	assert.Equal(t, 5, m.Get("c"))
	assert.Equal(t, []string{"a", "b", "c"}, m.Keys())
}
