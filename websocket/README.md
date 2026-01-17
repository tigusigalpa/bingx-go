# BingX WebSocket API

This package provides WebSocket streaming functionality for BingX API, supporting both public market data and private account data streams.

## Features

- **Market Data Stream**: Real-time public market data (trades, klines, depth, tickers)
- **Account Data Stream**: Real-time private account updates (balance, positions, orders)
- **Automatic Ping/Pong**: Handles WebSocket keep-alive automatically
- **GZIP Decompression**: Automatically decompresses gzipped messages
- **Thread-Safe**: Safe for concurrent use

## Installation

```bash
go get github.com/tigusigalpa/bingx-go
```

## Market Data Stream

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/tigusigalpa/bingx-go"
)

func main() {
    client := bingx.NewClient("", "")
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
    
    // Subscribe to trade updates
    stream.SubscribeTrade("BTC-USDT")
    
    // Start listening
    stream.Listen()
}
```

### Available Subscriptions

#### Trade Updates
```go
stream.SubscribeTrade("BTC-USDT")
stream.UnsubscribeTrade("BTC-USDT")
```

#### Kline/Candlestick Updates
```go
// Intervals: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 12h, 1d, 3d, 1w, 1M
stream.SubscribeKline("BTC-USDT", "1m")
stream.UnsubscribeKline("BTC-USDT", "1m")
```

#### Depth/Orderbook Updates
```go
// Levels: 5, 10, 20, 50, 100
stream.SubscribeDepth("BTC-USDT", 20)
stream.UnsubscribeDepth("BTC-USDT", 20)
```

#### 24hr Ticker Updates
```go
stream.SubscribeTicker("BTC-USDT")
stream.UnsubscribeTicker("BTC-USDT")
```

#### Book Ticker Updates (Best Bid/Ask)
```go
stream.SubscribeBookTicker("BTC-USDT")
stream.UnsubscribeBookTicker("BTC-USDT")
```

## Account Data Stream

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/tigusigalpa/bingx-go"
)

func main() {
    apiKey := os.Getenv("BINGX_API_KEY")
    apiSecret := os.Getenv("BINGX_API_SECRET")
    
    client := bingx.NewClient(apiKey, apiSecret)
    
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
    
    // Listen for all account updates
    stream.OnAccountUpdate(func(eventType string, data map[string]interface{}) {
        fmt.Printf("Event [%s]: %+v\n", eventType, data)
    })
    
    // Start listening
    stream.Listen()
}
```

### Available Event Handlers

#### All Account Updates
```go
stream.OnAccountUpdate(func(eventType string, data map[string]interface{}) {
    // eventType: "account", "order", or "unknown"
    fmt.Printf("Event [%s]: %+v\n", eventType, data)
})
```

#### Balance Updates Only
```go
stream.OnBalanceUpdate(func(balances interface{}) {
    fmt.Printf("Balance: %+v\n", balances)
})
```

#### Position Updates Only
```go
stream.OnPositionUpdate(func(positions interface{}) {
    fmt.Printf("Positions: %+v\n", positions)
})
```

#### Order Updates Only
```go
stream.OnOrderUpdate(func(order interface{}) {
    fmt.Printf("Order: %+v\n", order)
})
```

## Advanced Usage

### Multiple Subscriptions

```go
stream := client.NewMarketDataStream()
stream.Connect()

// Subscribe to multiple data types
stream.SubscribeTrade("BTC-USDT")
stream.SubscribeTrade("ETH-USDT")
stream.SubscribeKline("BTC-USDT", "1m")
stream.SubscribeDepth("BTC-USDT", 20)

stream.Listen()
```

### Custom Request IDs

```go
// Provide custom request ID as second parameter
stream.SubscribeTrade("BTC-USDT", "my-custom-id")
```

### Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    stream := client.NewMarketDataStream()
    stream.Connect()
    defer stream.Disconnect()
    
    // Handle Ctrl+C
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        stream.Stop()
    }()
    
    stream.Listen()
}
```

### Listen Key Management

Listen keys expire after 60 minutes. You should extend them periodically:

```go
// Extend listen key every 30 minutes
ticker := time.NewTicker(30 * time.Minute)
go func() {
    for range ticker.C {
        client.ListenKey().Extend(listenKey)
    }
}()
```

## WebSocket Endpoints

- **Market Data**: `wss://open-api-swap.bingx.com/swap-market`
- **Account Data**: `wss://open-api-swap.bingx.com/swap-market?listenKey={listenKey}`

## Error Handling

```go
if err := stream.Connect(); err != nil {
    log.Printf("Connection error: %v", err)
    return
}

if err := stream.Listen(); err != nil {
    log.Printf("Listen error: %v", err)
}
```

## Thread Safety

All WebSocket client methods are thread-safe and can be called from multiple goroutines.

## Examples

See the `examples/` directory for complete working examples:
- `websocket_market_data.go` - Market data streaming example
- `websocket_account_data.go` - Account data streaming example

## Notes

- Messages are automatically decompressed if gzipped
- Ping/pong messages are handled automatically
- The client will continue listening until `Stop()` is called or an error occurs
- Multiple message handlers can be registered using `OnMessage()`
