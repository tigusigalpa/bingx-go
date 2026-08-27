package bingx

import (
	"sync"

	"github.com/tigusigalpa/bingx-go/v2/http"
	"github.com/tigusigalpa/bingx-go/v2/services"
	"github.com/tigusigalpa/bingx-go/v2/websocket"
)

// Client represents a BingX API component or value.
type Client struct {
	httpClient   *http.BaseHTTPClient
	market       *services.MarketService
	account      *services.AccountService
	trade        *services.TradeService
	contract     *services.ContractService
	listenKey    *services.ListenKeyService
	wallet       *services.WalletService
	spotAccount  *services.SpotAccountService
	spotTrade    *services.SpotTradeService
	subAccount   *services.SubAccountService
	copyTrading  *services.CopyTradingService
	coinMClient  *CoinMClient
	tradfiClient *TradFiClient
	lazyMu       sync.Mutex
}

// NewClient creates a new client or service instance.
func NewClient(apiKey, apiSecret string, options ...ClientOption) *Client {
	config := &ClientConfig{
		BaseURI:           "https://open-api.bingx.com",
		SignatureEncoding: "hex",
	}

	for _, opt := range options {
		opt(config)
	}

	httpClient := http.NewBaseHTTPClient(
		apiKey,
		apiSecret,
		config.BaseURI,
		config.SourceKey,
		config.SignatureEncoding,
	)

	client := &Client{
		httpClient: httpClient,
	}

	client.market = services.NewMarketService(httpClient)
	client.account = services.NewAccountService(httpClient)
	client.trade = services.NewTradeService(httpClient)
	client.contract = services.NewContractService(httpClient)
	client.listenKey = services.NewListenKeyService(httpClient)
	client.wallet = services.NewWalletService(httpClient)
	client.spotAccount = services.NewSpotAccountService(httpClient)
	client.spotTrade = services.NewSpotTradeService(httpClient)
	client.subAccount = services.NewSubAccountService(httpClient)
	client.copyTrading = services.NewCopyTradingService(httpClient)

	return client
}

// NewDemoClient creates a client configured for demo trading (VST environment)
func NewDemoClient(apiKey, apiSecret string, options ...ClientOption) *Client {
	// Prepend the demo environment option
	demoOptions := []ClientOption{WithDemoEnvironment()}
	demoOptions = append(demoOptions, options...)
	return NewClient(apiKey, apiSecret, demoOptions...)
}

// ClientConfig represents a BingX API component or value.
type ClientConfig struct {
	BaseURI           string
	SourceKey         string
	SignatureEncoding string
}

// ClientOption represents a BingX API component or value.
type ClientOption func(*ClientConfig)

// WithBaseURI performs the WithBaseURI operation.
func WithBaseURI(uri string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURI = uri
	}
}

// WithSourceKey performs the WithSourceKey operation.
func WithSourceKey(key string) ClientOption {
	return func(c *ClientConfig) {
		c.SourceKey = key
	}
}

// WithSignatureEncoding sets the HMAC signature encoding.
// "hex" is the BingX-compatible default and should be used for production calls.
// "base64" is supported only for backward compatibility.
func WithSignatureEncoding(encoding string) ClientOption {
	return func(c *ClientConfig) {
		c.SignatureEncoding = encoding
	}
}

// WithDemoEnvironment configures the client for demo trading (VST environment)
func WithDemoEnvironment() ClientOption {
	return func(c *ClientConfig) {
		c.BaseURI = "https://open-api-vst.bingx.com"
	}
}

// Market performs the Market operation.
func (c *Client) Market() *services.MarketService {
	return c.market
}

// Account performs the Account operation.
func (c *Client) Account() *services.AccountService {
	return c.account
}

// Trade performs the Trade operation.
func (c *Client) Trade() *services.TradeService {
	return c.trade
}

// Contract performs the Contract operation.
func (c *Client) Contract() *services.ContractService {
	return c.contract
}

// ListenKey performs the ListenKey operation.
func (c *Client) ListenKey() *services.ListenKeyService {
	return c.listenKey
}

// Wallet performs the Wallet operation.
func (c *Client) Wallet() *services.WalletService {
	return c.wallet
}

// SpotAccount performs the SpotAccount operation.
func (c *Client) SpotAccount() *services.SpotAccountService {
	return c.spotAccount
}

// SpotTrade performs the SpotTrade operation.
func (c *Client) SpotTrade() *services.SpotTradeService {
	return c.spotTrade
}

// SubAccount performs the SubAccount operation.
func (c *Client) SubAccount() *services.SubAccountService {
	return c.subAccount
}

// CopyTrading performs the CopyTrading operation.
func (c *Client) CopyTrading() *services.CopyTradingService {
	return c.copyTrading
}

// CoinM performs the CoinM operation.
func (c *Client) CoinM() *CoinMClient {
	c.lazyMu.Lock()
	defer c.lazyMu.Unlock()
	if c.coinMClient == nil {
		c.coinMClient = NewCoinMClient(c.httpClient)
	}
	return c.coinMClient
}

// TradFi returns the TradFi client for Traditional Finance instruments (stocks, forex, commodities, indices).
func (c *Client) TradFi() *TradFiClient {
	c.lazyMu.Lock()
	defer c.lazyMu.Unlock()
	if c.tradfiClient == nil {
		c.tradfiClient = NewTradFiClient(c.httpClient)
	}
	return c.tradfiClient
}

// GetHTTPClient performs the GetHTTPClient operation.
func (c *Client) GetHTTPClient() *http.BaseHTTPClient {
	return c.httpClient
}

// GetEndpoint performs the GetEndpoint operation.
func (c *Client) GetEndpoint() string {
	return c.httpClient.GetEndpoint()
}

// GetAPIKey performs the GetAPIKey operation.
func (c *Client) GetAPIKey() string {
	return c.httpClient.GetAPIKey()
}

// GetBalance performs the GetBalance operation.
func (c *Client) GetBalance() (map[string]interface{}, error) {
	return c.account.GetBalance()
}

// GetSymbols performs the GetSymbols operation.
func (c *Client) GetSymbols() (map[string]interface{}, error) {
	return c.market.GetFuturesSymbols()
}

// CreateOrder performs the CreateOrder operation.
func (c *Client) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return c.trade.CreateOrder(params)
}

// NewMarketDataStream creates a new client or service instance.
func (c *Client) NewMarketDataStream() *websocket.MarketDataStream {
	return websocket.NewMarketDataStream()
}

// NewAccountDataStream creates a new client or service instance.
func (c *Client) NewAccountDataStream(listenKey string) *websocket.AccountDataStream {
	return websocket.NewAccountDataStream(listenKey)
}
