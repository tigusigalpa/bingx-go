package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

const spotAccountTestSecret = "test-secret"

func expectedHexSignatureFor(input string) string {
	h := hmac.New(sha256.New, []byte(spotAccountTestSecret))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func newSpotAccountTestService(t *testing.T, handler http.HandlerFunc) (*SpotAccountService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := bxhttp.NewBaseHTTPClient("test-key", spotAccountTestSecret, srv.URL, "", "hex")
	return NewSpotAccountService(client), srv
}

func TestSpotAccountGetAccountOverview_NilAccountType(t *testing.T) {
	var gotPath, rawQuery string
	service, srv := newSpotAccountTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		rawQuery = r.URL.RawQuery

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("failed to parse raw query: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")

		ts := values.Get("timestamp")
		if ts == "" {
			t.Fatalf("timestamp is missing from query: %s", rawQuery)
		}
		if len(values) != 1 {
			t.Fatalf("expected only timestamp in query, got %v", values)
		}

		canonical := "timestamp=" + ts
		wantSig := expectedHexSignatureFor(canonical)
		if sig != wantSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, wantSig)
		}
		if !strings.HasSuffix(rawQuery, "&signature="+wantSig) {
			t.Errorf("signature is not last in raw query: %s", rawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"code":0,"data":[]}`)
	})
	defer srv.Close()

	_, err := service.GetAccountOverview(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/openApi/account/v1/allAccountBalance" {
		t.Errorf("path = %s, want /openApi/account/v1/allAccountBalance", gotPath)
	}
	if rawQuery == "" {
		t.Errorf("raw query should not be empty")
	}
}

func TestSpotAccountGetAccountOverview_WithAccountType(t *testing.T) {
	var gotPath, rawQuery string
	service, srv := newSpotAccountTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		rawQuery = r.URL.RawQuery

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("failed to parse raw query: %v", err)
		}

		sig := values.Get("signature")
		values.Del("signature")

		accountType := values.Get("accountType")
		ts := values.Get("timestamp")
		if accountType != AccountTypeUSDTMPerp {
			t.Errorf("accountType = %q, want %q", accountType, AccountTypeUSDTMPerp)
		}
		if ts == "" {
			t.Fatalf("timestamp is missing from query: %s", rawQuery)
		}

		canonical := "accountType=" + accountType + "&timestamp=" + ts
		wantSig := expectedHexSignatureFor(canonical)
		if sig != wantSig {
			t.Errorf("signature mismatch: got %s, want %s", sig, wantSig)
		}
		if !strings.HasSuffix(rawQuery, "&signature="+wantSig) {
			t.Errorf("signature is not last in raw query: %s", rawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"code":0,"data":[]}`)
	})
	defer srv.Close()

	_, err := service.GetAccountOverview(strPtr(AccountTypeUSDTMPerp))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/openApi/account/v1/allAccountBalance" {
		t.Errorf("path = %s, want /openApi/account/v1/allAccountBalance", gotPath)
	}
	if !strings.Contains(rawQuery, "accountType=USDTMPerp") {
		t.Errorf("raw query should contain accountType=USDTMPerp: %s", rawQuery)
	}
}

func TestSpotAccountGetFundBalance_Retired(t *testing.T) {
	called := false
	service, srv := newSpotAccountTestService(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"code":0,"data":[]}`)
	})
	defer srv.Close()

	_, err := service.GetFundBalance()
	if err == nil {
		t.Fatal("expected error from retired GetFundBalance, got nil")
	}
	if !strings.Contains(err.Error(), "GetFundBalance is retired") {
		t.Errorf("expected retirement error, got %q", err.Error())
	}
	if called {
		t.Error("GetFundBalance should not make an HTTP request")
	}
}
