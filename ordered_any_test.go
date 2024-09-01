package ordmap

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOrderedAny_YAMLComplex(t *testing.T) {
	// test more corner cases:
	// nested structures, arbitrary types (both in keys and values) ...
	payload := `
- key: value1
- sub1: true
  73.57:
    - 1
    - true
    - three
    - keyZ: 1
      keyA: 2
      keyY: 3
      keyB: 4
  sub2: null
  123: null
- null: value5
- -123.4:
    "1": 1
    null: null
    "0": "0"
    12: foo
- false: value4`

	var x Any
	err := yaml.Unmarshal([]byte(payload), &x)
	require.NoError(t, err)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err = enc.Encode(x)
	require.NoError(t, err)

	// there will be no indentation issue if we remove only leading and trailing spaces,
	// apply TrimSpace to remove leading and trailing "\n"
	got := strings.TrimSpace(buf.String())
	expected := strings.TrimSpace(payload)

	// we want to check that the order of the keys in each mapping is preserved
	//
	// as a proxy: we compare the serialized strings, so we also add a dependency to the exact serialization
	//   format chosen by gopkg.in/yaml.v3
	assert.Equal(t, expected, got)
}
