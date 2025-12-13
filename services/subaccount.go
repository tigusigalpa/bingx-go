package services

import "github.com/tigusigalpa/bingx-go/http"

type SubAccountService struct {
	client *http.BaseHTTPClient
}

func NewSubAccountService(client *http.BaseHTTPClient) *SubAccountService {
	return &SubAccountService{client: client}
}

func (s *SubAccountService) CreateSubAccount(subAccountString string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/create", map[string]interface{}{
		"subAccountString": subAccountString,
	})
}

func (s *SubAccountService) GetAccountUID() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/subAccount/v1/uid", nil)
}

func (s *SubAccountService) GetSubAccountList(subAccountString *string, current, size int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"current": current,
		"size":    size,
	}

	if subAccountString != nil {
		params["subAccountString"] = *subAccountString
	}

	return s.client.Request("GET", "/openApi/subAccount/v1/list", params)
}

func (s *SubAccountService) GetSubAccountAssets(subUID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/subAccount/v1/assets", map[string]interface{}{
		"subUid": subUID,
	})
}

func (s *SubAccountService) UpdateSubAccountStatus(subAccountString string, status int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/status", map[string]interface{}{
		"subAccountString": subAccountString,
		"status":           status,
	})
}

func (s *SubAccountService) GetAllSubAccountBalances() (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/subAccount/v1/allBalances", nil)
}

func (s *SubAccountService) CreateSubAccountAPIKey(subAccountString, label string, permissions map[string]bool, ip *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"subAccountString": subAccountString,
		"label":            label,
		"permissions":      permissions,
	}

	if ip != nil {
		params["ip"] = *ip
	}

	return s.client.Request("POST", "/openApi/subAccount/v1/apiKey/create", params)
}

func (s *SubAccountService) QueryAPIKey(subAccountString string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/subAccount/v1/apiKey/query", map[string]interface{}{
		"subAccountString": subAccountString,
	})
}

func (s *SubAccountService) EditSubAccountAPIKey(subAccountString, apiKey string, permissions map[string]bool, ip *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"subAccountString": subAccountString,
		"apiKey":           apiKey,
		"permissions":      permissions,
	}

	if ip != nil {
		params["ip"] = *ip
	}

	return s.client.Request("POST", "/openApi/subAccount/v1/apiKey/edit", params)
}

func (s *SubAccountService) DeleteSubAccountAPIKey(subAccountString, apiKey string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/apiKey/delete", map[string]interface{}{
		"subAccountString": subAccountString,
		"apiKey":           apiKey,
	})
}

func (s *SubAccountService) AuthorizeSubAccountInternalTransfer(subAccountString string, authorize int) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/innerTransfer/authorize", map[string]interface{}{
		"subAccountString": subAccountString,
		"authorize":        authorize,
	})
}

func (s *SubAccountService) SubAccountInternalTransfer(coin, walletType string, amount float64, transferType string, fromSubUID, toSubUID, clientID *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":         coin,
		"walletType":   walletType,
		"amount":       amount,
		"transferType": transferType,
	}

	if fromSubUID != nil {
		params["fromSubUid"] = *fromSubUID
	}
	if toSubUID != nil {
		params["toSubUid"] = *toSubUID
	}
	if clientID != nil {
		params["clientId"] = *clientID
	}

	return s.client.Request("POST", "/openApi/subAccount/v1/innerTransfer/apply", params)
}

func (s *SubAccountService) GetSubAccountInternalTransferRecords(startTime, endTime *int64, current, size int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"current": current,
		"size":    size,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/subAccount/v1/innerTransfer/records", params)
}

func (s *SubAccountService) SubAccountAssetTransfer(subUID, transferType, asset string, amount float64) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/transfer", map[string]interface{}{
		"subUid": subUID,
		"type":   transferType,
		"asset":  asset,
		"amount": amount,
	})
}

func (s *SubAccountService) GetSubAccountTransferSupportedCoins(subUID string) (map[string]interface{}, error) {
	return s.client.Request("GET", "/openApi/subAccount/v1/transfer/supportCoins", map[string]interface{}{
		"subUid": subUID,
	})
}

func (s *SubAccountService) GetSubAccountAssetTransferHistory(subUID, transferType string, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"subUid": subUID,
		"type":   transferType,
		"limit":  limit,
	}

	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}

	return s.client.Request("GET", "/openApi/subAccount/v1/transfer/history", params)
}

func (s *SubAccountService) CreateSubAccountDepositAddress(coin, network, subUID string) (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/subAccount/v1/capital/deposit/address", map[string]interface{}{
		"coin":    coin,
		"network": network,
		"subUid":  subUID,
	})
}

func (s *SubAccountService) GetSubAccountDepositAddress(coin, subUID string, network *string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"coin":   coin,
		"subUid": subUID,
	}

	if network != nil {
		params["network"] = *network
	}

	return s.client.Request("GET", "/openApi/subAccount/v1/capital/deposit/address", params)
}

func (s *SubAccountService) GetSubAccountDepositHistory(subUID, coin string, status *int, startTime, endTime *int64, limit int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"subUid": subUID,
		"coin":   coin,
		"limit":  limit,
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

	return s.client.Request("GET", "/openApi/subAccount/v1/capital/deposit/history", params)
}
