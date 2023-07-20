package dal

import (
	"context"
	"go-test-examples/monkeypatching/db"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var usersTable = "users-table"

type User struct {
	ID    string
	Name  string
	Age   int
	Email string
}

func GetUser(id string) (*User, error) {
	pk, _ := attributevalue.Marshal(id)
	key := map[string]types.AttributeValue{"id": pk}
	output, err := db.GetItem(context.TODO(), key, usersTable)
	if err != nil {
		return nil, err
	}
	var user User
	attributevalue.UnmarshalMap(output, &user)
	return &user, nil
}
