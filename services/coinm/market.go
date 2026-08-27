package coinm

import "github.com/tigusigalpa/bingx-go/v2/http"

// MarketService represents a BingX API component or value.
type MarketService struct {
	client *http.BaseHTTPClient
}

// NewMarketService creates a new client or service instance.
func NewMarketService(client *http.BaseHTTPClient) *MarketService {
	return &MarketService{client: client}
}

// GetContracts performs the GetContracts operation.
func (s *MarketService) GetContracts() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/contracts", nil)
}

// GetTicker performs the GetTicker operation.
func (s *MarketService) GetTicker(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/ticker", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetDepth performs the GetDepth operation.
func (s *MarketService) GetDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/depth", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetKlines performs the GetKlines operation.
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

	return s.client.Request("GET", "/openApi/cswap/v1/market/klines", params)
}

// GetOpenInterest performs the GetOpenInterest operation.
func (s *MarketService) GetOpenInterest(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/openInterest", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetFundingRate performs the GetFundingRate operation.
func (s *MarketService) GetFundingRate(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/premiumIndex", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetFundingRateHistory performs the GetFundingRateHistory operation.
func (s *MarketService) GetFundingRateHistory(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/fundingRate", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetMarkPrice performs the GetMarkPrice operation.
func (s *MarketService) GetMarkPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/premiumIndex", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetIndexPrice performs the GetIndexPrice operation.
func (s *MarketService) GetIndexPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/premiumIndex", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetRecentTrades performs the GetRecentTrades operation.
func (s *MarketService) GetRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/cswap/v1/market/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}
