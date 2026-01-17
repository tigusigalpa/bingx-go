package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tigusigalpa/bingx-go"
)

func main() {
	apiKey := os.Getenv("BINGX_API_KEY")
	apiSecret := os.Getenv("BINGX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Please set BINGX_API_KEY and BINGX_API_SECRET environment variables")
	}

	client := bingx.NewClient(apiKey, apiSecret)

	listenKeyResp, err := client.ListenKey().Generate()
	if err != nil {
		log.Fatalf("Failed to generate listen key: %v", err)
	}

	listenKey, ok := listenKeyResp["listenKey"].(string)
	if !ok {
		log.Fatal("Failed to get listen key from response")
	}

	stream := client.NewAccountDataStream(listenKey)

	err = stream.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer stream.Disconnect()

	stream.OnAccountUpdate(func(eventType string, data map[string]interface{}) {
		fmt.Printf("Account Update [%s]: %+v\n", eventType, data)
	})

	stream.OnBalanceUpdate(func(balances interface{}) {
		fmt.Printf("Balance Update: %+v\n", balances)
	})

	stream.OnPositionUpdate(func(positions interface{}) {
		fmt.Printf("Position Update: %+v\n", positions)
	})

	stream.OnOrderUpdate(func(order interface{}) {
		fmt.Printf("Order Update: %+v\n", order)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		stream.Stop()
	}()

	fmt.Println("Listening for account data... Press Ctrl+C to stop")
	if err := stream.Listen(); err != nil {
		log.Printf("Listen error: %v", err)
	}
}
