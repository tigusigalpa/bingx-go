package coinm

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
		name   string
		path   string
		params map[string]string
		call   func(*MarketService) error
	}{
		{"contracts", "/openApi/cswap/v1/market/contracts", nil, func(s *MarketService) error { _, err := s.GetContracts(); return err }},
		{"ticker", "/openApi/cswap/v1/market/ticker", map[string]string{"symbol": "BTC-USD"}, func(s *MarketService) error { _, err := s.GetTicker("BTC-USD"); return err }},
		{"depth", "/openApi/cswap/v1/market/depth", map[string]string{"symbol": "BTC-USD", "limit": "20"}, func(s *MarketService) error { _, err := s.GetDepth("BTC-USD", 20); return err }},
		{"klines", "/openApi/cswap/v1/market/klines", map[string]string{"symbol": "BTC-USD", "interval": "1h", "limit": "100", "startTime": "1", "endTime": "2"}, func(s *MarketService) error {
			start, end := int64(1), int64(2)
			_, err := s.GetKlines("BTC-USD", "1h", 100, &start, &end)
			return err
		}},
		{"recent trades", "/openApi/cswap/v1/market/trades", map[string]string{"symbol": "BTC-USD", "limit": "10"}, func(s *MarketService) error { _, err := s.GetRecentTrades("BTC-USD", 10); return err }},
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

func TestTradeServiceRoutes(t *testing.T) {
	orderID, clientOrderID := "123", "client-123"
	tests := []struct {
		name   string
		method string
		path   string
		params map[string]string
		call   func(*TradeService) error
	}{
		{"create order", http.MethodPost, "/openApi/cswap/v1/trade/order", map[string]string{"symbol": "BTC-USD", "side": "BUY"}, func(s *TradeService) error {
			_, err := s.CreateOrder(map[string]interface{}{"symbol": "BTC-USD", "side": "BUY"})
			return err
		}},
		{"cancel order", http.MethodDelete, "/openApi/cswap/v1/trade/cancelOrder", map[string]string{"symbol": "BTC-USD", "orderId": orderID, "clientOrderId": clientOrderID}, func(s *TradeService) error { _, err := s.CancelOrder("BTC-USD", &orderID, &clientOrderID); return err }},
		{"set leverage", http.MethodPost, "/openApi/cswap/v1/trade/leverage", map[string]string{"symbol": "BTC-USD", "side": "LONG", "leverage": "5"}, func(s *TradeService) error { _, err := s.SetLeverage("BTC-USD", "LONG", 5); return err }},
		{"income history", http.MethodGet, "/openApi/cswap/v1/user/income", map[string]string{"symbol": "BTC-USD", "incomeType": "REALIZED_PNL", "limit": "20", "recvWindow": "5000"}, func(s *TradeService) error {
			symbol, incomeType, recvWindow := "BTC-USD", "REALIZED_PNL", int64(5000)
			_, err := s.GetIncomeHistory(&symbol, &incomeType, nil, nil, 20, &recvWindow)
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, srv := newTestClient(t, tt.method, tt.path, tt.params)
			defer srv.Close()
			if err := tt.call(NewTradeService(client)); err != nil {
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
		{"generate", http.MethodPost, nil, func(s *ListenKeyService) error { _, err := s.Generate(); return err }},
		{"extend", http.MethodPut, map[string]string{"listenKey": "key"}, func(s *ListenKeyService) error { _, err := s.Extend("key"); return err }},
		{"delete", http.MethodDelete, map[string]string{"listenKey": "key"}, func(s *ListenKeyService) error { _, err := s.Delete("key"); return err }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, srv := newTestClient(t, tt.method, "/openApi/user/auth/userDataStream", tt.params)
			defer srv.Close()
			if err := tt.call(NewListenKeyService(client)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
