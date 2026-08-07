# Trading

This is where the action happens. Here's everything you need to place orders, manage positions, and use the fancy new v3 features like TWAP and trailing stops.

## What's Here

- [Basic Orders](#basic-orders) — Market, limit, stop-loss
- [Advanced Orders](#advanced-order-types-v3) — Trailing stops, trigger limits *(v3)*
- [TWAP Orders](#twap-orders) — Execute large orders over time *(v3)*
- [Managing Orders](#order-management) — Query, modify, cancel
- [Position Control](#position-management) — Close, reverse, manage
- [Leverage & Margin](#leverage--margin) — Adjust your risk

---

## Basic Orders

### Market Order — "Just buy it now"

Executes immediately at whatever the current price is. Fast, but you might get slippage on large orders.

```go
import "github.com/tigusigalpa/bingx-go/v2/services"

order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeMarket,
    "positionSide": "LONG",
    "quantity":     0.001,
})
```

### Limit Order — "Buy it, but only at my price"

Only executes if the market reaches your price. You might wait a while (or forever).

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeLimit,
    "positionSide": "LONG",
    "price":        50000.0,    // Your target price
    "quantity":     0.001,
    "timeInForce":  "GTC",      // Good 'til canceled
})
```

**Time in Force options:**
- `GTC` — Stays open until filled or canceled
- `IOC` — Fill what you can immediately, cancel the rest
- `FOK` — Fill everything or nothing

### Stop-Loss — "Get me out if it drops"

Protect yourself from big losses:

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "SELL",
    "type":         services.OrderTypeStopMarket,
    "positionSide": "LONG",
    "stopPrice":    48000.0,  // If price hits this, sell immediately
    "quantity":     0.001,
})
```

### Take Profit — "Lock in my gains"

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "SELL",
    "type":         services.OrderTypeTakeProfitMarket,
    "positionSide": "LONG",
    "stopPrice":    52000.0,  // If price hits this, sell and take profit
    "quantity":     0.001,
})
```

---

## Advanced Order Types *(v3)*

These are the new toys in v3. They let you build more sophisticated strategies.

### Trigger Limit — "Wait for my price, then use a limit"

Like a stop order, but executes as a limit instead of market. Better price control.

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeTriggerLimit,
    "positionSide": "LONG",
    "stopPrice":    49500.0,    // When price hits this...
    "price":        50000.0,    // ...place a limit order at this
    "quantity":     0.1,
})
```

**When to use**: You want to buy the dip, but you don't want to chase it if it bounces too fast.

### Trailing Stop — "Follow the price up, catch me when it falls"

This is the one everyone asks about. The stop price moves with the market, locking in gains as price rises.

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":           "ETH-USDT",
    "side":             "SELL",
    "type":             services.OrderTypeTrailingStopMarket,
    "positionSide":     "LONG",
    "activationPrice":  3000.0,    // Start trailing at this price
    "callbackRate":     1.0,       // Stay 1% behind the high
    "quantity":         1.0,
})
```

**Here's how it plays out:**

| Price | Stop Price | What Happens |
|-------|------------|-------------|
| 2900 | — | Waiting to activate |
| 3000 | 2970 | Activated! Stop is 1% below |
| 3100 | 3069 | Stop moves up |
| 3200 | 3168 | Stop moves up again |
| 3168 | 3168 | **SOLD!** Locked in $168 of gains |

### Trailing TP/SL — "Set it and forget it"

Take-profit and stop-loss in one order, with trailing behavior:

```go
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":               "BTC-USDT",
    "side":                 "BUY",
    "type":                 services.OrderTypeTrailingTPSL,
    "positionSide":         "LONG",
    "quantity":             0.1,
    "takeProfitPrice":      52000.0,
    "stopLossPrice":        48000.0,
    "trailingStopPercent":  0.5,
})
```

---

## TWAP Orders *(v3)*

Got a big order that would move the market? TWAP breaks it into smaller pieces and executes them over time. Institutional traders use this all the time.

### Placing a TWAP Order

```go
// Buy 10 BTC over 1 hour, executing a piece every minute
twap, err := client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "positionSide": "LONG",
    "quantity":     10.0,
    "duration":     3600,      // 1 hour total
    "interval":     60,        // Execute every 60 seconds
})

fmt.Printf("TWAP started: %v\n", twap["orderId"])
```

### Checking Progress

```go
twap, _ := client.Trade().GetTWAPOrder("order_id", nil)
fmt.Printf("Status: %v, Progress: %v%%\n", twap["status"], twap["progress"])
```

### Canceling a TWAP

```go
err := client.Trade().CancelTWAPOrder("order_id", nil)
```

### Complete TWAP Example

```go
func buyWithTWAP(client *bingx.Client, symbol string, qty float64) {
    // Spread over 1 hour, execute every 2 minutes
    twap, err := client.Trade().PlaceTWAPOrder(map[string]interface{}{
        "symbol":       symbol,
        "side":         "BUY",
        "positionSide": "LONG",
        "quantity":     qty,
        "duration":     3600,
        "interval":     120,
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    orderId := twap["orderId"].(string)
    fmt.Printf("📈 TWAP started for %.2f %s\n", qty, symbol)
    
    // Check progress every 30 seconds
    for {
        time.Sleep(30 * time.Second)
        
        status, _ := client.Trade().GetTWAPOrder(orderId, nil)
        
        if status["status"] == "FINISHED" {
            fmt.Println("✅ TWAP complete!")
            return
        }
        
        fmt.Printf("   Progress: %v%%\n", status["progress"])
    }
}
```

---

## Order Management

### Check an Order

```go
order, _ := client.Trade().GetOrder("BTC-USDT", "order_id")
fmt.Printf("Status: %v, Filled: %v\n", order["status"], order["executedQty"])
```

### See All Open Orders

```go
// Everything that's waiting to fill
openOrders, _ := client.Trade().GetOpenOrders(nil, 100)
fmt.Printf("You have %d open orders\n", len(openOrders))

// Just for one symbol
symbol := "BTC-USDT"
btcOrders, _ := client.Trade().GetOpenOrders(&symbol, 100)
```

### Order History

```go
symbol := "BTC-USDT"
history, _ := client.Trade().GetOrderHistory(&symbol, 100, nil, nil)

for _, order := range history {
    o := order.(map[string]interface{})
    fmt.Printf("%s %s @ %v\n", o["side"], o["symbol"], o["price"])
}
```

### Cancel an Order

```go
// By order ID
orderID := "123456789"
client.Trade().CancelOrder("BTC-USDT", &orderID, nil, nil, nil)

// By your custom client order ID
clientID := "my_order_001"
client.Trade().CancelOrder("BTC-USDT", nil, &clientID, nil, nil)
```

### Cancel Everything

```go
// Cancel all BTC orders
symbol := "BTC-USDT"
client.Trade().CancelAllOrders(nil, &symbol, nil, nil)

// Nuclear option: cancel ALL orders
client.Trade().CancelAllOrders(nil, nil, nil, nil)
```

### Modify an Order

Change price or quantity without canceling and re-placing:

```go
orderID := "123456789"
client.Trade().ModifyOrder(
    "BTC-USDT",
    &orderID,
    nil,
    51000.0,    // New price
    0.002,      // New quantity
    nil,
)
```

---

## Position Management

### Your Trade History

```go
symbol := "BTC-USDT"
trades, _ := client.Trade().GetUserTrades(&symbol, 100, nil, nil)

for _, t := range trades {
    trade := t.(map[string]interface{})
    fmt.Printf("%v %s @ %v\n", trade["qty"], trade["side"], trade["price"])
}
```

### Flip Your Position *(v3)*

Going from long to short (or vice versa) in one atomic operation:

```go
// LONG → SHORT or SHORT → LONG, instantly
client.Trade().OneClickReversePosition("BTC-USDT", nil)
```

No need to close then open — it happens in one step.

### Close Everything

Panic button:

```go
func closeAllPositions(client *bingx.Client) {
    positions, _ := client.Account().GetPositions(nil)
    
    for _, pos := range positions {
        p := pos.(map[string]interface{})
        amt := p["positionAmt"].(float64)
        
        if amt == 0 {
            continue
        }
        
        side := "SELL"
        if amt < 0 {
            side = "BUY"
            amt = -amt
        }
        
        client.Trade().CreateOrder(map[string]interface{}{
            "symbol":       p["symbol"],
            "side":         side,
            "type":         "MARKET",
            "positionSide": p["positionSide"],
            "quantity":     amt,
        })
        
        fmt.Printf("❌ Closed %s\n", p["symbol"])
    }
}
```

---

## Leverage & Margin

### Setting Leverage

```go
// 20x leverage on BTC
client.Trade().SetLeverage("BTC-USDT", 20, nil, nil)

// Different leverage for longs vs shorts
long := "LONG"
short := "SHORT"
client.Trade().SetLeverage("BTC-USDT", 10, &long, nil)
client.Trade().SetLeverage("BTC-USDT", 5, &short, nil)
```

### Auto Add Margin *(v3)*

Automatically add margin when you're getting close to liquidation:

```go
// Turn it on for longs
client.Trade().SetAutoAddMargin("BTC-USDT", "LONG", true, nil)

// Turn it off for shorts
client.Trade().SetAutoAddMargin("BTC-USDT", "SHORT", false, nil)
```

### Multi-Assets Mode *(v3)*

Use your entire portfolio as collateral instead of per-position margin:

```go
// Enable it
client.Trade().SwitchMultiAssetsMode(true, nil)

// Check status
mode, _ := client.Trade().GetMultiAssetsMode(nil)
margin, _ := client.Trade().GetMultiAssetsMargin(nil)
```

---

## Tips & Tricks

### Always Include Required Fields

```go
// ✅ This will work
order, _ := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         services.OrderTypeLimit,
    "positionSide": "LONG",
    "price":        50000.0,
    "quantity":     0.001,
    "timeInForce":  "GTC",
})

// ❌ This will fail - missing required params
order, _ = client.Trade().CreateOrder(map[string]interface{}{
    "symbol": "BTC-USDT",
    "side":   "BUY",
})
```

### Handle Errors Properly

```go
import "github.com/tigusigalpa/bingx-go/v2/errors"

order, err := client.Trade().CreateOrder(params)
if err != nil {
    switch err.(type) {
    case *errors.InsufficientBalanceException:
        fmt.Println("💸 Not enough money!")
    case *errors.RateLimitException:
        fmt.Println("⏳ Slow down, you're hitting rate limits")
    default:
        fmt.Printf("❌ Error: %v\n", err)
    }
    return
}
```

### Test Before You Trade

Use `CreateTestOrder` to validate your order without risking money:

```go
// This won't actually execute
test, err := client.Trade().CreateTestOrder(params)
if err != nil {
    fmt.Println("Order would fail:", err)
    return
}

// Now place the real one
real, _ := client.Trade().CreateOrder(params)
```

### Calculate Position Size Based on Risk

```go
func positionSize(balance, riskPct, entry, stop float64) float64 {
    riskAmount := balance * (riskPct / 100)
    priceDiff := math.Abs(entry - stop)
    return riskAmount / priceDiff
}

// Risk 2% of $10k account
qty := positionSize(10000, 2.0, 50000, 48000)
fmt.Printf("Buy %.4f BTC\n", qty) // ~0.001 BTC
```

---

## Putting It Together

Here's a complete example that buys BTC with a stop-loss:

```go
package main

import (
    "fmt"
    "os"
    
    bingx "github.com/tigusigalpa/bingx-go/v2"
    "github.com/tigusigalpa/bingx-go/v2/services"
)

func main() {
    client := bingx.NewClient(
        os.Getenv("BINGX_API_KEY"),
        os.Getenv("BINGX_API_SECRET"),
    )
    
    symbol := "BTC-USDT"
    
    // What's the current price?
    price, _ := client.Market().GetLatestPrice(symbol)
    currentPrice := price.(float64)
    fmt.Printf("BTC is at $%.2f\n", currentPrice)
    
    // Place a limit buy 1% below market
    buyPrice := currentPrice * 0.99
    order, err := client.Trade().CreateOrder(map[string]interface{}{
        "symbol":       symbol,
        "side":         "BUY",
        "type":         services.OrderTypeLimit,
        "positionSide": "LONG",
        "price":        buyPrice,
        "quantity":     0.001,
        "timeInForce":  "GTC",
    })
    
    if err != nil {
        fmt.Printf("❌ Order failed: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Buy order placed at $%.2f\n", buyPrice)
    
    // Set a 5% stop-loss
    stopPrice := currentPrice * 0.95
    client.Trade().CreateOrder(map[string]interface{}{
        "symbol":       symbol,
        "side":         "SELL",
        "type":         services.OrderTypeStopMarket,
        "positionSide": "LONG",
        "stopPrice":    stopPrice,
        "quantity":     0.001,
    })
    
    fmt.Printf("🛡️ Stop-loss set at $%.2f\n", stopPrice)
}
```

---

## Keep Going

- **[Account](Account-Service)** — Check balance and positions
- **[Market Data](Market-Service)** — Prices, charts, order books
- **[v3 Features](API-v3-Features)** — TWAP, trailing stops, and more
- **[Error Handling](Error-Handling)** — When things go wrong

---

**Questions?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues) — we're here to help.
