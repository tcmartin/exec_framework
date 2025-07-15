package nodes

import "framework"

// ManualTrigger seeds the flow
type ManualTrigger struct {
    Payload []map[string]interface{}
}

func (n *ManualTrigger) Execute(ctx *framework.Context, input []map[string]interface{}) ([]map[string]interface{}, error) {
    return n.Payload, nil
}
