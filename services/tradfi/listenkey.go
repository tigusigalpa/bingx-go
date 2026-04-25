package tradfi

import "github.com/tigusigalpa/bingx-go/v2/http"

type ListenKeyService struct {
	client *http.BaseHTTPClient
}

func NewListenKeyService(client *http.BaseHTTPClient) *ListenKeyService {
	return &ListenKeyService{client: client}
}

// Create generates a new listen key for TradFi WebSocket streams
func (s *ListenKeyService) Create() (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/swap/v2/user/listenKey", nil)
}

// Extend extends the validity of an existing listen key
func (s *ListenKeyService) Extend(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("PUT", "/openApi/swap/v2/user/listenKey", map[string]interface{}{
		"listenKey": listenKey,
	})
}

// Delete deletes a listen key
func (s *ListenKeyService) Delete(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("DELETE", "/openApi/swap/v2/user/listenKey", map[string]interface{}{
		"listenKey": listenKey,
	})
}
