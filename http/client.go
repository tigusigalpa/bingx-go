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

	"github.com/tigusigalpa/bingx-go/errors"
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

func (c *BaseHTTPClient) buildQuery(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := url.Values{}
	for _, k := range keys {
		v := params[k]
		switch val := v.(type) {
		case string:
			values.Add(k, val)
		case int:
			values.Add(k, strconv.Itoa(val))
		case int64:
			values.Add(k, strconv.FormatInt(val, 10))
		case float64:
			values.Add(k, strconv.FormatFloat(val, 'f', -1, 64))
		case bool:
			values.Add(k, strconv.FormatBool(val))
		default:
			values.Add(k, fmt.Sprintf("%v", val))
		}
	}

	return values.Encode()
}

func (c *BaseHTTPClient) signString(str string) string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(str))

	if c.signatureEncoding == "hex" {
		return hex.EncodeToString(h.Sum(nil))
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
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
	case "100001", "100002", "100003", "100004":
		return errors.NewAuthenticationException(message, response)
	case "100005":
		return errors.NewRateLimitException(message, response)
	case "200001":
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

	query := c.buildQuery(params)
	signature := c.signString(query)

	var req *http.Request
	var err error

	fullURL := c.baseURI + path

	if method == "GET" || method == "DELETE" {
		params["signature"] = signature
		queryWithSig := c.buildQuery(params)
		fullURL = fullURL + "?" + queryWithSig
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		params["signature"] = signature
		formData := c.buildQuery(params)
		req, err = http.NewRequest(method, fullURL, bytes.NewBufferString(formData))
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
