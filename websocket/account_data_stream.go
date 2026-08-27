package websocket

import (
	"fmt"
	"net/url"
)

// AccountDataStreamBaseURL is the WebSocket endpoint for account data streams (v3 compatible)
const AccountDataStreamBaseURL = "wss://open-api-swap.bingx.com/swap-market"

// AccountDataStream represents a BingX API component or value.
type AccountDataStream struct {
	*Client
}

// NewAccountDataStream creates a new client or service instance.
func NewAccountDataStream(listenKey string) *AccountDataStream {
	endpoint := fmt.Sprintf("%s?listenKey=%s", AccountDataStreamBaseURL, url.QueryEscape(listenKey))
	return &AccountDataStream{
		Client: NewClient(endpoint),
	}
}

// AccountUpdateCallback represents a BingX API component or value.
type AccountUpdateCallback func(eventType string, data map[string]interface{})

// OnAccountUpdate performs the OnAccountUpdate operation.
func (a *AccountDataStream) OnAccountUpdate(callback AccountUpdateCallback) {
	a.OnMessage(func(data map[string]interface{}) {
		if eventType, ok := data["e"].(string); ok {
			switch eventType {
			case "ACCOUNT_UPDATE":
				callback("account", data)
			case "ORDER_TRADE_UPDATE":
				callback("order", data)
			default:
				callback("unknown", data)
			}
		}
	})
}

// BalanceUpdateCallback represents a BingX API component or value.
type BalanceUpdateCallback func(balances interface{})

// OnBalanceUpdate performs the OnBalanceUpdate operation.
func (a *AccountDataStream) OnBalanceUpdate(callback BalanceUpdateCallback) {
	a.OnMessage(func(data map[string]interface{}) {
		if eventType, ok := data["e"].(string); ok && eventType == "ACCOUNT_UPDATE" {
			if accountData, ok := data["a"].(map[string]interface{}); ok {
				if balances, ok := accountData["B"]; ok {
					callback(balances)
				}
			}
		}
	})
}

// PositionUpdateCallback represents a BingX API component or value.
type PositionUpdateCallback func(positions interface{})

// OnPositionUpdate performs the OnPositionUpdate operation.
func (a *AccountDataStream) OnPositionUpdate(callback PositionUpdateCallback) {
	a.OnMessage(func(data map[string]interface{}) {
		if eventType, ok := data["e"].(string); ok && eventType == "ACCOUNT_UPDATE" {
			if accountData, ok := data["a"].(map[string]interface{}); ok {
				if positions, ok := accountData["P"]; ok {
					callback(positions)
				}
			}
		}
	})
}

// OrderUpdateCallback represents a BingX API component or value.
type OrderUpdateCallback func(order interface{})

// OnOrderUpdate performs the OnOrderUpdate operation.
func (a *AccountDataStream) OnOrderUpdate(callback OrderUpdateCallback) {
	a.OnMessage(func(data map[string]interface{}) {
		if eventType, ok := data["e"].(string); ok && eventType == "ORDER_TRADE_UPDATE" {
			if order, ok := data["o"]; ok {
				callback(order)
			}
		}
	})
}
