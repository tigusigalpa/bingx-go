package http

import (
	"math"
	stdhttp "net/http"
	"time"
)

// ClientOption configures a BaseHTTPClient at construction time.
type ClientOption func(*BaseHTTPClient)

// WithHTTPClient injects a client, including its transport and timeout policy.
// Configure it before construction and do not mutate it concurrently with use.
// A nil client leaves the SDK's 30-second default in place.
func WithHTTPClient(client *stdhttp.Client) ClientOption {
	return func(c *BaseHTTPClient) {
		if client != nil {
			copyClient := *client
			c.httpClient = &copyClient
		}
	}
}

// WithMaxResponseBytes sets a positive response-body limit. Invalid limits cause
// requests to return ErrInvalidResponseLimit before any network I/O.
func WithMaxResponseBytes(limit int64) ClientOption {
	return func(c *BaseHTTPClient) {
		c.maxResponseBytes = limit
	}
}

// NewBaseHTTPClientWithOptions supports transport injection and body limits while
// NewBaseHTTPClient retains its original signature.
func NewBaseHTTPClientWithOptions(apiKey, apiSecret, baseURI, sourceKey, signatureEncoding string, options ...ClientOption) *BaseHTTPClient {
	c := &BaseHTTPClient{
		apiKey: apiKey, apiSecret: apiSecret, baseURI: baseURI,
		sourceKey: sourceKey, signatureEncoding: signatureEncoding,
		httpClient:       &stdhttp.Client{Timeout: 30 * time.Second},
		maxResponseBytes: DefaultMaxResponseBytes,
	}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	return c
}

func (c *BaseHTTPClient) validateLimit() error {
	if c.maxResponseBytes <= 0 || c.maxResponseBytes == math.MaxInt64 {
		return ErrInvalidResponseLimit
	}
	return nil
}
