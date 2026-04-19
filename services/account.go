package services

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

func (s *AccountService) GetBalance() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v3/user/balance", nil)
}

func (s *AccountService) GetPositions(symbol *string) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/positions", params)
}

func (s *AccountService) GetAccountInfo() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/account", nil)
}

func (s *AccountService) GetTradingFees(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/tradingFees", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *AccountService) GetMarginMode(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/getMarginMode", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *AccountService) SetMarginMode(symbol, marginMode string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/setMarginMode", map[string]interface{}{
		"symbol":     symbol,
		"marginMode": marginMode,
	})
}

func (s *AccountService) GetLeverage(symbol string, recvWindow *int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/trade/leverage", params)
}

func (s *AccountService) SetLeverage(symbol, side string, leverage int, recvWindow *int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"side":      side,
		"leverage":  leverage,
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/trade/leverage", params)
}

func (s *AccountService) GetPositionMargin(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/getPositionMargin", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *AccountService) SetPositionMargin(symbol, positionSide string, amount float64, marginType int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/setPositionMargin", map[string]interface{}{
		"symbol":       symbol,
		"positionSide": positionSide,
		"amount":       amount,
		"type":         marginType,
	})
}

func (s *AccountService) GetBalanceHistory(coin string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/balanceHistory", map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	})
}

func (s *AccountService) GetAccountPermissions() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/apiPermissions", nil)
}

func (s *AccountService) GetAPIKey() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/apiKey", nil)
}

func (s *AccountService) GetUserCommissionRates(symbol string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/commissionRate", map[string]interface{}{
		"symbol": symbol,
	})
}

func (s *AccountService) GetAPIRateLimits() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/apiRateLimits", nil)
}

func (s *AccountService) GetDepositHistory(coin string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/depositHistory", map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	})
}

func (s *AccountService) GetWithdrawHistory(coin string, limit int) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/withdrawHistory", map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	})
}

func (s *AccountService) GetAssetDetails(asset string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/assetDetails", map[string]interface{}{
		"asset": asset,
	})
}

func (s *AccountService) GetAllAssets() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/allAssets", nil)
}

func (s *AccountService) GetFundingWallet(asset string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/swap/v2/user/fundingWallet", map[string]interface{}{
		"asset": asset,
	})
}

func (s *AccountService) DustTransfer(assets []string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/dustTransfer", map[string]interface{}{
		"assets": assets,
	})
}

func (s *AccountService) GetPositionRisk(symbol *string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/positionRisk", params)
}

func (s *AccountService) GetIncomeHistory(symbol *string, incomeType *string, startTime, endTime *int64, limit int, recvWindow *int64) (map[string]interface{}, error) {
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
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/income", params)
}

func (s *AccountService) GetCommissionHistory(symbol string, startTime, endTime *int64, limit int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"symbol":    symbol,
		"timestamp": time.Now().UnixMilli(),
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
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/commissionHistory", params)
}

func (s *AccountService) GetForceOrders(symbol *string, autoCloseType *string, startTime, endTime *int64, limit int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if symbol != nil {
		params["symbol"] = *symbol
	}
	if autoCloseType != nil {
		params["autoCloseType"] = *autoCloseType
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
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/forceOrders", params)
}

func (s *AccountService) GetPositionMode(recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/swap/v2/user/positionSide/dual", params)
}

func (s *AccountService) SetPositionMode(dualSidePosition bool, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"dualSidePosition": dualSidePosition,
		"timestamp":        time.Now().UnixMilli(),
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/swap/v2/user/positionSide/dual", params)
}
