package wss

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Cookies      []*http.Cookie
	Url          string
	Origin       string
	conn         *websocket.Conn
}

func (w *WebSocket) Connect() error {
	jar, err := w.getCookieJar()
	if err != nil {
		return fmt.Errorf("failed to get cookie jar | %w", err)
	}

	tlsConfig, err := getTlsConfig()
	if err != nil {
		return fmt.Errorf("failed to get tls config | %w", err)
	}

	dialer := websocket.Dialer{
		Jar:              jar,
		TLSClientConfig:  tlsConfig,
		HandshakeTimeout: 15 * time.Second,
	}

	conn, _, err := dialer.Dial(w.Url, w.getHeaders())
	if err != nil {
		for i := 0; i < 5; i++ {
			fmt.Printf(`failed to connect to websocket, error : %v. Retrying...`, err)
			time.Sleep(time.Second)
			conn, _, err = dialer.Dial(w.Url, w.getHeaders())
			if err == nil {
				break;
			}
		}
		if err != nil {
			return fmt.Errorf("failed to connect to websocket | %w", err)
		}
	}
	fmt.Println("Connected to websocket!")
	w.conn = conn
	return nil
}

func (w *WebSocket) ReadMessage() (string, error) {
	for {
		messageType, message, err := w.conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("error reading WebSocket message: %w", err)
		}
		if messageType == websocket.PingMessage {
			if err := w.SendPong(message); err != nil {
				return "", fmt.Errorf("error responding to PING frame: %w", err)
			}
			continue
		}
		if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
			return string(message), nil
		}
	}
}

func (w *WebSocket) SendMessage(message string) error {
	return w.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (w *WebSocket) Close() {
	w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	w.conn.Close()
}

func (w *WebSocket) SendPing(ping string) error {
	return w.conn.WriteMessage(websocket.TextMessage, []byte(ping))
}

func (w *WebSocket) SendPong(payload []byte) error {
	return w.conn.WriteControl(websocket.PongMessage, payload, time.Now().Add(time.Second))
}
