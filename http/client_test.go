package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewBaseHTTPClient(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "test-api-secret"
	baseURI := "https://api.test.com"
	sourceKey := "source-key"
	signatureEncoding := "base64"

	client := NewBaseHTTPClient(apiKey, apiSecret, baseURI, sourceKey, signatureEncoding)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.apiSecret != apiSecret {
		t.Errorf("Expected API secret %s, got %s", apiSecret, client.apiSecret)
	}

	if client.baseURI != baseURI {
		t.Errorf("Expected base URI %s, got %s", baseURI, client.baseURI)
	}

	if client.sourceKey != sourceKey {
		t.Errorf("Expected source key %s, got %s", sourceKey, client.sourceKey)
	}

	if client.signatureEncoding != signatureEncoding {
		t.Errorf("Expected signature encoding %s, got %s", signatureEncoding, client.signatureEncoding)
	}
}

func TestTimestamp(t *testing.T) {
	client := NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")

	ts1 := client.timestamp()
	if ts1 == "" {
		t.Error("Timestamp should not be empty")
	}

	if len(ts1) < 13 {
		t.Error("Timestamp should be at least 13 characters (milliseconds)")
	}
}

func TestBuildQuery(t *testing.T) {
	client := NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")

	tests := []struct {
		name     string
		params   map[string]interface{}
		contains []string
	}{
		{
			name:     "Empty params",
			params:   map[string]interface{}{},
			contains: []string{},
		},
		{
			name: "String param",
			params: map[string]interface{}{
				"symbol": "BTC-USDT",
			},
			contains: []string{"symbol=BTC-USDT"},
		},
		{
			name: "Int param",
			params: map[string]interface{}{
				"limit": 100,
			},
			contains: []string{"limit=100"},
		},
		{
			name: "Float param",
			params: map[string]interface{}{
				"price": 50000.5,
			},
			contains: []string{"price=50000.5"},
		},
		{
			name: "Bool param",
			params: map[string]interface{}{
				"test": true,
			},
			contains: []string{"test=true"},
		},
		{
			name: "Multiple params",
			params: map[string]interface{}{
				"symbol": "BTC-USDT",
				"limit":  100,
			},
			contains: []string{"symbol=BTC-USDT", "limit=100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := client.buildQuery(tt.params)

			if len(tt.params) == 0 && query != "" {
				t.Errorf("Expected empty query, got %s", query)
			}

			for _, expected := range tt.contains {
				if query == "" {
					t.Errorf("Expected query to contain %s, got empty string", expected)
					continue
				}
			}
		})
	}
}

func TestSignString(t *testing.T) {
	tests := []struct {
		name              string
		apiSecret         string
		signatureEncoding string
		input             string
		expectedLen       int
	}{
		{
			name:              "Base64 encoding",
			apiSecret:         "test-secret",
			signatureEncoding: "base64",
			input:             "test-string",
			expectedLen:       44,
		},
		{
			name:              "Hex encoding",
			apiSecret:         "test-secret",
			signatureEncoding: "hex",
			input:             "test-string",
			expectedLen:       64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewBaseHTTPClient("key", tt.apiSecret, "https://api.test.com", "", tt.signatureEncoding)
			signature := client.signString(tt.input)

			if len(signature) != tt.expectedLen {
				t.Errorf("Expected signature length %d, got %d", tt.expectedLen, len(signature))
			}

			if signature == "" {
				t.Error("Signature should not be empty")
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		sourceKey string
		expected  map[string]string
	}{
		{
			name:      "Without source key",
			apiKey:    "test-api-key",
			sourceKey: "",
			expected: map[string]string{
				"X-BX-APIKEY":  "test-api-key",
				"Content-Type": "application/x-www-form-urlencoded",
			},
		},
		{
			name:      "With source key",
			apiKey:    "test-api-key",
			sourceKey: "source-123",
			expected: map[string]string{
				"X-BX-APIKEY":  "test-api-key",
				"Content-Type": "application/x-www-form-urlencoded",
				"X-SOURCE-KEY": "source-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewBaseHTTPClient(tt.apiKey, "secret", "https://api.test.com", tt.sourceKey, "base64")
			headers := client.headers()

			for key, expectedValue := range tt.expected {
				if headers[key] != expectedValue {
					t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, headers[key])
				}
			}
		})
	}
}

func TestHandleAPIError(t *testing.T) {
	client := NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")

	tests := []struct {
		name        string
		response    map[string]interface{}
		expectError bool
		errorType   string
	}{
		{
			name:        "No error",
			response:    map[string]interface{}{"data": "success"},
			expectError: false,
		},
		{
			name: "Authentication error",
			response: map[string]interface{}{
				"code": "100001",
				"msg":  "Invalid API key",
			},
			expectError: true,
			errorType:   "AuthenticationException",
		},
		{
			name: "Rate limit error",
			response: map[string]interface{}{
				"code": "100005",
				"msg":  "Rate limit exceeded",
			},
			expectError: true,
			errorType:   "RateLimitException",
		},
		{
			name: "Insufficient balance error",
			response: map[string]interface{}{
				"code": "200001",
				"msg":  "Insufficient balance",
			},
			expectError: true,
			errorType:   "InsufficientBalanceException",
		},
		{
			name: "Generic API error",
			response: map[string]interface{}{
				"code": "999999",
				"msg":  "Unknown error",
			},
			expectError: true,
			errorType:   "APIException",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.handleAPIError(tt.response)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestGetEndpoint(t *testing.T) {
	baseURI := "https://custom-api.bingx.com"
	client := NewBaseHTTPClient("key", "secret", baseURI, "", "base64")

	if client.GetEndpoint() != baseURI {
		t.Errorf("Expected endpoint %s, got %s", baseURI, client.GetEndpoint())
	}
}

func TestGetAPIKey(t *testing.T) {
	apiKey := "my-api-key"
	client := NewBaseHTTPClient(apiKey, "secret", "https://api.test.com", "", "base64")

	if client.GetAPIKey() != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.GetAPIKey())
	}
}

const testSignatureSecret = "test-secret"

// knownHexSignature is the expected HMAC-SHA256 hex signature for the
// signing string "symbol=BTC-USDT&timestamp=1702731500000" with secret
// "test-secret". It is used as a deterministic test vector.
const knownHexSignature = "36f610a8da9e21e9c417413745bcca4dddad81cea81e6d66baa3ef0c2229021c"

func expectedHexSignature(signingString string) string {
	h := hmac.New(sha256.New, []byte(testSignatureSecret))
	h.Write([]byte(signingString))
	return hex.EncodeToString(h.Sum(nil))
}

func TestSignString_HexKnownVector(t *testing.T) {
	client := NewBaseHTTPClient("test-key", testSignatureSecret, "https://api.test.com", "", "hex")
	got := client.signString("symbol=BTC-USDT&timestamp=1702731500000")
	if got != knownHexSignature {
		t.Errorf("hex signature mismatch: got %s, want %s", got, knownHexSignature)
	}
	if len(got) != 64 {
		t.Errorf("expected 64-character lowercase hex signature, got %d", len(got))
	}
}

func TestSignString_Base64StillSupported(t *testing.T) {
	client := NewBaseHTTPClient("test-key", testSignatureSecret, "https://api.test.com", "", "base64")
	got := client.signString("symbol=BTC-USDT&timestamp=1702731500000")
	if len(got) != 44 {
		t.Errorf("expected 44-character base64 signature, got %d", len(got))
	}
}

func TestSignString_UnknownEncodingFallsBackToHexNotBase64(t *testing.T) {
	client := NewBaseHTTPClient("test-key", testSignatureSecret, "https://api.test.com", "", "unknown")
	got := client.signString("test")
	if len(got) != 64 {
		t.Errorf("unknown encoding should produce 64-character hex, got %d: %s", len(got), got)
	}
}

func TestSignString_EmptyEncodingDefaultsToHex(t *testing.T) {
	client := NewBaseHTTPClient("test-key", testSignatureSecret, "https://api.test.com", "", "")
	got := client.signString("test")
	if len(got) != 64 {
		t.Errorf("empty encoding should default to 64-character hex, got %d: %s", len(got), got)
	}
}

func TestRequestSignatureIsVerifiedByServer(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"GET", "GET"},
		{"POST", "POST"},
		{"DELETE", "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var received string
			srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				if tt.method == "POST" {
					b, _ := io.ReadAll(r.Body)
					received = string(b)
				} else {
					received = r.URL.RawQuery
				}

				values, err := url.ParseQuery(received)
				if err != nil {
					t.Fatalf("failed to parse received %s: %v", tt.method, err)
				}

				sig := values.Get("signature")
				values.Del("signature")
				signingString := values.Encode()

				if signingString != "symbol=BTC-USDT&timestamp=1702731500000" {
					t.Errorf("unexpected signing string: %s", signingString)
				}

				expected := expectedHexSignature(signingString)
				if sig != expected {
					t.Errorf("signature mismatch for %s: got %s, want %s", tt.method, sig, expected)
				}
				if len(sig) != 64 {
					t.Errorf("expected 64-character hex signature for %s, got %d", tt.method, len(sig))
				}

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"code":0}`)
			}))
			defer srv.Close()

			client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
			_, err := client.Request(tt.method, "/test", map[string]interface{}{
				"symbol":    "BTC-USDT",
				"timestamp": "1702731500000",
			})
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", tt.method, err)
			}
		})
	}
}
