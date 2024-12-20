package binance

type DepthUpdate struct {
    Stream string `json:"stream"`
    Data   struct {
        Bids          [][2]string      `json:"bids"`
        Asks          [][2]string      `json:"asks"`
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
