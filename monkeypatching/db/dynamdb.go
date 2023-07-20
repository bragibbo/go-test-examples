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

var (
	GetItem      = getItem
	dynamoClient *dynamodb.Client
)

func initDynamoDB(ctx context.Context) (*dynamodb.Client, error) {
	if dynamoClient != nil {
		return dynamoClient, nil
	}
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	dynamodb := dynamodb.NewFromConfig(cfg)
	return dynamodb, nil
}

func getItem(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error) {
	dc, _ := initDynamoDB(ctx)
	output, err := dc.GetItem(ctx,
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
