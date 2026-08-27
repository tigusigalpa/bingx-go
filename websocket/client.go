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

// MessageCallback represents a BingX API component or value.
type MessageCallback func(data map[string]interface{})

// Client manages a connection to a BingX WebSocket endpoint.
type Client struct {
	url       string
	conn      *websocket.Conn
	callbacks []MessageCallback
	running   bool
	mu        sync.RWMutex
	// writeMu serializes WriteMessage calls. gorilla/websocket allows only
	// one concurrent writer per connection. Held independently of mu so a
	// long network write can't pin the state mutex and block reconnects.
	writeMu sync.Mutex
	done    chan struct{}
}

// NewClient creates a WebSocket client for url.
func NewClient(url string) *Client {
	return &Client{
		url:       url,
		callbacks: make([]MessageCallback, 0),
		done:      make(chan struct{}),
	}
}

// Connect performs the Connect operation.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 60 * time.Second,
	}

	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	// Disconnect closes done to wake a listener. A fresh connection needs a
	// fresh signal channel so the client can be reused.
	select {
	case <-c.done:
		c.done = make(chan struct{})
	default:
	}
	return nil
}

// Disconnect performs the Disconnect operation.
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.running = false
	select {
	case <-c.done:
		// Disconnect is intentionally idempotent.
	default:
		close(c.done)
	}

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}

	return nil
}

// Send performs the Send operation.
func (c *Client) Send(message map[string]interface{}) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("WebSocket client is not connected")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Serialize writes via writeMu (not c.mu) so a stalled WriteMessage on
	// a half-broken connection can't pin the state mutex and starve the
	// Disconnect/Connect path that wants to swap c.conn during reconnect.
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, data)
}

// Subscribe performs the Subscribe operation.
func (c *Client) Subscribe(id, dataType string) error {
	return c.Send(map[string]interface{}{
		"id":       id,
		"reqType":  "sub",
		"dataType": dataType,
	})
}

// Unsubscribe performs the Unsubscribe operation.
func (c *Client) Unsubscribe(id, dataType string) error {
	return c.Send(map[string]interface{}{
		"id":       id,
		"reqType":  "unsub",
		"dataType": dataType,
	})
}

// OnMessage performs the OnMessage operation.
func (c *Client) OnMessage(callback MessageCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks = append(c.callbacks, callback)
}

// Listen performs the Listen operation.
func (c *Client) Listen() error {
	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return fmt.Errorf("WebSocket client is not connected")
	}
	c.running = true
	conn := c.conn
	done := c.done
	c.mu.Unlock()

	for c.isRunning() {
		select {
		case <-done:
			return nil
		default:
			messageType, message, err := conn.ReadMessage()
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
					if err := c.Send(map[string]interface{}{"pong": ping}); err != nil {
						return fmt.Errorf("failed to respond to WebSocket ping: %w", err)
					}
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

// decompressMessage performs the decompressMessage operation.
func (c *Client) decompressMessage(message []byte) ([]byte, error) {
	if len(message) >= 2 && message[0] == 0x1f && message[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			return nil, err
		}
		defer func() { _ = reader.Close() }()

		decompressed, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return decompressed, nil
	}

	return message, nil
}

// IsConnected performs the IsConnected operation.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil
}

// Stop performs the Stop operation.
func (c *Client) Stop() {
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()
}

// isRunning performs the isRunning operation.
func (c *Client) isRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}
