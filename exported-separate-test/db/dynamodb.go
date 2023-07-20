package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrNotFound = errors.New("resource not found")

type dynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

type DynamoDBInter interface {
	BatchGetItem(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error)
}

type DynamoDB struct {
	client dynamoAPI
}

func InitDynamoDB(ctx context.Context, region string) (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	dynamodb := dynamodb.NewFromConfig(cfg)
	return &DynamoDB{client: dynamodb}, nil
}

func (d *DynamoDB) BatchGetItem(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error) {
	var err error
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	out := []map[string]types.AttributeValue{}
	for start < len(keys) {
		var getReqs []map[string]types.AttributeValue
		if end > len(keys) {
			end = len(keys)
		}
		getReqs = append(getReqs, keys[start:end]...)
		output, err := d.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{table: {Keys: getReqs}}})
		if err != nil {
			fmt.Println(err, "table", table)
		} else {
			out = append(out, output.Responses[table]...)
		}
		start = end
		end += batchSize
	}

	return out, err
}
