package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	bh "github.com/tigusigalpa/bingx-go/v2/http"
)

func TestEveryRawMarketMethodPreservesContextAndEvidence(t *testing.T) {
	start, end, fromID := int64(123), int64(456), int64(789)
	methods := []struct {
		name, path string
		call       func(*MarketService, context.Context) (*RawResponse, error)
	}{
		{"klines", "/openApi/swap/v3/quote/klines", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetKlinesRaw(ctx, "BTC-USDT", "1m", 10, &start, &end)
		}},
		{"recent", "/openApi/swap/v2/quote/trades", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetRecentTradesRaw(ctx, "BTC-USDT", 10)
		}},
		{"aggregate-legacy", "/openApi/swap/v2/market/aggTrades", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetAggregateTradesRaw(ctx, "BTC-USDT", 10, &fromID, &start, &end)
		}},
		{"OI", "/openApi/swap/v2/quote/openInterest", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetOpenInterestRaw(ctx, "BTC-USDT")
		}},
		{"premium", "/openApi/swap/v2/quote/premiumIndex", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetPremiumIndexRaw(ctx, "BTC-USDT")
		}},
		{"funding", "/openApi/swap/v2/quote/fundingRate", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetFundingRatesRaw(ctx, "BTC-USDT", &start, &end, 10)
		}},
		{"mark", "/openApi/swap/v1/market/markPriceKlines", func(s *MarketService, ctx context.Context) (*RawResponse, error) {
			return s.GetMarkPriceKlinesRaw(ctx, "BTC-USDT", "1m", 10, &start, &end)
		}},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			var calls atomic.Int64
			body, err := os.ReadFile("testdata/contracts/klines.synthetic.json")
			if err != nil {
				t.Fatal(err)
			}
			if method.name == "funding" {
				body, err = os.ReadFile("testdata/contracts/funding.synthetic.json")
				if err != nil {
					t.Fatal(err)
				}
			}
			server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				calls.Add(1)
				if r.URL.Path != method.path {
					t.Errorf("route changed: %s", r.URL.Path)
				}
				if r.URL.Query().Get("symbol") != "BTC-USDT" {
					t.Error("missing symbol")
				}
				_, _ = w.Write(body)
			}))
			defer server.Close()
			service := NewMarketService(bh.NewBaseHTTPClientWithOptions("key", "secret", server.URL, "", "hex", bh.WithMaxResponseBytes(1024)))
			receipt, err := method.call(service, context.Background())
			if err != nil || string(receipt.BodyBytes()) != string(body) || receipt.Route() != method.path {
				t.Fatalf("invalid receipt: %v / %v", receipt, err)
			}
			for _, deadline := range []bool{false, true} {
				ctx, cancel := context.WithCancel(context.Background())
				want := context.Canceled
				if deadline {
					ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
					want = context.DeadlineExceeded
				} else {
					cancel()
				}
				blocked, err := method.call(service, ctx)
				cancel()
				if blocked != nil || !errors.Is(err, want) {
					t.Fatalf("raw method lost context: %v / %v", blocked, err)
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("unexpected retries or canceled fetches: %d", calls.Load())
			}
		})
	}
}

func TestMarketLimitAppliesToLegacyMethod(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) { _, _ = io.WriteString(w, `{"code":0,"data":[]}`) }))
	defer server.Close()
	service := NewMarketService(bh.NewBaseHTTPClientWithOptions("key", "secret", server.URL, "", "hex", bh.WithMaxResponseBytes(5)))
	if _, err := service.GetKlines("BTC-USDT", "1m", 10, nil, nil); !errors.Is(err, bh.ErrResponseTooLarge) {
		t.Fatalf("legacy method ignored limit: %v", err)
	}
}

// The loopback replay tests SDK handling, not the exchange's live behavior.
func TestFundingDocumentationExamplePreserved(t *testing.T) {
	fixtureBytes, err := os.ReadFile("testdata/contracts/swap-market.documentation.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Funding struct {
			ResponseExample json.RawMessage `json:"response_example"`
		} `json:"funding"`
	}
	if err := json.Unmarshal(fixtureBytes, &fixture); err != nil {
		t.Fatal(err)
	}
	body := fixture.Funding.ResponseExample
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.URL.Path != "/openApi/swap/v2/quote/fundingRate" || r.URL.Query().Get("symbol") != "QNT-USDT" {
			t.Error("funding route or selector changed")
		}
		_, _ = w.Write(body)
	}))
	defer server.Close()
	service := NewMarketService(bh.NewBaseHTTPClientWithOptions("", "", server.URL, "", "hex", bh.WithMaxResponseBytes(1024)))
	receipt, err := service.GetFundingRatesRaw(context.Background(), "QNT-USDT", nil, nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(receipt.BodyBytes(), body) || receipt.Selectors()["limit"] != "1" {
		t.Fatal("documentation example or selectors were changed")
	}
	var response struct {
		Data []struct {
			Symbol      string `json:"symbol"`
			FundingRate string `json:"fundingRate"`
			FundingTime int64  `json:"fundingTime"`
		} `json:"data"`
	}
	if err := json.Unmarshal(receipt.BodyBytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Data) != 1 || response.Data[0].Symbol != "QNT-USDT" || response.Data[0].FundingRate != "0.00027100" || response.Data[0].FundingTime != 1702713600000 {
		t.Fatal("documented array, decimal string or millisecond time lost")
	}
}
