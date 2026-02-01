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

package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguage_IsValid(t *testing.T) {
	t.Run("returns true for English", func(t *testing.T) {
		assert.True(t, English.IsValid())
	})

	t.Run("returns true for Indonesian", func(t *testing.T) {
		assert.True(t, Indonesian.IsValid())
	})

	t.Run("returns false for unsupported language", func(t *testing.T) {
		assert.False(t, Language("fr").IsValid())
	})

	t.Run("returns false for empty language", func(t *testing.T) {
		assert.False(t, Language("").IsValid())
	})
}

func TestLanguage_String(t *testing.T) {
	assert.Equal(t, "en", English.String())
	assert.Equal(t, "id", Indonesian.String())
}

func TestGet(t *testing.T) {
	t.Run("returns English message for English language", func(t *testing.T) {
		msg := Get(English, MsgInvalidTransactionID)
		assert.Equal(t, "transaction ID must be 5-20 characters", msg)
	})

	t.Run("returns Indonesian message for Indonesian language", func(t *testing.T) {
		msg := Get(Indonesian, MsgInvalidTransactionID)
		assert.Equal(t, "ID transaksi harus 5-20 karakter", msg)
	})

	t.Run("falls back to English for unknown language", func(t *testing.T) {
		msg := Get(Language("fr"), MsgInvalidAmount)
		assert.Equal(t, "invalid payment amount", msg)
	})

	t.Run("returns key for unknown message key", func(t *testing.T) {
		msg := Get(English, MessageKey("unknown_key"))
		assert.Equal(t, "unknown_key", msg)
	})

	t.Run("returns all validation messages in English", func(t *testing.T) {
		assert.Equal(t, "minimum amount is 10000 IDR", Get(English, MsgMinAmountIDR))
		assert.Equal(t, "minimum amount is 1.00 USDT", Get(English, MsgMinAmountUSDT))
		assert.Equal(t, "minimum payout amount is 10000 IDR", Get(English, MsgMinPayoutAmountIDR))
		assert.Equal(t, "invalid amount format", Get(English, MsgInvalidAmountFormat))
	})

	t.Run("returns all validation messages in Indonesian", func(t *testing.T) {
		assert.Equal(t, "jumlah minimum adalah 10000 IDR", Get(Indonesian, MsgMinAmountIDR))
		assert.Equal(t, "jumlah minimum adalah 1.00 USDT", Get(Indonesian, MsgMinAmountUSDT))
		assert.Equal(t, "jumlah pembayaran minimum adalah 10000 IDR", Get(Indonesian, MsgMinPayoutAmountIDR))
		assert.Equal(t, "format jumlah tidak valid", Get(Indonesian, MsgInvalidAmountFormat))
	})
}
