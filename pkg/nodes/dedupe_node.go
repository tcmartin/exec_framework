package nodes

import (
	"go-workflow/pkg/framework"
)

// DedupeNode deduplicates records based on a specified key.
type DedupeNode struct {
	Key string
}

// NewDedupeNode creates a new DedupeNode.
func NewDedupeNode(key string) *DedupeNode {
	return &DedupeNode{Key: key}
}

// Execute deduplicates the input records.
func (n *DedupeNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	seen := make(map[interface{}]bool)
	var outputs []map[string]interface{}

	for _, input := range inputs {
		value, ok := input[n.Key]
		if !ok {
			// If the key is missing, treat this record as unique or handle as an error.
			// For now, we'll treat it as unique.
			outputs = append(outputs, input)
			continue
		}

		if !seen[value] {
			seen[value] = true
			outputs = append(outputs, input)
		}
	}

	return outputs, nil
}
