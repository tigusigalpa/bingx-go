package websocket

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageCallback func(data map[string]interface{})

type WebSocketClient struct {
	url       string
	conn      *websocket.Conn
	callbacks []MessageCallback
	running   bool
	mu        sync.RWMutex
	done      chan struct{}
}

func NewWebSocketClient(url string) *WebSocketClient {
	return &WebSocketClient{
		url:       url,
		callbacks: make([]MessageCallback, 0),
		done:      make(chan struct{}),
	}
}

func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dialer := websocket.Dialer{
		HandshakeTimeout: 60 * time.Second,
	}

	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	return nil
}

func (c *WebSocketClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.running = false
	close(c.done)

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}

	return nil
}

func (c *WebSocketClient) Send(message map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return fmt.Errorf("WebSocket client is not connected")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WebSocketClient) Subscribe(id, dataType string) error {
	return c.Send(map[string]interface{}{
		"id":       id,
		"reqType":  "sub",
		"dataType": dataType,
	})
}

func (c *WebSocketClient) Unsubscribe(id, dataType string) error {
	return c.Send(map[string]interface{}{
		"id":       id,
		"reqType":  "unsub",
		"dataType": dataType,
	})
}

func (c *WebSocketClient) OnMessage(callback MessageCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks = append(c.callbacks, callback)
}

func (c *WebSocketClient) Listen() error {
	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return fmt.Errorf("WebSocket client is not connected")
	}
	c.running = true
	c.mu.Unlock()

	for c.isRunning() {
		select {
		case <-c.done:
			return nil
		default:
			messageType, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return fmt.Errorf("WebSocket connection closed unexpectedly: %w", err)
				}
				return err
			}

			if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
				data, err := c.decompressMessage(message)
				if err != nil {
					continue
				}

				var parsed map[string]interface{}
				if err := json.Unmarshal(data, &parsed); err != nil {
					continue
				}

				if ping, ok := parsed["ping"]; ok {
					c.Send(map[string]interface{}{"pong": ping})
					continue
				}

				c.mu.RLock()
				callbacks := make([]MessageCallback, len(c.callbacks))
				copy(callbacks, c.callbacks)
				c.mu.RUnlock()

				for _, callback := range callbacks {
					callback(parsed)
				}
			}
		}
	}

	return nil
}

func (c *WebSocketClient) decompressMessage(message []byte) ([]byte, error) {
	if len(message) >= 2 && message[0] == 0x1f && message[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		decompressed, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return decompressed, nil
	}

	return message, nil
}

func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil
}

func (c *WebSocketClient) Stop() {
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()
}

func (c *WebSocketClient) isRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}
