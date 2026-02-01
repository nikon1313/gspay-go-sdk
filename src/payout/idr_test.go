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

package payout

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
	t.Run("creates payout successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Contains(t, r.URL.Path, "/idr/payout")

			var req idrAPIRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "TXN123456789", req.TransactionID)
			assert.Equal(t, "user123", req.Username)
			assert.Equal(t, "John Doe", req.AccountName)
			assert.Equal(t, "1234567890", req.AccountNumber)
			assert.Equal(t, int64(50000), req.Amount)
			assert.Equal(t, "BCA", req.BankTarget)
			assert.NotEmpty(t, req.Signature)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayout_id":123,"status":0}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		resp, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        50000,
			BankCode:      "BCA",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, json.Number("123"), resp.IDRPayoutID)
	})

	t.Run("validates bank code", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        50000,
			BankCode:      "INVALID",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrInvalidBankCode)
	})

	t.Run("validates minimum amount", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        5000, // Less than 10000
			BankCode:      "BCA",
		})

		require.Error(t, err)
		valErr := errors.GetValidationError(err)
		require.NotNil(t, valErr)
		assert.Equal(t, "amount", valErr.Field)
	})

	t.Run("normalizes bank code to uppercase", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req idrAPIRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "BCA", req.BankTarget) // Should be uppercase

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayout_id":123,"status":0}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		_, err := svc.Create(t.Context(), &IDRRequest{
			TransactionID: "TXN123456789",
			Username:      "user123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        50000,
			BankCode:      "bca", // lowercase
		})

		require.NoError(t, err)
	})
}

func TestIDRService_GetStatus(t *testing.T) {
	t.Run("gets payout status successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "TXN123456789", r.URL.Query().Get("transaction_id"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"idrpayout_id":123,"transaction_id":"TXN123456789","account_name":"John Doe","account_number":"1234567890","amount":50000.00,"status":1,"completed":true,"payout_success":true,"remark":"success","signature":"sig"}`,
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewIDRService(c)

		resp, err := svc.GetStatus(t.Context(), "TXN123456789")

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, json.Number("123"), resp.IDRPayoutID)
		assert.Equal(t, "TXN123456789", resp.TransactionID)
		assert.Equal(t, "John Doe", resp.AccountName)
		assert.Equal(t, "1234567890", resp.AccountNumber)
		assert.Equal(t, json.Number("50000.00"), resp.Amount)
		assert.Equal(t, constants.StatusSuccess, resp.Status)
		assert.True(t, resp.Completed)
		assert.True(t, resp.PayoutSuccess)
		assert.Equal(t, "success", resp.Remark)
		assert.Equal(t, "sig", resp.Signature)
	})
}

func TestIDRService_VerifyCallback(t *testing.T) {
	c := client.New("auth-key", "test-secret-key")
	svc := NewIDRService(c)

	t.Run("verifies valid callback signature", func(t *testing.T) {
		// Generate valid signature: idrpayout_id + account_number + amount + transaction_id + secret_key
		callback := &IDRCallback{
			IDRPayoutID:   "123",
			TransactionID: "TXN123456789",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        "50000.00",
			Completed:     true,
			PayoutSuccess: true,
			Remark:        "Payment completed successfully",
			Signature:     signature.Generate("123123456789050000.00TXN123456789test-secret-key"),
		}

		err := svc.VerifyCallback(callback)
		assert.NoError(t, err)
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		callback := &IDRCallback{
			IDRPayoutID:   "123",
			TransactionID: "TXN123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        "50000.00",
			Completed:     false,
			PayoutSuccess: false,
			Remark:        "Payment failed",
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
				name: "missing idrpayout_id",
				callback: &IDRCallback{
					TransactionID: "TXN123456789",
					AccountName:   "John Doe",
					AccountNumber: "1234567890",
					Amount:        "50000.00",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
					Signature:     "sig",
				},
			},
			{
				name: "missing account_number",
				callback: &IDRCallback{
					IDRPayoutID:   "123",
					TransactionID: "TXN123456789",
					AccountName:   "John Doe",
					Amount:        "50000.00",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
					Signature:     "sig",
				},
			},
			{
				name: "missing account_number",
				callback: &IDRCallback{
					IDRPayoutID:   "123",
					TransactionID: "TXN123456789",
					AccountName:   "John Doe",
					Amount:        "50000.00",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
					Signature:     "sig",
				},
			},
			{
				name: "missing amount",
				callback: &IDRCallback{
					IDRPayoutID:   "123",
					TransactionID: "TXN123456789",
					AccountName:   "John Doe",
					AccountNumber: "1234567890",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
					Signature:     "sig",
				},
			},
			{
				name: "missing transaction_id",
				callback: &IDRCallback{
					IDRPayoutID:   "123",
					AccountName:   "John Doe",
					AccountNumber: "1234567890",
					Amount:        "50000.00",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
					Signature:     "sig",
				},
			},
			{
				name: "missing signature",
				callback: &IDRCallback{
					IDRPayoutID:   "123",
					TransactionID: "TXN123456789",
					AccountName:   "John Doe",
					AccountNumber: "1234567890",
					Amount:        "50000.00",
					Completed:     true,
					PayoutSuccess: true,
					Remark:        "Success",
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
}

func TestIDRService_VerifyCallbackWithIP(t *testing.T) {
	t.Run("verifies callback with whitelisted IP", func(t *testing.T) {
		c := client.New("auth-key", "secret-key", client.WithCallbackIPWhitelist("192.168.1.1"))
		svc := NewIDRService(c)
		callback := &IDRCallback{
			IDRPayoutID:   "123",
			TransactionID: "TXN123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        "50000.00",
			Completed:     true,
			PayoutSuccess: true,
			Remark:        "Success",
			Signature:     signature.Generate("123123456789050000.00TXN123secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "192.168.1.1")
		assert.NoError(t, err)
	})

	t.Run("rejects callback with non-whitelisted IP", func(t *testing.T) {
		c := client.New("auth-key", "secret-key", client.WithCallbackIPWhitelist("192.168.1.1"))
		svc := NewIDRService(c)

		callback := &IDRCallback{
			IDRPayoutID:   "123",
			TransactionID: "TXN123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        "50000.00",
			Completed:     true,
			PayoutSuccess: true,
			Remark:        "Success",
			Signature:     signature.Generate("123123456789050000.00TXN123secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "192.168.1.2")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrIPNotWhitelisted)
	})

	t.Run("skips IP check when no whitelist configured", func(t *testing.T) {
		c := client.New("auth-key", "secret-key")
		svc := NewIDRService(c)

		callback := &IDRCallback{
			IDRPayoutID:   "123",
			TransactionID: "TXN123",
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			Amount:        "50000.00",
			Completed:     true,
			PayoutSuccess: true,
			Remark:        "Success",
			Signature:     signature.Generate("123123456789050000.00TXN123secret-key"),
		}

		err := svc.VerifyCallbackWithIP(callback, "any.ip")
		assert.NoError(t, err)
	})
}
