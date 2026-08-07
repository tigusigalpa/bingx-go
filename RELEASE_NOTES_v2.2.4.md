# Release Notes v2.2.4 - Fix Default Signature Encoding to Hex

**Release Date:** August 7, 2026

This patch release fixes the default HMAC signature encoding used by `NewClient()` and `NewDemoClient()`. The SDK previously defaulted to `base64`, which caused authenticated BingX requests to fail with `Signature verification failed due to signature mismatch`. BingX expects the HMAC-SHA256 digest as a lowercase hexadecimal string.

---

## 🐛 Bug Fixes

### Signature Encoding Default

- `ClientConfig.SignatureEncoding` now defaults to `"hex"` in `NewClient()` and `NewDemoClient()`.
- The shared `BaseHTTPClient` encodes HMAC-SHA256 signatures as 64-character lowercase hex by default.
- `WithSignatureEncoding("base64")` is still supported for backward compatibility, but any unrecognized or empty encoding safely falls back to hex rather than silently producing base64.

### Affected Authenticated Calls

This fix applies to every signed request, including:

```go
client.SpotAccount().GetBalance()
client.SpotAccount().GetFundBalance()
client.Account().GetBalance()
client.Trade().CreateOrder(...)
client.SpotTrade().CreateOrderRequest(...)
```

---

## 🔧 Technical Details

- `client.go`: default `SignatureEncoding` changed from `"base64"` to `"hex"`.
- `http/client.go`: `signString()` now treats hex as the default and only produces base64 when explicitly requested.
- `WithSignatureEncoding()` is documented as hex-first; base64 is retained for backward compatibility.
- The signing and transport logic is unchanged: GET, DELETE, and form POST requests continue to build the canonical query/form payload, sign it, and append `signature=<hex digest>` before sending.

---

## 🧪 Testing

New deterministic regression tests were added to `http/client_test.go` and `client_test.go`:

- Default `NewClient()` and `NewDemoClient()` produce 64-character lowercase hex signatures.
- Fixed test vector: `symbol=BTC-USDT&timestamp=1702731500000` signed with `test-secret` yields `36f610a8da9e21e9c417413745bcca4dddad81cea81e6d66baa3ef0c2229021c`.
- `httptest.Server` verifies GET, POST, and DELETE signatures against the received query/body.
- Explicit `WithSignatureEncoding("hex")` is honored.
- `WithSignatureEncoding("base64")` still works.
- Unknown/empty encodings fall back to hex, never to base64.

No real API keys or network calls are used in the test suite.

---

## 🔄 Migration Guide

No code changes are required for most users.

- If you were **not** explicitly passing `bingx.WithSignatureEncoding("base64")`, delete that option (if present) and the new default `hex` will be used.
- If you were explicitly using `bingx.WithSignatureEncoding("hex")`, behavior is identical.
- If you rely on base64 signatures for a non-BingX environment, keep `bingx.WithSignatureEncoding("base64")`.

---

## 📚 Documentation

- `README.md` quick-start and configuration examples updated to use hex.
- `skill.md` updated to describe hex as the default.
- `wiki/Getting-Started.md` troubleshooting section updated.
- `examples/basic/main.go` and `examples/trading/main.go` updated to remove explicit `WithSignatureEncoding("base64")`.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.2.3...v2.2.4
