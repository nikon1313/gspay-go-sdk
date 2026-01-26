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

// Package client provides the HTTP client for the GSPAY2 API.
package client

import (
	"net/http"
	"strings"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
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
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout >= 5*time.Second {
			c.Timeout = timeout
		}
	}
}

// WithRetries sets the number of retry attempts for transient failures.
func WithRetries(retries int) Option {
	return func(c *Client) {
		if retries >= 0 {
			c.Retries = retries
		}
	}
}

// WithRetryWait sets the minimum and maximum wait times between retries.
func WithRetryWait(min, max time.Duration) Option {
	return func(c *Client) {
		c.RetryWaitMin = min
		c.RetryWaitMax = max
	}
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
func (c *Client) GenerateSignature(data string) string {
	return signature.Generate(data)
}

// VerifySignature verifies a callback signature.
func (c *Client) VerifySignature(expected, actual string) bool {
	return signature.Verify(expected, actual)
}
