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

type GetProductsResponse struct {
	Code          string    `json:"code"`
	Message       *string   `json:"message"`
	MessageDetail *string   `json:"messageDetail"`
	Data          []GetProductsSymbol `json:"data"`
	Success       bool      `json:"success"`
}

type GetProductsSymbol struct {
	Symbol                string   `json:"s"`
	Status                string   `json:"st"`
	BaseAsset             string   `json:"b"`
	QuoteAsset            string   `json:"q"`
	BaseAssetFullName     string   `json:"an"`
	QuoteAssetFullName    string   `json:"qn"`
	OpenPrice             string   `json:"o"`
	HighPrice             string   `json:"h"`
	LowPrice              string   `json:"l"`
	ClosePrice            string   `json:"c"`
	Volume                string   `json:"v"`
	Tags                  []string `json:"tags"`
	IsETF                 bool     `json:"etf"`
}
