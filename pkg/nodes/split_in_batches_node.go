package nodes

import (
	"fmt"

	"go-workflow/pkg/framework"
)

// SplitInBatchesNode splits a large array of items into chunks of N items.
type SplitInBatchesNode struct {
	BatchSize int
}

// NewSplitInBatchesNode creates a new SplitInBatchesNode.
func NewSplitInBatchesNode(batchSize int) *SplitInBatchesNode {
	return &SplitInBatchesNode{BatchSize: batchSize}
}

// Execute splits the input records into batches.
func (n *SplitInBatchesNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	if n.BatchSize <= 0 {
		return nil, fmt.Errorf("batch size must be greater than 0")
	}

	outputs := make([]map[string]interface{}, 0)
	for i := 0; i < len(inputs); i += n.BatchSize {
		end := i + n.BatchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		// Each batch is a single record in the output, containing a slice of the original records.
		outputs = append(outputs, map[string]interface{}{"batch": inputs[i:end]})
	}

	return outputs, nil
}
