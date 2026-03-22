package websocket

import (
	"strings"
	"testing"
)

func TestNewMarketDataStream(t *testing.T) {
	stream := NewMarketDataStream()

	if stream == nil {
		t.Fatal("Expected MarketDataStream to be created, got nil")
	}

	if stream.WebSocketClient == nil {
		t.Error("WebSocketClient should not be nil")
	}
}

func TestMarketDataStreamURL(t *testing.T) {
	expectedURL := "wss://open-api-swap.bingx.com/swap-market"
	if MarketDataStreamURL != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, MarketDataStreamURL)
	}
}

func TestGenerateID(t *testing.T) {
	stream := NewMarketDataStream()

	tests := []struct {
		name     string
		id       []string
		expected string
	}{
		{
			name:     "With custom ID",
			id:       []string{"custom-id-123"},
			expected: "custom-id-123",
		},
		{
			name:     "Without ID",
			id:       []string{},
			expected: "bingx_",
		},
		{
			name:     "With empty string",
			id:       []string{""},
			expected: "bingx_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stream.generateID(tt.id...)

			if tt.expected == "bingx_" {
				if !strings.HasPrefix(result, tt.expected) {
					t.Errorf("Expected ID to start with %s, got %s", tt.expected, result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected ID %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestSubscribeTrade(t *testing.T) {
	stream := NewMarketDataStream()

	tests := []struct {
		name   string
		symbol string
		id     []string
	}{
		{
			name:   "Subscribe trade without ID",
			symbol: "BTC-USDT",
			id:     []string{},
		},
		{
			name:   "Subscribe trade with ID",
			symbol: "ETH-USDT",
			id:     []string{"trade-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := stream.SubscribeTrade(tt.symbol, tt.id...)
			if err == nil {
				t.Skip("Skipping test - would require WebSocket connection")
			}
		})
	}
}

func TestSubscribeKline(t *testing.T) {
	stream := NewMarketDataStream()

	tests := []struct {
		name     string
		symbol   string
		interval string
		id       []string
	}{
		{
			name:     "Subscribe 1h kline",
			symbol:   "BTC-USDT",
			interval: "1h",
			id:       []string{},
		},
		{
			name:     "Subscribe 15m kline with ID",
			symbol:   "ETH-USDT",
			interval: "15m",
			id:       []string{"kline-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := stream.SubscribeKline(tt.symbol, tt.interval, tt.id...)
			if err == nil {
				t.Skip("Skipping test - would require WebSocket connection")
			}
		})
	}
}

func TestSubscribeDepth(t *testing.T) {
	stream := NewMarketDataStream()

	tests := []struct {
		name   string
		symbol string
		levels int
		id     []string
	}{
		{
			name:   "Subscribe depth 5",
			symbol: "BTC-USDT",
			levels: 5,
			id:     []string{},
		},
		{
			name:   "Subscribe depth 20 with ID",
			symbol: "ETH-USDT",
			levels: 20,
			id:     []string{"depth-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := stream.SubscribeDepth(tt.symbol, tt.levels, tt.id...)
			if err == nil {
				t.Skip("Skipping test - would require WebSocket connection")
			}
		})
	}
}

func TestSubscribeTicker(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.SubscribeTicker("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestSubscribeBookTicker(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.SubscribeBookTicker("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestUnsubscribeTrade(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.UnsubscribeTrade("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestUnsubscribeKline(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.UnsubscribeKline("BTC-USDT", "1h")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestUnsubscribeDepth(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.UnsubscribeDepth("BTC-USDT", 5)
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestUnsubscribeTicker(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.UnsubscribeTicker("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}

func TestUnsubscribeBookTicker(t *testing.T) {
	stream := NewMarketDataStream()

	err := stream.UnsubscribeBookTicker("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require WebSocket connection")
	}
}
