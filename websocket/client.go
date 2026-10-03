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

const maxWebSocketMessageSize = 16 << 20

// MessageCallback represents a BingX API component or value.
type MessageCallback func(data map[string]interface{})

// RawMessageHandler receives a copy of the decompressed payload before JSON decoding.
type RawMessageHandler func(payload []byte, receivedAt time.Time)

// Client manages a connection to a BingX WebSocket endpoint.
type Client struct {
	url         string
	conn        *websocket.Conn
	callbacks   []MessageCallback
	rawHandlers []RawMessageHandler
	running     bool
	mu          sync.RWMutex
	// writeMu serializes WriteMessage calls. gorilla/websocket allows only
	// one concurrent writer per connection. Held independently of mu so a
	// long network write can't pin the state mutex and block reconnects.
	writeMu sync.Mutex
	done    chan struct{}
}

// NewClient creates a WebSocket client for url.
func NewClient(url string) *Client {
	return &Client{
		url:         url,
		callbacks:   make([]MessageCallback, 0),
		rawHandlers: make([]RawMessageHandler, 0),
		done:        make(chan struct{}),
	}
}

// OnRawMessage registers a handler for decompressed text and binary messages.
// The handler receives its own immutable copy of the payload and must return promptly.
func (c *Client) OnRawMessage(handler RawMessageHandler) {
	if handler == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.rawHandlers = append(c.rawHandlers, handler)
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
	conn.SetReadLimit(maxWebSocketMessageSize)

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
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("WebSocket client is already listening")
	}
	c.running = true
	conn := c.conn
	done := c.done
	c.mu.Unlock()
	defer c.finishListen(conn)

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
					return fmt.Errorf("failed to decompress WebSocket message: %w", err)
				}

				receivedAt := time.Now()
				if bytes.Equal(bytes.TrimSpace(data), []byte("Ping")) {
					if err := c.sendText([]byte("Pong")); err != nil {
						return fmt.Errorf("failed to respond to WebSocket ping: %w", err)
					}
					continue
				}

				c.dispatchRawMessage(data, receivedAt)

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

func (c *Client) sendText(message []byte) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return fmt.Errorf("WebSocket client is not connected")
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) dispatchRawMessage(data []byte, receivedAt time.Time) {
	c.mu.RLock()
	handlers := make([]RawMessageHandler, len(c.rawHandlers))
	copy(handlers, c.rawHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		handler(append([]byte(nil), data...), receivedAt)
	}
}

// decompressMessage performs the decompressMessage operation.
func (c *Client) decompressMessage(message []byte) ([]byte, error) {
	if len(message) >= 2 && message[0] == 0x1f && message[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			return nil, err
		}
		defer func() { _ = reader.Close() }()

		decompressed, err := io.ReadAll(io.LimitReader(reader, maxWebSocketMessageSize+1))
		if err != nil {
			return nil, err
		}
		if len(decompressed) > maxWebSocketMessageSize {
			return nil, fmt.Errorf("decompressed WebSocket message exceeds %d bytes", maxWebSocketMessageSize)
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

func (c *Client) finishListen(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == conn {
		c.running = false
	}
}

// isRunning performs the isRunning operation.
func (c *Client) isRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}
