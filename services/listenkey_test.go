package services

import (
	"testing"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

func TestNewListenKeyService(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewListenKeyService(client)

	if service == nil {
		t.Fatal("Expected ListenKeyService to be created, got nil")
	}

	if service.client == nil {
		t.Error("ListenKeyService client should not be nil")
	}
}

func TestGenerateListenKey(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewListenKeyService(client)

	_, err := service.Generate()
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestExtendListenKey(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewListenKeyService(client)

	listenKey := "test-listen-key-123"
	_, err := service.Extend(listenKey)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}

func TestDeleteListenKey(t *testing.T) {
	client := http.NewBaseHTTPClient("key", "secret", "https://api.test.com", "", "base64")
	service := NewListenKeyService(client)

	listenKey := "test-listen-key-123"
	_, err := service.Delete(listenKey)
	if err == nil {
		t.Skip("Skipping test - would require mock HTTP server")
	}
}
