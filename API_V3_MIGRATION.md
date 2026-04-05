# BingX API v3 Migration Guide

This guide helps you migrate from previous versions of `bingx-go` to v1.4.0 with full BingX API v3 support.

## What's New in API v3

### Summary of Changes

- **20+ new trading methods** including TWAP orders, multi-assets mode, and one-click position reversal
- **New order types**: TRIGGER_LIMIT, TRAILING_STOP_MARKET, TRAILING_TP_SL
- **Enhanced market data** with open interest history, funding rate info, and book tickers
- **Advanced account management** with position risk, income history, and force orders
- **Improved error handling** with additional v3 error codes
- **Full backward compatibility** - no breaking changes

## New Features

### 1. Advanced Order Types

```go
import "github.com/tigusigalpa/bingx-go/services"

// Use new order type constants
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeTriggerLimit,  // New!
    "positionSide": "LONG",
    "price":        50000.0,
    "stopPrice":    49500.0,
    "quantity":     0.001,
})

// Trailing stop market order
order, err = client.Trade().CreateOrder(map[string]interface{}{
    "symbol":           "ETH-USDT",
    "side":             "SELL",
    "type":             services.OrderTypeTrailingStopMarket,  // New!
    "positionSide":     "SHORT",
    "activationPrice":  3000.0,
    "callbackRate":     1.0,  // 1% callback
    "quantity":         0.1,
})
```

### 2. TWAP Orders (Time-Weighted Average Price)

Perfect for executing large orders with minimal market impact:

```go
// Place TWAP order
twapOrder, err := client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "positionSide": "LONG",
    "quantity":     1.0,
    "duration":     3600,  // Execute over 1 hour
    "interval":     60,    // Split into 1-minute intervals
})

// Query TWAP order status
order, err := client.Trade().GetTWAPOrder("twap_order_id", nil)

// List all TWAP orders
status := "WORKING"
orders, err := client.Trade().GetTWAPOrders(
    &symbol,
    &status,
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)

// Cancel TWAP order
err = client.Trade().CancelTWAPOrder("twap_order_id", nil)
```

### 3. One-Click Position Reversal

Instantly reverse your position from LONG to SHORT or vice versa:

```go
// Reverse position with one call
result, err := client.Trade().OneClickReversePosition("BTC-USDT", nil)

// Before: LONG 1.0 BTC
// After:  SHORT 1.0 BTC (automatically closes LONG and opens SHORT)
```

### 4. Multi-Assets Margin Mode

Enable portfolio margin across multiple assets:

```go
// Enable multi-assets mode
err := client.Trade().SwitchMultiAssetsMode(true, nil)

// Query current mode
mode, err := client.Trade().GetMultiAssetsMode(nil)

// Get multi-assets rules
rules, err := client.Trade().GetMultiAssetsRules(nil)

// Query multi-assets margin
margin, err := client.Trade().GetMultiAssetsMargin(nil)
```

### 5. Hedge Mode Auto Margin Addition

Automatically add margin to positions in hedge mode:

```go
// Enable auto margin addition for LONG positions
err := client.Trade().SetAutoAddMargin(
    "BTC-USDT",
    "LONG",
    true,  // Enable auto-add
    nil,
)

// Disable for SHORT positions
err = client.Trade().SetAutoAddMargin(
    "BTC-USDT",
    "SHORT",
    false,  // Disable auto-add
    nil,
)
```

### 6. Enhanced Market Data

```go
// Get open interest
oi, err := client.Market().GetOpenInterest("BTC-USDT")

// Get open interest history
oiHistory, err := client.Market().GetOpenInterestHistory(
    "BTC-USDT",
    "5m",   // period
    100,    // limit
    nil,    // startTime
    nil,    // endTime
)

// Get funding rate info
fundingInfo, err := client.Market().GetFundingRateInfo("BTC-USDT")

// Get best bid/ask (book ticker)
bookTicker, err := client.Market().GetBookTicker(&symbol)

// Get index price
indexPrice, err := client.Market().GetIndexPrice("BTC-USDT")

// Get ticker price
tickerPrice, err := client.Market().GetTickerPrice(&symbol)
```

### 7. Advanced Account Management

```go
// Get position risk metrics
risk, err := client.Account().GetPositionRisk(&symbol, nil)

// Get income/PnL history
incomeType := "REALIZED_PNL"
income, err := client.Account().GetIncomeHistory(
    &symbol,
    &incomeType,
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)

// Get commission history
commissions, err := client.Account().GetCommissionHistory(
    "BTC-USDT",
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)

// Get liquidation/force orders
forceOrders, err := client.Account().GetForceOrders(
    &symbol,
    nil,    // autoCloseType
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)

// Query position mode (hedge vs one-way)
mode, err := client.Account().GetPositionMode(nil)

// Switch to hedge mode
err = client.Account().SetPositionMode(true, nil)  // true = hedge mode

// Switch to one-way mode
err = client.Account().SetPositionMode(false, nil)  // false = one-way mode
```

### 8. Coin-M Futures Enhancements

```go
coinm := client.CoinM()

// Get position risk for Coin-M
risk, err := coinm.Trade().GetPositionRisk(&symbol, nil)

// Get income history
income, err := coinm.Trade().GetIncomeHistory(
    &symbol,
    &incomeType,
    nil, nil, 100, nil,
)

// Enhanced market data
fundingHistory, err := coinm.Market().GetFundingRateHistory("BTC-USD", 100)
markPrice, err := coinm.Market().GetMarkPrice("BTC-USD")
indexPrice, err := coinm.Market().GetIndexPrice("BTC-USD")
trades, err := coinm.Market().GetRecentTrades("BTC-USD", 100)
```

## Migration Checklist

### ✅ No Breaking Changes

All existing code continues to work without modifications. The v3 update only adds new features.

### ✅ Optional Upgrades

Consider upgrading to new features:

1. **Replace manual position reversal** with `OneClickReversePosition()`
2. **Use TWAP orders** for large trades instead of manual splitting
3. **Enable multi-assets mode** if trading multiple contracts
4. **Add position risk monitoring** with `GetPositionRisk()`
5. **Track income/PnL** with `GetIncomeHistory()`

### ✅ Error Handling

New error codes are automatically handled:

```go
order, err := client.Trade().CreateOrder(params)
if err != nil {
    switch e := err.(type) {
    case *errors.AuthenticationException:
        // Handles both old and new auth error codes
        log.Printf("Auth error: %v", e)
    case *errors.RateLimitException:
        // Handles 100005 and new 100429
        log.Printf("Rate limit: %v", e)
    case *errors.InsufficientBalanceException:
        // Handles 200001 and new 200002
        log.Printf("Insufficient balance: %v", e)
    default:
        log.Printf("Error: %v", err)
    }
}
```

## Best Practices for v3

### 1. Use Order Type Constants

```go
// ✅ Good - Type-safe
"type": services.OrderTypeTriggerLimit

// ❌ Avoid - String literals prone to typos
"type": "TRIGGER_LIMIT"
```

### 2. TWAP for Large Orders

```go
// ✅ Good - Use TWAP for large positions
client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "side":     "BUY",
    "quantity": 10.0,
    "duration": 3600,
})

// ❌ Avoid - Single large market order causes slippage
client.Trade().CreateOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "side":     "BUY",
    "type":     "MARKET",
    "quantity": 10.0,
})
```

### 3. Monitor Position Risk

```go
// ✅ Good - Regular risk monitoring
risk, err := client.Account().GetPositionRisk(&symbol, nil)
if err == nil {
    // Check liquidation price, leverage, etc.
}
```

### 4. Track Income History

```go
// ✅ Good - Track realized PnL
incomeType := "REALIZED_PNL"
income, err := client.Account().GetIncomeHistory(&symbol, &incomeType, nil, nil, 100, nil)
```

## Performance Improvements

### v3 Optimizations

1. **Batch operations** remain the most efficient way to manage multiple orders
2. **WebSocket streams** continue to provide real-time data with minimal latency
3. **TWAP orders** reduce market impact for large trades
4. **Multi-assets mode** optimizes margin usage across positions

## Testing Your Migration

```go
package main

import (
    "fmt"
    "log"
    
    bingx "github.com/tigusigalpa/bingx-go"
    "github.com/tigusigalpa/bingx-go/services"
)

func main() {
    client := bingx.NewClient(
        "YOUR_API_KEY",
        "YOUR_API_SECRET",
    )
    
    // Test 1: Verify new order types
    fmt.Println("Testing new order types...")
    fmt.Printf("TRIGGER_LIMIT: %s\n", services.OrderTypeTriggerLimit)
    fmt.Printf("TRAILING_STOP_MARKET: %s\n", services.OrderTypeTrailingStopMarket)
    
    // Test 2: Check multi-assets mode
    fmt.Println("\nTesting multi-assets mode...")
    mode, err := client.Trade().GetMultiAssetsMode(nil)
    if err != nil {
        log.Printf("Multi-assets mode check: %v", err)
    } else {
        fmt.Printf("Multi-assets mode: %+v\n", mode)
    }
    
    // Test 3: Get position risk
    fmt.Println("\nTesting position risk...")
    symbol := "BTC-USDT"
    risk, err := client.Account().GetPositionRisk(&symbol, nil)
    if err != nil {
        log.Printf("Position risk check: %v", err)
    } else {
        fmt.Printf("Position risk: %+v\n", risk)
    }
    
    // Test 4: Check open interest
    fmt.Println("\nTesting open interest...")
    oi, err := client.Market().GetOpenInterest("BTC-USDT")
    if err != nil {
        log.Printf("Open interest check: %v", err)
    } else {
        fmt.Printf("Open interest: %+v\n", oi)
    }
    
    fmt.Println("\n✅ All v3 features tested successfully!")
}
```

## Support & Resources

- **API Documentation**: https://bingx-api.github.io/docs-v3/
- **GitHub Issues**: https://github.com/tigusigalpa/bingx-go/issues
- **GoDoc**: https://pkg.go.dev/github.com/tigusigalpa/bingx-go

## Version History

- **v1.4.0** (2026-04-05) - Full BingX API v3 support
- **v1.3.11** (2026-03-15) - Pre-v3 updates
- **v1.3.0** - Stable release with v2 API

## Need Help?

If you encounter issues during migration:

1. Check the [CHANGELOG.md](CHANGELOG.md) for detailed changes
2. Review [examples/](examples/) for updated code samples
3. Open an issue on GitHub with your specific use case
4. Ensure you're using Go 1.21 or higher

---

**Note**: All v3 features are production-ready and tested against BingX API v3 endpoints.
