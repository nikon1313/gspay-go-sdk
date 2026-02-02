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

package amount

import (
	"testing"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Run("formats integer amount to 2 decimal places", func(t *testing.T) {
		result, err := Format("10000", i18n.English)
		require.NoError(t, err)
		assert.Equal(t, "10000.00", result)
	})

	t.Run("formats decimal amount to 2 decimal places", func(t *testing.T) {
		result, err := Format("10000.5", i18n.English)
		require.NoError(t, err)
		assert.Equal(t, "10000.50", result)
	})

	t.Run("truncates extra decimal places", func(t *testing.T) {
		result, err := Format("10000.999", i18n.English)
		require.NoError(t, err)
		assert.Equal(t, "10001.00", result)
	})

	t.Run("handles large amounts", func(t *testing.T) {
		result, err := Format("1000000000", i18n.English)
		require.NoError(t, err)
		assert.Equal(t, "1000000000.00", result)
	})

	t.Run("returns error for invalid amount string", func(t *testing.T) {
		_, err := Format("invalid", i18n.English)
		require.Error(t, err)
		assert.True(t, errors.IsValidationError(err))
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		_, err := Format("", i18n.English)
		require.Error(t, err)
		assert.True(t, errors.IsValidationError(err))
	})

	t.Run("returns localized error for Indonesian language", func(t *testing.T) {
		_, err := Format("invalid", i18n.Indonesian)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "format jumlah tidak valid")
	})
}
