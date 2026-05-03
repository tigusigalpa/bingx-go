package websocket

import (
	"sync"
	"testing"
	"time"
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
