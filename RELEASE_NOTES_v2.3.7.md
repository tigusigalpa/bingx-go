# Release Notes v2.3.7 - Client Reliability Fixes

**Release Date:** August 11, 2026

This patch improves the reliability of shared HTTP and WebSocket client infrastructure. It introduces no breaking API changes.

---

## Fixed

### HTTP Client

- `BaseHTTPClient.Request()` and `RequestJSON()` no longer add a generated `timestamp` to the caller-provided parameter map. Reusing the same map now produces a fresh timestamp for each request unless the caller explicitly supplies one.
- Non-2xx HTTP responses that do not contain a BingX API error code now return `*errors.BingXException` with the HTTP status code and decoded response payload. BingX API error mapping remains unchanged when an API error code is present.

### WebSocket Client

- `WebSocketClient.Disconnect()` is now idempotent and no longer panics when it is called more than once.
- A `WebSocketClient` can now connect again after `Disconnect()`; each new connection receives a fresh listener shutdown signal.
- `Listen()` captures its active connection before reading, preventing it from observing a concurrently cleared connection during disconnect.
- Concurrent first access to `Client.CoinM()` or `Client.TradFi()` is synchronized, preventing duplicate lazy initialization and data races.
- Account-stream `listenKey` values are URL-escaped before the WebSocket URL is constructed.

## Compatibility

- No public methods, method signatures, endpoints, request signing rules, or default settings were changed.
- Existing error types for BingX API responses are preserved. The new HTTP-status error is returned only when a non-2xx response has no mapped API error code.
- Calling `Connect()` on an already connected WebSocket client remains safe and leaves the active connection unchanged.

## Tests

- Added regression tests for immutable HTTP request parameters and HTTP failures without an API code.
- Added WebSocket regression tests for idempotent disconnect and reconnect-after-disconnect behavior.
- Verified with `go test ./...`, `go vet ./...`, and `go test -race ./...`.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.3.6...v2.3.7
