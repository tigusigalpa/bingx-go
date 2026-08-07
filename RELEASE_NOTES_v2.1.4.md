# BingX Go SDK v2.1.4 Release Notes

## 🎉 New Features

### ✨ Demo Trading Support (VST Environment)

We're excited to introduce comprehensive demo trading support through BingX's VST (Virtual Simulation Trading) environment! This feature allows you to test trading strategies without financial risk.

#### Key Additions:

**🏗️ New Client Constructor**
- `NewDemoClient()` - One-line demo client creation
- Automatically configures VST endpoint (`https://open-api-vst.bingx.com`)
- Full compatibility with all existing API methods

**⚙️ Configuration Options**
- `WithDemoEnvironment()` - Configure any client for demo trading
- Seamless switching between live and demo environments

**📊 VST Status API**
- `GetVst()` method to retrieve demo trading information
- Verify VST environment status and capabilities

#### Usage Examples:

```go
// Create demo client (recommended)
demoClient := bingx.NewDemoClient("API_KEY", "API_SECRET")

// Alternative: Configure demo environment
demoClient := bingx.NewClient("API_KEY", "API_SECRET", 
    bingx.WithDemoEnvironment())

// Check environment
fmt.Printf("Endpoint: %s\n", demoClient.GetEndpoint())
// Output: https://open-api-vst.bingx.com

// Get VST information
vstInfo, err := demoClient.Trade().GetVst(nil)

// Place demo orders
order, err := demoClient.Trade().CreateOrder(map[string]interface{}{
    "symbol": "BTC-USDT",
    "side": "BUY", 
    "type": "LIMIT",
    "price": 50000.0,
    "quantity": 0.001,
})
```

#### 🎯 Benefits:

- **Risk-Free Testing**: Practice trading strategies without real money
- **Full API Access**: All 240+ API methods available in demo mode
- **Real Market Data**: Uses live market prices for realistic simulation
- **Easy Migration**: Seamless switch between demo and live environments
- **Development Friendly**: Perfect for testing and debugging

#### 📁 New Examples:

- `examples/demo_trading/main.go` - Complete demo trading example
- `examples/demo_trading/README.md` - Comprehensive demo trading guide

---

## 📋 API Reference Updates

### New Methods

#### TradeService
- `GetVst(recvWindow *int64)` - Retrieve VST (Virtual Simulation Trading) information

#### Client Methods
- `NewDemoClient(apiKey, apiSecret string, options...)` - Create demo trading client
- `WithDemoEnvironment()` - Configure client for VST environment

### Updated Documentation

- **Configuration Section**: Added demo trading setup instructions
- **API Reference**: Updated with new demo trading methods
- **Examples**: Added comprehensive demo trading examples

---

## 🔧 Technical Details

### Environment Endpoints

| Environment | Endpoint | Purpose |
|-------------|----------|---------|
| **Live** | `https://open-api.bingx.com` | Real trading with actual funds |
| **Demo (VST)** | `https://open-api-vst.bingx.com` | Simulated trading with virtual funds |

### Backward Compatibility

✅ **100% Backward Compatible** - No breaking changes introduced

All existing code continues to work without modification. Demo trading features are opt-in and don't affect existing functionality.

---

## 🚀 Getting Started

### Installation

```bash
go get github.com/tigusigalpa/bingx-go/v2@v2.1.4
```

### Quick Demo Setup

```go
package main

import (
    "fmt"
    bingx "github.com/tigusigalpa/bingx-go/v2"
)

func main() {
    // Create demo client
    client := bingx.NewDemoClient("YOUR_API_KEY", "YOUR_API_SECRET")
    
    // Start demo trading
    price, _ := client.Market().GetLatestPrice("BTC-USDT")
    fmt.Printf("BTC Price: %v\n", price)
    
    // Place demo order
    order, _ := client.Trade().CreateOrder(map[string]interface{}{
        "symbol": "BTC-USDT",
        "side": "BUY",
        "type": "MARKET", 
        "quantity": 0.001,
    })
    fmt.Printf("Demo Order: %v\n", order)
}
```

---

## 🐛 Bug Fixes

- Fixed lint error in demo trading example (removed unused import)
- Improved error handling for VST endpoint configuration

---

## 📚 Documentation

- **New**: Demo Trading Guide (`examples/demo_trading/README.md`)
- **Updated**: Main README with demo trading configuration
- **Enhanced**: API reference with new methods

---

## 🔍 Migration Notes

### From v2.1.3 to v2.1.4

No code changes required! This is a feature-only release.

To start using demo trading:

```go
// Existing code continues to work
client := bingx.NewClient("API_KEY", "API_SECRET")

// New demo functionality (optional)
demoClient := bingx.NewDemoClient("API_KEY", "API_SECRET")
```

---

## 🙏 Acknowledgments

Special thanks to the BingX team for providing the VST (Virtual Simulation Trading) environment, enabling risk-free strategy testing and development.

---

## 📞 Support

- 📖 [Documentation](https://github.com/tigusigalpa/bingx-go/wiki)
- 🐛 [Issues](https://github.com/tigusigalpa/bingx-go/issues)
- 💬 [Discussions](https://github.com/tigusigalpa/bingx-go/discussions)

---

**Full Changelog**: [CHANGELOG.md](CHANGELOG.md)  
**Previous Release**: [v2.1.3](https://github.com/tigusigalpa/bingx-go/releases/tag/v2.1.3)
