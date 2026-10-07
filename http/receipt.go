package http

import (
	"errors"
	"fmt"
	"mime"
	stdhttp "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultMaxResponseBytes bounds every response body, including error bodies.
const DefaultMaxResponseBytes int64 = 16 << 20

// ErrResponseTooLarge identifies a response exceeding the configured byte limit.
var ErrResponseTooLarge = errors.New("response body exceeds byte limit")

// ErrIncompleteResponse identifies a body that could not be read to completion.
var ErrIncompleteResponse = errors.New("response body is incomplete")

// ErrInvalidResponseLimit identifies nonpositive or overflowing response limits.
var ErrInvalidResponseLimit = errors.New("maximum response bytes must be positive and less than MaxInt64")

// ResponseTooLargeError rejects a body without retaining partial response bytes.
type ResponseTooLargeError struct {
	Limit int64
}

func (e *ResponseTooLargeError) Error() string {
	return fmt.Sprintf("response body exceeds %d bytes", e.Limit)
}

// Is supports errors.Is(err, ErrResponseTooLarge).
func (e *ResponseTooLargeError) Is(target error) bool { return target == ErrResponseTooLarge }

// IncompleteResponseError reports a failed body read, including cancellation.
type IncompleteResponseError struct {
	cause error
}

func (e *IncompleteResponseError) Error() string { return ErrIncompleteResponse.Error() }

// Unwrap preserves context cancellation, deadlines and I/O error identity.
func (e *IncompleteResponseError) Unwrap() error { return e.cause }

// Is supports errors.Is(err, ErrIncompleteResponse).
func (e *IncompleteResponseError) Is(target error) bool { return target == ErrIncompleteResponse }

// RequestError hides signed URLs in request/transport diagnostic messages.
type RequestError struct {
	operation string
	cause     error
}

func (e *RequestError) Error() string { return e.operation }

// Unwrap preserves errors.Is/As without copying a transport error into logs.
func (e *RequestError) Unwrap() error { return e.cause }

// RawResponse stores a complete, bounded HTTP response. Legacy exported fields
// are independent copies; use accessors for immutable receipt evidence.
type RawResponse struct {
	Body        []byte
	RetrievedAt time.Time
	body        []byte
	completedAt time.Time
	statusCode  int
	headers     stdhttp.Header
	route       string
	selectors   map[string]string
}

// BodyBytes returns a copy of the exact body, independent of legacy Body edits.
func (r *RawResponse) BodyBytes() []byte { return append([]byte(nil), r.body...) }

// CompletedAt returns local time immediately after a successful complete read.
// It is not the exchange's event time or HTTP Date.
func (r *RawResponse) CompletedAt() time.Time { return r.completedAt }

// StatusCode returns the HTTP status associated with the receipt.
func (r *RawResponse) StatusCode() int { return r.statusCode }

// Headers returns a copy of the allowlisted response headers.
func (r *RawResponse) Headers() stdhttp.Header { return r.headers.Clone() }

// Route returns the request path without query, origin, userinfo or fragment.
func (r *RawResponse) Route() string { return r.route }

// Selectors returns a copy of allowlisted public market selectors.
func (r *RawResponse) Selectors() map[string]string {
	result := make(map[string]string, len(r.selectors))
	for key, value := range r.selectors {
		result[key] = value
	}
	return result
}

// String omits the body and request parameters from diagnostic output.
func (r RawResponse) String() string {
	return fmt.Sprintf("HTTP receipt: status=%d bytes=%d", r.statusCode, len(r.body))
}

// GoString also keeps %#v diagnostics free of payloads and signed URLs.
func (r RawResponse) GoString() string { return r.String() }

// GoString keeps transport and body errors safe in %#v diagnostic output.
func (e RequestError) GoString() string { return e.operation }

// GoString omits the underlying reader error from diagnostics.
func (e IncompleteResponseError) GoString() string { return ErrIncompleteResponse.Error() }

// GoString omits provider payloads from diagnostics.
func (e ResponseError) GoString() string { return e.Error() }

func (r *RawResponse) clone() *RawResponse {
	if r == nil {
		return nil
	}
	return &RawResponse{
		Body: r.BodyBytes(), RetrievedAt: r.completedAt,
		body: r.body, completedAt: r.completedAt, statusCode: r.statusCode,
		headers: r.headers, route: r.route, selectors: r.selectors,
	}
}

// ResponseError retains complete bounded response evidence for raw API errors.
// errors.As can recover the original BingX provider exception via Unwrap.
type ResponseError struct {
	cause   error
	receipt *RawResponse
}

func (e *ResponseError) Error() string {
	if e.receipt == nil {
		return "HTTP response rejected"
	}
	return fmt.Sprintf("HTTP response rejected: status=%d", e.receipt.statusCode)
}

// Unwrap preserves the underlying provider or JSON error.
func (e *ResponseError) Unwrap() error { return e.cause }

// Receipt returns an independent view of complete bounded response evidence.
func (e *ResponseError) Receipt() *RawResponse { return e.receipt.clone() }

func safeRoute(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return ""
	}
	return u.EscapedPath()
}

func (c *BaseHTTPClient) safeSelectors(params map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for _, key := range []string{"symbol", "interval", "period", "limit", "fromId", "startTime", "endTime"} {
		value, ok := c.paramValueToString(params[key])
		if !ok || len(value) > 64 || value == c.apiKey || value == c.apiSecret {
			continue
		}
		if key == "symbol" || key == "interval" || key == "period" {
			if value == "" || strings.IndexFunc(value, func(r rune) bool {
				return (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '_'
			}) >= 0 {
				continue
			}
		} else if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			continue
		}
		result[key] = value
	}
	return result
}

func safeHeaders(headers stdhttp.Header) stdhttp.Header {
	result := make(stdhttp.Header)
	for _, key := range []string{"Date", "Retry-After", "Content-Type", "Content-Length", "Content-Encoding"} {
		value := headers.Get(key)
		if value == "" || len(value) > 256 {
			continue
		}
		switch key {
		case "Date":
			if _, err := stdhttp.ParseTime(value); err != nil {
				continue
			}
		case "Retry-After":
			if seconds, err := strconv.ParseUint(value, 10, 64); err != nil || seconds > 86400*365 {
				if _, err := stdhttp.ParseTime(value); err != nil {
					continue
				}
			}
		case "Content-Length":
			if _, err := strconv.ParseUint(value, 10, 64); err != nil {
				continue
			}
		case "Content-Type":
			mediaType, _, err := mime.ParseMediaType(value)
			if err != nil || strings.Count(mediaType, "/") != 1 {
				continue
			}
			value = mediaType // omit arbitrary MIME parameters
		case "Content-Encoding":
			if value != "gzip" && value != "identity" && value != "deflate" && value != "br" {
				continue
			}
		}
		result.Set(key, value)
	}
	return result
}
