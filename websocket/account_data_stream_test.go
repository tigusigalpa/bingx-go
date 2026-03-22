package websocket

import (
	"strings"
	"testing"
)

func TestNewAccountDataStream(t *testing.T) {
	listenKey := "test-listen-key-123"
	stream := NewAccountDataStream(listenKey)

	if stream == nil {
		t.Fatal("Expected AccountDataStream to be created, got nil")
	}

	if stream.WebSocketClient == nil {
		t.Error("WebSocketClient should not be nil")
	}
}

func TestAccountDataStreamURL(t *testing.T) {
	listenKey := "test-key"
	stream := NewAccountDataStream(listenKey)

	if stream == nil {
		t.Fatal("Stream should not be nil")
	}

	expectedURL := AccountDataStreamBaseURL + "?listenKey=" + listenKey
	if !strings.Contains(expectedURL, listenKey) {
		t.Errorf("Expected URL to contain listen key %s", listenKey)
	}
}

func TestOnAccountUpdate(t *testing.T) {
	stream := NewAccountDataStream("test-key")

	callbackCalled := false
	stream.OnAccountUpdate(func(eventType string, data map[string]interface{}) {
		callbackCalled = true
	})

	if stream == nil {
		t.Error("Stream should not be nil")
	}

	if callbackCalled {
		t.Error("Callback should not be called without connection")
	}
}

func TestOnBalanceUpdate(t *testing.T) {
	stream := NewAccountDataStream("test-key")

	callbackCalled := false
	stream.OnBalanceUpdate(func(balances interface{}) {
		callbackCalled = true
	})

	if stream == nil {
		t.Error("Stream should not be nil")
	}

	if callbackCalled {
		t.Error("Callback should not be called without connection")
	}
}

func TestOnPositionUpdate(t *testing.T) {
	stream := NewAccountDataStream("test-key")

	callbackCalled := false
	stream.OnPositionUpdate(func(positions interface{}) {
		callbackCalled = true
	})

	if stream == nil {
		t.Error("Stream should not be nil")
	}

	if callbackCalled {
		t.Error("Callback should not be called without connection")
	}
}

func TestOnOrderUpdate(t *testing.T) {
	stream := NewAccountDataStream("test-key")

	callbackCalled := false
	stream.OnOrderUpdate(func(order interface{}) {
		callbackCalled = true
	})

	if stream == nil {
		t.Error("Stream should not be nil")
	}

	if callbackCalled {
		t.Error("Callback should not be called without connection")
	}
}
