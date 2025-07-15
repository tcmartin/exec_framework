package nodes

import (
	"context"
	"go-workflow/pkg/framework"
	"testing"
)

func TestManualTrigger_Execute(t *testing.T) {
	payload := []map[string]interface{}{{"key": "value"}}
	node := &ManualTrigger{Payload: payload}

	ctx := &framework.Context{
		Ctx: context.Background(),
	}

	out, err := node.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 output, got %d", len(out))
	}
	if out[0]["key"] != "value" {
		t.Errorf("expected key value, got %v", out[0]["key"])
	}
}
