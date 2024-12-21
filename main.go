package main

import (
	"algotrade/pkg/algo"
	"algotrade/pkg/binance"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	marketDetails := binance.NewMarketDetails(ctx)
	marketDetails.SetProducts()
	marketDetails.InitDepthStreams()
	marketDetails.StartUpdateHandler()
	marketDetails.StartDepthStreams()

	algo := algo.NewArbitrageDetector(ctx, marketDetails)
	go algo.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	marketDetails.StopStreams()
	fmt.Println("Program exiting gracefully...")
}
