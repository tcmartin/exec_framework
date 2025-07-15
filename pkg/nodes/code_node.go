package nodes

import "go-workflow/pkg/framework"

// CodeNode applies a transform function
type CodeNode struct{ Fn func([]map[string]interface{}) []map[string]interface{} }

func (n *CodeNode) Execute(ctx *framework.Context, input []map[string]interface{}) ([]map[string]interface{}, error) {
    return n.Fn(input), nil
}