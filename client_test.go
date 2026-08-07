package bingx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "test-api-secret"

	client := NewClient(apiKey, apiSecret)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.GetAPIKey() != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.GetAPIKey())
	}

	if client.GetEndpoint() != "https://open-api.bingx.com" {
		t.Errorf("Expected default endpoint https://open-api.bingx.com, got %s", client.GetEndpoint())
	}
}

func TestNewDemoClient(t *testing.T) {
	client := NewDemoClient("demo-key", "demo-secret")

	if client == nil {
		t.Fatal("Expected demo client to be created, got nil")
	}

	const demoBaseURI = "https://open-api-vst.bingx.com"

	if client.GetEndpoint() != demoBaseURI {
		t.Errorf("Expected demo endpoint %s, got %s", demoBaseURI, client.GetEndpoint())
	}

	// All services share the same signed httpClient, so SpotTrade (and every
	// other service) must transparently use the same demo base URI.
	if client.SpotTrade() == nil {
		t.Fatal("Expected SpotTrade service to be created on demo client, got nil")
	}
	if client.httpClient.GetEndpoint() != demoBaseURI {
		t.Errorf("Expected SpotTrade's underlying client endpoint %s, got %s", demoBaseURI, client.httpClient.GetEndpoint())
	}
}

func TestNewClientWithOptions(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "test-api-secret"
	customURI := "https://custom-api.bingx.com"
	sourceKey := "test-source-key"

	client := NewClient(
		apiKey,
		apiSecret,
		WithBaseURI(customURI),
		WithSourceKey(sourceKey),
		WithSignatureEncoding("hex"),
	)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.GetEndpoint() != customURI {
		t.Errorf("Expected custom endpoint %s, got %s", customURI, client.GetEndpoint())
	}
}

func TestClientServices(t *testing.T) {
	client := NewClient("test-key", "test-secret")

	tests := []struct {
		name    string
		service interface{}
	}{
		{"Market", client.Market()},
		{"Account", client.Account()},
		{"Trade", client.Trade()},
		{"Contract", client.Contract()},
		{"ListenKey", client.ListenKey()},
		{"Wallet", client.Wallet()},
		{"SpotAccount", client.SpotAccount()},
		{"SpotTrade", client.SpotTrade()},
		{"SubAccount", client.SubAccount()},
		{"CopyTrading", client.CopyTrading()},
		{"TradFi", client.TradFi()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.service == nil {
				t.Errorf("%s service should not be nil", tt.name)
			}
		})
	}
}

func TestClientCoinM(t *testing.T) {
	client := NewClient("test-key", "test-secret")

	coinm := client.CoinM()
	if coinm == nil {
		t.Fatal("CoinM client should not be nil")
	}

	coinm2 := client.CoinM()
	if coinm != coinm2 {
		t.Error("CoinM should return the same instance")
	}
}

func TestClientWebSocketStreams(t *testing.T) {
	client := NewClient("test-key", "test-secret")

	marketStream := client.NewMarketDataStream()
	if marketStream == nil {
		t.Error("Market data stream should not be nil")
	}

	accountStream := client.NewAccountDataStream("test-listen-key")
	if accountStream == nil {
		t.Error("Account data stream should not be nil")
	}
}

func TestClientOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   ClientOption
		checkFn  func(*ClientConfig) bool
		expected bool
	}{
		{
			name:   "WithBaseURI",
			option: WithBaseURI("https://test.com"),
			checkFn: func(c *ClientConfig) bool {
				return c.BaseURI == "https://test.com"
			},
			expected: true,
		},
		{
			name:   "WithSourceKey",
			option: WithSourceKey("source-123"),
			checkFn: func(c *ClientConfig) bool {
				return c.SourceKey == "source-123"
			},
			expected: true,
		},
		{
			name:   "WithSignatureEncoding",
			option: WithSignatureEncoding("hex"),
			checkFn: func(c *ClientConfig) bool {
				return c.SignatureEncoding == "hex"
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ClientConfig{}
			tt.option(config)
			if tt.checkFn(config) != tt.expected {
				t.Errorf("Option %s did not apply correctly", tt.name)
			}
		})
	}
}

func TestCoinMClient(t *testing.T) {
	client := NewClient("test-key", "test-secret")
	coinm := client.CoinM()

	if coinm.Market() == nil {
		t.Error("CoinM Market service should not be nil")
	}

	if coinm.Trade() == nil {
		t.Error("CoinM Trade service should not be nil")
	}

	if coinm.ListenKey() == nil {
		t.Error("CoinM ListenKey service should not be nil")
	}
}

func TestTradFiClient(t *testing.T) {
	client := NewClient("test-key", "test-secret")
	tradfi := client.TradFi()

	if tradfi == nil {
		t.Fatal("TradFi client should not be nil")
	}

	// Test singleton pattern
	tradfi2 := client.TradFi()
	if tradfi != tradfi2 {
		t.Error("TradFi should return the same instance")
	}

	// Test all services
	if tradfi.Market() == nil {
		t.Error("TradFi Market service should not be nil")
	}

	if tradfi.Trade() == nil {
		t.Error("TradFi Trade service should not be nil")
	}

	if tradfi.Account() == nil {
		t.Error("TradFi Account service should not be nil")
	}

	if tradfi.ListenKey() == nil {
		t.Error("TradFi ListenKey service should not be nil")
	}
}

const testClientSecret = "test-secret"

func expectedClientHexSignature(signingString string) string {
	h := hmac.New(sha256.New, []byte(testClientSecret))
	h.Write([]byte(signingString))
	return hex.EncodeToString(h.Sum(nil))
}

func newSignatureTestServer(t *testing.T) (*httptest.Server, chan string) {
	received := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string
		if r.Method == "POST" {
			b, _ := io.ReadAll(r.Body)
			body = string(b)
		} else {
			body = r.URL.RawQuery
		}
		received <- body

		values, err := url.ParseQuery(body)
		if err != nil {
			t.Errorf("failed to parse received query/body: %v", err)
		} else {
			sig := values.Get("signature")
			values.Del("signature")
			signingString := values.Encode()
			expected := expectedClientHexSignature(signingString)
			if sig != expected {
				t.Errorf("signature mismatch: got %s, want %s", sig, expected)
			}
			if len(sig) != 64 {
				t.Errorf("expected 64-character hex signature, got %d", len(sig))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0}`)
	}))
	return srv, received
}

func TestNewClient_UsesHexSignatureByDefault(t *testing.T) {
	srv, received := newSignatureTestServer(t)
	defer srv.Close()

	client := NewClient("test-key", testClientSecret, WithBaseURI(srv.URL))
	_, err := client.GetHTTPClient().Request("GET", "/test", map[string]interface{}{
		"symbol":    "BTC-USDT",
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case body := <-received:
		if !hasSignature(body) {
			t.Errorf("expected signature in request, got %s", body)
		}
	default:
		t.Fatal("server did not receive a request")
	}
}

func TestNewDemoClient_UsesHexSignatureByDefault(t *testing.T) {
	srv, received := newSignatureTestServer(t)
	defer srv.Close()

	client := NewDemoClient("test-key", testClientSecret, WithBaseURI(srv.URL))
	_, err := client.GetHTTPClient().Request("GET", "/test", map[string]interface{}{
		"symbol":    "BTC-USDT",
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case body := <-received:
		if !hasSignature(body) {
			t.Errorf("expected signature in request, got %s", body)
		}
	default:
		t.Fatal("server did not receive a request")
	}
}

func TestWithSignatureEncoding_HexStillWorks(t *testing.T) {
	srv, received := newSignatureTestServer(t)
	defer srv.Close()

	client := NewClient("test-key", testClientSecret, WithBaseURI(srv.URL), WithSignatureEncoding("hex"))
	_, err := client.GetHTTPClient().Request("GET", "/test", map[string]interface{}{
		"symbol":    "BTC-USDT",
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case body := <-received:
		if !hasSignature(body) {
			t.Errorf("expected signature in request, got %s", body)
		}
	default:
		t.Fatal("server did not receive a request")
	}
}

func TestDefaultClient_SpotAccountGetBalance_SendsHexSignature(t *testing.T) {
	srv, received := newSignatureTestServer(t)
	defer srv.Close()

	client := NewClient("test-key", testClientSecret, WithBaseURI(srv.URL))
	_, err := client.SpotAccount().GetBalance()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case body := <-received:
		if !hasSignature(body) {
			t.Errorf("expected signature in request, got %s", body)
		}
	default:
		t.Fatal("server did not receive a request")
	}
}

func TestDefaultClient_AccountGetBalance_SendsHexSignature(t *testing.T) {
	srv, received := newSignatureTestServer(t)
	defer srv.Close()

	client := NewClient("test-key", testClientSecret, WithBaseURI(srv.URL))
	_, err := client.Account().GetBalance()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case body := <-received:
		if !hasSignature(body) {
			t.Errorf("expected signature in request, got %s", body)
		}
	default:
		t.Fatal("server did not receive a request")
	}
}

func hasSignature(body string) bool {
	values, err := url.ParseQuery(body)
	if err != nil {
		return false
	}
	sig := values.Get("signature")
	return len(sig) == 64
}
