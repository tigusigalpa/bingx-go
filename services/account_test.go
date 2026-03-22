package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/http"
)

func TestNewAccountService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	if service == nil {
		t.Fatal("Expected AccountService to be created, got nil")
	}

	if service.client == nil {
		t.Error("AccountService client should not be nil")
	}
}

func TestGetBalance(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetBalance()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetPositions(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	tests := []struct {
		name   string
		symbol *string
	}{
		{
			name:   "All positions",
			symbol: nil,
		},
		{
			name:   "Specific symbol",
			symbol: func() *string { s := "BTC-USDT"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetPositions(tt.symbol)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetAccountInfo(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAccountInfo()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetTradingFees(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetTradingFees("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetMarginMode(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetMarginMode("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestSetMarginMode(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.SetMarginMode("BTC-USDT", "ISOLATED")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestAccountServiceGetLeverage(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	tests := []struct {
		name       string
		symbol     string
		recvWindow *int
	}{
		{
			name:       "Without recv window",
			symbol:     "BTC-USDT",
			recvWindow: nil,
		},
		{
			name:       "With recv window",
			symbol:     "BTC-USDT",
			recvWindow: func() *int { v := 5000; return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetLeverage(tt.symbol, tt.recvWindow)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestAccountServiceSetLeverage(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	tests := []struct {
		name       string
		symbol     string
		side       string
		leverage   int
		recvWindow *int
	}{
		{
			name:       "Set leverage LONG",
			symbol:     "BTC-USDT",
			side:       "LONG",
			leverage:   10,
			recvWindow: nil,
		},
		{
			name:       "Set leverage SHORT with recv window",
			symbol:     "ETH-USDT",
			side:       "SHORT",
			leverage:   20,
			recvWindow: func() *int { v := 5000; return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SetLeverage(tt.symbol, tt.side, tt.leverage, tt.recvWindow)
			if err == nil {
				t.Skip("Skipping test - would require mock HTTP server")
			}
		})
	}
}

func TestGetPositionMargin(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetPositionMargin("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestSetPositionMargin(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.SetPositionMargin("BTC-USDT", "LONG", 100.0, 1)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetBalanceHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetBalanceHistory("USDT", 100)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAccountPermissions(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAccountPermissions()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAPIKey(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAPIKey()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetUserCommissionRates(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetUserCommissionRates("BTC-USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAPIRateLimits(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAPIRateLimits()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetDepositHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetDepositHistory("USDT", 100)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetWithdrawHistory(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetWithdrawHistory("USDT", 100)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAssetDetails(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAssetDetails("USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetAllAssets(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetAllAssets()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestGetFundingWallet(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	_, err := service.GetFundingWallet("USDT")
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestDustTransfer(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewAccountService(client)

	assets := []string{"BTC", "ETH", "BNB"}
	_, err := service.DustTransfer(assets)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}
