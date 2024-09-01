package ordmap

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
	return m.m[key]
}

func (m *Map[K, V]) Get2(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

func (m *Map[K, V]) Set(key K, value V) {
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
	return len(m.m)
}

func (m *Map[K, V]) Delete(key K) bool {
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
	m.m = nil
	m.keys = nil
}

func (m *Map[K, V]) Clone() *Map[K, V] {
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
	res := make([]K, len(m.keys))
	copy(res, m.keys)
	return res
}
