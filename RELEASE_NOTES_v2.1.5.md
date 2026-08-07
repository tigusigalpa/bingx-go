# Release Notes v2.1.5 - TradFi (Traditional Finance) Support

**Release Date:** April 25, 2026

We're excited to announce the addition of **Traditional Finance (TradFi)** support to the BingX Go SDK! This major feature expansion allows you to trade stock tokens, forex pairs, commodities, and stock indices through BingX's perpetual swap markets.

---

## 🚀 What's New

### TradFi (Traditional Finance) - Full Support

Access traditional financial instruments traded as perpetual swaps on BingX:

- **📈 Stock Tokens** - Trade popular stocks like TSLA, AAPL, NVDA, MSFT, AMZN, GOOGL
- **💱 Forex Pairs** - Trade major currency pairs: EUR-USD, GBP-USD, USD-JPY, USD-CHF
- **🥇 Commodities** - Trade GOLD, SILVER, OIL, NATURALGAS and more
- **📊 Stock Indices** - Trade SPX (S&P 500), DJI (Dow Jones), NDX (Nasdaq 100)

---

## 📦 New API Components

### 1. TradFi Client
```go
client := bingx.NewClient(apiKey, apiSecret)
tradfi := client.TradFi()
```

### 2. Market Service (15+ methods)
- `GetSymbols()` - All TradFi instruments
- `GetStockSymbols()` - Stock tokens only
- `GetForexSymbols()` - Forex pairs only
- `GetCommoditySymbols()` - Commodities only
- `GetIndexSymbols()` - Stock indices only
- `GetTicker(symbol)` - 24h ticker statistics
- `GetLatestPrice(symbol)` - Current price
- `GetDepth(symbol, limit)` - Order book depth
- `GetKlines(symbol, interval, limit, startTime, endTime)` - Candlestick data
- `GetMarkPrice(symbol)` - Mark price
- `GetFundingRate(symbol)` - Current funding rate
- `GetFundingRateHistory(symbol, limit)` - Historical funding rates
- `GetOpenInterest(symbol)` - Open interest data
- `GetRecentTrades(symbol, limit)` - Recent trades
- `GetBookTicker(symbol)` - Best bid/ask prices
- `GetTradingRules(symbol)` - Trading specifications

### 3. Trade Service (20+ methods)
All standard trading operations supported:
- `CreateOrder(params)` - Place orders (MARKET, LIMIT, STOP, etc.)
- `CreateTestOrder(params)` - Test order validation
- `CancelOrder(symbol, orderID, clientOrderID)` - Cancel by ID
- `CancelAllOrders(symbol)` - Cancel all open orders
- `GetOrder(symbol, orderID)` - Order details
- `GetOpenOrders(symbol, limit)` - List open orders
- `GetOrderHistory(symbol, limit, startTime, endTime)` - Historical orders
- `GetUserTrades(symbol, limit, startTime, endTime)` - Trade history
- `SetLeverage(symbol, leverage, side)` - Configure leverage
- `GetLeverage(symbol)` - Current leverage
- `SetMarginType(symbol, marginType)` - ISOLATED/CROSSED
- `GetMarginType(symbol)` - Current margin type
- `SetPositionMargin(...)` - Add/reduce margin
- `OneClickReversePosition(symbol)` - Reverse position side
- `ModifyOrder(...)` - Modify existing orders
- `SetAutoAddMargin(...)` - Auto margin addition
- **TWAP Orders** - Time-Weighted Average Price for large trades:
  - `PlaceTWAPOrder(params)`
  - `GetTWAPOrder(orderID)`
  - `GetTWAPOrders(symbol, status, limit)`
  - `CancelTWAPOrder(orderID)`

### 4. Account Service (15+ methods)
- `GetBalance()` - Account balance
- `GetAccountInfo()` - Comprehensive account info
- `GetPositions(symbol)` - Open positions
- `GetPositionRisk(symbol)` - Risk metrics (liquidation price, margin ratio)
- `GetIncomeHistory(...)` - PNL, funding fees, commissions
- `GetCommissionHistory(...)` - Commission records
- `GetForceOrders(...)` - Liquidation history
- `GetPositionMode()` / `SetPositionMode(hedgeMode)` - Hedge/One-way mode
- `GetMarginMode(symbol)` / `SetMarginMode(symbol, mode)` - Margin configuration
- `GetTradingFees(symbol)` - Fee rates
- `GetMultiAssetsMode()` / `SetMultiAssetsMode(enabled)` - Multi-assets margin
- `GetMultiAssetsMargin()` - Multi-assets details
- `GetAPIPermissions()` - API key permissions
- `GetBalanceHistory(coin, limit)` - Balance history

### 5. Listen Key Service (3 methods)
- `Create()` - Generate WebSocket listen key
- `Extend(listenKey)` - Extend key validity
- `Delete(listenKey)` - Delete listen key

---

## 💡 Usage Examples

### Stock Trading
```go
// Get TradFi client
tradfi := client.TradFi()

// Get available stocks
stocks, _ := tradfi.Market().GetStockSymbols()

// Check Tesla price
ticker, _ := tradfi.Market().GetTicker("TSLA-USDT")

// Buy Tesla stock with limit order
order, _ := tradfi.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "TSLA-USDT",
    "side":         "BUY",
    "type":         "LIMIT",
    "positionSide": "LONG",
    "price":        250.0,
    "quantity":     10.0,
    "timeInForce":  "GTC",
})
```

### Forex Trading
```go
// Get forex pairs
forex, _ := tradfi.Market().GetForexSymbols()

// Get EUR/USD price
price, _ := tradfi.Market().GetLatestPrice("EUR-USDT")

// Trade EUR/USD with 50x leverage
tradfi.Trade().SetLeverage("EUR-USDT", 50, nil)
order, _ := tradfi.Trade().CreateOrder(map[string]interface{}{
    "symbol": "EUR-USDT",
    "side":   "BUY",
    "type":   "MARKET",
    "quantity": 1000,
})
```

### Commodity Trading
```go
// Get commodity symbols
commodities, _ := tradfi.Market().GetCommoditySymbols()

// Check Gold price
price, _ := tradfi.Market().GetLatestPrice("GOLD-USDT")

// Get Gold klines (1-hour candles)
klines, _ := tradfi.Market().GetKlines("GOLD-USDT", "1h", 100, nil, nil)
```

### TWAP Orders (Large Trades)
```go
// Execute large Apple stock order over 1 hour
twap, _ := tradfi.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "AAPL-USDT",
    "side":     "BUY",
    "quantity": 500.0,  // Large quantity
    "duration": 3600,  // Spread over 1 hour
    "interval": 60,    // Order every 60 seconds
})
```

---

## ⚠️ Important Trading Notes

### Trading Hours
| Instrument Type | Trading Hours |
|----------------|---------------|
| **Stocks** | Mon-Fri 09:30-16:00 ET (US market hours) |
| **Forex** | Sun 17:00 ET - Fri 17:00 ET (24/5) |
| **Commodities** | Nearly 24h (varies by instrument) |
| **Indices** | Based on underlying market hours |

### Leverage Limits
- **Stocks**: Typically 5x-20x
- **Forex**: Typically 50x-100x
- **Commodities**: Typically 10x-50x
- **Indices**: Typically 10x-25x

### Funding Rates
- Applied every 8 hours: 00:00, 08:00, 16:00 UTC
- Positive rate = longs pay shorts
- Negative rate = shorts pay longs

---

## 📊 Statistics

| Metric | Value |
|--------|-------|
| New Methods | 50+ |
| Total API Methods | 260+ |
| New Services | 4 (Market, Trade, Account, ListenKey) |
| New Files | 5 (tradfi_client.go + 4 service files) |
| Test Coverage | 100% for new components |

---

## 🔧 Technical Details

### Architecture
```
services/tradfi/
├── market.go      # Market data endpoints
├── trade.go       # Trading operations
├── account.go     # Account management
└── listenkey.go   # WebSocket authentication

tradfi_client.go   # Main TradFi client
```

### Endpoint Prefix
- TradFi uses the same endpoints as USDT-M perpetuals: `/openApi/swap/v2/`
- All existing crypto perpetual features work with TradFi instruments

### Authentication
- Same API key/secret as crypto trading
- No additional permissions required
- WebSocket streaming supported

---

## 🔄 Migration Guide

No migration needed! This release is **100% backward compatible**.

Existing code continues to work unchanged. TradFi features are opt-in via the new `client.TradFi()` accessor.

```go
// Existing crypto trading - unchanged
btcPrice, _ := client.Market().GetLatestPrice("BTC-USDT")

// New TradFi trading
 tslaPrice, _ := client.TradFi().Market().GetLatestPrice("TSLA-USDT")
```

---

## 🧪 Testing

All new components include comprehensive tests:

```bash
cd /path/to/bingx-go
go test -v -run "TestTradFi"
```

Test coverage:
- ✅ Client initialization
- ✅ Singleton pattern (caching)
- ✅ All 4 services (Market, Trade, Account, ListenKey)
- ✅ Integration with main Client

---

## 📚 Documentation

- **README.md** - Updated with TradFi section
- **skill.md** - AI skill reference updated
- **Full Examples** - See README.md "TradFi (Traditional Finance)" section

---

## 🎯 What's Next

We're continuously expanding TradFi capabilities:
- Additional stock markets (EU, Asia)
- More forex pairs
- Expanded commodity offerings
- Enhanced WebSocket streams for TradFi

---

## 🐛 Bug Fixes

None - this is a feature release.

---

## 🙏 Contributors

Thanks to the community for feedback and feature requests that led to this release!

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.1.4...v2.1.5
