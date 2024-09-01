package ordmap

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func strYamlKind(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.MappingNode:
		return "mapping"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"

	default:
		return fmt.Sprintf("<unknown %d>", kind)
	}
}

func (m *Map[K, V]) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("invalid yaml value: expected a mapping, got a %s", strYamlKind(value.Kind))
	}

	myMap := make(map[K]V)
	var myKeys []K
	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valueNode := value.Content[i+1]

		var key K
		if err := keyNode.Decode(&key); err != nil {
			return fmt.Errorf("failed to decode key at index %d: %w", i, err)
		}

		var val V
		if err := valueNode.Decode(&val); err != nil {
			return fmt.Errorf("failed to decode value at index %d: %w", i+1, err)
		}

		myMap[key] = val
		myKeys = append(myKeys, key)
	}

	m.m = myMap
	m.keys = myKeys
	return nil
}

func (m Map[K, V]) MarshalYAML() (any, error) {
	if len(m.keys) == 0 {
		return struct{}{}, nil
	}

	node := &yaml.Node{}
	node.Kind = yaml.MappingNode
	node.Content = make([]*yaml.Node, 0, len(m.keys)*2)

	for _, key := range m.keys {
		var keyNode yaml.Node
		var valueNode yaml.Node

		keyBytes, err := yaml.Marshal(key)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key: %w", err)
		}
		err = yaml.Unmarshal(keyBytes, &keyNode)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key: %w", err)
		}
		if keyNode.Kind == yaml.DocumentNode {
			keyNode = *keyNode.Content[0]
		}

		valueBytes, err := yaml.Marshal(m.m[key])
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value: %w", err)
		}
		err = yaml.Unmarshal(valueBytes, &valueNode)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value: %w", err)
		}
		if valueNode.Kind == yaml.DocumentNode {
			valueNode = *valueNode.Content[0]
		}

		node.Content = append(node.Content, &keyNode, &valueNode)
	}

	return node, nil
}
