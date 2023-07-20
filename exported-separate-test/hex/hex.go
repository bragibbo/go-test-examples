package hex

import (
	"context"

	"go-test-examples/exported-separate-test/db"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Hex struct {
	ID         int64  `json:"hex" dynamodbav:"hex"`
	H303       int64  `json:"h3_03" dynamodbav:"h3_03"`
	CityName   string `json:"city_name" dynamodbav:"place_name"`
	StateAbbr  string `json:"state_abbr" dynamodbav:"state_abbr"`
	Resolution int    `json:"resolution" dynamodbav:"resolution"`
}

type HexQueryServiceConfig struct {
	Region string
	Table  string
}

type HexQueryService struct {
	config   HexQueryServiceConfig
	HexCache map[int64]Hex
	Client   db.DynamoDBInter
}

func NewHexQueryService(c context.Context, config HexQueryServiceConfig) (*HexQueryService, error) {
	client, err := db.InitDynamoDB(context.TODO(), config.Region)
	if err != nil {
		return nil, err
	}

	svc := HexQueryService{
		Client:   client,
		HexCache: map[int64]Hex{},
		config:   config,
	}
	return &svc, nil
}

func (h *HexQueryService) cacheHexes(hexes []Hex) {
	for _, hex := range hexes {
		h.HexCache[hex.ID] = hex
	}
}

func (h *HexQueryService) uncachedHexes(hexIDs []int64) []int64 {
	uncachedHexIDs := []int64{}
	for _, id := range hexIDs {
		if _, exists := h.HexCache[id]; !exists {
			uncachedHexIDs = append(uncachedHexIDs, id)
		}
	}
	return uncachedHexIDs
}

func (h *HexQueryService) hexes(hexIDs []int64) []Hex {
	hexes := []Hex{}
	for _, id := range hexIDs {
		if hex, exists := h.HexCache[id]; exists {
			hexes = append(hexes, hex)
		}
	}
	return hexes
}

func (h *HexQueryService) cacheUncachedHexes(hexIDs []int64) error {
	uncachedHexIDs := h.uncachedHexes(hexIDs)
	keys := []map[string]types.AttributeValue{}
	for _, id := range uncachedHexIDs {
		pk, _ := attributevalue.Marshal(id)
		keys = append(keys, map[string]types.AttributeValue{"hex": pk})
	}
	output, err := h.Client.BatchGetItem(context.TODO(), keys, h.config.Table)
	if err != nil {
		return err
	}
	hexes := make([]Hex, len(output))
	attributevalue.UnmarshalListOfMaps(output, &hexes)
	h.cacheHexes(hexes)
	return nil
}

func (h *HexQueryService) GetHexes(hexIDs []int64) ([]Hex, error) {
	err := h.cacheUncachedHexes(hexIDs)
	if err != nil {
		return nil, err
	}

	hexes := h.hexes(hexIDs)
	return hexes, nil
}
