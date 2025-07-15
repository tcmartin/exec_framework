package framework

import (
    "time"
)

// Node is the interface for a workflow step
// defined in node.go

// Workflow defines nodes and their connections
type Workflow struct {
    Nodes       map[string]Node
    Connections map[string][]string
}

// Run executes the workflow starting at startNode. Supports branching/loops.
func (w *Workflow) Run(ctx *Context, startNode string) error {
    data := map[string][]map[string]interface{}{}
    var exec func(name string) error
    exec = func(name string) error {
        node := w.Nodes[name]
        inputs := data[name]
        start := time.Now()
        outputs, err := node.Execute(ctx, inputs)
        duration := time.Since(start).Seconds()
        ctx.Metrics.NodeDuration.WithLabelValues(name).Observe(duration)
        if err != nil {
            ctx.Metrics.NodeErrors.WithLabelValues(name).Inc()
            ctx.Logger.Errorf("Node %s error: %v", name, err)
            return err
        }
        for _, child := range w.Connections[name] {
            data[child] = append(data[child], outputs...)
        }
        for _, child := range w.Connections[name] {
            if err := exec(child); err != nil {
                return err
            }
        }
        return nil
    }
    data[startNode] = nil
    return exec(startNode)
}
