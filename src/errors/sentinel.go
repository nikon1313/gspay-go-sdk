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

// sentinelMessages maps sentinel errors to their message keys.
var sentinelMessages = map[error]i18n.MessageKey{
	ErrInvalidTransactionID: MsgInvalidTransactionID,
	ErrInvalidAmount:        MsgInvalidAmount,
	ErrInvalidBankCode:      MsgInvalidBankCode,
	ErrInvalidSignature:     MsgInvalidSignature,
	ErrMissingCallbackField: MsgMissingCallbackField,
	ErrEmptyResponse:        MsgEmptyResponse,
	ErrInvalidJSON:          MsgInvalidJSON,
	ErrRequestFailed:        MsgRequestFailed,
	ErrIPNotWhitelisted:     MsgIPNotWhitelisted,
	ErrInvalidIPAddress:     MsgInvalidIPAddress,
}
