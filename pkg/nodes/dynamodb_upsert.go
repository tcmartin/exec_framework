package nodes

import (
    "go-workflow/pkg/framework"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBUpsert writes items
type DynamoDBUpsert struct{ TableName string }

func (n *DynamoDBUpsert) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    for _, rec := range inputs {
        item := map[string]types.AttributeValue{}
        for k, v := range rec {
            if s, ok := v.(string); ok {
                item[k] = &types.AttributeValueMemberS{Value: s}
            }
        }
        _, err := ctx.DynamoDBClient.PutItem(ctx.Ctx, &dynamodb.PutItemInput{TableName: &n.TableName, Item: item})
        if err != nil {
            return nil, err
        }
    }
    return inputs, nil
}