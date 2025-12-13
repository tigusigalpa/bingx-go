package coinm

import "github.com/tigusigalpa/bingx-go/http"

type MarketService struct {
	client *http.BaseHTTPClient
}

func NewMarketService(client *http.BaseHTTPClient) *MarketService {
	return &MarketService{client: client}
}

func (s *MarketService) GetContracts() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/market/contracts", nil)
}

func (s *MarketService) GetTicker(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/market/ticker", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *MarketService) GetDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/market/depth", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

func (s *MarketService) GetKlines(symbol, interval string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":   symbol,
		"interval": interval,
		"limit":    limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v1/market/kline", params)
}

func (s *MarketService) GetOpenInterest(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/market/openInterest", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *MarketService) GetFundingRate(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/market/fundingRate", map[string]interface{}{
		"symbol": symbol,
	})
}
