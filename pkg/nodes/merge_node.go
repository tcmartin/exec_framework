package nodes

import (
	"go-workflow/pkg/framework"
)

// MergeNode merges records based on a specified key.
type MergeNode struct {
	Key string
}

// NewMergeNode creates a new MergeNode.
func NewMergeNode(key string) *MergeNode {
	return &MergeNode{Key: key}
}

// Execute merges the input records based on the specified key.
// If multiple records have the same key, their fields are merged,
// with values from later records overwriting earlier ones in case of conflicts.
func (n *MergeNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	mergedRecords := make(map[interface{}]map[string]interface{})
	var uniqueRecords []map[string]interface{}

	for _, input := range inputs {
		keyValue, ok := input[n.Key]
		if !ok {
			// If the key is missing, treat this record as unique and add it directly.
			uniqueRecords = append(uniqueRecords, input)
			continue
		}

		if existingRecord, found := mergedRecords[keyValue]; found {
			// Merge current input into existing record
			for k, v := range input {
				existingRecord[k] = v
			}
		} else {
			// Add a copy of the current input as a new merged record
			newRecord := make(map[string]interface{})
			for k, v := range input {
				newRecord[k] = v
			}
			mergedRecords[keyValue] = newRecord
		}
	}

	// Convert the map of merged records back to a slice
	outputs := make([]map[string]interface{}, 0, len(mergedRecords)+len(uniqueRecords))
	for _, record := range mergedRecords {
		outputs = append(outputs, record)
	}
	outputs = append(outputs, uniqueRecords...)

	return outputs, nil
}
