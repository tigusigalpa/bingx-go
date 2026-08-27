package coinm

import (
	"time"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

// TradeService represents a BingX API component or value.
type TradeService struct {
	client *http.BaseHTTPClient
}

// NewTradeService creates a new client or service instance.
func NewTradeService(client *http.BaseHTTPClient) *TradeService {
	return &TradeService{client: client}
}

// CreateOrder performs the CreateOrder operation.
func (s *TradeService) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/cswap/v1/trade/order", params)
}

// CancelOrder performs the CancelOrder operation.
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

	return s.client.Request("DELETE", "/openApi/cswap/v1/trade/cancelOrder", params)
}

// CancelAllOrders performs the CancelAllOrders operation.
func (s *TradeService) CancelAllOrders(symbol string) (map[string]interface{}, error) {
	return s.client.Request("DELETE", "/openApi/cswap/v1/trade/allOpenOrders", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetOrder performs the GetOrder operation.
func (s *TradeService) GetOrder(symbol, orderID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/trade/orderDetail", map[string]interface{}{
		"symbol":  symbol,
		"orderId": orderID,
	})
}

// GetOpenOrders performs the GetOpenOrders operation.
func (s *TradeService) GetOpenOrders(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/cswap/v1/trade/openOrders", params)
}

// GetPositions performs the GetPositions operation.
func (s *TradeService) GetPositions(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/cswap/v1/user/positions", params)
}

// GetBalance performs the GetBalance operation.
func (s *TradeService) GetBalance() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/user/balance", nil)
}

// GetLeverage performs the GetLeverage operation.
func (s *TradeService) GetLeverage(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/trade/leverage", map[string]interface{}{
		"symbol": symbol,
	})
}

// SetLeverage performs the SetLeverage operation.
func (s *TradeService) SetLeverage(symbol, side string, leverage int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/cswap/v1/trade/leverage", map[string]interface{}{
		"symbol":    symbol,
		"side":      side,
		"leverage":  leverage,
		"timestamp": time.Now().UnixMilli(),
	})
}

// GetMarginType performs the GetMarginType operation.
func (s *TradeService) GetMarginType(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/trade/marginType", map[string]interface{}{
		"symbol": symbol,
	})
}

// SetMarginType performs the SetMarginType operation.
func (s *TradeService) SetMarginType(symbol, marginType string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/cswap/v1/trade/marginType", map[string]interface{}{
		"symbol":     symbol,
		"marginType": marginType,
		"timestamp":  time.Now().UnixMilli(),
	})
}

// SetPositionMargin performs the SetPositionMargin operation.
func (s *TradeService) SetPositionMargin(symbol, positionSide string, amount float64, marginType int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/cswap/v1/trade/positionMargin", map[string]interface{}{
		"symbol":       symbol,
		"positionSide": positionSide,
		"amount":       amount,
		"type":         marginType,
	})
}

// GetOrderHistory performs the GetOrderHistory operation.
func (s *TradeService) GetOrderHistory(symbol string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/cswap/v1/trade/orderHistory", params)
}

// GetUserTrades performs the GetUserTrades operation.
func (s *TradeService) GetUserTrades(symbol string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/cswap/v1/trade/allFillOrders", params)
}

// GetPositionRisk performs the GetPositionRisk operation.
func (s *TradeService) GetPositionRisk(symbol *string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/cswap/v1/user/positions", params)
}

// GetIncomeHistory performs the GetIncomeHistory operation.
func (s *TradeService) GetIncomeHistory(symbol *string, incomeType *string, startTime, endTime *int64, limit int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}
	if incomeType != nil {
		params["incomeType"] = *incomeType
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/cswap/v1/user/income", params)
}
