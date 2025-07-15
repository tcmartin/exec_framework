package nodes

import (
    "framework"
    "time"
    "math/rand"
)

// WaitNode sleeps for up to MaxSeconds
type WaitNode struct{ MaxSeconds int }

func (n *WaitNode) Execute(ctx *framework.Context, input []map[string]interface{}) ([]map[string]interface{}, error) {
    time.Sleep(time.Duration(rand.Intn(n.MaxSeconds)) * time.Second)
    return input, nil
}
