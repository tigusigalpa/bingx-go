package services

import "github.com/tigusigalpa/bingx-go/http"

type CopyTradingService struct {
	client *http.BaseHTTPClient
}

func NewCopyTradingService(client *http.BaseHTTPClient) *CopyTradingService {
	return &CopyTradingService{client: client}
}

func (s *CopyTradingService) GetCurrentTrackOrders(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/currentTrack", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *CopyTradingService) CloseTrackOrder(orderNumber string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/trader/closeTrack", map[string]interface{}{
		"orderNumber": orderNumber,
	})
}

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

func (s *CopyTradingService) GetTraderDetail() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/detail", nil)
}

func (s *CopyTradingService) GetProfitSummary() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/profitSummary", nil)
}

func (s *CopyTradingService) GetProfitDetail(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/profitDetail", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}

func (s *CopyTradingService) SetCommission(commission float64) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/trader/setCommission", map[string]interface{}{
		"commission": commission,
	})
}

func (s *CopyTradingService) GetTradingPairs() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/trader/tradingPairs", nil)
}

func (s *CopyTradingService) SellSpotOrder(buyOrderID string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/copy/v1/spot/trader/sell", map[string]interface{}{
		"buyOrderId": buyOrderID,
	})
}

func (s *CopyTradingService) GetSpotTraderDetail() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/detail", nil)
}

func (s *CopyTradingService) GetSpotProfitSummary() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/profitSummary", nil)
}

func (s *CopyTradingService) GetSpotProfitDetail(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/profitDetail", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}

func (s *CopyTradingService) GetSpotHistoryOrders(pageIndex, pageSize int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/copy/v1/spot/trader/historyOrders", map[string]interface{}{
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}
