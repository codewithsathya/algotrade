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
	updates           chan DepthUpdate
}

func NewDepthStream(subscribedSymbols []string) (*DepthStream) {
	return &DepthStream{
		subscribedSymbols: subscribedSymbols,
		ws: &wss.WebSocket{
			Url: "wss://stream.binance.com/stream",
		},
		updates: make(chan DepthUpdate),
	}
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

	if err := d.startReading(ctx); err != nil {
		fmt.Printf("error in depth stream reading: %v\n", err)
		return err
	}

	return nil
}

func (d *DepthStream) GetUpdatesChannel() <-chan DepthUpdate {
	return d.updates
}

func (d *DepthStream) Stop() {
	if d.ws != nil {
		d.ws.Close()
	}
	close(d.updates)
}

func (d *DepthStream) constructPayload() (string, error) {
	params := make([]string, len(d.subscribedSymbols))
	for i, symbol := range d.subscribedSymbols {
		params[i] = strings.ToLower(symbol) + "@depth5"
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

func (d *DepthStream) startReading(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := d.ws.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("error reading WebSocket message: %w", err)
			}
			var depthUpdate DepthUpdate
			if err := json.Unmarshal([]byte(message), &depthUpdate); err != nil {
				fmt.Printf("failed to unmarshal depth update: %v\n", err)
				continue
			}

			select {
			case d.updates <- depthUpdate:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
