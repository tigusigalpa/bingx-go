package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	bxerrors "github.com/tigusigalpa/bingx-go/v2/errors"
	bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func newBookTickerTestService(t *testing.T, handler http.HandlerFunc) (*MarketService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := bxhttp.NewBaseHTTPClient("test-key", "test-secret", srv.URL, "", "hex")
	return NewMarketService(client), srv
}

func writeBookTickerResponse(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(response))
}

func TestGetBookTickerData(t *testing.T) {
	fixture, err := os.ReadFile(filepath.Join("testdata", "futures_book_ticker.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	symbol := "ETH-USDT"
	service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/openApi/swap/v2/quote/bookTicker" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("symbol"); got != symbol {
			t.Errorf("symbol = %q, want %q", got, symbol)
		}
		writeBookTickerResponse(w, string(fixture))
	})
	defer srv.Close()

	ticker, err := service.GetBookTickerData(&symbol)
	if err != nil {
		t.Fatalf("GetBookTickerData() error = %v", err)
	}

	if ticker.Symbol != "ETH-USDT" || ticker.BidPrice != "1916.79" || ticker.BidQuantity != "657.08" || ticker.AskPrice != "1916.8" || ticker.AskQuantity != "393.12" || ticker.LastUpdateID != "533536907" || ticker.Timestamp != "1786280077438" {
		t.Errorf("GetBookTickerData() = %+v", ticker)
	}
}

func TestGetBookTickerDataPreservesStringNumbers(t *testing.T) {
	service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
		writeBookTickerResponse(w, `{"code":0,"data":{"book_ticker":{"symbol":"ETH-USDT","bid_price":"1916.790000000000000001","bid_qty":"657.080000000000000001","ask_price":"1916.800000000000000001","ask_qty":"393.120000000000000001","lastUpdateId":"533536907","time":"1786280077438"}}}`)
	})
	defer srv.Close()

	ticker, err := service.GetBookTickerData(nil)
	if err != nil {
		t.Fatalf("GetBookTickerData() error = %v", err)
	}
	if ticker.BidPrice != "1916.790000000000000001" || ticker.BidQuantity != "657.080000000000000001" || ticker.AskPrice != "1916.800000000000000001" || ticker.AskQuantity != "393.120000000000000001" {
		t.Errorf("decimal strings were not preserved: %+v", ticker)
	}
}

func TestGetBookTickerDataMalformedResponses(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{name: "missing data", response: `{"code":0}`},
		{name: "missing book ticker", response: `{"code":0,"data":{}}`},
		{name: "invalid field type", response: `{"code":0,"data":{"book_ticker":{"symbol":"ETH-USDT","bid_price":true,"bid_qty":657.08,"ask_price":1916.8,"ask_qty":393.12,"lastUpdateId":533536907,"time":1786280077438}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
				writeBookTickerResponse(w, tt.response)
			})
			defer srv.Close()

			if _, err := service.GetBookTickerData(nil); err == nil {
				t.Fatal("GetBookTickerData() error = nil")
			}
		})
	}
}

func TestGetBookTickerDataPreservesAPIErrors(t *testing.T) {
	service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
		writeBookTickerResponse(w, `{"code":100005,"msg":"Rate limit exceeded"}`)
	})
	defer srv.Close()

	_, err := service.GetBookTickerData(nil)
	var apiErr *bxerrors.RateLimitException
	if !errors.As(err, &apiErr) {
		t.Fatalf("GetBookTickerData() error = %T %v, want RateLimitException", err, err)
	}
}

func TestGetBookTickerRawResponseRemainsUnchanged(t *testing.T) {
	fixture, err := os.ReadFile(filepath.Join("testdata", "futures_book_ticker.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
		writeBookTickerResponse(w, string(fixture))
	})
	defer srv.Close()

	response, err := service.GetBookTicker(nil)
	if err != nil {
		t.Fatalf("GetBookTicker() error = %v", err)
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("data = %T, want map[string]interface{}", response["data"])
	}
	bookTicker, ok := data["book_ticker"].(map[string]interface{})
	if !ok {
		t.Fatalf("data.book_ticker = %T, want map[string]interface{}", data["book_ticker"])
	}
	if bidPrice, ok := bookTicker["bid_price"].(float64); !ok || bidPrice != 1916.79 {
		t.Errorf("data.book_ticker.bid_price = %#v, want float64(1916.79)", bookTicker["bid_price"])
	}
}

func TestGetSpotBookTickerData(t *testing.T) {
	symbol := "ETH-USDT"
	service, srv := newBookTickerTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/openApi/spot/v1/ticker/bookTicker" {
			t.Errorf("path = %q", r.URL.Path)
		}
		writeBookTickerResponse(w, `{"code":0,"data":[{"eventType":"bookTicker","time":1786280077438,"symbol":"ETH-USDT","bidPrice":"1916.790000000000000001","bidVolume":"657.080000000000000001","askPrice":"1916.800000000000000001","askVolume":"393.120000000000000001"}]}`)
	})
	defer srv.Close()

	ticker, err := service.GetSpotBookTickerData(&symbol)
	if err != nil {
		t.Fatalf("GetSpotBookTickerData() error = %v", err)
	}
	if ticker.Symbol != symbol || ticker.BidPrice != "1916.790000000000000001" || ticker.BidQuantity != "657.080000000000000001" || ticker.AskPrice != "1916.800000000000000001" || ticker.AskQuantity != "393.120000000000000001" || ticker.Timestamp != "1786280077438" {
		t.Errorf("GetSpotBookTickerData() = %+v", ticker)
	}
}
