package ordmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	// test basic actions on a Map
	var m Map[string, int]

	assert.Equal(t, 0, m.Len())

	m.Set("a", 1)
	m.Set("c", 2)
	m.Set("b", 3)

	assert.Equal(t, []string{"a", "c", "b"}, m.Keys())
	assert.Equal(t, 1, m.Get("a"))
	assert.Equal(t, 2, m.Get("c"))
	assert.Equal(t, 3, m.Get("b"))

	m.Set("c", 4)
	assert.Equal(t, 4, m.Get("c"))
	assert.Equal(t, []string{"a", "c", "b"}, m.Keys())
	assert.Equal(t, m.Len(), 3)

	m.Delete("c")
	assert.Equal(t, 0, m.Get("c"))
	assert.Equal(t, []string{"a", "b"}, m.Keys())
	assert.Equal(t, m.Len(), 2)

	m.Set("c", 5)
	assert.Equal(t, 5, m.Get("c"))
	assert.Equal(t, []string{"a", "b", "c"}, m.Keys())
	assert.Equal(t, m.Len(), 3)

	var (
		value int
		ok    bool
	)

	value, ok = m.Get2("b")
	assert.Equal(t, 3, value)
	assert.True(t, ok)

	value, ok = m.Get2("d")
	assert.Equal(t, 0, value)
	assert.False(t, ok)
}

func TestClone(t *testing.T) {
	var m = &Map[string, int]{}
	m.Set("a", 1)
	m.Set("d", 2)
	m.Set("c", 3)
	m.Set("b", 4)

	clone := m.Clone()

	// delete a key in m:
	m.Delete("a")
	// update a key in m:
	m.Set("d", 5)
	// add a new key in m:
	m.Set("e", 6)
	// delete then re-add a key in m:
	m.Delete("c")
	m.Set("c", 7)

	assert.Equal(t, []string{"a", "d", "c", "b"}, clone.Keys())
	assert.Equal(t, []string{"d", "b", "e", "c"}, m.Keys())

	assert.Equal(t, 2, clone.Get("d"))
	assert.Equal(t, 5, m.Get("d"))
}

func TestForAll(t *testing.T) {
	var m Map[string, int]
	m.Set("a", 1)
	m.Set("d", 2)
	m.Set("c", 3)
	m.Set("b", 4)

	var keys []string
	var values []int
	for key, value := range m.ForAll() {
		keys = append(keys, key)
		values = append(values, value)
	}
	assert.Equal(t, []string{"a", "d", "c", "b"}, keys)
	assert.Equal(t, []int{1, 2, 3, 4}, values)
}

func TestNilMap(t *testing.T) {
	var m *Map[string, int]

	assert.NotPanics(t, func() {
		// the following operations should not panic on a nil *Map:
		assert.Equal(t, 0, m.Get("foo"))

		value, ok := m.Get2("foo")
		assert.Equal(t, 0, value)
		assert.False(t, ok)

		assert.Equal(t, 0, m.Len())

		m.Clear()

		assert.Nil(t, m.Clone())

		assert.Len(t, m.Keys(), 0)
		for range m.Values() {
			panic("Values should be an empty loop")
		}
		for range m.ForAll() {
			panic("ForAll should be an empty loop")
		}

	})

	assert.Panics(t, func() {
		// .Set() should panic on a nil *Map:
		m.Set("key", 42)
	})
}

func TestOrdmapKeys(t *testing.T) {
	var m Map[string, int]
	m.Set("a", 1)
	m.Set("a.bis", 2)
	m.Set("d", 3)
	m.Set("d.bis", 4)
	m.Set("c", 5)
	m.Set("c.bis", 6)
	m.Set("b", 7)
	m.Set("b.bis", 8)

	var seen []string

	for _, key := range m.Keys() {
		seen = append(seen, key)
	}
	assert.Equal(t, []string{"a", "a.bis", "d", "d.bis", "c", "c.bis", "b", "b.bis"}, seen)

	assert.NotPanics(t, func() {
		var seen []string
		for _, key := range m.Keys() {
			seen = append(seen, key)

			// delete key while iterating
			m.Delete(key + ".bis")
		}
		// t.Logf("seen: %v", seen)
	})
}

func TestOrdmapValues(t *testing.T) {
	var m Map[string, int]
	m.Set("a", 1)
	m.Set("a.bis", 2)
	m.Set("d", 3)
	m.Set("d.bis", 4)

	var seen []int

	for _, val := range m.Values() {
		seen = append(seen, val)
	}
	assert.Equal(t, []int{1, 2, 3, 4}, seen)

	assert.NotPanics(t, func() {
		var seen []int
		for _, val := range m.Values() {
			seen = append(seen, val)

			// delete keys while iterating
			switch val {
			case 1:
				m.Delete("a.bis")
			case 3:
				m.Delete("d.bis")
			}
		}
		// t.Logf("seen: %v", seen)
	})
}

func TestOrdmapForAll(t *testing.T) {
	var m Map[string, int]
	m.Set("a", 1)
	m.Set("a.bis", 2)
	m.Set("d", 3)
	m.Set("d.bis", 4)

	var keys []string
	var values []int

	for k, val := range m.ForAll() {
		keys = append(keys, k)
		values = append(values, val)
	}

	assert.Equal(t, []string{"a", "a.bis", "d", "d.bis"}, keys)
	assert.Equal(t, []int{1, 2, 3, 4}, values)
}
