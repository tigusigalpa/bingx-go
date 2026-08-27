package services

import (
	"errors"
	"time"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

// OrderTypeMarket defines a BingX API constant.
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

// TradeService represents a BingX API component or value.
type TradeService struct {
	client *http.BaseHTTPClient
}

// FuturesCommissionRate defines a BingX API constant.
const FuturesCommissionRate = 0.00045

// NewTradeService creates a new client or service instance.
func NewTradeService(client *http.BaseHTTPClient) *TradeService {
	return &TradeService{client: client}
}

// CommissionResult represents a BingX API component or value.
type CommissionResult struct {
	Margin                float64 `json:"margin"`
	Leverage              int     `json:"leverage"`
	PositionValue         float64 `json:"position_value"`
	CommissionRate        float64 `json:"commission_rate"`
	CommissionRatePercent float64 `json:"commission_rate_percent"`
	Commission            float64 `json:"commission"`
	CommissionRounded     float64 `json:"commission_rounded"`
	NetPositionValue      float64 `json:"net_position_value"`
}

// CalculateFuturesCommission performs the CalculateFuturesCommission operation.
func (s *TradeService) CalculateFuturesCommission(margin float64, leverage int, commissionRate *float64) CommissionResult {
	rate := FuturesCommissionRate
	if commissionRate != nil {
		rate = *commissionRate
	}

	positionValue := margin * float64(leverage)
	commission := positionValue * rate

	return CommissionResult{
		Margin:                margin,
		Leverage:              leverage,
		PositionValue:         positionValue,
		CommissionRate:        rate,
		CommissionRatePercent: rate * 100,
		Commission:            commission,
		CommissionRounded:     float64(int(commission*1000000)) / 1000000,
		NetPositionValue:      positionValue - commission,
	}
}

// GetCommissionAmount performs the GetCommissionAmount operation.
func (s *TradeService) GetCommissionAmount(margin float64, leverage int) float64 {
	return margin * float64(leverage) * FuturesCommissionRate
}

// CreateOrder performs the CreateOrder operation.
func (s *TradeService) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/order", params)
}

// ModifyOrder performs the ModifyOrder operation.
func (s *TradeService) ModifyOrder(symbol string, quantity float64, orderID, clientOrderID *string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	if orderID == nil && clientOrderID == nil {
		return nil, errors.New("modifyOrder requires either orderID or clientOrderID")
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"quantity": quantity,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if orderID != nil {
		params["orderId"] = *orderID
	}

	if clientOrderID != nil {
		params["clientOrderId"] = *clientOrderID
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v1/trade/amend", params)
}

// CreateTestOrder performs the CreateTestOrder operation.
func (s *TradeService) CreateTestOrder(params map[string]interface{}) (map[string]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	if _, exists := params["timestamp"]; !exists {
		params["timestamp"] = time.Now().UnixMilli()
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/order/test", params)
}

// CloseAllPositions performs the CloseAllPositions operation.
func (s *TradeService) CloseAllPositions(symbol string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/closeAllPositions", params)
}

// GetMarginType performs the GetMarginType operation.
func (s *TradeService) GetMarginType(symbol string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/marginType", params)
}

// ChangeMarginType performs the ChangeMarginType operation.
func (s *TradeService) ChangeMarginType(symbol, marginType string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	if marginType != "ISOLATED" && marginType != "CROSSED" {
		return nil, errors.New("margin type must be ISOLATED or CROSSED")
	}

	params := map[string]interface{}{
		"symbol":     symbol,
		"marginType": marginType,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/marginType", params)
}

// GetLeverage performs the GetLeverage operation.
func (s *TradeService) GetLeverage(symbol string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/leverage", params)
}

// SetLeverage performs the SetLeverage operation.
func (s *TradeService) SetLeverage(symbol string, leverage int, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	if leverage < 1 || leverage > 125 {
		return nil, errors.New("leverage must be between 1 and 125")
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"leverage": leverage,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/leverage", params)
}

// CreateBatchOrders performs the CreateBatchOrders operation.
func (s *TradeService) CreateBatchOrders(orders []map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/trade/batchOrders", map[string]interface{}{
		"orders": orders,
	})
}

// CancelOrder performs the CancelOrder operation.
func (s *TradeService) CancelOrder(symbol string, orderID, clientOrderID *string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if orderID != nil {
		params["orderId"] = *orderID
	}

	if clientOrderID != nil {
		params["clientOrderId"] = *clientOrderID
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/order", params)
}

// CancelAllOrders performs the CancelAllOrders operation.
func (s *TradeService) CancelAllOrders(timestamp *int64, symbol, orderType *string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}

	if orderType != nil {
		params["type"] = *orderType
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/allOpenOrders", params)
}

// CancelBatchOrders performs the CancelBatchOrders operation.
func (s *TradeService) CancelBatchOrders(symbol string, orderIDs []string, clientOrderIDs []string, timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
	}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if len(orderIDs) > 0 {
		orderIDsStr := ""
		for i, id := range orderIDs {
			if i > 0 {
				orderIDsStr += ","
			}
			orderIDsStr += id
		}
		params["orderIds"] = orderIDsStr
	}

	if len(clientOrderIDs) > 0 {
		clientIDsStr := ""
		for i, id := range clientOrderIDs {
			if i > 0 {
				clientIDsStr += ","
			}
			clientIDsStr += id
		}
		params["clientOrderIds"] = clientIDsStr
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/batchOrders", params)
}

// GetOrder performs the GetOrder operation.
func (s *TradeService) GetOrder(symbol, orderID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/trade/order", map[string]interface{}{
		"symbol":  symbol,
		"orderId": orderID,
	})
}

// GetOpenOrders performs the GetOpenOrders operation.
func (s *TradeService) GetOpenOrders(symbol *string, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/openOrders", params)
}

// GetOrderHistory performs the GetOrderHistory operation.
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

// GetFilledOrders performs the GetFilledOrders operation.
func (s *TradeService) GetFilledOrders(symbol *string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/trade/filledOrders", params)
}

// GetUserTrades performs the GetUserTrades operation.
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

	return s.client.Request("GET", "/openApi/swap/v2/trade/userTrades", params)
}

// ChangeLeverage performs the ChangeLeverage operation.
func (s *TradeService) ChangeLeverage(symbol, side string, leverage int, recvWindow *int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"side":      side,
		"leverage":  leverage,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/leverage", params)
}

// OneClickReversePosition performs the OneClickReversePosition operation.
func (s *TradeService) OneClickReversePosition(symbol string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/oneClickReversePosition", params)
}

// SetAutoAddMargin performs the SetAutoAddMargin operation.
func (s *TradeService) SetAutoAddMargin(symbol, positionSide string, autoAddMargin bool, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":        symbol,
		"positionSide":  positionSide,
		"autoAddMargin": autoAddMargin,
		"timestamp":     time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/autoAddMargin", params)
}

// SwitchMultiAssetsMode performs the SwitchMultiAssetsMode operation.
func (s *TradeService) SwitchMultiAssetsMode(multiAssetsMargin bool, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"multiAssetsMargin": multiAssetsMargin,
		"timestamp":         time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/multiAssetsMode", params)
}

// GetMultiAssetsMode performs the GetMultiAssetsMode operation.
func (s *TradeService) GetMultiAssetsMode(recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/multiAssetsMode", params)
}

// GetMultiAssetsRules performs the GetMultiAssetsRules operation.
func (s *TradeService) GetMultiAssetsRules(recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/multiAssetsRules", params)
}

// GetMultiAssetsMargin performs the GetMultiAssetsMargin operation.
func (s *TradeService) GetMultiAssetsMargin(recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/multiAssetsMargin", params)
}

// PlaceTWAPOrder performs the PlaceTWAPOrder operation.
func (s *TradeService) PlaceTWAPOrder(params map[string]interface{}) (map[string]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	if _, exists := params["timestamp"]; !exists {
		params["timestamp"] = time.Now().UnixMilli()
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/twapOrder", params)
}

// CancelTWAPOrder performs the CancelTWAPOrder operation.
func (s *TradeService) CancelTWAPOrder(orderID string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"orderId":   orderID,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("DELETE", "/openApi/swap/v2/trade/twapOrder", params)
}

// GetTWAPOrder performs the GetTWAPOrder operation.
func (s *TradeService) GetTWAPOrder(orderID string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"orderId":   orderID,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/twapOrder", params)
}

// GetTWAPOrders performs the GetTWAPOrders operation.
func (s *TradeService) GetTWAPOrders(symbol *string, status *string, startTime, endTime *int64, limit int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}
	if status != nil {
		params["status"] = *status
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

	return s.client.Request("GET", "/openApi/swap/v2/trade/twapOrders", params)
}

// GetVst retrieves VST (Virtual Simulation Trading) information
// This enables demo trading functionality in the simulated environment
func (s *TradeService) GetVst(recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/getVst", params)
}
