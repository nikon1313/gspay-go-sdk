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
	"net/http"
	"net/url"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/helper/gc"
)

// Response represents a generic API response structure.
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// IsSuccess checks if the API response indicates success.
func (r *Response) IsSuccess() bool {
	return r.Code == 200
}

// DoRequest performs an HTTP request with retry logic.
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body interface{}) (*Response, error) {
	fullURL := c.BaseURL + endpoint

	var reqBody io.Reader
	var reqBuf gc.Buffer
	if body != nil {
		reqBuf = gc.Default.Get()
		if err := json.NewEncoder(reqBuf).Encode(body); err != nil {
			reqBuf.Reset()
			gc.Default.Put(reqBuf)
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(reqBuf.Bytes())
	}
	defer func() {
		if reqBuf != nil {
			reqBuf.Reset()
			gc.Default.Put(reqBuf)
		}
	}()

	var lastErr error
	for attempt := 0; attempt <= c.Retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			waitTime := c.RetryWaitMin * time.Duration(attempt)
			if waitTime > c.RetryWaitMax {
				waitTime = c.RetryWaitMax
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}

			// Reset body reader for retry
			if body != nil {
				reqBody = bytes.NewReader(reqBuf.Bytes())
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", constants.UserAgent())
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			// Retry on transient network errors
			if attempt < c.Retries {
				continue
			}
			return nil, fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
		}

		respBuf := gc.Default.Get()
		_, err = respBuf.ReadFrom(resp.Body)
		resp.Body.Close()

		if err != nil {
			respBuf.Reset()
			gc.Default.Put(respBuf)
			lastErr = err
			if attempt < c.Retries {
				continue
			}
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Handle HTTP errors - retry on server errors (5xx) or 404
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = &errors.APIError{
				Code:        resp.StatusCode,
				Message:     fmt.Sprintf("HTTP Error: %d", resp.StatusCode),
				Endpoint:    endpoint,
				RawResponse: string(respBuf.Bytes()),
			}
			respBuf.Reset()
			gc.Default.Put(respBuf)
			if (resp.StatusCode >= 500 || resp.StatusCode == 404) && attempt < c.Retries {
				continue
			}
			return nil, lastErr
		}

		// Handle empty response
		if respBuf.Len() == 0 {
			respBuf.Reset()
			gc.Default.Put(respBuf)
			lastErr = errors.ErrEmptyResponse
			if attempt < c.Retries {
				continue
			}
			return nil, lastErr
		}

		// Parse response
		var apiResp Response
		if err := json.Unmarshal(respBuf.Bytes(), &apiResp); err != nil {
			respBuf.Reset()
			gc.Default.Put(respBuf)
			return nil, fmt.Errorf("%w: %v", errors.ErrInvalidJSON, err)
		}

		// Check for API-level errors
		if !apiResp.IsSuccess() {
			err := &errors.APIError{
				Code:        apiResp.Code,
				Message:     apiResp.Message,
				Endpoint:    endpoint,
				RawResponse: string(respBuf.Bytes()),
			}
			respBuf.Reset()
			gc.Default.Put(respBuf)
			return nil, err
		}

		// Success, return response and clean up buffer
		respBuf.Reset()
		gc.Default.Put(respBuf)
		return &apiResp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("request failed after %d retries", c.Retries)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}) (*Response, error) {
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
func ParseData[T any](data json.RawMessage) (*T, error) {
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

	// Try to unmarshal as single object
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse data field: %w", err)
	}

	return &result, nil
}
