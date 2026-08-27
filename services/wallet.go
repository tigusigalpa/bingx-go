package services

import "github.com/tigusigalpa/bingx-go/v2/http"

// WalletService represents a BingX API component or value.
type WalletService struct {
	client *http.BaseHTTPClient
}

// NewWalletService creates a new client or service instance.
func NewWalletService(client *http.BaseHTTPClient) *WalletService {
	return &WalletService{client: client}
}

// GetDepositHistory performs the GetDepositHistory operation.
func (s *WalletService) GetDepositHistory(coin string, status *int, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	}

	if status != nil {
		params["status"] = *status
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/wallets/v1/capital/deposit/history", params)
}

// GetDepositAddress performs the GetDepositAddress operation.
func (s *WalletService) GetDepositAddress(coin, network string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/wallets/v1/capital/deposit/address", map[string]interface{}{
		"coin":    coin,
		"network": network,
	})
}

// GetWithdrawalHistory performs the GetWithdrawalHistory operation.
func (s *WalletService) GetWithdrawalHistory(coin string, status *int, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	}

	if status != nil {
		params["status"] = *status
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/wallets/v1/capital/withdraw/history", params)
}

// Withdraw performs the Withdraw operation.
func (s *WalletService) Withdraw(coin, address string, amount float64, network string, addressTag *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":    coin,
		"address": address,
		"amount":  amount,
		"network": network,
	}

	if addressTag != nil {
		params["addressTag"] = *addressTag
	}

	return s.client.Request("POST", "/openApi/wallets/v1/capital/withdraw/apply", params)
}

// GetAllCoinInfo performs the GetAllCoinInfo operation.
func (s *WalletService) GetAllCoinInfo() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/wallets/v1/capital/config/getall", nil)
}

// GetMainAccountTransferHistory performs the GetMainAccountTransferHistory operation.
func (s *WalletService) GetMainAccountTransferHistory(coin string, transferType *string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":  coin,
		"limit": limit,
	}

	if transferType != nil {
		params["type"] = *transferType
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/wallets/v1/capital/transfer/history", params)
}
