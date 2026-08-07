# Account Management

Your account is where the money lives. This guide covers checking balances, managing positions, adjusting leverage, and keeping an eye on risk.

## What's Here

- [Your Balance](#your-balance)
- [Your Positions](#your-positions)
- [Leverage & Margin](#leverage--margin)
- [Risk Monitoring](#risk-monitoring-v3) *(v3)*
- [P&L Tracking](#pnl-tracking-v3) *(v3)*
- [Account Settings](#account-settings)

---

## Your Balance

### How Much Do I Have?

```go
balance, err := client.Account().GetBalance()

fmt.Printf("Available: %v USDT\n", balance["availableBalance"])
fmt.Printf("Total: %v USDT\n", balance["balance"])
fmt.Printf("Unrealized P&L: %v USDT\n", balance["unrealizedProfit"])
```

### Full Account Overview

```go
info, err := client.Account().GetAccountInfo()

fmt.Printf("Equity: %v\n", info["totalEquity"])
fmt.Printf("Used Margin: %v\n", info["totalMargin"])
fmt.Printf("Free Margin: %v\n", info["availableMargin"])
fmt.Printf("Margin Level: %v%%\n", info["marginLevel"])
```

### Low Balance Alert

Don't get caught with an empty account:

```go
func watchBalance(client *bingx.Client, minBalance float64) {
    for {
        balance, err := client.Account().GetBalance()
        if err != nil {
            time.Sleep(30 * time.Second)
            continue
        }
        
        available := balance["availableBalance"].(float64)
        
        if available < minBalance {
            fmt.Printf("🚨 Low balance! Only %.2f USDT left\n", available)
            // Send alert, pause trading, etc.
        }
        
        time.Sleep(30 * time.Second)
    }
}
```

---

## Your Positions

### See Everything You're Holding

```go
positions, err := client.Account().GetPositions(nil)

for _, pos := range positions {
    p := pos.(map[string]interface{})
    
    // Skip empty positions
    if p["positionAmt"].(float64) == 0 {
        continue
    }
    
    pnl := p["unrealizedProfit"].(float64)
    emoji := "🟢"
    if pnl < 0 {
        emoji = "🔴"
    }
    
    fmt.Printf("%s %s %s: %.4f @ %.2f (%s %.2f)\n",
        emoji,
        p["positionSide"],
        p["symbol"],
        p["positionAmt"],
        p["avgPrice"],
        emoji,
        pnl)
}
```

### Check a Specific Position

```go
symbol := "BTC-USDT"
positions, err := client.Account().GetPositions(&symbol)
```

### Quick Portfolio Summary

```go
func portfolioSnapshot(client *bingx.Client) {
    positions, _ := client.Account().GetPositions(nil)
    
    var totalPnL float64
    var longs, shorts int
    
    for _, pos := range positions {
        p := pos.(map[string]interface{})
        amt := p["positionAmt"].(float64)
        
        if amt == 0 {
            continue
        }
        
        totalPnL += p["unrealizedProfit"].(float64)
        
        if amt > 0 {
            longs++
        } else {
            shorts++
        }
    }
    
    fmt.Printf("Portfolio: %d longs, %d shorts, P&L: %.2f USDT\n", 
        longs, shorts, totalPnL)
}
```

---

## Leverage & Margin

Leverage is a double-edged sword. Higher leverage = bigger gains AND bigger losses.

### Check Current Leverage

```go
leverage, _ := client.Account().GetLeverage("BTC-USDT", nil)
fmt.Printf("Currently using %vx leverage\n", leverage)
```

### Change Leverage

```go
// Set 20x for all positions on BTC
err := client.Account().SetLeverage("BTC-USDT", "BOTH", 20, nil)

// Or set different leverage for longs vs shorts
client.Account().SetLeverage("BTC-USDT", "LONG", 10, nil)
client.Account().SetLeverage("BTC-USDT", "SHORT", 5, nil)
```

### Margin Mode: Isolated vs Cross

**Isolated**: Each position has its own margin. If liquidated, only that position's margin is lost.

**Cross**: All positions share margin. More capital efficient, but one bad trade can affect everything.

```go
// Check current mode
mode, _ := client.Account().GetMarginMode("BTC-USDT")
fmt.Printf("Mode: %v\n", mode) // ISOLATED or CROSSED

// Switch modes
client.Account().SetMarginMode("BTC-USDT", "ISOLATED")
client.Account().SetMarginMode("BTC-USDT", "CROSSED")
```

### Hedge Mode vs One-Way *(v3)*

**One-way**: One position per symbol. Buy closes your short.

**Hedge**: Can hold long AND short simultaneously. Good for hedging.

```go
// Check current mode
mode, _ := client.Account().GetPositionMode(nil)

// Enable hedge mode
client.Account().SetPositionMode(true, nil)

// Back to one-way
client.Account().SetPositionMode(false, nil)
```

> ⚠️ You can't switch modes while you have open positions.

---

## Risk Monitoring *(v3)*

Knowing your liquidation price before you get liquidated is... important.

### Check Your Risk

```go
symbol := "BTC-USDT"
risk, _ := client.Account().GetPositionRisk(&symbol, nil)

fmt.Printf("BTC Position Risk:\n")
fmt.Printf("  Liquidation at: $%.2f\n", risk["liquidationPrice"])
fmt.Printf("  Margin ratio: %.1f%%\n", risk["marginRatio"].(float64)*100)
fmt.Printf("  Unrealized P&L: $%.2f\n", risk["unrealizedProfit"])
```

### Simple Risk Alert

```go
func checkRisk(client *bingx.Client, symbol string) {
    risk, err := client.Account().GetPositionRisk(&symbol, nil)
    if err != nil {
        return
    }
    
    marginRatio := risk["marginRatio"].(float64)
    liqPrice := risk["liquidationPrice"].(float64)
    
    if marginRatio > 0.8 {
        fmt.Printf("🚨 DANGER! %s margin at %.0f%%, liq price: %.2f\n", 
            symbol, marginRatio*100, liqPrice)
    } else if marginRatio > 0.5 {
        fmt.Printf("⚠️ Warning: %s margin at %.0f%%\n", symbol, marginRatio*100)
    }
}
```

### Auto Risk Management

Here's a more complete example that actually does something when risk gets high:

```go
func manageRisk(client *bingx.Client, symbol string, maxRatio float64) {
    risk, _ := client.Account().GetPositionRisk(&symbol, nil)
    
    marginRatio := risk["marginRatio"].(float64)
    if marginRatio <= maxRatio {
        return // We're fine
    }
    
    // Risk too high - reduce position by half
    posAmt := risk["positionAmt"].(float64)
    if posAmt == 0 {
        return
    }
    
    reduceBy := math.Abs(posAmt) * 0.5
    side := "SELL"
    if posAmt < 0 {
        side = "BUY"
    }
    
    fmt.Printf("🚨 Reducing %s position by 50%%...\n", symbol)
    
    client.Trade().CreateOrder(map[string]interface{}{
        "symbol":       symbol,
        "side":         side,
        "type":         "MARKET",
        "positionSide": risk["positionSide"],
        "quantity":     reduceBy,
    })
}
```

---

## P&L Tracking *(v3)*

Want to know if you're actually making money? These endpoints tell you exactly where your profits (and losses) are coming from.

### Your Trading History

```go
symbol := "BTC-USDT"
incomeType := "REALIZED_PNL"

income, _ := client.Account().GetIncomeHistory(
    &symbol,
    &incomeType,
    nil, nil, 100, nil,
)

var total float64
for _, record := range income {
    r := record.(map[string]interface{})
    total += r["income"].(float64)
}

if total >= 0 {
    fmt.Printf("🟢 You've made %.2f USDT on BTC\n", total)
} else {
    fmt.Printf("🔴 You've lost %.2f USDT on BTC\n", -total)
}
```

### What Gets Tracked

- `REALIZED_PNL` — Actual profit/loss from closed trades
- `FUNDING_FEE` — Those 8-hour funding payments
- `COMMISSION` — Trading fees
- `TRANSFER` — Money moved between accounts
- `INSURANCE_CLEAR` — Insurance fund (rare)

### How Much Are You Paying in Fees?

```go
commissions, _ := client.Account().GetCommissionHistory(
    "BTC-USDT", nil, nil, 100, nil,
)

var total float64
for _, c := range commissions {
    total += c.(map[string]interface{})["commission"].(float64)
}

fmt.Printf("You've paid %.4f in fees\n", total)
```

### Liquidation History

Hopefully this list is empty, but if you've been liquidated, here's how to see it:

```go
symbol := "BTC-USDT"
liquidations, _ := client.Account().GetForceOrders(
    &symbol, nil, nil, nil, 100, nil,
)

if len(liquidations) == 0 {
    fmt.Println("✅ No liquidations - nice!")
} else {
    fmt.Printf("💥 You've been liquidated %d times\n", len(liquidations))
}
```

### Trading Performance Report

Here's a handy function that gives you a full breakdown:

```go
func tradingReport(client *bingx.Client, symbol string, days int) {
    endTime := time.Now().UnixMilli()
    startTime := time.Now().AddDate(0, 0, -days).UnixMilli()
    incomeType := "REALIZED_PNL"
    
    income, _ := client.Account().GetIncomeHistory(
        &symbol, &incomeType, &startTime, &endTime, 1000, nil,
    )
    
    var wins, losses int
    var totalWins, totalLosses float64
    
    for _, r := range income {
        pnl := r.(map[string]interface{})["income"].(float64)
        if pnl > 0 {
            wins++
            totalWins += pnl
        } else if pnl < 0 {
            losses++
            totalLosses += pnl
        }
    }
    
    total := wins + losses
    if total == 0 {
        fmt.Println("No trades in this period")
        return
    }
    
    winRate := float64(wins) / float64(total) * 100
    netPnL := totalWins + totalLosses
    
    fmt.Printf("📊 %s Performance (%d days):\n", symbol, days)
    fmt.Printf("   Trades: %d (%d wins, %d losses)\n", total, wins, losses)
    fmt.Printf("   Win Rate: %.1f%%\n", winRate)
    fmt.Printf("   Net P&L: %.2f USDT\n", netPnL)
    
    if totalLosses != 0 {
        profitFactor := totalWins / math.Abs(totalLosses)
        fmt.Printf("   Profit Factor: %.2f\n", profitFactor)
    }
}
```

---

## Account Settings

### Balance History

```go
history, _ := client.Account().GetBalanceHistory(100, nil, nil)

for _, h := range history {
    record := h.(map[string]interface{})
    fmt.Printf("%v: %v (change: %v)\n",
        record["updateTime"], record["balance"], record["balanceChange"])
}
```

### Check Your API Permissions

```go
perms, _ := client.Account().GetAPIPermissions()

fmt.Printf("Can trade: %v\n", perms["enableTrading"])
fmt.Printf("Can withdraw: %v\n", perms["enableWithdrawals"])
fmt.Printf("Can read: %v\n", perms["enableReading"])
```

---

## Putting It All Together

### Health Check Function

Run this periodically to catch problems early:

```go
func healthCheck(client *bingx.Client) error {
    // Check balance
    balance, err := client.Account().GetBalance()
    if err != nil {
        return fmt.Errorf("can't get balance: %w", err)
    }
    
    available := balance["availableBalance"].(float64)
    if available < 100 {
        return fmt.Errorf("🚨 low balance: %.2f USDT", available)
    }
    
    // Check for risky positions
    positions, _ := client.Account().GetPositions(nil)
    for _, pos := range positions {
        p := pos.(map[string]interface{})
        if p["positionAmt"].(float64) == 0 {
            continue
        }
        
        symbol := p["symbol"].(string)
        risk, _ := client.Account().GetPositionRisk(&symbol, nil)
        
        if risk["marginRatio"].(float64) > 0.7 {
            return fmt.Errorf("🚨 %s margin at %.0f%%", 
                symbol, risk["marginRatio"].(float64)*100)
        }
    }
    
    return nil
}
```

### Simple Risk Rules

```go
func canOpenNewPosition(client *bingx.Client) bool {
    const maxPositions = 5
    
    positions, _ := client.Account().GetPositions(nil)
    
    active := 0
    for _, p := range positions {
        if p.(map[string]interface{})["positionAmt"].(float64) != 0 {
            active++
        }
    }
    
    if active >= maxPositions {
        fmt.Printf("⚠️ Already have %d positions, max is %d\n", active, maxPositions)
        return false
    }
    
    return true
}
```

---

## Keep Going

- **[Trading](Trade-Service)** — Place and manage orders
- **[Market Data](Market-Service)** — Prices, charts, order books
- **[v3 Features](API-v3-Features)** — TWAP, trailing stops, and more
- **[Error Handling](Error-Handling)** — When things go wrong

---

**Questions?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues) — we're here to help.
