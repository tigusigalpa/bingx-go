# What's New in API v3

Version 1.4.0 is a big one. We've added 26 new methods that bring institutional-grade features to everyone — and the best part? Your existing code keeps working exactly as before.

Here's what you can do now:

- **TWAP Orders** — Execute large trades over time without moving the market
- **Multi-Assets Margin** — Use your whole portfolio as collateral
- **One-Click Reversal** — Flip from LONG to SHORT instantly
- **Better Risk Tools** — Track P&L, commissions, and liquidation risk
- **More Market Data** — Open interest history, funding rates, best bid/ask
- **Trailing Stops** — Dynamic stop-losses that follow price movements

---

## TWAP Orders

Ever tried to buy 10 BTC at once and watched the price spike against you? That's slippage, and it's the bane of large traders.

**TWAP (Time-Weighted Average Price)** solves this by breaking your big order into smaller pieces and executing them over time. Instead of one massive market-moving trade, you get dozens of small ones that blend into normal market activity.

**When to use it:**
- You're trading size that would move the market
- You want to build or unwind a position gradually
- You care more about average price than instant execution

### How It Works

```go
import bingx "github.com/tigusigalpa/bingx-go/v2"

client := bingx.NewClient(apiKey, apiSecret)

// Place TWAP order - Execute 10 BTC over 1 hour
twapOrder, err := client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "positionSide": "LONG",
    "quantity":     10.0,
    "duration":     3600,  // 1 hour in seconds
    "interval":     60,    // Execute every 60 seconds
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("TWAP Order ID: %v\n", twapOrder["orderId"])
```

### Query TWAP Order

```go
// Get specific TWAP order
twap, err := client.Trade().GetTWAPOrder("twap_order_id", nil)

// Get all TWAP orders
symbol := "BTC-USDT"
status := "WORKING"  // WORKING, FINISHED, CANCELLED
twapOrders, err := client.Trade().GetTWAPOrders(
    &symbol,
    &status,
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)
```

### Cancel TWAP Order

```go
err := client.Trade().CancelTWAPOrder("twap_order_id", nil)
```

### TWAP Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `symbol` | string | Yes | Trading pair (e.g., "BTC-USDT") |
| `side` | string | Yes | "BUY" or "SELL" |
| `positionSide` | string | Yes | "LONG" or "SHORT" |
| `quantity` | float64 | Yes | Total quantity to execute |
| `duration` | int | Yes | Total execution time in seconds |
| `interval` | int | Yes | Interval between executions in seconds |
| `price` | float64 | No | Limit price (optional) |

### Getting It Right

```go
// ✅ This makes sense - spread 100 BTC over 2 hours
client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "quantity": 100.0,
    "duration": 7200,  // 2 hours
    "interval": 120,   // A piece every 2 minutes
})

// ❌ This defeats the purpose - might as well use a regular order
client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "quantity": 100.0,
    "duration": 60,    // Only 1 minute total?
    "interval": 10,    // That's just 6 orders
})
```

**Rule of thumb**: If your TWAP would execute in under 10 minutes, you probably don't need TWAP.

---

## Multi-Assets Margin Mode

Normally, each position needs its own margin. Got a BTC long and an ETH long? That's two separate margin pools.

**Multi-Assets Mode** changes this. Your entire portfolio becomes the collateral, and positions can offset each other. A BTC long and a BTC short? They partially cancel out, so you need less margin overall.

**Why you'd want this:**
- Lower margin requirements (often significantly)
- More efficient use of your capital
- Better for portfolio strategies with hedged positions

### Turning It On

```go
// Enable multi-assets margin
err := client.Trade().SwitchMultiAssetsMode(true, nil)
if err != nil {
    log.Fatal(err)
}

// Disable multi-assets margin
err = client.Trade().SwitchMultiAssetsMode(false, nil)
```

### Query Multi-Assets Status

```go
// Get current mode
mode, err := client.Trade().GetMultiAssetsMode(nil)
fmt.Printf("Multi-Assets Mode: %v\n", mode)

// Get trading rules
rules, err := client.Trade().GetMultiAssetsRules(nil)
fmt.Printf("Rules: %v\n", rules)

// Get margin details
margin, err := client.Trade().GetMultiAssetsMargin(nil)
fmt.Printf("Margin: %v\n", margin)
```

### Example: Portfolio Margin Strategy

```go
// Enable multi-assets mode
client.Trade().SwitchMultiAssetsMode(true, nil)

// Open multiple positions with optimized margin
positions := []map[string]interface{}{
    {
        "symbol": "BTC-USDT",
        "side": "BUY",
        "quantity": 1.0,
    },
    {
        "symbol": "ETH-USDT",
        "side": "BUY",
        "quantity": 10.0,
    },
    {
        "symbol": "SOL-USDT",
        "side": "SELL",
        "quantity": 100.0,
    },
}

for _, pos := range positions {
    order, err := client.Trade().CreateOrder(pos)
    if err != nil {
        log.Printf("Error: %v", err)
        continue
    }
    fmt.Printf("Order placed: %v\n", order)
}

// Check total margin usage
margin, _ := client.Trade().GetMultiAssetsMargin(nil)
```

---

## One-Click Position Reversal

Sometimes the market turns and you need to flip your position *now*. Instead of closing your long and then opening a short (two orders, two chances for slippage), you can do it in one atomic operation.

**Perfect for:**
- Trend-following strategies that need to flip direction
- "Stop and reverse" systems
- Those "oh no" moments when you realize you're on the wrong side

### How to Use It

```go
// Reverse position: LONG → SHORT or SHORT → LONG
result, err := client.Trade().OneClickReversePosition("BTC-USDT", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Position reversed: %v\n", result)
```

### Example: Trend Reversal Strategy

```go
func reverseTrend(client *bingx.Client, symbol string) {
    // Get current position
    positions, err := client.Account().GetPositions(&symbol)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if we have a position
    if len(positions) > 0 {
        fmt.Println("Reversing position...")
        result, err := client.Trade().OneClickReversePosition(symbol, nil)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Position reversed: %v\n", result)
    } else {
        fmt.Println("No position to reverse")
    }
}
```

### What Actually Happens

Say you're long 1 BTC and you call `OneClickReversePosition`:

1. Your LONG 1.0 BTC gets closed
2. A new SHORT 1.0 BTC opens immediately
3. Both happen atomically — no gap where you're flat and exposed

It's like a U-turn for your position.

---

## New Order Types

v3 adds three order types that let you build more sophisticated strategies.

### TRIGGER_LIMIT

A trigger order that executes as a limit (not market) when the trigger price is hit. You get the automation of a trigger with the price control of a limit.

```go
import "github.com/tigusigalpa/bingx-go/v2/services"

order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeTriggerLimit,
    "positionSide": "LONG",
    "price":        50000.0,    // Limit price
    "stopPrice":    49500.0,    // Trigger price
    "quantity":     0.1,
})
```

### TRAILING_STOP_MARKET

This is the one everyone asks about. A trailing stop follows the price as it moves in your favor, then triggers when it reverses by a certain percentage.

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":           "ETH-USDT",
    "side":             "SELL",
    "type":             services.OrderTypeTrailingStopMarket,
    "positionSide":     "LONG",
    "activationPrice":  3000.0,    // Start trailing when price hits 3000
    "callbackRate":     1.0,       // Trail 1% behind the high
    "quantity":         1.0,
})
```

**Here's how it plays out:**
1. Price hits 3000 → trailing stop activates
2. Price rises to 3100 → stop moves to 3069 (1% below)
3. Price rises to 3200 → stop moves to 3168 (1% below)
4. Price drops to 3168 → **SOLD!**

You locked in most of that $200 move without having to guess the top.

### TRAILING_TP_SL

Take-profit and stop-loss in one order, with trailing behavior. Set it and forget it.

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":               "BTC-USDT",
    "side":                 "BUY",
    "type":                 services.OrderTypeTrailingTPSL,
    "positionSide":         "LONG",
    "quantity":             0.1,
    "takeProfitPrice":      52000.0,
    "stopLossPrice":        48000.0,
    "trailingStopPercent":  0.5,  // 0.5% trailing
})
```

---

## Position Risk Management

Knowing your liquidation price *before* you get liquidated is kind of important. The new risk endpoints give you everything you need to monitor your exposure.

### Checking Your Risk

```go
symbol := "BTC-USDT"
risk, err := client.Account().GetPositionRisk(&symbol, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Position Risk:\n")
fmt.Printf("  Liquidation Price: %v\n", risk["liquidationPrice"])
fmt.Printf("  Leverage: %v\n", risk["leverage"])
fmt.Printf("  Margin Ratio: %v\n", risk["marginRatio"])
fmt.Printf("  Unrealized PnL: %v\n", risk["unrealizedProfit"])
```

### Building a Risk Monitor

Here's a simple example that alerts you when things get spicy:

```go
func monitorRisk(client *bingx.Client, symbol string, maxMarginRatio float64) {
    risk, err := client.Account().GetPositionRisk(&symbol, nil)
    if err != nil {
        log.Printf("Couldn't get risk data: %v", err)
        return
    }
    
    marginRatio := risk["marginRatio"].(float64)
    
    if marginRatio > maxMarginRatio {
        log.Printf("🚨 Danger zone! Margin ratio at %.2f%% (limit: %.2f%%)", 
            marginRatio*100, maxMarginRatio*100)
        
        // Time to do something about it
        reducePosition(client, symbol)
    }
}
```

In production, you'd run this on a timer and maybe send alerts to Slack or Telegram.

---

## Income & PnL Tracking

Want to know how much you've actually made (or lost)? These endpoints break down your P&L by type — realized gains, funding fees, commissions, everything.

### Get Income History

```go
symbol := "BTC-USDT"
incomeType := "REALIZED_PNL"  // REALIZED_PNL, FUNDING_FEE, COMMISSION, etc.

income, err := client.Account().GetIncomeHistory(
    &symbol,
    &incomeType,
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)

for _, record := range income {
    fmt.Printf("Income: %v at %v\n", record["income"], record["time"])
}
```

### What Each Type Means

- `REALIZED_PNL` — Actual profit/loss from closed trades
- `FUNDING_FEE` — Those 8-hour funding payments (can be positive or negative!)
- `COMMISSION` — Trading fees you've paid
- `TRANSFER` — Money moved between accounts
- `INSURANCE_CLEAR` — Insurance fund stuff (rare)

### Commission History

```go
commissions, err := client.Account().GetCommissionHistory(
    "BTC-USDT",
    nil,    // startTime
    nil,    // endTime
    100,    // limit
    nil,    // recvWindow
)
```

### Quick P&L Summary

```go
func howAmIDoing(client *bingx.Client, symbol string) {
    incomeType := "REALIZED_PNL"
    income, err := client.Account().GetIncomeHistory(&symbol, &incomeType, nil, nil, 1000, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    var totalPnL float64
    for _, record := range income {
        totalPnL += record["income"].(float64)
    }
    
    if totalPnL >= 0 {
        fmt.Printf("🟢 You're up %.2f USDT on %s\n", totalPnL, symbol)
    } else {
        fmt.Printf("🔴 You're down %.2f USDT on %s\n", -totalPnL, symbol)
    }
}
```

---

## More Market Data

v3 gives you access to data that used to require expensive third-party services.

### Open Interest

Open interest tells you how much money is in the market. Rising OI with rising price? Bulls are piling in. Rising OI with falling price? Bears are getting aggressive.

```go
// Current open interest
oi, err := client.Market().GetOpenInterest("BTC-USDT")
fmt.Printf("Open Interest: %v\n", oi)

// Historical open interest
oiHistory, err := client.Market().GetOpenInterestHistory(
    "BTC-USDT",
    "5m",   // period: 5m, 15m, 30m, 1h, 4h, 1d
    100,    // limit
    nil,    // startTime
    nil,    // endTime
)
```

### Funding Rate Information

Funding rates can make or break your position if you're holding overnight. Now you can check them programmatically:

```go
fundingInfo, err := client.Market().GetFundingRateInfo("BTC-USDT")
fmt.Printf("Current rate: %v\n", fundingInfo["fundingRate"])
fmt.Printf("Next payment in: %v\n", fundingInfo["fundingTime"])

// Positive rate = longs pay shorts
// Negative rate = shorts pay longs
```

### Book Ticker (Best Bid/Ask)

Need the top of the order book without fetching the whole thing? Book ticker is your friend:

```go
symbol := "BTC-USDT"
ticker, err := client.Market().GetBookTickerData(&symbol)
if err != nil {
    return err
}

fmt.Println("Best bid:", ticker.BidPrice)
fmt.Println("Best ask:", ticker.AskPrice)

// GetBookTicker remains available when you need the raw BingX envelope.
// The current futures payload is nested at data.book_ticker with bid_price and ask_price.
```

### Index & Mark Price

```go
// Index price
indexPrice, err := client.Market().GetIndexPrice("BTC-USDT")

// Ticker price
tickerPrice, err := client.Market().GetTickerPrice(&symbol)
```

---

## Auto Add Margin

In hedge mode, you can tell BingX to automatically add margin to a position when it gets close to liquidation. It's like a safety net that uses your available balance to keep positions alive.

```go
// Enable auto-add margin for LONG positions
err := client.Trade().SetAutoAddMargin(
    "BTC-USDT",
    "LONG",
    true,   // enable
    nil,
)

// Disable for SHORT positions
err = client.Trade().SetAutoAddMargin(
    "BTC-USDT",
    "SHORT",
    false,  // disable
    nil,
)
```

---

## Position Mode: Hedge vs One-Way

**One-way mode**: You can only have one position per symbol. A buy closes your short (or opens a long).

**Hedge mode**: You can have both a long AND a short on the same symbol. Useful for hedging or complex strategies.

```go
// What mode am I in?
mode, err := client.Account().GetPositionMode(nil)

// Switch to hedge mode
err = client.Account().SetPositionMode(true, nil)

// Switch to one-way mode
err = client.Account().SetPositionMode(false, nil)
```

**Note**: You can't switch modes while you have open positions.

---

## Upgrading from Older Versions

Good news: **you don't have to change anything**. All your existing code keeps working.

v3 features are completely opt-in. Use them when you're ready:

```go
// Your old code still works fine
client := bingx.NewClient(apiKey, apiSecret)
balance, _ := client.Account().GetBalance()
order, _ := client.Trade().CreateOrder(params)

// When you're ready, start using v3 features
import "github.com/tigusigalpa/bingx-go/v2/services"

// Try a trailing stop
order, _ := client.Trade().CreateOrder(map[string]interface{}{
    "type": services.OrderTypeTrailingStopMarket,
    // ...
})

// Or TWAP for big orders
twap, _ := client.Trade().PlaceTWAPOrder(params)

// Or check your risk
risk, _ := client.Account().GetPositionRisk(&symbol, nil)
```

For the full migration guide with more examples, see [API_V3_MIGRATION.md](https://github.com/tigusigalpa/bingx-go/blob/main/API_V3_MIGRATION.md)

---

## Keep Exploring

- **[Trading Guide](Trade-Service)** — All the order types and trading operations
- **[Account Management](Account-Service)** — Balance, positions, settings
- **[Market Data](Market-Service)** — Prices, charts, order books
- **[Error Handling](Error-Handling)** — When things go wrong

---

**Got questions?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues) — we're happy to help.
