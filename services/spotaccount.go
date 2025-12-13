package services

import "github.com/tigusigalpa/bingx-go/http"

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

func (s *SpotAccountService) InternalTransfer(coin, walletType string, amount float64, transferType, subUID string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/wallets/v1/capital/innerTransfer/apply", map[string]interface{}{
		"coin":         coin,
		"walletType":   walletType,
		"amount":       amount,
		"transferType": transferType,
		"subUid":       subUID,
	})
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
