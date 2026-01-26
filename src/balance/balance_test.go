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

package balance

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Get(t *testing.T) {
	t.Run("gets balance successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/get/balance")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "success",
				"data":    []map[string]float64{{"balance": 100000.00, "usdt_balance": 0.0}},
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewService(c)

		resp, err := svc.Get(context.Background())

		require.NoError(t, err)
		assert.Equal(t, "100000.00", resp)
	})

	t.Run("handles API error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"code":    500,
				"message": "internal server error",
			})
		}))
		defer server.Close()

		c := client.New("auth-key", "secret-key", client.WithBaseURL(server.URL))
		svc := NewService(c)

		_, err := svc.Get(context.Background())

		require.Error(t, err)
	})
}
