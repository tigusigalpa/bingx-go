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
	client := bingx.NewClient("", "")

	stream := client.NewMarketDataStream()

	err := stream.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer stream.Disconnect()

	stream.OnMessage(func(data map[string]interface{}) {
		fmt.Printf("Received message: %+v\n", data)
	})

	err = stream.SubscribeTrade("BTC-USDT")
	if err != nil {
		log.Fatalf("Failed to subscribe to trade: %v", err)
	}

	err = stream.SubscribeKline("BTC-USDT", "1m")
	if err != nil {
		log.Fatalf("Failed to subscribe to kline: %v", err)
	}

	err = stream.SubscribeDepth("BTC-USDT", 20)
	if err != nil {
		log.Fatalf("Failed to subscribe to depth: %v", err)
	}

	err = stream.SubscribeTicker("BTC-USDT")
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	err = stream.SubscribeBookTicker("BTC-USDT")
	if err != nil {
		log.Fatalf("Failed to subscribe to book ticker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		stream.Stop()
	}()

	fmt.Println("Listening for market data... Press Ctrl+C to stop")
	if err := stream.Listen(); err != nil {
		log.Printf("Listen error: %v", err)
	}
}
