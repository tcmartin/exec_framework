package framework

import (
    "context"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
    "github.com/tmc/langchaingo/llms"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "go.uber.org/zap"
)

// LangChainClient defines only the needed ChatCompletion method
// implemented by *langchaingo.Client
type LangChainClient interface {
    GenerateFromSinglePrompt(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

// DynamoDBPutItemAPI defines the interface for the PutItem function.
// We use this interface to test the code without needing a real DynamoDB instance.
type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// DynamoDBClientFactory defines a function type for creating DynamoDB clients dynamically.
type DynamoDBClientFactory func(ctx context.Context, region, accessKeyID, secretAccessKey string) (DynamoDBPutItemAPI, error)

// Context holds shared clients and config
type Context struct {
    Ctx            context.Context
    HTTPClient     *retryablehttp.Client
    LangChain      LangChainClient
    DynamoDBClient DynamoDBPutItemAPI
    DynamoDBClientFactory DynamoDBClientFactory
    Logger         *zap.SugaredLogger
    Metrics        *Metrics
    Env            map[string]string
}

// Node represents a workflow step
// Execute runs the step with the given input and returns output records
type Node interface {
    Execute(ctx *Context, input []map[string]interface{}) ([]map[string]interface{}, error)
}