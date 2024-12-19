package main

import (
	"algotrade/pkg/binance"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	subscribedSymbols := []string{"btcusdt", "ethusdt", "bnbusdt"}
	depthStream, err := binance.NewDepthStream(subscribedSymbols)
	if err != nil {
		fmt.Println(err)
	}
	depthStream.Start(ctx)
	defer depthStream.Stop()
}