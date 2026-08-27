package websocket

import (
	"fmt"
	"time"
)

// MarketDataStreamURL is the WebSocket endpoint for market data streams (v3 compatible)
const MarketDataStreamURL = "wss://open-api-swap.bingx.com/swap-market"

// MarketDataStream represents a BingX API component or value.
type MarketDataStream struct {
	*Client
}

// NewMarketDataStream creates a new client or service instance.
func NewMarketDataStream() *MarketDataStream {
	return &MarketDataStream{
		Client: NewClient(MarketDataStreamURL),
	}
}

// SubscribeTrade performs the SubscribeTrade operation.
func (m *MarketDataStream) SubscribeTrade(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@trade", symbol))
}

// SubscribeKline performs the SubscribeKline operation.
func (m *MarketDataStream) SubscribeKline(symbol, interval string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@kline_%s", symbol, interval))
}

// SubscribeDepth performs the SubscribeDepth operation.
func (m *MarketDataStream) SubscribeDepth(symbol string, levels int, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@depth%d", symbol, levels))
}

// SubscribeTicker performs the SubscribeTicker operation.
func (m *MarketDataStream) SubscribeTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@ticker", symbol))
}

// SubscribeBookTicker performs the SubscribeBookTicker operation.
func (m *MarketDataStream) SubscribeBookTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@bookTicker", symbol))
}

// UnsubscribeTrade performs the UnsubscribeTrade operation.
func (m *MarketDataStream) UnsubscribeTrade(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@trade", symbol))
}

// UnsubscribeKline performs the UnsubscribeKline operation.
func (m *MarketDataStream) UnsubscribeKline(symbol, interval string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@kline_%s", symbol, interval))
}

// UnsubscribeDepth performs the UnsubscribeDepth operation.
func (m *MarketDataStream) UnsubscribeDepth(symbol string, levels int, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@depth%d", symbol, levels))
}

// UnsubscribeTicker performs the UnsubscribeTicker operation.
func (m *MarketDataStream) UnsubscribeTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@ticker", symbol))
}

// UnsubscribeBookTicker performs the UnsubscribeBookTicker operation.
func (m *MarketDataStream) UnsubscribeBookTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@bookTicker", symbol))
}

// generateID performs the generateID operation.
func (m *MarketDataStream) generateID(id ...string) string {
	if len(id) > 0 && id[0] != "" {
		return id[0]
	}
	return fmt.Sprintf("bingx_%d", time.Now().UnixNano())
}
