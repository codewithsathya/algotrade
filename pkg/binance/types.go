package binance

type DepthUpdate struct {
    Stream string `json:"stream"`
    Data   struct {
        Event         string          `json:"e"`
        EventTime     int64           `json:"E"`
        Symbol        string          `json:"s"`
        FirstUpdateID int64           `json:"U"`
        FinalUpdateID int64           `json:"u"`
        Bids          [][]string      `json:"b"`
        Asks          [][]string      `json:"a"`
    } `json:"data"`
}