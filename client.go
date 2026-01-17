package bingx

import (
	"github.com/tigusigalpa/bingx-go/http"
	"github.com/tigusigalpa/bingx-go/services"
	"github.com/tigusigalpa/bingx-go/websocket"
)

type Client struct {
	httpClient  *http.BaseHTTPClient
	market      *services.MarketService
	account     *services.AccountService
	trade       *services.TradeService
	contract    *services.ContractService
	listenKey   *services.ListenKeyService
	wallet      *services.WalletService
	spotAccount *services.SpotAccountService
	subAccount  *services.SubAccountService
	copyTrading *services.CopyTradingService
	coinMClient *CoinMClient
}

func NewClient(apiKey, apiSecret string, options ...ClientOption) *Client {
	config := &ClientConfig{
		BaseURI:           "https://open-api.bingx.com",
		SignatureEncoding: "base64",
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
	client.subAccount = services.NewSubAccountService(httpClient)
	client.copyTrading = services.NewCopyTradingService(httpClient)

	return client
}

type ClientConfig struct {
	BaseURI           string
	SourceKey         string
	SignatureEncoding string
}

type ClientOption func(*ClientConfig)

func WithBaseURI(uri string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURI = uri
	}
}

func WithSourceKey(key string) ClientOption {
	return func(c *ClientConfig) {
		c.SourceKey = key
	}
}

func WithSignatureEncoding(encoding string) ClientOption {
	return func(c *ClientConfig) {
		c.SignatureEncoding = encoding
	}
}

func (c *Client) Market() *services.MarketService {
	return c.market
}

func (c *Client) Account() *services.AccountService {
	return c.account
}

func (c *Client) Trade() *services.TradeService {
	return c.trade
}

func (c *Client) Contract() *services.ContractService {
	return c.contract
}

func (c *Client) ListenKey() *services.ListenKeyService {
	return c.listenKey
}

func (c *Client) Wallet() *services.WalletService {
	return c.wallet
}

func (c *Client) SpotAccount() *services.SpotAccountService {
	return c.spotAccount
}

func (c *Client) SubAccount() *services.SubAccountService {
	return c.subAccount
}

func (c *Client) CopyTrading() *services.CopyTradingService {
	return c.copyTrading
}

func (c *Client) CoinM() *CoinMClient {
	if c.coinMClient == nil {
		c.coinMClient = NewCoinMClient(c.httpClient)
	}
	return c.coinMClient
}

func (c *Client) GetHTTPClient() *http.BaseHTTPClient {
	return c.httpClient
}

func (c *Client) GetEndpoint() string {
	return c.httpClient.GetEndpoint()
}

func (c *Client) GetAPIKey() string {
	return c.httpClient.GetAPIKey()
}

func (c *Client) GetBalance() (map[string]interface{}, error) {
	return c.account.GetBalance()
}

func (c *Client) GetSymbols() (map[string]interface{}, error) {
	return c.market.GetFuturesSymbols()
}

func (c *Client) CreateOrder(params map[string]interface{}) (map[string]interface{}, error) {
	return c.trade.CreateOrder(params)
}

func (c *Client) NewMarketDataStream() *websocket.MarketDataStream {
	return websocket.NewMarketDataStream()
}

func (c *Client) NewAccountDataStream(listenKey string) *websocket.AccountDataStream {
	return websocket.NewAccountDataStream(listenKey)
}
