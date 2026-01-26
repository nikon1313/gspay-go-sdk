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

package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates client with default values", func(t *testing.T) {
		c := New("auth-key", "secret-key")

		assert.Equal(t, "auth-key", c.AuthKey)
		assert.Equal(t, "secret-key", c.SecretKey)
		assert.Equal(t, constants.DefaultBaseURL, c.BaseURL)
		assert.Equal(t, time.Duration(constants.DefaultTimeout)*time.Second, c.Timeout)
		assert.Equal(t, constants.DefaultRetries, c.Retries)
		assert.NotNil(t, c.HTTPClient)
	})

	t.Run("applies custom options", func(t *testing.T) {
		customHTTPClient := &http.Client{Timeout: 60 * time.Second}

		c := New(
			"auth-key",
			"secret-key",
			WithBaseURL("https://custom.api.com/"),
			WithTimeout(60*time.Second),
			WithRetries(5),
			WithRetryWait(1*time.Second, 5*time.Second),
			WithHTTPClient(customHTTPClient),
		)

		assert.Equal(t, "https://custom.api.com", c.BaseURL) // trailing slash removed
		assert.Equal(t, 60*time.Second, c.Timeout)
		assert.Equal(t, 5, c.Retries)
		assert.Equal(t, 1*time.Second, c.RetryWaitMin)
		assert.Equal(t, 5*time.Second, c.RetryWaitMax)
		assert.Same(t, customHTTPClient, c.HTTPClient)
	})

	t.Run("ignores invalid timeout", func(t *testing.T) {
		c := New("auth-key", "secret-key", WithTimeout(1*time.Second))
		assert.Equal(t, time.Duration(constants.DefaultTimeout)*time.Second, c.Timeout)
	})

	t.Run("ignores negative retries", func(t *testing.T) {
		c := New("auth-key", "secret-key", WithRetries(-1))
		assert.Equal(t, constants.DefaultRetries, c.Retries)
	})
}

func TestGenerateSignature(t *testing.T) {
	c := New("auth-key", "secret-key")

	t.Run("generates correct MD5 signature", func(t *testing.T) {
		// MD5("test") = 098f6bcd4621d373cade4e832627b4f6
		sig := c.GenerateSignature("test")
		assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", sig)
	})

	t.Run("generates consistent signatures", func(t *testing.T) {
		sig1 := c.GenerateSignature("payment123user45000secret")
		sig2 := c.GenerateSignature("payment123user45000secret")
		assert.Equal(t, sig1, sig2)
	})
}

func TestVerifySignature(t *testing.T) {
	c := New("auth-key", "secret-key")

	t.Run("returns true for matching signatures", func(t *testing.T) {
		sig := c.GenerateSignature("test")
		assert.True(t, c.VerifySignature(sig, sig))
	})

	t.Run("returns false for non-matching signatures", func(t *testing.T) {
		sig := c.GenerateSignature("test")
		assert.False(t, c.VerifySignature(sig, "invalid"))
	})
}

func TestGenerateTransactionID(t *testing.T) {
	t.Run("generates transaction ID with prefix", func(t *testing.T) {
		txnID := GenerateTransactionID("TXN")
		require.NotEmpty(t, txnID)
		assert.True(t, len(txnID) <= 20)
		assert.Contains(t, txnID, "TXN")
	})

	t.Run("truncates long prefix", func(t *testing.T) {
		txnID := GenerateTransactionID("VERYLONGPREFIX")
		require.NotEmpty(t, txnID)
		assert.True(t, len(txnID) <= 20)
		assert.Contains(t, txnID, "VER")
	})

	t.Run("generates valid format IDs", func(t *testing.T) {
		// Generate multiple IDs and verify format
		for i := 0; i < 10; i++ {
			id := GenerateTransactionID("TXN")
			assert.True(t, len(id) >= 17 && len(id) <= 20, "ID length should be 17-20, got %d", len(id))
			assert.Regexp(t, `^TXN\d{14,17}$`, id)
		}
	})
}

func TestBuildReturnURL(t *testing.T) {
	t.Run("appends return parameter with ?", func(t *testing.T) {
		result := BuildReturnURL("https://pay.example.com/payment/123", "https://mysite.com/complete")
		assert.Equal(t, "https://pay.example.com/payment/123?return=https%3A%2F%2Fmysite.com%2Fcomplete", result)
	})

	t.Run("appends return parameter with &", func(t *testing.T) {
		result := BuildReturnURL("https://pay.example.com/payment/123?foo=bar", "https://mysite.com/complete")
		assert.Equal(t, "https://pay.example.com/payment/123?foo=bar&return=https%3A%2F%2Fmysite.com%2Fcomplete", result)
	})
}

func TestFormatAmountIDR(t *testing.T) {
	tests := []struct {
		amount   int64
		expected string
	}{
		{100, "Rp 100"},
		{1000, "Rp 1.000"},
		{10000, "Rp 10.000"},
		{50000, "Rp 50.000"},
		{100000, "Rp 100.000"},
		{1000000, "Rp 1.000.000"},
		{10000000, "Rp 10.000.000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatAmountIDR(tt.amount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatAmountUSDT(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{1.00, "1.00 USDT"},
		{10.50, "10.50 USDT"},
		{100.00, "100.00 USDT"},
		{1234.56, "1234.56 USDT"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatAmountUSDT(tt.amount)
			assert.Equal(t, tt.expected, result)
		})
	}
}
