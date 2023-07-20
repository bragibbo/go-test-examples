package dal

import (
	"context"
	"errors"
	"go-test-examples/interfaces/db"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDB struct {
	mockGetItem func(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error)
	mockBatchGetItem func() error
}

func (mock *mockDynamoDB) GetItem(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error) {
	return mock.mockGetItem(ctx, key, table)
}

func (mock *mockDynamoDB) BatchGetItem() error {
	return mock.mockBatchGetItem()
}

func TestDAL_GetUser(t *testing.T) {
	type fields struct {
		dynamo db.DynamoDBInter
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Returns expected user struct from id",
			fields: fields{
				dynamo: &mockDynamoDB{
					mockGetItem: func(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error) {
						userId := "freddy"
						if v, exist := key["id"]; !exist {
							var id string
							if attributevalue.Unmarshal(v, id); id != userId {
								return map[string]types.AttributeValue{}, nil
							}
						}
						return map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberS{Value: userId},
							"name":  &types.AttributeValueMemberS{Value: "freddy"},
							"age":   &types.AttributeValueMemberN{Value: "55"},
							"email": &types.AttributeValueMemberS{Value: "freddy@mercury.com"},
						}, nil
					},
					mockBatchGetItem: func() error {
						return errors.New("hello")
					},
				},
			},
			args: args{
				id: "freddy",
			},
			want: &User{
				ID:    "freddy",
				Name:  "freddy",
				Age:   55,
				Email: "freddy@mercury.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dal := &DAL{
				dynamo: tt.fields.dynamo,
			}
			got, err := dal.GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DAL.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DAL.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
