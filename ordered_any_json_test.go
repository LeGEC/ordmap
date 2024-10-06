package ordmap

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderedAny_Json(t *testing.T) {
	// table based test cases for basic values
	type testCase struct{ input, expected string }
	table := []testCase{
		{`null`, `null`},
		{`true`, `true`},
		{`false`, `false`},
		{`0`, `0`},
		{`10`, `10`},
		{`-3`, `-3`},
		{`2.3`, `2.3`},
		{`-4.5`, `-4.5`},
		{`1.2e3`, "1200"},
		{`-3.4e-5`, "-0.000034"},
		{`"toto"`, `"toto"`},
		{`[]`, `[]`},
		{`{}`, `{}`},
		{`[1,2,3]`, `[1,2,3]`},
		{`{"a":1,"b":2,"c":3}`, `{"a":1,"b":2,"c":3}`},
		// key order must be preserved:
		{`{"c":1,"a":2,"b":3}`, `{"c":1,"a":2,"b":3}`},
		// nested objects
		{`{"key":[{"a2":"one","a0":true,"a1":3},{"a3":false,"a5":5,"a4":"six"}]}`, `{"key":[{"a2":"one","a0":true,"a1":3},{"a3":false,"a5":5,"a4":"six"}]}`},
	}

	for i, tc := range table {
		var x Any
		err := x.UnmarshalJSON([]byte(tc.input))
		require.NoError(t, err)

		bs, err := json.Marshal(x)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, string(bs), "test %d: input: %s", i, tc.input)
	}
}

func TestOrderedAny_IncorrectJson(t *testing.T) {
	// table based test cases for basic values
	type testCase struct{ input string }
	table := []testCase{
		{`{`},
		{`[`},
		{`fal`},
		{`tru`},
		{`00`},
		{`+1`},
		{`1e1000`}, // technically a valid json string, but too large to Unmarshal into a float64 (note: stdlib also chokes on it when unmarshaling to an 'any')
		{`"foo`},
		{`[`},
		{`{`},
		{`[1,2,3,]`},
		{`{"a":1,"b":2,"c":}`},
		{`{"a":1,"b":2,"c":3,}`},
		{`{"a":[1,2,3}`},
		{`{"a":[1,2,3`},
	}

	for i, tc := range table {
		var x Any
		err := x.UnmarshalJSON([]byte(tc.input))
		require.Error(t, err, "test %d: input: %s", i, tc.input)
	}
}

// copyJson[T] will do an Encode -> Decode roundtrip, and Decode to a value of type T
func copyJson[T any](t *testing.T, x any) T {
	bs, err := json.Marshal(x)
	require.NoError(t, err)

	var target T
	err = json.Unmarshal(bs, &target)
	require.NoError(t, err)
	return target
}

func TestOrderedAny_UnmarshalJSON_Complex(t *testing.T) {
	input := `[
	{
		"z":1,
		"y":[
			{"x": true, "w" :null},
			{"v": false, "u": "string"}
		]
	},
	{"t": 33, "s": null, "r": 42}
]`
	expected := []any{
		map[string]any{
			"z": float64(1),
			"y": []any{
				map[string]any{"x": true, "w": nil},
				map[string]any{"v": false, "u": "string"},
			},
		},
		map[string]any{"t": float64(33), "s": nil, "r": float64(42)},
	}

	var obj Any
	err := json.Unmarshal([]byte(input), &obj)
	require.NoError(t, err)

	got := copyJson[any](t, obj.v)
	assert.Equal(t, expected, got)
}

func TestOrderedAny_MarshalJSON_Complex(t *testing.T) {
	expected := `[{"z":1,"y":[{"x":true,"w":null},{"v":false,"u":"string"}]},{"t":33,"s":null,"r":42}]`

	var oa Any
	err := json.Unmarshal([]byte(expected), &oa)
	require.NoError(t, err)

	bs, err := json.Marshal(oa)
	require.NoError(t, err)

	// check that the marshalled value is the exact same as the original,
	// including the order of the keys in the maps
	assert.Equal(t, expected, string(bs))
}

func FuzzAnyUnmarshalJSON(f *testing.F) {
	f.Add(`null`)
	f.Add(`2`)
	f.Add(`-1.2e3`)
	f.Add(`false`)
	f.Add(`[]`)
	f.Add(`[1,true,"three"]`)
	f.Add(`{}`)
	f.Add(`{"a":1,"c":2,"b":3}`)
	// incorrect json:
	f.Add(``)
	f.Add(`[`)
	f.Add(`{`)
	f.Add(`[1,2,3,]`)
	f.Add(`{"a":["foo","bar"`)
	f.Fuzz(func(t *testing.T, str string) {
		var goX any
		errGo := json.Unmarshal([]byte(str), &goX)

		var x Any
		err := x.UnmarshalJSON([]byte(str))

		if errGo == nil {
			assert.NoError(t, err, "should not trigger an error: |%s|", str)
		} else {
			assert.Error(t, err, "should trigger an error: |%s|", str)
		}
	})
}

func FuzzAnyJsonDecoder(f *testing.F) {
	f.Add(`null`)
	f.Add(`2`)
	f.Add(`-1.2e3`)
	f.Add(`false`)
	f.Add(`[]`)
	f.Add(`[1,true,"three"]`)
	f.Add(`{}`)
	f.Add(`{"a":1,"c":2,"b":3}`)
	// incorrect json:
	f.Add(``)
	f.Add(`[`)
	f.Add(`{`)
	f.Add(`[1,2,3,]`)
	f.Add(`{"a":["foo","bar"`)
	f.Fuzz(func(t *testing.T, str string) {
		var dec = json.NewDecoder(strings.NewReader(str))
		var goX any
		errGo := dec.Decode(&goX)

		dec = json.NewDecoder(strings.NewReader(str))
		var x Any
		err := dec.Decode(&x)

		if errGo == nil {
			assert.NoError(t, err, "should not trigger an error: |%s|", str)
		} else {
			assert.Error(t, err, "should trigger an error: |%s|", str)
		}
	})
}

func jsonCompact(input string) string {
	var x any
	err := json.Unmarshal([]byte(input), &x)
	if err != nil {
		return ""
	}

	bs, err := json.Marshal(x)
	if err != nil {
		return ""
	}

	return string(bs)
}

func FuzzAnyJsonUnmarshal(f *testing.F) {
	f.Add(`null`)
	f.Add(`2`)
	f.Add(`-1.2e3`)
	f.Add(`false`)
	f.Add(`[]`)
	f.Add(`[1,true,"three"]`)
	f.Add(`{}`)
	f.Add(`{"a":1,"c":2,"b":3}`)

	f.Fuzz(func(t *testing.T, input string) {
		input = jsonCompact(input)
		if input == "" {
			return
		}

		var x Any
		err := json.Unmarshal([]byte(input), &x)
		require.NoError(t, err)

		jsonEnd, err := json.Marshal(x)
		require.NoError(t, err)
		assert.Equal(t, input, string(jsonEnd), "input: |%s|")
	})
}
