package main

import (
	"fmt"

	bingx "github.com/tigusigalpa/bingx-go/v2"
)

func main() {
	fmt.Println("=== BingX Demo Trading (VST) Example ===\n")

	// Create a demo trading client
	// This automatically uses the VST (Virtual Simulation Trading) environment
	demoClient := bingx.NewDemoClient(
		"YOUR_API_KEY",
		"YOUR_API_SECRET",
		bingx.WithSourceKey("demo-trading-example"),
	)

	fmt.Printf("Demo client created with endpoint: %s\n", demoClient.GetEndpoint())
	fmt.Println("Note: This is the VST (Virtual Simulation Trading) environment\n")

	// 1. Check VST status
	fmt.Println("1. Checking VST status...")
	vstInfo, err := demoClient.Trade().GetVst(nil)
	if err != nil {
		fmt.Printf("   Error getting VST info: %v\n", err)
	} else {
		fmt.Printf("   VST Info: %+v\n", vstInfo)
	}

	// 2. Get market data (works in demo mode)
	fmt.Println("\n2. Getting market data...")
	price, err := demoClient.Market().GetLatestPrice("BTC-USDT")
	if err != nil {
		fmt.Printf("   Error getting price: %v\n", err)
	} else {
		fmt.Printf("   BTC-USDT Price: %v\n", price)
	}

	// 3. Get account balance (demo balance)
	fmt.Println("\n3. Getting demo account balance...")
	balance, err := demoClient.Account().GetBalance()
	if err != nil {
		fmt.Printf("   Error getting balance: %v\n", err)
	} else {
		fmt.Printf("   Account Balance: %+v\n", balance)
	}

	// 4. Place a demo order
	fmt.Println("\n4. Placing a demo order...")
	order, err := demoClient.Trade().CreateOrder(map[string]interface{}{
		"symbol":       "BTC-USDT",
		"side":         "BUY",
		"type":         "LIMIT",
		"positionSide": "LONG",
		"price":        50000.0,
		"quantity":     0.001,
	})
	if err != nil {
		fmt.Printf("   Error creating order: %v\n", err)
	} else {
		fmt.Printf("   Demo Order Created: %+v\n", order)
	}

	// 5. Place a test order (doesn't execute)
	fmt.Println("\n5. Placing a test order (no execution)...")
	testOrder, err := demoClient.Trade().CreateTestOrder(map[string]interface{}{
		"symbol":       "ETH-USDT",
		"side":         "SELL",
		"type":         "LIMIT",
		"positionSide": "SHORT",
		"price":        3000.0,
		"quantity":     0.01,
	})
	if err != nil {
		fmt.Printf("   Error creating test order: %v\n", err)
	} else {
		fmt.Printf("   Test Order Result: %+v\n", testOrder)
	}

	// 6. Get open orders
	fmt.Println("\n6. Getting open demo orders...")
	orders, err := demoClient.Trade().GetOpenOrders(nil, 100)
	if err != nil {
		fmt.Printf("   Error getting open orders: %v\n", err)
	} else {
		fmt.Printf("   Open Orders: %+v\n", orders)
	}

	// 7. Get positions
	fmt.Println("\n7. Getting demo positions...")
	positions, err := demoClient.Account().GetPositions(nil)
	if err != nil {
		fmt.Printf("   Error getting positions: %v\n", err)
	} else {
		fmt.Printf("   Positions: %+v\n", positions)
	}

	fmt.Println("\n=== Demo Trading Examples Complete ===")
	fmt.Println("\nKey Points:")
	fmt.Println("- All operations are performed in the VST (Virtual Simulation Trading) environment")
	fmt.Println("- No real money is used - this is simulated trading")
	fmt.Println("- The VST environment mirrors the live production environment")
	fmt.Println("- Use this environment for testing strategies without financial risk")
}

// Example of switching between live and demo environments
func demonstrateEnvironmentSwitching() {
	fmt.Println("\n=== Environment Switching Example ===")

	apiKey := "YOUR_API_KEY"
	apiSecret := "YOUR_API_SECRET"

	// Live production client
	liveClient := bingx.NewClient(apiKey, apiSecret)
	fmt.Printf("Live Client Endpoint: %s\n", liveClient.GetEndpoint())

	// Demo client
	demoClient := bingx.NewDemoClient(apiKey, apiSecret)
	fmt.Printf("Demo Client Endpoint: %s\n", demoClient.GetEndpoint())

	// Manual demo configuration
	manualDemoClient := bingx.NewClient(
		apiKey,
		apiSecret,
		bingx.WithDemoEnvironment(),
	)
	fmt.Printf("Manual Demo Client Endpoint: %s\n", manualDemoClient.GetEndpoint())
}
