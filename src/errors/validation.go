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

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
)

// ValidationError represents a validation error for request parameters.
type ValidationError struct {
	Field   string
	Message string
	Lang    i18n.Language
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	format := i18n.Get(e.Lang, i18n.MsgValidationErrorFormat)
	return fmt.Sprintf(format, e.Field, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(lang i18n.Language, field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message, Lang: lang}
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
