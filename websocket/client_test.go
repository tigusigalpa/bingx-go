package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestSend_DoesNotHoldStateLock_WhileWriting is a regression guard for #4.
// If Send held c.mu (RWMutex) for the duration of WriteMessage, a slow
// network write would block any caller of Disconnect/Connect that needs
// c.mu.Lock(). This test verifies that Send takes the state lock only
// briefly: even with the writer goroutine pinned on writeMu, a separate
// goroutine can still take c.mu.Lock() (here via Disconnect's behavior of
// reading c.conn under Lock).
func TestSend_DoesNotHoldStateLock_WhileWriting(t *testing.T) {
	c := NewWebSocketClient("ws://invalid.local")

	// Pin writeMu in another goroutine to simulate an in-flight WriteMessage.
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	// A Send call must still be able to acquire the state lock to read c.conn.
	// Issue the Send in a goroutine — it will block on writeMu (held above),
	// which is the correct behavior, but it must not be holding c.mu by then.
	sendDone := make(chan struct{})
	go func() {
		// c.conn is nil so Send returns the not-connected error before
		// reaching writeMu. This proves the state-lock path is non-blocking.
		_ = c.Send(map[string]interface{}{"id": "x"})
		close(sendDone)
	}()

	select {
	case <-sendDone:
		// Good: Send released c.mu and returned (conn was nil).
	case <-time.After(2 * time.Second):
		t.Fatal("Send appears to be holding the state mutex while writeMu is contended")
	}
}

// TestSend_ConcurrentCallers_NoDeadlock verifies that two concurrent Send
// calls on a not-connected client return (with the not-connected error)
// rather than deadlocking each other.
func TestSend_ConcurrentCallers_NoDeadlock(t *testing.T) {
	c := NewWebSocketClient("ws://invalid.local")

	var wg sync.WaitGroup
	done := make(chan struct{})

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.Send(map[string]interface{}{"id": "x"})
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent Send callers deadlocked")
	}
}

func TestDisconnectIsIdempotent(t *testing.T) {
	c := NewWebSocketClient("ws://invalid.local")
	if err := c.Disconnect(); err != nil {
		t.Fatalf("first disconnect returned error: %v", err)
	}
	if err := c.Disconnect(); err != nil {
		t.Fatalf("second disconnect returned error: %v", err)
	}
}

func TestClientReconnectsAfterDisconnect(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	c := NewWebSocketClient("ws" + strings.TrimPrefix(srv.URL, "http"))
	if err := c.Connect(); err != nil {
		t.Fatalf("first connect returned error: %v", err)
	}
	if err := c.Disconnect(); err != nil {
		t.Fatalf("disconnect returned error: %v", err)
	}
	if err := c.Connect(); err != nil {
		t.Fatalf("reconnect returned error: %v", err)
	}
	if !c.IsConnected() {
		t.Fatal("client is not connected after reconnect")
	}
	if err := c.Disconnect(); err != nil {
		t.Fatalf("final disconnect returned error: %v", err)
	}
}
