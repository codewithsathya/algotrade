package binance

import (
	"algotrade/pkg/httpclient"
	"encoding/json"
	"fmt"
)

type Binance struct {
	client *httpclient.HttpClient
}

func (b *Binance) New() {
	b.client = &httpclient.HttpClient{}
	b.client.New()
}

func (b *Binance) GetAllProducts() (GetProductsResponse, error) {
	data, err := b.client.Get("https://www.binance.com/bapi/asset/v2/public/asset-service/product/get-products?includeEtf=true");
	if err != nil {
		return GetProductsResponse{}, fmt.Errorf(`failed to get all products, error: %v`, err)
	}
	var response GetProductsResponse;
	err = json.Unmarshal([]byte(data), &response);
	if err != nil {
		return GetProductsResponse{}, fmt.Errorf(`failed to map get products response, error: %v`, err);
	}
	return response, nil;
}

