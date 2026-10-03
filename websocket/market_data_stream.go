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

// SubscribeDepthAt subscribes to a depth stream at one of BingX's supported intervals.
func (m *MarketDataStream) SubscribeDepthAt(symbol string, levels int, interval time.Duration, id ...string) error {
	if err := validateDepthInterval(symbol, interval); err != nil {
		return err
	}
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@depth%d@%dms", symbol, levels, interval.Milliseconds()))
}

// SubscribeIncrementalDepth subscribes to incremental order-book updates.
func (m *MarketDataStream) SubscribeIncrementalDepth(symbol string, id ...string) error {
	return m.Subscribe(m.generateID(id...), fmt.Sprintf("%s@incrDepth", symbol))
}

// SubscribeLastPrice subscribes to last-price updates.
func (m *MarketDataStream) SubscribeLastPrice(symbol string, id ...string) error {
	return m.Subscribe(m.generateID(id...), fmt.Sprintf("%s@lastPrice", symbol))
}

// SubscribeMarkPrice subscribes to mark-price updates.
func (m *MarketDataStream) SubscribeMarkPrice(symbol string, id ...string) error {
	return m.Subscribe(m.generateID(id...), fmt.Sprintf("%s@markPrice", symbol))
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

// UnsubscribeDepthAt unsubscribes from a depth stream at one of BingX's supported intervals.
func (m *MarketDataStream) UnsubscribeDepthAt(symbol string, levels int, interval time.Duration, id ...string) error {
	if err := validateDepthInterval(symbol, interval); err != nil {
		return err
	}
	return m.Unsubscribe(m.generateID(id...), fmt.Sprintf("%s@depth%d@%dms", symbol, levels, interval.Milliseconds()))
}

func validateDepthInterval(symbol string, interval time.Duration) error {
	if interval != 200*time.Millisecond && interval != 500*time.Millisecond {
		return fmt.Errorf("unsupported depth interval %s: use 200ms or 500ms", interval)
	}
	if interval == 200*time.Millisecond && symbol != "BTC-USDT" && symbol != "ETH-USDT" {
		return fmt.Errorf("unsupported 200ms depth interval for %s: use 500ms", symbol)
	}
	return nil
}

// UnsubscribeIncrementalDepth unsubscribes from incremental order-book updates.
func (m *MarketDataStream) UnsubscribeIncrementalDepth(symbol string, id ...string) error {
	return m.Unsubscribe(m.generateID(id...), fmt.Sprintf("%s@incrDepth", symbol))
}

// UnsubscribeLastPrice unsubscribes from last-price updates.
func (m *MarketDataStream) UnsubscribeLastPrice(symbol string, id ...string) error {
	return m.Unsubscribe(m.generateID(id...), fmt.Sprintf("%s@lastPrice", symbol))
}

// UnsubscribeMarkPrice unsubscribes from mark-price updates.
func (m *MarketDataStream) UnsubscribeMarkPrice(symbol string, id ...string) error {
	return m.Unsubscribe(m.generateID(id...), fmt.Sprintf("%s@markPrice", symbol))
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
