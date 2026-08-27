package errors

import "fmt"

// BingXException represents a BingX API component or value.
type BingXException struct {
	Message  string
	Code     int
	Response map[string]interface{}
}

// Error performs the Error operation.
func (e *BingXException) Error() string {
	return e.Message
}

// GetResponse performs the GetResponse operation.
func (e *BingXException) GetResponse() map[string]interface{} {
	return e.Response
}

// NewBingXException creates a new client or service instance.
func NewBingXException(message string, code int, response map[string]interface{}) *BingXException {
	return &BingXException{
		Message:  message,
		Code:     code,
		Response: response,
	}
}

// APIException represents a BingX API component or value.
type APIException struct {
	*BingXException
	APICode string
}

// Error performs the Error operation.
func (e *APIException) Error() string {
	return fmt.Sprintf("API Error [%s]: %s", e.APICode, e.Message)
}

// NewAPIException creates a new client or service instance.
func NewAPIException(message, apiCode string, response map[string]interface{}) *APIException {
	return &APIException{
		BingXException: NewBingXException(message, 0, response),
		APICode:        apiCode,
	}
}

// AuthenticationException represents a BingX API component or value.
type AuthenticationException struct {
	*BingXException
}

// Error performs the Error operation.
func (e *AuthenticationException) Error() string {
	return fmt.Sprintf("Authentication Error: %s", e.Message)
}

// NewAuthenticationException creates a new client or service instance.
func NewAuthenticationException(message string, response map[string]interface{}) *AuthenticationException {
	return &AuthenticationException{
		BingXException: NewBingXException(message, 0, response),
	}
}

// RateLimitException represents a BingX API component or value.
type RateLimitException struct {
	*BingXException
}

// Error performs the Error operation.
func (e *RateLimitException) Error() string {
	return fmt.Sprintf("Rate Limit Exceeded: %s", e.Message)
}

// NewRateLimitException creates a new client or service instance.
func NewRateLimitException(message string, response map[string]interface{}) *RateLimitException {
	return &RateLimitException{
		BingXException: NewBingXException(message, 0, response),
	}
}

// InsufficientBalanceException represents a BingX API component or value.
type InsufficientBalanceException struct {
	*BingXException
}

// Error performs the Error operation.
func (e *InsufficientBalanceException) Error() string {
	return fmt.Sprintf("Insufficient Balance: %s", e.Message)
}

// NewInsufficientBalanceException creates a new client or service instance.
func NewInsufficientBalanceException(message string, response map[string]interface{}) *InsufficientBalanceException {
	return &InsufficientBalanceException{
		BingXException: NewBingXException(message, 0, response),
	}
}
