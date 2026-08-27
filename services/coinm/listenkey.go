package coinm

import "github.com/tigusigalpa/bingx-go/v2/http"

// ListenKeyService represents a BingX API component or value.
type ListenKeyService struct {
	client *http.BaseHTTPClient
}

// NewListenKeyService creates a new client or service instance.
func NewListenKeyService(client *http.BaseHTTPClient) *ListenKeyService {
	return &ListenKeyService{client: client}
}

// Generate performs the Generate operation.
func (s *ListenKeyService) Generate() (map[string]interface{}, error) {
	return s.client.Request("POST", "/openApi/user/auth/userDataStream", nil)
}

// Extend performs the Extend operation.
func (s *ListenKeyService) Extend(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("PUT", "/openApi/user/auth/userDataStream", map[string]interface{}{
		"listenKey": listenKey,
	})
}

// Delete performs the Delete operation.
func (s *ListenKeyService) Delete(listenKey string) (map[string]interface{}, error) {
	return s.client.Request("DELETE", "/openApi/user/auth/userDataStream", map[string]interface{}{
		"listenKey": listenKey,
	})
}
