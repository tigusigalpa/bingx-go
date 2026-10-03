package services

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	bingxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func TestMarketRawMethodsUseVerifiedRoutes(t *testing.T) {
	paths := make(chan string, 2)
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		paths <- r.URL.Path
		if _, err := w.Write([]byte(`{"code":0,"data":[]}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	service := NewMarketService(bingxhttp.NewBaseHTTPClient("key", "secret", server.URL, "", "hex"))
	if _, err := service.GetFundingRatesRaw(context.Background(), "BTC-USDT", nil, nil, 100); err != nil {
		t.Fatalf("GetFundingRatesRaw returned error: %v", err)
	}
	if got := <-paths; got != "/openApi/swap/v2/quote/fundingRate" {
		t.Fatalf("unexpected funding route: %s", got)
	}
	if _, err := service.GetMarkPriceKlinesRaw(context.Background(), "BTC-USDT", "1h", 10, nil, nil); err != nil {
		t.Fatalf("GetMarkPriceKlinesRaw returned error: %v", err)
	}
	if got := <-paths; got != "/openApi/swap/v1/market/markPriceKlines" {
		t.Fatalf("unexpected mark-candle route: %s", got)
	}
}
