package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	bingxerr "github.com/tigusigalpa/bingx-go/v2/errors"
)

// All responses in this file are synthetic loopback/transport fixtures.
type roundTripFunc func(*stdhttp.Request) (*stdhttp.Response, error)

func (f roundTripFunc) RoundTrip(r *stdhttp.Request) (*stdhttp.Response, error) { return f(r) }

func newLimitedClient(base string, limit int64, transport stdhttp.RoundTripper) *BaseHTTPClient {
	options := []ClientOption{WithMaxResponseBytes(limit)}
	if transport != nil {
		options = append(options, WithHTTPClient(&stdhttp.Client{Transport: transport}))
	}
	return NewBaseHTTPClientWithOptions("test-key", "test-secret", base, "", "hex", options...)
}

func TestBoundedResponsesAtLimitAndLimitPlusOne(t *testing.T) {
	const limit = 128
	for _, status := range []int{200, 503} {
		for _, code := range []int{0, 100429} {
			for _, size := range []int{limit, limit + 1} {
				t.Run(fmt.Sprintf("status%d/code%d/bytes%d", status, code, size), func(t *testing.T) {
					prefix := fmt.Sprintf(`{"code":%d,"data":"`, code)
					body := prefix + strings.Repeat("x", size-len(prefix)-2) + `"}`
					server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
						w.WriteHeader(status)
						_, _ = io.WriteString(w, body)
					}))
					defer server.Close()
					client := newLimitedClient(server.URL, limit, nil)
					receipt, err := client.RequestRaw(context.Background(), "GET", "/bounded", nil)
					if size > limit {
						var oversized *ResponseTooLargeError
						if receipt != nil || !errors.Is(err, ErrResponseTooLarge) || !errors.As(err, &oversized) || oversized.Limit != limit {
							t.Fatalf("expected oversized body rejection, got %v / %v", receipt, err)
						}
						return
					}
					if receipt == nil || string(receipt.BodyBytes()) != body || receipt.StatusCode() != status {
						t.Fatalf("exact-limit receipt lost: %v / %v", receipt, err)
					}
					if status == 200 && code == 0 {
						if err != nil {
							t.Fatal(err)
						}
					} else {
						var responseError *ResponseError
						if !errors.As(err, &responseError) || string(responseError.Receipt().BodyBytes()) != body {
							t.Fatalf("complete error evidence missing: %v", err)
						}
						if code != 0 {
							var rateLimit *bingxerr.RateLimitException
							if !errors.As(err, &rateLimit) {
								t.Fatalf("provider error type lost: %v", err)
							}
						}
					}
				})
			}
		}
	}
}

func TestStreamingBodyStopsAtLimit(t *testing.T) {
	const limit = 64 * 1024
	var sent atomic.Int64
	stopped := make(chan struct{})
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		defer close(stopped)
		chunk := strings.Repeat("x", 4096)
		for i := 0; i < 1<<18; i++ { // potentially 1 GiB; never allocate the whole stream
			if _, err := io.WriteString(w, chunk); err != nil {
				return
			}
			sent.Add(int64(len(chunk)))
			w.(stdhttp.Flusher).Flush()
			select {
			case <-r.Context().Done():
				return
			default:
			}
		}
	}))
	defer server.Close()
	client := newLimitedClient(server.URL, limit, nil)
	receipt, err := client.RequestRaw(context.Background(), "GET", "/stream", nil)
	if receipt != nil || !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("unexpected result: %v / %v", receipt, err)
	}
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("oversized stream was not closed")
	}
	if sent.Load() >= 1<<30 {
		t.Fatal("SDK consumed the full oversized stream")
	}
}

type countingBody struct {
	count  int64
	closed bool
}

func (b *countingBody) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}
	b.count += int64(len(p))
	return len(p), nil
}
func (b *countingBody) Close() error { b.closed = true; return nil }

func TestReaderConsumesOnlyLimitPlusOne(t *testing.T) {
	const limit int64 = 4096
	body := &countingBody{}
	client := newLimitedClient("http://loopback.invalid", limit, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{StatusCode: 503, Body: body, ContentLength: -1, Header: make(stdhttp.Header)}, nil
	}))
	_, err := client.RequestRaw(context.Background(), "GET", "/stream", nil)
	if !errors.Is(err, ErrResponseTooLarge) || body.count != limit+1 || !body.closed {
		t.Fatalf("read boundary failed: count=%d closed=%v err=%v", body.count, body.closed, err)
	}
}

func TestTruncatedResponseNeverReturnsReceipt(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = io.WriteString(w, `{"code":0}`)
	}))
	defer server.Close()
	client := newLimitedClient(server.URL, 128, nil)
	receipt, err := client.RequestRaw(context.Background(), "GET", "/truncated", nil)
	var incomplete *IncompleteResponseError
	if receipt != nil || !errors.As(err, &incomplete) || !errors.Is(err, ErrIncompleteResponse) || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated receipt accepted: %v / %v", receipt, err)
	}
}

func TestMalformedJSONRetainsCompleteReceipt(t *testing.T) {
	for _, status := range []int{200, 502} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			const body = `<html>synthetic gateway error</html>`
			server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				w.WriteHeader(status)
				_, _ = io.WriteString(w, body)
			}))
			defer server.Close()
			receipt, err := newLimitedClient(server.URL, 128, nil).RequestRaw(context.Background(), "GET", "/invalid", nil)
			var responseError *ResponseError
			if !errors.As(err, &responseError) || receipt == nil || string(responseError.Receipt().BodyBytes()) != body || receipt.StatusCode() != status {
				t.Fatalf("missing malformed response evidence: %v / %v", receipt, err)
			}
		})
	}
}

func TestReceiptImmutableAccessorsAndRepeatedRequests(t *testing.T) {
	const payload = ` {"code":0,"data":{"value":"1.230000","id":9007199254740993}} `
	var calls atomic.Int64
	var completed atomic.Int64
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Retry-After", "3")
		w.Header().Set("Set-Cookie", "secret-cookie")
		w.Header().Set("Location", "https://example.com/?signature=secret")
		w.Header().Set("X-BX-APIKEY", "secret-api-key")
		_, _ = io.WriteString(w, payload[:10])
		w.(stdhttp.Flusher).Flush()
		time.Sleep(time.Millisecond)
		completed.Store(time.Now().UnixNano())
		_, _ = io.WriteString(w, payload[10:])
	}))
	defer server.Close()
	client := newLimitedClient(server.URL, 128, nil)
	params := map[string]interface{}{"symbol": "BTC-USDT", "startTime": int64(123), "limit": 10, "apiKey": "secret-selector", "signature": "secret-signature"}
	receipt, err := client.RequestRaw(context.Background(), "GET", "/receipt", params)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.CompletedAt().UnixNano() < completed.Load() || receipt.CompletedAt().After(time.Now()) {
		t.Fatal("timestamp was not recorded at body completion")
	}
	if receipt.Route() != "/receipt" || receipt.Selectors()["symbol"] != "BTC-USDT" || receipt.Selectors()["startTime"] != "123" {
		t.Fatal("safe request provenance missing")
	}
	for _, secret := range []string{"apiKey", "signature", "timestamp"} {
		if _, ok := receipt.Selectors()[secret]; ok {
			t.Fatalf("unsafe selector %s", secret)
		}
	}
	for _, key := range []string{"Set-Cookie", "Location", "X-BX-APIKEY"} {
		if receipt.Headers().Get(key) != "" {
			t.Fatalf("unsafe header %s", key)
		}
	}
	if receipt.Headers().Get("Retry-After") != "3" || receipt.Headers().Get("Content-Type") != "application/json" {
		t.Fatal("safe headers missing")
	}
	bytes := receipt.BodyBytes()
	bytes[0] = 'x'
	receipt.Body[0] = 'y'
	receipt.RetrievedAt = time.Time{}
	selectors := receipt.Selectors()
	selectors["symbol"] = "changed"
	headers := receipt.Headers()
	headers.Set("Retry-After", "99")
	params["symbol"] = "ETH-USDT"
	if string(receipt.BodyBytes()) != payload || receipt.CompletedAt().IsZero() || receipt.Selectors()["symbol"] != "BTC-USDT" || receipt.Headers().Get("Retry-After") != "3" {
		t.Fatal("receipt accessor data mutated")
	}
	second, err := client.RequestRaw(context.Background(), "GET", "/receipt", params)
	if err != nil || string(second.Body) != payload || calls.Load() != 2 || second.Selectors()["symbol"] != "ETH-USDT" {
		t.Fatalf("repeated request failed: %v / %v", second, err)
	}
}

func TestCancellationAndDeadlineDuringBody(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		t.Run(fmt.Sprint(deadline), func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				_, _ = io.WriteString(w, `{"code":0,"data":"`)
				w.(stdhttp.Flusher).Flush()
				<-r.Context().Done()
			}))
			defer server.Close()
			ctx, cancel := context.WithCancel(context.Background())
			want := context.Canceled
			if deadline {
				ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
				want = context.DeadlineExceeded
			}
			defer cancel()
			transport := stdhttp.DefaultTransport.(*stdhttp.Transport).Clone()
			defer transport.CloseIdleConnections()
			client := newLimitedClient(server.URL, 128, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
				response, err := transport.RoundTrip(r)
				if err == nil {
					response.Body = &notifyingBody{ReadCloser: response.Body, started: started}
				}
				return response, err
			}))
			done := make(chan error, 1)
			go func() {
				receipt, err := client.RequestRaw(ctx, "GET", "/slow-body", nil)
				if receipt != nil {
					t.Error("partial receipt returned on cancellation")
				}
				done <- err
			}()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("body read never started")
			}
			if !deadline {
				cancel()
			}
			select {
			case err := <-done:
				if !errors.Is(err, want) || !errors.Is(err, ErrIncompleteResponse) {
					t.Fatalf("cancellation identity lost: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("body read did not cancel")
			}
		})
	}
}

type notifyingBody struct {
	io.ReadCloser
	started chan struct{}
	once    sync.Once
}

func (b *notifyingBody) Read(p []byte) (int, error) {
	b.once.Do(func() { close(b.started) })
	return b.ReadCloser.Read(p)
}

func TestCustomTransportIsolationAndLegacyMethods(t *testing.T) {
	var customCalls, otherCalls atomic.Int64
	transport := roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		customCalls.Add(1)
		return &stdhttp.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"code":0}`)), ContentLength: -1, Header: make(stdhttp.Header)}, nil
	})
	originalDefault := stdhttp.DefaultTransport
	injected := &stdhttp.Client{Transport: transport, Timeout: time.Second}
	custom := NewBaseHTTPClientWithOptions("key", "secret", "http://custom.invalid", "", "hex", WithHTTPClient(injected), WithMaxResponseBytes(32))
	injected.Transport = roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) { return nil, errors.New("mutated client") })
	if data, err := custom.Request("GET", "/legacy", nil); err != nil || data["code"] != float64(0) {
		t.Fatalf("legacy Request failed: %v", err)
	}
	var result map[string]interface{}
	if err := custom.RequestJSON("POST", "/typed", nil, &result); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		otherCalls.Add(1)
		_, _ = io.WriteString(w, `{"code":0}`)
	}))
	defer server.Close()
	if _, err := NewBaseHTTPClient("key", "secret", server.URL, "", "hex").Request("GET", "/default", nil); err != nil {
		t.Fatal(err)
	}
	if stdhttp.DefaultTransport != originalDefault || customCalls.Load() != 2 || otherCalls.Load() != 1 {
		t.Fatal("custom client leaked across instances")
	}
}

func TestInvalidLimitsFailBeforeIO(t *testing.T) {
	for _, limit := range []int64{0, -1, math.MaxInt64} {
		client := newLimitedClient("http://custom.invalid", limit, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
			t.Fatal("invalid limit performed I/O")
			return nil, nil
		}))
		if receipt, err := client.RequestRaw(context.Background(), "GET", "/invalid", nil); receipt != nil || !errors.Is(err, ErrInvalidResponseLimit) {
			t.Fatalf("invalid limit %d accepted: %v", limit, err)
		}
	}
}

func TestTransportDiagnosticsDoNotExposeSignedURL(t *testing.T) {
	client := newLimitedClient("http://custom.invalid", 32, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		return nil, fmt.Errorf("synthetic transport: %s key=%s", r.URL, r.Header.Get("X-BX-APIKEY"))
	}))
	_, err := client.RequestRaw(context.Background(), "GET", "/private", nil)
	if err == nil {
		t.Fatal("expected transport error")
	}
	for _, diagnostic := range []string{err.Error(), fmt.Sprintf("%+v", err), fmt.Sprintf("%#v", err)} {
		if strings.Contains(diagnostic, "signature=") || strings.Contains(diagnostic, "test-key") {
			t.Fatalf("sensitive diagnostic: %s", diagnostic)
		}
	}
}

func TestErrorReceiptIsIndependent(t *testing.T) {
	const body = `{"code":100429,"msg":"try later"}`
	transport := roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{StatusCode: 429, Body: io.NopCloser(strings.NewReader(body)), ContentLength: -1, Header: stdhttp.Header{"Retry-After": []string{"2"}}}, nil
	})
	raw, err := newLimitedClient("http://custom.invalid", 128, transport).RequestRaw(context.Background(), "GET", "/error", nil)
	var responseError *ResponseError
	if !errors.As(err, &responseError) {
		t.Fatalf("expected ResponseError: %v", err)
	}
	raw.Body[0] = 'x'
	one := responseError.Receipt()
	one.Body[0] = 'y'
	one.Headers().Set("Retry-After", "99")
	if string(responseError.Receipt().BodyBytes()) != body || responseError.Receipt().Headers().Get("Retry-After") != "2" {
		t.Fatal("error evidence mutated")
	}
}

func TestProviderDiagnosticsRedactSecrets(t *testing.T) {
	client := NewBaseHTTPClient("test-key", "test-secret", "http://custom.invalid", "", "hex")
	for _, message := range []string{"bad test-key test-secret", "https://example.com/?signature=secret-value"} {
		err := client.handleAPIError(map[string]interface{}{"code": 999, "msg": message})
		if err == nil || strings.Contains(err.Error(), "test-key") || strings.Contains(err.Error(), "test-secret") || strings.Contains(err.Error(), "signature=") {
			t.Fatalf("unsafe provider diagnostic: %v", err)
		}
	}
}

func TestMetadataAllowlists(t *testing.T) {
	client := NewBaseHTTPClient("private-key", "private-secret", "http://custom.invalid", "", "hex")
	selectors := client.safeSelectors(map[string]interface{}{"symbol": "private-key", "interval": "1m?signature=secret", "limit": "invalid", "fromId": int64(0), "period": strings.Repeat("x", 65), "startTime": int64(1), "endTime": int64(2)})
	if len(selectors) != 3 || selectors["fromId"] != "0" || selectors["startTime"] != "1" || selectors["endTime"] != "2" {
		t.Fatalf("unsafe selectors retained: %#v", selectors)
	}
	for _, path := range []string{"/quote?signature=secret#fragment", "https://userinfo@example.com/quote?signature=secret"} {
		if safeRoute(path) != "/quote" {
			t.Fatalf("unsafe route: %s", safeRoute(path))
		}
	}
	if safeRoute("%invalid") != "" {
		t.Fatal("invalid route retained")
	}
	headers := safeHeaders(stdhttp.Header{"Date": []string{"invalid"}, "Retry-After": []string{"signature=secret"}, "Content-Length": []string{"-1"}, "Content-Type": []string{"invalid;"}, "Content-Encoding": []string{"private-key"}})
	if len(headers) != 0 {
		t.Fatalf("unsafe header retained: %#v", headers)
	}
	headers = safeHeaders(stdhttp.Header{"Date": []string{"Wed, 07 Oct 2026 12:00:00 GMT"}, "Retry-After": []string{"Wed, 07 Oct 2026 12:00:00 GMT"}, "Content-Length": []string{"10"}, "Content-Encoding": []string{"gzip"}})
	if len(headers) != 4 {
		t.Fatalf("valid header lost: %#v", headers)
	}
}

func TestLegacyProviderTypeRemainsDirect(t *testing.T) {
	transport := roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{StatusCode: 429, ContentLength: -1, Body: io.NopCloser(strings.NewReader(`{"code":100429,"msg":"try later"}`)), Header: make(stdhttp.Header)}, nil
	})
	client := newLimitedClient("http://custom.invalid", 128, transport)
	_, err := client.Request("GET", "/legacy-error", nil)
	if _, ok := err.(*bingxerr.RateLimitException); !ok {
		t.Fatalf("legacy exception was wrapped: %T", err)
	}
}

func TestReceiptDiagnosticsOmitPayload(t *testing.T) {
	receipt := RawResponse{Body: []byte("test-key signature=secret"), body: []byte("test-key signature=secret"), statusCode: 200}
	for _, diagnostic := range []string{fmt.Sprintf("%v", receipt), fmt.Sprintf("%+v", &receipt), fmt.Sprintf("%#v", receipt), fmt.Sprintf("%#v", &receipt)} {
		if strings.Contains(diagnostic, "test-key") || strings.Contains(diagnostic, "signature=") {
			t.Fatalf("unsafe receipt diagnostic: %s", diagnostic)
		}
	}
}

func TestAllBoundaryErrorDiagnosticsAreSafe(t *testing.T) {
	const secret = "test-key signature=secret"
	for _, err := range []error{
		&ResponseTooLargeError{Limit: 128},
		&IncompleteResponseError{cause: errors.New(secret)},
		&ResponseError{cause: errors.New(secret), receipt: &RawResponse{body: []byte(secret), statusCode: 503}},
		bingxerr.NewBingXException("Invalid JSON response from API", 0, map[string]interface{}{"raw": secret}),
		bingxerr.NewAPIException("Provider request rejected", "100001", map[string]interface{}{"raw": secret}),
	} {
		for _, diagnostic := range []string{err.Error(), fmt.Sprintf("%+v", err), fmt.Sprintf("%#v", err)} {
			if diagnostic == "" || strings.Contains(diagnostic, "test-key") || strings.Contains(diagnostic, "signature=") {
				t.Fatalf("unsafe boundary diagnostic: %s", diagnostic)
			}
		}
	}
}

func TestDeclaredLengthMismatchAndRequestCreationError(t *testing.T) {
	client := newLimitedClient("http://custom.invalid", 128, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{StatusCode: 200, ContentLength: 100, Body: io.NopCloser(strings.NewReader(`{"code":0}`)), Header: make(stdhttp.Header)}, nil
	}))
	if receipt, err := client.RequestRaw(context.Background(), "GET", "/mismatch", nil); receipt != nil || !errors.Is(err, ErrIncompleteResponse) {
		t.Fatalf("declared mismatch accepted: %v / %v", receipt, err)
	}
	client = NewBaseHTTPClient("test-key", "secret", "://invalid", "", "hex")
	if receipt, err := client.RequestRaw(context.Background(), "GET", "/bad", nil); receipt != nil || err == nil || strings.Contains(err.Error(), "signature=") {
		t.Fatalf("invalid URL not safely rejected: %v / %v", receipt, err)
	}
}

type failingBoundaryBody struct {
	cause  error
	closed bool
}

func (b *failingBoundaryBody) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = ' '
	}
	return len(p), b.cause
}

func (b *failingBoundaryBody) Close() error { b.closed = true; return nil }

func TestReadErrorIdentityAtSizeBoundary(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded, io.ErrUnexpectedEOF} {
		t.Run(cause.Error(), func(t *testing.T) {
			body := &failingBoundaryBody{cause: cause}
			client := newLimitedClient("http://custom.invalid", 128, roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{StatusCode: 200, ContentLength: -1, Body: body, Header: make(stdhttp.Header)}, nil
			}))
			receipt, err := client.RequestRaw(context.Background(), "GET", "/boundary", nil)
			if receipt != nil || !errors.Is(err, ErrIncompleteResponse) || !errors.Is(err, cause) || !body.closed {
				t.Fatalf("boundary read masked its error or returned evidence: %v / %v", receipt, err)
			}
		})
	}
}
