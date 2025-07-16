package framework

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// NodeFactory is a function that creates a Node instance from a YAML node definition.
type NodeFactory func(nodeDef *yaml.Node) (Node, error)

// nodeFactories stores registered NodeFactory functions.
var nodeFactories = make(map[string]NodeFactory)

// RegisterNodeFactory registers a NodeFactory for a given node type.
func RegisterNodeFactory(nodeType string, factory NodeFactory) {
	nodeFactories[nodeType] = factory
}

// CreateNode creates a Node instance based on the provided YAML node definition.
func CreateNode(nodeDef *yaml.Node) (Node, error) {
	var nodeType struct {
		Type string `yaml:"type"`
	}
	if err := nodeDef.Decode(&nodeType); err != nil {
		return nil, fmt.Errorf("failed to decode node type: %w", err)
	}

	factory, ok := nodeFactories[nodeType.Type]
	if !ok {
		return nil, fmt.Errorf("unknown node type: %s", nodeType.Type)
	}

	return factory(nodeDef)
}
