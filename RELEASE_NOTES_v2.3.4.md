# Release Notes v2.3.4 - Fix BingX Signed REST Request Algorithm

**Release Date:** August 7, 2026

This release fixes the signed REST request construction used by `BaseHTTPClient.Request()`. The previous implementation used `url.Values.Encode()` to build the query, then re-sorted all parameters together with the signature. This produced URLs where `signature` appeared before `timestamp` and other parameters, causing BingX to respond with `100001 Signature verification failed` even when the HMAC digest was encoded as hex.

---

## 🐛 Bug Fixes

### Canonical Signing String

- `BaseHTTPClient.Request()` now builds a separate **canonical signing string** before creating the wire query/body:
  - All business parameters + `timestamp` are included.
  - Keys are sorted in ASCII ascending order.
  - The string is built as raw `key=value&key=value`.
  - Values are **not** URL-encoded.
  - The `signature` parameter is **never** part of the signing string.
- `HMAC-SHA256(apiSecret, canonicalString)` is computed and returned as a 64-character lowercase hex digest by default.
- `WithSignatureEncoding("base64")` still works as an explicit legacy option.

### GET and DELETE Requests

- Final URL format is now strictly `?canonicalString&signature=<hex>`.
- `signature` is always appended **last**, not sorted with business parameters.
- Values containing `[` or `{` (e.g., batch order JSON arrays) are URL-escaped in the actual query string, while the canonical signing string stays raw.

### POST / PUT Requests

- Request body is `canonicalString&signature=<hex>`.
- `Content-Type` remains `application/x-www-form-urlencoded`.

### Affected Authenticated Calls

This fix applies to every signed request, including:

```go
client.Account().GetBalance()
client.SpotAccount().GetBalance()
client.Trade().CreateOrder(...)
client.Trade().CreateBatchOrders(...)
client.SpotTrade().CreateOrderRequest(...)
client.SpotTrade().CancelBatchOrders(...)
```

---

## 🔧 Technical Details

- `http/client.go`:
  - `buildCanonicalString()` builds the raw, sorted, unencoded signing payload.
  - `buildSignedString()` builds the final GET/DELETE URL query or POST body, appending the signature last and URL-escaping values that contain `[` or `{` only for query strings.
  - `Request()` no longer re-adds `signature` to the parameter map or re-sorts via `url.Values.Encode()`.
- `paramValueToString()` now serializes complex values (slices, maps) as JSON while keeping scalar values raw.

---

## 🧪 Testing

New regression tests were added in `http/client_signature_test.go` and `http/client_test.go`:

- GET `/openApi/swap/v3/user/balance` produces `timestamp=<fixed>&signature=<expected_hex>`.
- GET with `symbol=BTC-USDT` produces `symbol=BTC-USDT&timestamp=<fixed>&signature=<expected_hex>`.
- POST body is exactly `canonicalString&signature=<expected_hex>`.
- Batch `orders=[...]` value is URL-encoded in the query but raw in the canonical signature.
- `buildCanonicalString` and `buildSignedString` unit tests verify sorting, raw canonical, conditional URL-encoding, and signature-last ordering.

All tests run with fixed secrets and `httptest.Server` — no real API keys or network calls are required.

---

## 🔄 Migration Guide

No code changes are required for users of `NewClient()` or `NewDemoClient()`.

- Remove any manual `signature` parameters from calls; the SDK adds and positions it automatically.
- If you were passing pre-URL-encoded batch values, pass them as raw JSON strings or `[]map[string]interface{}`; the SDK encodes them in the URL but signs the raw canonical string.

---

## 📚 Documentation

- `CHANGELOG.md` updated with the v2.3.4 release notes.
- `README.md`, `skill.md`, and `wiki/Getting-Started.md` already document hex signatures and the corrected query ordering.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.3.3...v2.3.4
