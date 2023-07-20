package dal

import (
	"context"
	"go-test-examples/monkeypatching/db"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Simple
func TestGetUserSimple(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Returns expected user struct from id",
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
		db.GetItem = func(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error) {
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
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Alternate
func TestGetUserAlternate(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		want        *User
		wantErr     bool
		mockGetItem func(ctx context.Context, key map[string]types.AttributeValue, table string) (map[string]types.AttributeValue, error)
	}{
		{
			name: "Returns expected user struct from id",
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
		},
	}
	for _, tt := range tests {
		db.GetItem = tt.mockGetItem
		
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
