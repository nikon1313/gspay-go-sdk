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
	"net"
	"net/http"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client/logger"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/internal/signature"
)

// Client is the GSPAY2 API client.
type Client struct {
	// AuthKey is the operator authentication key (used in URL path).
	AuthKey string
	// SecretKey is the operator secret key (used for signature generation).
	SecretKey string
	// BaseURL is the API base URL.
	BaseURL string
	// HTTPClient is the underlying HTTP client.
	HTTPClient *http.Client
	// Timeout is the request timeout duration.
	Timeout time.Duration
	// Retries is the number of retry attempts for transient failures.
	Retries int
	// RetryWaitMin is the minimum wait time between retries.
	RetryWaitMin time.Duration
	// RetryWaitMax is the maximum wait time between retries.
	RetryWaitMax time.Duration
	// CallbackIPWhitelist contains allowed IP addresses/CIDR ranges for callbacks.
	// If empty, IP validation is skipped.
	CallbackIPWhitelist []string
	// parsedIPNets contains parsed CIDR networks for efficient IP checking.
	parsedIPNets []*net.IPNet
	// Debug enables debug logging of API requests and responses.
	// When true and no custom logger is set, uses the default logger.
	Debug bool
	// parsedIPs contains parsed individual IP addresses.
	parsedIPs []net.IP
	// Language is the language for SDK error messages.
	// Default is English.
	Language i18n.Language
	// logger is the structured logger for the client.
	// Default is logger.Nop (no logging).
	logger logger.Handler
}

// Logger returns the configured logger instance.
//
// This allows services and other packages to use the same logger configuration
// as the client for consistent logging across the SDK.
func (c *Client) Logger() logger.Handler {
	return c.logger
}

// I18n returns the localized message for the given key using the client's language setting.
// This is a convenience method for i18n support in logging and error messages.
func (c *Client) I18n(key i18n.MessageKey) string {
	return i18n.Get(c.Language, key)
}

// Error creates a localized error wrapping the provided sentinel error.
// This is a convenience method that uses the client's language setting.
func (c *Client) Error(sentinel error, args ...any) error {
	return errors.New(c.Language, sentinel, args...)
}

// New creates a new GSPAY2 API client.
//
// Parameters:
//   - authKey: Operator authentication key (used in URL path)
//   - secretKey: Operator secret key (used for signature generation)
//   - opts: Optional configuration options
func New(authKey, secretKey string, opts ...Option) *Client {
	c := &Client{
		AuthKey:      authKey,
		SecretKey:    secretKey,
		BaseURL:      constants.DefaultBaseURL,
		Timeout:      time.Duration(constants.DefaultTimeout) * time.Second,
		Retries:      constants.DefaultRetries,
		RetryWaitMin: time.Duration(constants.DefaultRetryWaitMin) * time.Millisecond,
		RetryWaitMax: time.Duration(constants.DefaultRetryWaitMax) * time.Millisecond,
		Language:     i18n.English,
		logger:       logger.Nop{},
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: c.Timeout,
		}
	}

	return c
}

// GenerateSignature generates an MD5 signature for API requests.
func (c *Client) GenerateSignature(data string) string { return signature.Generate(data) }

// VerifySignature verifies a callback signature.
func (c *Client) VerifySignature(expected, actual string) bool {
	return signature.Verify(expected, actual)
}

// parseIPWhitelist parses the IP whitelist into net.IP and net.IPNet for efficient checking.
func (c *Client) parseIPWhitelist() {
	c.parsedIPNets = nil
	c.parsedIPs = nil

	for _, ipStr := range c.CallbackIPWhitelist {
		// Try parsing as CIDR first
		if _, ipNet, err := net.ParseCIDR(ipStr); err == nil {
			c.parsedIPNets = append(c.parsedIPNets, ipNet)
			continue
		}

		// Try parsing as individual IP
		if ip := net.ParseIP(ipStr); ip != nil {
			c.parsedIPs = append(c.parsedIPs, ip)
		}
	}
}

// IsIPWhitelisted checks if the given IP address is in the whitelist.
//
// Returns true if:
//   - The whitelist is empty (IP validation disabled)
//   - The IP matches an individual whitelisted IP
//   - The IP falls within a whitelisted CIDR range
//
// The ipStr parameter can include a port (e.g., "192.168.1.1:8080"),
// which will be automatically stripped before validation.
func (c *Client) IsIPWhitelisted(ipStr string) bool {
	// If no whitelist configured, allow all IPs
	if len(c.CallbackIPWhitelist) == 0 {
		return true
	}

	// Strip port if present (handles both IPv4 and IPv6)
	host := ipStr
	if h, _, err := net.SplitHostPort(ipStr); err == nil {
		host = h
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	// Check individual IPs
	for _, whitelistedIP := range c.parsedIPs {
		if whitelistedIP.Equal(ip) {
			return true
		}
	}

	// Check CIDR ranges
	for _, ipNet := range c.parsedIPNets {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// VerifyCallbackIP verifies that the callback request originates from a whitelisted IP.
//
// Returns nil if the IP is whitelisted or if the whitelist is empty.
// Returns ErrIPNotWhitelisted if the IP is not in the whitelist.
// Returns ErrInvalidIPAddress if the IP address format is invalid.
func (c *Client) VerifyCallbackIP(ipStr string) error {
	// If no whitelist configured, skip IP validation
	if len(c.CallbackIPWhitelist) == 0 {
		return nil
	}

	// Strip port if present
	host := ipStr
	if h, _, err := net.SplitHostPort(ipStr); err == nil {
		host = h
	}

	// Validate IP format
	if net.ParseIP(host) == nil {
		return c.Error(errors.ErrInvalidIPAddress)
	}

	// Check whitelist
	if !c.IsIPWhitelisted(ipStr) {
		return c.Error(errors.ErrIPNotWhitelisted)
	}

	return nil
}
