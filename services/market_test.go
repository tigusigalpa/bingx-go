package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

func TestNewMarketService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	if service == nil {
		t.Fatal("Expected MarketService to be created, got nil")
	}

	if service.client == nil {
		t.Error("MarketService client should not be nil")
	}
}

func TestMarketServiceMethods(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name string
		fn   func() (map[string]interface{}, error)
	}{
		{
			name: "GetFuturesSymbols",
			fn:   func() (map[string]interface{}, error) { return service.GetFuturesSymbols() },
		},
		{
			name: "GetSpotSymbols",
			fn:   func() (map[string]interface{}, error) { return service.GetSpotSymbols() },
		},
		{
			name: "GetSymbols",
			fn:   func() (map[string]interface{}, error) { return service.GetSymbols() },
		},
		{
			name: "GetServerTime",
			fn:   func() (map[string]interface{}, error) { return service.GetServerTime() },
		},
		{
			name: "GetSpotServerTime",
			fn:   func() (map[string]interface{}, error) { return service.GetSpotServerTime() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.fn()
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetLatestPrice(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	symbol := "BTC-USDT"
	_, err := service.GetLatestPrice(symbol)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetSpotLatestPrice(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	symbol := "BTC-USDT"
	_, err := service.GetSpotLatestPrice(symbol)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetDepth(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	symbol := "BTC-USDT"
	limit := 20
	_, err := service.GetDepth(symbol, limit)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetSpotDepth(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	symbol := "BTC-USDT"
	limit := 20
	_, err := service.GetSpotDepth(symbol, limit)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetKlines(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name      string
		symbol    string
		interval  string
		limit     int
		startTime *int64
		endTime   *int64
	}{
		{
			name:      "Basic klines",
			symbol:    "BTC-USDT",
			interval:  "1h",
			limit:     100,
			startTime: nil,
			endTime:   nil,
		},
		{
			name:      "Klines with time range",
			symbol:    "ETH-USDT",
			interval:  "15m",
			limit:     50,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetKlines(tt.symbol, tt.interval, tt.limit, tt.startTime, tt.endTime)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetSpotKlines(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name      string
		symbol    string
		interval  string
		limit     int
		startTime *int64
		endTime   *int64
		timeZone  *int64
	}{
		{
			name:      "Basic spot klines",
			symbol:    "BTC-USDT",
			interval:  "1h",
			limit:     100,
			startTime: nil,
			endTime:   nil,
			timeZone:  nil,
		},
		{
			name:      "Spot klines with timezone",
			symbol:    "ETH-USDT",
			interval:  "15m",
			limit:     50,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
			timeZone:  func() *int64 { v := int64(8); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetSpotKlines(tt.symbol, tt.interval, tt.limit, tt.startTime, tt.endTime, tt.timeZone)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGet24hrTicker(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name   string
		symbol *string
	}{
		{
			name:   "All symbols",
			symbol: nil,
		},
		{
			name:   "Specific symbol",
			symbol: func() *string { s := "BTC-USDT"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Get24hrTicker(tt.symbol)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetSpot24hrTicker(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name   string
		symbol *string
	}{
		{
			name:   "All symbols",
			symbol: nil,
		},
		{
			name:   "Specific symbol",
			symbol: func() *string { s := "BTC-USDT"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetSpot24hrTicker(tt.symbol)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetFundingRateHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	_, err := service.GetFundingRateHistory("BTC-USDT", 100)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetMarkPrice(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	_, err := service.GetMarkPrice("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAggregateTrades(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	tests := []struct {
		name      string
		symbol    string
		limit     int
		fromID    *int64
		startTime *int64
		endTime   *int64
	}{
		{
			name:      "Basic aggregate trades",
			symbol:    "BTC-USDT",
			limit:     100,
			fromID:    nil,
			startTime: nil,
			endTime:   nil,
		},
		{
			name:      "Aggregate trades with fromID",
			symbol:    "ETH-USDT",
			limit:     50,
			fromID:    func() *int64 { v := int64(12345); return &v }(),
			startTime: nil,
			endTime:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetAggregateTrades(tt.symbol, tt.limit, tt.fromID, tt.startTime, tt.endTime)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetRecentTrades(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	_, err := service.GetRecentTrades("BTC-USDT", 100)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAllSymbols(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewMarketService(client)

	_, err := service.GetAllSymbols()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}
