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
		bingx.WithBaseURI("https://open-api.bingx.com"),
		bingx.WithSignatureEncoding("base64"),
	)

	price, err := client.Market().GetLatestPrice("BTC-USDT")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("BTC-USDT Price: %v\n", price)

	balance, err := client.Account().GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account Balance: %v\n", balance)

	symbols, err := client.Market().GetFuturesSymbols()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Available Symbols: %v\n", symbols)
}
