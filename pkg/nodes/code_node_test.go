package nodes

import (
	"context"
	"go-workflow/pkg/framework"
	"testing"
)

func TestCodeNode_Execute(t *testing.T) {
	transformFn := func(inputs []map[string]interface{}) []map[string]interface{} {
		var outputs []map[string]interface{}
		for _, input := range inputs {
			input["transformed"] = true
			outputs = append(outputs, input)
		}
		return outputs
	}

	node := &CodeNode{Fn: transformFn}

	ctx := &framework.Context{
		Ctx: context.Background(),
	}

	inputs := []map[string]interface{}{{"key": "value"}}
	out, err := node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 output, got %d", len(out))
	}
	if out[0]["transformed"] != true {
		t.Errorf("expected transformed to be true, got %v", out[0]["transformed"])
	}
}
