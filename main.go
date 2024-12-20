package main

import (
	"algotrade/pkg/binance"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	marketDetails := binance.NewMarketDetails(ctx)
	marketDetails.SetProducts()
	marketDetails.InitDepthStreams()
	marketDetails.StartUpdateHandler()
	marketDetails.StartDepthStreams()

	go func(marketDetails *binance.MarketDetails) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Using orderBook at:", time.Now())

			asks := marketDetails.GetAsks()
			for symbol, ask := range asks {
				fmt.Printf("Asks for %s: %v\n", symbol, ask)
			}

			bids := marketDetails.GetBids()
			for symbol, bid := range bids {
				fmt.Printf("Bids for %s: %v\n", symbol, bid)
			}
		}
	}(marketDetails)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	marketDetails.StopStreams()
	fmt.Println("Program exiting gracefully...")
}
