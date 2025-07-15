package store

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient returns a DynamoDB client
func NewClient(ctx context.Context, region, accessKeyID, secretAccessKey string) (*dynamodb.Client, error) {
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
    )
    if err != nil {
        return nil, err
    }
    return dynamodb.NewFromConfig(cfg), nil
}