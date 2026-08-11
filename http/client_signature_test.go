package http

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestGET_TimestampOnly(t *testing.T) {
	// Known test vector for "timestamp=1702731500000" signed with test-secret.
	expectedSig := expectedHexSignature("timestamp=1702731500000")

	var rawQuery string
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		rawQuery = r.URL.RawQuery

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("failed to parse raw query: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")
		if len(values) != 1 || values.Get("timestamp") != "1702731500000" {
			t.Errorf("unexpected query values: %v", values)
		}
		if sig != expectedSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, expectedSig)
		}
		if !strings.HasSuffix(rawQuery, "&signature="+expectedSig) {
			t.Errorf("signature is not last in raw query: %s", rawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	_, err := client.Request("GET", "/openApi/swap/v3/user/balance", map[string]interface{}{
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "timestamp=1702731500000&signature=" + expectedSig
	if rawQuery != want {
		t.Errorf("raw query = %q, want %q", rawQuery, want)
	}
}

func TestRequestGET_WithSymbol(t *testing.T) {
	expectedSig := expectedHexSignature("symbol=BTC-USDT&timestamp=1702731500000")

	var rawQuery string
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		rawQuery = r.URL.RawQuery

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("failed to parse raw query: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")
		if values.Get("symbol") != "BTC-USDT" || values.Get("timestamp") != "1702731500000" {
			t.Errorf("unexpected query values: %v", values)
		}
		if sig != expectedSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, expectedSig)
		}
		if !strings.HasSuffix(rawQuery, "&signature="+expectedSig) {
			t.Errorf("signature is not last in raw query: %s", rawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	_, err := client.Request("GET", "/test", map[string]interface{}{
		"symbol":    "BTC-USDT",
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "symbol=BTC-USDT&timestamp=1702731500000&signature=" + expectedSig
	if rawQuery != want {
		t.Errorf("raw query = %q, want %q", rawQuery, want)
	}
}

func TestRequestPOST_BodyIsCanonicalPlusSignature(t *testing.T) {
	expectedSig := expectedHexSignature("symbol=BTC-USDT&timestamp=1702731500000")

	var gotBody string
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)

		values, err := url.ParseQuery(gotBody)
		if err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")
		if values.Get("symbol") != "BTC-USDT" || values.Get("timestamp") != "1702731500000" {
			t.Errorf("unexpected body values: %v", values)
		}
		if sig != expectedSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, expectedSig)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	_, err := client.Request("POST", "/test", map[string]interface{}{
		"symbol":    "BTC-USDT",
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "symbol=BTC-USDT&timestamp=1702731500000&signature=" + expectedSig
	if gotBody != want {
		t.Errorf("body = %q, want %q", gotBody, want)
	}
}

func TestRequestGET_BatchParamValueIsURLEncodedInQueryButRawInSignature(t *testing.T) {
	orders := `[{"symbol":"BTC-USDT","side":"Buy","type":"Limit","price":"30000","qty":"0.01"}]`
	canonical := "orders=" + orders + "&timestamp=1702731500000"
	expectedSig := expectedHexSignature(canonical)

	var rawQuery string
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		rawQuery = r.URL.RawQuery

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("failed to parse raw query: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")
		if values.Get("orders") != orders {
			t.Errorf("decoded orders value = %q, want %q", values.Get("orders"), orders)
		}
		if values.Get("timestamp") != "1702731500000" {
			t.Errorf("decoded timestamp = %q, want %s", values.Get("timestamp"), "1702731500000")
		}
		if sig != expectedSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, expectedSig)
		}
		if !strings.HasSuffix(rawQuery, "&signature="+expectedSig) {
			t.Errorf("signature is not last in raw query: %s", rawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	_, err := client.Request("GET", "/test", map[string]interface{}{
		"orders":    orders,
		"timestamp": "1702731500000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The raw query must URL-encode the JSON value because it contains [ and {,
	// while the canonical/signature remain raw.
	if strings.Contains(rawQuery, "orders=[") {
		t.Errorf("raw query left orders value unencoded: %s", rawQuery)
	}
	if !strings.Contains(rawQuery, "orders=%5B") {
		t.Errorf("raw query does not contain URL-encoded orders value: %s", rawQuery)
	}
}

func TestRequestDoesNotMutateParams(t *testing.T) {
	params := map[string]interface{}{"symbol": "BTC-USDT"}
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	if _, err := client.Request("GET", "/test", params); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := params["timestamp"]; ok {
		t.Fatal("Request must not add timestamp to the caller's params")
	}
}

func TestRequestReturnsErrorForHTTPFailureWithoutAPICode(t *testing.T) {
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusInternalServerError)
		w.Write([]byte(`{"message":"temporary failure"}`))
	}))
	defer srv.Close()

	client := NewBaseHTTPClient("test-key", testSignatureSecret, srv.URL, "", "hex")
	if _, err := client.Request("GET", "/test", nil); err == nil {
		t.Fatal("expected an error for an HTTP 500 response")
	}
}
