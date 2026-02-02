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

// Package amount provides utility functions for formatting monetary amounts.
//
// This package is used for callback signature verification where amounts are
// formatted with 2 decimal places (e.g., "10000.00").
package amount

import (
	"fmt"
	"strconv"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"
)

// Format formats an amount string to exactly 2 decimal places.
//
// This is used for callback signature verification where amounts are
// formatted with 2 decimal places (e.g., "10000.00").
//
// Note: Uses float64 parsing which may have precision limitations for
// extremely large amounts (> 2^53). For typical payment amounts, this
// is not a concern.
func Format(amountStr string, lang i18n.Language) (string, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return "", errors.NewValidationError("amount",
			errors.GetMessage(errors.Language(lang), errors.KeyInvalidAmountFormat))
	}
	return fmt.Sprintf("%.2f", amount), nil
}
