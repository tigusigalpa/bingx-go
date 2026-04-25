# bingx-go SDK — AI Skill Reference

> Full-featured Go SDK for the BingX cryptocurrency exchange API.
> Module: `github.com/tigusigalpa/bingx-go/v2`
> Requires: Go 1.21+, dependency: `gorilla/websocket`

---

## Installation

```bash
go get github.com/tigusigalpa/bingx-go/v2
go mod tidy
```

---

## Client Initialization

```go
import bingx "github.com/tigusigalpa/bingx-go/v2"

// Standard client
client := bingx.NewClient("API_KEY", "API_SECRET")

// With options
client := bingx.NewClient(apiKey, apiSecret,
    bingx.WithBaseURI("https://open-api.bingx.com"),   // default
    bingx.WithSignatureEncoding("hex"),                 // or "base64" (default)
    bingx.WithSourceKey("MyBot"),                       // optional tag
)

// Demo/VST environment (paper trading)
client := bingx.NewDemoClient(apiKey, apiSecret)
```

**Client options:**
- `WithBaseURI(uri string)` — override API endpoint
- `WithSignatureEncoding(enc string)` — `"base64"` (default) or `"hex"`
- `WithSourceKey(key string)` — label requests for debugging
- `WithDemoEnvironment()` — sets base URI to `https://open-api-vst.bingx.com`

**Service accessors on `*Client`:**
```go
client.Market()      // *services.MarketService
client.Account()     // *services.AccountService
client.Trade()       // *services.TradeService
client.Wallet()      // *services.WalletService
client.SpotAccount() // *services.SpotAccountService
client.SubAccount()  // *services.SubAccountService
client.CopyTrading() // *services.CopyTradingService
client.Contract()    // *services.ContractService
client.ListenKey()   // *services.ListenKeyService
client.CoinM()       // *CoinMClient  (Coin-M futures)
client.TradFi()      // *TradFiClient (Traditional Finance - stocks, forex, commodities)
```

Shortcut methods on client root: `GetBalance()`, `GetSymbols()`, `CreateOrder()`, `NewMarketDataStream()`, `NewAccountDataStream(listenKey)`.

---

## Market Service — `client.Market()`

All market endpoints. Most do **not** require authentication.

### Symbols / Trading Pairs
```go
client.Market().GetFuturesSymbols()         // USDT-M perpetual pairs
client.Market().GetSpotSymbols()            // Spot pairs
client.Market().GetAllSymbols()             // All pairs combined
```

### Prices
```go
client.Market().GetLatestPrice("BTC-USDT")         // current price
client.Market().GetIndexPrice("BTC-USDT")           // index price (v3)
sym := "BTC-USDT"
client.Market().GetTickerPrice(&sym)                // ticker price (v3)
```

### Order Book
```go
client.Market().GetDepth("BTC-USDT", 20)            // depth levels (5/10/20/50/100/500/1000)
client.Market().GetBookTicker(&sym)                  // best bid/ask (v3)
client.Market().GetSpotBookTicker(&sym)              // spot best bid/ask (v3)
```

### Klines / Candlesticks
```go
client.Market().GetKlines("BTC-USDT", "1h", 100, startTime, endTime)
// timeframes: 1m 3m 5m 15m 30m 1h 2h 4h 6h 12h 1d 3d 1w 1M
// startTime/endTime are *int64 (UnixMilli), pass nil for latest

client.Market().GetSpotKlines("BTC-USDT", "1h", 100, nil, nil, &timeZone)
// timeZone is *int64 (UTC offset hours)
```

### 24h Ticker Statistics
```go
client.Market().Get24hrTicker(&sym)   // single symbol
client.Market().Get24hrTicker(nil)    // all symbols
// fields: priceChange, priceChangePercent, highPrice, lowPrice, volume
```

### Funding Rates
```go
client.Market().GetFundingRateHistory("BTC-USDT", 100)
client.Market().GetFundingRateInfo("BTC-USDT")    // current rate + next payment time (v3)
// positive rate = longs pay shorts; negative = shorts pay longs
```

### Open Interest (v3)
```go
client.Market().GetOpenInterest("BTC-USDT")
client.Market().GetOpenInterestHistory("BTC-USDT", "1h", 100, nil, nil)
// periods: 5m 15m 30m 1h 4h 1d
```

### Trades & Sentiment (v3)
```go
client.Market().GetRecentTrades("BTC-USDT", 100)
client.Market().GetAggregateTrades("BTC-USDT", 100, fromId, startTime, endTime)
client.Market().GetLongShortRatio("BTC-USDT", "5m", 100)
client.Market().GetBasis("BTC-USDT", "CURRENT_QUARTER", 100, nil, nil)
```

---

## Trade Service — `client.Trade()`

Requires authentication.

### Order Type Constants (`services` package)
```go
import "github.com/tigusigalpa/bingx-go/v2/services"

services.OrderTypeMarket             // "MARKET"
services.OrderTypeLimit              // "LIMIT"
services.OrderTypeStop               // "STOP"
services.OrderTypeStopMarket         // "STOP_MARKET"
services.OrderTypeTakeProfit         // "TAKE_PROFIT"
services.OrderTypeTakeProfitMarket   // "TAKE_PROFIT_MARKET"
services.OrderTypeTriggerLimit       // "TRIGGER_LIMIT"          (v3)
services.OrderTypeTrailingStopMarket // "TRAILING_STOP_MARKET"   (v3)
services.OrderTypeTrailingTPSL       // "TRAILING_TP_SL"         (v3)
```

### Place / Test Orders
```go
// Real order
client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",          // BUY | SELL
    "type":         services.OrderTypeLimit,
    "positionSide": "LONG",         // LONG | SHORT | BOTH
    "price":        50000.0,        // required for LIMIT
    "quantity":     0.001,
    "timeInForce":  "GTC",          // GTC | IOC | FOK
})

// Test order — validates without executing
client.Trade().CreateTestOrder(params)
```

### Trailing Stop (v3)
```go
client.Trade().CreateOrder(map[string]interface{}{
    "symbol":          "ETH-USDT",
    "side":            "SELL",
    "type":            services.OrderTypeTrailingStopMarket,
    "positionSide":    "LONG",
    "activationPrice": 3000.0,   // start trailing at this price
    "callbackRate":    1.0,      // trail % behind the high
    "quantity":        1.0,
})
```

### Trailing TP/SL (v3)
```go
client.Trade().CreateOrder(map[string]interface{}{
    "symbol":              "BTC-USDT",
    "side":               "BUY",
    "type":               services.OrderTypeTrailingTPSL,
    "positionSide":       "LONG",
    "quantity":           0.1,
    "takeProfitPrice":    52000.0,
    "stopLossPrice":      48000.0,
    "trailingStopPercent": 0.5,
})
```

### TWAP Orders (v3)
```go
// Place TWAP — execute large order over time
client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "positionSide": "LONG",
    "quantity":     10.0,
    "duration":     3600,  // total seconds
    "interval":     60,    // seconds between sub-orders
    // "price": 50000.0,  // optional limit price
})

client.Trade().GetTWAPOrder("twap_order_id", nil)
client.Trade().GetTWAPOrders(&symbol, &status, nil, nil, 100, nil)
// status: "WORKING" | "FINISHED" | "CANCELLED"

client.Trade().CancelTWAPOrder("twap_order_id", nil)
```

### Order Management
```go
client.Trade().GetOrder("BTC-USDT", "order_id")
client.Trade().GetOpenOrders(nil, 100)            // all symbols
client.Trade().GetOpenOrders(&symbol, 100)        // one symbol
client.Trade().GetOrderHistory(&symbol, 100, nil, nil)
client.Trade().GetUserTrades(&symbol, 100, nil, nil)

// Cancel by ID or clientOrderID
client.Trade().CancelOrder("BTC-USDT", &orderID, &clientID, nil, nil)
client.Trade().CancelAllOrders(nil, &symbol, nil, nil)   // by symbol
client.Trade().CancelAllOrders(nil, nil, nil, nil)       // all orders

// Modify price/quantity without cancel-replace
client.Trade().ModifyOrder("BTC-USDT", &orderID, nil, 51000.0, 0.002, nil)
```

### Position Management
```go
// One-click reversal: LONG↔SHORT atomically (v3)
client.Trade().OneClickReversePosition("BTC-USDT", nil)
```

### Leverage & Margin
```go
client.Trade().SetLeverage("BTC-USDT", 20, nil, nil)
client.Trade().SetLeverage("BTC-USDT", 10, &side, nil)  // side: *"LONG"|*"SHORT"

// Auto add margin when near liquidation (v3, hedge mode only)
client.Trade().SetAutoAddMargin("BTC-USDT", "LONG", true, nil)

// Multi-assets margin mode (v3)
client.Trade().SwitchMultiAssetsMode(true, nil)
client.Trade().GetMultiAssetsMode(nil)
client.Trade().GetMultiAssetsRules(nil)
client.Trade().GetMultiAssetsMargin(nil)
```

### Commission Calculation (local, no API call)
```go
result := client.Trade().CalculateFuturesCommission(margin, leverage, nil)
// result.Commission, result.PositionValue, result.NetPositionValue, etc.
// Default rate: FuturesCommissionRate = 0.00045 (0.045%)

amount := client.Trade().GetCommissionAmount(margin, leverage)
```

---

## Account Service — `client.Account()`

### Balance & Account Info
```go
client.Account().GetBalance()
// fields: availableBalance, balance, unrealizedProfit

client.Account().GetAccountInfo()
// fields: totalEquity, totalMargin, availableMargin, marginLevel
```

### Positions
```go
client.Account().GetPositions(nil)       // all positions
client.Account().GetPositions(&symbol)   // specific symbol
// fields: positionAmt, positionSide, avgPrice, unrealizedProfit, symbol
```

### Leverage & Margin Mode
```go
client.Account().GetLeverage("BTC-USDT", nil)
client.Account().SetLeverage("BTC-USDT", "BOTH", 20, nil)
client.Account().SetLeverage("BTC-USDT", "LONG", 10, nil)

client.Account().GetMarginMode("BTC-USDT")
client.Account().SetMarginMode("BTC-USDT", "ISOLATED")   // or "CROSSED"
```

### Position Mode (v3)
```go
client.Account().GetPositionMode(nil)
client.Account().SetPositionMode(true, nil)   // true = hedge, false = one-way
// Cannot switch with open positions
```

### Risk Monitoring (v3)
```go
client.Account().GetPositionRisk(&symbol, nil)
// fields: liquidationPrice, marginRatio, leverage, unrealizedProfit
```

### P&L & Income History (v3)
```go
client.Account().GetIncomeHistory(&symbol, &incomeType, &startTime, &endTime, 100, nil)
// incomeType: "REALIZED_PNL" | "FUNDING_FEE" | "COMMISSION" | "TRANSFER" | "INSURANCE_CLEAR"

client.Account().GetCommissionHistory("BTC-USDT", nil, nil, 100, nil)

client.Account().GetForceOrders(&symbol, nil, nil, nil, 100, nil)  // liquidation history
```

### Other Account Methods
```go
client.Account().GetBalanceHistory(100, nil, nil)
client.Account().GetAPIPermissions()
// fields: enableTrading, enableWithdrawals, enableReading
```

---

## Wallet Service — `client.Wallet()`

```go
client.Wallet().GetDepositHistory(coin, &status, &startTime, &endTime, limit)
client.Wallet().GetDepositAddress(coin, network)
client.Wallet().GetWithdrawalHistory(coin, &status, &startTime, &endTime, limit)
client.Wallet().Withdraw(coin, address, amount, network, &addressTag)
client.Wallet().GetAllCoinInfo()
client.Wallet().GetMainAccountTransferHistory(coin, &transferType, &startTime, &endTime, limit)
```

---

## Spot Account Service — `client.SpotAccount()`

Manages spot account balance and spot-related operations.

---

## Sub-Account Service — `client.SubAccount()`

### Wallet Type Constants
```go
services.SubAccountWalletTypeFund             = 1   // Fund Account
services.SubAccountWalletTypeStandardFutures  = 2   // Standard Futures
services.SubAccountWalletTypePerpetualFutures = 3   // Perpetual Futures
services.SubAccountWalletTypeSpot             = 15  // Spot Account
```

### Methods
```go
client.SubAccount().CreateSubAccount(subAccountString)
client.SubAccount().GetAccountUID()
client.SubAccount().GetSubAccountList(&subAccountString, current, size)
client.SubAccount().GetSubAccountAssets(subUID)
client.SubAccount().UpdateSubAccountStatus(subAccountString, status)
client.SubAccount().GetAllSubAccountBalances()
```

---

## Copy Trading Service — `client.CopyTrading()`

### Perpetual (Swap) Copy Trading
```go
client.CopyTrading().GetCurrentTrackOrders(symbol)
client.CopyTrading().CloseTrackOrder(orderNumber)
client.CopyTrading().SetTPSL(positionID, &stopLoss, &takeProfit)
client.CopyTrading().GetTraderDetail()
client.CopyTrading().GetProfitSummary()
client.CopyTrading().GetProfitDetail(pageIndex, pageSize)
client.CopyTrading().SetCommission(commission)
client.CopyTrading().GetTradingPairs()
```

### Spot Copy Trading
```go
client.CopyTrading().SellSpotOrder(buyOrderID)
client.CopyTrading().GetSpotTraderDetail()
client.CopyTrading().GetSpotProfitSummary()
client.CopyTrading().GetSpotProfitDetail(pageIndex, pageSize)
client.CopyTrading().GetSpotHistoryOrders(pageIndex, pageSize)
```

---

## Coin-M Futures — `client.CoinM()`

Coin-margined (inverse) perpetual contracts. Endpoint prefix: `/openApi/cswap/v1/`

```go
coinM := client.CoinM()

// Market data
coinM.Market().GetSymbols()
coinM.Market().GetDepth(symbol, limit)
coinM.Market().GetKlines(symbol, interval, limit, startTime, endTime)

// Trading
coinM.Trade().CreateOrder(params)
coinM.Trade().CancelOrder(symbol, &orderID, &clientOrderID)
coinM.Trade().CancelAllOrders(symbol)
coinM.Trade().GetOrder(symbol, orderID)
coinM.Trade().GetOpenOrders(&symbol)
coinM.Trade().GetPositions(&symbol)
coinM.Trade().GetBalance()
coinM.Trade().GetLeverage(symbol)
coinM.Trade().SetLeverage(symbol, side, leverage)
coinM.Trade().GetMarginType(symbol)
coinM.Trade().SetMarginType(symbol, marginType)
coinM.Trade().SetPositionMargin(symbol, positionSide, amount, marginType)
coinM.Trade().GetOrderHistory(symbol, limit, startTime, endTime)
coinM.Trade().GetUserTrades(symbol, limit, startTime, endTime)
coinM.Trade().GetPositionRisk(&symbol, &recvWindow)
coinM.Trade().GetIncomeHistory(&symbol, &incomeType, startTime, endTime, limit, &recvWindow)

// Listen Key
coinM.ListenKey().Create()
coinM.ListenKey().Extend(listenKey)
coinM.ListenKey().Delete(listenKey)
```

---

## TradFi (Traditional Finance) — `client.TradFi()`

TradFi provides access to traditional financial instruments traded as perpetual swaps on BingX:
- **Stock Tokens**: TSLA-USDT, AAPL-USDT, NVDA-USDT, etc.
- **Forex Pairs**: EUR-USD, GBP-USD, USD-JPY, etc.
- **Commodities**: GOLD-USDT, SILVER-USDT, OIL-USDT, etc.
- **Stock Indices**: SPX-USDT, DJI-USDT, NDX-USDT, etc.

Endpoint prefix: `/openApi/swap/v2/` (same as USDT-M perpetuals)

```go
tradfi := client.TradFi()

// Market Data
tradfi.Market().GetSymbols()                    // All TradFi symbols
tradfi.Market().GetStockSymbols()               // Stock tokens only
tradfi.Market().GetForexSymbols()               // Forex pairs only
tradfi.Market().GetCommoditySymbols()           // Commodities only
tradfi.Market().GetIndexSymbols()                // Stock indices only
tradfi.Market().GetTicker("TSLA-USDT")          // 24h ticker
tradfi.Market().GetLatestPrice("EUR-USDT")      // Latest price
tradfi.Market().GetDepth("GOLD-USDT", 20)       // Order book depth
tradfi.Market().GetKlines("AAPL-USDT", "1h", 100, nil, nil)  // Candlesticks
tradfi.Market().GetMarkPrice("TSLA-USDT")       // Mark price
tradfi.Market().GetFundingRate("EUR-USDT")      // Current funding rate
tradfi.Market().GetFundingRateHistory("SPX-USDT", 100)
tradfi.Market().GetOpenInterest("GOLD-USDT")    // Open interest
tradfi.Market().GetRecentTrades("TSLA-USDT", 50)
tradfi.Market().GetBookTicker("AAPL-USDT")      // Best bid/ask
tradfi.Market().GetTradingRules("EUR-USDT")     // Trading specifications

// Trading - Same interface as crypto perpetuals
tradfi.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "TSLA-USDT",
    "side":         "BUY",
    "type":         "LIMIT",
    "positionSide": "LONG",
    "price":        250.0,
    "quantity":     1.0,
})
tradfi.Trade().CreateTestOrder(params)
tradfi.Trade().CancelOrder("TSLA-USDT", &orderID, nil)
tradfi.Trade().CancelAllOrders(&symbol)
tradfi.Trade().GetOrder("EUR-USDT", orderID)
tradfi.Trade().GetOpenOrders(&symbol, 100)
tradfi.Trade().GetOrderHistory(&symbol, 100, nil, nil)
tradfi.Trade().GetUserTrades(&symbol, 100, nil, nil)
tradfi.Trade().SetLeverage("GOLD-USDT", 20, nil)
tradfi.Trade().GetLeverage("AAPL-USDT")
tradfi.Trade().SetMarginType("EUR-USDT", "ISOLATED")
tradfi.Trade().GetMarginType("TSLA-USDT")
tradfi.Trade().SetPositionMargin(symbol, positionSide, amount, marginType)
tradfi.Trade().OneClickReversePosition("TSLA-USDT")
tradfi.Trade().ModifyOrder(symbol, &orderID, nil, newPrice, newQty)
tradfi.Trade().SetAutoAddMargin(symbol, positionSide, true)

// TWAP Orders (Time-Weighted Average Price)
tradfi.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "AAPL-USDT",
    "side":     "BUY",
    "quantity": 100.0,
    "duration": 3600,
    "interval": 60,
})
tradfi.Trade().GetTWAPOrder(orderID)
tradfi.Trade().GetTWAPOrders(&symbol, &status, 100)
tradfi.Trade().CancelTWAPOrder(orderID)

// Account
tradfi.Account().GetBalance()
tradfi.Account().GetAccountInfo()
tradfi.Account().GetPositions(&symbol)
tradfi.Account().GetPositionRisk(&symbol)
tradfi.Account().GetIncomeHistory(&symbol, &incomeType, startTime, endTime, 100)
tradfi.Account().GetCommissionHistory(symbol, startTime, endTime, 100)
tradfi.Account().GetForceOrders(&symbol, startTime, endTime, 100)
tradfi.Account().GetPositionMode()
tradfi.Account().SetPositionMode(true)  // hedge mode
tradfi.Account().GetMarginMode(symbol)
tradfi.Account().SetMarginMode(symbol, "ISOLATED")
tradfi.Account().GetTradingFees(symbol)
tradfi.Account().GetMultiAssetsMode()
tradfi.Account().SetMultiAssetsMode(true)

// Listen Key (WebSocket)
tradfi.ListenKey().Create()
tradfi.ListenKey().Extend(listenKey)
tradfi.ListenKey().Delete(listenKey)
```

### TradFi Trading Notes

- **Trading Hours**: Stock tokens follow US market hours (Mon-Fri 09:30-16:00 ET)
- **Forex Hours**: 24/5 (Sun 17:00 ET - Fri 17:00 ET)
- **Commodities**: Vary by instrument (GOLD, SILVER trade nearly 24h)
- **Leverage**: Typically lower than crypto (5x-20x for stocks, 50x-100x for forex)
- **Funding Rates**: Applied every 8 hours (00:00, 08:00, 16:00 UTC)

---

## WebSocket Streams

### Market Data Stream
```go
stream := client.NewMarketDataStream()
// URL: wss://open-api-swap.bingx.com/swap-market

stream.Connect()
stream.SubscribeTrade("BTC-USDT")
stream.SubscribeKline("BTC-USDT", "1m")
stream.SubscribeDepth("BTC-USDT", 20)
stream.SubscribeTicker("BTC-USDT")
stream.SubscribeBookTicker("BTC-USDT")

stream.UnsubscribeTrade("BTC-USDT")
stream.UnsubscribeKline("BTC-USDT", "1m")
stream.UnsubscribeDepth("BTC-USDT", 20)
stream.UnsubscribeTicker("BTC-USDT")
stream.UnsubscribeBookTicker("BTC-USDT")
```

### Account Data Stream
```go
// 1. Get listen key
lk, _ := client.ListenKey().Create()
listenKey := lk["listenKey"].(string)

// 2. Open stream
stream := client.NewAccountDataStream(listenKey)
stream.Connect()

// 3. Keep-alive (every ~30 min)
client.ListenKey().Extend(listenKey)

// 4. Cleanup
client.ListenKey().Delete(listenKey)
```

Stream events: order updates, balance changes, position changes.

---

## Error Handling

Package: `github.com/tigusigalpa/bingx-go/v2/errors`

```go
import bingxerr "github.com/tigusigalpa/bingx-go/v2/errors"

order, err := client.Trade().CreateOrder(params)
if err != nil {
    switch e := err.(type) {
    case *bingxerr.InsufficientBalanceException:
        // not enough funds
    case *bingxerr.RateLimitException:
        // slow down, add backoff
    case *bingxerr.AuthenticationException:
        // bad keys, wrong IP whitelist, clock skew
    case *bingxerr.APIException:
        fmt.Println(e.APICode, e.Message)
    case *bingxerr.BingXException:
        fmt.Println(e.Code, e.Message)
    default:
        // network / unknown
    }
}
```

**Error types hierarchy:**
- `BingXException` — base (Message, Code, Response)
  - `APIException` — adds APICode
  - `AuthenticationException`
  - `RateLimitException`
  - `InsufficientBalanceException`

---

## Key Patterns & Best Practices

### Credentials — never hardcode
```go
apiKey    := os.Getenv("BINGX_API_KEY")
apiSecret := os.Getenv("BINGX_API_SECRET")
```

### Rate Limiting
```go
import "golang.org/x/time/rate"
limiter := rate.NewLimiter(10, 1) // 10 req/sec
limiter.Wait(context.Background())
```

### Retry with backoff
```go
for attempt := 1; attempt <= 3; attempt++ {
    result, err := client.Market().GetLatestPrice(symbol)
    if err == nil { return result, nil }
    time.Sleep(time.Duration(attempt) * time.Second)
}
```

### Test orders before real orders
```go
_, err := client.Trade().CreateTestOrder(params)
if err != nil { /* order would fail */ return }
client.Trade().CreateOrder(params)
```

### Close all positions (panic button)
```go
positions, _ := client.Account().GetPositions(nil)
for _, pos := range positions {
    p := pos.(map[string]interface{})
    amt := p["positionAmt"].(float64)
    if amt == 0 { continue }
    side := "SELL"
    if amt < 0 { side = "BUY"; amt = -amt }
    client.Trade().CreateOrder(map[string]interface{}{
        "symbol": p["symbol"], "side": side,
        "type": "MARKET", "positionSide": p["positionSide"], "quantity": amt,
    })
}
```

---

## API Endpoints Reference

| Service | Prefix |
|---------|--------|
| USDT-M Perpetual Swap | `/openApi/swap/v2/` |
| Coin-M Perpetual Swap | `/openApi/cswap/v1/` |
| TradFi (Stocks, Forex, Commodities) | `/openApi/swap/v2/` |
| Spot | `/openApi/spot/v1/` |
| Wallet | `/openApi/wallets/v1/` |
| Sub-Account | `/openApi/subAccount/v1/` |
| Copy Trading | `/openApi/copy/v1/` |
| Base URL (live) | `https://open-api.bingx.com` |
| Base URL (demo) | `https://open-api-vst.bingx.com` |
| WebSocket (market) | `wss://open-api-swap.bingx.com/swap-market` |

---

## What Requires Auth vs Public

| Category | Auth Required |
|----------|:---:|
| Market data (prices, klines, depth, OI) | No |
| Symbols list | No |
| TradFi market data | No |
| Account balance & positions | Yes |
| Place / cancel / modify orders | Yes |
| TradFi trading | Yes |
| TWAP orders | Yes |
| Listen keys | Yes |
| Wallet deposit/withdraw | Yes |
| Sub-account management | Yes |
| Copy trading | Yes |

---

## v3 Feature Summary (API v3 / SDK v1.4+)

| Feature | Method |
|---------|--------|
| TWAP Orders | `Trade().PlaceTWAPOrder()` / `GetTWAPOrder()` / `CancelTWAPOrder()` |
| Trailing Stop | `OrderTypeTrailingStopMarket` in `CreateOrder()` |
| Trailing TP/SL | `OrderTypeTrailingTPSL` in `CreateOrder()` |
| Trigger Limit | `OrderTypeTriggerLimit` in `CreateOrder()` |
| One-Click Reversal | `Trade().OneClickReversePosition()` |
| Multi-Assets Margin | `Trade().SwitchMultiAssetsMode()` / `GetMultiAssetsMode()` |
| Auto Add Margin | `Trade().SetAutoAddMargin()` |
| Position Risk | `Account().GetPositionRisk()` |
| Income History | `Account().GetIncomeHistory()` |
| Commission History | `Account().GetCommissionHistory()` |
| Liquidation History | `Account().GetForceOrders()` |
| Open Interest | `Market().GetOpenInterest()` / `GetOpenInterestHistory()` |
| Funding Rate Info | `Market().GetFundingRateInfo()` |
| Book Ticker | `Market().GetBookTicker()` / `GetSpotBookTicker()` |
| Index / Ticker Price | `Market().GetIndexPrice()` / `GetTickerPrice()` |
| Long/Short Ratio | `Market().GetLongShortRatio()` |
| Aggregate Trades | `Market().GetAggregateTrades()` |
| Hedge Mode | `Account().GetPositionMode()` / `SetPositionMode()` |
| TradFi Support | `client.TradFi()` — Stocks, Forex, Commodities, Indices |
| TradFi Symbol Filter | `TradFi().Market().GetStockSymbols()` / `GetForexSymbols()` / `GetCommoditySymbols()` |

---

## Links

- GitHub: https://github.com/tigusigalpa/bingx-go
- Wiki: https://github.com/tigusigalpa/bingx-go/wiki
- BingX Official API Docs: https://bingx-api.github.io/docs-v3/
- Migration Guide (v3): `API_V3_MIGRATION.md`
- Issues: https://github.com/tigusigalpa/bingx-go/issues
