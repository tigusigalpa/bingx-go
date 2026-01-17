package websocket

import (
	"fmt"
	"time"
)

const MarketDataStreamURL = "wss://open-api-swap.bingx.com/swap-market"

type MarketDataStream struct {
	*WebSocketClient
}

func NewMarketDataStream() *MarketDataStream {
	return &MarketDataStream{
		WebSocketClient: NewWebSocketClient(MarketDataStreamURL),
	}
}

func (m *MarketDataStream) SubscribeTrade(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@trade", symbol))
}

func (m *MarketDataStream) SubscribeKline(symbol, interval string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@kline_%s", symbol, interval))
}

func (m *MarketDataStream) SubscribeDepth(symbol string, levels int, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@depth%d", symbol, levels))
}

func (m *MarketDataStream) SubscribeTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@ticker", symbol))
}

func (m *MarketDataStream) SubscribeBookTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Subscribe(requestID, fmt.Sprintf("%s@bookTicker", symbol))
}

func (m *MarketDataStream) UnsubscribeTrade(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@trade", symbol))
}

func (m *MarketDataStream) UnsubscribeKline(symbol, interval string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@kline_%s", symbol, interval))
}

func (m *MarketDataStream) UnsubscribeDepth(symbol string, levels int, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@depth%d", symbol, levels))
}

func (m *MarketDataStream) UnsubscribeTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@ticker", symbol))
}

func (m *MarketDataStream) UnsubscribeBookTicker(symbol string, id ...string) error {
	requestID := m.generateID(id...)
	return m.Unsubscribe(requestID, fmt.Sprintf("%s@bookTicker", symbol))
}

func (m *MarketDataStream) generateID(id ...string) string {
	if len(id) > 0 && id[0] != "" {
		return id[0]
	}
	return fmt.Sprintf("bingx_%d", time.Now().UnixNano())
}
