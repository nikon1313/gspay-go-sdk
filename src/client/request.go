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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/helper/gc"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/internal/sanitize"
)

// Response represents a generic API response structure.
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// IsSuccess checks if the API response indicates success.
func (r *Response) IsSuccess() bool { return r.Code == 200 }

// LogEndpoint returns the endpoint for logging, sanitized unless debug mode is enabled.
//
// In debug mode, the full endpoint (including auth keys) is returned for troubleshooting.
// In production mode, auth keys are redacted (e.g., "/operators/[REDACTED]/idr/payment").
func (c *Client) LogEndpoint(endpoint string) string {
	if c.Debug {
		return endpoint
	}
	return sanitize.Endpoint(endpoint)
}

// LogAccountNumber returns the account number for logging, sanitized unless debug mode is enabled.
//
// In debug mode, the full account number is returned for troubleshooting.
// In production mode, only the last 4 digits are shown (e.g., "****7890").
func (c *Client) LogAccountNumber(accountNumber string) string {
	if c.Debug {
		return accountNumber
	}
	return sanitize.AccountNumber(accountNumber)
}

// LogAccountName returns the account name for logging, sanitized unless debug mode is enabled.
//
// In debug mode, the full account name is returned for troubleshooting.
// In production mode, only initials are shown (e.g., "J*** D***").
func (c *Client) LogAccountName(accountName string) string {
	if c.Debug {
		return accountName
	}
	return sanitize.AccountName(accountName)
}

// parseRetryAfter parses the Retry-After header value and returns the suggested wait duration.
// It supports both seconds format (e.g., "120") and HTTP-date format (e.g., "Wed, 21 Oct 2025 07:28:00 GMT").
// Returns 0 if the header is empty or cannot be parsed.
func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}

	// Try parsing as seconds (most common)
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		if seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
		return 0
	}

	// Try parsing as HTTP-date (RFC 1123 format)
	if t, err := time.Parse(time.RFC1123, value); err == nil {
		duration := time.Until(t)
		if duration > 0 {
			return duration
		}
	}

	return 0
}

// responseResult holds the result of processing an HTTP response.
type responseResult struct {
	Response   *Response
	Retry      bool
	RetryAfter time.Duration // Server-suggested wait time from Retry-After header (0 means use manual backoff)
	Err        error
}

// prepareRequestBody prepares the request body for HTTP requests.
func (c *Client) prepareRequestBody(body any) (io.Reader, gc.Buffer, func(), error) {
	if body == nil {
		return nil, nil, func() {}, nil
	}

	buf := gc.Default.Get()
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		buf.Reset()
		gc.Default.Put(buf)
		return nil, nil, func() {}, errors.New(c.Language, errors.ErrInvalidJSON, err)
	}

	reqBody := bytes.NewReader(buf.Bytes())
	cleanup := func() {
		buf.Reset()
		gc.Default.Put(buf)
	}

	return reqBody, buf, cleanup, nil
}

// createHTTPRequest creates an HTTP request with appropriate headers.
func (c *Client) createHTTPRequest(ctx context.Context, method, fullURL string, reqBody io.Reader, hasBody bool) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, errors.New(c.Language, errors.ErrRequestFailed, err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", constants.UserAgent())
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// processResponse processes the HTTP response and returns parsed data or error.
func (c *Client) processResponse(resp *http.Response, endpoint string) responseResult {
	defer resp.Body.Close()

	respBuf := gc.Default.Get()
	_, err := respBuf.ReadFrom(resp.Body)

	if err != nil {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return responseResult{Retry: true, Err: errors.New(c.Language, errors.ErrRequestFailed, err)}
	}

	// Handle HTTP errors - retry on server errors (5xx), 404, or 429
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &errors.APIError{
			Code:        resp.StatusCode,
			Message:     fmt.Sprintf(c.I18n(i18n.MsgHTTPError), resp.StatusCode),
			Endpoint:    endpoint,
			RawResponse: string(respBuf.Bytes()),
			Lang:        c.Language,
		}
		// Retry on 5xx server errors, 404s, and 429 (rate limit).
		// Note: 404 is included because the GSPAY API may transiently return 404
		// during service deployments or load balancer routing changes.
		// 429 indicates rate limiting - retry with backoff.
		retry := (resp.StatusCode >= 500 || resp.StatusCode == 404 || resp.StatusCode == 429)

		// Log error
		c.logger.Error(c.I18n(i18n.LogHTTPErrorResponse),
			"endpoint", c.LogEndpoint(endpoint),
			"statusCode", resp.StatusCode,
			"retryable", retry,
		)

		respBuf.Reset()
		gc.Default.Put(respBuf)

		// Return specific error for rate limiting with Retry-After support
		if resp.StatusCode == 429 {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			return responseResult{
				Retry:      retry,
				RetryAfter: retryAfter,
				Err:        errors.New(c.Language, errors.ErrRateLimited),
			}
		}

		return responseResult{Retry: retry, Err: apiErr}
	}

	// Handle empty response
	if respBuf.Len() == 0 {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return responseResult{Retry: true, Err: errors.New(c.Language, errors.ErrEmptyResponse)}
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(respBuf.Bytes(), &apiResp); err != nil {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return responseResult{Err: errors.New(c.Language, errors.ErrInvalidJSON, err)}
	}

	// Debug logging
	c.logger.Debug(c.I18n(i18n.LogAPIResponseReceived),
		"endpoint", c.LogEndpoint(endpoint),
		"status", resp.StatusCode,
		"body", string(respBuf.Bytes()),
	)

	// Check for API-level errors
	if !apiResp.IsSuccess() {
		apiErr := &errors.APIError{
			Code:        apiResp.Code,
			Message:     apiResp.Message,
			Endpoint:    endpoint,
			RawResponse: string(respBuf.Bytes()),
			Lang:        c.Language,
		}
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return responseResult{Err: apiErr}
	}

	// Clean up buffer
	respBuf.Reset()
	gc.Default.Put(respBuf)

	return responseResult{Response: &apiResp}
}

// requestParams holds the parameters for a single request attempt.
type requestParams struct {
	Method   string
	FullURL  string
	Endpoint string
	Body     io.Reader
	HasBody  bool
	Attempt  int
}

// retryParams holds the parameters for request execution with retry logic.
type retryParams struct {
	requestParams
	// BodyBuffer is the original body buffer for resetting on retry.
	BodyBuffer gc.Buffer
}

// performRequest executes a single HTTP request attempt.
func (c *Client) performRequest(ctx context.Context, params requestParams) responseResult {
	req, err := c.createHTTPRequest(ctx, params.Method, params.FullURL, params.Body, params.HasBody)
	if err != nil {
		return responseResult{Err: err}
	}

	// Log outgoing request
	c.logger.Debug(c.I18n(i18n.LogSendingRequest),
		"method", params.Method,
		"endpoint", c.LogEndpoint(params.Endpoint),
		"attempt", params.Attempt,
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Log error
		c.logger.Error(c.I18n(i18n.LogRequestFailed),
			"endpoint", c.LogEndpoint(params.Endpoint),
			"attempt", params.Attempt,
			"error", err.Error(),
		)
		// Retry on transient network errors
		return responseResult{Retry: true, Err: errors.New(c.Language, errors.ErrRequestFailed, err)}
	}

	result := c.processResponse(resp, params.Endpoint)
	if result.Err != nil {
		return result
	}

	// Log success
	c.logger.Info(c.I18n(i18n.LogRequestCompleted),
		"endpoint", c.LogEndpoint(params.Endpoint),
		"attempts", params.Attempt+1,
	)

	return result
}

// executeWithRetry executes the HTTP request with retry logic.
func (c *Client) executeWithRetry(ctx context.Context, params retryParams) (*Response, error) {
	var lastErr error
	var actualAttempts int
	var suggestedWait time.Duration // Server-suggested wait time from Retry-After header

	for attempt := 0; attempt <= c.Retries; attempt++ {
		actualAttempts = attempt
		if attempt > 0 {
			// Log retry attempt
			c.logger.Warn(c.I18n(i18n.LogRetryingRequest),
				"endpoint", c.LogEndpoint(params.Endpoint),
				"attempt", attempt,
				"maxRetries", c.Retries,
			)

			// Wait with server-suggested time or fallback to exponential backoff with jitter
			if err := c.waitBackoff(ctx, attempt, suggestedWait); err != nil {
				return nil, err
			}

			// Reset suggested wait for next iteration
			suggestedWait = 0

			// Reset body reader for retry
			if params.HasBody {
				params.Body = bytes.NewReader(params.BodyBuffer.Bytes())
			}
		}

		// Update attempt number and call performRequest
		params.Attempt = attempt
		result := c.performRequest(ctx, params.requestParams)
		if result.Err == nil {
			return result.Response, nil
		}

		lastErr = result.Err
		suggestedWait = result.RetryAfter

		if result.Retry && attempt < c.Retries {
			// Log retryable error with rate limit info if applicable
			if suggestedWait > 0 {
				c.logger.Warn(c.I18n(i18n.LogRateLimitedRetry),
					"endpoint", c.LogEndpoint(params.Endpoint),
					"attempt", attempt,
					"retryAfter", suggestedWait.String(),
				)
			} else {
				c.logger.Warn(c.I18n(i18n.LogRetryableError),
					"endpoint", c.LogEndpoint(params.Endpoint),
					"attempt", attempt,
					"error", result.Err.Error(),
				)
			}
			continue
		}
		break
	}

	// lastErr is always non-nil here because:
	// 1. The loop only exits via break when err != nil
	// 2. Successful requests return early
	return nil, fmt.Errorf(c.I18n(i18n.MsgRequestFailedAfterRetries)+": %w", actualAttempts, lastErr)
}

// waitBackoff waits before retrying a request.
// If suggestedWait is provided (> 0), it uses the server-suggested Retry-After duration.
// Otherwise, it falls back to exponential backoff with jitter to prevent thundering herd.
func (c *Client) waitBackoff(ctx context.Context, attempt int, suggestedWait time.Duration) error {
	var waitTime time.Duration

	if suggestedWait > 0 {
		// Use server-suggested wait time (from Retry-After header)
		// Cap at RetryWaitMax to prevent excessively long waits
		waitTime = min(suggestedWait, c.RetryWaitMax)
	} else {
		// Fallback to manual exponential backoff with jitter
		baseWait := min(c.RetryWaitMin*time.Duration(1<<(attempt-1)), c.RetryWaitMax)
		// Add up to 25% jitter
		var jitter time.Duration
		if jitterMax := int64(baseWait / 4); jitterMax > 0 {
			jitter = time.Duration(rand.Int64N(jitterMax))
		}
		waitTime = baseWait + jitter
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(waitTime):
		return nil
	}
}

// DoRequest performs an HTTP request with retry logic.
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body any) (*Response, error) {
	fullURL := c.BaseURL + endpoint
	hasBody := body != nil

	reqBody, reqBuf, cleanup, err := c.prepareRequestBody(body)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	return c.executeWithRetry(ctx, retryParams{
		requestParams: requestParams{
			Method:   method,
			FullURL:  fullURL,
			Endpoint: endpoint,
			Body:     reqBody,
			HasBody:  hasBody,
		},
		BodyBuffer: reqBuf,
	})
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, endpoint string, body any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPost, endpoint, body)
}

// Get performs a GET request with query parameters.
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) (*Response, error) {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		endpoint = endpoint + "?" + values.Encode()
	}
	return c.DoRequest(ctx, http.MethodGet, endpoint, nil)
}

// ParseData parses the data field from an API response.
// GSPAY2 API returns data as a JSON string that needs to be decoded.
func ParseData[T any](data json.RawMessage, lang i18n.Language) (*T, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// First, try to unmarshal as a string (JSON encoded string)
	var jsonStr string
	if err := json.Unmarshal(data, &jsonStr); err == nil {
		// It was a JSON string, now unmarshal the string content
		data = json.RawMessage(jsonStr)
	}

	// Try to unmarshal as array first
	var arr []T
	if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
		return &arr[0], nil
	}

	// Try to unmarshal as array of strings
	var strArr []string
	if err := json.Unmarshal(data, &strArr); err == nil && len(strArr) > 0 {
		// Try to unmarshal the first string as T if it's JSON
		var result T
		if err := json.Unmarshal([]byte(strArr[0]), &result); err == nil {
			return &result, nil
		}
	}

	// Try to unmarshal as single object
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.New(lang, errors.ErrInvalidJSON, err)
	}

	return &result, nil
}
