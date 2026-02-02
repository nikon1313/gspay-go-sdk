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
	"fmt"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
)

// New creates a localized error wrapping the provided sentinel error.
// It automatically resolves the correct message key for the given sentinel.
//
// Optional arguments:
//   - If an error is provided, it is wrapped as the cause.
//   - If a string is provided, it is added as context (e.g., field name).
func New(lang i18n.Language, sentinel error, args ...any) error {
	key, ok := sentinelMessages[sentinel]
	if !ok {
		return sentinel
	}

	msg := i18n.Get(lang, key)

	if len(args) > 0 {
		// Handle cause (error)
		if cause, ok := args[0].(error); ok && cause != nil {
			// Wrapping sentinel error inside a new error with localized message
			baseErr := fmt.Errorf("%s: %w", msg, sentinel)
			// Returning the base error wrapped with the cause
			return fmt.Errorf("%w: %v", baseErr, cause)
		}
		// Handle context (string)
		if contextStr, ok := args[0].(string); ok && contextStr != "" {
			return fmt.Errorf("%s: %s: %w", msg, contextStr, sentinel)
		}
	}

	return fmt.Errorf("%s: %w", msg, sentinel)
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
