package tradfi

import (
	"time"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

type TradeService struct {
	client *http.BaseHTTPClient
}

func NewTradeService(client *http.BaseHTTPClient) *TradeService {
	return &TradeService{client: client}
}

const (
	OrderTypeMarket             = "MARKET"
	OrderTypeLimit              = "LIMIT"
	OrderTypeStop               = "STOP"
	OrderTypeStopMarket         = "STOP_MARKET"
	OrderTypeTakeProfit         = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket   = "TAKE_PROFIT_MARKET"
	OrderTypeTriggerLimit       = "TRIGGER_LIMIT"
	OrderTypeTrailingStopMarket = "TRAILING_STOP_MARKET"
	OrderTypeTrailingTPSL       = "TRAILING_TP_SL"
)

// CreateOrder places a new order for TradFi instruments
func (s *TradeService) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/order", params)
}

// CreateTestOrder validates an order without executing it
func (s *TradeService) CreateTestOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/order/test", params)
}

// CancelOrder cancels an order by ID or client order ID
func (s *TradeService) CancelOrder(symbol string, orderID *string, clientOrderID *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if orderID != nil {
		params["orderId"] = *orderID
	}
	if clientOrderID != nil {
		params["clientOrderId"] = *clientOrderID
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/order", params)
}

// CancelAllOrders cancels all open orders for a symbol
func (s *TradeService) CancelAllOrders(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/allOpenOrders", params)
}

// GetOrder retrieves order details
func (s *TradeService) GetOrder(symbol, orderID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/trade/orderDetail", map[string]interface{}{
		"symbol":  symbol,
		"orderId": orderID,
	})
}

// GetOpenOrders retrieves all open orders
func (s *TradeService) GetOpenOrders(symbol *string, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/openOrders", params)
}

// GetOrderHistory retrieves historical orders
func (s *TradeService) GetOrderHistory(symbol *string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/orderHistory", params)
}

// GetUserTrades retrieves user trade history
func (s *TradeService) GetUserTrades(symbol *string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/fillHistory", params)
}

// SetLeverage sets leverage for a TradFi symbol
func (s *TradeService) SetLeverage(symbol string, leverage int, side *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"leverage":  leverage,
		"timestamp": time.Now().UnixMilli(),
	}
	if side != nil {
		params["side"] = *side
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/leverage", params)
}

// GetLeverage retrieves current leverage for a TradFi symbol
func (s *TradeService) GetLeverage(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/trade/leverage", map[string]interface{}{
		"symbol": symbol,
	})
}

// SetMarginType sets margin type (ISOLATED or CROSSED)
func (s *TradeService) SetMarginType(symbol, marginType string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/marginType", map[string]interface{}{
		"symbol":     symbol,
		"marginType": marginType,
		"timestamp":  time.Now().UnixMilli(),
	})
}

// GetMarginType retrieves current margin type
func (s *TradeService) GetMarginType(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/trade/marginType", map[string]interface{}{
		"symbol": symbol,
	})
}

// SetPositionMargin adds or reduces position margin
func (s *TradeService) SetPositionMargin(symbol, positionSide string, amount float64, marginType int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/positionMargin", map[string]interface{}{
		"symbol":       symbol,
		"positionSide": positionSide,
		"amount":       amount,
		"type":         marginType,
	})
}

// OneClickReversePosition reverses position side (LONG <-> SHORT)
func (s *TradeService) OneClickReversePosition(symbol string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/oneClickReversePosition", map[string]interface{}{
		"symbol": symbol,
	})
}

// PlaceTWAPOrder places a TWAP (Time-Weighted Average Price) order
func (s *TradeService) PlaceTWAPOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/twap/order", params)
}

// GetTWAPOrder retrieves TWAP order details
func (s *TradeService) GetTWAPOrder(orderID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/trade/twap/order", map[string]interface{}{
		"orderId": orderID,
	})
}

// GetTWAPOrders retrieves all TWAP orders
func (s *TradeService) GetTWAPOrders(symbol *string, status *string, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if status != nil {
		params["status"] = *status
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/twap/orders", params)
}

// CancelTWAPOrder cancels a TWAP order
func (s *TradeService) CancelTWAPOrder(orderID string) (map[string]interface{}, error) {
	return s.client.Request("DELETE", "/openApi/swap/v2/trade/twap/order", map[string]interface{}{
		"orderId": orderID,
	})
}

// ModifyOrder modifies an existing order's price and/or quantity
func (s *TradeService) ModifyOrder(symbol string, orderID *string, clientOrderID *string, price, quantity float64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":   symbol,
		"price":    price,
		"quantity": quantity,
	}
	if orderID != nil {
		params["orderId"] = *orderID
	}
	if clientOrderID != nil {
		params["clientOrderId"] = *clientOrderID
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/order", params)
}

// SetAutoAddMargin enables/disables auto margin addition
func (s *TradeService) SetAutoAddMargin(symbol, positionSide string, enabled bool) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/autoAddMargin", map[string]interface{}{
		"symbol":       symbol,
		"positionSide": positionSide,
		"enabled":      enabled,
	})
}
