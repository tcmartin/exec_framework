package nodes

import (
	"fmt"
	"time"

	"go-workflow/pkg/framework"
)

// Sleeper interface for mocking time.Sleep
type Sleeper interface {
	Sleep(d time.Duration)
}

// realSleeper implements Sleeper using time.Sleep
type realSleeper struct{}

func (rs *realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Clock interface for mocking time.Now
type Clock interface {
	Now() time.Time
}

// realClock implements Clock using time.Now
type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

// WaitForNode pauses processing of an individual item until a specified timestamp.
type WaitForNode struct {
	TimestampKey string
	Sleeper      Sleeper
	Clock        Clock
}

// NewWaitForNode creates a new WaitForNode.
func NewWaitForNode(timestampKey string) *WaitForNode {
	return &WaitForNode{TimestampKey: timestampKey, Sleeper: &realSleeper{}, Clock: &realClock{}}
}

// Execute pauses processing for each input record until the timestamp specified by TimestampKey.
func (n *WaitForNode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	var outputs []map[string]interface{}

	for _, input := range inputs {
		timestampVal, ok := input[n.TimestampKey]
		if !ok {
			return nil, fmt.Errorf("timestamp key '%s' not found in input record", n.TimestampKey)
		}

		timestampStr, isString := timestampVal.(string)
		if !isString {
			return nil, fmt.Errorf("timestamp value for key '%s' is not a string", n.TimestampKey)
		}

		t, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp '%s' for key '%s': %w", timestampStr, n.TimestampKey, err)
        }

        duration := t.Sub(n.Clock.Now())
		if duration > 0 {
			n.Sleeper.Sleep(duration)
		}

		outputs = append(outputs, input)
	}

	return outputs, nil
}
