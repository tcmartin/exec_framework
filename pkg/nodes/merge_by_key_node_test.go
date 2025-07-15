package nodes

import (
	"context"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
)

func TestMergeByKeyNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("merge records by string key", func(t *testing.T) {
		node := NewMergeByKeyNode("id")
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A", "value": 10},
			{"id": "2", "name": "B", "value": 20},
			{"id": "1", "age": 30, "value": 15}, // Should merge with first, value should be 15
			{"id": "3", "name": "C"},
			{"id": "2", "city": "NY"}, // Should merge with second
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1", "name": "A", "value": 15, "age": 30},
			{"id": "2", "name": "B", "value": 20, "city": "NY"},
			{"id": "3", "name": "C"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Convert to a map for easier comparison (order doesn't matter for merge)
		outputMap := make(map[interface{}]map[string]interface{})
		for _, item := range out {
			outputMap[item["id"]] = item
		}

		for _, expected := range expectedOutputs {
			actual, ok := outputMap[expected["id"]]
			if !ok {
				t.Errorf("expected item with id %v not found in output", expected["id"])
				continue
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("for id %v, expected %v, got %v", expected["id"], expected, actual)
			}
		}
		if len(out) != len(expectedOutputs) {
			t.Errorf("expected %d outputs, got %d", len(expectedOutputs), len(out))
		}
	})

	t.Run("handle missing key", func(t *testing.T) {
		node := NewMergeByKeyNode("nonExistentKey")
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
			{"id": "2", "name": "B"},
		}
		// Records with missing key should be skipped
		expectedOutputs := []map[string]interface{}{}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("empty inputs", func(t *testing.T) {
		node := NewMergeByKeyNode("id")
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
}
