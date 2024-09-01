package ordmap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderedAny_UnmarshalJSON(t *testing.T) {
	jsonInput := `{"key1":"value1","key3":true,"key2":2}`

	var oa Any

	err := json.Unmarshal([]byte(jsonInput), &oa)
	require.NoError(t, err)

	// validation:
	require.IsType(t, (*Map[string, any])(nil), oa.v)

	v := oa.v.(*Map[string, any])
	require.Equal(t, 3, v.Len())

	assert.Equal(t, "value1", v.Get("key1"))
	assert.EqualValues(t, 2, v.Get("key2"))
	assert.Equal(t, true, v.Get("key3"))

	assert.Equal(t, "key1", v.keys[0])
	assert.Equal(t, "key3", v.keys[1])
	assert.Equal(t, "key2", v.keys[2])
}

func TestOrderedAny_MarshalJSON(t *testing.T) {
	jsonInput := `{"key1":"value1","key3":true,"key2":2}`

	var oa Any
	err := json.Unmarshal([]byte(jsonInput), &oa)
	require.NoError(t, err)

	marshalled, err := json.Marshal(oa)
	require.NoError(t, err)

	// check that the marshalled value is the exact same as the original,
	// including the order of the keys in the maps
	assert.Equal(t, jsonInput, string(marshalled))
}

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
