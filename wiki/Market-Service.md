# Market Data

Everything you need to know about fetching prices, charts, order books, and more. Most of these endpoints don't require authentication, so they're great for getting started.

## What's Covered

- [Available Trading Pairs](#available-trading-pairs)
- [Prices](#prices)
- [Order Book Depth](#order-book-depth)
- [Candlestick Charts](#candlestick-charts-klines)
- [24h Statistics](#24h-statistics)
- [Funding Rates](#funding-rates)
- [Open Interest](#open-interest-v3) *(v3)*
- [More Market Data](#more-market-data-v3) *(v3)*

---

## Available Trading Pairs

Before you can trade, you need to know what's available.

### Futures Pairs

```go
symbols, err := client.Market().GetFuturesSymbols()
if err != nil {
    log.Fatal(err)
}

for _, symbol := range symbols {
    fmt.Printf("Symbol: %v\n", symbol["symbol"])
    fmt.Printf("  Contract: %v\n", symbol["contractId"])
    fmt.Printf("  Status: %v\n", symbol["status"])
}
```

### Spot Pairs

```go
spotSymbols, err := client.Market().GetSpotSymbols()
```

### Everything at Once

```go
allSymbols, err := client.Market().GetAllSymbols()
```

---

## Prices

The most common thing you'll do — check what something is trading at.

### Current Price

```go
price, err := client.Market().GetLatestPrice("BTC-USDT")
fmt.Printf("BTC is at: %v\n", price)
```

### Check Multiple Pairs

```go
watchlist := []string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}

for _, symbol := range watchlist {
    price, err := client.Market().GetLatestPrice(symbol)
    if err != nil {
        log.Printf("Couldn't get %s: %v", symbol, err)
        continue
    }
    fmt.Printf("%s: %v\n", symbol, price)
}
```

### Index Price *(v3)*

The index price is an average across multiple exchanges — useful for understanding "fair value":

```go
indexPrice, err := client.Market().GetIndexPrice("BTC-USDT")
```

### Ticker Price *(v3)*

```go
symbol := "BTC-USDT"
tickerPrice, err := client.Market().GetTickerPrice(&symbol)
```

---

## Order Book Depth

The order book shows you all the buy and sell orders waiting to be filled. Essential for understanding liquidity and potential slippage.

### Full Order Book

```go
// Grab the top 20 price levels
depth, err := client.Market().GetDepth("BTC-USDT", 20)

// Bids = people wanting to buy (sorted high to low)
fmt.Println("Buyers:")
for _, bid := range depth["bids"].([]interface{}) {
    b := bid.([]interface{})
    fmt.Printf("  %v BTC wanted @ %v\n", b[1], b[0])
}

// Asks = people wanting to sell (sorted low to high)
fmt.Println("Sellers:")
for _, ask := range depth["asks"].([]interface{}) {
    a := ask.([]interface{})
    fmt.Printf("  %v BTC offered @ %v\n", a[1], a[0])
}
```

### Just the Top of Book *(v3)*

If you only need the best bid and ask (faster than fetching the whole book):

```go
symbol := "BTC-USDT"
bookTicker, err := client.Market().GetBookTicker(&symbol)

fmt.Printf("Best bid: %v @ %v\n", bookTicker["bidQty"], bookTicker["bidPrice"])
fmt.Printf("Best ask: %v @ %v\n", bookTicker["askQty"], bookTicker["askPrice"])

// The spread tells you about liquidity
spread := bookTicker["askPrice"].(float64) - bookTicker["bidPrice"].(float64)
fmt.Printf("Spread: %.2f\n", spread)
```

### Spot Book Ticker *(v3)*

```go
spotBookTicker, err := client.Market().GetSpotBookTicker(&symbol)
```

---

## Candlestick Charts (Klines)

Klines (candlesticks) are the foundation of technical analysis. Each candle shows the open, high, low, close, and volume for a time period.

### Basic Usage

```go
// Get the last 100 hourly candles
klines, err := client.Market().GetKlines(
    "BTC-USDT",
    "1h",       // timeframe
    100,        // how many candles
    nil,        // startTime
    nil,        // endTime
)

for _, kline := range klines {
    k := kline.(map[string]interface{})
    fmt.Printf("Time: %v, Open: %v, High: %v, Low: %v, Close: %v, Volume: %v\n",
        k["time"], k["open"], k["high"], k["low"], k["close"], k["volume"])
}
```

### Available Timeframes

| Short-term | Medium-term | Long-term |
|------------|-------------|----------|
| `1m`, `3m`, `5m` | `1h`, `2h`, `4h` | `1d`, `3d` |
| `15m`, `30m` | `6h`, `12h` | `1w`, `1M` |

### Fetching a Specific Time Range

```go
import "time"

// Get klines for last 24 hours
endTime := time.Now().UnixMilli()
startTime := time.Now().Add(-24 * time.Hour).UnixMilli()

klines, err := client.Market().GetKlines(
    "BTC-USDT",
    "1h",
    100,
    &startTime,
    &endTime,
)
```

### Spot Klines with Timezone

Spot klines can be adjusted for your timezone:

```go
// Get candles in UTC+8 (e.g., Singapore/Hong Kong time)
timeZone := int64(8)
spotKlines, err := client.Market().GetSpotKlines(
    "BTC-USDT",
    "1h",
    100,
    nil,
    nil,
    &timeZone,
)
```

---

## 24h Statistics

Want to know how a coin performed today? The 24-hour ticker has everything.

### Single Symbol

```go
symbol := "BTC-USDT"
ticker, err := client.Market().Get24hrTicker(&symbol)

fmt.Printf("BTC in the last 24 hours:\n")
fmt.Printf("  Change: %v (%v%%)\n", ticker["priceChange"], ticker["priceChangePercent"])
fmt.Printf("  High: %v\n", ticker["highPrice"])
fmt.Printf("  Low: %v\n", ticker["lowPrice"])
fmt.Printf("  Volume: %v\n", ticker["volume"])
```

### Scan the Whole Market

```go
// Pass nil to get everything
tickers, err := client.Market().Get24hrTicker(nil)

// Find the biggest movers
for _, ticker := range tickers {
    t := ticker.(map[string]interface{})
    change := t["priceChangePercent"].(float64)
    
    if change > 10 || change < -10 {
        fmt.Printf("🚨 %s moved %.1f%% today!\n", t["symbol"], change)
    }
}
```

---

## Funding Rates

Funding rates are how perpetual futures stay pegged to spot price. Every 8 hours, longs pay shorts (or vice versa) based on this rate. If you're holding positions overnight, you need to watch this.

### Historical Rates

```go
fundingRates, err := client.Market().GetFundingRateHistory("BTC-USDT", 100)

for _, rate := range fundingRates {
    r := rate.(map[string]interface{})
    fmt.Printf("%v: %v\n", r["fundingTime"], r["fundingRate"])
}
```

### Current Rate & Next Payment *(v3)*

```go
fundingInfo, err := client.Market().GetFundingRateInfo("BTC-USDT")

rate := fundingInfo["fundingRate"].(float64)
fmt.Printf("Current rate: %.4f%%\n", rate*100)
fmt.Printf("Next payment: %v\n", fundingInfo["fundingTime"])

// Positive = longs pay shorts (bullish sentiment)
// Negative = shorts pay longs (bearish sentiment)
if rate > 0.001 {
    fmt.Println("🟢 High positive funding - consider shorting")
} else if rate < -0.001 {
    fmt.Println("🔴 High negative funding - consider longing")
}
```

---

## Open Interest *(v3)*

Open interest = total value of all open positions. It tells you how much money is "in the game."

**Reading OI:**
- OI rising + price rising = new longs entering (bullish)
- OI rising + price falling = new shorts entering (bearish)
- OI falling = positions closing (trend weakening)

### Current Open Interest

```go
oi, err := client.Market().GetOpenInterest("BTC-USDT")
fmt.Printf("Total open positions: %v\n", oi["openInterest"])
```

### Historical Open Interest

```go
oiHistory, err := client.Market().GetOpenInterestHistory(
    "BTC-USDT",
    "1h",       // granularity: 5m, 15m, 30m, 1h, 4h, 1d
    100,        // how many data points
    nil, nil,   // time range (optional)
)
```

### Practical Example: OI Change Alert

```go
func checkOIChange(client *bingx.Client, symbol string) {
    currentOI, _ := client.Market().GetOpenInterest(symbol)
    oiHistory, _ := client.Market().GetOpenInterestHistory(symbol, "1h", 24, nil, nil)
    
    if len(oiHistory) == 0 {
        return
    }
    
    // Compare current to 24h ago
    oldOI := oiHistory[0].(map[string]interface{})["openInterest"].(float64)
    newOI := currentOI["openInterest"].(float64)
    change := ((newOI - oldOI) / oldOI) * 100
    
    fmt.Printf("OI changed %.1f%% in 24h\n", change)
    
    if change > 15 {
        fmt.Println("🚨 Big OI spike - expect volatility!")
    } else if change < -15 {
        fmt.Println("🚨 Big OI drop - positions unwinding")
    }
}
```

---

## More Market Data *(v3)*

### Recent Trades

See the actual trades happening in real-time:

```go
trades, err := client.Market().GetRecentTrades("BTC-USDT", 100)

for _, trade := range trades {
    t := trade.(map[string]interface{})
    fmt.Printf("%v BTC @ %v\n", t["qty"], t["price"])
}
```

### Aggregate Trades

```go
aggTrades, err := client.Market().GetAggregateTrades(
    "BTC-USDT",
    100,
    nil,  // fromId
    nil,  // startTime
    nil,  // endTime
)
```

### Long/Short Ratio

See how traders are positioned. When everyone's long, maybe it's time to be careful:

```go
ratio, err := client.Market().GetLongShortRatio("BTC-USDT", "5m", 100)

for _, r := range ratio {
    record := r.(map[string]interface{})
    longPct := record["longAccount"].(float64)
    shortPct := record["shortAccount"].(float64)
    
    fmt.Printf("Longs: %.1f%% | Shorts: %.1f%%\n", longPct, shortPct)
    
    // Extreme readings can signal reversals
    if longPct > 70 {
        fmt.Println("⚠️ Crowded long - watch for squeeze")
    }
}
```

### Basis Data

```go
basis, err := client.Market().GetBasis(
    "BTC-USDT",
    "CURRENT_QUARTER",  // contractType
    100,                // limit
    nil,                // startTime
    nil,                // endTime
)
```

---

## Real-World Examples

Here are some patterns you'll actually use.

### Price Alert Bot

```go
type PriceAlert struct {
    Symbol      string
    TargetPrice float64
    Triggered   bool
}

func monitorPrice(client *bingx.Client, alert *PriceAlert) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        price, err := client.Market().GetLatestPrice(alert.Symbol)
        if err != nil {
            continue // Keep trying
        }
        
        currentPrice := price.(float64)
        
        if currentPrice >= alert.TargetPrice && !alert.Triggered {
            fmt.Printf("🔔 %s hit %.2f!\n", alert.Symbol, currentPrice)
            alert.Triggered = true
            // Send to Telegram, Slack, etc.
            return
        }
    }
}
```

### Find Today's Biggest Movers

```go
func findBigMovers(client *bingx.Client) {
    tickers, _ := client.Market().Get24hrTicker(nil)
    
    type Mover struct {
        Symbol string
        Change float64
    }
    
    var gainers, losers []Mover
    
    for _, ticker := range tickers {
        t := ticker.(map[string]interface{})
        change := t["priceChangePercent"].(float64)
        symbol := t["symbol"].(string)
        
        if change > 5 {
            gainers = append(gainers, Mover{symbol, change})
        } else if change < -5 {
            losers = append(losers, Mover{symbol, change})
        }
    }
    
    fmt.Println("🟢 Top Gainers:")
    for _, g := range gainers[:min(5, len(gainers))] {
        fmt.Printf("  %s: +%.1f%%\n", g.Symbol, g.Change)
    }
    
    fmt.Println("🔴 Top Losers:")
    for _, l := range losers[:min(5, len(losers))] {
        fmt.Printf("  %s: %.1f%%\n", l.Symbol, l.Change)
    }
}
```

### Volume Spike Detector

High volume often precedes big moves:

```go
func checkVolumeSpike(client *bingx.Client, symbol string) bool {
    klines, _ := client.Market().GetKlines(symbol, "1h", 24, nil, nil)
    
    // Calculate average volume
    var totalVol float64
    for _, k := range klines[:len(klines)-1] {
        totalVol += k.(map[string]interface{})["volume"].(float64)
    }
    avgVol := totalVol / float64(len(klines)-1)
    
    // Compare to current
    currentVol := klines[len(klines)-1].(map[string]interface{})["volume"].(float64)
    
    if currentVol > avgVol*2 {
        fmt.Printf("🚨 %s volume is %.1fx average!\n", symbol, currentVol/avgVol)
        return true
    }
    return false
}
```

### Quick Multi-Timeframe Check

```go
func quickMTF(client *bingx.Client, symbol string) {
    timeframes := []string{"15m", "1h", "4h", "1d"}
    
    fmt.Printf("%s trend check:\n", symbol)
    
    for _, tf := range timeframes {
        klines, _ := client.Market().GetKlines(symbol, tf, 2, nil, nil)
        if len(klines) < 2 {
            continue
        }
        
        prev := klines[0].(map[string]interface{})["close"].(float64)
        curr := klines[1].(map[string]interface{})["close"].(float64)
        
        arrow := "🟢"
        if curr < prev {
            arrow = "🔴"
        }
        
        fmt.Printf("  %s %s\n", tf, arrow)
    }
}
```

---

## Pro Tips

### Don't Hammer the API

Add a simple rate limiter to avoid getting blocked:

```go
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(10, 1) // 10 req/sec

func getPrice(client *bingx.Client, symbol string) float64 {
    limiter.Wait(context.Background())
    price, _ := client.Market().GetLatestPrice(symbol)
    return price.(float64)
}
```

### Cache When You Can

If you're checking the same price multiple times per second, cache it:

```go
var (
    priceCache = make(map[string]float64)
    cacheTime  = make(map[string]time.Time)
    cacheMu    sync.RWMutex
)

func getCachedPrice(client *bingx.Client, symbol string) float64 {
    cacheMu.RLock()
    if t, ok := cacheTime[symbol]; ok && time.Since(t) < time.Second {
        price := priceCache[symbol]
        cacheMu.RUnlock()
        return price
    }
    cacheMu.RUnlock()
    
    // Fetch fresh
    price, _ := client.Market().GetLatestPrice(symbol)
    p := price.(float64)
    
    cacheMu.Lock()
    priceCache[symbol] = p
    cacheTime[symbol] = time.Now()
    cacheMu.Unlock()
    
    return p
}
```

### Retry on Failure

Network hiccups happen. Don't let them crash your bot:

```go
func getWithRetry(client *bingx.Client, symbol string) (float64, error) {
    for attempt := 1; attempt <= 3; attempt++ {
        price, err := client.Market().GetLatestPrice(symbol)
        if err == nil {
            return price.(float64), nil
        }
        time.Sleep(time.Duration(attempt) * time.Second)
    }
    return 0, fmt.Errorf("gave up after 3 tries")
}
```

---

## Keep Going

- **[Trading](Trade-Service)** — Put this data to use
- **[Account](Account-Service)** — Check your balance and positions
- **[WebSockets](WebSocket-Service)** — Get data pushed to you in real-time
- **[v3 Features](API-v3-Features)** — The new stuff

---

**Questions?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues) — we're here to help.
