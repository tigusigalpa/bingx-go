package services

import "github.com/tigusigalpa/bingx-go/v2/http"

// MarketService represents a BingX API component or value.
type MarketService struct {
	client *http.BaseHTTPClient
}

// NewMarketService creates a new client or service instance.
func NewMarketService(client *http.BaseHTTPClient) *MarketService {
	return &MarketService{client: client}
}

// GetFuturesSymbols performs the GetFuturesSymbols operation.
func (s *MarketService) GetFuturesSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", nil)
}

// GetSpotSymbols retrieves spot trading symbols
// Response includes maxMarketNotional (max notional for market orders)
// and status field (0=Offline, 1=Online, 5=Pre-open, 10=Accessed, 25=Suspended, 29=Pre-Delisted, 30=Delisted)
func (s *MarketService) GetSpotSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/common/symbols", nil)
}

// GetAllSymbols performs the GetAllSymbols operation.
func (s *MarketService) GetAllSymbols() (map[string]interface{}, error) {
	spot, err := s.GetSpotSymbols()
	if err != nil {
		return nil, err
	}

	futures, err := s.GetFuturesSymbols()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"spot":    spot,
		"futures": futures,
	}, nil
}

// GetSymbols performs the GetSymbols operation.
func (s *MarketService) GetSymbols() (map[string]interface{}, error) {
	return s.GetFuturesSymbols()
}

// GetLatestPrice performs the GetLatestPrice operation.
func (s *MarketService) GetLatestPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/latestPrice", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetSpotLatestPrice performs the GetSpotLatestPrice operation.
func (s *MarketService) GetSpotLatestPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/ticker/price", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetDepth performs the GetDepth operation.
func (s *MarketService) GetDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/depth", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetSpotDepth performs the GetSpotDepth operation.
func (s *MarketService) GetSpotDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/depth", map[string]interface{}{
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

	return s.client.Request("GET", "/openApi/swap/v3/quote/klines", params)
}

// GetSpotKlines retrieves spot K-line (candlestick) data
// timeZone: optional timezone offset (0=UTC (default), 8=UTC+8)
func (s *MarketService) GetSpotKlines(symbol, interval string, limit int, startTime, endTime, timeZone *int64) (map[string]interface{}, error) {
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
	if timeZone != nil {
		params["timeZone"] = *timeZone
	}

	return s.client.Request("GET", "/openApi/spot/v2/market/kline", params)
}

// Get24hrTicker performs the Get24hrTicker operation.
func (s *MarketService) Get24hrTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/quote/ticker", params)
}

// GetSpot24hrTicker performs the GetSpot24hrTicker operation.
func (s *MarketService) GetSpot24hrTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/spot/v1/market/ticker/24hr", params)
}

// GetFundingRateHistory performs the GetFundingRateHistory operation.
func (s *MarketService) GetFundingRateHistory(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/fundingRate/history", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetMarkPrice performs the GetMarkPrice operation.
func (s *MarketService) GetMarkPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/premiumIndex", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetPremiumIndexKlines performs the GetPremiumIndexKlines operation.
func (s *MarketService) GetPremiumIndexKlines(symbol, interval string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/market/premiumIndexKline", params)
}

// GetAggregateTrades performs the GetAggregateTrades operation.
func (s *MarketService) GetAggregateTrades(symbol string, limit int, fromID, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}

	if fromID != nil {
		params["fromId"] = *fromID
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/market/aggTrades", params)
}

// GetRecentTrades performs the GetRecentTrades operation.
func (s *MarketService) GetRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetSpotAggregateTrades performs the GetSpotAggregateTrades operation.
func (s *MarketService) GetSpotAggregateTrades(symbol string, limit int, fromID *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}

	if fromID != nil {
		params["fromId"] = *fromID
	}

	return s.client.Request("GET", "/openApi/spot/v1/market/aggTrades", params)
}

// GetSpotRecentTrades performs the GetSpotRecentTrades operation.
func (s *MarketService) GetSpotRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetServerTime performs the GetServerTime operation.
func (s *MarketService) GetServerTime() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/time", nil)
}

// GetSpotServerTime performs the GetSpotServerTime operation.
func (s *MarketService) GetSpotServerTime() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/time", nil)
}

// GetContinuousKlines performs the GetContinuousKlines operation.
func (s *MarketService) GetContinuousKlines(symbol, interval string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/market/continuousKline", params)
}

// GetIndexPriceKlines performs the GetIndexPriceKlines operation.
func (s *MarketService) GetIndexPriceKlines(symbol, interval string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/market/indexPriceKline", params)
}

// GetTopLongShortRatio performs the GetTopLongShortRatio operation.
func (s *MarketService) GetTopLongShortRatio(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/topLongShortRatio", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetTopTradersPositionRatio performs the GetTopTradersPositionRatio operation.
func (s *MarketService) GetTopTradersPositionRatio(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/topTraderPositionRatio", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetHistoricalTopLongShortRatio performs the GetHistoricalTopLongShortRatio operation.
func (s *MarketService) GetHistoricalTopLongShortRatio(symbol string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/market/topLongShortAccount", params)
}

// GetTopTradersLongShortRatio performs the GetTopTradersLongShortRatio operation.
func (s *MarketService) GetTopTradersLongShortRatio(symbol string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/market/topLongShortPosition", params)
}

// GetBasis performs the GetBasis operation.
func (s *MarketService) GetBasis(symbol, contractType string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":       symbol,
		"contractType": contractType,
		"limit":        limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/market/basis", params)
}

// GetOpenInterest performs the GetOpenInterest operation.
func (s *MarketService) GetOpenInterest(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/openInterest", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetOpenInterestHistory performs the GetOpenInterestHistory operation.
func (s *MarketService) GetOpenInterestHistory(symbol, period string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"period": period,
		"limit":  limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/market/openInterest/history", params)
}

// GetFundingRateInfo performs the GetFundingRateInfo operation.
func (s *MarketService) GetFundingRateInfo(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/fundingRate", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetBookTicker performs the GetBookTicker operation.
func (s *MarketService) GetBookTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/quote/bookTicker", params)
}

// GetSpotBookTicker performs the GetSpotBookTicker operation.
func (s *MarketService) GetSpotBookTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/spot/v1/market/bookTicker", params)
}

// GetIndexPrice performs the GetIndexPrice operation.
func (s *MarketService) GetIndexPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/indexPrice", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetTickerPrice performs the GetTickerPrice operation.
func (s *MarketService) GetTickerPrice(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/market/ticker/price", params)
}
