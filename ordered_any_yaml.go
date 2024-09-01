package ordmap

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

func (x Any) MarshalYAML() (interface{}, error) {
	return x.v, nil
}

func (x *Any) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.DocumentNode {
		node = node.Content[0]
	}

	var v any
	var err error
	switch node.Kind {
	case yaml.DocumentNode:
		err = errors.New("unexpected document node")
	case yaml.AliasNode:
		err = errors.New("unexpected alias node")

	case yaml.MappingNode:
		v, err = yamlUnmarshalAnyObject(node)
	case yaml.SequenceNode:
		v, err = yamlUnmarshalAnyArray(node)
	default:
		err = node.Decode(&v)
	}
	if err != nil {
		return err
	}

	x.v = v
	return nil
}

func yamlUnmarshalAnyObject(node *yaml.Node) (*Map[any, any], error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("error when decoding object: expected mapping node, got %v", node.Kind)
	}

	var m Map[any, any]
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("error when decoding object key: expected scalar node, got %v", strYamlKind(keyNode.Kind))
		}

		var (
			key   any
			value any
			err   error
		)
		err = keyNode.Decode(&key)
		if err != nil {
			return nil, fmt.Errorf("error when decoding object key: %w", err)
		}

		switch valueNode.Kind {
		case yaml.DocumentNode:
			return nil, fmt.Errorf("error when decoding object value: expected a value node, got %v", strYamlKind(keyNode.Kind))
		case yaml.AliasNode:
			return nil, fmt.Errorf("error when decoding object value: unhadled alias node")

		case yaml.MappingNode:
			value, err = yamlUnmarshalAnyObject(valueNode)
		case yaml.SequenceNode:
			value, err = yamlUnmarshalAnyArray(valueNode)

		case yaml.ScalarNode:
			err = valueNode.Decode(&value)
		}

		if err != nil {
			return nil, err
		}

		m.Set(key, value)
	}
	return &m, nil
}

func yamlUnmarshalAnyArray(node *yaml.Node) ([]any, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("error when decoding array: expected sequence node, got %v", node.Kind)
	}

	var res []any
	for _, valueNode := range node.Content {
		var value any
		var err error
		switch valueNode.Kind {
		case yaml.DocumentNode:
			return nil, fmt.Errorf("error when decoding array value: expected a value node, got %v", strYamlKind(valueNode.Kind))
		case yaml.AliasNode:
			return nil, fmt.Errorf("error when decoding array value: unhandled alias node")

		case yaml.MappingNode:
			value, err = yamlUnmarshalAnyObject(valueNode)
		case yaml.SequenceNode:
			value, err = yamlUnmarshalAnyArray(valueNode)

		case yaml.ScalarNode:
			err = valueNode.Decode(&value)
		}

		if err != nil {
			return nil, err
		}

		res = append(res, value)
	}

	return res, nil
}
