package errors

import "fmt"

type BingXException struct {
	Message  string
	Code     int
	Response map[string]interface{}
}

func (e *BingXException) Error() string {
	return e.Message
}

func (e *BingXException) GetResponse() map[string]interface{} {
	return e.Response
}

func NewBingXException(message string, code int, response map[string]interface{}) *BingXException {
	return &BingXException{
		Message:  message,
		Code:     code,
		Response: response,
	}
}

type APIException struct {
	*BingXException
	APICode string
}

func (e *APIException) Error() string {
	return fmt.Sprintf("API Error [%s]: %s", e.APICode, e.Message)
}

func NewAPIException(message, apiCode string, response map[string]interface{}) *APIException {
	return &APIException{
		BingXException: NewBingXException(message, 0, response),
		APICode:        apiCode,
	}
}

type AuthenticationException struct {
	*BingXException
}

func (e *AuthenticationException) Error() string {
	return fmt.Sprintf("Authentication Error: %s", e.Message)
}

func NewAuthenticationException(message string, response map[string]interface{}) *AuthenticationException {
	return &AuthenticationException{
		BingXException: NewBingXException(message, 0, response),
	}
}

type RateLimitException struct {
	*BingXException
}

func (e *RateLimitException) Error() string {
	return fmt.Sprintf("Rate Limit Exceeded: %s", e.Message)
}

func NewRateLimitException(message string, response map[string]interface{}) *RateLimitException {
	return &RateLimitException{
		BingXException: NewBingXException(message, 0, response),
	}
}

type InsufficientBalanceException struct {
	*BingXException
}

func (e *InsufficientBalanceException) Error() string {
	return fmt.Sprintf("Insufficient Balance: %s", e.Message)
}

func NewInsufficientBalanceException(message string, response map[string]interface{}) *InsufficientBalanceException {
	return &InsufficientBalanceException{
		BingXException: NewBingXException(message, 0, response),
	}
}
