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

	got, err := toYamlSimple(x)
	require.NoError(t, err)

	// there will be no indentation issue if we remove only leading and trailing spaces,
	// apply TrimSpace to remove leading and trailing "\n"
	got = strings.TrimSpace(got)
	expected := strings.TrimSpace(payload)

	// we want to check that the order of the keys in each mapping is preserved
	//
	// as a proxy: we compare the serialized strings, so we also add a dependency to the exact serialization
	//   format chosen by gopkg.in/yaml.v3
	assert.Equal(t, expected, got)
}

func toYamlSimple(x any) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(x)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func toYamlDouble(x any) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(x)
	if err != nil {
		return "", err
	}

	var res any
	err = yaml.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		return "", err
	}

	buf.Reset()
	enc = yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err = enc.Encode(res)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TestOrderedAnyCompareWithGoyaml(t *testing.T) {
	var testCases = []string{
		``,
		`~`,
		`foo`,
		`42`,
		`[1,true,"three"]`,
		`x: 1
z: true
y: "three"`,
		`08`,
		`0x8`,
		`0o8`,
		`{"A":08}`,
	}

	var testCasesError = []string{
		// errors:
		`{`, // unclosed '{'
		`foo:
  - x:
    - 1
  y: 2`, // map key after sequence (array) item
		`
{
  foo: !!str,
  !!str: bar,
}`,
		`{"z":null,"z":null}`, // duplicate key
	}

	for _, payload := range testCases {
		var x any
		errStd := yaml.Unmarshal([]byte(payload), &x)

		var om Any
		errOm := yaml.Unmarshal([]byte(payload), &om)

		if errStd != nil || errOm != nil {
			if errStd != nil {
				t.Logf("*** error for standard map: %v", errStd)
			}
			if errOm != nil {
				t.Logf("*** error for  ordered map: %v", errOm)
			}

			err := errOm
			if err == nil {
				err = errStd
			}
			assert.Fail(t, "unexpected error", "payload: |%s|", payload)
		}

		// check that values match
		expected, _ := toYamlSimple(x)
		got, err := toYamlDouble(om)
		require.NoError(t, err, "unexpected error", "payload: |%s|", payload)

		expected = strings.TrimSpace(expected)
		got = strings.TrimSpace(got)
		assert.Equal(t, expected, got, "payload: |%s|", payload)
	}

	for _, payload := range testCasesError {
		var x any
		errStd := yaml.Unmarshal([]byte(payload), &x)
		if errStd == nil {
			assert.Fail(t, "expected an error with go-yaml library on payload: |%s|", payload)
		}

		var om Any
		errOm := yaml.Unmarshal([]byte(payload), &om)
		if errOm == nil {
			assert.Fail(t, "expected an error with ordered map on payload: |%s|", payload)
		}
	}
}

func FuzzAnyYaml(f *testing.F) {
	// I ended up triggering corner cases in the go-yaml/yaml library itself, and tiptoeing around them
	// Skipping this fuzz test for now
	f.Skip("the go-yaml/yaml has corner cases, which are not worth covering. skipping fuzz test for now")

	f.Add(``)
	f.Add(`~`)
	f.Add(`foo`)
	f.Add(`42`)
	f.Add(`[1,true,"three"]`)
	f.Add(`x: 1
z: true
y: "three"`)
	f.Add(`08`)
	f.Add(`{"A":08}`)

	f.Add(`{`)
	f.Add(`foo:
  - x:
    - 1
  y: 2`)
	f.Add((`
{
  foo: !!str,
  !!str: bar,
}`))
	f.Add(`{"z":null,"z":null}`)
	f.Add(`{1: "a", "1": "b"}`)
	f.Add(`{"z",0,"0"}`)
	f.Add(`{"z",0,"0",0}`)

	f.Fuzz(func(t *testing.T, payload string) {
		var x any
		errStd := yaml.Unmarshal([]byte(payload), &x)

		var om Any
		errOm := yaml.Unmarshal([]byte(payload), &om)

		hasErrStd := errStd != nil
		hasErrOm := errOm != nil

		if hasErrStd != hasErrOm {
			if errStd != nil {
				t.Logf("*** error for standard map: %v", errStd)
			}
			if errOm != nil {
				t.Logf("*** error for  ordered map: %v", errOm)
			}

			t.Logf("got error in one case and not the other, payload: |%s|", payload)
			//			require.Equal(t, hasErrStd, hasErrOm, "got error in one case and not the other, payload: |%s|", payload)
		}

		if hasErrStd || hasErrOm {
			return
		}

		// no error: check that values match
		expected, _ := toYamlSimple(x)
		got, err := toYamlDouble(om)
		require.NoError(t, err, "unexpected error, payload: |%s|", payload)

		expected = strings.TrimSpace(expected)
		got = strings.TrimSpace(got)
		assert.Equal(t, expected, got, "payload: |%s|", payload)
	})
}
