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

func TestPaymentStatus_String(t *testing.T) {
	tests := []struct {
		status   PaymentStatus
		expected string
	}{
		{StatusPending, "Pending/Expired"},
		{StatusSuccess, "Success"},
		{StatusFailed, "Timeout/Failed"},
		{StatusTimeout, "Timeout/Failed"},
		{PaymentStatus(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestPaymentStatus_IsSuccess(t *testing.T) {
	assert.True(t, StatusSuccess.IsSuccess())
	assert.False(t, StatusPending.IsSuccess())
	assert.False(t, StatusFailed.IsSuccess())
	assert.False(t, StatusTimeout.IsSuccess())
}

func TestPaymentStatus_IsFailed(t *testing.T) {
	assert.True(t, StatusFailed.IsFailed())
	assert.True(t, StatusTimeout.IsFailed())
	assert.False(t, StatusPending.IsFailed())
	assert.False(t, StatusSuccess.IsFailed())
}

func TestPaymentStatus_IsPending(t *testing.T) {
	assert.True(t, StatusPending.IsPending())
	assert.False(t, StatusSuccess.IsPending())
	assert.False(t, StatusFailed.IsPending())
	assert.False(t, StatusTimeout.IsPending())
}

func TestParsePaymentStatus(t *testing.T) {
	tests := []struct {
		input    int
		expected PaymentStatus
	}{
		{0, StatusPending},
		{1, StatusSuccess},
		{2, StatusFailed},
		{4, StatusTimeout},
		{99, PaymentStatus(99)},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, ParsePaymentStatus(tt.input))
		})
	}
}
