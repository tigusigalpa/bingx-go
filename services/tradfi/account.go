package tradfi

import (
	"time"

	"github.com/tigusigalpa/bingx-go/v2/http"
)

type AccountService struct {
	client *http.BaseHTTPClient
}

func NewAccountService(client *http.BaseHTTPClient) *AccountService {
	return &AccountService{client: client}
}

// GetBalance retrieves TradFi account balance
func (s *AccountService) GetBalance() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v3/user/balance", nil)
}

// GetAccountInfo retrieves comprehensive account information
func (s *AccountService) GetAccountInfo() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/account", nil)
}

// GetPositions retrieves open positions for TradFi instruments
func (s *AccountService) GetPositions(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/positions", params)
}

// GetPositionRisk retrieves position risk data including liquidation price
func (s *AccountService) GetPositionRisk(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/positionRisk", params)
}

// GetIncomeHistory retrieves income history (PNL, funding fees, commissions)
func (s *AccountService) GetIncomeHistory(symbol *string, incomeType *string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if incomeType != nil {
		params["incomeType"] = *incomeType
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/income", params)
}

// GetCommissionHistory retrieves commission history for TradFi trades
func (s *AccountService) GetCommissionHistory(symbol string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/swap/v2/user/commissionRate", params)
}

// GetForceOrders retrieves liquidation/force order history
func (s *AccountService) GetForceOrders(symbol *string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"limit": limit,
	}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/forceOrders", params)
}

// GetPositionMode retrieves current position mode (hedge or one-way)
func (s *AccountService) GetPositionMode() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/positionMode", nil)
}

// SetPositionMode sets position mode (true = hedge mode, false = one-way mode)
func (s *AccountService) SetPositionMode(hedgeMode bool) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/positionMode", map[string]interface{}{
		"positionMode": map[bool]string{true: "HEDGE", false: "ONEWAY"}[hedgeMode],
	})
}

// GetMarginMode retrieves margin mode for a symbol
func (s *AccountService) GetMarginMode(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/getMarginMode", map[string]interface{}{
		"symbol": symbol,
	})
}

// SetMarginMode sets margin mode (ISOLATED or CROSSED)
func (s *AccountService) SetMarginMode(symbol, marginMode string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/setMarginMode", map[string]interface{}{
		"symbol":     symbol,
		"marginMode": marginMode,
	})
}

// GetTradingFees retrieves trading fee rates for TradFi instruments
func (s *AccountService) GetTradingFees(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/tradingFees", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetUserCommissionRates retrieves user commission rates
func (s *AccountService) GetUserCommissionRates(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/commissionRate", map[string]interface{}{
		"symbol": symbol,
	})
}

// GetMultiAssetsMode retrieves multi-assets margin mode status
func (s *AccountService) GetMultiAssetsMode() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/multiAssetsMode", nil)
}

// SetMultiAssetsMode enables/disables multi-assets margin mode
func (s *AccountService) SetMultiAssetsMode(enabled bool) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/multiAssetsMode", map[string]interface{}{
		"multiAssetsMode": enabled,
	})
}

// GetMultiAssetsMargin retrieves multi-assets margin details
func (s *AccountService) GetMultiAssetsMargin() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/multiAssetsMargin", nil)
}

// GetAPIPermissions retrieves API key permissions
func (s *AccountService) GetAPIPermissions() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/apiPermissions", nil)
}

// GetBalanceHistory retrieves balance history
func (s *AccountService) GetBalanceHistory(coin string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/balanceHistory", map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	})
}
