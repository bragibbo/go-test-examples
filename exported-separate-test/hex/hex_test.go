package hex_test

import (
	"context"
	"errors"
	"go-test-examples/exported-separate-test/db"
	"go-test-examples/exported-separate-test/hex"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDB struct {
	mockBatchGetItem func(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error)
}

func (db *mockDynamoDB) BatchGetItem(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error) {
	return db.mockBatchGetItem(ctx, keys, table)
}

func TestHexQueryService_GetHexes(t *testing.T) {
	type fields struct {
		hexCache map[int64]hex.Hex
		client   db.DynamoDBInter
	}
	type args struct {
		hexIDs []int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []hex.Hex
		wantErr bool
	}{
		{
			name: "Returns correct number of hexes and respects the cache",
			fields: fields{
				client: &mockDynamoDB{
					mockBatchGetItem: func(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error) {
						want := []map[string]types.AttributeValue{
							{"hex": &types.AttributeValueMemberN{Value: "613168376124014591"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375440343039"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375941562367"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375935270911"}},
						}
						if !reflect.DeepEqual(keys, want) {
							t.Errorf("HexQueryService.GetHexes() = %v, want %v", keys, want)
						}
						return []map[string]types.AttributeValue{
							{"hex": &types.AttributeValueMemberN{Value: "613168376124014591"}, "place_name": &types.AttributeValueMemberS{Value: "layton"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375440343039"}, "place_name": &types.AttributeValueMemberS{Value: "sandy"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375941562367"}, "place_name": &types.AttributeValueMemberS{Value: "bluffdale"}},
							{"hex": &types.AttributeValueMemberN{Value: "613168375935270911"}, "place_name": &types.AttributeValueMemberS{Value: "draper"}},
						}, nil
					},
				},
				hexCache: map[int64]hex.Hex{
					613168375937368063: {
						ID:       613168375937368063,
						CityName: "lehi",
					},
					613168375945756671: {
						ID:       613168375945756671,
						CityName: "provo",
					},
					613168376132403199: {
						ID:       613168376132403199,
						CityName: "riverton",
					},
				},
			},
			args: args{
				hexIDs: []int64{
					// "Cached" hex ids
					613168375937368063,
					613168375945756671,
					613168376132403199,
					// "Uncached" hex ids
					613168376124014591,
					613168375440343039,
					613168375941562367,
					613168375935270911},
			},
			want: []hex.Hex{
				{ID: 613168375937368063, CityName: "lehi"},
				{ID: 613168375945756671, CityName: "provo"},
				{ID: 613168376132403199, CityName: "riverton"},
				{ID: 613168376124014591, CityName: "layton"},
				{ID: 613168375440343039, CityName: "sandy"},
				{ID: 613168375941562367, CityName: "bluffdale"},
				{ID: 613168375935270911, CityName: "draper"},
			},
			wantErr: false,
		},
		{
			name: "Returns error if db query fails",
			fields: fields{
				client: &mockDynamoDB{
					mockBatchGetItem: func(ctx context.Context, keys []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error) {
						return nil, errors.New("some error")
					},
				},
				hexCache: map[int64]hex.Hex{},
			},
			args: args{
				hexIDs: []int64{608664766742790143, 608664766205919231},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hex.HexQueryService{
				HexCache: tt.fields.hexCache,
				Client:   tt.fields.client,
			}
			got, err := h.GetHexes(tt.args.hexIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexQueryService.GetHexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexQueryService.GetHexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
