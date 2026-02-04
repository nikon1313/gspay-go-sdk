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

package payment

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/internal/signature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIDRService_Create(t *testing.T) {
	t.Run("creates payment successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Contains(t, r.URL.Path, "/idr/payment")

			var req idrAPIRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "TXN123456789", req.TransactionID)
			assert.Equal(t, "user123", req.Username)
			assert.Equal(t, int64(50000), req.Amount)
			assert.Equal(t, "QRIS", req.Channel)
			assert.NotEmpty(t, req.Signature)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayment_id":"PAY123","transaction_id":"TXN123456789","amount":"50000","expire_date":"2026-01-26 15:00:00","status":"0","payment_url":"https://pay.example.com"}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		resp, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			Amount:        50000,
			Channel:       constants.ChannelQRIS,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "https://pay.example.com", resp.PaymentURL)
		assert.Equal(t, "PAY123", resp.IDRPaymentID)
	})

	t.Run("validates transaction ID length", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		// Too short
		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN",
			Username:      "user123",
			Amount:        50000,
		})
		valErr := errors.GetValidationError(err)
		require.NotNil(t, valErr)
		assert.Equal(t, "transaction_id", valErr.Field)

		// Too long
		_, err = svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN12345678901234567890",
			Username:      "user123",
			Amount:        50000,
		})
		valErr = errors.GetValidationError(err)
		require.NotNil(t, valErr)
		assert.Equal(t, "transaction_id", valErr.Field)
	})

	t.Run("validates minimum amount", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			Amount:        5000, // Less than 10000
		})

		require.Error(t, err)
		valErr := errors.GetValidationError(err)
		require.NotNil(t, valErr)
		assert.Equal(t, "amount", valErr.Field)
	})

	t.Run("handles API error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    400,
				"message": "invalid signature",
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			Amount:        50000,
		})

		require.Error(t, err)
		apiErr := errors.GetAPIError(err)
		require.NotNil(t, apiErr)
		assert.Equal(t, 400, apiErr.Code)
		assert.Equal(t, "invalid signature", apiErr.Message)
	})

	t.Run("ignores invalid channel", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req idrAPIRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Empty(t, req.Channel, "channel should be empty")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayment_id":"PAY123","transaction_id":"TXN123456789","amount":"50000","expire_date":"2026-01-26 15:00:00","status":"0","payment_url":"https://pay.example.com"}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			Amount:        50000,
			Channel:       "INVALID_CHANNEL",
		})

		require.NoError(t, err)
	})

	t.Run("normalizes channel to uppercase", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req idrAPIRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "QRIS", req.Channel) // Should be uppercase

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayment_id":"PAY123","transaction_id":"TXN123456789","amount":"50000","expire_date":"2026-01-26 15:00:00","status":"0","payment_url":"https://pay.example.com"}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			Amount:        50000,
			Channel:       "qris", // lowercase
		})

		require.NoError(t, err)
	})
}

func TestIDRService_GetStatus(t *testing.T) {
	t.Run("gets payment status successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "TXN123456789", r.URL.Query().Get("transaction_id"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayment_id":123,"transaction_id":"TXN123456789","player_username":"demo_user","status":1,"amount":50000.00,"completed":true,"success":true,"remark":"success","signature":"sig"}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		resp, err := svc.GetStatus(t.Context(), "TXN123456789")

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, json.Number("123"), resp.IDRPaymentID)
		assert.Equal(t, "TXN123456789", resp.TransactionID)
		assert.Equal(t, "demo_user", resp.PlayerUsername)
		assert.Equal(t, constants.StatusSuccess, resp.Status)
		assert.Equal(t, json.Number("50000.00"), resp.Amount)
		assert.True(t, resp.Completed)
		assert.True(t, resp.Success)
		assert.Equal(t, "success", resp.Remark)
		assert.Equal(t, "sig", resp.Signature)
	})
}

func TestIDRService_VerifyStatusSignature(t *testing.T) {
	c := client.New("auth-key", "test-secret-key")
	svc := NewIDRService(c)

	t.Run("verifies valid status signature", func(t *testing.T) {
		status := &IDRStatusResponse{
			IDRPaymentID:   "123",
			TransactionID:  "TXN123456789",
			PlayerUsername: "demo_user",
			Status:         1,
			Amount:         "50000.00",
			Completed:      true,
			Success:        true,
			Remark:         "success",
		}
		// Generate correct signature
		status.Signature = signature.Generate("12350000.00TXN1234567891test-secret-key")

		err := svc.VerifyStatusSignature(status)
		assert.NoError(t, err)
	})

	t.Run("rejects invalid status signature", func(t *testing.T) {
		status := &IDRStatusResponse{
			IDRPaymentID:   "123",
			TransactionID:  "TXN123456789",
			PlayerUsername: "demo_user",
			Status:         1,
			Amount:         "50000.00",
			Completed:      true,
			Success:        true,
			Remark:         "success",
			Signature:      "invalid",
		}

		err := svc.VerifyStatusSignature(status)
		assert.ErrorIs(t, err, errors.ErrInvalidSignature)
	})
}

func TestIDRService_VerifyCallback(t *testing.T) {
	c := client.New("auth-key", "test-secret-key")
	svc := NewIDRService(c)

	t.Run("verifies valid callback signature", func(t *testing.T) {
		// Generate valid signature
		signatureData := "PAY12350000.00TXN1234567891test-secret-key"
		validSignature := signature.Generate(signatureData)

		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "50000.00",
			TransactionID: "TXN123456789",
			Status:        constants.StatusSuccess,
			Signature:     validSignature,
		}

		err := svc.VerifyCallback(callback)
		assert.NoError(t, err)
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "50000.00",
			TransactionID: "TXN123456789",
			Status:        constants.StatusSuccess,
			Signature:     "invalid-signature",
		}

		err := svc.VerifyCallback(callback)
		assert.ErrorIs(t, err, errors.ErrInvalidSignature)
	})

	t.Run("rejects missing required fields", func(t *testing.T) {
		testCases := []struct {
			name     string
			callback *IDRCallback
		}{
			{
				name: "missing idrpayment_id",
				callback: &IDRCallback{
					Amount:        "50000.00",
					TransactionID: "TXN123456789",
					Status:        1,
					Signature:     "sig",
				},
			},
			{
				name: "missing amount",
				callback: &IDRCallback{
					IDRPaymentID:  "PAY123",
					TransactionID: "TXN123456789",
					Status:        1,
					Signature:     "sig",
				},
			},
			{
				name: "missing transaction_id",
				callback: &IDRCallback{
					IDRPaymentID: "PAY123",
					Amount:       "50000.00",
					Status:       1,
					Signature:    "sig",
				},
			},
			{
				name: "missing signature",
				callback: &IDRCallback{
					IDRPaymentID:  "PAY123",
					Amount:        "50000.00",
					TransactionID: "TXN123456789",
					Status:        1,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := svc.VerifyCallback(tc.callback)
				assert.ErrorIs(t, err, errors.ErrMissingCallbackField)
			})
		}
	})

	t.Run("rejects invalid amount format", func(t *testing.T) {
		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "invalid",
			TransactionID: "TXN123456789",
			Status:        1,
			Signature:     "sig",
		}

		err := svc.VerifyCallback(callback)
		valErr := errors.GetValidationError(err)
		require.NotNil(t, valErr)
		assert.Equal(t, "amount", valErr.Field)
	})
}

func TestIDRService_VerifyCallbackWithIP(t *testing.T) {
	t.Run("verifies callback with whitelisted IP", func(t *testing.T) {
		c := client.New("auth-key", "secret-key", client.WithCallbackIPWhitelist("192.168.1.1"))
		svc := NewIDRService(c)

		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "50000.00",
			TransactionID: "TXN123",
			Status:        1,
			Signature:     signature.Generate("PAY12350000.00TXN1231secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "192.168.1.1")
		assert.NoError(t, err)
	})

	t.Run("rejects callback with non-whitelisted IP", func(t *testing.T) {
		c := client.New("auth-key", "secret-key", client.WithCallbackIPWhitelist("192.168.1.1"))
		svc := NewIDRService(c)

		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "50000.00",
			TransactionID: "TXN123",
			Status:        1,
			Signature:     signature.Generate("PAY12350000.00TXN1231secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "192.168.1.2")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrIPNotWhitelisted)
	})

	t.Run("skips IP check when no whitelist configured", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		callback := &IDRCallback{
			IDRPaymentID:  "PAY123",
			Amount:        "50000.00",
			TransactionID: "TXN123",
			Status:        1,
			Signature:     signature.Generate("PAY12350000.00TXN1231secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "any.ip")
		assert.NoError(t, err)
	})
}
