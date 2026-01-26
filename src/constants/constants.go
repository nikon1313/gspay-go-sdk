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

// Package constants provides constant values for the GSPAY2 API.
package constants

// DefaultBaseURL is the default GSPAY2 API base URL.
const DefaultBaseURL = "https://api.thegspay.com"

// Default client configuration values.
const (
	DefaultTimeout      = 30 // seconds
	DefaultRetries      = 3
	DefaultRetryWaitMin = 500  // milliseconds
	DefaultRetryWaitMax = 2000 // milliseconds
)

// Minimum amount constraints.
const (
	MinAmountIDR  = 10000 // Minimum IDR amount
	MinAmountUSDT = 1.00  // Minimum USDT amount
)

// Transaction ID constraints.
const (
	MinTransactionIDLength = 5
	MaxTransactionIDLength = 20
)
