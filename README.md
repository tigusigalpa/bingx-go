# BingX Go SDK

<div align="center">

![BingX Golang SDK](https://github.com/user-attachments/assets/61babf48-746f-469e-aeed-5058a9337d86)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/bingx-go?style=flat-square)](https://goreportcard.com/report/github.com/tigusigalpa/bingx-go)
[![GitHub Release](https://img.shields.io/github/v/release/tigusigalpa/bingx-go?style=flat-square)](https://github.com/tigusigalpa/bingx-go/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/tigusigalpa/bingx-go)

</div>

Go client for the [BingX](https://bingx.com) cryptocurrency exchange API. Covers USDT-M and Coin-M perpetual futures, spot trading, copy trading, sub-accounts, and WebSocket streaming. 220+ API methods.

**Package:** [pkg.go.dev/github.com/tigusigalpa/bingx-go](https://pkg.go.dev/github.com/tigusigalpa/bingx-go)

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
    - [Market Data](#market-service---market-data)
    - [Account Management](#account-service---account-management)
    - [Trading Operations](#trade-service---trading-operations)
    - [Wallet Management](#wallet-service---wallet-management)
    - [Coin-M Futures](#coin-m-perpetual-futures)
    - [Sub-Accounts](#sub-account-management)
    - [Copy Trading](#copy-trading-operations)
    - [WebSocket Streaming](#websocket-streaming)
- [Advanced Features](#advanced-features)
- [Configuration](#configuration)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## Features

### Service coverage

<table>
<tr>
<td width="50%" valign="top">

#### USDT-M Perpetual Futures

- **Market Service** (40+ methods)
    - Real-time & historical prices
    - Market depth & order books
    - Candlestick data (K-lines)
    - 24hr tickers & statistics
    - Funding rates & premium index
    - Aggregate & recent trades
    - Long/short ratios
    - Basis data & sentiment analysis

- **Account Service** (20+ methods)
    - Balance & asset management
    - Position tracking & monitoring
    - Leverage configuration
    - Margin mode management
    - Trading fees & commissions
    - API permissions & rate limits
    - Balance history & deposits

- **Trade Service** (25+ methods)
    - Order creation & management
    - Batch order operations
    - Order modification & cancellation
    - Position management
    - Trade history & analytics
    - Test orders (sandbox)
    - Commission calculations

- **Wallet Service** (6+ methods)
    - Deposit & withdrawal management
    - Address generation
    - Transaction history
    - Multi-coin support
    - Network selection

</td>
<td width="50%" valign="top">

#### **Advanced Features**

- **Spot Account Service** (8+ methods)
    - Spot balance management
    - Universal transfers
    - Internal transfers
    - Transfer history
    - Multi-account support

- **Sub-Account Service** (20+ methods)
    - Sub-account creation & management
    - API key management
    - Asset transfers
    - Deposit address management
    - Authorization controls

- **Copy Trading Service** (13+ methods)
    - Futures copy trading
    - Spot copy trading
    - Profit tracking
    - Commission management
    - Trading pair configuration

- **Contract Service** (3+ methods)
    - Standard contract positions
    - Order history
    - Balance queries

- **Listen Key Service** (3+ methods)
    - WebSocket authentication
    - Key generation & extension
    - Session management

#### **Coin-M Perpetual Futures**

- **Coin-M Market** (6+ methods)
    - Contract specifications
    - Ticker & price data
    - Market depth
    - K-line data
    - Open interest
    - Funding rates

- **Coin-M Trade** (17+ methods)
    - Order management
    - Position tracking
    - Leverage & margin
    - Balance queries
    - Trade history

</td>
</tr>
</table>

### Key capabilities

| Feature              | Description                                                |
|----------------------|------------------------------------------------------------|
| **Real-Time Data**   | WebSocket support for live market data and account updates |
| **Order Types**      | MARKET, LIMIT, STOP, STOP_MARKET, TAKE_PROFIT, and more    |
| **Position Modes**   | One-way and Hedge mode support                             |
| **Margin Types**     | ISOLATED and CROSSED margin modes                          |
| **Time in Force**    | GTC, IOC, FOK order execution options                      |
| **Risk Management**  | Stop-loss, take-profit, trailing stops                     |
| **Batch Operations** | Create, modify, and cancel multiple orders at once         |
| **Historical Data**  | Access to comprehensive historical trading data            |

---

## Architecture

```
bingx-go/
├── client.go              # Main client with service orchestration
├── coinm_client.go        # Coin-M futures client
├── http/
│   └── client.go          # HTTP client with HMAC signing
├── errors/
│   └── errors.go          # Custom error types
├── services/
│   ├── market.go          # Market data service
│   ├── account.go         # Account management service
│   ├── trade.go           # Trading operations service
│   ├── wallet.go          # Wallet operations service
│   ├── spotaccount.go     # Spot account service
│   ├── subaccount.go      # Sub-account service
│   ├── copytrading.go     # Copy trading service
│   ├── contract.go        # Contract service
│   ├── listenkey.go       # WebSocket auth service
│   └── coinm/             # Coin-M specific services
│       ├── market.go
│       ├── trade.go
│       └── listenkey.go
└── examples/              # Usage examples
```

### Design Principles

- **Single Responsibility** - Each service handles a specific domain
- **Dependency Injection** - HTTP client injected for testability
- **Functional Options** - Flexible client configuration
- **Error Transparency** - Rich error types for debugging
- **Immutability** - Thread-safe by design

---

## Installation

- Go 1.21+
- BingX account with API keys (for authenticated endpoints)

```bash
go get github.com/tigusigalpa/bingx-go
```

Only external dependency: `github.com/gorilla/websocket` (for WebSocket features)

---

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    bingx "github.com/tigusigalpa/bingx-go"
)

func main() {
    // Create client
    client := bingx.NewClient(
        "YOUR_API_KEY",
        "YOUR_API_SECRET",
        bingx.WithBaseURI("https://open-api.bingx.com"),
        bingx.WithSignatureEncoding("base64"),
    )
    
    // Get current price
    price, err := client.Market().GetLatestPrice("BTC-USDT")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC price: %v\n", price)
    
    // Get account balance
    balance, err := client.Account().GetBalance()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %v\n", balance)
    
    // Create an order
    order, err := client.Trade().CreateOrder(map[string]interface{}{
        "symbol":       "BTC-USDT",
        "side":         "BUY",
        "type":         "LIMIT",
        "positionSide": "LONG",
        "price":        50000.0,
        "quantity":     0.001,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order created: %v\n", order)
}
```

---

## Usage Examples

### Market Service - Market Data

```go
// Get futures symbols
symbols, err := client.Market().GetFuturesSymbols()

// Get spot symbols
spotSymbols, err := client.Market().GetSpotSymbols()

// Get all symbols
allSymbols, err := client.Market().GetAllSymbols()

// Get latest price
price, err := client.Market().GetLatestPrice("BTC-USDT")

// Get market depth
depth, err := client.Market().GetDepth("BTC-USDT", 20)

// Get K-lines
klines, err := client.Market().GetKlines("BTC-USDT", "1h", 100, nil, nil)

// Get 24hr ticker
ticker, err := client.Market().Get24hrTicker(nil) // nil for all symbols

// Get funding rate history
fundingRate, err := client.Market().GetFundingRateHistory("BTC-USDT", 100)
```

### Account Service - Account Management

```go
// Get account balance
balance, err := client.Account().GetBalance()

// Get positions
positions, err := client.Account().GetPositions(nil) // nil for all positions

// Get positions for specific symbol
btcPositions, err := client.Account().GetPositions(&symbol)

// Get account info
accountInfo, err := client.Account().GetAccountInfo()

// Get/Set leverage
leverage, err := client.Account().GetLeverage("BTC-USDT", nil)
err = client.Account().SetLeverage("BTC-USDT", "BOTH", 20, nil)

// Get/Set margin mode
marginMode, err := client.Account().GetMarginMode("BTC-USDT")
err = client.Account().SetMarginMode("BTC-USDT", "ISOLATED")
```

### Trade Service - Trading Operations

```go
// Create order
order, err := client.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         "MARKET",
    "positionSide": "LONG",
    "quantity":     0.001,
})

// Cancel order
orderID := "123456789"
err = client.Trade().CancelOrder("BTC-USDT", &orderID, nil, nil, nil)

// Cancel all orders
err = client.Trade().CancelAllOrders(nil, &symbol, nil, nil)

// Get order details
order, err := client.Trade().GetOrder("BTC-USDT", "123456789")

// Get open orders
openOrders, err := client.Trade().GetOpenOrders(nil, 100)

// Get order history
history, err := client.Trade().GetOrderHistory(&symbol, 100, nil, nil)

// Get user trades
trades, err := client.Trade().GetUserTrades(&symbol, 100, nil, nil)
```

### Wallet Service - Wallet Management

```go
// Get deposit history
deposits, err := client.Wallet().GetDepositHistory("USDT", nil, nil, nil, 100)

// Get deposit address
address, err := client.Wallet().GetDepositAddress("USDT", "TRC20")

// Get withdrawal history
withdrawals, err := client.Wallet().GetWithdrawalHistory("USDT", nil, nil, nil, 100)

// Create withdrawal
withdrawal, err := client.Wallet().Withdraw("USDT", "TXxx...xxx", 100.0, "TRC20", nil)

// Get all coin info
coins, err := client.Wallet().GetAllCoinInfo()
```

### Coin-M Perpetual Futures

```go
// Get Coin-M market data
ticker, err := client.CoinM().Market().GetTicker("BTC-USD")

// Get Coin-M contracts
contracts, err := client.CoinM().Market().GetContracts()

// Create Coin-M order
order, err := client.CoinM().Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USD",
    "side":         "BUY",
    "positionSide": "LONG",
    "type":         "MARKET",
    "quantity":     100,
})

// Get Coin-M positions
positions, err := client.CoinM().Trade().GetPositions(nil)

// Get Coin-M balance
balance, err := client.CoinM().Trade().GetBalance()
```

### Sub-Account Management

```go
// Create sub-account
result, err := client.SubAccount().CreateSubAccount("sub_account_001")

// Get sub-account list
subAccounts, err := client.SubAccount().GetSubAccountList(nil, 1, 10)

// Create API key for sub-account
apiKey, err := client.SubAccount().CreateSubAccountAPIKey(
    "sub_account_001",
    "Trading Bot",
    map[string]bool{"spot": true, "futures": true},
    nil, // IP whitelist (optional)
)

// Transfer assets to sub-account
transfer, err := client.SubAccount().SubAccountInternalTransfer(
    "USDT",
    "SPOT",
    100.0,
    "FROM_MAIN_TO_SUB",
    nil,
    &subUID,
    nil,
)

// Get sub-account assets
assets, err := client.SubAccount().GetSubAccountAssets("12345678")
```

### Copy Trading Operations

```go
// Get current track orders
orders, err := client.CopyTrading().GetCurrentTrackOrders("BTC-USDT")

// Close track order
result, err := client.CopyTrading().CloseTrackOrder("1252864099381234567")

// Set take profit and stop loss
stopLoss := 48000.0
takeProfit := 52000.0
err = client.CopyTrading().SetTPSL("1252864099381234567", &stopLoss, &takeProfit)

// Get profit summary
summary, err := client.CopyTrading().GetProfitSummary()

// Set commission rate
err = client.CopyTrading().SetCommission(5.0) // 5% commission
```

### WebSocket Streaming

#### Market Data Stream

```go
// Create market data stream
stream := client.NewMarketDataStream()

// Connect to WebSocket
if err := stream.Connect(); err != nil {
    log.Fatal(err)
}
defer stream.Disconnect()

// Register message handler
stream.OnMessage(func(data map[string]interface{}) {
    fmt.Printf("Received: %+v\n", data)
})

// Subscribe to various data types
stream.SubscribeTrade("BTC-USDT")
stream.SubscribeKline("BTC-USDT", "1m")
stream.SubscribeDepth("BTC-USDT", 20)
stream.SubscribeTicker("BTC-USDT")
stream.SubscribeBookTicker("BTC-USDT")

// Start listening (blocking)
stream.Listen()
```

#### Account Data Stream

```go
// Generate listen key
resp, err := client.ListenKey().Generate()
if err != nil {
    log.Fatal(err)
}
listenKey := resp["listenKey"].(string)

// Create account data stream
stream := client.NewAccountDataStream(listenKey)

if err := stream.Connect(); err != nil {
    log.Fatal(err)
}
defer stream.Disconnect()

// Listen for account updates
stream.OnAccountUpdate(func(eventType string, data map[string]interface{}) {
    fmt.Printf("Event [%s]: %+v\n", eventType, data)
})

// Listen for specific update types
stream.OnBalanceUpdate(func(balances interface{}) {
    fmt.Printf("Balance: %+v\n", balances)
})

stream.OnPositionUpdate(func(positions interface{}) {
    fmt.Printf("Positions: %+v\n", positions)
})

stream.OnOrderUpdate(func(order interface{}) {
    fmt.Printf("Order: %+v\n", order)
})

// Start listening
stream.Listen()
```

For detailed WebSocket documentation, see [websocket/README.md](websocket/README.md).

---

## Advanced Features

### Commission Calculation

```go
// Calculate futures commission
commission := client.Trade().CalculateFuturesCommission(100.0, 10, nil)
fmt.Printf("Commission: %.6f USDT\n", commission.Commission)
fmt.Printf("Net Position Value: %.2f USDT\n", commission.NetPositionValue)

// Quick commission amount
amount := client.Trade().GetCommissionAmount(100.0, 10)
fmt.Printf("Commission Amount: %.6f USDT\n", amount)
```

### Batch Operations

```go
// Create multiple orders at once
orders := []map[string]interface{}{
    {
        "symbol":       "BTC-USDT",
        "side":         "BUY",
        "type":         "LIMIT",
        "positionSide": "LONG",
        "price":        50000.0,
        "quantity":     0.001,
    },
    {
        "symbol":       "ETH-USDT",
        "side":         "BUY",
        "type":         "LIMIT",
        "positionSide": "LONG",
        "price":        3000.0,
        "quantity":     0.01,
    },
}

result, err := client.Trade().CreateBatchOrders(orders)

// Cancel multiple orders
orderIDs := []string{"123456789", "987654321"}
err = client.Trade().CancelBatchOrders("BTC-USDT", orderIDs, nil, nil, nil)
```

### Time-Based Queries

```go
import "time"

// Get trades from last 24 hours
startTime := time.Now().Add(-24 * time.Hour).UnixMilli()
endTime := time.Now().UnixMilli()

trades, err := client.Trade().GetUserTrades(
    &symbol,
    100,
    &startTime,
    &endTime,
)

// Get historical klines
klines, err := client.Market().GetKlines(
    "BTC-USDT",
    "1h",
    100,
    &startTime,
    &endTime,
)
```

---

## Configuration

```go
client := bingx.NewClient(
    apiKey,
    apiSecret,
    // Set custom base URI (default: https://open-api.bingx.com)
    bingx.WithBaseURI("https://open-api.bingx.com"),
    
    // Set source key for tracking (optional)
    bingx.WithSourceKey("my-trading-bot-v1"),
    
    // Set signature encoding: "base64" or "hex" (default: base64)
    bingx.WithSignatureEncoding("base64"),
)
```

### Environment Variables

```bash
# Recommended: Store credentials in environment variables
export BINGX_API_KEY="your_api_key_here"
export BINGX_API_SECRET="your_api_secret_here"
export BINGX_SOURCE_KEY="optional_source_key"
```

```go
import "os"

client := bingx.NewClient(
    os.Getenv("BINGX_API_KEY"),
    os.Getenv("BINGX_API_SECRET"),
    bingx.WithSourceKey(os.Getenv("BINGX_SOURCE_KEY")),
)
```

### Timeout Configuration

```go
// The HTTP client uses default timeouts:
// - Connection timeout: 10 seconds
// - Request timeout: 30 seconds

// For custom timeout, modify the HTTP client in http/client.go
```

---

## API Reference

### Client Methods

| Method                                            | Description                       | Returns                |
|---------------------------------------------------|-----------------------------------|------------------------|
| `NewClient(apiKey, apiSecret string, options...)` | Create new BingX client           | `*Client`              |
| `client.Market()`                                 | Access market data service        | `*MarketService`       |
| `client.Account()`                                | Access account management service | `*AccountService`      |
| `client.Trade()`                                  | Access trading operations service | `*TradeService`        |
| `client.Contract()`                               | Access contract service           | `*ContractService`     |
| `client.ListenKey()`                              | Access WebSocket auth service     | `*ListenKeyService`    |
| `client.Wallet()`                                 | Access wallet operations service  | `*WalletService`       |
| `client.SpotAccount()`                            | Access spot account service       | `*SpotAccountService`  |
| `client.SubAccount()`                             | Access sub-account service        | `*SubAccountService`   |
| `client.CopyTrading()`                            | Access copy trading service       | `*CopyTradingService`  |
| `client.CoinM()`                                  | Access Coin-M futures client      | `*CoinMClient`         |
| `client.NewMarketDataStream()`                    | Create market data WebSocket      | `*MarketDataStream`    |
| `client.NewAccountDataStream(listenKey)`          | Create account data WebSocket     | `*AccountDataStream`   |
| `client.GetHTTPClient()`                          | Get underlying HTTP client        | `*http.BaseHTTPClient` |
| `client.GetEndpoint()`                            | Get API endpoint URL              | `string`               |
| `client.GetAPIKey()`                              | Get configured API key            | `string`               |

### Service Overview

<details>
<summary><b>Market Service Methods</b></summary>

- `GetFuturesSymbols()` - Get all futures trading symbols
- `GetSpotSymbols()` - Get all spot trading symbols
- `GetAllSymbols()` - Get both spot and futures symbols
- `GetLatestPrice(symbol)` - Get latest price for a symbol
- `GetSpotLatestPrice(symbol)` - Get latest spot price
- `GetDepth(symbol, limit)` - Get order book depth
- `GetSpotDepth(symbol, limit)` - Get spot order book
- `GetKlines(symbol, interval, limit, startTime, endTime)` - Get candlestick data
- `GetSpotKlines(...)` - Get spot candlestick data
- `Get24hrTicker(symbol)` - Get 24hr ticker statistics
- `GetSpot24hrTicker(symbol)` - Get spot 24hr ticker
- `GetFundingRateHistory(symbol, limit)` - Get funding rate history
- `GetMarkPrice(symbol)` - Get mark price
- `GetPremiumIndexKlines(...)` - Get premium index klines
- `GetAggregateTrades(...)` - Get aggregate trades
- `GetRecentTrades(symbol, limit)` - Get recent trades
- `GetServerTime()` - Get server time
- `GetContinuousKlines(...)` - Get continuous contract klines
- `GetIndexPriceKlines(...)` - Get index price klines
- `GetTopLongShortRatio(...)` - Get long/short ratio
- `GetHistoricalTopLongShortRatio(...)` - Get historical ratios
- `GetBasis(...)` - Get basis data

</details>

<details>
<summary><b>Account Service Methods</b></summary>

- `GetBalance()` - Get account balance
- `GetPositions(symbol)` - Get positions
- `GetAccountInfo()` - Get account information
- `GetTradingFees(symbol)` - Get trading fees
- `GetMarginMode(symbol)` - Get margin mode
- `SetMarginMode(symbol, mode)` - Set margin mode
- `GetLeverage(symbol, recvWindow)` - Get leverage
- `SetLeverage(symbol, side, leverage, recvWindow)` - Set leverage
- `GetPositionMargin(symbol)` - Get position margin
- `SetPositionMargin(symbol, positionSide, amount, type)` - Set position margin
- `GetBalanceHistory(coin, limit)` - Get balance history
- `GetAccountPermissions()` - Get API permissions
- `GetAPIKey()` - Get API key info
- `GetUserCommissionRates(symbol)` - Get commission rates
- `GetAPIRateLimits()` - Get rate limits
- `GetDepositHistory(coin, limit)` - Get deposit history
- `GetWithdrawHistory(coin, limit)` - Get withdrawal history
- `GetAssetDetails(asset)` - Get asset details
- `GetAllAssets()` - Get all assets
- `GetFundingWallet(asset)` - Get funding wallet
- `DustTransfer(assets)` - Convert dust to BNB

</details>

<details>
<summary><b>Trade Service Methods</b></summary>

- `CreateOrder(params)` - Create new order
- `ModifyOrder(...)` - Modify existing order
- `CreateTestOrder(params)` - Create test order (no execution)
- `CloseAllPositions(symbol, ...)` - Close all positions
- `GetMarginType(symbol, ...)` - Get margin type
- `ChangeMarginType(symbol, type, ...)` - Change margin type
- `GetLeverage(symbol, ...)` - Get leverage
- `SetLeverage(symbol, leverage, ...)` - Set leverage
- `CreateBatchOrders(orders)` - Create multiple orders
- `CancelOrder(symbol, orderID, ...)` - Cancel order
- `CancelAllOrders(...)` - Cancel all orders
- `CancelBatchOrders(...)` - Cancel multiple orders
- `GetOrder(symbol, orderID)` - Get order details
- `GetOpenOrders(symbol, limit)` - Get open orders
- `GetOrderHistory(...)` - Get order history
- `GetFilledOrders(...)` - Get filled orders
- `GetUserTrades(...)` - Get user trades
- `ChangeLeverage(symbol, side, leverage, ...)` - Change leverage
- `CalculateFuturesCommission(margin, leverage, rate)` - Calculate commission
- `GetCommissionAmount(margin, leverage)` - Get commission amount

</details>

### Response Format

All API methods return `(map[string]interface{}, error)`:

```go
response, err := client.Market().GetLatestPrice("BTC-USDT")
if err != nil {
    // Handle error
}

// Access response data
if data, ok := response["data"].(map[string]interface{}); ok {
    price := data["price"]
    fmt.Printf("Price: %v\n", price)
}
```

---

## Error Handling

```go
import "github.com/tigusigalpa/bingx-go/errors"

// Handle errors
order, err := client.Trade().CreateOrder(params)
if err != nil {
    switch e := err.(type) {
    case *errors.AuthenticationException:
        log.Printf("Authentication error: %v", e)
    case *errors.RateLimitException:
        log.Printf("Rate limit exceeded: %v", e)
    case *errors.InsufficientBalanceException:
        log.Printf("Insufficient balance: %v", e)
    case *errors.APIException:
        log.Printf("API error [%s]: %v", e.APICode, e)
    default:
        log.Printf("Error: %v", err)
    }
    return
}
```

### Error Types

- `BingXException` - Base exception for all errors
- `APIException` - API-level errors with error codes
- `AuthenticationException` - Authentication failures
- `RateLimitException` - Rate limit exceeded
- `InsufficientBalanceException` - Insufficient balance errors

---

## Best Practices

### Security

```go
// DO: Use environment variables for credentials
apiKey := os.Getenv("BINGX_API_KEY")
apiSecret := os.Getenv("BINGX_API_SECRET")

// ❌ DON'T: Hardcode credentials
apiKey := "your_api_key_here" // Never do this!

// ✅ DO: Restrict API key permissions
// Only enable permissions you need (read-only for monitoring, etc.)

// ✅ DO: Use IP whitelist
// Configure IP restrictions in BingX API settings
```

### Error Handling

```go
// ✅ DO: Always check errors
result, err := client.Trade().CreateOrder(params)
if err != nil {
    // Handle specific error types
    switch e := err.(type) {
    case *errors.RateLimitException:
        // Implement exponential backoff
        time.Sleep(time.Second * 5)
    case *errors.InsufficientBalanceException:
        // Log and alert
        log.Error("Insufficient balance")
    default:
        // Generic error handling
        log.Error(err)
    }
    return
}

// ❌ DON'T: Ignore errors
result, _ := client.Trade().CreateOrder(params) // Bad practice
```

### Performance

```go
// ✅ DO: Reuse client instances
var globalClient *bingx.Client

func init() {
    globalClient = bingx.NewClient(apiKey, apiSecret)
}

// ❌ DON'T: Create new client for each request
func makeRequest() {
    client := bingx.NewClient(apiKey, apiSecret) // Wasteful
    client.Market().GetLatestPrice("BTC-USDT")
}

// ✅ DO: Use batch operations when possible
orders := []map[string]interface{}{...}
client.Trade().CreateBatchOrders(orders)

// ❌ DON'T: Make individual requests in loops
for _, order := range orders {
    client.Trade().CreateOrder(order) // Inefficient
}
```

### Concurrency

```go
// ✅ DO: Client is thread-safe
var wg sync.WaitGroup
for _, symbol := range symbols {
    wg.Add(1)
    go func(s string) {
        defer wg.Done()
        price, _ := client.Market().GetLatestPrice(s)
        // Process price
    }(symbol)
}
wg.Wait()
```

---

## Testing

### Unit testing

```go
package main

import (
    "testing"
    bingx "github.com/tigusigalpa/bingx-go"
)

func TestMarketData(t *testing.T) {
    client := bingx.NewClient("", "") // Empty credentials for public endpoints
    
    symbols, err := client.Market().GetFuturesSymbols()
    if err != nil {
        t.Fatalf("Failed to get symbols: %v", err)
    }
    
    if symbols == nil {
        t.Error("Expected symbols, got nil")
    }
}
```

### Integration Testing

```go
func TestOrderCreation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    client := bingx.NewClient(
        os.Getenv("BINGX_API_KEY"),
        os.Getenv("BINGX_API_SECRET"),
    )
    
    // Test order creation
    order, err := client.Trade().CreateTestOrder(params)
    if err != nil {
        t.Fatalf("Test order failed: %v", err)
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestMarketData

# Skip integration tests
go test -short ./...
```

---

## Performance

- **Request Latency**: < 50ms (network dependent)
- **Memory Usage**: ~2MB per client instance
- **Concurrent Requests**: Supports 1000+ concurrent goroutines
- **Signature Generation**: < 1ms per request

### Optimization Tips

1. **Connection Pooling**: HTTP client reuses connections automatically
2. **Batch Operations**: Use batch methods for multiple orders
3. **Caching**: Cache market data locally when appropriate
4. **Rate Limiting**: Implement client-side rate limiting to avoid API limits

```go
// Example: Rate-limited requests
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Limit(10), 1) // 10 requests per second

func makeRequest() {
    limiter.Wait(context.Background())
    client.Market().GetLatestPrice("BTC-USDT")
}
```

---

## Roadmap

### Planned (v1.1.0)

- Order builder with fluent interface
- Rate limiter middleware
- Request/response logging
- Retry with exponential backoff

### Future (v2.0.0)

- Built-in technical indicators
- Strategy backtesting framework
- Portfolio management tools
- Event-driven architecture

---

## Creating API Keys

1. Login to your BingX account
2. Go to [API Management](https://bingx.com/en-US/accounts/api)
3. Create API key, configure permissions (Read, Trade)
4. Set IP whitelist (recommended)
5. Save your API Key and Secret Key — the secret is shown only once

---

## Contributing

### Development setup

```bash
# Clone the repository
git clone https://github.com/tigusigalpa/bingx-go.git
cd bingx-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Add comments for exported functions

---

## License

MIT License — see [LICENSE](LICENSE).

## Disclaimer

This SDK is provided "as-is". Cryptocurrency trading carries significant financial risk. Always test in a sandbox environment before production use. The authors are not responsible for any financial losses.

## Links

- [GitHub Repository](https://github.com/tigusigalpa/bingx-go)
- [Issues](https://github.com/tigusigalpa/bingx-go/issues)
- [BingX API Documentation](https://bingx-api.github.io/docs/)
- [GoDoc](https://pkg.go.dev/github.com/tigusigalpa/bingx-go)
