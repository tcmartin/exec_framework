package nodes

import (
	"context"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
)

func TestSwitchNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("route by single condition", func(t *testing.T) {
		conditions := map[string]Condition{
			"branchA": {Field: "status", Value: "pending"},
		}
		node := NewSwitchNode(conditions)

		inputs := []map[string]interface{}{
			{"id": 1, "status": "pending"},
			{"id": 2, "status": "approved"},
			{"id": 3, "status": "pending"},
		}
		expectedOutputs := []map[string]interface{}{
			{"branchA": []map[string]interface{}{
				{"id": 1, "status": "pending"},
				{"id": 3, "status": "pending"},
			}},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Order of branches in the output slice might vary, so check content flexibly
		if len(out) != len(expectedOutputs) {
			t.Fatalf("expected %d outputs, got %d", len(expectedOutputs), len(out))
		}
		// This comparison is simplified and assumes only one branch in this test case
		if !reflect.DeepEqual(out[0]["branchA"], expectedOutputs[0]["branchA"]) {
			t.Errorf("expected branchA %v, got %v", expectedOutputs[0]["branchA"], out[0]["branchA"])
		}
	})

	t.Run("route by multiple conditions", func(t *testing.T) {
		conditions := map[string]Condition{
			"branchA": {Field: "type", Value: "email"},
			"branchB": {Field: "type", Value: "sms"},
		}
		node := NewSwitchNode(conditions)

		inputs := []map[string]interface{}{
			{"id": 1, "type": "email"},
			{"id": 2, "type": "sms"},
			{"id": 3, "type": "push"},
			{"id": 4, "type": "email"},
		}
		expectedBranchA := []map[string]interface{}{
			{"id": 1, "type": "email"},
			{"id": 4, "type": "email"},
		}
		expectedBranchB := []map[string]interface{}{
			{"id": 2, "type": "sms"},
		}

		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if the number of output branches is correct
		if len(out) != 2 {
			t.Fatalf("expected 2 output branches, got %d", len(out))
		}

		// Extract actual branch outputs into maps for easier comparison
		actualBranches := make(map[string][]map[string]interface{})
		for _, branchMap := range out {
			for k, v := range branchMap {
				actualBranches[k] = v.([]map[string]interface{})
			}
		}

		if !reflect.DeepEqual(actualBranches["branchA"], expectedBranchA) {
			t.Errorf("expected branchA %v, got %v", expectedBranchA, actualBranches["branchA"])
		}
		if !reflect.DeepEqual(actualBranches["branchB"], expectedBranchB) {
			t.Errorf("expected branchB %v, got %v", expectedBranchB, actualBranches["branchB"])
		}
	})

	t.Run("no matching condition", func(t *testing.T) {
		conditions := map[string]Condition{
			"branchA": {Field: "status", Value: "pending"},
		}
		node := NewSwitchNode(conditions)

		inputs := []map[string]interface{}{
			{"id": 1, "status": "approved"},
		}
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
		conditions := map[string]Condition{
			"branchA": {Field: "status", Value: "pending"},
		}
		node := NewSwitchNode(conditions)

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
