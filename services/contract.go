package services

import (
	"time"

	"github.com/tigusigalpa/bingx-go/http"
)

type ContractService struct {
	client *http.BaseHTTPClient
}

func NewContractService(client *http.BaseHTTPClient) *ContractService {
	return &ContractService{client: client}
}

func (s *ContractService) GetAllPositions(timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/contract/v1/allPosition", params)
}

func (s *ContractService) GetAllOrders(symbol string, limit int, startTime, endTime *int64) (map[string]interface{}, error) {
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

	return s.client.Request("GET", "/openApi/contract/v1/allOrders", params)
}

func (s *ContractService) GetBalance(timestamp, recvWindow *int64) (map[string]interface{}, error) {
	params := map[string]interface{}{}

	if timestamp != nil {
		params["timestamp"] = *timestamp
	} else {
		params["timestamp"] = time.Now().UnixMilli()
	}

	if recvWindow != nil {
		params["recvWindow"] = *recvWindow
	}

	return s.client.Request("GET", "/openApi/contract/v1/balance", params)
}
