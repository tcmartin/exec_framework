package nodes

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-workflow/pkg/framework"
)

// mockSleeper implements Sleeper for testing.
type mockSleeper struct {
	sleptFor time.Duration
}

func (ms *mockSleeper) Sleep(d time.Duration) {
	ms.sleptFor = d
}

// mockClock implements Clock for testing.
type mockClock struct {
	now time.Time
}

func (mc *mockClock) Now() time.Time {
	return mc.now
}

func TestWaitForNode_Execute(t *testing.T) {
	ctx := &framework.Context{Ctx: context.Background()}

	t.Run("wait for future timestamp", func(t *testing.T) {
		mockSleeper := &mockSleeper{}
		// Set mockClock.now to a time well before the futureTime
		mockClock := &mockClock{now: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}
		node := &WaitForNode{TimestampKey: "triggerTime", Sleeper: mockSleeper, Clock: mockClock}

		// Set futureTime to be 100ms after mockClock.now
		futureTime := mockClock.Now().Add(100 * time.Millisecond).Format(time.RFC3339Nano)
		inputs := []map[string]interface{}{
			{"id": 1, "triggerTime": futureTime},
		}

		_, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check that it attempted to sleep for approximately the correct duration
		if mockSleeper.sleptFor < 90*time.Millisecond || mockSleeper.sleptFor > 110*time.Millisecond {
			t.Errorf("expected to sleep for ~100ms, but slept for %v", mockSleeper.sleptFor)
		}
	})

	t.Run("do not wait for past timestamp", func(t *testing.T) {
		mockSleeper := &mockSleeper{}
		mockClock := &mockClock{now: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).Add(200 * time.Millisecond)}
		node := &WaitForNode{TimestampKey: "triggerTime", Sleeper: mockSleeper, Clock: mockClock}
		pastTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
		inputs := []map[string]interface{}{
			{"id": 1, "triggerTime": pastTime},
		}

		_, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check that it did not sleep
		if mockSleeper.sleptFor != 0 {
			t.Errorf("expected to not sleep, but slept for %v", mockSleeper.sleptFor)
		}
	})

	t.Run("handle missing timestamp key", func(t *testing.T) {
		node := NewWaitForNode("nonExistentTime")
		inputs := []map[string]interface{}{
			{"id": 1},
		}
		_, err := node.Execute(ctx, inputs)
		if err == nil || err.Error() != "timestamp key 'nonExistentTime' not found in input record" {
			t.Errorf("expected error for missing key, got %v", err)
		}
	})

	t.Run("handle invalid timestamp format", func(t *testing.T) {
		node := NewWaitForNode("triggerTime")
		inputs := []map[string]interface{}{
			{"id": 1, "triggerTime": "not-a-timestamp"},
		}
		_, err := node.Execute(ctx, inputs)
		if err == nil || !strings.Contains(err.Error(), "failed to parse timestamp") {
			t.Errorf("expected error for invalid timestamp, got %v", err)
		}
	})

	t.Run("handle non-string timestamp value", func(t *testing.T) {
		node := NewWaitForNode("triggerTime")
		inputs := []map[string]interface{}{
			{"id": 1, "triggerTime": 12345},
		}
		_, err := node.Execute(ctx, inputs)
		if err == nil || err.Error() != "timestamp value for key 'triggerTime' is not a string" {
			t.Errorf("expected error for non-string timestamp, got %v", err)
		}
	})

	t.Run("empty inputs", func(t *testing.T) {
		node := NewWaitForNode("triggerTime")
		inputs := []map[string]interface{}{}
		out, err := node.Execute(ctx, inputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) != 0 {
			t.Errorf("expected empty output, got %v", out)
		}
	})
}
