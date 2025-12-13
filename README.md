# BingX Go SDK

![BingX Go SDK](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

**Full-Featured Go SDK for BingX API**

USDT-M & Coin-M Futures | Market Data | WebSocket Streams | Production Ready

---

## 📖 Table of Contents

- [About](#-about)
- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Usage Examples](#-usage-examples)
- [Documentation](#-documentation)
- [Error Handling](#-error-handling)
- [License](#-license)

---

## ✨ About

**BingX Go SDK** is a professional, feature-rich library for working with the BingX API (USDT-M & Coin-M Perpetual Futures).

Built with modern Go practices and provides:

- ✅ **100% coverage** of USDT-M Perpetual Futures API
- ✅ **Coin-M Perpetual Futures** fully implemented
- ✅ **Modular architecture** with separate services
- ✅ **Advanced error handling** with custom exceptions
- ✅ **WebSocket support** for streaming data
- ✅ **Complete security** with HMAC-SHA256 signatures
- ✅ **Type-safe** with proper Go interfaces
- ✅ **Production ready** with comprehensive error handling

---

## 🚀 Features

### 📊 Supported Services

| Service                      | Description                                         |
|------------------------------|-----------------------------------------------------|
| **USDT-M Perpetual Futures** |                                                     |
| 🏪 **Market Service**        | Market data, symbols, prices, candles               |
| 👤 **Account Service**       | Balance, positions, leverage, margin, assets        |
| 🔄 **Trade Service**         | Orders, trade history, position management          |
| 💰 **Wallet Service**        | Deposits, withdrawals, wallet addresses             |
| 💵 **Spot Account Service**  | Spot balance, transfers, internal transfers         |
| 👥 **Sub-Account Service**   | Sub-account management, API keys, transfers         |
| 🔄 **Copy Trading Service**  | Copy trading for futures and spot                   |
| 📋 **Contract Service**      | Standard contract API                               |
| 🔐 **Listen Key Service**    | WebSocket authentication                            |
| **Coin-M Perpetual Futures** |                                                     |
| 🪙 **Coin-M Market**         | Contract info, ticker, depth, klines, open interest |
| 🪙 **Coin-M Trade**          | Orders, positions, leverage, margin, balance        |
| 🪙 **Coin-M Listen Key**     | WebSocket authentication for Coin-M                 |

---

## 📦 Installation

### Requirements

- Go 1.21 or higher

### Install Package

```bash
go get github.com/tigusigalpa/bingx-go
```

---

## ⚡ Quick Start

### Basic Usage

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

## 📚 Usage Examples

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

---

## 📖 Documentation

### Client Options

```go
// Create client with custom options
client := bingx.NewClient(
    apiKey,
    apiSecret,
    bingx.WithBaseURI("https://open-api.bingx.com"),
    bingx.WithSourceKey("optional-source-key"),
    bingx.WithSignatureEncoding("base64"), // or "hex"
)
```

### Available Services

- `client.Market()` - Market data operations
- `client.Account()` - Account management
- `client.Trade()` - Trading operations
- `client.Contract()` - Standard contract operations
- `client.ListenKey()` - WebSocket authentication
- `client.Wallet()` - Wallet operations
- `client.SpotAccount()` - Spot account operations
- `client.SubAccount()` - Sub-account management
- `client.CopyTrading()` - Copy trading operations
- `client.CoinM()` - Coin-M perpetual futures

---

## 🛡️ Error Handling

The SDK provides custom error types for different scenarios:

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

## 🔑 Creating API Keys

1. Go to [BingX API Settings](https://bingx.com/en-US/accounts/api)
2. Click "Create API"
3. Save your **API Key** and **Secret Key** in a secure location
4. Configure access rights
5. ⚠️ Secret Key is displayed only once!

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## ⚠️ Disclaimer

This SDK is provided as-is. Trading cryptocurrencies carries risk. Always test thoroughly before using in production.

---

## 📧 Support

For issues and questions:
- GitHub Issues: [https://github.com/tigusigalpa/bingx-go/issues](https://github.com/tigusigalpa/bingx-go/issues)
- BingX API Documentation: [https://bingx-api.github.io/docs/](https://bingx-api.github.io/docs/)
