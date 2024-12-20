package binance

import "sync"

type OrderBook struct {
    Bids map[string][][2]string
    Asks map[string][][2]string
    mu   sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: make(map[string][][2]string),
		Bids: make(map[string][][2]string),
	}
}

func (o *OrderBook) Update(symbol string, update DepthUpdate) {
    o.mu.Lock()
    defer o.mu.Unlock()

	o.Asks[symbol] = update.Data.Asks
	o.Bids[symbol] = update.Data.Bids
}

func (ob *OrderBook) GetAsks() map[string][][2]string {
    ob.mu.Lock()
    defer ob.mu.Unlock()

    asksCopy := make(map[string][][2]string)
    for k, v := range ob.Asks {
        asksCopy[k] = v
    }
    return asksCopy
}

func (ob *OrderBook) GetBids() map[string][][2]string {
    ob.mu.Lock()
    defer ob.mu.Unlock()

    bidsCopy := make(map[string][][2]string)
    for k, v := range ob.Bids {
        bidsCopy[k] = v
    }
    return bidsCopy
}