package services

import "github.com/tigusigalpa/bingx-go/http"

type MarketService struct {
	client *http.BaseHTTPClient
}

func NewMarketService(client *http.BaseHTTPClient) *MarketService {
	return &MarketService{client: client}
}

func (s *MarketService) GetFuturesSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/symbols", nil)
}

func (s *MarketService) GetSpotSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/symbols", nil)
}

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

func (s *MarketService) GetSymbols() (map[string]interface{}, error) {
	return s.GetFuturesSymbols()
}

func (s *MarketService) GetLatestPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/latestPrice", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *MarketService) GetSpotLatestPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/ticker/price", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *MarketService) GetDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/depth", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

func (s *MarketService) GetSpotDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/depth", map[string]interface{}{
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

	return s.client.Request("GET", "/openApi/swap/v2/market/kline", params)
}

func (s *MarketService) GetSpotKlines(symbol, interval string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/spot/v1/market/klines", params)
}

func (s *MarketService) Get24hrTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/market/ticker24hr", params)
}

func (s *MarketService) GetSpot24hrTicker(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/spot/v1/market/ticker/24hr", params)
}

func (s *MarketService) GetFundingRateHistory(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/fundingRate/history", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

func (s *MarketService) GetMarkPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/markPrice", map[string]interface{}{
		"symbol": symbol,
	})
}

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

func (s *MarketService) GetRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

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

func (s *MarketService) GetSpotRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

func (s *MarketService) GetServerTime() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/time", nil)
}

func (s *MarketService) GetSpotServerTime() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/market/time", nil)
}

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

func (s *MarketService) GetTopLongShortRatio(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/topLongShortRatio", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

func (s *MarketService) GetTopTradersPositionRatio(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/topTraderPositionRatio", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

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
