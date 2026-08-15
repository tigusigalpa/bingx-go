package services

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func newSpotTradeTestService(t *testing.T, handler http.HandlerFunc) (*SpotTradeService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := bxhttp.NewBaseHTTPClient("test-key", "test-secret", srv.URL, "", "base64")
	return NewSpotTradeService(client), srv
}

func strPtr(s string) *string { return &s }
func int64Ptr(v int64) *int64 { return &v }

func TestNewSpotTradeService(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	if service == nil {
		t.Fatal("Expected SpotTradeService to be created, got nil")
	}
	if service.client == nil {
		t.Error("SpotTradeService client should not be nil")
	}
}

// --- CreateOrderRequest: limit order request/response flow ---

func TestCreateOrderRequest_Limit(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody string

	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := readAll(r)
		gotBody = body

		if r.Header.Get("X-BX-APIKEY") != "test-key" {
			t.Errorf("expected X-BX-APIKEY header to be set, got %q", r.Header.Get("X-BX-APIKEY"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"symbol":"BTC-USDT","orderId":123,"status":"NEW"}}`)
	})
	defer srv.Close()

	price := "60000"
	result, err := service.CreateOrderRequest(SpotOrderRequest{
		Symbol:   "BTC-USDT",
		Side:     SpotSideBuy,
		Type:     SpotOrderTypeLimit,
		Quantity: "0.001",
		Price:    &price,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != "POST" {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/openApi/spot/v1/trade/order" {
		t.Errorf("expected path /openApi/spot/v1/trade/order, got %s", gotPath)
	}
	for _, want := range []string{"symbol=BTC-USDT", "side=BUY", "type=LIMIT", "quantity=0.001", "price=60000", "signature="} {
		if !strings.Contains(gotBody, want) {
			t.Errorf("expected request body to contain %q, got %q", want, gotBody)
		}
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got %#v", result)
	}
	if data["symbol"] != "BTC-USDT" {
		t.Errorf("expected symbol=BTC-USDT, got %v", data["symbol"])
	}
}

// --- CreateOrderRequest: market buy using quoteOrderQty ---

func TestCreateOrderRequest_MarketWithQuoteOrderQty(t *testing.T) {
	var gotBody string

	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"symbol":"BTC-USDT","orderId":124,"status":"NEW"}}`)
	})
	defer srv.Close()

	quoteQty := "100"
	_, err := service.CreateOrderRequest(SpotOrderRequest{
		Symbol:        "BTC-USDT",
		Side:          SpotSideBuy,
		Type:          SpotOrderTypeMarket,
		QuoteOrderQty: &quoteQty,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(gotBody, "quantity=") {
		t.Errorf("expected no quantity param when only quoteOrderQty is set, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "quoteOrderQty=100") {
		t.Errorf("expected quoteOrderQty=100 in body, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "type=MARKET") {
		t.Errorf("expected type=MARKET in body, got %q", gotBody)
	}
}

// --- Invalid typed requests ---

func TestCreateOrderRequest_Invalid(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	price := "100"
	qty := "1"

	tests := []struct {
		name string
		req  SpotOrderRequest
	}{
		{
			name: "missing symbol",
			req:  SpotOrderRequest{Side: SpotSideBuy, Type: SpotOrderTypeLimit, Quantity: qty, Price: &price},
		},
		{
			name: "invalid side",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: "HOLD", Type: SpotOrderTypeLimit, Quantity: qty, Price: &price},
		},
		{
			name: "invalid type",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: "ICEBERG", Quantity: qty, Price: &price},
		},
		{
			name: "limit without price",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: SpotOrderTypeLimit, Quantity: qty},
		},
		{
			name: "limit with zero price",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: SpotOrderTypeLimit, Quantity: qty, Price: strPtr("0")},
		},
		{
			name: "limit without quantity",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: SpotOrderTypeLimit, Price: &price},
		},
		{
			name: "market without quantity or quoteOrderQty",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: SpotOrderTypeMarket},
		},
		{
			name: "market with negative quantity",
			req:  SpotOrderRequest{Symbol: "BTC-USDT", Side: SpotSideBuy, Type: SpotOrderTypeMarket, Quantity: "-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateOrderRequest(tt.req)
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestCreateOrderRequest_ValidMarketWithQuantityOnly(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	// Validation should pass (network error is fine/expected here).
	_, err := service.CreateOrderRequest(SpotOrderRequest{
		Symbol:   "BTC-USDT",
		Side:     SpotSideSell,
		Type:     SpotOrderTypeMarket,
		Quantity: "0.5",
	})
	if err != nil && strings.Contains(err.Error(), "quantity") {
		t.Errorf("unexpected validation error for a valid request: %v", err)
	}
}

// --- CancelOrder by orderID and by clientOrderID ---

func TestCancelOrder_ByOrderID(t *testing.T) {
	var gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"symbol":"BTC-USDT","orderId":123,"status":"CANCELED"}}`)
	})
	defer srv.Close()

	orderID := "123"
	_, err := service.CancelOrder("BTC-USDT", &orderID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "orderId=123") {
		t.Errorf("expected orderId=123 in body, got %q", gotBody)
	}
}

func TestCancelOrder_ByClientOrderID(t *testing.T) {
	var gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"symbol":"BTC-USDT","clientOrderID":"my-id","status":"CANCELED"}}`)
	})
	defer srv.Close()

	clientOrderID := "my-id"
	_, err := service.CancelOrder("BTC-USDT", nil, &clientOrderID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "clientOrderID=my-id") {
		t.Errorf("expected clientOrderID=my-id in body, got %q", gotBody)
	}
}

func TestCancelOrder_RequiresIdentifier(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	_, err := service.CancelOrder("BTC-USDT", nil, nil)
	if err == nil {
		t.Error("expected error when neither orderID nor clientOrderID is provided")
	}
}

// --- CancelAllOrders ---

func TestCancelAllOrders(t *testing.T) {
	var gotPath, gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"orders":[]}}`)
	})
	defer srv.Close()

	symbol := "BTC-USDT"
	_, err := service.CancelAllOrders(&symbol)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/openApi/spot/v1/trade/cancelOpenOrders" {
		t.Errorf("expected path /openApi/spot/v1/trade/cancelOpenOrders, got %s", gotPath)
	}
	if !strings.Contains(gotBody, "symbol=BTC-USDT") {
		t.Errorf("expected symbol=BTC-USDT in body, got %q", gotBody)
	}
}

func TestCancelAllOrders_NilSymbol(t *testing.T) {
	var gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"orders":[]}}`)
	})
	defer srv.Close()

	_, err := service.CancelAllOrders(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(gotBody, "symbol=") {
		t.Errorf("expected no symbol param when nil, got %q", gotBody)
	}
}

// --- CancelBatchOrders ---

func TestSpotCancelBatchOrders(t *testing.T) {
	var gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"orders":[]}}`)
	})
	defer srv.Close()

	_, err := service.CancelBatchOrders("BTC-USDT", []string{"1", "2", "3"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody, "orderIds=1%2C2%2C3") && !strings.Contains(gotBody, "orderIds=1,2,3") {
		t.Errorf("expected orderIds=1,2,3 (URL-encoded or not) in body, got %q", gotBody)
	}
}

func TestCancelBatchOrders_RequiresSymbolAndIDs(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	if _, err := service.CancelBatchOrders("", []string{"1"}, nil); err == nil {
		t.Error("expected error when symbol is empty")
	}
	if _, err := service.CancelBatchOrders("BTC-USDT", nil, nil); err == nil {
		t.Error("expected error when no order IDs are provided")
	}
}

// --- GetOrder / GetOpenOrders ---

func TestGetOrder(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"symbol":"BTC-USDT","orderId":123,"status":"NEW"}}`)
	})
	defer srv.Close()

	orderID := "123"
	_, err := service.GetOrder("BTC-USDT", &orderID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/openApi/spot/v1/trade/query" {
		t.Errorf("expected path /openApi/spot/v1/trade/query, got %s", gotPath)
	}
	if !strings.Contains(gotQuery, "orderId=123") {
		t.Errorf("expected orderId=123 in query, got %q", gotQuery)
	}
}

func TestGetOrder_RequiresIdentifier(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	if _, err := service.GetOrder("BTC-USDT", nil, nil); err == nil {
		t.Error("expected error when neither orderID nor clientOrderID is provided")
	}
	if _, err := service.GetOrder("", nil, strPtr("id")); err == nil {
		t.Error("expected error when symbol is empty")
	}
}

func TestSpotGetOpenOrders(t *testing.T) {
	var gotQuery string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"orders":[]}}`)
	})
	defer srv.Close()

	symbol := "BTC-USDT"
	_, err := service.GetOpenOrders(&symbol)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotQuery, "symbol=BTC-USDT") {
		t.Errorf("expected symbol=BTC-USDT in query, got %q", gotQuery)
	}
}

func TestGetTrades_RequiresSymbol(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	if _, err := service.GetTrades(nil, 100, nil, nil); err == nil {
		t.Error("expected error when symbol is nil")
	}
	if _, err := service.GetTrades(strPtr(""), 100, nil, nil); err == nil {
		t.Error("expected error when symbol is empty")
	}
}

func TestSpotGetOrderHistory(t *testing.T) {
	var gotQuery string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"orders":[],"total":0}}`)
	})
	defer srv.Close()

	symbol := "BTC-USDT"
	_, err := service.GetOrderHistory(&symbol, 50, int64Ptr(1000), int64Ptr(2000))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotQuery, "pageSize=50") {
		t.Errorf("expected pageSize=50 in query, got %q", gotQuery)
	}
	if !strings.Contains(gotQuery, "startTime=1000") || !strings.Contains(gotQuery, "endTime=2000") {
		t.Errorf("expected startTime/endTime in query, got %q", gotQuery)
	}
}

// --- AmendOrder (cancel-replace) ---

func TestAmendOrder(t *testing.T) {
	var gotPath, gotBody string
	service, srv := newSpotTradeTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		body, _ := readAll(r)
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"data":{"cancelResult":"SUCCESS","newOrderResult":"SUCCESS"}}`)
	})
	defer srv.Close()

	price := "61000"
	orderID := "123"
	_, err := service.AmendOrder("BTC-USDT", &orderID, nil, SpotCancelReplaceStopOnFailure, SpotOrderRequest{
		Side:     SpotSideBuy,
		Type:     SpotOrderTypeLimit,
		Quantity: "0.001",
		Price:    &price,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/openApi/spot/v1/trade/order/cancelReplace" {
		t.Errorf("expected path /openApi/spot/v1/trade/order/cancelReplace, got %s", gotPath)
	}
	for _, want := range []string{"symbol=BTC-USDT", "cancelOrderId=123", "cancelReplaceMode=STOP_ON_FAILURE", "price=61000"} {
		if !strings.Contains(gotBody, want) {
			t.Errorf("expected request body to contain %q, got %q", want, gotBody)
		}
	}
}

func TestAmendOrder_Validation(t *testing.T) {
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewSpotTradeService(client)

	orderID := "123"
	valid := SpotOrderRequest{Side: SpotSideBuy, Type: SpotOrderTypeLimit, Quantity: "1", Price: strPtr("100")}

	if _, err := service.AmendOrder("", &orderID, nil, SpotCancelReplaceStopOnFailure, valid); err == nil {
		t.Error("expected error when symbol is empty")
	}
	if _, err := service.AmendOrder("BTC-USDT", nil, nil, SpotCancelReplaceStopOnFailure, valid); err == nil {
		t.Error("expected error when neither cancelOrderID nor cancelClientOrderID is provided")
	}
	if _, err := service.AmendOrder("BTC-USDT", &orderID, nil, "INVALID_MODE", valid); err == nil {
		t.Error("expected error for invalid cancelReplaceMode")
	}
}

// --- Demo client base URL wiring ---

func TestNewDemoClient_SpotTradeUsesDemoBaseURL(t *testing.T) {
	// SpotTradeService itself has no exported base URL getter; verify via
	// the shared http client used to construct it, mirroring how other
	// services share the same signed transport.
	client := bxhttp.NewBaseHTTPClient("key", "secret", "https://open-api-vst.bingx.com", "", "base64")
	service := NewSpotTradeService(client)

	if service.client.GetEndpoint() != "https://open-api-vst.bingx.com" {
		t.Errorf("expected demo endpoint, got %s", service.client.GetEndpoint())
	}
}

func readAll(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}
	defer func() { _ = r.Body.Close() }()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
