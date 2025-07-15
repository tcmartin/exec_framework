package framework

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

// mockNode is a mock implementation of the Node interface for testing.
type mockNode struct {
	name     string
	execute  func(ctx *Context, inputs []map[string]interface{}) ([]map[string]interface{}, error)
	executed bool
}

func (n *mockNode) Execute(ctx *Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	n.executed = true
	if n.execute != nil {
		return n.execute(ctx, inputs)
	}
	return inputs, nil
}

func TestWorkflow_Run(t *testing.T) {
	// Create a new registry for each test to avoid duplicate metric registration.
	metrics := NewMetrics()

	ctx := &Context{
		Ctx:     context.Background(),
		Logger:  zap.NewNop().Sugar(),
		Metrics: metrics,
	}

	t.Run("linear workflow", func(t *testing.T) {
		node1 := &mockNode{name: "node1"}
		node2 := &mockNode{name: "node2"}
		node3 := &mockNode{name: "node3"}

		workflow := &Workflow{
			Nodes: map[string]Node{
				"node1": node1,
				"node2": node2,
				"node3": node3,
			},
			Connections: map[string][]string{
				"node1": {"node2"},
				"node2": {"node3"},
			},
		}

		err := workflow.Run(ctx, "node1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !node1.executed {
			t.Error("node1 was not executed")
		}
		if !node2.executed {
			t.Error("node2 was not executed")
		}
		if !node3.executed {
			t.Error("node3 was not executed")
		}
	})

	t.Run("workflow with error", func(t *testing.T) {
		node1 := &mockNode{name: "node1"}
		errorNode := &mockNode{
			name: "errorNode",
			execute: func(ctx *Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
				return nil, errors.New("test error")
			},
		}
		node3 := &mockNode{name: "node3"}

		workflow := &Workflow{
			Nodes: map[string]Node{
				"node1":     node1,
				"errorNode": errorNode,
				"node3":     node3,
			},
			Connections: map[string][]string{
				"node1": {"errorNode"},
				"errorNode": {"node3"},
			},
		}

		err := workflow.Run(ctx, "node1")
		if err == nil {
			t.Fatal("expected an error, but got nil")
		}

		if !node1.executed {
			t.Error("node1 was not executed")
		}
		if !errorNode.executed {
			t.Error("errorNode was not executed")
		}
		if node3.executed {
			t.Error("node3 should not have been executed")
		}
	})
}

func TestWorkflow_Run_DataPassing(t *testing.T) {
	metrics := NewMetrics()

	ctx := &Context{
		Ctx:     context.Background(),
		Logger:  zap.NewNop().Sugar(),
		Metrics: metrics,
	}

	node1 := &mockNode{
		name: "node1",
		execute: func(ctx *Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
			return []map[string]interface{}{{"key": "value1"}}, nil
		},
	}
	node2 := &mockNode{
		name: "node2",
		execute: func(ctx *Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
			if len(inputs) != 1 {
				t.Errorf("expected 1 input, got %d", len(inputs))
			}
			if inputs[0]["key"] != "value1" {
				t.Errorf("expected input key 'value1', got %v", inputs[0]["key"])
			}
			return []map[string]interface{}{{"key": "value2"}}, nil
		},
	}

	workflow := &Workflow{
		Nodes: map[string]Node{
			"node1": node1,
			"node2": node2,
		},
		Connections: map[string][]string{
			"node1": {"node2"},
		},
	}

	err := workflow.Run(ctx, "node1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !node1.executed {
		t.Error("node1 was not executed")
	}
	if !node2.executed {
		t.Error("node2 was not executed")
	}
}