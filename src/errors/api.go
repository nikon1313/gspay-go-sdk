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

package errors

import (
	"errors"
	"fmt"
	"strings"
)

// APIError represents an error returned by the GSPAY2 API.
type APIError struct {
	// Code is the HTTP status code or API error code.
	Code int `json:"code"`
	// Message is the error message from the API.
	Message string `json:"message"`
	// Endpoint is the API endpoint that was called.
	Endpoint string `json:"-"`
	// RawResponse contains the raw response body for debugging.
	RawResponse string `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Endpoint != "" {
		sanitizedEndpoint := sanitizeEndpoint(e.Endpoint)
		return fmt.Sprintf("gspay: API error %d on %s: %s", e.Code, sanitizedEndpoint, e.Message)
	}
	return fmt.Sprintf("gspay: API error %d: %s", e.Code, e.Message)
}

// sanitizeEndpoint redacts sensitive information like auth keys from endpoint URLs.
func sanitizeEndpoint(endpoint string) string {
	// Redact auth key in operator endpoints:
	// - /v2/integrations/operator/{authkey}/...  (singular - e.g., balance)
	// - /v2/integrations/operators/{authkey}/... (plural - e.g., USDT)
	//
	// Path structure after split:
	// parts[0] = "" (empty, from leading slash)
	// parts[1] = "v2"
	// parts[2] = "integrations"
	// parts[3] = "operator" or "operators"
	// parts[4] = authkey (to be redacted)
	// parts[5+] = remaining path segments
	parts := strings.Split(endpoint, "/")
	if len(parts) >= 5 && parts[1] == "v2" && parts[2] == "integrations" && len(parts[4]) > 0 {
		if parts[3] == "operator" || parts[3] == "operators" {
			parts[4] = "[REDACTED]"
			return strings.Join(parts, "/")
		}
	}
	return endpoint
}

// IsAPIError checks if an error is an APIError.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// GetAPIError extracts an APIError from an error.
// Returns nil if the error is not an APIError.
func GetAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}
