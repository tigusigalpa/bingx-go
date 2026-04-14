package main

import (
	"fmt"
	"log"
	"os"
	"time"

	bingx "github.com/tigusigalpa/bingx-go/v2"
	"github.com/tigusigalpa/bingx-go/v2/services"
)

func main() {
	apiKey := os.Getenv("BINGX_API_KEY")
	apiSecret := os.Getenv("BINGX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Please set BINGX_API_KEY and BINGX_API_SECRET environment variables")
	}

	client := bingx.NewClient(
		apiKey,
		apiSecret,
		bingx.WithBaseURI("https://open-api.bingx.com"),
	)

	fmt.Println("=== BingX API v3 Features Demo ===\n")

	// 1. New Order Types
	demonstrateNewOrderTypes(client)

	// 2. TWAP Orders
	demonstrateTWAPOrders(client)

	// 3. Position Management
	demonstratePositionManagement(client)

	// 4. Multi-Assets Mode
	demonstrateMultiAssetsMode(client)

	// 5. Enhanced Market Data
	demonstrateMarketData(client)

	// 6. Advanced Account Features
	demonstrateAccountFeatures(client)

	fmt.Println("\n=== Demo Complete ===")
}

func demonstrateNewOrderTypes(client *bingx.Client) {
	fmt.Println("1. New Order Types (v3)")
	fmt.Println("------------------------")

	// Trigger Limit Order
	fmt.Println("Available order types:")
	fmt.Printf("  - TRIGGER_LIMIT: %s\n", services.OrderTypeTriggerLimit)
	fmt.Printf("  - TRAILING_STOP_MARKET: %s\n", services.OrderTypeTrailingStopMarket)
	fmt.Printf("  - TRAILING_TP_SL: %s\n", services.OrderTypeTrailingTPSL)

	// Example: Create a trigger limit order (test mode)
	testOrder, err := client.Trade().CreateTestOrder(map[string]interface{}{
		"symbol":       "BTC-USDT",
		"side":         "BUY",
		"type":         services.OrderTypeTriggerLimit,
		"positionSide": "LONG",
		"price":        50000.0,
		"stopPrice":    49500.0,
		"quantity":     0.001,
	})

	if err != nil {
		fmt.Printf("  Test order error: %v\n", err)
	} else {
		fmt.Printf("  Test order created: %+v\n", testOrder)
	}

	fmt.Println()
}

func demonstrateTWAPOrders(client *bingx.Client) {
	fmt.Println("2. TWAP Orders (Time-Weighted Average Price)")
	fmt.Println("----------------------------------------------")

	symbol := "BTC-USDT"

	// List existing TWAP orders
	twapOrders, err := client.Trade().GetTWAPOrders(
		&symbol,
		nil, // status
		nil, // startTime
		nil, // endTime
		10,  // limit
		nil, // recvWindow
	)

	if err != nil {
		fmt.Printf("  Error fetching TWAP orders: %v\n", err)
	} else {
		fmt.Printf("  TWAP orders: %+v\n", twapOrders)
	}

	// Example TWAP order parameters
	fmt.Println("\n  Example TWAP order:")
	fmt.Println("  {")
	fmt.Println("    symbol: BTC-USDT,")
	fmt.Println("    side: BUY,")
	fmt.Println("    positionSide: LONG,")
	fmt.Println("    quantity: 1.0,")
	fmt.Println("    duration: 3600,  // Execute over 1 hour")
	fmt.Println("    interval: 60     // Split into 1-minute intervals")
	fmt.Println("  }")

	fmt.Println()
}

func demonstratePositionManagement(client *bingx.Client) {
	fmt.Println("3. Advanced Position Management")
	fmt.Println("--------------------------------")

	symbol := "BTC-USDT"

	// Get position mode
	mode, err := client.Account().GetPositionMode(nil)
	if err != nil {
		fmt.Printf("  Error getting position mode: %v\n", err)
	} else {
		fmt.Printf("  Current position mode: %+v\n", mode)
	}

	// Get position risk
	risk, err := client.Account().GetPositionRisk(&symbol, nil)
	if err != nil {
		fmt.Printf("  Error getting position risk: %v\n", err)
	} else {
		fmt.Printf("  Position risk: %+v\n", risk)
	}

	// One-Click Reverse Position (example - commented out for safety)
	fmt.Println("\n  One-Click Reverse Position:")
	fmt.Println("  // Instantly reverse from LONG to SHORT or vice versa")
	fmt.Println("  // result, err := client.Trade().OneClickReversePosition(\"BTC-USDT\", nil)")

	// Auto Add Margin (example)
	fmt.Println("\n  Auto Add Margin (Hedge Mode):")
	fmt.Println("  // Automatically add margin to positions")
	fmt.Println("  // err := client.Trade().SetAutoAddMargin(\"BTC-USDT\", \"LONG\", true, nil)")

	fmt.Println()
}

func demonstrateMultiAssetsMode(client *bingx.Client) {
	fmt.Println("4. Multi-Assets Mode")
	fmt.Println("--------------------")

	// Get current multi-assets mode
	mode, err := client.Trade().GetMultiAssetsMode(nil)
	if err != nil {
		fmt.Printf("  Error getting multi-assets mode: %v\n", err)
	} else {
		fmt.Printf("  Multi-assets mode: %+v\n", mode)
	}

	// Get multi-assets rules
	rules, err := client.Trade().GetMultiAssetsRules(nil)
	if err != nil {
		fmt.Printf("  Error getting multi-assets rules: %v\n", err)
	} else {
		fmt.Printf("  Multi-assets rules: %+v\n", rules)
	}

	// Get multi-assets margin
	margin, err := client.Trade().GetMultiAssetsMargin(nil)
	if err != nil {
		fmt.Printf("  Error getting multi-assets margin: %v\n", err)
	} else {
		fmt.Printf("  Multi-assets margin: %+v\n", margin)
	}

	fmt.Println("\n  Enable multi-assets mode:")
	fmt.Println("  // err := client.Trade().SwitchMultiAssetsMode(true, nil)")

	fmt.Println()
}

func demonstrateMarketData(client *bingx.Client) {
	fmt.Println("5. Enhanced Market Data (v3)")
	fmt.Println("----------------------------")

	symbol := "BTC-USDT"

	// Get open interest
	oi, err := client.Market().GetOpenInterest(symbol)
	if err != nil {
		fmt.Printf("  Error getting open interest: %v\n", err)
	} else {
		fmt.Printf("  Open interest: %+v\n", oi)
	}

	// Get funding rate info
	fundingInfo, err := client.Market().GetFundingRateInfo(symbol)
	if err != nil {
		fmt.Printf("  Error getting funding rate: %v\n", err)
	} else {
		fmt.Printf("  Funding rate info: %+v\n", fundingInfo)
	}

	// Get book ticker (best bid/ask)
	bookTicker, err := client.Market().GetBookTicker(&symbol)
	if err != nil {
		fmt.Printf("  Error getting book ticker: %v\n", err)
	} else {
		fmt.Printf("  Book ticker: %+v\n", bookTicker)
	}

	// Get index price
	indexPrice, err := client.Market().GetIndexPrice(symbol)
	if err != nil {
		fmt.Printf("  Error getting index price: %v\n", err)
	} else {
		fmt.Printf("  Index price: %+v\n", indexPrice)
	}

	// Get open interest history
	startTime := time.Now().Add(-24 * time.Hour).UnixMilli()
	endTime := time.Now().UnixMilli()
	oiHistory, err := client.Market().GetOpenInterestHistory(
		symbol,
		"5m",
		100,
		&startTime,
		&endTime,
	)
	if err != nil {
		fmt.Printf("  Error getting OI history: %v\n", err)
	} else {
		fmt.Printf("  Open interest history (last 24h): %d records\n", len(oiHistory))
	}

	fmt.Println()
}

func demonstrateAccountFeatures(client *bingx.Client) {
	fmt.Println("6. Advanced Account Features")
	fmt.Println("----------------------------")

	symbol := "BTC-USDT"

	// Get income history
	incomeType := "REALIZED_PNL"
	income, err := client.Account().GetIncomeHistory(
		&symbol,
		&incomeType,
		nil, // startTime
		nil, // endTime
		10,  // limit
		nil, // recvWindow
	)
	if err != nil {
		fmt.Printf("  Error getting income history: %v\n", err)
	} else {
		fmt.Printf("  Income history: %+v\n", income)
	}

	// Get commission history
	commissions, err := client.Account().GetCommissionHistory(
		symbol,
		nil, // startTime
		nil, // endTime
		10,  // limit
		nil, // recvWindow
	)
	if err != nil {
		fmt.Printf("  Error getting commission history: %v\n", err)
	} else {
		fmt.Printf("  Commission history: %+v\n", commissions)
	}

	// Get force orders (liquidations)
	forceOrders, err := client.Account().GetForceOrders(
		&symbol,
		nil, // autoCloseType
		nil, // startTime
		nil, // endTime
		10,  // limit
		nil, // recvWindow
	)
	if err != nil {
		fmt.Printf("  Error getting force orders: %v\n", err)
	} else {
		fmt.Printf("  Force orders: %+v\n", forceOrders)
	}

	fmt.Println("\n  Additional features:")
	fmt.Println("  - Position risk metrics")
	fmt.Println("  - Income/PnL tracking by type")
	fmt.Println("  - Commission history analysis")
	fmt.Println("  - Liquidation order tracking")
	fmt.Println("  - Position mode switching (hedge/one-way)")

	fmt.Println()
}
