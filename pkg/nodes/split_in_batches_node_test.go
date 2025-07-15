package nodes

import (
	"context"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
)

func TestSplitInBatchesNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("split into batches of 2", func(t *testing.T) {
		node := NewSplitInBatchesNode(2)
		inputs := []map[string]interface{}{
			{"id": 1},
			{"id": 2},
			{"id": 3},
			{"id": 4},
			{"id": 5},
		}
		expectedOutputs := []map[string]interface{}{
			{"batch": []map[string]interface{}{{"id": 1}, {"id": 2}}},
			{"batch": []map[string]interface{}{{"id": 3}, {"id": 4}}},
			{"batch": []map[string]interface{}{{"id": 5}}},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("empty inputs", func(t *testing.T) {
		node := NewSplitInBatchesNode(2)
		inputs := []map[string]interface{}{}
		expectedOutputs := []map[string]interface{}{}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("batch size larger than inputs", func(t *testing.T) {
		node := NewSplitInBatchesNode(10)
		inputs := []map[string]interface{}{
			{"id": 1},
			{"id": 2},
		}
		expectedOutputs := []map[string]interface{}{
			{"batch": []map[string]interface{}{{"id": 1}, {"id": 2}}},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("batch size of 1", func(t *testing.T) {
		node := NewSplitInBatchesNode(1)
		inputs := []map[string]interface{}{
			{"id": 1},
			{"id": 2},
		}
		expectedOutputs := []map[string]interface{}{
			{"batch": []map[string]interface{}{{"id": 1}}},
			{"batch": []map[string]interface{}{{"id": 2}}},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("batch size of 0 should return error", func(t *testing.T) {
		node := NewSplitInBatchesNode(0)
		inputs := []map[string]interface{}{{"id": 1}}
		_, err := node.Execute(ctx, inputs)
		if err == nil || err.Error() != "batch size must be greater than 0" {
			t.Errorf("expected error for batch size 0, got %v", err)
		}
	})
}
