package nodes

import "go-workflow/pkg/framework"

// ManualTrigger seeds initial data
type ManualTrigger struct{ Payload []map[string]interface{} }

func (n *ManualTrigger) Execute(ctx *framework.Context, input []map[string]interface{}) ([]map[string]interface{}, error) {
    return n.Payload, nil
}