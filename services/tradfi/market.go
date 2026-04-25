package tradfi

import "github.com/tigusigalpa/bingx-go/v2/http"

type MarketService struct {
	client *http.BaseHTTPClient
}

func NewMarketService(client *http.BaseHTTPClient) *MarketService {
	return &MarketService{client: client}
}

// GetSymbols retrieves all TradFi trading symbols (stocks, forex, commodities, indices)
func (s *MarketService) GetSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", nil)
}

// GetStockSymbols retrieves stock token symbols (TSLA, AAPL, etc.)
func (s *MarketService) GetStockSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", map[string]interface{}{
		"assetType": "STOCK",
	})
}

// GetForexSymbols retrieves forex trading symbols (EUR-USD, GBP-USD, etc.)
func (s *MarketService) GetForexSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", map[string]interface{}{
		"assetType": "FOREX",
	})
}

// GetCommoditySymbols retrieves commodity symbols (GOLD, SILVER, OIL, etc.)
func (s *MarketService) GetCommoditySymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", map[string]interface{}{
		"assetType": "COMMODITY",
	})
}

// GetIndexSymbols retrieves stock index symbols (SPX, DJI, etc.)
func (s *MarketService) GetIndexSymbols() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/contracts", map[string]interface{}{
		"assetType": "INDEX",
	})
}

// GetTicker retrieves 24h ticker data for a TradFi symbol
func (s *MarketService) GetTicker(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/ticker", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetLatestPrice retrieves latest price for a TradFi symbol
func (s *MarketService) GetLatestPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/latestPrice", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetDepth retrieves order book depth for a TradFi symbol
func (s *MarketService) GetDepth(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/depth", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetKlines retrieves kline/candlestick data for a TradFi symbol
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

// GetMarkPrice retrieves mark price for a TradFi symbol
func (s *MarketService) GetMarkPrice(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/premiumIndex", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetFundingRate retrieves current funding rate for a TradFi perpetual
func (s *MarketService) GetFundingRate(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/fundingRate", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetFundingRateHistory retrieves historical funding rates
func (s *MarketService) GetFundingRateHistory(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/market/fundingRate/history", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetOpenInterest retrieves open interest for a TradFi symbol
func (s *MarketService) GetOpenInterest(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/openInterest", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetRecentTrades retrieves recent public trades
func (s *MarketService) GetRecentTrades(symbol string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/trades", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	})
}

// GetBookTicker retrieves best bid/ask price and quantity
func (s *MarketService) GetBookTicker(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/quote/bookTicker", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetTradingRules retrieves trading rules and specifications for TradFi symbols
func (s *MarketService) GetTradingRules(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v1/tradingRules", map[string]interface{}{
		"symbol": symbol,
	})
}
