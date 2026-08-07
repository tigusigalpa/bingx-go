package http

import (
	"bytes"
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

type BaseHTTPClient struct {
	apiKey            string
	apiSecret         string
	baseURI           string
	sourceKey         string
	signatureEncoding string
	httpClient        *http.Client
}

func NewBaseHTTPClient(apiKey, apiSecret, baseURI, sourceKey, signatureEncoding string) *BaseHTTPClient {
	return &BaseHTTPClient{
		apiKey:            apiKey,
		apiSecret:         apiSecret,
		baseURI:           baseURI,
		sourceKey:         sourceKey,
		signatureEncoding: signatureEncoding,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *BaseHTTPClient) timestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

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

func (c *BaseHTTPClient) handleAPIError(response map[string]interface{}) error {
	code, hasCode := response["code"]
	if !hasCode {
		return nil
	}

	codeStr := fmt.Sprintf("%v", code)
	message := "Unknown API error"
	if msg, ok := response["msg"].(string); ok {
		message = msg
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

func (c *BaseHTTPClient) Request(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	method = strings.ToUpper(method)

	if params == nil {
		params = make(map[string]interface{})
	}

	if _, exists := params["timestamp"]; !exists {
		params["timestamp"] = c.timestamp()
	}

	// The canonical string is the raw, sorted, URL-unencoded "key=value&..."
	// payload that is signed. It must never include the signature parameter.
	canonical := c.buildCanonicalString(params)
	signature := c.signString(canonical)

	var req *http.Request
	var err error

	fullURL := c.baseURI + path

	if method == "GET" || method == "DELETE" {
		// Signature is appended last, after the canonical query. Values that
		// contain '[' or '{' are URL-escaped in the actual query string, while
		// the canonical signing string remains raw.
		query := c.buildSignedString(params, signature, true)
		fullURL = fullURL + "?" + query
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		// POST/PUT bodies are form-urlencoded. The canonical string is sent as-is
		// (raw) followed by the signature.
		body := c.buildSignedString(params, signature, false)
		req, err = http.NewRequest(method, fullURL, bytes.NewBufferString(body))
	}

	if err != nil {
		return nil, errors.NewBingXException("Failed to create request: "+err.Error(), 0, nil)
	}

	for k, v := range c.headers() {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewBingXException("HTTP request failed: "+err.Error(), 0, nil)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewBingXException("Failed to read response: "+err.Error(), 0, nil)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.NewBingXException("Invalid JSON response from API", 0, map[string]interface{}{"raw": string(body)})
	}

	if err := c.handleAPIError(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *BaseHTTPClient) GetEndpoint() string {
	return c.baseURI
}

func (c *BaseHTTPClient) GetAPIKey() string {
	return c.apiKey
}
