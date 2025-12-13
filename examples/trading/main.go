package main

import (
	"fmt"
	"log"

	bingx "github.com/tigusigalpa/bingx-go"
)

func main() {
	client := bingx.NewClient(
		"YOUR_API_KEY",
		"YOUR_API_SECRET",
		bingx.WithSignatureEncoding("base64"),
	)

	order, err := client.Trade().CreateOrder(map[string]interface{}{
		"symbol":       "BTC-USDT",
		"side":         "BUY",
		"type":         "LIMIT",
		"positionSide": "LONG",
		"price":        50000.0,
		"quantity":     0.001,
		"timeInForce":  "GTC",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order Created: %v\n", order)

	openOrders, err := client.Trade().GetOpenOrders(nil, 100)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Open Orders: %v\n", openOrders)

	positions, err := client.Account().GetPositions(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Positions: %v\n", positions)
}
