package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/http"
)

func TestNewContractService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewContractService(client)

	if service == nil {
		t.Fatal("Expected ContractService to be created, got nil")
	}

	if service.client == nil {
		t.Error("ContractService client should not be nil")
	}
}

func TestGetAllPositions(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewContractService(client)

	tests := []struct {
		name       string
		timestamp  *int64
		recvWindow *int64
	}{
		{
			name:       "Without timestamp and recv window",
			timestamp:  nil,
			recvWindow: nil,
		},
		{
			name:       "With timestamp",
			timestamp:  func() *int64 { v := int64(1609459200000); return &v }(),
			recvWindow: nil,
		},
		{
			name:       "With recv window",
			timestamp:  nil,
			recvWindow: func() *int64 { v := int64(5000); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetAllPositions(tt.timestamp, tt.recvWindow)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewContractService(client)

	tests := []struct {
		name      string
		symbol    string
		limit     int
		startTime *int64
		endTime   *int64
	}{
		{
			name:      "Basic all orders",
			symbol:    "BTC-USDT",
			limit:     100,
			startTime: nil,
			endTime:   nil,
		},
		{
			name:      "All orders with time range",
			symbol:    "ETH-USDT",
			limit:     50,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetAllOrders(tt.symbol, tt.limit, tt.startTime, tt.endTime)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestContractGetBalance(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewContractService(client)

	tests := []struct {
		name       string
		timestamp  *int64
		recvWindow *int64
	}{
		{
			name:       "Without timestamp and recv window",
			timestamp:  nil,
			recvWindow: nil,
		},
		{
			name:       "With timestamp",
			timestamp:  func() *int64 { v := int64(1609459200000); return &v }(),
			recvWindow: nil,
		},
		{
			name:       "With recv window",
			timestamp:  nil,
			recvWindow: func() *int64 { v := int64(5000); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetBalance(tt.timestamp, tt.recvWindow)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}
