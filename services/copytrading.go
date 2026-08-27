package services

import "github.com/tigusigalpa/bingx-go/v2/http"

// CopyTradingService represents a BingX API component or value.
type CopyTradingService struct {
	client *http.BaseHTTPClient
}

// NewCopyTradingService creates a new client or service instance.
func NewCopyTradingService(client *http.BaseHTTPClient) *CopyTradingService {
	return &CopyTradingService{client: client}
}

// GetCurrentTrackOrders performs the GetCurrentTrackOrders operation.
func (s *CopyTradingService) GetCurrentTrackOrders(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/currentTrack", map[string]interface{}{
		"symbol": symbol,
	})
}

// CloseTrackOrder performs the CloseTrackOrder operation.
func (s *CopyTradingService) CloseTrackOrder(orderNumber string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/trader/closeTrack", map[string]interface{}{
		"orderNumber": orderNumber,
	})
}

// SetTPSL performs the SetTPSL operation.
func (s *CopyTradingService) SetTPSL(positionID string, stopLoss, takeProfit *float64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"positionId": positionID,
	}

	if stopLoss != nil {
		params["stopLoss"] = *stopLoss
	}
	if takeProfit != nil {
		params["takeProfit"] = *takeProfit
	}

	return s.client.Request("POST", "/openApi/copy/v1/trader/setTPSL", params)
}

// GetTraderDetail performs the GetTraderDetail operation.
func (s *CopyTradingService) GetTraderDetail() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/detail", nil)
}

// GetProfitSummary performs the GetProfitSummary operation.
func (s *CopyTradingService) GetProfitSummary() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/profitSummary", nil)
}

// GetProfitDetail performs the GetProfitDetail operation.
func (s *CopyTradingService) GetProfitDetail(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/profitDetail", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}

// SetCommission performs the SetCommission operation.
func (s *CopyTradingService) SetCommission(commission float64) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/trader/setCommission", map[string]interface{}{
		"commission": commission,
	})
}

// GetTradingPairs performs the GetTradingPairs operation.
func (s *CopyTradingService) GetTradingPairs() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/tradingPairs", nil)
}

// SellSpotOrder performs the SellSpotOrder operation.
func (s *CopyTradingService) SellSpotOrder(buyOrderID string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/spot/trader/sell", map[string]interface{}{
		"buyOrderId": buyOrderID,
	})
}

// GetSpotTraderDetail performs the GetSpotTraderDetail operation.
func (s *CopyTradingService) GetSpotTraderDetail() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/detail", nil)
}

// GetSpotProfitSummary performs the GetSpotProfitSummary operation.
func (s *CopyTradingService) GetSpotProfitSummary() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/profitSummary", nil)
}

// GetSpotProfitDetail performs the GetSpotProfitDetail operation.
func (s *CopyTradingService) GetSpotProfitDetail(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/profitDetail", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}

// GetSpotHistoryOrders performs the GetSpotHistoryOrders operation.
func (s *CopyTradingService) GetSpotHistoryOrders(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/historyOrders", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}
