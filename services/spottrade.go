package services

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

// Spot order side, type, and time-in-force constants.
//
// Only LIMIT and MARKET are exposed as constants because those are the two
// order types supported by the typed SpotOrderRequest / CreateOrderRequest
// convenience API. The BingX Spot API also supports TAKE_STOP_LIMIT,
// TAKE_STOP_MARKET, TRIGGER_LIMIT, and TRIGGER_MARKET order types; use
// CreateOrder with a raw map for those.
const (
	SpotOrderTypeLimit  = "LIMIT"
	SpotOrderTypeMarket = "MARKET"

	SpotSideBuy  = "BUY"
	SpotSideSell = "SELL"

	SpotTimeInForceGTC      = "GTC"
	SpotTimeInForceIOC      = "IOC"
	SpotTimeInForceFOK      = "FOK"
	SpotTimeInForcePostOnly = "PostOnly"
)

// Cancel-replace modes for AmendOrder.
// See POST /openApi/spot/v1/trade/order/cancelReplace.
const (
	SpotCancelReplaceStopOnFailure = "STOP_ON_FAILURE"
	SpotCancelReplaceAllowFailure  = "ALLOW_FAILURE"
)

// SpotOrderRequest is a typed convenience wrapper for placing spot LIMIT or
// MARKET orders. For advanced order types (TAKE_STOP_LIMIT,
// TAKE_STOP_MARKET, TRIGGER_LIMIT, TRIGGER_MARKET) or exchange-specific
// parameters, use CreateOrder with a raw map instead.
//
// Quantity, Price, and QuoteOrderQty are strings by design: BingX expects
// decimal values as exact strings on the wire. Representing them as float64
// risks silent precision loss (e.g. binary floating point rounding) that can
// lead to rejected orders or incorrect fills. Callers should format decimals
// themselves (e.g. with a decimal library) rather than using
// strconv.FormatFloat.
type SpotOrderRequest struct {
	Symbol        string
	Side          string
	Type          string
	Quantity      string
	Price         *string
	QuoteOrderQty *string
	TimeInForce   *string
	ClientOrderID *string
}

// SpotTradeService provides access to the BingX Spot Trading REST API
// (https://openApi/spot/v1/trade/*). It reuses the shared, already-signed
// http.BaseHTTPClient used by the other services, so no HMAC signing or
// error handling logic is duplicated here.
type SpotTradeService struct {
	client *http.BaseHTTPClient
}

// NewSpotTradeService creates a new SpotTradeService bound to the given
// signed HTTP client.
func NewSpotTradeService(client *http.BaseHTTPClient) *SpotTradeService {
	return &SpotTradeService{client: client}
}

// CreateOrder places a spot order using raw exchange parameters. Use this
// for advanced order types not covered by CreateOrderRequest, or to pass
// exchange-specific parameters not modeled by SpotOrderRequest.
//
// POST /openApi/spot/v1/trade/order
func (s *SpotTradeService) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/spot/v1/trade/order", params)
}

// CreateOrderRequest places a typed LIMIT or MARKET spot order after
// validating the request client-side. See SpotOrderRequest for field
// documentation.
func (s *SpotTradeService) CreateOrderRequest(req SpotOrderRequest) (map[string]interface{}, error) {
	if err := validateSpotOrderRequest(req); err != nil {
		return nil, err
	}

	return s.CreateOrder(spotOrderRequestToParams(req))
}

// CreateTestOrder is intentionally not implemented.
//
// Unlike the Futures API (POST /openApi/swap/v2/trade/order/test), the
// BingX Spot API does not expose a dry-run "test order" endpoint that
// validates signature/parameters without executing a trade. Implementing
// this method would either silently call the real order-placement endpoint
// (misleading callers into thinking no trade occurred) or fabricate an
// endpoint that does not exist. Validate orders client-side via
// CreateOrderRequest (which performs the same checks the exchange would
// reject on) instead.

// CancelOrder cancels a single spot order. Exactly one of orderID or
// clientOrderID must be provided to identify the order.
//
// POST /openApi/spot/v1/trade/cancel
func (s *SpotTradeService) CancelOrder(symbol string, orderID, clientOrderID *string) (map[string]interface{}, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	if orderID == nil && clientOrderID == nil {
		return nil, errors.New("cancelOrder requires either orderID or clientOrderID")
	}

	params := map[string]interface{}{
		"symbol": symbol,
	}
	if orderID != nil {
		params["orderId"] = *orderID
	}
	if clientOrderID != nil {
		params["clientOrderID"] = *clientOrderID
	}

	return s.client.Request("POST", "/openApi/spot/v1/trade/cancel", params)
}

// CancelAllOrders cancels every open spot order for the given symbol. If
// symbol is nil, all open orders across all trading pairs are cancelled.
//
// POST /openApi/spot/v1/trade/cancelOpenOrders
func (s *SpotTradeService) CancelAllOrders(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("POST", "/openApi/spot/v1/trade/cancelOpenOrders", params)
}

// CancelBatchOrders cancels a batch of spot orders identified by orderIDs
// and/or clientOrderIDs. At least one ID must be provided.
//
// POST /openApi/spot/v1/trade/cancelOrders
func (s *SpotTradeService) CancelBatchOrders(symbol string, orderIDs, clientOrderIDs []string) (map[string]interface{}, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	if len(orderIDs) == 0 && len(clientOrderIDs) == 0 {
		return nil, errors.New("cancelBatchOrders requires at least one orderID or clientOrderID")
	}

	params := map[string]interface{}{
		"symbol": symbol,
	}
	if len(orderIDs) > 0 {
		params["orderIds"] = strings.Join(orderIDs, ",")
	}
	if len(clientOrderIDs) > 0 {
		params["clientOrderIDs"] = strings.Join(clientOrderIDs, ",")
	}

	return s.client.Request("POST", "/openApi/spot/v1/trade/cancelOrders", params)
}

// AmendOrder amends an existing spot order using BingX's native
// cancel-and-replace endpoint.
//
// BingX Spot does not support true in-place order modification the way
// Futures does (POST /openApi/swap/v1/trade/amend, which only changes the
// quantity of an existing order). Instead, the Spot API exposes
// POST /openApi/spot/v1/trade/order/cancelReplace, which atomically cancels
// the specified order and places a brand new one in a single signed
// request. This is safer than a manual CancelOrder + CreateOrder sequence
// (which has a window where neither or both orders could be resting on the
// book) and is therefore surfaced here as AmendOrder.
//
// Exactly one of cancelOrderID or cancelClientOrderID must identify the
// order to replace. cancelReplaceMode controls what happens if the cancel
// step fails: SpotCancelReplaceStopOnFailure aborts without placing the new
// order, while SpotCancelReplaceAllowFailure places the new order
// regardless of whether the cancel succeeded.
//
// POST /openApi/spot/v1/trade/order/cancelReplace
func (s *SpotTradeService) AmendOrder(symbol string, cancelOrderID, cancelClientOrderID *string, cancelReplaceMode string, newOrder SpotOrderRequest) (map[string]interface{}, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	if cancelOrderID == nil && cancelClientOrderID == nil {
		return nil, errors.New("amendOrder requires either cancelOrderID or cancelClientOrderID")
	}
	if cancelReplaceMode != SpotCancelReplaceStopOnFailure && cancelReplaceMode != SpotCancelReplaceAllowFailure {
		return nil, errors.New("cancelReplaceMode must be STOP_ON_FAILURE or ALLOW_FAILURE")
	}

	if newOrder.Symbol == "" {
		newOrder.Symbol = symbol
	}
	if err := validateSpotOrderRequest(newOrder); err != nil {
		return nil, err
	}

	params := spotOrderRequestToParams(newOrder)
	params["cancelReplaceMode"] = cancelReplaceMode
	if cancelOrderID != nil {
		params["cancelOrderId"] = *cancelOrderID
	}
	if cancelClientOrderID != nil {
		params["cancelClientOrderID"] = *cancelClientOrderID
	}

	return s.client.Request("POST", "/openApi/spot/v1/trade/order/cancelReplace", params)
}

// GetOrder queries a single spot order. Exactly one of orderID or
// clientOrderID must be provided to identify the order.
//
// GET /openApi/spot/v1/trade/query
func (s *SpotTradeService) GetOrder(symbol string, orderID, clientOrderID *string) (map[string]interface{}, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	if orderID == nil && clientOrderID == nil {
		return nil, errors.New("getOrder requires either orderID or clientOrderID")
	}

	params := map[string]interface{}{
		"symbol": symbol,
	}
	if orderID != nil {
		params["orderId"] = *orderID
	}
	if clientOrderID != nil {
		params["clientOrderID"] = *clientOrderID
	}

	return s.client.Request("GET", "/openApi/spot/v1/trade/query", params)
}

// GetOpenOrders returns current open spot orders. If symbol is nil, open
// orders across all trading pairs are returned.
//
// Note: unlike the Futures TradeService.GetOpenOrders, the BingX Spot
// "openOrders" endpoint does not accept a limit/pageSize parameter, so this
// method intentionally has no limit argument.
//
// GET /openApi/spot/v1/trade/openOrders
func (s *SpotTradeService) GetOpenOrders(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/spot/v1/trade/openOrders", params)
}

// GetOrderHistory returns historical spot orders. limit is sent as the
// API's pageSize parameter (max 100) if greater than zero.
//
// GET /openApi/spot/v1/trade/historyOrders
func (s *SpotTradeService) GetOrderHistory(symbol *string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if limit > 0 {
		params["pageSize"] = limit
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/spot/v1/trade/historyOrders", params)
}

// GetTrades returns filled trade details for a symbol. The BingX Spot API
// requires symbol for this endpoint (unlike GetOpenOrders/GetOrderHistory),
// so symbol must be non-nil and non-empty.
//
// GET /openApi/spot/v1/trade/myTrades
func (s *SpotTradeService) GetTrades(symbol *string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	if symbol == nil || *symbol == "" {
		return nil, errors.New("symbol is required")
	}

	params := map[string]interface{}{
		"symbol": *symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/spot/v1/trade/myTrades", params)
}

// spotOrderRequestToParams converts a validated SpotOrderRequest into raw
// exchange parameters, omitting fields that were left unset.
func spotOrderRequestToParams(req SpotOrderRequest) map[string]interface{} {
	params := map[string]interface{}{
		"symbol": req.Symbol,
		"side":   req.Side,
		"type":   req.Type,
	}

	if req.Quantity != "" {
		params["quantity"] = req.Quantity
	}
	if req.Price != nil {
		params["price"] = *req.Price
	}
	if req.QuoteOrderQty != nil {
		params["quoteOrderQty"] = *req.QuoteOrderQty
	}
	if req.TimeInForce != nil {
		params["timeInForce"] = *req.TimeInForce
	}
	if req.ClientOrderID != nil {
		params["newClientOrderId"] = *req.ClientOrderID
	}

	return params
}

// validateSpotOrderRequest validates a SpotOrderRequest client-side before
// it is sent to the exchange, per the BingX Spot API's order rules:
//   - symbol, side, and type are mandatory.
//   - side must be BUY or SELL.
//   - type must be LIMIT or MARKET (use CreateOrder for other types).
//   - a LIMIT order requires a positive price and a positive quantity.
//   - a MARKET order requires a positive quantity or quoteOrderQty.
func validateSpotOrderRequest(req SpotOrderRequest) error {
	if req.Symbol == "" {
		return errors.New("symbol is required")
	}
	if req.Side != SpotSideBuy && req.Side != SpotSideSell {
		return errors.New("side must be BUY or SELL")
	}
	if req.Type != SpotOrderTypeLimit && req.Type != SpotOrderTypeMarket {
		return errors.New("type must be LIMIT or MARKET")
	}

	switch req.Type {
	case SpotOrderTypeLimit:
		if req.Price == nil || !isPositiveDecimalString(*req.Price) {
			return errors.New("a LIMIT order requires a positive price")
		}
		if !isPositiveDecimalString(req.Quantity) {
			return errors.New("a LIMIT order requires a positive quantity")
		}
	case SpotOrderTypeMarket:
		hasQuantity := isPositiveDecimalString(req.Quantity)
		hasQuoteQty := req.QuoteOrderQty != nil && isPositiveDecimalString(*req.QuoteOrderQty)
		if !hasQuantity && !hasQuoteQty {
			return errors.New("a MARKET order requires a positive quantity or quoteOrderQty")
		}
	}

	return nil
}

// isPositiveDecimalString reports whether v parses as a decimal number
// greater than zero. It is used only for client-side validation; the
// original string is always what gets sent on the wire, never a
// float64-round-tripped value.
func isPositiveDecimalString(v string) bool {
	if v == "" {
		return false
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false
	}
	return f > 0
}
