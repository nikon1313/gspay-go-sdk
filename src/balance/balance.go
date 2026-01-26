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

// Package balance provides balance query functionality for the GSPAY2 SDK.
package balance

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
)

// Response represents the response from querying operator balance.
type Response struct {
	// Balance is the operator's IDR balance.
	Balance float64 `json:"balance"`
	// UsdtBalance is the operator's USDT balance.
	UsdtBalance float64 `json:"usdt_balance"`
}

// Service handles balance operations.
type Service struct{ client *client.Client }

// NewService creates a new balance service.
func NewService(c *client.Client) *Service { return &Service{client: c} }

// Get queries the operator's available settlement balance.
func (s *Service) Get(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("/v2/integrations/operator/%s/get/balance", s.client.AuthKey)
	resp, err := s.client.Get(ctx, endpoint, nil)
	if err != nil {
		return "", err
	}

	result, err := client.ParseData[Response](resp.Data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f", (*result).Balance), nil
}
