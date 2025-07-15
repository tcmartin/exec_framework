package nodes

import (
	"context"
	"errors"
	"go-workflow/pkg/framework"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}

	var capturedInput *dynamodb.PutItemInput
	mockClient := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			capturedInput = params
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return mockClient, nil
		},
	}

	inputs := []map[string]interface{}{
		{
			"tableName":        "test-table",
			"awsRegion":        "us-east-1",
			"awsAccessKeyID":   "test-access-key",
			"awsSecretAccessKey": "test-secret-key",
			"id":               "123",
			"name":             "test",
		},
	}
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

func TestDynamoDBUpsert_Execute_PutItemError(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}

	mockClient := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("PutItem error")
		},
	}

	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return mockClient, nil
		},
	}

	inputs := []map[string]interface{}{
		{
			"tableName":        "test-table",
			"awsRegion":        "us-east-1",
			"awsAccessKeyID":   "test-access-key",
			"awsSecretAccessKey": "test-secret-key",
			"id":               "123",
			"name":             "test",
		},
	}
	_, err := node.Execute(ctx, inputs)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
	if err.Error() != "PutItem error" {
		t.Errorf("expected PutItem error, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_MissingTableNameKey(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return &mockDynamoDBClient{}, nil
		},
	}
	inputs := []map[string]interface{}{{"awsRegion": "us-east-1", "awsAccessKeyID": "test", "awsSecretAccessKey": "test"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "table name not found or not a string in input record for key tableName" {
		t.Errorf("expected error for missing TableNameKey, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_MissingAWSRegionKey(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return &mockDynamoDBClient{}, nil
		},
	}
	inputs := []map[string]interface{}{{"tableName": "test-table", "awsAccessKeyID": "test", "awsSecretAccessKey": "test"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "AWS region not found or not a string in input record for key awsRegion" {
		t.Errorf("expected error for missing AWSRegionKey, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_MissingAWSAccessKeyIDKey(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return &mockDynamoDBClient{}, nil
		},
	}
	inputs := []map[string]interface{}{{"tableName": "test-table", "awsRegion": "us-east-1", "awsSecretAccessKey": "test"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "AWS access key ID not found or not a string in input record for key awsAccessKeyID" {
		t.Errorf("expected error for missing AWSAccessKeyIDKey, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_MissingAWSSecretAccessKeyKey(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return &mockDynamoDBClient{}, nil
		},
	}
	inputs := []map[string]interface{}{{"tableName": "test-table", "awsRegion": "us-east-1", "awsAccessKeyID": "test-access-key"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "AWS secret access key not found or not a string in input record for key awsSecretAccessKey" {
		t.Errorf("expected error for missing AWSSecretAccessKeyKey, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_ClientFactoryError(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{
		Ctx: context.Background(),
		DynamoDBClientFactory: func(ctx context.Context, region, accessKeyID, secretAccessKey string) (framework.DynamoDBPutItemAPI, error) {
			return nil, errors.New("client factory error")
		},
	}
	inputs := []map[string]interface{}{{"tableName": "test-table", "awsRegion": "us-east-1", "awsAccessKeyID": "test-access-key", "awsSecretAccessKey": "test-secret-key"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "failed to create DynamoDB client: client factory error" {
		t.Errorf("expected client factory error, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_NoClientAvailable(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey:        "tableName",
		AWSRegionKey:        "awsRegion",
		AWSAccessKeyIDKey:   "awsAccessKeyID",
		AWSSecretAccessKeyKey: "awsSecretAccessKey",
	}
	ctx := &framework.Context{Ctx: context.Background(), DynamoDBClient: nil, DynamoDBClientFactory: nil}
	inputs := []map[string]interface{}{{"tableName": "test-table", "awsRegion": "us-east-1", "awsAccessKeyID": "test-access-key", "awsSecretAccessKey": "test-secret-key"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "no DynamoDB client or factory available" {
		t.Errorf("expected error for no client available, got %v", err)
	}
}

func TestDynamoDBUpsert_Execute_DefaultClient(t *testing.T) {
	node := &DynamoDBUpsert{
		TableNameKey: "tableName",
	}
	var capturedInput *dynamodb.PutItemInput
	mockClient := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			capturedInput = params
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	ctx := &framework.Context{Ctx: context.Background(), DynamoDBClient: mockClient, DynamoDBClientFactory: nil}
	inputs := []map[string]interface{}{{"tableName": "test-table", "id": "123", "name": "test"}}
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
}
