package store

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Use static credentials for testing to avoid making real AWS calls.
	ctx := context.TODO()

	// Temporarily set the AWS_CONFIG_FILE environment variable to a non-existent file
	// to prevent the test from loading credentials from the user's machine.
	t.Setenv("AWS_CONFIG_FILE", "/dev/null")

	_, err := NewClient(ctx, "us-west-2", "AKID", "SECRET_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}