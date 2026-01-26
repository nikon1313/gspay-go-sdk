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

// Package errors provides error types for the GSPAY2 SDK.
package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common error conditions.
var (
	// ErrInvalidTransactionID is returned when the transaction ID is invalid.
	ErrInvalidTransactionID = errors.New("transaction ID must be 5-20 characters")
	// ErrInvalidAmount is returned when the payment amount is invalid.
	ErrInvalidAmount = errors.New("invalid payment amount")
	// ErrInvalidBankCode is returned when the bank code is not recognized.
	ErrInvalidBankCode = errors.New("invalid bank code")
	// ErrInvalidSignature is returned when callback signature verification fails.
	ErrInvalidSignature = errors.New("invalid callback signature")
	// ErrMissingCallbackField is returned when a required callback field is missing.
	ErrMissingCallbackField = errors.New("missing required callback field")
	// ErrEmptyResponse is returned when the API returns an empty response.
	ErrEmptyResponse = errors.New("empty response from API")
	// ErrInvalidJSON is returned when the API response is not valid JSON.
	ErrInvalidJSON = errors.New("invalid JSON response")
	// ErrRequestFailed is returned when the HTTP request fails.
	ErrRequestFailed = errors.New("request failed")
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
		return fmt.Sprintf("gspay: API error %d on %s: %s", e.Code, e.Endpoint, e.Message)
	}
	return fmt.Sprintf("gspay: API error %d: %s", e.Code, e.Message)
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

// ValidationError represents a validation error for request parameters.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("gspay: validation error for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// GetValidationError extracts a ValidationError from an error.
// Returns nil if the error is not a ValidationError.
func GetValidationError(err error) *ValidationError {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		return valErr
	}
	return nil
}
