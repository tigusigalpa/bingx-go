package tradfi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func newTestClient(t *testing.T, method, path string, want map[string]string) (*bxhttp.BaseHTTPClient, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != path {
			t.Errorf("path = %s, want %s", r.URL.Path, path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		for key, value := range want {
			if got := r.Form.Get(key); got != value {
				t.Errorf("%s = %q, want %q", key, got, value)
			}
		}
		if r.Form.Get("timestamp") == "" || r.Form.Get("signature") == "" {
			t.Error("signed request is missing timestamp or signature")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"code":0,"data":{}}`)
	}))
	return bxhttp.NewBaseHTTPClient("key", "secret", srv.URL, "", "hex"), srv
}

func TestMarketServiceRoutes(t *testing.T) {
	tests := []struct {
		name, path string
		params     map[string]string
		call       func(*MarketService) error
	}{
		{"stock symbols", "/openApi/swap/v2/quote/contracts", map[string]string{"assetType": "STOCK"}, func(s *MarketService) error { _, err := s.GetStockSymbols(); return err }},
		{"ticker", "/openApi/swap/v2/quote/ticker", map[string]string{"symbol": "TSLA-USDT"}, func(s *MarketService) error { _, err := s.GetTicker("TSLA-USDT"); return err }},
		{"klines", "/openApi/swap/v3/quote/klines", map[string]string{"symbol": "TSLA-USDT", "interval": "1h", "limit": "50"}, func(s *MarketService) error { _, err := s.GetKlines("TSLA-USDT", "1h", 50, nil, nil); return err }},
		{"trading rules", "/openApi/swap/v1/tradingRules", map[string]string{"symbol": "TSLA-USDT"}, func(s *MarketService) error { _, err := s.GetTradingRules("TSLA-USDT"); return err }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, srv := newTestClient(t, http.MethodGet, tt.path, tt.params)
			defer srv.Close()
			if err := tt.call(NewMarketService(client)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAccountAndTradeServiceRoutes(t *testing.T) {
	tests := []struct {
		name, method, path string
		params             map[string]string
		call               func(*bxhttp.BaseHTTPClient) error
	}{
		{"account position mode", http.MethodPost, "/openApi/swap/v2/user/positionMode", map[string]string{"positionMode": "HEDGE"}, func(c *bxhttp.BaseHTTPClient) error { _, err := NewAccountService(c).SetPositionMode(true); return err }},
		{"account income", http.MethodGet, "/openApi/swap/v2/user/income", map[string]string{"symbol": "TSLA-USDT", "incomeType": "REALIZED_PNL", "limit": "10"}, func(c *bxhttp.BaseHTTPClient) error {
			symbol, incomeType := "TSLA-USDT", "REALIZED_PNL"
			_, err := NewAccountService(c).GetIncomeHistory(&symbol, &incomeType, nil, nil, 10)
			return err
		}},
		{"cancel order", http.MethodDelete, "/openApi/swap/v2/trade/order", map[string]string{"symbol": "TSLA-USDT", "orderId": "123"}, func(c *bxhttp.BaseHTTPClient) error {
			orderID := "123"
			_, err := NewTradeService(c).CancelOrder("TSLA-USDT", &orderID, nil)
			return err
		}},
		{"set leverage", http.MethodPost, "/openApi/swap/v2/trade/leverage", map[string]string{"symbol": "TSLA-USDT", "leverage": "3", "side": "LONG"}, func(c *bxhttp.BaseHTTPClient) error {
			side := "LONG"
			_, err := NewTradeService(c).SetLeverage("TSLA-USDT", 3, &side)
			return err
		}},
		{"cancel twap", http.MethodDelete, "/openApi/swap/v2/trade/twap/order", map[string]string{"orderId": "twap-1"}, func(c *bxhttp.BaseHTTPClient) error {
			_, err := NewTradeService(c).CancelTWAPOrder("twap-1")
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, srv := newTestClient(t, tt.method, tt.path, tt.params)
			defer srv.Close()
			if err := tt.call(client); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestListenKeyServiceRoutes(t *testing.T) {
	for _, tt := range []struct {
		name, method string
		params       map[string]string
		call         func(*ListenKeyService) error
	}{
		{"create", http.MethodPost, nil, func(s *ListenKeyService) error { _, err := s.Create(); return err }},
		{"extend", http.MethodPut, map[string]string{"listenKey": "key"}, func(s *ListenKeyService) error { _, err := s.Extend("key"); return err }},
		{"delete", http.MethodDelete, map[string]string{"listenKey": "key"}, func(s *ListenKeyService) error { _, err := s.Delete("key"); return err }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, srv := newTestClient(t, tt.method, "/openApi/swap/v2/user/listenKey", tt.params)
			defer srv.Close()
			if err := tt.call(NewListenKeyService(client)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
