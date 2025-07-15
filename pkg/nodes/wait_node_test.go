package nodes

import (
	"context"
	"go-workflow/pkg/framework"
	"testing"
	"time"
)

func TestWaitNode_Execute(t *testing.T) {
	node := &WaitNode{MaxSeconds: 1}

	ctx := &framework.Context{
		Ctx: context.Background(),
	}

	start := time.Now()
	input := []map[string]interface{}{{"key": "value"}}
	out, err := node.Execute(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration := time.Since(start)

	if duration > time.Second {
		t.Errorf("expected wait time to be less than 1 second, got %v", duration)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 output, got %d", len(out))
	}
	if out[0]["key"] != "value" {
		t.Errorf("expected key value, got %v", out[0]["key"])
	}
}
