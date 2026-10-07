package http

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tigusigalpa/bingx-go/v2/errors"
)

// BaseHTTPClient represents a BingX API component or value.
type BaseHTTPClient struct {
	apiKey            string
	apiSecret         string
	baseURI           string
	sourceKey         string
	signatureEncoding string
	httpClient        *http.Client
	maxResponseBytes  int64
}

// NewBaseHTTPClient creates a new client or service instance.
func NewBaseHTTPClient(apiKey, apiSecret, baseURI, sourceKey, signatureEncoding string) *BaseHTTPClient {
	return NewBaseHTTPClientWithOptions(apiKey, apiSecret, baseURI, sourceKey, signatureEncoding)
}

// timestamp performs the timestamp operation.
func (c *BaseHTTPClient) timestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// sortedKeys performs the sortedKeys operation.
func (c *BaseHTTPClient) sortedKeys(params map[string]interface{}) []string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// paramValueToString performs the paramValueToString operation.
func (c *BaseHTTPClient) paramValueToString(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}

	switch val := v.(type) {
	case string:
		return val, true
	case *string:
		if val == nil {
			return "", false
		}
		return *val, true
	case int:
		return strconv.Itoa(val), true
	case *int:
		if val == nil {
			return "", false
		}
		return strconv.Itoa(*val), true
	case int64:
		return strconv.FormatInt(val, 10), true
	case *int64:
		if val == nil {
			return "", false
		}
		return strconv.FormatInt(*val, 10), true
	case uint:
		return strconv.FormatUint(uint64(val), 10), true
	case uint64:
		return strconv.FormatUint(val, 10), true
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), true
	case *float64:
		if val == nil {
			return "", false
		}
		return strconv.FormatFloat(*val, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(val), true
	case *bool:
		if val == nil {
			return "", false
		}
		return strconv.FormatBool(*val), true
	case []byte:
		return string(val), true
	}

	// Complex values (slices, maps, structs) are serialized as JSON.
	b, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// buildCanonicalString returns the raw, URL-unencoded, sorted "key=value&..."
// string that is signed. It never includes the signature parameter.
func (c *BaseHTTPClient) buildCanonicalString(params map[string]interface{}) string {
	keys := c.sortedKeys(params)
	parts := make([]string, 0, len(keys))

	for _, k := range keys {
		v, ok := c.paramValueToString(params[k])
		if !ok {
			continue
		}
		parts = append(parts, k+"="+v)
	}

	return strings.Join(parts, "&")
}

// buildSignedString builds the wire representation used in the final URL query
// or POST body. The canonical part is sorted in the same way as the signing
// string. Values containing '[' or '{' are URL-escaped only when forURL is true
// (GET/DELETE query strings). The signature is always appended last and is
// URL-escaped so legacy base64 signatures remain valid on the wire.
func (c *BaseHTTPClient) buildSignedString(params map[string]interface{}, signature string, forURL bool) string {
	keys := c.sortedKeys(params)
	parts := make([]string, 0, len(keys)+1)

	for _, k := range keys {
		v, ok := c.paramValueToString(params[k])
		if !ok {
			continue
		}
		if forURL && (strings.Contains(v, "[") || strings.Contains(v, "{")) {
			v = url.QueryEscape(v)
		}
		parts = append(parts, k+"="+v)
	}

	parts = append(parts, "signature="+url.QueryEscape(signature))
	return strings.Join(parts, "&")
}

// signString performs the signString operation.
func (c *BaseHTTPClient) signString(str string) string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(str))

	// Hex is the only encoding compatible with BingX's signature authentication.
	// Keep base64 as an explicit backward-compatible option, but never silently
	// fall back to base64 for empty or unrecognized encodings.
	if c.signatureEncoding == "base64" {
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	return hex.EncodeToString(h.Sum(nil))
}

// headers performs the headers operation.
func (c *BaseHTTPClient) headers() map[string]string {
	headers := map[string]string{
		"X-BX-APIKEY":  c.apiKey,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if c.sourceKey != "" {
		headers["X-SOURCE-KEY"] = c.sourceKey
	}

	return headers
}

// handleAPIError performs the handleAPIError operation.
func (c *BaseHTTPClient) handleAPIError(response map[string]interface{}) error {
	code, hasCode := response["code"]
	if !hasCode {
		return nil
	}

	codeStr := fmt.Sprintf("%v", code)
	if _, err := strconv.ParseInt(codeStr, 10, 64); err != nil {
		codeStr = "unknown"
	}
	message := "Unknown API error"
	if msg, ok := response["msg"].(string); ok {
		message = msg
	}
	if strings.Contains(strings.ToLower(message), "signature=") || strings.Contains(strings.ToLower(message), "signature%3d") {
		message = "Provider request rejected"
	}
	for _, secret := range []string{c.apiKey, c.apiSecret} {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, "[redacted]")
		}
	}

	switch codeStr {
	case "0":
		return nil
	case "100001", "100002", "100003", "100004", "100412":
		return errors.NewAuthenticationException(message, response)
	case "100005", "100429":
		return errors.NewRateLimitException(message, response)
	case "200001", "200002":
		return errors.NewInsufficientBalanceException(message, response)
	default:
		return errors.NewAPIException(message, codeStr, response)
	}
}

// Request performs the Request operation.
func (c *BaseHTTPClient) Request(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	receipt, err := c.requestBody(context.Background(), method, path, params)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(receipt.body, &data); err != nil {
		return nil, errors.NewBingXException("Invalid JSON response from API", 0, map[string]interface{}{"raw": string(receipt.body)})
	}

	return data, nil
}

// RequestJSON performs the RequestJSON operation.
func (c *BaseHTTPClient) RequestJSON(method, path string, params map[string]interface{}, result interface{}) error {
	receipt, err := c.requestBody(context.Background(), method, path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(receipt.body, result); err != nil {
		return errors.NewBingXException("Invalid JSON response from API", 0, map[string]interface{}{"raw": string(receipt.body)})
	}

	return nil
}

// RequestRaw executes one HTTP request with ctx and a bounded body read.
// Complete rejected responses return a receipt AND an error; incomplete or
// oversized responses return nil. The SDK does not retry requests.
func (c *BaseHTTPClient) RequestRaw(ctx context.Context, method, path string, params map[string]interface{}) (*RawResponse, error) {
	receipt, err := c.requestBody(ctx, method, path, params)
	if err != nil {
		if receipt != nil {
			return receipt.clone(), &ResponseError{cause: err, receipt: receipt}
		}
		return nil, err
	}
	return receipt.clone(), nil
}

// requestBody performs the requestBody operation.
func (c *BaseHTTPClient) requestBody(ctx context.Context, method, path string, params map[string]interface{}) (*RawResponse, error) {
	if err := c.validateLimit(); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	method = strings.ToUpper(method)

	// Do not add the generated timestamp to the caller's map. Reusing a map
	// across requests should produce a fresh timestamp and must not cause a
	// surprising mutation outside the client.
	requestParams := make(map[string]interface{}, len(params)+1)
	for key, value := range params {
		requestParams[key] = value
	}

	if _, exists := requestParams["timestamp"]; !exists {
		requestParams["timestamp"] = c.timestamp()
	}

	// The canonical string is the raw, sorted, URL-unencoded "key=value&..."
	// payload that is signed. It must never include the signature parameter.
	canonical := c.buildCanonicalString(requestParams)
	signature := c.signString(canonical)

	var req *http.Request
	var err error

	fullURL := c.baseURI + path

	if method == "GET" || method == "DELETE" {
		// Signature is appended last, after the canonical query. Values that
		// contain '[' or '{' are URL-escaped in the actual query string, while
		// the canonical signing string remains raw.
		query := c.buildSignedString(requestParams, signature, true)
		fullURL = fullURL + "?" + query
		req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
	} else {
		// POST/PUT bodies are form-urlencoded. The canonical string is sent as-is
		// (raw) followed by the signature.
		body := c.buildSignedString(requestParams, signature, false)
		req, err = http.NewRequestWithContext(ctx, method, fullURL, bytes.NewBufferString(body))
	}

	if err != nil {
		return nil, &RequestError{operation: "failed to create HTTP request", cause: err}
	}

	for k, v := range c.headers() {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{operation: "HTTP request failed", cause: err}
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, c.maxResponseBytes+1))
	if int64(len(body)) > c.maxResponseBytes {
		return nil, &ResponseTooLargeError{Limit: c.maxResponseBytes}
	}
	if err != nil {
		return nil, &IncompleteResponseError{cause: err}
	}
	retrievedAt := time.Now()
	if resp.ContentLength >= 0 && int64(len(body)) != resp.ContentLength {
		return nil, &IncompleteResponseError{cause: io.ErrUnexpectedEOF}
	}
	receipt := &RawResponse{
		body: body, completedAt: retrievedAt, statusCode: resp.StatusCode,
		headers: safeHeaders(resp.Header), route: safeRoute(path), selectors: c.safeSelectors(requestParams),
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return receipt, errors.NewBingXException("Invalid JSON response from API", 0, map[string]interface{}{"raw": string(body)})
	}

	if err := c.handleAPIError(data); err != nil {
		return receipt, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return receipt, errors.NewBingXException(
			fmt.Sprintf("HTTP request failed with status %d", resp.StatusCode),
			resp.StatusCode,
			data,
		)
	}

	return receipt, nil
}

// GetEndpoint performs the GetEndpoint operation.
func (c *BaseHTTPClient) GetEndpoint() string {
	return c.baseURI
}

// GetAPIKey performs the GetAPIKey operation.
func (c *BaseHTTPClient) GetAPIKey() string {
	return c.apiKey
}
