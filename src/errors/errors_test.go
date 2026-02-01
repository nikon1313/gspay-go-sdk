// Copyright 2026 H0llyW00dzZ
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"testing"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	t.Run("formats error without endpoint", func(t *testing.T) {
		err := &APIError{
			Code:    400,
			Message: "bad request",
		}
		expected := "gspay: API error 400: bad request"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("formats error with endpoint", func(t *testing.T) {
		err := &APIError{
			Code:     500,
			Message:  "internal server error",
			Endpoint: "/api/test",
		}
		expected := "gspay: API error 500 on /api/test: internal server error"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("sanitizes auth key in endpoint", func(t *testing.T) {
		err := &APIError{
			Code:     401,
			Message:  "unauthorized",
			Endpoint: "/v2/integrations/operators/secretkey123/idr/payment",
		}
		expected := "gspay: API error 401 on /v2/integrations/operators/[REDACTED]/idr/payment: unauthorized"
		assert.Equal(t, expected, err.Error())
	})
}

func TestIsAPIError(t *testing.T) {
	t.Run("returns true for APIError", func(t *testing.T) {
		err := &APIError{Code: 400, Message: "test"}
		assert.True(t, IsAPIError(err))
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("regular error")
		assert.False(t, IsAPIError(err))
	})

	t.Run("returns false for nil", func(t *testing.T) {
		assert.False(t, IsAPIError(nil))
	})
}

func TestGetAPIError(t *testing.T) {
	t.Run("extracts APIError", func(t *testing.T) {
		original := &APIError{Code: 404, Message: "not found"}
		wrapped := GetAPIError(original)
		assert.Equal(t, original, wrapped)
	})

	t.Run("returns nil for non-APIError", func(t *testing.T) {
		err := errors.New("regular error")
		assert.Nil(t, GetAPIError(err))
	})
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "amount",
		Message: "must be positive",
	}
	expected := "gspay: validation error for amount: must be positive"
	assert.Equal(t, expected, err.Error())
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("bank_code", "invalid code")
	assert.Equal(t, "bank_code", err.Field)
	assert.Equal(t, "invalid code", err.Message)
}

func TestIsValidationError(t *testing.T) {
	t.Run("returns true for ValidationError", func(t *testing.T) {
		err := &ValidationError{Field: "test", Message: "error"}
		assert.True(t, IsValidationError(err))
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("regular error")
		assert.False(t, IsValidationError(err))
	})
}

func TestGetValidationError(t *testing.T) {
	t.Run("extracts ValidationError", func(t *testing.T) {
		original := &ValidationError{Field: "test", Message: "error"}
		extracted := GetValidationError(original)
		assert.Equal(t, original, extracted)
	})

	t.Run("returns nil for non-ValidationError", func(t *testing.T) {
		err := errors.New("regular error")
		assert.Nil(t, GetValidationError(err))
	})
}

func TestSentinelErrors(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"ErrInvalidTransactionID", ErrInvalidTransactionID},
		{"ErrInvalidAmount", ErrInvalidAmount},
		{"ErrInvalidBankCode", ErrInvalidBankCode},
		{"ErrInvalidSignature", ErrInvalidSignature},
		{"ErrMissingCallbackField", ErrMissingCallbackField},
		{"ErrEmptyResponse", ErrEmptyResponse},
		{"ErrInvalidJSON", ErrInvalidJSON},
		{"ErrRequestFailed", ErrRequestFailed},
		{"ErrIPNotWhitelisted", ErrIPNotWhitelisted},
		{"ErrInvalidIPAddress", ErrInvalidIPAddress},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.err)
			assert.NotEmpty(t, tc.err.Error())
		})
	}
}

func TestLocalizedError_Error(t *testing.T) {
	t.Run("returns English message", func(t *testing.T) {
		err := NewLocalizedError(i18n.English, KeyInvalidTransactionID)
		assert.Equal(t, "transaction ID must be 5-20 characters", err.Error())
	})

	t.Run("returns Indonesian message", func(t *testing.T) {
		err := NewLocalizedError(i18n.Indonesian, KeyInvalidTransactionID)
		assert.Equal(t, "ID transaksi harus 5-20 karakter", err.Error())
	})

	t.Run("falls back to English for unknown language", func(t *testing.T) {
		err := NewLocalizedError(i18n.Language("fr"), KeyInvalidAmount)
		assert.Equal(t, "invalid payment amount", err.Error())
	})
}

func TestLocalizedError_Key(t *testing.T) {
	err := NewLocalizedError(i18n.English, KeyMinAmountIDR)
	assert.Equal(t, KeyMinAmountIDR, err.Key())
}

func TestNewLocalizedError(t *testing.T) {
	err := NewLocalizedError(i18n.Indonesian, KeyInvalidSignature)
	assert.NotNil(t, err)
	assert.Equal(t, KeyInvalidSignature, err.Key())
	assert.Equal(t, "tanda tangan tidak valid", err.Error())
}

func TestIsLocalizedError(t *testing.T) {
	t.Run("returns true for LocalizedError", func(t *testing.T) {
		err := NewLocalizedError(i18n.English, KeyRequestFailed)
		assert.True(t, IsLocalizedError(err))
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("regular error")
		assert.False(t, IsLocalizedError(err))
	})

	t.Run("returns false for nil", func(t *testing.T) {
		assert.False(t, IsLocalizedError(nil))
	})
}

func TestGetLocalizedError(t *testing.T) {
	t.Run("extracts LocalizedError", func(t *testing.T) {
		original := NewLocalizedError(i18n.Indonesian, KeyEmptyResponse)
		extracted := GetLocalizedError(original)
		assert.Equal(t, original, extracted)
	})

	t.Run("returns nil for non-LocalizedError", func(t *testing.T) {
		err := errors.New("regular error")
		assert.Nil(t, GetLocalizedError(err))
	})
}

func TestGetMessage(t *testing.T) {
	t.Run("returns English message", func(t *testing.T) {
		msg := GetMessage(i18n.English, KeyMinAmountIDR)
		assert.Equal(t, "minimum amount is 10000 IDR", msg)
	})

	t.Run("returns Indonesian message", func(t *testing.T) {
		msg := GetMessage(i18n.Indonesian, KeyMinPayoutAmountIDR)
		assert.Equal(t, "jumlah pembayaran minimum adalah 10000 IDR", msg)
	})

	t.Run("falls back to English for unknown language", func(t *testing.T) {
		msg := GetMessage(i18n.Language("de"), KeyInvalidJSON)
		assert.Equal(t, "invalid JSON response", msg)
	})
}

func TestLocalizedErrorMessageKeys(t *testing.T) {
	// Verify all re-exported keys work correctly
	testCases := []struct {
		key      MessageKey
		expected string
	}{
		{KeyInvalidTransactionID, "transaction ID must be 5-20 characters"},
		{KeyInvalidAmount, "invalid payment amount"},
		{KeyInvalidBankCode, "invalid bank code"},
		{KeyInvalidSignature, "invalid signature"},
		{KeyMissingCallbackField, "missing required callback field"},
		{KeyEmptyResponse, "empty response from API"},
		{KeyInvalidJSON, "invalid JSON response"},
		{KeyRequestFailed, "request failed"},
		{KeyIPNotWhitelisted, "IP address not whitelisted"},
		{KeyInvalidIPAddress, "invalid IP address format"},
		{KeyMinAmountIDR, "minimum amount is 10000 IDR"},
		{KeyMinAmountUSDT, "minimum amount is 1.00 USDT"},
		{KeyMinPayoutAmountIDR, "minimum payout amount is 10000 IDR"},
		{KeyInvalidAmountFormat, "invalid amount format"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.key), func(t *testing.T) {
			msg := GetMessage(i18n.English, tc.key)
			assert.Equal(t, tc.expected, msg)
		})
	}
}
