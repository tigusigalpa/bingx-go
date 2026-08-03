# Release Notes v2.2.3 - Spot Trading Support

**Release Date:** August 3, 2026

This release adds full **Spot Trading** support to the BingX Go SDK. Until now, `SpotAccount()` only covered spot balances and transfers — there was no way to create, cancel, amend, or query spot orders. `SpotTrade()` closes that gap with a dedicated service for the BingX Spot Trading REST API.

---

## 🚀 What's New

### Spot Trading (`SpotTradeService`) - Full Support

```go
client := bingx.NewClient(apiKey, apiSecret)
spot := client.SpotTrade()
```

`SpotTrade()` is for **spot** orders (buying/selling the underlying asset directly, no leverage or position side) — as opposed to `Trade()`, which is for **perpetual futures** orders. It reuses the SDK's existing signed HTTP client, so there is no duplicated HMAC signing or error-handling logic.

---

## 📦 New API Components

### 1. Order Placement
- `CreateOrder(params map[string]interface{})` - Place any spot order using raw exchange parameters (`POST /openApi/spot/v1/trade/order`). Use this for advanced order types not covered by `CreateOrderRequest` (`TAKE_STOP_LIMIT`, `TAKE_STOP_MARKET`, `TRIGGER_LIMIT`, `TRIGGER_MARKET`).
- `CreateOrderRequest(req SpotOrderRequest)` - Typed, client-side-validated convenience method for `LIMIT`/`MARKET` orders.

### 2. Cancellation
- `CancelOrder(symbol, orderID, clientOrderID *string)` - Cancel a single order (`POST /openApi/spot/v1/trade/cancel`).
- `CancelAllOrders(symbol *string)` - Cancel all open orders for a symbol, or across all symbols if `nil` (`POST /openApi/spot/v1/trade/cancelOpenOrders`).
- `CancelBatchOrders(symbol string, orderIDs, clientOrderIDs []string)` - Cancel a batch of orders (`POST /openApi/spot/v1/trade/cancelOrders`).

### 3. Amendment
- `AmendOrder(symbol string, cancelOrderID, cancelClientOrderID *string, cancelReplaceMode string, newOrder SpotOrderRequest)` - Atomically cancel an order and place a new one in a single signed request (`POST /openApi/spot/v1/trade/order/cancelReplace`).

  BingX Spot does not support true in-place order modification the way Futures does (`Trade().ModifyOrder`, which only changes an existing order's quantity). `AmendOrder` wraps the exchange's native cancel-and-replace endpoint instead, which is safer than a manual `CancelOrder` + `CreateOrder` sequence because there is never a window where neither or both orders are resting on the book.

### 4. Queries
- `GetOrder(symbol string, orderID, clientOrderID *string)` - Query a single order (`GET /openApi/spot/v1/trade/query`).
- `GetOpenOrders(symbol *string)` - Current open orders (`GET /openApi/spot/v1/trade/openOrders`).
- `GetOrderHistory(symbol *string, limit int, startTime, endTime *int64)` - Historical orders; `limit` is sent as the API's `pageSize` parameter (`GET /openApi/spot/v1/trade/historyOrders`).
- `GetTrades(symbol *string, limit int, startTime, endTime *int64)` - Filled trade details; `symbol` is required by this endpoint (`GET /openApi/spot/v1/trade/myTrades`).

### 5. Types & Constants
```go
type SpotOrderRequest struct {
    Symbol        string
    Side          string
    Type          string
    Quantity      string
    Price         *string
    QuoteOrderQty *string
    TimeInForce   *string
    ClientOrderID *string
}

const (
    SpotOrderTypeLimit  = "LIMIT"
    SpotOrderTypeMarket = "MARKET"

    SpotSideBuy  = "BUY"
    SpotSideSell = "SELL"

    SpotTimeInForceGTC      = "GTC"
    SpotTimeInForceIOC      = "IOC"
    SpotTimeInForceFOK      = "FOK"
    SpotTimeInForcePostOnly = "PostOnly"

    SpotCancelReplaceStopOnFailure = "STOP_ON_FAILURE"
    SpotCancelReplaceAllowFailure  = "ALLOW_FAILURE"
)
```

**Decimal safety:** `Quantity`, `Price`, and `QuoteOrderQty` are `string`, never `float64`. Representing decimals as `float64` risks silent binary floating-point rounding (e.g. `0.1 + 0.2 != 0.3`), which can lead to rejected orders or subtly wrong fills. Always format decimal amounts as exact strings before passing them in.

---

## 💡 Usage Examples

### Limit Buy Order
```go
price := "60000"
quantity := "0.001"

order, err := client.SpotTrade().CreateOrderRequest(services.SpotOrderRequest{
    Symbol:   "BTC-USDT",
    Side:     services.SpotSideBuy,
    Type:     services.SpotOrderTypeLimit,
    Quantity: quantity,
    Price:    &price,
})
```

### Market Buy Order (spend 100 USDT)
```go
quoteOrderQty := "100"

order, err := client.SpotTrade().CreateOrderRequest(services.SpotOrderRequest{
    Symbol:        "BTC-USDT",
    Side:          services.SpotSideBuy,
    Type:          services.SpotOrderTypeMarket,
    QuoteOrderQty: &quoteOrderQty,
})
```

### Cancel an Order
```go
orderID := "123456789"
_, err := client.SpotTrade().CancelOrder("BTC-USDT", &orderID, nil)
```

### Amend an Order (cancel-and-replace)
```go
newPrice := "61000"
_, err := client.SpotTrade().AmendOrder(
    "BTC-USDT",
    &orderID, nil,
    services.SpotCancelReplaceStopOnFailure,
    services.SpotOrderRequest{
        Side:     services.SpotSideBuy,
        Type:     services.SpotOrderTypeLimit,
        Quantity: "0.001",
        Price:    &newPrice,
    },
)
```

---

## ⚠️ Not Implemented — Spot Test Orders

Unlike `Trade().CreateTestOrder` (Futures `POST /openApi/swap/v2/trade/order/test`), the BingX Spot API does not expose a dry-run order-validation endpoint. `SpotTradeService` therefore has **no** `CreateTestOrder` method — adding one would either silently place a real order or fabricate a non-existent endpoint. Validate orders client-side via `CreateOrderRequest`'s built-in checks instead:

- `symbol`, `side`, and `type` are mandatory.
- `side` must be `BUY` or `SELL`.
- `type` must be `LIMIT` or `MARKET` (use `CreateOrder` for other types).
- a `LIMIT` order requires a positive `price` and a positive `quantity`.
- a `MARKET` order requires a positive `quantity` or `quoteOrderQty`.

---

## 🔧 Technical Details

### Architecture
```
services/
├── spottrade.go       # New: SpotTradeService, SpotOrderRequest, constants
└── spottrade_test.go  # New: httptest-based unit tests

client.go              # Added spotTrade field + Client.SpotTrade() accessor
```

### Demo/VST Support
`NewDemoClient()` creates `SpotTradeService` against the same demo base URL (`https://open-api-vst.bingx.com`) as every other service, since all services share the client's single signed HTTP transport. No additional wiring was needed.

---

## 🧪 Testing

New tests use `httptest.NewServer` to verify request path, method, body, and signature flow without any real API key or network access:

```bash
cd /path/to/bingx-go
go test ./services/... -run "CreateOrderRequest|CancelOrder|CancelAllOrders|CancelBatchOrders|GetOrder|GetTrades|AmendOrder"
go test . -run TestNewDemoClient
```

Coverage:
- ✅ Limit order placement (path/method/body/signature)
- ✅ Market order via `quoteOrderQty`
- ✅ Invalid typed requests (8 validation cases)
- ✅ Cancellation by `orderID` and by `clientOrderID`
- ✅ Cancel-all and batch-cancel
- ✅ `GetOrder`, `GetOpenOrders`, `GetOrderHistory`
- ✅ `AmendOrder` (cancel-replace) and its validation
- ✅ Demo client base URL propagation
- ✅ No regressions in existing Futures/Coin-M/TradFi/Copy Trading tests

---

## 🔄 Migration Guide

No migration needed! This release is **100% backward compatible**. All existing methods, signatures, and endpoints (Futures, Coin-M, TradFi, Copy Trading, `SpotAccount`) are unchanged. Spot Trading is purely additive via the new `client.SpotTrade()` accessor.

---

## 📚 Documentation

- **README.md** - New "Spot Trading - Spot Order Management" section with full examples and a spot-vs-futures explanation.
- **skill.md** - Not yet updated for this release; see follow-up work.

---

## 🐛 Bug Fixes

None - this is a feature release.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.1.5...v2.2.3
