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

package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBankName(t *testing.T) {
	t.Run("returns IDR bank name", func(t *testing.T) {
		assert.Equal(t, "Bank BCA", GetBankName("BCA", CurrencyIDR))
		assert.Equal(t, "Bank BRI", GetBankName("BRI", CurrencyIDR))
		assert.Equal(t, "DANA", GetBankName("DANA", CurrencyIDR))
	})

	t.Run("returns MYR bank name", func(t *testing.T) {
		assert.Equal(t, "MAYBANK", GetBankName("MBB", CurrencyMYR))
		assert.Equal(t, "CIMB", GetBankName("CIMB", CurrencyMYR))
	})

	t.Run("returns THB bank name", func(t *testing.T) {
		assert.Equal(t, "KASIKORNBANK PUBLIC COMPANY LIMITED", GetBankName("KBANK", CurrencyTHB))
	})

	t.Run("returns empty string for unknown bank", func(t *testing.T) {
		assert.Empty(t, GetBankName("UNKNOWN", CurrencyIDR))
	})

	t.Run("returns empty string for unknown currency", func(t *testing.T) {
		assert.Empty(t, GetBankName("BCA", Currency("XXX")))
	})
}

func TestGetBankCodes(t *testing.T) {
	t.Run("returns IDR bank codes", func(t *testing.T) {
		codes := GetBankCodes(CurrencyIDR)
		assert.NotEmpty(t, codes)
		assert.Contains(t, codes, "BCA")
		assert.Contains(t, codes, "BRI")
		assert.Contains(t, codes, "DANA")
	})

	t.Run("returns MYR bank codes", func(t *testing.T) {
		codes := GetBankCodes(CurrencyMYR)
		assert.NotEmpty(t, codes)
		assert.Contains(t, codes, "MBB")
		assert.Contains(t, codes, "CIMB")
	})

	t.Run("returns THB bank codes", func(t *testing.T) {
		codes := GetBankCodes(CurrencyTHB)
		assert.NotEmpty(t, codes)
		assert.Contains(t, codes, "KBANK")
		assert.Contains(t, codes, "SCB")
	})

	t.Run("returns nil for unknown currency", func(t *testing.T) {
		codes := GetBankCodes(Currency("XXX"))
		assert.Nil(t, codes)
	})
}

func TestIsValidBankIDR(t *testing.T) {
	t.Run("returns true for valid banks", func(t *testing.T) {
		assert.True(t, IsValidBankIDR("BCA"))
		assert.True(t, IsValidBankIDR("BRI"))
		assert.True(t, IsValidBankIDR("MANDIRI"))
		assert.True(t, IsValidBankIDR("DANA"))
		assert.True(t, IsValidBankIDR("OVO"))
	})

	t.Run("returns false for invalid banks", func(t *testing.T) {
		assert.False(t, IsValidBankIDR("UNKNOWN"))
		assert.False(t, IsValidBankIDR(""))
		assert.False(t, IsValidBankIDR("bca")) // case-sensitive
	})
}
