package nodes

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go-workflow/pkg/framework"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// mockErrorNode is a mock node that can be configured to return an error.
type mockErrorNode struct {
	shouldError bool
	name        string
}

func (m *mockErrorNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	if m.shouldError {
		return nil, errors.New("simulated error from " + m.name)
	}
	return inputs, nil
}

func TestErrorHandlerNode(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := framework.NewMetrics(reg)

	ctx := &framework.Context{
		Ctx:     context.Background(),
		Logger:  zap.NewNop().Sugar(),
		Metrics: metrics,
	}

	t.Run("error is routed to error handler", func(t *testing.T) {
		errorNode := &mockErrorNode{shouldError: true, name: "errorProducer"}
		handlerNode := NewErrorHandlerNode()

		workflow := &framework.Workflow{
			Nodes: map[string]framework.Node{
				"errorProducer": errorNode,
				"errorHandler":  handlerNode,
			},
			Connections: map[string][]string{
				"errorProducer": {},
			},
			ErrorConnections: map[string]string{
				"errorProducer": "errorHandler",
			},
		}

		input := []map[string]interface{}{{"data": "test"}}
		// Manually set initial data for the start node
		workflow.Nodes["errorProducer"].(*mockErrorNode).Execute(ctx, input) // Simulate initial input

		err := workflow.Run(ctx, "errorProducer", input)
		if err != nil {
			// The workflow should not return an error if it's handled by an error node
			t.Fatalf("expected no error from workflow.Run, got %v", err)
		}

		if len(handlerNode.ReceivedErrors) != 1 {
			t.Fatalf("expected 1 error record, got %d", len(handlerNode.ReceivedErrors))
		}

		expectedErrorRecord := map[string]interface{}{
			"original_input": input,
			"error":          "simulated error from errorProducer",
			"node":           "errorProducer",
		}

		// Compare received error record, ignoring order of map keys
		if !reflect.DeepEqual(handlerNode.ReceivedErrors[0], expectedErrorRecord) {
			t.Errorf("expected error record %v, got %v", expectedErrorRecord, handlerNode.ReceivedErrors[0])
		}
	})

	t.Run("error is not routed if no error connection", func(t *testing.T) {
		errorNode := &mockErrorNode{shouldError: true, name: "errorProducer"}
		handlerNode := NewErrorHandlerNode()

		workflow := &framework.Workflow{
			Nodes: map[string]framework.Node{
				"errorProducer": errorNode,
				"errorHandler":  handlerNode,
			},
			Connections: map[string][]string{
				"errorProducer": {},
			},
			ErrorConnections: map[string]string{},
		}

		input := []map[string]interface{}{{"data": "test"}}
		workflow.Nodes["errorProducer"].(*mockErrorNode).Execute(ctx, input) // Simulate initial input

		err := workflow.Run(ctx, "errorProducer", input)
		if err == nil {
			t.Fatal("expected an error from workflow.Run, but got nil")
		}
		if err.Error() != "simulated error from errorProducer" {
			t.Errorf("expected error 'simulated error from errorProducer', got %v", err)
		}

		if len(handlerNode.ReceivedErrors) != 0 {
			t.Errorf("expected 0 error records, got %d", len(handlerNode.ReceivedErrors))
		}
	})
}
