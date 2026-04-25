package bingx

import (
	"github.com/tigusigalpa/bingx-go/v2/http"
	"github.com/tigusigalpa/bingx-go/v2/services/tradfi"
)

// TradFiClient provides access to Traditional Finance (TradFi) instruments on BingX.
// TradFi includes stock tokens (TSLA, AAPL), forex pairs (EUR-USD), commodities (GOLD, OIL),
// and stock indices (SPX, DJI) traded as perpetual swaps.
type TradFiClient struct {
	httpClient *http.BaseHTTPClient
	market     *tradfi.MarketService
	trade      *tradfi.TradeService
	account    *tradfi.AccountService
	listenKey  *tradfi.ListenKeyService
}

// NewTradFiClient creates a new TradFi client using the provided HTTP client.
func NewTradFiClient(httpClient *http.BaseHTTPClient) *TradFiClient {
	return &TradFiClient{
		httpClient: httpClient,
		market:     tradfi.NewMarketService(httpClient),
		trade:      tradfi.NewTradeService(httpClient),
		account:    tradfi.NewAccountService(httpClient),
		listenKey:  tradfi.NewListenKeyService(httpClient),
	}
}

// Market returns the TradFi market data service.
func (c *TradFiClient) Market() *tradfi.MarketService {
	return c.market
}

// Trade returns the TradFi trading service.
func (c *TradFiClient) Trade() *tradfi.TradeService {
	return c.trade
}

// Account returns the TradFi account service.
func (c *TradFiClient) Account() *tradfi.AccountService {
	return c.account
}

// ListenKey returns the TradFi listen key service for WebSocket streams.
func (c *TradFiClient) ListenKey() *tradfi.ListenKeyService {
	return c.listenKey
}
