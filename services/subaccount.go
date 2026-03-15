package services

import "github.com/tigusigalpa/bingx-go/http"

// Sub-account wallet type constants
const (
	SubAccountWalletTypeFund             = 1  // Fund Account
	SubAccountWalletTypeStandardFutures  = 2  // Standard Futures Account
	SubAccountWalletTypePerpetualFutures = 3  // Perpetual Futures Account
	SubAccountWalletTypeSpot             = 15 // Spot Account
)

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

// SubAccountInternalTransfer performs sub-account internal transfer
// walletType: 1=Fund, 2=Standard Futures, 3=Perpetual Futures, 15=Spot
// userAccountType: 1=UID, 2=Phone number, 3=Email
func (s *SubAccountService) SubAccountInternalTransfer(coin string, walletType int, amount float64, userAccountType int, userAccount string, callingCode, transferClientID *string, recvWindow *int64) (map[string]interface{}, error) {
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

	return s.client.Request("POST", "/openApi/wallets/v1/capital/subAccountInnerTransfer/apply", params)
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

// SubMotherAccountAssetTransfer performs asset transfer between parent and sub-accounts
// Note: This endpoint is only available to the master account
// fromAccountType/toAccountType: 1=Funding, 2=Standard futures, 3=Perpetual U-based, 15=Spot
func (s *SubAccountService) SubMotherAccountAssetTransfer(assetName string, transferAmount float64, fromUID int64, fromType int, fromAccountType int, toUID int64, toType int, toAccountType int, remark string, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"assetName":       assetName,
		"transferAmount":  transferAmount,
		"fromUid":         fromUID,
		"fromType":        fromType,
		"fromAccountType": fromAccountType,
		"toUid":           toUID,
		"toType":          toType,
		"toAccountType":   toAccountType,
		"remark":          remark,
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/account/transfer/v1/subAccount/transferAsset", params)
}

// GetSubMotherAccountTransferableAmount queries supported coins and available transferable amounts
// Note: This endpoint is only available to the master account
func (s *SubAccountService) GetSubMotherAccountTransferableAmount(fromUID int64, fromAccountType int, toUID int64, toAccountType int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"fromUid":         fromUID,
		"fromAccountType": fromAccountType,
		"toUid":           toUID,
		"toAccountType":   toAccountType,
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("POST", "/openApi/account/transfer/v1/subAccount/transferAsset/supportCoins", params)
}

// GetSubMotherAccountTransferHistory queries transfer history between sub-accounts and parent account
// Note: This endpoint is only available to the master account
func (s *SubAccountService) GetSubMotherAccountTransferHistory(uid int64, transferType, tranID *string, startTime, endTime *int64, pageID, pagingSize *int, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"uid": uid,
	}

	if transferType != nil {
		params["type"] = *transferType
	}
	if tranID != nil {
		params["tranId"] = *tranID
	}
	if startTime != nil {
		params["startTime"] = *startTime
	}
	if endTime != nil {
		params["endTime"] = *endTime
	}
	if pageID != nil {
		params["pageId"] = *pageID
	}
	if pagingSize != nil {
		params["pagingSize"] = *pagingSize
	}
	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/account/transfer/v1/subAccount/asset/transferHistory", params)
}
