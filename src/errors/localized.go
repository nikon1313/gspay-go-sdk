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
