package binance

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