package nodes

import (
    "go-workflow/pkg/framework"
    "math/rand"
    "time"
)

// WaitNode waits up to MaxSeconds randomly
type WaitNode struct{ MaxSeconds int }

func (n *WaitNode) Execute(ctx *framework.Context, input []map[string]interface{}) ([]map[string]interface{}, error) {
    time.Sleep(time.Duration(rand.Intn(n.MaxSeconds)) * time.Second)
    return input, nil
}