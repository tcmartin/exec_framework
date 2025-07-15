package store

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func TestNewClient(t *testing.T) {
	// Use static credentials for testing to avoid making real AWS calls.
	ctx := context.WithValue(context.Background(), "aws.config", &aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider("AKID", "SECRET_KEY", "SESSION_TOKEN"),
		Region:      "us-west-2",
	})

	// Temporarily set the AWS_CONFIG_FILE environment variable to a non-existent file
	// to prevent the test from loading credentials from the user's machine.
	t.Setenv("AWS_CONFIG_FILE", "/dev/null")

	_, err := NewClient(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
