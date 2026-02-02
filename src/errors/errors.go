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
	"strings"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
)

// Sentinel errors for common error conditions.
// These use i18n for their default English messages.
var (
	// ErrInvalidTransactionID is returned when the transaction ID is invalid.
	ErrInvalidTransactionID = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidTransactionID))
	// ErrInvalidAmount is returned when the payment amount is invalid.
	ErrInvalidAmount = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidAmount))
	// ErrInvalidBankCode is returned when the bank code is not recognized.
	ErrInvalidBankCode = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidBankCode))
	// ErrInvalidSignature is returned when signature verification fails.
	ErrInvalidSignature = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidSignature))
	// ErrMissingCallbackField is returned when a required callback field is missing.
	ErrMissingCallbackField = errors.New(i18n.Get(i18n.English, i18n.MsgMissingCallbackField))
	// ErrEmptyResponse is returned when the API returns an empty response.
	ErrEmptyResponse = errors.New(i18n.Get(i18n.English, i18n.MsgEmptyResponse))
	// ErrInvalidJSON is returned when the API response is not valid JSON.
	ErrInvalidJSON = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidJSON))
	// ErrRequestFailed is returned when the HTTP request fails.
	ErrRequestFailed = errors.New(i18n.Get(i18n.English, i18n.MsgRequestFailed))
	// ErrIPNotWhitelisted is returned when the callback IP is not in the whitelist.
	ErrIPNotWhitelisted = errors.New(i18n.Get(i18n.English, i18n.MsgIPNotWhitelisted))
	// ErrInvalidIPAddress is returned when the IP address format is invalid.
	ErrInvalidIPAddress = errors.New(i18n.Get(i18n.English, i18n.MsgInvalidIPAddress))
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

// LocalizedError represents an error with language-specific messages.
type LocalizedError struct {
	key  i18n.MessageKey
	lang i18n.Language
}

// Error implements the error interface.
func (e *LocalizedError) Error() string {
	return i18n.Get(e.lang, e.key)
}

// Key returns the message key of the error.
func (e *LocalizedError) Key() i18n.MessageKey {
	return e.key
}

// NewLocalizedError creates a new localized error with the specified language and message key.
func NewLocalizedError(lang i18n.Language, key i18n.MessageKey) *LocalizedError {
	return &LocalizedError{key: key, lang: lang}
}

// NewLocalizedWrappedError wraps a sentinel error with a localized message.
// This maintains errors.Is() compatibility while providing localized messages.
// Example: NewLocalizedWrappedError(ErrMissingCallbackField, lang, KeyMissingCallbackField, "id")
func NewLocalizedWrappedError(sentinel error, lang i18n.Language, key i18n.MessageKey, fieldName string) error {
	localizedMsg := i18n.Get(lang, key)
	if fieldName != "" {
		return fmt.Errorf("%s: %s: %w", localizedMsg, fieldName, sentinel)
	}
	return fmt.Errorf("%s: %w", localizedMsg, sentinel)
}

// NewMissingFieldError creates a localized error for a missing callback field.
// Wraps ErrMissingCallbackField while maintaining errors.Is() compatibility.
func NewMissingFieldError(lang i18n.Language, fieldName string) error {
	return NewLocalizedWrappedError(ErrMissingCallbackField, lang, MsgMissingCallbackField, fieldName)
}

// NewInvalidSignatureError creates a localized error for invalid signature.
// Wraps ErrInvalidSignature while maintaining errors.Is() compatibility.
func NewInvalidSignatureError(lang i18n.Language) error {
	return NewLocalizedWrappedError(ErrInvalidSignature, lang, MsgInvalidSignature, "")
}

// IsLocalizedError checks if an error is a LocalizedError.
func IsLocalizedError(err error) bool {
	var locErr *LocalizedError
	return errors.As(err, &locErr)
}

// GetLocalizedError extracts a LocalizedError from an error.
// Returns nil if the error is not a LocalizedError.
func GetLocalizedError(err error) *LocalizedError {
	var locErr *LocalizedError
	if errors.As(err, &locErr) {
		return locErr
	}
	return nil
}

// GetMessage is a convenience function that delegates to i18n.Get.
// It retrieves a message for the specified language and key.
func GetMessage(lang i18n.Language, key i18n.MessageKey) string {
	return i18n.Get(lang, key)
}

// Re-export i18n types and constants for convenience
type (
	// Language represents a supported language.
	Language = i18n.Language
	// MessageKey identifies a translatable message.
	MessageKey = i18n.MessageKey
)

// Re-export language constants
const (
	English    = i18n.English
	Indonesian = i18n.Indonesian
)

// Re-export message keys
const (
	// Sentinel error message keys
	MsgInvalidTransactionID = i18n.MsgInvalidTransactionID
	MsgInvalidAmount        = i18n.MsgInvalidAmount
	MsgInvalidBankCode      = i18n.MsgInvalidBankCode
	MsgInvalidSignature     = i18n.MsgInvalidSignature
	MsgMissingCallbackField = i18n.MsgMissingCallbackField
	MsgEmptyResponse        = i18n.MsgEmptyResponse
	MsgInvalidJSON          = i18n.MsgInvalidJSON
	MsgRequestFailed        = i18n.MsgRequestFailed
	MsgIPNotWhitelisted     = i18n.MsgIPNotWhitelisted
	MsgInvalidIPAddress     = i18n.MsgInvalidIPAddress

	// Validation error message keys
	KeyMinAmountIDR        = i18n.MsgMinAmountIDR
	KeyMinAmountUSDT       = i18n.MsgMinAmountUSDT
	KeyMinPayoutAmountIDR  = i18n.MsgMinPayoutAmountIDR
	KeyInvalidAmountFormat = i18n.MsgInvalidAmountFormat

	// Request retry message keys
	MsgRequestFailedAfterRetries = i18n.MsgRequestFailedAfterRetries
)
