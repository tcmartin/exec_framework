package nodes

import (
	"go-workflow/pkg/framework"
)

// SetNode adds, updates, or removes fields on each item.
type SetNode struct {
	SetValues  map[string]interface{}
	RemoveKeys []string
}

// NewSetNode creates a new SetNode.
func NewSetNode(setValues map[string]interface{}, removeKeys []string) *SetNode {
	return &SetNode{
		SetValues:  setValues,
		RemoveKeys: removeKeys,
	}
}

// Execute applies the set/remove operations to the input records.
func (n *SetNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	outputs := make([]map[string]interface{}, 0)
	for _, input := range inputs {
		// Create a copy to avoid modifying the original input map directly if it's reused elsewhere
		// (though in this framework, inputs are typically new for each node execution)
		output := make(map[string]interface{})
		for k, v := range input {
			output[k] = v
		}

		// Apply SetValues
		for k, v := range n.SetValues {
			output[k] = v
		}

		// Apply RemoveKeys
		for _, k := range n.RemoveKeys {
			delete(output, k)
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}
