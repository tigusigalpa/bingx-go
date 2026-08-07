# Getting Started

Let's get you up and running with bingx-go. By the end of this guide, you'll have a working connection to BingX and be ready to fetch data or place your first trade.

## What You'll Need

Before we dive in, make sure you have:

- **Go 1.21+** installed on your machine
- A **BingX account** (you'll need API keys for anything beyond public data)
- Some basic Go knowledge (if you can write a "Hello World", you're good)

## Installation

Grab the SDK with a single command:

```bash
go get github.com/tigusigalpa/bingx-go/v2
```

Then tidy up your dependencies:

```bash
go mod tidy
```

That's it! The only external dependency is `gorilla/websocket` for real-time streaming, and Go will handle that automatically.

## Getting Your API Keys

To do anything beyond fetching public market data, you'll need API credentials:

1. Log in to [BingX](https://bingx.com) and head to **API Management**
2. Hit **Create API Key**
3. Set up permissions:
   - ✅ **Read** — so you can see your account data
   - ✅ **Trade** — so you can place orders
   - ⚠️ **Withdraw** — skip this unless you really need it (your bot probably doesn't)
4. **Add your IP to the whitelist** — seriously, do this. It's the single best thing you can do for security
5. Copy your **API Key** and **Secret Key** somewhere safe

> 🚨 **Heads up**: Never, ever commit your API keys to Git or share them with anyone. If someone gets your keys, they get access to your account.

## Quick Start

Let's make sure everything works. Here's the simplest possible example:

```go
package main

import (
    "fmt"
    "log"
    
    bingx "github.com/tigusigalpa/bingx-go/v2"
)

func main() {
    // Plug in your credentials
    client := bingx.NewClient(
        "YOUR_API_KEY",
        "YOUR_API_SECRET",
    )
    
    // Let's see what BTC is trading at
    price, err := client.Market().GetLatestPrice("BTC-USDT")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("BTC is at: %v\n", price)
}
```

Run it. If you see a price, you're in business!

### Tweaking the Client

The defaults work fine for most people, but you can customize things if needed:

```go
client := bingx.NewClient(
    apiKey,
    apiSecret,

    // Different API endpoint (rarely needed)
    bingx.WithBaseURI("https://open-api.bingx.com"),

    // HMAC signature encoding: "hex" is the BingX-compatible default.
    // "base64" is still supported for backward compatibility only.
    bingx.WithSignatureEncoding("hex"),

    // Tag your requests for easier debugging
    bingx.WithSourceKey("MyAwesomeBot"),
)
```

## Building Your First Bot

Ready for something more substantial? Here's a complete example that shows off the main features:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    bingx "github.com/tigusigalpa/bingx-go/v2"
    "github.com/tigusigalpa/bingx-go/v2/services"
)

func main() {
    // Load API credentials from environment variables
    apiKey := os.Getenv("BINGX_API_KEY")
    apiSecret := os.Getenv("BINGX_API_SECRET")
    
    if apiKey == "" || apiSecret == "" {
        log.Fatal("Please set BINGX_API_KEY and BINGX_API_SECRET environment variables")
    }
    
    // Create client
    client := bingx.NewClient(apiKey, apiSecret)
    
    // 1. Get Market Data
    fmt.Println("=== Market Data ===")
    price, err := client.Market().GetLatestPrice("BTC-USDT")
    if err != nil {
        log.Printf("Error getting price: %v", err)
    } else {
        fmt.Printf("BTC-USDT Price: %v\n", price)
    }
    
    // 2. Check Account Balance
    fmt.Println("\n=== Account Balance ===")
    balance, err := client.Account().GetBalance()
    if err != nil {
        log.Printf("Error getting balance: %v", err)
    } else {
        fmt.Printf("Balance: %v\n", balance)
    }
    
    // 3. Get Current Positions
    fmt.Println("\n=== Positions ===")
    positions, err := client.Account().GetPositions(nil)
    if err != nil {
        log.Printf("Error getting positions: %v", err)
    } else {
        fmt.Printf("Positions: %v\n", positions)
    }
    
    // 4. Place a Test Order (doesn't execute)
    fmt.Println("\n=== Test Order ===")
    testOrder, err := client.Trade().CreateTestOrder(map[string]interface{}{
        "symbol":       "BTC-USDT",
        "side":         "BUY",
        "type":         services.OrderTypeLimit,
        "positionSide": "LONG",
        "price":        50000.0,
        "quantity":     0.001,
    })
    if err != nil {
        log.Printf("Error creating test order: %v", err)
    } else {
        fmt.Printf("Test Order: %v\n", testOrder)
    }
    
    fmt.Println("\n✅ Bot initialized successfully!")
}
```

### Running It

```bash
# Set your credentials
export BINGX_API_KEY="your_api_key"
export BINGX_API_SECRET="your_api_secret"

# Fire it up
go run main.go
```

## Managing Your API Keys

Hardcoding credentials is a recipe for disaster. Here's how to do it properly.

### On Mac/Linux

Add these to your `~/.bashrc` or `~/.zshrc`:

```bash
export BINGX_API_KEY="your_api_key_here"
export BINGX_API_SECRET="your_api_secret_here"
```

Then reload: `source ~/.bashrc`

### On Windows (PowerShell)

```powershell
$env:BINGX_API_KEY="your_api_key_here"
$env:BINGX_API_SECRET="your_api_secret_here"
```

### The Better Way: Use a .env File

This is what most developers do. First, grab the godotenv package:

```bash
go get github.com/joho/godotenv
```

Create a `.env` file in your project root:

```
BINGX_API_KEY=your_api_key_here
BINGX_API_SECRET=your_api_secret_here
```

Load it at startup:

```go
import "github.com/joho/godotenv"

func init() {
    godotenv.Load() // Silently fails if no .env file, which is fine
}
```

> 🚨 **Don't forget**: Add `.env` to your `.gitignore` right now. Seriously, do it before you forget.

## Making Sure It All Works

### Step 1: Test Public Data (No Keys Needed)

This should work even without API credentials:

```go
// What trading pairs are available?
symbols, err := client.Market().GetFuturesSymbols()

// What's the current BTC price?
price, err := client.Market().GetLatestPrice("BTC-USDT")

// What does the order book look like?
depth, err := client.Market().GetDepth("BTC-USDT", 20)
```

### Step 2: Test Your Credentials

If this works, your API keys are set up correctly:

```go
balance, err := client.Account().GetBalance()
if err != nil {
    log.Fatal("Hmm, authentication failed:", err)
}
fmt.Println("✅ You're connected!")
```

### Step 3: Test Trading (Without Risking Money)

BingX has a test endpoint that validates your order without actually placing it:

```go
testOrder, err := client.Trade().CreateTestOrder(map[string]interface{}{
    "symbol":       "BTC-USDT",
    "side":         "BUY",
    "type":         "LIMIT",
    "positionSide": "LONG",
    "price":        50000.0,
    "quantity":     0.001,
})
// If this succeeds, you're ready to trade for real
```

## Troubleshooting

### "Authentication failed" or "Invalid signature"

This is the most common issue. Try these fixes in order:

1. **Veoifubtlecsignaturh eeck your key`bingx.NewClient` now defsults to hex (`"hex"`). Only use — copy-paste errors are sneakybas64if have a specific backwad-ompatibity ned
2. **Try hex encoding** — add `bingx.WithSignatureEncoding("hex")` to your client
3. **Sync your clock** — if your system time is off by more than a few seconds, signatures will fail
4. **Check your IP whitelist** — if you set one up, make sure your current IP is on it

### "Too many requests"

You're hitting the rate limit. A few ways to fix this:

- **Add a rate limiter** to your code:

```go
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Limit(10), 1) // 10 requests per second max

func makeRequest() {
    limiter.Wait(context.Background())
    // Now make your API call
}
```

- **Use WebSockets** for real-time data instead of polling
- **Batch your requests** when possible

### Connection Timeouts

If you can't connect at all:

- Check your internet connection (obvious, but worth mentioning)
- Try accessing `https://open-api.bingx.com` in your browser
- Some regions may need a VPN
- Add retry logic with exponential backoff for production code

## Where to Go From Here

You've got the basics down. Here's a suggested learning path:

1. **[Market Data](Market-Service)** — Learn to fetch prices, charts, and order books
2. **[Account Management](Account-Service)** — Check your balance, positions, and settings
3. **[Trading](Trade-Service)** — Place, modify, and cancel orders
4. **[v3 Features](API-v3-Features)** — Try TWAP orders, multi-assets mode, and more
5. **[WebSockets](WebSocket-Service)** — Get real-time updates instead of polling

## Tips From the Trenches

### Keep It Secure

- Environment variables for credentials, always
- IP whitelist on BingX — it's free protection
- Read-only API keys for monitoring-only bots
- Rotate your keys periodically

### During Development

- Start with tiny amounts (or test orders)
- Log everything — you'll thank yourself later
- Handle errors properly from day one
- Keep an eye on rate limits

### In Production

- Add circuit breakers so one failure doesn't cascade
- Set up health checks and monitoring
- Use structured logging (JSON) for easier debugging
- Alert on errors — don't wait to discover problems
- Handle shutdowns gracefully (close positions if needed)

## A Clean Project Structure

As your bot grows, you'll want to organize things. Here's a structure that works well:

```
my-trading-bot/
├── main.go              # Entry point
├── go.mod
├── go.sum
├── .env                 # Your secrets (NEVER commit this)
├── .gitignore           # Make sure .env is in here
├── config/
│   └── config.go        # Load and validate configuration
├── strategies/
│   ├── strategy.go      # Interface that all strategies implement
│   └── simple.go        # Your first strategy
├── handlers/
│   ├── market.go        # React to market data
│   └── trading.go       # Execute trades
└── utils/
    ├── logger.go        # Consistent logging
    └── errors.go        # Custom error types
```

## More Reading

- **[v3 Features](API-v3-Features)** — The new hotness
- **[Error Handling](Error-Handling)** — Because things will go wrong
- **[Advanced Patterns](Advanced-Features)** — Rate limiting, batching, testing
- **[Official BingX Docs](https://bingx-api.github.io/docs-v3/)** — The authoritative source

---

**Stuck?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues) and we'll help you out.
