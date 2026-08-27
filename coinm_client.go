package bingx

import (
	"github.com/tigusigalpa/bingx-go/v2/http"
	"github.com/tigusigalpa/bingx-go/v2/services/coinm"
)

// CoinMClient represents a BingX API component or value.
type CoinMClient struct {
	httpClient *http.BaseHTTPClient
	market     *coinm.MarketService
	trade      *coinm.TradeService
	listenKey  *coinm.ListenKeyService
}

// NewCoinMClient creates a new client or service instance.
func NewCoinMClient(httpClient *http.BaseHTTPClient) *CoinMClient {
	return &CoinMClient{
		httpClient: httpClient,
		market:     coinm.NewMarketService(httpClient),
		trade:      coinm.NewTradeService(httpClient),
		listenKey:  coinm.NewListenKeyService(httpClient),
	}
}

// Market performs the Market operation.
func (c *CoinMClient) Market() *coinm.MarketService {
	return c.market
}

// Trade performs the Trade operation.
func (c *CoinMClient) Trade() *coinm.TradeService {
	return c.trade
}

// ListenKey performs the ListenKey operation.
func (c *CoinMClient) ListenKey() *coinm.ListenKeyService {
	return c.listenKey
}
