package main

import (
	"algotrade/pkg/binance"
	"algotrade/pkg/wss"
	"encoding/json"

	"context"
	"fmt"
	"time"
)

func main() {
	ws := wss.WebSocket{
		Url: "wss://stream.binance.com/stream",
	}
	err := ws.Connect();
	if err != nil {
		fmt.Printf(`failed to connect to websocket: error %v`, err)
	}
	defer ws.Close();

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	ws.SendMessage(`{"method":"SUBSCRIBE","params":["bnbbtc@depth"],"id":1}`);
	startReading(ctx, &ws);
}

func startReading(ctx context.Context, ws *wss.WebSocket) error {
	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
		default:
			message, err := ws.ReadMessage()
			if err != nil {
				return err
			}
			var depthUpdate binance.DepthUpdate
			err = json.Unmarshal([]byte(message), &depthUpdate);
			if err != nil {
				continue
			}
			fmt.Println(depthUpdate.Data.Asks)
		}
	}
}