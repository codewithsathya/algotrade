package binance

import (
	"algotrade/pkg/utils"
	"context"
	"fmt"
	"time"
)

type MarketDetails struct {
	products          []GetProductsSymbol
	subscribedSymbols []string
	orderBook         *OrderBook
	streams           []*DepthStream
    streamsRunning    bool
	ctx               context.Context
	cancelFunc        context.CancelFunc
	streamsCtx        context.Context
	streamsCancelFunc context.CancelFunc
	client            *BinanceClient
	errCh             chan error
	updateChannel     chan DepthUpdate
	batches           [][]string
	batchSize         int
}

func NewMarketDetails(ctx context.Context) *MarketDetails {
	return &MarketDetails{
		products:          []GetProductsSymbol{},
		subscribedSymbols: []string{},
		orderBook:         NewOrderBook(),
		streams:           []*DepthStream{},
        streamsRunning:    false,
		ctx:               ctx,
		cancelFunc:        nil,
		client:            NewBinanceClient(),
		errCh:             make(chan error),
		updateChannel:     make(chan DepthUpdate),
		batchSize:         50,
	}
}

func (md *MarketDetails) SetProducts() error {
	resp, err := md.client.GetAllProducts()
	if err != nil {
		return fmt.Errorf("failed to get all products | %v", err)
	}
	md.products = resp.Data
	subscribedSymbols := []string{}
	for _, product := range resp.Data {
		subscribedSymbols = append(subscribedSymbols, product.BaseAsset+product.QuoteAsset)
	}
	md.subscribedSymbols = subscribedSymbols
	return nil
}

func (md *MarketDetails) InitDepthStreams() {
	md.batches = utils.ChunkSlice(md.subscribedSymbols, 50)
	md.streams = make([]*DepthStream, len(md.batches))

	for i, batch := range md.batches {
		md.streams[i] = NewDepthStream(batch)
		go func(stream *DepthStream) {
			for update := range stream.GetUpdatesChannel() {
				select {
				case md.updateChannel <- update:
				case <-md.ctx.Done():
					return
				}
			}
		}(md.streams[i])
	}
}

func (md *MarketDetails) StartUpdateHandler() {
	go func() {
		for {
			select {
			case err := <-md.errCh:
				fmt.Printf("error occurred in websocket goroutine: %v\n", err)
				if md.streamsCancelFunc != nil {
                    md.streamsCancelFunc()
				}
                md.streamsRunning = false
			case <-md.ctx.Done():
				return
			}
		}
	}()
	go func(ob *OrderBook) {
		for {
			select {
			case update := <-md.updateChannel:
				ob.Update(update.Stream, update)
			case <-md.ctx.Done():
				return
			}
		}
	}(md.orderBook)
}

func (md *MarketDetails) StartDepthStreams() {
    if md.streamsRunning {
        md.streamsCancelFunc()
    }
    md.streamsCtx, md.streamsCancelFunc = context.WithCancel(context.Background())
    for i, stream := range md.streams {
        go func(stream *DepthStream) {
            if err := stream.Start(md.streamsCtx); err != nil {
                md.errCh <- fmt.Errorf("error occurred in stream: %w", err)
            }
        }(stream)
        if i < len(md.streams) - 1 {
            time.Sleep(1 * time.Second)
        }
    }
    md.streamsRunning = true
}

func (md *MarketDetails) StopStreams() {
    if md.streamsRunning {
        md.streamsCancelFunc()
        md.streamsRunning = false
    }
}

func (md *MarketDetails) GetAsks() map[string][][2]string {
    return md.orderBook.GetAsks()
}

func (md *MarketDetails) GetBids() map[string][][2]string {
    return md.orderBook.GetBids()
}
