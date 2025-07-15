package nodes

import (
	"context"
	"go-workflow/pkg/framework"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBPutItemAPI defines the interface for the PutItem function.
// We use this interface to test the code without needing a real DynamoDB instance.
type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// mockDynamoDBClient is a mock implementation of the DynamoDB client for testing.
type mockDynamoDBClient struct {
	putItem func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m *mockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.putItem != nil {
		return m.putItem(ctx, params, optFns...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func TestDynamoDBUpsert_Execute(t *testing.T) {
	node := &DynamoDBUpsert{TableName: "test-table"}

	var capturedInput *dynamodb.PutItemInput
	mockClient := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			capturedInput = params
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	ctx := &framework.Context{
		Ctx:            context.Background(),
		DynamoDBClient: mockClient,
	}

	inputs := []map[string]interface{}{{"id": "123", "name": "test"}}
	_, err := node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedInput == nil {
		t.Fatal("PutItem was not called")
	}
	if *capturedInput.TableName != "test-table" {
		t.Errorf("expected table name test-table, got %s", *capturedInput.TableName)
	}
	if _, ok := capturedInput.Item["id"].(*types.AttributeValueMemberS); !ok {
		t.Error("expected id to be a string attribute")
	}
	if capturedInput.Item["id"].(*types.AttributeValueMemberS).Value != "123" {
		t.Errorf("expected id to be 123, got %s", capturedInput.Item["id"].(*types.AttributeValueMemberS).Value)
	}
}
