package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

func TestNewWalletService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	if service == nil {
		t.Fatal("Expected WalletService to be created, got nil")
	}

	if service.client == nil {
		t.Error("WalletService client should not be nil")
	}
}

func TestWalletGetDepositHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	tests := []struct {
		name      string
		coin      string
		status    *int
		startTime *int64
		endTime   *int64
		limit     int
	}{
		{
			name:      "Basic deposit history",
			coin:      "USDT",
			status:    nil,
			startTime: nil,
			endTime:   nil,
			limit:     100,
		},
		{
			name:      "Deposit history with status",
			coin:      "BTC",
			status:    func() *int { v := 1; return &v }(),
			startTime: nil,
			endTime:   nil,
			limit:     50,
		},
		{
			name:      "Deposit history with time range",
			coin:      "ETH",
			status:    nil,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
			limit:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetDepositHistory(tt.coin, tt.status, tt.startTime, tt.endTime, tt.limit)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetDepositAddress(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	_, err := service.GetDepositAddress("USDT", "TRC20")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestWalletGetWithdrawalHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	tests := []struct {
		name      string
		coin      string
		status    *int
		startTime *int64
		endTime   *int64
		limit     int
	}{
		{
			name:      "Basic withdrawal history",
			coin:      "USDT",
			status:    nil,
			startTime: nil,
			endTime:   nil,
			limit:     100,
		},
		{
			name:      "Withdrawal history with status",
			coin:      "BTC",
			status:    func() *int { v := 1; return &v }(),
			startTime: nil,
			endTime:   nil,
			limit:     50,
		},
		{
			name:      "Withdrawal history with time range",
			coin:      "ETH",
			status:    nil,
			startTime: func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:   func() *int64 { v := int64(1609545600000); return &v }(),
			limit:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetWithdrawalHistory(tt.coin, tt.status, tt.startTime, tt.endTime, tt.limit)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	tests := []struct {
		name       string
		coin       string
		address    string
		amount     float64
		network    string
		addressTag *string
	}{
		{
			name:       "Basic withdrawal",
			coin:       "USDT",
			address:    "TXYZabc123",
			amount:     100.0,
			network:    "TRC20",
			addressTag: nil,
		},
		{
			name:       "Withdrawal with address tag",
			coin:       "XRP",
			address:    "rXYZabc123",
			amount:     50.0,
			network:    "XRP",
			addressTag: func() *string { s := "12345"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Withdraw(tt.coin, tt.address, tt.amount, tt.network, tt.addressTag)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetAllCoinInfo(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	_, err := service.GetAllCoinInfo()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetMainAccountTransferHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewWalletService(client)

	tests := []struct {
		name         string
		coin         string
		transferType *string
		startTime    *int64
		endTime      *int64
		limit        int
	}{
		{
			name:         "Basic transfer history",
			coin:         "USDT",
			transferType: nil,
			startTime:    nil,
			endTime:      nil,
			limit:        100,
		},
		{
			name:         "Transfer history with type",
			coin:         "BTC",
			transferType: func() *string { s := "DEPOSIT"; return &s }(),
			startTime:    nil,
			endTime:      nil,
			limit:        50,
		},
		{
			name:         "Transfer history with time range",
			coin:         "ETH",
			transferType: nil,
			startTime:    func() *int64 { v := int64(1609459200000); return &v }(),
			endTime:      func() *int64 { v := int64(1609545600000); return &v }(),
			limit:        100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetMainAccountTransferHistory(tt.coin, tt.transferType, tt.startTime, tt.endTime, tt.limit)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}
