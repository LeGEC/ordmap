package ordmap_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeGEC/ordmap"
	"gopkg.in/yaml.v3"
)

func ExampleAny_json() {
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

	var x ordmap.Any
	_ = json.Unmarshal([]byte(input), &x)

	// ordmap.Any: all objects, including nested, are unmarshalled as ordmap.Map,
	//   the order of keys is preserved everywhere
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(x)

	fmt.Println()

	obj := x.V().(*ordmap.Map[string, any])
	// the type for any json object is '*ordmap.Map[string, any]':
	fmt.Printf("type for skills: %T\n", obj.Get("skills"))
	// the type for other fields is the regular go type for generic unmarshalling
	fmt.Printf("type for first_name: %T\n", obj.Get("first_name"))

	// Output:
	// {
	//   "last_name": "Doe",
	//   "first_name": "John",
	//   "age": 42,
	//   "skills": {
	//     "go": 5,
	//     "python": 3,
	//     "ada": 2,
	//     "rust": 1
	//   }
	// }
	//
	// type for skills: *ordmap.Map[string,interface {}]
	// type for first_name: string
}

func ExampleAny_yaml() {
	input := `
last_name: Doe
first_name: John
age: 42
skills:
  go: 5
  python: 3
  ada: 2
  rust: 1`

	var x ordmap.Any
	_ = yaml.Unmarshal([]byte(input), &x)

	// ordmap.Any: all objects, including nested, are unmarshalled as ordmap.Map,
	//   the order of keys is preserved everywhere
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	enc.Encode(x)

	// Output:
	//
	// last_name: Doe
	// first_name: John
	// age: 42
	// skills:
	//   go: 5
	//   python: 3
	//   ada: 2
	//   rust: 1
}
