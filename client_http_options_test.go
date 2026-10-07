package bingx

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"

	bh "github.com/tigusigalpa/bingx-go/v2/http"
)

type testHTTPTransport func(*stdhttp.Request) (*stdhttp.Response, error)

func (f testHTTPTransport) RoundTrip(r *stdhttp.Request) (*stdhttp.Response, error) { return f(r) }

func TestMainClientPropagatesHTTPOptions(t *testing.T) {
	calls := 0
	injected := &stdhttp.Client{Transport: testHTTPTransport(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		calls++
		if r.URL.Host != "custom.invalid" {
			t.Errorf("incorrect endpoint: %s", r.URL.Host)
		}
		return &stdhttp.Response{StatusCode: 200, ContentLength: -1, Body: io.NopCloser(strings.NewReader(`{"code":0}`)), Header: make(stdhttp.Header)}, nil
	})}
	client := NewClient("key", "secret", WithBaseURI("http://custom.invalid"), WithHTTPClient(injected), WithMaxResponseBytes(10))
	if _, err := client.Market().GetKlinesRaw(context.Background(), "BTC-USDT", "1m", 10, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Account().GetBalance(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CoinM().Market().GetContracts(); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("shared client options not propagated: %d", calls)
	}
	client = NewDemoClient("key", "secret", WithBaseURI("http://custom.invalid"), WithHTTPClient(injected), WithMaxResponseBytes(8))
	if _, err := client.Market().GetOpenInterestRaw(context.Background(), "BTC-USDT"); !errors.Is(err, bh.ErrResponseTooLarge) {
		t.Fatalf("demo client lost limit: %v", err)
	}
}
