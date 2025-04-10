package ordmap_test

import (
	"encoding/json"
	"os"

	"github.com/LeGEC/ordmap"
	"gopkg.in/yaml.v3"
)

func ExampleMap_json() {
	input := `{
		"last_name": "Doe",
		"first_name": "John",
		"age": 42,
		"skills": {
			"go": 5,
			"python": 3,
			"ada": 2,
			"rust": 1
		}
	}`

	var x ordmap.Map[string, any]
	_ = json.Unmarshal([]byte(input), &x)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// ordmap.Map[string, any]: the order of the keys in the root object is preserved,
	//   but with values of type 'any', the order in nested objects is not preserved.
	enc.Encode(x)

	// Output:
	//
	// {
	//   "last_name": "Doe",
	//   "first_name": "John",
	//   "age": 42,
	//   "skills": {
	//     "ada": 2,
	//     "go": 5,
	//     "python": 3,
	//     "rust": 1
	//   }
	// }
}

func ExampleMap_yaml() {
	input := `
last_name: Doe
first_name: John
age: 42
skills:
  go: 5
  python: 3
  ada: 2
  rust: 1`

	var x ordmap.Map[string, any]
	_ = yaml.Unmarshal([]byte(input), &x)

	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)

	// ordmap.Map[string, any]: the order of the keys in the root object is preserved,
	//   but with an 'any' type as value, the order in nested objects is not preserved.
	enc.Encode(x)

	// Output:
	//
	// last_name: Doe
	// first_name: John
	// age: 42
	// skills:
	//   ada: 2
	//   go: 5
	//   python: 3
	//   rust: 1
}
