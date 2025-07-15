package nodes

import (
	"context"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
)

func TestSetNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("set new values", func(t *testing.T) {
		node := NewSetNode(map[string]interface{}{"newField": "newValue"}, nil)
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1", "name": "A", "newField": "newValue"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("update existing values", func(t *testing.T) {
		node := NewSetNode(map[string]interface{}{"name": "UpdatedA"}, nil)
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1", "name": "UpdatedA"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("remove keys", func(t *testing.T) {
		node := NewSetNode(nil, []string{"name"})
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A"},
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(out, expectedOutputs) {
			t.Errorf("expected %v, got %v", expectedOutputs, out)
		}
	})

	t.Run("set and remove keys", func(t *testing.T) {
		node := NewSetNode(map[string]interface{}{"newField": "newValue"}, []string{"name"})
		inputs := []map[string]interface{}{
			{"id": "1", "name": "A", "oldField": "oldValue"},
		}
		expectedOutputs := []map[string]interface{}{
			{"id": "1", "oldField": "oldValue", "newField": "newValue"},
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
		node := NewSetNode(map[string]interface{}{"newField": "newValue"}, nil)
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
