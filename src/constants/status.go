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

// PaymentStatus represents the status of a payment or payout.
type PaymentStatus int

const (
	// StatusPending indicates the payment is pending or expired.
	StatusPending PaymentStatus = 0
	// StatusSuccess indicates the payment was successful.
	StatusSuccess PaymentStatus = 1
	// StatusFailed indicates the payment failed.
	StatusFailed PaymentStatus = 2
	// StatusTimeout indicates the payment timed out.
	StatusTimeout PaymentStatus = 4
)

// String returns the human-readable label for a payment status.
func (s PaymentStatus) String() string {
	switch s {
	case StatusPending:
		return "Pending/Expired"
	case StatusSuccess:
		return "Success"
	case StatusFailed, StatusTimeout:
		return "Timeout/Failed"
	default:
		return "Unknown"
	}
}

// IsSuccess returns true if the status indicates a successful payment.
func (s PaymentStatus) IsSuccess() bool {
	return s == StatusSuccess
}

// IsFailed returns true if the status indicates a failed or timed out payment.
func (s PaymentStatus) IsFailed() bool {
	return s == StatusFailed || s == StatusTimeout
}

// IsPending returns true if the status indicates a pending payment.
func (s PaymentStatus) IsPending() bool {
	return s == StatusPending
}

// ParsePaymentStatus converts an integer status code to PaymentStatus type.
func ParsePaymentStatus(status int) PaymentStatus {
	switch status {
	case 0:
		return StatusPending
	case 1:
		return StatusSuccess
	case 2:
		return StatusFailed
	case 4:
		return StatusTimeout
	default:
		return PaymentStatus(status)
	}
}
