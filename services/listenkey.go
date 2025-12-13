package services

import "github.com/tigusigalpa/bingx-go/http"

type ListenKeyService struct {
	client *http.BaseHTTPClient
}

func NewListenKeyService(client *http.BaseHTTPClient) *ListenKeyService {
	return &ListenKeyService{client: client}
}

func (s *ListenKeyService) Generate() (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/user/auth/userDataStream", nil)
}

func (s *ListenKeyService) Extend(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("PUT", "/openApi/user/auth/userDataStream", map[string]interface{}{
		"listenKey": listenKey,
	})
}

func (s *ListenKeyService) Delete(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("DELETE", "/openApi/user/auth/userDataStream", map[string]interface{}{
		"listenKey": listenKey,
	})
}
