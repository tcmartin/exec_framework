package nodes

import (
	"go-workflow/pkg/framework"
)

// MergeByKeyNode groups incoming items by a key and outputs one merged record per key.
type MergeByKeyNode struct {
	Key string
}

// NewMergeByKeyNode creates a new MergeByKeyNode.
func NewMergeByKeyNode(key string) *MergeByKeyNode {
	return &MergeByKeyNode{Key: key}
}

// Execute groups and merges input records by the specified key.
// For each unique key, all records sharing that key are merged into a single output record,
// with values from later records overwriting earlier ones in case of conflicts.
func (n *MergeByKeyNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	mergedGroups := make(map[interface{}]map[string]interface{})

	for _, input := range inputs {
		keyValue, ok := input[n.Key]
		if !ok {
			// If the key is missing, this record is not part of any group to be merged.
			// For now, we'll skip it. Depending on requirements, it could be passed through
			// as a separate output or cause an error.
			continue
		}

		currentGroup, found := mergedGroups[keyValue]
		if !found {
			// First record for this key, initialize the group with a copy of the record
			currentGroup = make(map[string]interface{})
			for k, v := range input {
				currentGroup[k] = v
			}
			mergedGroups[keyValue] = currentGroup
		} else {
			// Merge current input into the existing group record
			for k, v := range input {
				currentGroup[k] = v
			}
		}
	}

	// Convert the map of merged groups back to a slice of records
	outputs := make([]map[string]interface{}, 0, len(mergedGroups))
	for _, record := range mergedGroups {
		outputs = append(outputs, record)
	}

	return outputs, nil
}
