package framework

import (
    "context"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
    "github.com/tmc/langchaingo"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "go.uber.org/zap"
)

// LangChainClient defines the subset of langchaingo.Client used by the workflow
type LangChainClient interface {
    ChatCompletion(ctx context.Context, req langchaingo.ChatCompletionRequest) (langchaingo.ChatCompletionResponse, error)
}

// Context holds shared clients and configuration for nodes
type Context struct {
    Ctx            context.Context
    HTTPClient     *retryablehttp.Client
    LangChain      LangChainClient
    DynamoDBClient *dynamodb.Client
    Logger         *zap.SugaredLogger
    Metrics        *Metrics
    Env            map[string]string
}
