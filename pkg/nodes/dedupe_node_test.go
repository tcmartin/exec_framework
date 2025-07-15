package nodes

import (
	"context"
	"testing"

	"go-workflow/pkg/framework"
)

func TestDedupeNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("deduplicate by string key", func(t *testing.T) {
		node := NewDedupeNode("id")
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
			{"id": "2", "name": "B"},
			{"id": "1", "name": "C"},
			{"id": "3", "name": "D"},
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
			{"id": "2", "name": "B"},
			{"id": "3", "name": "D"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(out) != len(expectedOutputs) {
			t.Fatalf("expected %d outputs, got %d", len(expectedOutputs), len(out))
		}

		// Convert to a map for easier comparison (order doesn't matter for dedupe)
		outputMap := make(map[interface{}]map[string]interface{})
		for _, item := range out {
			outputMap[item["id"]] = item
		}

		for _, expected := range expectedOutputs {
			if _, ok := outputMap[expected["id"]]; !ok {
				t.Errorf("expected item with id %v not found in output", expected["id"])
			}
		}
	})

	t.Run("deduplicate by int key", func(t *testing.T) {
		node := NewDedupeNode("value")
		inputs := []map[string]interface{}{
			{"value": 10, "data": "X"},
			{"value": 20, "data": "Y"},
			{"value": 10, "data": "Z"},
		}
		expectedOutputs := []map[string]interface{}{
			{"value": 10, "data": "X"},
			{"value": 20, "data": "Y"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(out) != len(expectedOutputs) {
			t.Fatalf("expected %d outputs, got %d", len(expectedOutputs), len(out))
		}

		outputMap := make(map[interface{}]map[string]interface{})
		for _, item := range out {
			outputMap[item["value"]] = item
		}

		for _, expected := range expectedOutputs {
			if _, ok := outputMap[expected["value"]]; !ok {
				t.Errorf("expected item with value %v not found in output", expected["value"])
			}
		}
	})

	t.Run("handle missing key", func(t *testing.T) {
		node := NewDedupeNode("nonExistentKey")
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
			{"id": "2", "name": "B"},
		}
		// All inputs should be returned as unique if the key is missing
		expectedOutputs := inputs

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(out) != len(expectedOutputs) {
			t.Fatalf("expected %d outputs, got %d", len(expectedOutputs), len(out))
		}

		// For missing key, order matters as we append directly
		for i, expected := range expectedOutputs {
			if expected["id"] != out[i]["id"] {
				t.Errorf("expected %v, got %v", expected, out[i])
			}
		}
	})
}
