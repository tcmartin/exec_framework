package nodes

import (
	"context"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
)

func TestWebhookTrigger_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}
	node := NewWebhookTrigger()

	inputs := []map[string]interface{}{
		{"key1": "value1"},
		{"key2": "value2"},
	}
	expectedOutputs := inputs // WebhookTrigger should just pass through inputs

	out, err := node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(out, expectedOutputs) {
		t.Errorf("expected %v, got %v", expectedOutputs, out)
	}

	// Test with empty inputs
	inputs = []map[string]interface{}{}
	expectedOutputs = inputs
	out, err = node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(out, expectedOutputs) {
		t.Errorf("expected %v, got %v", expectedOutputs, out)
	}
}
