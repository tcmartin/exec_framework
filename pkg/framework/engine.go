package framework

import (
    "time"
)

// Workflow manages nodes and their connections
type Workflow struct {
    Nodes           map[string]Node
    Connections     map[string][]string
    ErrorConnections map[string]string // Map from node name to error handler node name
}

// Run executes starting at startNode, supporting branching and loops
func (w *Workflow) Run(ctx *Context, startNode string, initialInput []map[string]interface{}) error {
    data := map[string][]map[string]interface{}{}
    data[startNode] = initialInput
    var exec func(name string) error
    exec = func(name string) error {
        node := w.Nodes[name]
        inputs := data[name]
        start := time.Now()
        outputs, err := node.Execute(ctx, inputs)
        elapsed := time.Since(start).Seconds()
        ctx.Metrics.NodeDuration.WithLabelValues(name).Observe(elapsed)
        if err != nil {
            ctx.Metrics.NodeErrors.WithLabelValues(name).Inc()
            ctx.Logger.Errorf("node %s error: %v", name, err)

            if errorNodeName, ok := w.ErrorConnections[name]; ok {
                // Route error to the specified error handling node
                errorRecord := map[string]interface{}{
                    "original_input": inputs,
                    "error":          err.Error(),
                    "node":           name,
                }
                data[errorNodeName] = append(data[errorNodeName], errorRecord)
                // Continue execution from the error node, but don't pass outputs to regular children
                return exec(errorNodeName)
            } else {
                // No error connection, propagate the error
                return err
            }
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
    return exec(startNode)
}