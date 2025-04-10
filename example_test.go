package ordmap_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeGEC/ordmap"
)

func Example_compared() {
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

	var std any
	var oMap ordmap.Map[string, any]
	var oAny ordmap.Any

	_ = json.Unmarshal([]byte(input), &std)
	_ = json.Unmarshal([]byte(input), &oMap)
	_ = json.Unmarshal([]byte(input), &oAny)

	var output = json.NewEncoder(os.Stdout)
	output.SetIndent("", "  ")
	output.SetEscapeHTML(false)

	fmt.Println("// standard 'any' object: order of keys is not preserved")
	output.Encode(std)
	fmt.Println()

	fmt.Println("// ordmap.Map[string, any]: order of keys is preserved in the root object (last_name, first_name, age, skills),")
	fmt.Println("//     but with values of 'any' type, not in nested objects:")
	fmt.Println("//       inside 'skills' object: the order (go, python, ada, rust) is not preserved")
	output.Encode(oMap)
	fmt.Println()

	fmt.Println("// ordmap.Any: order of keys is preserved everywhere")
	output.Encode(oAny)

	// Output:
	// // standard 'any' object: order of keys is not preserved
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
	//
	// // ordmap.Map[string, any]: order of keys is preserved in the root object (last_name, first_name, age, skills),
	// //     but with values of 'any' type, not in nested objects:
	// //       inside 'skills' object: the order (go, python, ada, rust) is not preserved
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
	//
	// // ordmap.Any: order of keys is preserved everywhere
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
}

func Example_openapiUseCase() {
	// a Schema declaration, to match a 'schema' entry in an OpenAPI document,
	// limited to the fields visible in the example below (excerpt from the github openapi specification)
	type Schema struct {
		Ref   string    `yaml:"$ref,omitempty" json:"$ref,omitempty"`
		AnyOf []*Schema `yaml:"anyOf,omitempty" json:"anyOf,omitempty"`
		OneOf []*Schema `yaml:"oneOf,omitempty" json:"oneOf,omitempty"`

		Title       string  `yaml:"title,omitempty" json:"title,omitempty"`
		Description string  `yaml:"description,omitempty" json:"description,omitempty"`
		Type        any     `yaml:"type,omitempty" json:"type,omitempty"`
		Enum        []any   `yaml:"enum,omitempty" json:"enum,omitempty"`
		Default     any     `yaml:"default,omitempty" json:"default,omitempty"`
		Format      string  `yaml:"format,omitempty" json:"format,omitempty"`
		Items       *Schema `yaml:"items,omitempty" json:"items,omitempty"`

		// The 'example' field may hold any value, be it a simple scalar or a fully fledged object.
		// Using 'ordpmap.Any' here allows to keep the original order if you are into unmarshalling/marshalling
		// a schema, and want to have a stable value here
		Example *ordmap.Any `yaml:"example,omitempty" json:"example,omitempty"`

		Examples []*ordmap.Any `yaml:"examples,omitempty" json:"examples,omitempty"`

		// keep Properties in same order as in the original schema
		Properties *ordmap.Map[string, *Schema] `yaml:"properties,omitempty" json:"properties,omitempty"`

		ReadOnly  bool     `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
		Required  []string `yaml:"required,omitempty" json:"required,omitempty"`
		WriteOnly bool     `yaml:"writeOnly,omitempty" json:"writeOnly,omitempty"`
	}

	type Components struct {
		Schemas *ordmap.Map[string, *Schema] `yaml:"schemas,omitempty" json:"schemas,omitempty"`
		// other fields ...
	}

	type Document struct {
		Components *Components `yaml:"components,omitempty" json:"components,omitempty"`
		// other fields ...
	}

	// file testdata/github.excerpt.3.14.json:
	// excerpt from the github openapi specification at:
	//  https://raw.githubusercontent.com/github/rest-api-description/main/descriptions-next/ghes-3.14/ghes-3.14.2022-11-28.json
	//
	// you will note that the order of the "properties", for example, is preserved,
	// as well as, in the "authentication-token" schema, the "examples" section for "properties.permissions"
	inputFile := "testdata/github.excerpt.3.14.json"
	input, err := os.ReadFile(inputFile)

	var doc Document
	err = json.Unmarshal(input, &doc)
	if err != nil {
		panic(err)
	}

	marshalled, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		panic(err)
	}

	if bytes.Equal(marshalled, input) {
		fmt.Println("fields order from input is preserved")
	} else {
		fmt.Println("!!! fields order from input is not preserved")
	}

	// Output:
	//
	// fields order from input is preserved
}
