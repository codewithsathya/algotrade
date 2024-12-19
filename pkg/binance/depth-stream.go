package binance

import (
	"algotrade/pkg/wss"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type DepthStream struct {
	subscribedSymbols []string
	ws                *wss.WebSocket
}

func NewDepthStream(subscribedSymbols []string) (*DepthStream, error) {
	if len(subscribedSymbols) == 0 {
		return nil, fmt.Errorf("no symbols provided for subscription")
	}

	return &DepthStream{
		subscribedSymbols: subscribedSymbols,
		ws: &wss.WebSocket{
			Url: "wss://stream.binance.com/stream",
		},
	}, nil
}

func (d *DepthStream) Start(ctx context.Context) error {
	if err := d.ws.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Binance stream: %w", err)
	}

	defer d.ws.Close()

	payload, err := d.constructPayload()
	if err != nil {
		return fmt.Errorf("failed to construct payload: %w", err)
	}

	if err := d.ws.SendMessage(payload); err != nil {
		return fmt.Errorf("failed to send subscription payload: %w", err)
	}

	return d.startReading(ctx)
}

func (d *DepthStream) startReading(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := d.ws.ReadMessage()
			if err != nil {
				return fmt.Errorf("error reading WebSocket message: %w", err)
			}
			var depthUpdate DepthUpdate
			if err := json.Unmarshal([]byte(message), &depthUpdate); err != nil {
				fmt.Printf("failed to unmarshal depth update: %v\n", err)
				continue
			}
			fmt.Println(depthUpdate.Data.Symbol)
			if(len(depthUpdate.Data.Asks) > 0) {
				fmt.Println(`Buy Price`, depthUpdate.Data.Asks[0])
			}
			if(len(depthUpdate.Data.Bids) > 0) {
				fmt.Println(`Sell Price`, depthUpdate.Data.Bids[0])
			}
		}
	}
}

func (d *DepthStream) Stop() {
	if d.ws != nil {
		d.ws.Close()
	}
}

func (d *DepthStream) constructPayload() (string, error) {
	params := make([]string, len(d.subscribedSymbols))
	for i, symbol := range d.subscribedSymbols {
		params[i] = strings.ToLower(symbol) + "@depth"
	}

	payload := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": params,
		"id":     1,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling subscription payload: %w", err)
	}

	return string(payloadJSON), nil
}
