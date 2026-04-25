# BingX Demo Trading (VST) Examples

This directory contains examples of using the BingX Go SDK for demo trading in the VST (Virtual Simulation Trading) environment.

## What is VST?

VST stands for **Virtual Simulation Trading**, which is BingX's demo trading environment. It allows you to:
- Test trading strategies without real money
- Practice using the full BingX API functionality
- Develop and debug applications in a risk-free environment
- Access all the same features as the live production environment

## Key Differences

| Feature | Live Environment | VST (Demo) Environment |
|---------|------------------|------------------------|
| Base URL | `https://open-api.bingx.com` | `https://open-api-vst.bingx.com` |
| Real Money | Yes | No (simulated) |
| Data Feed | Real market data | Real market data |
| Order Execution | Real trades | Simulated trades |
| Risk | Financial risk | No financial risk |

## Quick Start

### 1. Create a Demo Client

```go
package main

import (
    "fmt"
    bingx "github.com/tigusigalpa/bingx-go/v2"
)

func main() {
    // Create a demo trading client
    demoClient := bingx.NewDemoClient(
        "YOUR_API_KEY",
        "YOUR_API_SECRET",
        bingx.WithSourceKey("my-demo-app"),
    )
    
    fmt.Printf("Demo endpoint: %s\n", demoClient.GetEndpoint())
    // Output: https://open-api-vst.bingx.com
}
```

### 2. Alternative: Manual Demo Configuration

```go
// Method 1: Using WithDemoEnvironment option
demoClient := bingx.NewClient(
    "YOUR_API_KEY",
    "YOUR_API_SECRET",
    bingx.WithDemoEnvironment(),
)

// Method 2: Setting VST URL directly
demoClient := bingx.NewClient(
    "YOUR_API_KEY",
    "YOUR_API_SECRET",
    bingx.WithBaseURI("https://open-api-vst.bingx.com"),
)
```

### 3. Check VST Status

```go
// Get VST information
vstInfo, err := demoClient.Trade().GetVst(nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("VST Info: %+v\n", vstInfo)
```

## Available Operations in Demo Mode

All operations that work in the live environment also work in demo mode:

### Market Data
```go
// Get market prices (same as live)
price, err := demoClient.Market().GetLatestPrice("BTC-USDT")
depth, err := demoClient.Market().GetDepth("BTC-USDT", 20)
klines, err := demoClient.Market().GetKlines("BTC-USDT", "1h", 100, nil, nil)
```

### Account Management
```go
// Get demo account balance
balance, err := demoClient.Account().GetBalance()

// Get demo positions
positions, err := demoClient.Account().GetPositions(nil)

// Set leverage (demo)
err = demoClient.Account().SetLeverage("BTC-USDT", "BOTH", 10, nil)
```

### Trading Operations
```go
// Place a demo order
order, err := demoClient.Trade().CreateOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         "LIMIT",
    "positionSide": "LONG",
    "price":        50000.0,
    "quantity":     0.001,
})

// Place a test order (no execution)
testOrder, err := demoClient.Trade().CreateTestOrder(map[string]interface{}{
    "symbol":       "ETH-USDT",
    "side":         "SELL",
    "type":         "LIMIT",
    "positionSide": "SHORT",
    "price":        3000.0,
    "quantity":     0.01,
})

// Cancel demo orders
err = demoClient.Trade().CancelOrder("BTC-USDT", &orderID, nil, nil, nil)
```

### WebSocket Streaming
```go
// Create market data stream (same as live)
stream := demoClient.NewMarketDataStream()
if err := stream.Connect(); err != nil {
    log.Fatal(err)
}
defer stream.Disconnect()

// Subscribe to data
stream.SubscribeTrade("BTC-USDT")
stream.SubscribeKline("BTC-USDT", "1m")
stream.Listen()
```

## Best Practices

### 1. Use Separate API Keys
- Create dedicated API keys for demo trading
- Label them clearly in your BingX account
- Never use live API keys in demo code

### 2. Environment Detection
```go
func isDemoEnvironment(client *bingx.Client) bool {
    return client.GetEndpoint() == "https://open-api-vst.bingx.com"
}

// Usage
if isDemoEnvironment(demoClient) {
    fmt.Println("Running in DEMO mode - no real money at risk")
} else {
    fmt.Println("Running in LIVE mode - real trades will be executed")
}
```

### 3. Configuration Management
```go
type Config struct {
    APIKey     string
    APISecret  string
    IsDemo     bool
    SourceKey  string
}

func NewClientFromConfig(cfg Config) *bingx.Client {
    if cfg.IsDemo {
        return bingx.NewDemoClient(cfg.APIKey, cfg.APISecret, 
            bingx.WithSourceKey(cfg.SourceKey))
    }
    return bingx.NewClient(cfg.APIKey, cfg.APISecret,
        bingx.WithSourceKey(cfg.SourceKey))
}
```

### 4. Testing Strategy
1. **Development**: Use demo environment exclusively
2. **Staging**: Test with demo environment using production-like data
3. **Production**: Switch to live environment after thorough testing

## Common Use Cases

### 1. Strategy Development
```go
// Test a moving average crossover strategy
func testMAStrategy(client *bingx.Client) {
    // Get historical data
    klines, _ := client.Market().GetKlines("BTC-USDT", "1h", 200, nil, nil)
    
    // Implement strategy logic
    // Place demo orders based on signals
    
    // Monitor performance without risk
}
```

### 2. Risk Management Testing
```go
// Test stop-loss and take-profit orders
func testRiskManagement(client *bingx.Client) {
    // Place orders with SL/TP
    order, _ := client.Trade().CreateOrder(map[string]interface{}{
        "symbol":       "BTC-USDT",
        "side":         "BUY",
        "type":         "MARKET",
        "positionSide": "LONG",
        "quantity":     0.001,
        "stopLoss": map[string]interface{}{
            "type":      "STOP_MARKET",
            "stopPrice": 49000.0,
        },
        "takeProfit": map[string]interface{}{
            "type":      "TAKE_PROFIT_MARKET",
            "stopPrice": 51000.0,
        },
    })
}
```

### 3. API Integration Testing
```go
// Test all API endpoints
func testAPIIntegration(client *bingx.Client) {
    // Test market data
    client.Market().GetLatestPrice("BTC-USDT")
    
    // Test account info
    client.Account().GetBalance()
    
    // Test trading
    client.Trade().CreateTestOrder(map[string]interface{}{
        "symbol": "BTC-USDT",
        "side":   "BUY",
        "type":   "MARKET",
        "quantity": 0.001,
    })
}
```

## Running the Examples

1. **Install dependencies**:
   ```bash
   cd examples/demo_trading
   go mod init demo-trading
   go mod tidy
   ```

2. **Set your API credentials**:
   ```bash
   export BINGX_API_KEY="your_api_key"
   export BINGX_API_SECRET="your_api_secret"
   ```

3. **Run the example**:
   ```bash
   go run main.go
   ```

## Important Notes

- **No Real Money**: Demo trading uses simulated funds only
- **Real Market Data**: Prices and market data are real-time
- **Same API**: All API methods work identically to live environment
- **Rate Limits**: Demo environment has the same rate limits as live
- **No Slippage**: Demo orders may have different execution characteristics

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Ensure API keys are valid for demo trading
   - Check that the keys have the necessary permissions

2. **Connection Issues**
   - Verify the endpoint is `https://open-api-vst.bingx.com`
   - Check network connectivity

3. **Order Rejections**
   - Ensure sufficient demo balance
   - Check order parameters (price, quantity, etc.)

### Getting Help

- Check the [BingX API Documentation](https://bingx-api.github.io/docs/)
- Review the main [README.md](../../README.md)
- Open an issue on GitHub for SDK-specific problems

## Next Steps

1. Explore the full example in `main.go`
2. Implement your own trading strategy
3. Test thoroughly in demo environment
4. Consider production deployment after successful testing
