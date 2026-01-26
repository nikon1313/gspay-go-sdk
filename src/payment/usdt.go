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
	"context"
	"fmt"
	"strconv"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/internal/signature"
)

// USDTRequest represents a request to create a USDT payment.
type USDTRequest struct {
	// TransactionID is a unique transaction ID.
	TransactionID string `json:"transaction_id"`
	// Username is the customer ID or username.
	Username string `json:"player_username"`
	// Amount is the payment amount in USDT (2 decimal places).
	Amount float64 `json:"amount"`
}

// usdtAPIRequest is the internal API request structure.
type usdtAPIRequest struct {
	TransactionID string `json:"transaction_id"`
	Username      string `json:"player_username"`
	Amount        string `json:"amount"`
	Signature     string `json:"signature"`
}

// USDTResponse represents the response from creating a USDT payment.
type USDTResponse struct {
	// PaymentURL is the URL to redirect the user for payment.
	PaymentURL string `json:"payment_url"`
	// CryptoPaymentID is the unique payment ID assigned by GSPAY2.
	CryptoPaymentID string `json:"cryptopayment_id"`
	// ExpireDate is the payment expiration date/time.
	ExpireDate string `json:"expire_date"`
}

// USDTCallback represents the callback data received from GSPAY2 for USDT payments.
type USDTCallback struct {
	// CryptoPaymentID is the unique payment ID.
	CryptoPaymentID string `json:"cryptopayment_id"`
	// Amount is the payment amount (with 2 decimal places).
	Amount string `json:"amount"`
	// TransactionID is the original transaction ID.
	TransactionID string `json:"transaction_id"`
	// Status is the payment status.
	Status constants.PaymentStatus `json:"status"`
	// Signature is the callback signature for verification.
	Signature string `json:"signature"`
}

// USDTService handles USDT payment operations.
type USDTService struct {
	client *client.Client
}

// NewUSDTService creates a new USDT payment service.
func NewUSDTService(c *client.Client) *USDTService {
	return &USDTService{client: c}
}

// Create creates a new USDT payment order using TRC20 network.
//
// The generated order expires after approximately 2 minutes.
//
// Signature formula: MD5(transaction_id + player_username + amount + operator_secret_key)
func (s *USDTService) Create(ctx context.Context, req *USDTRequest) (*USDTResponse, error) {
	// Validate amount (minimum 1.00 USDT)
	if req.Amount < constants.MinAmountUSDT {
		return nil, errors.NewValidationError("amount", "minimum amount is 1.00 USDT")
	}

	// Format amount with 2 decimal places
	formattedAmount := fmt.Sprintf("%.2f", req.Amount)

	// Generate signature: transaction_id + player_username + amount + secret_key
	signatureData := fmt.Sprintf("%s%s%s%s",
		req.TransactionID,
		req.Username,
		formattedAmount,
		s.client.SecretKey,
	)
	sig := signature.Generate(signatureData)

	// Build API request
	apiReq := usdtAPIRequest{
		TransactionID: req.TransactionID,
		Username:      req.Username,
		Amount:        formattedAmount,
		Signature:     sig,
	}

	endpoint := fmt.Sprintf("/v2/integrations/operators/%s/cryptocurrency/trc20/usdt", s.client.AuthKey)
	resp, err := s.client.Post(ctx, endpoint, apiReq)
	if err != nil {
		return nil, err
	}

	result, err := client.ParseData[USDTResponse](resp.Data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// VerifyCallback verifies the signature of a USDT payment callback.
//
// Callback Signature formula: MD5(cryptopayment_id + amount + transaction_id + status + secret_key)
//
// This method only verifies the signature. To also verify the source IP,
// use [USDTService.VerifyCallbackWithIP] instead.
func (s *USDTService) VerifyCallback(callback *USDTCallback) error {
	return s.verifyCallbackSignature(callback)
}

// VerifyCallbackWithIP verifies both the signature and source IP of a USDT payment callback.
//
// The sourceIP parameter should be the IP address of the callback request,
// typically obtained from [http.Request.RemoteAddr] or the X-Forwarded-For header.
//
// If the client was configured with [WithCallbackIPWhitelist], this method will
// verify that the source IP is in the whitelist before verifying the signature.
// If no whitelist was configured, IP verification is skipped.
func (s *USDTService) VerifyCallbackWithIP(callback *USDTCallback, sourceIP string) error {
	// Verify IP first (fast fail)
	if err := s.client.VerifyCallbackIP(sourceIP); err != nil {
		return err
	}

	// Then verify signature
	return s.verifyCallbackSignature(callback)
}

// verifyCallbackSignature performs the actual signature verification.
func (s *USDTService) verifyCallbackSignature(callback *USDTCallback) error {
	// Check required fields
	if callback.CryptoPaymentID == "" {
		return fmt.Errorf("%w: cryptopayment_id", errors.ErrMissingCallbackField)
	}
	if callback.Amount == "" {
		return fmt.Errorf("%w: amount", errors.ErrMissingCallbackField)
	}
	if callback.TransactionID == "" {
		return fmt.Errorf("%w: transaction_id", errors.ErrMissingCallbackField)
	}
	if callback.Signature == "" {
		return fmt.Errorf("%w: signature", errors.ErrMissingCallbackField)
	}

	// Format amount with 2 decimal places
	amount, err := strconv.ParseFloat(callback.Amount, 64)
	if err != nil {
		return errors.NewValidationError("amount", "invalid amount format")
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
	expectedSignature := signature.Generate(signatureData)

	// Constant-time comparison to prevent timing attacks
	if !signature.Verify(expectedSignature, callback.Signature) {
		return errors.ErrInvalidSignature
	}

	return nil
}
