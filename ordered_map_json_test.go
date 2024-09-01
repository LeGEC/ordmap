package ordmap

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMapUnmarshalJson(t *testing.T) {
	p := `
{
	"aaa": {"a": 123, "b": 234, "name": "value1"},
	"ccc": {"a": 234, "b": 345, "name": "value2"},
	"bbb": {"a": 345, "b": 456, "name": "value3"}
}`

	type Entry struct {
		A    int    `json:"a"`
		B    int    `json:"b"`
		Name string `json:"name"`
	}

	var m Map[string, *Entry]

	err := json.Unmarshal([]byte(p), &m)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
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

func stripWhitespace(s string) string {
	var buf bytes.Buffer
	buf.Grow(len(s))
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			continue
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

func TestOrderedMapMarshalJson(t *testing.T) {
	type Entry struct {
		A    int    `json:"a"`
		B    int    `json:"b"`
		Name string `json:"name"`
	}

	m := Map[string, *Entry]{}
	m.Set("aaa", &Entry{123, 234, "value1"})
	m.Set("ccc", &Entry{234, 345, "value2"})
	m.Set("bbb", &Entry{345, 456, "value3"})

	got, err := json.Marshal(&m)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}

	expected := `{
	"aaa":{"a": 123, "b": 234, "name": "value1"},
	"ccc":{"a": 234, "b": 345, "name": "value2"},
	"bbb":{"a": 345, "b": 456, "name": "value3"}
}`
	expected = stripWhitespace(expected)
	if string(got) != expected {
		t.Fatalf("wrong output: expected:\n%s\n-----\ngot:\n%s\n-----\n", expected, got)
	}
}
