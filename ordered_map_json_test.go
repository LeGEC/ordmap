package ordmap

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderedMap_Json(t *testing.T) {
	// table based test cases for basic values
	type testCase struct{ input, expected string }
	table := []testCase{
		{`{}`, ""},
		{`{"a":1,"b":2,"c":3}`, ""},
		// key order must be preserved:
		{`{"c":1,"a":2,"b":3}`, ""},
		// nested objects: we will test deserializing to Map[string, any] (values are of type interface{}), so we expect the inner keys to be sorted:
		{`{"key3":[{"a1":"one","a0":true}],"key1":[{"a1":"one","a0":true}],"key2":[{"a1":"one","a0":true}]}`, `{"key3":[{"a0":true,"a1":"one"}],"key1":[{"a0":true,"a1":"one"}],"key2":[{"a0":true,"a1":"one"}]}`},
		// note: if you want to preserve the order of the inner keys, use ordmap.Any instead of 'any'
	}

	for i, tc := range table {
		var x Map[string, any]
		err := x.UnmarshalJSON([]byte(tc.input))
		require.NoError(t, err)

		bs, err := json.Marshal(x)
		require.NoError(t, err)

		expected := tc.expected
		if expected == "" {
			expected = tc.input
		}
		assert.Equal(t, expected, string(bs), "test %d: input: %s", i, tc.input)
	}
}

func TestOrderedMap_JsonErrors(t *testing.T) {
	// table based test cases for basic values
	type testCase struct{ input string }
	table := []testCase{
		{``},
		{`{`},
		{`{}0`},
		{`{"a":`},
		{`{"a":}`},
		{`{"a":1,`},
		{`{"a":1,}`},
		{`{1:"x"}`},
		{`{true:"x"}`},
	}

	for i, tc := range table {
		var x Map[string, any]
		err := x.UnmarshalJSON([]byte(tc.input))
		require.Error(t, err, "test %d: input: %s", i, tc.input)
	}
}

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

func FuzzMapUnmarshalJSON(f *testing.F) {
	f.Add(`{}`)
	f.Add(`{"a":1,"c":2,"b":3}`)
	// incorrect json:
	f.Add(``)
	f.Add(`{`)
	f.Add(`{"a":["foo","bar"`)
	f.Add(`{123:"x"}`)
	f.Add(`{}0`)
	f.Fuzz(func(t *testing.T, str string) {
		var goX map[string]any
		errGo := json.Unmarshal([]byte(str), &goX)

		var x Map[string, any]
		err := x.UnmarshalJSON([]byte(str))

		if errGo == nil {
			assert.NoError(t, err, "should not trigger an error: |%s|", str)
		} else {
			assert.Error(t, err, "should trigger an error: |%s|", str)
		}
	})
}

func FuzzMapJsonDecoder(f *testing.F) {
	f.Add(`{}`)
	f.Add(`{"a":1,"c":2,"b":3}`)
	// incorrect json:
	f.Add(``)
	f.Add(`{`)
	f.Add(`{"a":["foo","bar"`)
	f.Add(`{123:"x"}`)
	f.Fuzz(func(t *testing.T, str string) {
		var dec = json.NewDecoder(strings.NewReader(str))
		var goX map[string]any
		errGo := dec.Decode(&goX)

		dec = json.NewDecoder(strings.NewReader(str))
		var x Map[string, any]
		err := dec.Decode(&x)

		if errGo == nil {
			assert.NoError(t, err, "should not trigger an error: |%s|", str)
		} else {
			assert.Error(t, err, "should trigger an error: |%s|", str)
		}
	})
}
