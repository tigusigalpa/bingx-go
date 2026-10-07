package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestRawPreservesBodyAndContext(t *testing.T) {
	const body = `{"code":0,"data":{"price":"1.2300"}}`
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.URL.Path != "/raw" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewBaseHTTPClient("key", "secret", server.URL, "", "hex")
	response, err := client.RequestRaw(context.Background(), "GET", "/raw", nil)
	if err != nil {
		t.Fatalf("RequestRaw returned error: %v", err)
	}
	if string(response.Body) != body {
		t.Fatalf("unexpected raw body: %q", response.Body)
	}
	if response.RetrievedAt.IsZero() {
		t.Fatal("RetrievedAt is zero")
	}
}

func TestRequestRawHonorsCancellation(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewBaseHTTPClient("key", "secret", server.URL, "", "hex")
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		_, err := client.RequestRaw(ctx, "GET", "/slow", nil)
		errCh <- err
	}()
	<-started
	cancel()

	if err := <-errCh; !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}
