package ordmap

import "iter"

// Map is a map which preserves the order in which the keys were inserted.
//
// Its main feature is to implement the 4 following interfaces:
//   - `json.Marshaler` (from standard package `encoding/json`)
//   - `json.Unmarshaler` (from standard package `encoding/json`)
//   - `yaml.Marshaler` (from package `gopkg.in/yaml.v3`)
//   - `yaml.Unmarshaler` (from package `gopkg.in/yaml.v3`)
//
// in a way that preserves the order of the keys in the source data.
type Map[K comparable, V any] struct {
	m    map[K]V
	keys []K
}

func (m *Map[K, V]) Get(key K) V {
	if m == nil {
		var zero V
		return zero
	}
	return m.m[key]
}

func (m *Map[K, V]) Get2(key K) (V, bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	v, ok := m.m[key]
	return v, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	if m == nil {
		panic("ordmap: assigning to nil map")
	}
	if m.m == nil {
		m.m = make(map[K]V)
	}
	if _, ok := m.m[key]; ok {
		m.m[key] = value
		return
	}

	m.m[key] = value
	m.keys = append(m.keys, key)
}

func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.m)
}

func (m *Map[K, V]) Delete(key K) bool {
	if m == nil {
		return false
	}
	if _, ok := m.m[key]; !ok {
		return false
	}

	delete(m.m, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
	return true
}

func (m *Map[K, V]) Clear() {
	if m == nil {
		return
	}
	m.m = nil
	m.keys = nil
}

func (m *Map[K, V]) Clone() *Map[K, V] {
	if m == nil {
		return nil
	}
	res := &Map[K, V]{}
	if len(m.m) == 0 {
		return res
	}

	res.m = make(map[K]V)
	for k, v := range m.m {
		res.m[k] = v
	}
	res.keys = make([]K, len(m.keys))
	copy(res.keys, m.keys)
	return res
}

func (m *Map[K, V]) Keys() []K {
	if m == nil {
		return nil
	}
	res := make([]K, len(m.keys))
	copy(res, m.keys)
	return res
}

func (m *Map[K, V]) Values() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		if m == nil {
			return
		}
		for i, k := range m.keys {
			// this is an iterator, if the map is modified while iterating, add
			// a basic handling of "that key was deleted" case
			// note: this is
			v, ok := m.m[k]
			if !ok {
				continue
			}
			if !yield(i, v) {
				return
			}
		}
	}
}

func (m *Map[K, V]) ForAll() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for _, k := range m.keys {
			if !yield(k, m.m[k]) {
				return
			}
		}
	}
}
