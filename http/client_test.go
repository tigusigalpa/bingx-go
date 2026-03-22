package http

import (
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
