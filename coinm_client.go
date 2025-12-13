package bingx

import (
	"github.com/tigusigalpa/bingx-go/http"
	"github.com/tigusigalpa/bingx-go/services/coinm"
)

type CoinMClient struct {
	httpClient *http.BaseHTTPClient
	market     *coinm.MarketService
	trade      *coinm.TradeService
	listenKey  *coinm.ListenKeyService
}

func NewCoinMClient(httpClient *http.BaseHTTPClient) *CoinMClient {
	return &CoinMClient{
		httpClient: httpClient,
		market:     coinm.NewMarketService(httpClient),
		trade:      coinm.NewTradeService(httpClient),
		listenKey:  coinm.NewListenKeyService(httpClient),
	}
}

func (c *CoinMClient) Market() *coinm.MarketService {
	return c.market
}

func (c *CoinMClient) Trade() *coinm.TradeService {
	return c.trade
}

func (c *CoinMClient) ListenKey() *coinm.ListenKeyService {
	return c.listenKey
}
