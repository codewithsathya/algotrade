package main

import (
	"algotrade/pkg/binance"
	"algotrade/pkg/utils"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startDepthStreams(ctx context.Context, symbols []string, batchSize int) (<-chan binance.DepthUpdate, []*binance.DepthStream, <-chan error) {
    batches := utils.ChunkSlice(symbols, batchSize)
    streams := make([]*binance.DepthStream, len(batches))
    allUpdates := make(chan binance.DepthUpdate)
    errCh := make(chan error, len(batches))

    for i, batch := range batches {
        stream, err := binance.NewDepthStream(batch)
        if err != nil {
            errCh <- fmt.Errorf("failed to create stream for batch %v: %w", batch, err)
            continue
        }

        go func(stream *binance.DepthStream) {
            for update := range stream.Updates() {
                select {
                case allUpdates <- update:
                case <-ctx.Done():
                    return
                }
            }
        }(stream)

        go func(stream *binance.DepthStream) {
            if err := stream.Start(ctx); err != nil {
                errCh <- fmt.Errorf("error occurred in stream: %w", err)
            }
        }(stream)

        streams[i] = stream
        fmt.Println("Created stream for:", batch)

        if i < len(batches)-1 {
            time.Sleep(1 * time.Second)
        }
    }

    return allUpdates, streams, errCh
}

func main() {
	client := binance.NewBinanceClient()
	resp, err := client.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(resp.Data))

	orderBook := binance.NewOrderBook()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscribedSymbols := []string{}
	for _, product := range resp.Data {
		subscribedSymbols = append(subscribedSymbols, product.BaseAsset+product.QuoteAsset)
	}
	allUpdates, streams, errCh := startDepthStreams(ctx, subscribedSymbols, 50)
	go func() {
		for err := range errCh {
			fmt.Printf("Error: %v\n", err)
		}
	}()
	fmt.Printf("%d streams created\n", len(streams))

	go func(ob *binance.OrderBook) {
		for update := range allUpdates {
			// fmt.Println(update.Stream)
			// fmt.Println(len(update.Data.Asks))
			ob.Update(update.Stream, update)
		}
	}(orderBook)

	go func(orderBook *binance.OrderBook) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Using orderBook at:", time.Now())
			
			asks := orderBook.GetAsks()
			for symbol, ask := range asks {
				fmt.Printf("Asks for %s: %v\n", symbol, ask)
			}

			bids := orderBook.GetBids()
			for symbol, bid := range bids {
				fmt.Printf("Bids for %s: %v\n", symbol, bid)
			}
		}
	}(orderBook)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	fmt.Println("Program exiting gracefully...")
}