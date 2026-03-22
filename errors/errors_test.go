package errors

import (
	"testing"
)

func TestNewBingXException(t *testing.T) {
	message := "Test error message"
	code := 500
	response := map[string]interface{}{"error": "test"}

	err := NewBingXException(message, code, response)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}

	if err.Code != code {
		t.Errorf("Expected code %d, got %d", code, err.Code)
	}

	if err.Response == nil {
		t.Error("Expected response to be set")
	}
}

func TestBingXExceptionError(t *testing.T) {
	message := "Test error"
	err := NewBingXException(message, 500, nil)

	errorString := err.Error()
	if errorString != message {
		t.Errorf("Expected error string '%s', got '%s'", message, errorString)
	}
}

func TestNewAPIException(t *testing.T) {
	message := "API error"
	code := "100001"
	response := map[string]interface{}{"code": code}

	err := NewAPIException(message, code, response)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}

	if err.APICode != code {
		t.Errorf("Expected code '%s', got '%s'", code, err.APICode)
	}
}

func TestNewAuthenticationException(t *testing.T) {
	message := "Authentication failed"
	response := map[string]interface{}{"code": "100001"}

	err := NewAuthenticationException(message, response)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}
}

func TestNewRateLimitException(t *testing.T) {
	message := "Rate limit exceeded"
	response := map[string]interface{}{"code": "100005"}

	err := NewRateLimitException(message, response)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}
}

func TestNewInsufficientBalanceException(t *testing.T) {
	message := "Insufficient balance"
	response := map[string]interface{}{"code": "200001"}

	err := NewInsufficientBalanceException(message, response)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}
}
