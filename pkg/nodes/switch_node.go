package nodes

import (
	"go-workflow/pkg/framework"
)

// Condition defines a condition for routing.
type Condition struct {
	Field string
	Value interface{}
}

// SwitchNode routes items into different branches based on conditions.
type SwitchNode struct {
	Conditions map[string]Condition // Map of output branch name to condition
}

// NewSwitchNode creates a new SwitchNode.
func NewSwitchNode(conditions map[string]Condition) *SwitchNode {
	return &SwitchNode{Conditions: conditions}
}

// Execute routes the input records based on the defined conditions.
// The output will be a map where keys are branch names and values are slices of records
// that satisfy the condition for that branch.
func (n *SwitchNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	branchOutputs := make(map[string][]map[string]interface{})

	for _, input := range inputs {
		routed := false
		for branchName, cond := range n.Conditions {
			if val, ok := input[cond.Field]; ok && val == cond.Value {
				branchOutputs[branchName] = append(branchOutputs[branchName], input)
				routed = true
				break // Route to the first matching branch
			}
		}
		if !routed {
			// If no condition matches, add to a default 'unmatched' branch or similar
			// For now, we'll just not include it in any branch output.
			// This behavior might need to be explicitly defined in the workflow.
		}
	}

	// Convert the branchOutputs map into a slice of maps, where each map represents a branch.
	// This is a temporary representation and might need adjustment based on engine's branching.
	outputs := make([]map[string]interface{}, 0)
	for branchName, records := range branchOutputs {
		outputs = append(outputs, map[string]interface{}{branchName: records})
	}

	return outputs, nil
}
