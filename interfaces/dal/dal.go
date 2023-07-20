package dal

import (
	"context"
	"go-test-examples/interfaces/db"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var usersTable = "users-table"

type User struct {
	ID    string `dynamodbav:"id"`
	Name  string `dynamodbav:"name"`
	Age   int `dynamodbav:"age"`
	Email string `dynamodbav:"email"`
}

type DAL struct {
	dynamo db.DynamoDBInter
}

func InitDAL() *DAL {
	dynamoDB, _ := db.InitDynamoDB(context.TODO())
	return &DAL{dynamo: dynamoDB}
}

func (dal *DAL) GetUser(id string) (*User, error) {
	pk, _ := attributevalue.Marshal(id)
	key := map[string]types.AttributeValue{"id": pk}
	output, err := dal.dynamo.GetItem(context.TODO(), key, usersTable)
	if err != nil {
		return nil, err
	}
	var user User
	attributevalue.UnmarshalMap(output, &user)
	return &user, nil
}
