package ordmap_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeGEC/ordmap"
)

func ExampleStandard() {
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

	var x any
	_ = json.Unmarshal([]byte(input), &x)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(x)

	// regular go unmarshaling: the order of the keys is not preserved

	// Output:
	// {
	//   "age": 42,
	//   "first_name": "John",
	//   "last_name": "Doe",
	//   "skills": {
	//     "ada": 2,
	//     "go": 5,
	//     "python": 3,
	//     "rust": 1
	//   }
	// }
}

func ExampleMap() {
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
	enc.Encode(x)

	// OrderedMap: the order of the keys in the root object is preserved,
	//   but with an 'any' type as value, the order in nested objects is not preserved.

	// Output:
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

func ExampleAny() {
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
