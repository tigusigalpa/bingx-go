package services

import "github.com/tigusigalpa/bingx-go/v2/http"

// Wallet type constants for internal transfers
const (
	WalletTypeFund             = 1 // Fund Account
	WalletTypeStandardFutures  = 2 // Standard Futures Account
	WalletTypePerpetualFutures = 3 // Perpetual Futures Account
	WalletTypeSpot             = 4 // Spot Account
)

type SpotAccountService struct {
	client *http.BaseHTTPClient
}

func NewSpotAccountService(client *http.BaseHTTPClient) *SpotAccountService {
	return &SpotAccountService{client: client}
}

func (s *SpotAccountService) GetBalance() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/account/balance", nil)
}

func (s *SpotAccountService) GetFundBalance() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/wallets/v1/capital/fundBalance", nil)
}

func (s *SpotAccountService) UniversalTransfer(transferType, asset string, amount float64) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/wallets/v1/capital/transfer", map[string]interface{}{
		"type":   transferType,
		"asset":  asset,
		"amount": amount,
	})
}

func (s *SpotAccountService) GetAssetTransferRecords(transferType string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"type":  transferType,
		"limit": limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/wallets/v1/capital/transfer/records", params)
}

// InternalTransfer performs main account internal transfer
// walletType: 1=Fund, 2=Standard Futures, 3=Perpetual Futures, 4=Spot
// userAccountType: 1=UID, 2=Phone number, 3=Email
func (s *SpotAccountService) InternalTransfer(coin string, walletType int, amount float64, userAccountType int, userAccount string, callingCode, transferClientID *string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":            coin,
		"walletType":      walletType,
		"amount":          amount,
		"userAccountType": userAccountType,
		"userAccount":     userAccount,
	}

	if callingCode != nil {
		params["callingCode"] = *callingCode
	}
	if transferClientID != nil {
		params["transferClientId"] = *transferClientID
	}
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/wallets/v1/capital/innerTransfer/apply", params)
}

func (s *SpotAccountService) GetInternalTransferRecords(coin string, transferType *string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	}

	if transferType != nil {
		params["transferType"] = *transferType
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/wallets/v1/capital/innerTransfer/records", params)
}

func (s *SpotAccountService) GetAllAccountBalances() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/account/allBalances", nil)
}

func (s *SpotAccountService) GetAccountType() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/spot/v1/account/type", nil)
}
