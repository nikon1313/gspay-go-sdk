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
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
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
		for range 10 {
			id := GenerateTransactionID("TXN")
			assert.True(t, len(id) >= 17 && len(id) <= 20, "ID length should be 17-20, got %d", len(id))
			assert.Regexp(t, `^TXN\d{14,17}$`, id)
		}
	})
}

func TestGenerateUUIDTransactionID(t *testing.T) {
	t.Run("generates UUID transaction ID with prefix", func(t *testing.T) {
		txnID := GenerateUUIDTransactionID("TXN")
		require.NotEmpty(t, txnID)
		assert.True(t, len(txnID) <= 20)
		assert.Contains(t, txnID, "TXN")
	})

	t.Run("truncates long prefix", func(t *testing.T) {
		txnID := GenerateUUIDTransactionID("VERYLONGPREFIX")
		require.NotEmpty(t, txnID)
		assert.True(t, len(txnID) <= 20)
		assert.Contains(t, txnID, "VER")
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		ids := make(map[string]bool)
		for range 100 {
			id := GenerateUUIDTransactionID("TXN")
			assert.False(t, ids[id], "Generated ID should be unique")
			ids[id] = true
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
		{100000000, "Rp 100.000.000"},
		{1000000000, "Rp 1.000.000.000"},
		{1000000000000, "Rp 1.000.000.000.000"},
		{1234567890, "Rp 1.234.567.890"},
		{9999999999999, "Rp 9.999.999.999.999"},
		{19999999999999, "Rp 19.999.999.999.999"},
		{109999999999999, "Rp 109.999.999.999.999"},
		{1109999999999999, "Rp 1.109.999.999.999.999"},
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

func TestWithCallbackIPWhitelist(t *testing.T) {
	t.Run("configures IP whitelist", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist(
			"192.168.1.1",
			"10.0.0.0/8",
		))

		assert.Len(t, c.CallbackIPWhitelist, 2)
		assert.Contains(t, c.CallbackIPWhitelist, "192.168.1.1")
		assert.Contains(t, c.CallbackIPWhitelist, "10.0.0.0/8")
	})

	t.Run("parses IPs and CIDRs correctly", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist(
			"192.168.1.1",
			"10.0.0.0/8",
			"2001:db8::1",
			"2001:db8::/32",
		))

		assert.Len(t, c.parsedIPs, 2)    // 192.168.1.1 and 2001:db8::1
		assert.Len(t, c.parsedIPNets, 2) // 10.0.0.0/8 and 2001:db8::/32
	})
}

func TestIsIPWhitelisted(t *testing.T) {
	t.Run("allows all IPs when whitelist is empty", func(t *testing.T) {
		c := New("auth", "secret")

		assert.True(t, c.IsIPWhitelisted("192.168.1.1"))
		assert.True(t, c.IsIPWhitelisted("10.0.0.1"))
		assert.True(t, c.IsIPWhitelisted("any.invalid.ip"))
	})

	t.Run("matches exact IP address", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1", "10.0.0.1"))

		assert.True(t, c.IsIPWhitelisted("192.168.1.1"))
		assert.True(t, c.IsIPWhitelisted("10.0.0.1"))
		assert.False(t, c.IsIPWhitelisted("192.168.1.2"))
		assert.False(t, c.IsIPWhitelisted("172.16.0.1"))
	})

	t.Run("matches CIDR range", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("10.0.0.0/8", "192.168.0.0/16"))

		assert.True(t, c.IsIPWhitelisted("10.0.0.1"))
		assert.True(t, c.IsIPWhitelisted("10.255.255.255"))
		assert.True(t, c.IsIPWhitelisted("192.168.1.1"))
		assert.True(t, c.IsIPWhitelisted("192.168.255.255"))
		assert.False(t, c.IsIPWhitelisted("172.16.0.1"))
		assert.False(t, c.IsIPWhitelisted("11.0.0.1"))
	})

	t.Run("handles IP with port", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		assert.True(t, c.IsIPWhitelisted("192.168.1.1:8080"))
		assert.True(t, c.IsIPWhitelisted("192.168.1.1:443"))
		assert.False(t, c.IsIPWhitelisted("192.168.1.2:8080"))
	})

	t.Run("handles IPv6 addresses", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("2001:db8::1", "fe80::/10"))

		assert.True(t, c.IsIPWhitelisted("2001:db8::1"))
		assert.True(t, c.IsIPWhitelisted("fe80::1"))
		assert.True(t, c.IsIPWhitelisted("fe80:0000:0000:0000:0000:0000:0000:0001"))
		assert.False(t, c.IsIPWhitelisted("2001:db8::2"))
	})

	t.Run("handles IPv6 with port", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("2001:db8::1"))

		assert.True(t, c.IsIPWhitelisted("[2001:db8::1]:8080"))
		assert.False(t, c.IsIPWhitelisted("[2001:db8::2]:8080"))
	})

	t.Run("returns false for invalid IP", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		assert.False(t, c.IsIPWhitelisted("invalid-ip"))
		assert.False(t, c.IsIPWhitelisted(""))
	})
}

func TestVerifyCallbackIP(t *testing.T) {
	t.Run("returns nil when whitelist is empty", func(t *testing.T) {
		c := New("auth", "secret")

		err := c.VerifyCallbackIP("192.168.1.1")
		assert.NoError(t, err)
	})

	t.Run("returns nil for whitelisted IP", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		err := c.VerifyCallbackIP("192.168.1.1")
		assert.NoError(t, err)
	})

	t.Run("returns error for non-whitelisted IP", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		err := c.VerifyCallbackIP("192.168.1.2")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrIPNotWhitelisted)
	})

	t.Run("returns error for invalid IP format", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		err := c.VerifyCallbackIP("invalid-ip")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrInvalidIPAddress)
	})

	t.Run("handles IP with port correctly", func(t *testing.T) {
		c := New("auth", "secret", WithCallbackIPWhitelist("192.168.1.1"))

		err := c.VerifyCallbackIP("192.168.1.1:8080")
		assert.NoError(t, err)

		err = c.VerifyCallbackIP("192.168.1.2:8080")
		assert.ErrorIs(t, err, errors.ErrIPNotWhitelisted)
	})
}
