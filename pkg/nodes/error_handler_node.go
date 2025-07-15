package nodes

import (
	"go-workflow/pkg/framework"
)

// ErrorHandlerNode is a node that receives error records.
type ErrorHandlerNode struct {
	// For testing purposes, we can store the received errors.
	ReceivedErrors []map[string]interface{}
}

// NewErrorHandlerNode creates a new ErrorHandlerNode.
func NewErrorHandlerNode() *ErrorHandlerNode {
	return &ErrorHandlerNode{}
}

// Execute processes the incoming error records.
func (n *ErrorHandlerNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	// In a real scenario, this node would log, send alerts, retry, etc.
	// For testing, we'll just store them.
	n.ReceivedErrors = append(n.ReceivedErrors, inputs...)
	return nil, nil // Error handler nodes typically don't produce further output for the main flow
}
