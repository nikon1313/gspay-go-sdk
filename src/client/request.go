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
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/helper/gc"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
)

// Response represents a generic API response structure.
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// IsSuccess checks if the API response indicates success.
func (r *Response) IsSuccess() bool { return r.Code == 200 }

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
func (c *Client) processResponse(resp *http.Response, endpoint string) (*Response, bool, error) {
	respBuf := gc.Default.Get()
	_, err := respBuf.ReadFrom(resp.Body)
	resp.Body.Close()

	if err != nil {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return nil, true, errors.New(c.Language, errors.ErrRequestFailed, err)
	}

	// Handle HTTP errors - retry on server errors (5xx) or 404
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &errors.APIError{
			Code:        resp.StatusCode,
			Message:     fmt.Sprintf("HTTP Error: %d", resp.StatusCode),
			Endpoint:    endpoint,
			RawResponse: string(respBuf.Bytes()),
		}
		respBuf.Reset()
		gc.Default.Put(respBuf)
		// Retry on 5xx server errors and 404s.
		// Note: 404 is included because the GSPAY API may transiently return 404
		// during service deployments or load balancer routing changes.
		retry := (resp.StatusCode >= 500 || resp.StatusCode == 404)
		return nil, retry, apiErr
	}

	// Handle empty response
	if respBuf.Len() == 0 {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return nil, true, errors.New(c.Language, errors.ErrEmptyResponse)
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(respBuf.Bytes(), &apiResp); err != nil {
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return nil, false, errors.New(c.Language, errors.ErrInvalidJSON, err)
	}

	// Debug logging
	if c.Debug {
		fmt.Printf("DEBUG API Response for %s: %s\n", endpoint, string(respBuf.Bytes()))
	}

	// Check for API-level errors
	if !apiResp.IsSuccess() {
		apiErr := &errors.APIError{
			Code:        apiResp.Code,
			Message:     apiResp.Message,
			Endpoint:    endpoint,
			RawResponse: string(respBuf.Bytes()),
		}
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return nil, false, apiErr
	}

	// Clean up buffer
	respBuf.Reset()
	gc.Default.Put(respBuf)

	return &apiResp, false, nil
}

// executeWithRetry executes the HTTP request with retry logic.
func (c *Client) executeWithRetry(ctx context.Context, method, fullURL string, reqBody io.Reader, reqBuf gc.Buffer, hasBody bool, endpoint string) (*Response, error) {
	var lastErr error
	var actualAttempts int
	for attempt := 0; attempt <= c.Retries; attempt++ {
		actualAttempts = attempt
		if attempt > 0 {
			// Wait with exponential backoff and jitter
			if err := c.waitBackoff(ctx, attempt); err != nil {
				return nil, err
			}

			// Reset body reader for retry
			if hasBody {
				reqBody = bytes.NewReader(reqBuf.Bytes())
			}
		}

		req, err := c.createHTTPRequest(ctx, method, fullURL, reqBody, hasBody)
		if err != nil {
			return nil, err
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = errors.New(c.Language, errors.ErrRequestFailed, err)
			// Retry on transient network errors
			if attempt < c.Retries {
				continue
			}
			break
		}

		apiResp, retry, err := c.processResponse(resp, endpoint)
		if err != nil {
			lastErr = err
			if retry && attempt < c.Retries {
				continue
			}
			break
		}

		return apiResp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf(i18n.Get(c.Language, i18n.MsgRequestFailedAfterRetries)+": %w", actualAttempts, lastErr)
	}
	return nil, fmt.Errorf(i18n.Get(c.Language, i18n.MsgRequestFailedAfterRetries), actualAttempts)
}

// waitBackoff calculates the backoff duration and waits.
func (c *Client) waitBackoff(ctx context.Context, attempt int) error {
	// Exponential backoff with jitter to prevent thundering herd
	baseWait := min(c.RetryWaitMin*time.Duration(1<<(attempt-1)), c.RetryWaitMax)
	// Add up to 25% jitter
	var jitter time.Duration
	// TODO: Do we need to use crypto/rand instead of math/rand?
	if jitterMax := int64(baseWait / 4); jitterMax > 0 {
		jitter = time.Duration(rand.Int64N(jitterMax))
	}
	waitTime := baseWait + jitter

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

	return c.executeWithRetry(ctx, method, fullURL, reqBody, reqBuf, hasBody, endpoint)
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
