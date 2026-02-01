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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoRequest(t *testing.T) {
	t.Run("successful POST request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Equal(t, constants.UserAgent(), r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    `{"payment_url":"https://pay.example.com"}`,
			})
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL))
		resp, err := c.Post(t.Context(), "/test", map[string]string{"key": "value"})

		require.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "success", resp.Message)
	})

	t.Run("successful GET request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "value", r.URL.Query().Get("key"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
			})
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL))
		resp, err := c.Get(t.Context(), "/test", map[string]string{"key": "value"})

		require.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("handles API error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    400,
				"message": "invalid request",
			})
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL))
		_, err := c.Post(t.Context(), "/test", nil)

		require.Error(t, err)
		apiErr := errors.GetAPIError(err)
		require.NotNil(t, apiErr)
		assert.Equal(t, 400, apiErr.Code)
		assert.Equal(t, "invalid request", apiErr.Message)
	})

	t.Run("handles HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL), WithRetries(0))
		_, err := c.Post(t.Context(), "/test", nil)

		require.Error(t, err)
		apiErr := errors.GetAPIError(err)
		require.NotNil(t, apiErr)
		assert.Equal(t, 500, apiErr.Code)
	})

	t.Run("handles empty response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// Empty body
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL), WithRetries(0))
		_, err := c.Post(t.Context(), "/test", nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrEmptyResponse)
	})

	t.Run("handles invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not json"))
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL), WithRetries(0))
		_, err := c.Post(t.Context(), "/test", nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrInvalidJSON)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		c := New("auth-key", "secret-key", WithBaseURL(server.URL))

		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel immediately

		_, err := c.Post(ctx, "/test", nil)
		require.Error(t, err)
	})

	t.Run("retries on server error", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
			})
		}))
		defer server.Close()

		c := New(
			"auth-key",
			"secret-key",
			WithBaseURL(server.URL),
			WithRetries(3),
			WithRetryWait(10*time.Millisecond, 50*time.Millisecond),
		)
		resp, err := c.Post(t.Context(), "/test", nil)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, 3, attempts)
	})

	t.Run("exponential backoff timing", func(t *testing.T) {
		attemptTimes := make([]time.Time, 0, 3)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptTimes = append(attemptTimes, time.Now())
			if len(attemptTimes) < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
			})
		}))
		defer server.Close()

		c := New(
			"auth-key",
			"secret-key",
			WithBaseURL(server.URL),
			WithRetries(2),
			WithRetryWait(10*time.Millisecond, 100*time.Millisecond),
		)
		resp, err := c.Post(t.Context(), "/test", nil)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
		require.Len(t, attemptTimes, 3)
		diff1 := attemptTimes[1].Sub(attemptTimes[0])
		diff2 := attemptTimes[2].Sub(attemptTimes[1])
		assert.True(t, diff1 >= 10*time.Millisecond, "first retry delay should be at least 10ms")
		assert.True(t, diff2 >= 20*time.Millisecond, "second retry delay should be at least 20ms")
	})

	t.Run("fails after retries exhausted", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := New(
			"auth-key",
			"secret-key",
			WithBaseURL(server.URL),
			WithRetries(2),
			WithRetryWait(1*time.Millisecond, 10*time.Millisecond),
		)
		_, err := c.Post(t.Context(), "/test", nil)

		require.Error(t, err)
		assert.Equal(t, 3, attempts) // initial + 2 retries
		assert.Contains(t, err.Error(), "request failed after 2 retries")
		apiErr := errors.GetAPIError(err)
		require.NotNil(t, apiErr)
		assert.Equal(t, 500, apiErr.Code)
	})
}

func TestParseData(t *testing.T) {
	t.Run("parses JSON string data", func(t *testing.T) {
		data := json.RawMessage(`"{\"payment_url\":\"https://pay.example.com\"}"`)

		type testStruct struct {
			PaymentURL string `json:"payment_url"`
		}

		result, err := ParseData[testStruct](data)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "https://pay.example.com", result.PaymentURL)
	})

	t.Run("parses array data", func(t *testing.T) {
		data := json.RawMessage(`[{"payment_url":"https://pay.example.com"}]`)

		type testStruct struct {
			PaymentURL string `json:"payment_url"`
		}

		result, err := ParseData[testStruct](data)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "https://pay.example.com", result.PaymentURL)
	})

	t.Run("parses object data", func(t *testing.T) {
		data := json.RawMessage(`{"payment_url":"https://pay.example.com"}`)

		type testStruct struct {
			PaymentURL string `json:"payment_url"`
		}

		result, err := ParseData[testStruct](data)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "https://pay.example.com", result.PaymentURL)
	})

	t.Run("handles empty data", func(t *testing.T) {
		type testStruct struct{}
		result, err := ParseData[testStruct](nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestResponse_IsSuccess(t *testing.T) {
	t.Run("returns true for code 200", func(t *testing.T) {
		resp := &Response{Code: 200}
		assert.True(t, resp.IsSuccess())
	})

	t.Run("returns false for non-200 codes", func(t *testing.T) {
		testCases := []int{0, 100, 201, 400, 500}
		for _, code := range testCases {
			resp := &Response{Code: code}
			assert.False(t, resp.IsSuccess(), "expected false for code %d", code)
		}
	})
}
