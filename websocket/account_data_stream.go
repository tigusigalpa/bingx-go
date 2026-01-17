package websocket

import (
	"fmt"
)

const AccountDataStreamBaseURL = "wss://open-api-swap.bingx.com/swap-market"

type AccountDataStream struct {
	*WebSocketClient
}

func NewAccountDataStream(listenKey string) *AccountDataStream {
	url := fmt.Sprintf("%s?listenKey=%s", AccountDataStreamBaseURL, listenKey)
	return &AccountDataStream{
		WebSocketClient: NewWebSocketClient(url),
	}
}

type AccountUpdateCallback func(eventType string, data map[string]interface{})

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

type BalanceUpdateCallback func(balances interface{})

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

type PositionUpdateCallback func(positions interface{})

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

type OrderUpdateCallback func(order interface{})

func (a *AccountDataStream) OnOrderUpdate(callback OrderUpdateCallback) {
	a.OnMessage(func(data map[string]interface{}) {
		if eventType, ok := data["e"].(string); ok && eventType == "ORDER_TRADE_UPDATE" {
			if order, ok := data["o"]; ok {
				callback(order)
			}
		}
	})
}
