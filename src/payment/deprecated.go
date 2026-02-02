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

package payment

import (
	"fmt"
	"strconv"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
)

// verifyCallbackSignature performs the actual signature verification.
//
// Deprecated: Use VerifySignature directly instead.
func (s *USDTService) verifyCallbackSignature(callback *USDTCallback) error {
	lang := errors.Language(s.client.Language)

	// Check required fields
	if callback.CryptoPaymentID == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "cryptopayment_id")
	}
	if callback.Amount == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "amount")
	}
	if callback.TransactionID == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "transaction_id")
	}
	if callback.Signature == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "signature")
	}

	// Format amount with 2 decimal places
	amount, err := strconv.ParseFloat(callback.Amount, 64)
	if err != nil {
		return errors.NewValidationError(lang, "amount", errors.GetMessage(lang, errors.KeyInvalidAmountFormat))
	}
	formattedAmount := fmt.Sprintf("%.2f", amount)

	// Generate expected signature
	signatureData := fmt.Sprintf("%s%s%s%d%s",
		callback.CryptoPaymentID,
		formattedAmount,
		callback.TransactionID,
		callback.Status,
		s.client.SecretKey,
	)
	expectedSignature := s.client.GenerateSignature(signatureData)

	// Constant-time comparison to prevent timing attacks
	if !s.client.VerifySignature(expectedSignature, callback.Signature) {
		return errors.New(lang, errors.ErrInvalidSignature)
	}

	return nil
}

// verifyCallbackSignature performs the actual signature verification.
//
// Deprecated: Use VerifySignature directly instead.
func (s *IDRService) verifyCallbackSignature(callback *IDRCallback) error {
	return s.VerifySignature(
		string(callback.IDRPaymentID),
		string(callback.Amount),
		callback.TransactionID,
		callback.Status,
		callback.Signature,
	)
}
