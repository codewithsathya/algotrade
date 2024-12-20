package binance

import (
	"algotrade/pkg/httpclient"
	"encoding/json"
	"fmt"
)

type BinanceClient struct {
	client *httpclient.HttpClient
}

func NewBinanceClient() *BinanceClient {
	bClient := BinanceClient{
	}
	bClient.client = &httpclient.HttpClient{}
	bClient.client.New()
	return &bClient
}

func (b *BinanceClient) GetAllProducts() (GetProductsResponse, error) {
	data, err := b.client.Get("https://www.binance.com/bapi/asset/v2/public/asset-service/product/get-products?includeEtf=true")
	if err != nil {
		return GetProductsResponse{}, fmt.Errorf(`failed to get all products, error: %v`, err)
	}
	var response GetProductsResponse
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return GetProductsResponse{}, fmt.Errorf(`failed to map get products response, error: %v`, err)
	}
	return response, nil
}
