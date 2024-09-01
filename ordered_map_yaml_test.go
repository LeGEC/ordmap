package ordmap

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOrderedMapUnmarshalYaml(t *testing.T) {
	p := `
aaa:
  a: 123
  b: 234
  name: value1
ccc:
  a: 234
  b: 345
  name: value2
bbb:
  a: 345
  b: 456
  name: value3`

	type Entry struct {
		A    int    `yaml:"a"`
		B    int    `yaml:"b"`
		Name string `yaml:"name"`
	}

	var m Map[string, *Entry]

	err := yaml.Unmarshal([]byte(p), &m)
	if err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}

	// expected values:
	aaa := &Entry{123, 234, "value1"}
	ccc := &Entry{234, 345, "value2"}
	bbb := &Entry{345, 456, "value3"}

	assert.Equal(t, []string{"aaa", "ccc", "bbb"}, m.Keys())
	assert.Equal(t, aaa, m.Get("aaa"))
	assert.Equal(t, ccc, m.Get("ccc"))
	assert.Equal(t, bbb, m.Get("bbb"))
}

func TestOrderedMapMarshalYaml(t *testing.T) {
	type Entry struct {
		A    int    `yaml:"a"`
		B    int    `yaml:"b"`
		Name string `yaml:"name"`
	}

	var m Map[string, *Entry]

	m.Set("aaa", &Entry{123, 234, "value1"})
	m.Set("ccc", &Entry{234, 345, "value2"})
	m.Set("bbb", &Entry{345, 456, "value3"})

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(&m)
	require.NoError(t, err)

	expected := `
aaa:
  a: 123
  b: 234
  name: value1
ccc:
  a: 234
  b: 345
  name: value2
bbb:
  a: 345
  b: 456
  name: value3`

	got := strings.TrimSpace(buf.String())
	expected = strings.TrimSpace(expected)
	assert.Equal(t, expected, got)
}
