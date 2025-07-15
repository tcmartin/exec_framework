package nodes

import (
    "errors"
    "fmt"
    "go-workflow/pkg/framework"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBUpsert writes items
type DynamoDBUpsert struct{
    TableNameKey        string
    AWSRegionKey        string
    AWSAccessKeyIDKey   string
    AWSSecretAccessKeyKey string
}

func (n *DynamoDBUpsert) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    for _, rec := range inputs {
        tableName, ok := rec[n.TableNameKey].(string)
        if !ok {
            return nil, fmt.Errorf("table name not found or not a string in input record for key %s", n.TableNameKey)
        }

        var client framework.DynamoDBPutItemAPI
        if n.AWSRegionKey != "" && n.AWSAccessKeyIDKey != "" && n.AWSSecretAccessKeyKey != "" && ctx.DynamoDBClientFactory != nil {
            region, ok := rec[n.AWSRegionKey].(string)
            if !ok {
                return nil, fmt.Errorf("AWS region not found or not a string in input record for key %s", n.AWSRegionKey)
            }

            accessKeyID, ok := rec[n.AWSAccessKeyIDKey].(string)
            if !ok {
                return nil, fmt.Errorf("AWS access key ID not found or not a string in input record for key %s", n.AWSAccessKeyIDKey)
            }

            secretAccessKey, ok := rec[n.AWSSecretAccessKeyKey].(string)
            if !ok {
                return nil, fmt.Errorf("AWS secret access key not found or not a string in input record for key %s", n.AWSSecretAccessKeyKey)
            }
            var err error
            client, err = ctx.DynamoDBClientFactory(ctx.Ctx, region, accessKeyID, secretAccessKey)
            if err != nil {
                return nil, fmt.Errorf("failed to create DynamoDB client: %w", err)
            }
        } else if ctx.DynamoDBClient != nil {
            client = ctx.DynamoDBClient
        } else {
            return nil, errors.New("no DynamoDB client or factory available")
        }

        item := map[string]types.AttributeValue{}
        for k, v := range rec {
            if s, ok := v.(string); ok {
                item[k] = &types.AttributeValueMemberS{Value: s}
            }
        }
        _, err := client.PutItem(ctx.Ctx, &dynamodb.PutItemInput{TableName: &tableName, Item: item})
        if err != nil {
            return nil, err
        }
    }
    return inputs, nil
}