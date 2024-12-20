package main

import (
	"algotrade/pkg/binance"
	"context"
	"fmt"
	"time"
)

func main() {
	client := binance.NewBinanceClient()
	resp, err := client.GetAllProducts()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(resp.Data))

	orderBook := binance.NewOrderBook()

	ctx := context.Background()

	subscribedSymbols := []string{"btcusdt", "ethusdt", "bnbusdt"}
	depthStream, err := binance.NewDepthStream(subscribedSymbols)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		if err := depthStream.Start(ctx); err != nil {
			fmt.Printf("DepthStream encountered an error: %v\n", err)
		}
	}()

	for update := range depthStream.Updates() {
		orderBook.Update(update.Stream, update)
		fmt.Println(orderBook.Asks[update.Stream])
	}

	time.Sleep(10 * time.Second)
	depthStream.Stop()
}