# Welcome to bingx-go!

**Your Go-to SDK for trading on BingX — now with full API v3 support**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![API v3](https://img.shields.io/badge/API-v3-blue?style=flat-square)](https://bingx-api.github.io/docs-v3/)

Whether you're building a trading bot, analyzing market data, or managing your portfolio programmatically, this SDK has you covered. It wraps the entire BingX API in clean, idiomatic Go code that's easy to use and maintain.

## 🚀 Fresh in Version 2.0.5

We've just shipped a major update with over 26 new methods. Here's what you can do now:

- **TWAP Orders** — Break up large trades into smaller pieces over time, so you don't move the market against yourself
- **Multi-Assets Margin** — Use your entire portfolio as collateral instead of managing margin per-position
- **One-Click Reversal** — Flip from LONG to SHORT (or back) with a single call
- **Trailing Stops** — Set dynamic stop-losses that follow the price as it moves in your favor
- **Better Risk Tools** — Track your P&L, commissions, and liquidation risk in real-time
- **More Market Data** — Open interest history, funding rates, best bid/ask, and more

Already using an older version? Don't worry — everything is backward compatible. Check out the **[Migration Guide](https://github.com/tigusigalpa/bingx-go/blob/main/API_V3_MIGRATION.md)** when you're ready to explore the new features.

---

## 📚 Find What You Need

### Just Getting Started?
- **[Installation & Setup](Getting-Started)** — Get up and running in a few minutes
- **[Quick Start Guide](Getting-Started#quick-start)** — Build your first trading bot step by step
- **[What's New in v3](API-v3-Features)** — Explore all the shiny new features

### Core Functionality
- **[Market Data](Market-Service)** — Prices, order books, candles, funding rates
- **[Account Management](Account-Service)** — Check balances, positions, leverage settings
- **[Trading](Trade-Service)** — Place orders, manage positions, use TWAP
- **[Wallet Operations](Wallet-Service)** — Handle deposits, withdrawals, transfers

### Going Deeper
- **[WebSocket Streams](WebSocket-Service)** — Get real-time updates pushed to you
- **[Handling Errors](Error-Handling)** — Deal with problems gracefully
- **[Sub-Accounts](Sub-Account-Service)** — Manage multiple accounts under one master
- **[Copy Trading](Copy-Trading-Service)** — Follow other traders automatically
- **[Coin-M Futures](Coin-M-Futures)** — Trade coin-margined contracts
- **[Pro Tips](Advanced-Features)** — Rate limiting, batching, and testing strategies

---

## 🎯 Jump to What You Need

### "I want to..."

| Goal | Where to Look |
|------|---------------|
| Fetch prices and charts | [Market Data](Market-Service) |
| Place my first order | [Trading Basics](Trade-Service#basic-orders) |
| See how much money I have | [Account Service](Account-Service) |
| Execute a huge order without slippage | [TWAP Orders](API-v3-Features#twap-orders) |
| Monitor my risk exposure | [Position Risk](API-v3-Features#position-risk-management) |
| Get live price updates | [WebSocket Streams](WebSocket-Service) |
| Figure out why something broke | [Error Handling](Error-Handling) |

### Pick Your Path

**New to the SDK?** Start here:
1. [Install and configure](Getting-Started) the library
2. Follow the [Quick Start](Getting-Started#quick-start) tutorial
3. Learn to [fetch market data](Market-Service#basic-usage)
4. [Place your first order](Trade-Service#basic-orders)

**Ready for more?** Level up:
1. Master [account management](Account-Service)
2. Try [advanced order types](Trade-Service#advanced-order-types) like trailing stops
3. Set up [WebSocket streams](WebSocket-Service) for real-time data
4. Implement proper [error handling](Error-Handling)

**Building something serious?** Go pro:
1. Use [TWAP](API-v3-Features#twap-orders) for institutional-grade execution
2. Enable [Multi-Assets Mode](API-v3-Features#multi-assets-margin) for better capital efficiency
3. Build a [risk dashboard](API-v3-Features#position-risk-management)
4. Optimize with [batch operations](Advanced-Features#batch-operations)

---

## 📊 What's Supported Where?

Not every feature works on every market type. Here's the breakdown:

| Feature | USDT-M Futures | Coin-M Futures | Spot |
|---------|----------------|----------------|------|
| Market Data | ✅ Full | ✅ Full | ✅ Full |
| Trading | ✅ Full | ✅ Full | ⚠️ Basic |
| Account Info | ✅ Full | ✅ Full | ✅ Full |
| WebSockets | ✅ Full | ✅ Full | ✅ Full |
| TWAP Orders | ✅ New in v3 | ❌ Not available | ❌ Not available |
| Multi-Assets | ✅ New in v3 | ❌ Not available | ❌ Not available |
| Risk Metrics | ✅ New in v3 | ✅ New in v3 | — |

---

## 💡 See It in Action

### The Basics

Here's how simple it is to get started:

```go
import bingx "github.com/tigusigalpa/bingx-go/v2"

// Create a client with your API credentials
client := bingx.NewClient("API_KEY", "API_SECRET")

// What's BTC trading at right now?
price, _ := client.Market().GetLatestPrice("BTC-USDT")

// How much do I have in my account?
balance, _ := client.Account().GetBalance()

// Let's buy some Bitcoin
order, _ := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "side":     "BUY",
    "type":     "LIMIT",
    "price":    50000.0,
    "quantity": 0.001,
})
```

### The Cool New Stuff (v3)

Here's where it gets interesting:

```go
import "github.com/tigusigalpa/bingx-go/v2/services"

// Got a big order? Spread it out over an hour to avoid slippage
client.Trade().PlaceTWAPOrder(map[string]interface{}{
    "symbol":   "BTC-USDT",
    "quantity": 10.0,
    "duration": 3600,  // 1 hour
    "interval": 60,    // execute a piece every minute
})

// Market turned against you? Flip your position instantly
client.Trade().OneClickReversePosition("BTC-USDT", nil)

// Keep an eye on your risk
risk, _ := client.Account().GetPositionRisk(&symbol, nil)

// Use your whole portfolio as margin
client.Trade().SwitchMultiAssetsMode(true, nil)
```

---

## 🔗 Useful Links

- **[Source Code](https://github.com/tigusigalpa/bingx-go)** — Star us if you find this useful!
- **[Official BingX API Docs](https://bingx-api.github.io/docs-v3/)** — The source of truth
- **[GoDoc](https://pkg.go.dev/github.com/tigusigalpa/bingx-go/v2)** — Auto-generated API reference
- **[Report a Bug](https://github.com/tigusigalpa/bingx-go/issues)** — Found something broken? Let us know
- **[Changelog](https://github.com/tigusigalpa/bingx-go/blob/main/CHANGELOG.md)** — See what changed in each version

---

## 🤝 Want to Contribute?

We'd love your help! Whether it's fixing a typo, adding a feature, or improving docs — all contributions are welcome. Check out our [Contributing Guide](https://github.com/tigusigalpa/bingx-go/blob/main/CONTRIBUTING.md) to get started.

---

## 📄 License

This project is MIT licensed, which means you can use it for pretty much anything. See the [LICENSE](https://github.com/tigusigalpa/bingx-go/blob/main/LICENSE) file for the legal details.

---

*Last updated: April 2026 • Version 1.4.0*
