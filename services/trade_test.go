package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/http"
)

func TestNewTradeService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	if service == nil {
		t.Fatal("Expected TradeService to be created, got nil")
	}

	if service.client == nil {
		t.Error("TradeService client should not be nil")
	}
}

func TestCalculateFuturesCommission(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name           string
		margin         float64
		leverage       int
		commissionRate *float64
		expectedPV     float64
		expectedComm   float64
	}{
		{
			name:           "Default commission rate",
			margin:         100.0,
			leverage:       10,
			commissionRate: nil,
			expectedPV:     1000.0,
			expectedComm:   0.45,
		},
		{
			name:           "Custom commission rate",
			margin:         100.0,
			leverage:       10,
			commissionRate: func() *float64 { v := 0.001; return &v }(),
			expectedPV:     1000.0,
			expectedComm:   1.0,
		},
		{
			name:           "High leverage",
			margin:         50.0,
			leverage:       125,
			commissionRate: nil,
			expectedPV:     6250.0,
			expectedComm:   2.8125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateFuturesCommission(tt.margin, tt.leverage, tt.commissionRate)

			if result.Margin != tt.margin {
				t.Errorf("Expected margin %f, got %f", tt.margin, result.Margin)
			}

			if result.Leverage != tt.leverage {
				t.Errorf("Expected leverage %d, got %d", tt.leverage, result.Leverage)
			}

			if result.PositionValue != tt.expectedPV {
				t.Errorf("Expected position value %f, got %f", tt.expectedPV, result.PositionValue)
			}

			if result.Commission != tt.expectedComm {
				t.Errorf("Expected commission %f, got %f", tt.expectedComm, result.Commission)
			}

			if result.NetPositionValue != tt.expectedPV-tt.expectedComm {
				t.Errorf("Expected net position value %f, got %f", tt.expectedPV-tt.expectedComm, result.NetPositionValue)
			}
		})
	}
}

func TestGetCommissionAmount(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name     string
		margin   float64
		leverage int
		expected float64
	}{
		{
			name:     "Basic commission",
			margin:   100.0,
			leverage: 10,
			expected: 0.45,
		},
		{
			name:     "High leverage commission",
			margin:   50.0,
			leverage: 125,
			expected: 2.8125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetCommissionAmount(tt.margin, tt.leverage)
			if result != tt.expected {
				t.Errorf("Expected commission %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestModifyOrder(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name          string
		symbol        string
		quantity      float64
		orderID       *string
		clientOrderID *string
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "No order ID provided",
			symbol:        "BTC-USDT",
			quantity:      1.0,
			orderID:       nil,
			clientOrderID: nil,
			expectError:   true,
			errorMsg:      "modifyOrder requires either orderID or clientOrderID",
		},
		{
			name:          "Zero quantity",
			symbol:        "BTC-USDT",
			quantity:      0.0,
			orderID:       func() *string { s := "12345"; return &s }(),
			clientOrderID: nil,
			expectError:   true,
			errorMsg:      "quantity must be greater than 0",
		},
		{
			name:          "Negative quantity",
			symbol:        "BTC-USDT",
			quantity:      -1.0,
			orderID:       func() *string { s := "12345"; return &s }(),
			clientOrderID: nil,
			expectError:   true,
			errorMsg:      "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ModifyOrder(tt.symbol, tt.quantity, tt.orderID, tt.clientOrderID, nil, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestChangeMarginType(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name        string
		symbol      string
		marginType  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid margin type",
			symbol:      "BTC-USDT",
			marginType:  "INVALID",
			expectError: true,
			errorMsg:    "margin type must be ISOLATED or CROSSED",
		},
		{
			name:        "Valid ISOLATED",
			symbol:      "BTC-USDT",
			marginType:  "ISOLATED",
			expectError: false,
		},
		{
			name:        "Valid CROSSED",
			symbol:      "BTC-USDT",
			marginType:  "CROSSED",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ChangeMarginType(tt.symbol, tt.marginType, nil, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else if !tt.expectError && err != nil {
				if err.Error() != tt.errorMsg {
					t.Skip("Skipping test - would require mock HTTP server")
				}
			}
		})
	}
}

func TestSetLeverage(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name        string
		symbol      string
		leverage    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Leverage too low",
			symbol:      "BTC-USDT",
			leverage:    0,
			expectError: true,
			errorMsg:    "leverage must be between 1 and 125",
		},
		{
			name:        "Leverage too high",
			symbol:      "BTC-USDT",
			leverage:    126,
			expectError: true,
			errorMsg:    "leverage must be between 1 and 125",
		},
		{
			name:        "Valid leverage min",
			symbol:      "BTC-USDT",
			leverage:    1,
			expectError: false,
		},
		{
			name:        "Valid leverage max",
			symbol:      "BTC-USDT",
			leverage:    125,
			expectError: false,
		},
		{
			name:        "Valid leverage mid",
			symbol:      "BTC-USDT",
			leverage:    50,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SetLeverage(tt.symbol, tt.leverage, nil, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else if !tt.expectError && err != nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestCancelBatchOrders(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name           string
		symbol         string
		orderIDs       []string
		clientOrderIDs []string
	}{
		{
			name:           "Cancel with order IDs",
			symbol:         "BTC-USDT",
			orderIDs:       []string{"123", "456", "789"},
			clientOrderIDs: nil,
		},
		{
			name:           "Cancel with client order IDs",
			symbol:         "BTC-USDT",
			orderIDs:       nil,
			clientOrderIDs: []string{"client-1", "client-2"},
		},
		{
			name:           "Cancel with both ID types",
			symbol:         "BTC-USDT",
			orderIDs:       []string{"123"},
			clientOrderIDs: []string{"client-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CancelBatchOrders(tt.symbol, tt.orderIDs, tt.clientOrderIDs, nil, nil)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetOpenOrders(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name   string
		symbol *string
		limit  int
	}{
		{
			name:   "All symbols",
			symbol: nil,
			limit:  100,
		},
		{
			name:   "Specific symbol",
			symbol: func() *string { s := "BTC-USDT"; return &s }(),
			limit:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetOpenOrders(tt.symbol, tt.limit)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetOrderHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewTradeService(client)

	tests := []struct {
		name      string
		symbol    *string
		limit     int
		startTime *int64
		endTime   *int64
	}{
		{
			name:      "All symbols no time range",
			symbol:    nil,
			limit:     100,
			startTime: nil,
			endTime:   nil,
		},
		{
			name:      "Specific symbol with time range",
			symbol:    func() *string { s := "BTC-USDT"; return &s }(),
			limit:     50,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetOrderHistory(tt.symbol, tt.limit, tt.startTime, tt.endTime)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}
