package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrNotFound = errors.New("resource not found")

type dynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type DynamoDBInter interface {
	GetItem(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error)
}

type DynamoDB struct {
	client dynamoAPI
}

func InitDynamoDB(ctx context.Context) (*DynamoDB, error) {
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	dynamodb := dynamodb.NewFromConfig(cfg)
	return &DynamoDB{client: dynamodb}, nil
}

func (d *DynamoDB) GetItem(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error) {
	output, err := d.client.GetItem(ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(table),
			Key:       key,
		},
	)
	if err != nil {
		return nil, err
	}
	return output.Item, nil
}
