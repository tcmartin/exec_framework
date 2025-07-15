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

// Context holds shared clients and config
type Context struct {
    Ctx            context.Context
    HTTPClient     *retryablehttp.Client
    LangChain      LangChainClient
    DynamoDBClient *dynamodb.Client
    Logger         *zap.SugaredLogger
    Metrics        *Metrics
    Env            map[string]string
}

// Node represents a workflow step
// Execute runs the step with the given input and returns output records
type Node interface {
    Execute(ctx *Context, input []map[string]interface{}) ([]map[string]interface{}, error)
}