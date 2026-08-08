package main

import (
	"fmt"
	"log"

	bingx "github.com/tigusigalpa/bingx-go/v2"
	"github.com/tigusigalpa/bingx-go/v2/services"
)

func main() {
	client := bingx.NewClient(
		"YOUR_API_KEY",
		"YOUR_API_SECRET",
		bingx.WithBaseURI("https://open-api.bingx.com"),
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

	// Get an overview of all wallet types (spot, standard futures, USDT-M, etc.)
	overview, err := client.SpotAccount().GetAccountOverview(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All Account Overview: %v\n", overview)

	// Get only the USDT-M perpetual wallet overview
	usdtM := services.AccountTypeUSDTMPerp
	usdtOverview, err := client.SpotAccount().GetAccountOverview(&usdtM)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT-M Perp Overview: %v\n", usdtOverview)

	symbols, err := client.Market().GetFuturesSymbols()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Available Symbols: %v\n", symbols)
}
