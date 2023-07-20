package app

import (
	"context"
	"fmt"
	"go-test-examples/exported-separate-test/hex"
)

func PrintHexes() {
	hexSvc, err := hex.NewHexQueryService(context.TODO(), hex.HexQueryServiceConfig{})
	if err != nil {
		fmt.Println(err)
	}

	hexInfo, err := hexSvc.GetHexes([]int64{})
	if err != nil {
		fmt.Println(err)
	}

	for _, hex := range hexInfo {
		fmt.Println(hex)
	}
}
