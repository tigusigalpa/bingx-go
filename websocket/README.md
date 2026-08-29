# BingX WebSocket API

This package provides WebSocket streaming functionality for BingX API, supporting both public market data and private account data streams.

> ⚠️ **Migration for v2.4.x:** the direct client was renamed from `websocket.WebSocketClient` to `websocket.Client`. Replace `websocket.NewWebSocketClient(url)` with `websocket.NewClient(url)`. `bingx.Client.NewMarketDataStream()` and `bingx.Client.NewAccountDataStream()` are unchanged.

## Features

- **Market Data Stream**: Real-time public market data (trades, klines, depth, tickers)
- **Account Data Stream**: Real-time private account updates (balance, positions, orders)
- **Automatic Ping/Pong**: Handles WebSocket keep-alive automatically
- **GZIP Decompression**: Automatically decompresses gzipped messages
- **Thread-Safe**: Safe for concurrent use
- **Explicit lifecycle**: `Disconnect()` is idempotent and a client can be connected again after disconnecting

## Installation

```bash
go get github.com/tigusigalpa/bingx-go/v2
```

## Market Data Stream

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/tigusigalpa/bingx-go/v2"
)

func main() {
    client := bingx.NewClient("", "")
    stream := client.NewMarketDataStream()
    
    // Connect to WebSocket
    if err := stream.Connect(); err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := stream.Disconnect(); err != nil {
            log.Printf("disconnect: %v", err)
        }
    }()
    
    // Register message handler
    stream.OnMessage(func(data map[string]interface{}) {
        fmt.Printf("Received: %+v\n", data)
    })
    
    // Subscribe to trade updates
    if err := stream.SubscribeTrade("BTC-USDT"); err != nil {
        log.Fatal(err)
    }
    
    // Start listening
    if err := stream.Listen(); err != nil {
        log.Printf("stream stopped: %v", err)
    }
}
```

### Direct Client Usage

Use `websocket.Client` only when you need a custom endpoint or want to send raw subscription frames. For normal market and account streams, prefer the high-level helpers above.

```go
stream := websocket.NewClient("wss://open-api-swap.bingx.com/swap-market")
if err := stream.Connect(); err != nil {
    return err
}
defer func() { _ = stream.Disconnect() }()

stream.OnMessage(func(message map[string]interface{}) {
    // Process the decoded JSON message.
})
if err := stream.Subscribe("prices", "BTC-USDT@ticker"); err != nil {
    return err
}
return stream.Listen() // blocks until Disconnect, a read error, or a close frame
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
    
    "github.com/tigusigalpa/bingx-go/v2"
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
    defer func() {
        if err := stream.Disconnect(); err != nil {
            log.Printf("disconnect: %v", err)
        }
    }()
    
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
if err := stream.Connect(); err != nil {
    log.Fatal(err)
}
defer func() { _ = stream.Disconnect() }()

// Subscribe to multiple data types
if err := stream.SubscribeTrade("BTC-USDT"); err != nil { log.Fatal(err) }
if err := stream.SubscribeTrade("ETH-USDT"); err != nil { log.Fatal(err) }
if err := stream.SubscribeKline("BTC-USDT", "1m"); err != nil { log.Fatal(err) }
if err := stream.SubscribeDepth("BTC-USDT", 20); err != nil { log.Fatal(err) }

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
    if err := stream.Connect(); err != nil {
        log.Fatal(err)
    }
    
    // Handle Ctrl+C
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        // Disconnect closes the socket and unblocks a pending ReadMessage.
        if err := stream.Disconnect(); err != nil {
            log.Printf("disconnect: %v", err)
        }
    }()
    
    if err := stream.Listen(); err != nil {
        log.Printf("stream stopped: %v", err)
    }
}
```

### Listen Key Management

Listen keys expire after 60 minutes. You should extend them periodically:

```go
// Extend listen key every 30 minutes
ticker := time.NewTicker(30 * time.Minute)
go func() {
    for range ticker.C {
        if _, err := client.ListenKey().Extend(listenKey); err != nil {
            log.Printf("listen key extension failed: %v", err)
        }
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

`Listen()` is blocking. A server close, read error, or failed pong response is returned to the caller. The package does not reconnect automatically; create your own retry/backoff loop if the application requires persistent connectivity.

## Message and Subscription Protocol

`Subscribe` and `Unsubscribe` send BingX frames in the following form:

```json
{"id":"request-id","reqType":"sub","dataType":"BTC-USDT@trade"}
```

The typed market helpers construct `dataType` values such as `BTC-USDT@trade`, `BTC-USDT@kline_1m`, `BTC-USDT@depth20`, `BTC-USDT@ticker`, and `BTC-USDT@bookTicker`. Each helper accepts an optional request ID; otherwise the SDK creates one.

## Thread Safety

All WebSocket client methods are thread-safe and can be called from multiple goroutines. Writes are serialized so concurrent subscribe/unsubscribe calls do not violate Gorilla WebSocket's single-writer requirement. Register callbacks before calling `Listen()` when possible; callbacks run synchronously on the listener goroutine, so long-running work should be handed off to another goroutine or channel.

## Examples

See the `examples/` directory for complete working examples:
- `examples/websocket_market_data/main.go` - Market data streaming example
- `examples/websocket_account_data/main.go` - Account data streaming example

## Notes

- Messages are automatically decompressed if gzipped
- Ping/pong messages are handled automatically
- `Stop()` changes the listener state but does not interrupt an already-blocked network read; use `Disconnect()` for shutdown from another goroutine
- The client will continue listening until it is disconnected or a read error occurs
- Multiple message handlers can be registered using `OnMessage()`
