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

// Package payment provides payment-related functionality for the GSPAY2 SDK.
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

// IDRRequest represents a request to create an IDR payment.
type IDRRequest struct {
	// TransactionID is a unique transaction ID (5-20 characters).
	TransactionID string `json:"transaction_id"`
	// Username is the customer ID or username.
	Username string `json:"player_username"`
	// Amount is the payment amount in IDR (no decimals, e.g., 10000).
	Amount int64 `json:"amount"`
	// Channel is an optional payment channel (QRIS, DANA, BNI).
	// If omitted, user will select on the payment page.
	Channel constants.ChannelIDR `json:"channel,omitempty"`
}

// idrAPIRequest is the internal API request structure.
type idrAPIRequest struct {
	TransactionID string `json:"transaction_id"`
	Username      string `json:"player_username"`
	Amount        int64  `json:"amount"`
	Signature     string `json:"signature"`
	Channel       string `json:"channel,omitempty"`
}

// IDRResponse represents the response from creating an IDR payment.
type IDRResponse struct {
	// PaymentURL is the URL to redirect the user for payment.
	PaymentURL string `json:"payment_url"`
	// IDRPaymentID is the unique payment ID assigned by GSPAY2.
	IDRPaymentID string `json:"idrpayment_id"`
	// ExpireDate is the payment expiration date/time.
	ExpireDate string `json:"expire_date"`
	// Status is the initial payment status.
	Status constants.PaymentStatus `json:"status"`
}

// IDRStatusResponse represents the response from querying IDR payment status.
type IDRStatusResponse struct {
	// IDRPaymentID is the unique payment ID.
	IDRPaymentID string `json:"idrpayment_id"`
	// Status is the current payment status.
	Status constants.PaymentStatus `json:"status"`
	// Amount is the payment amount.
	Amount string `json:"amount"`
	// Remark contains additional information about the payment.
	Remark string `json:"remark"`
	// Signature is the response signature for verification.
	Signature string `json:"signature"`
}

// IDRCallback represents the callback data received from GSPAY2 for IDR payments.
type IDRCallback struct {
	// IDRPaymentID is the unique payment ID.
	IDRPaymentID string `json:"idrpayment_id"`
	// Amount is the payment amount (with 2 decimal places, e.g., "10000.00").
	Amount string `json:"amount"`
	// TransactionID is the original transaction ID.
	TransactionID string `json:"transaction_id"`
	// Status is the payment status.
	Status constants.PaymentStatus `json:"status"`
	// Signature is the callback signature for verification.
	Signature string `json:"signature"`
}

// IDRService handles IDR payment operations.
type IDRService struct {
	client *client.Client
}

// NewIDRService creates a new IDR payment service.
func NewIDRService(c *client.Client) *IDRService {
	return &IDRService{client: c}
}

// Create creates a new IDR payment order.
//
// The generated order expires after approximately 15 minutes.
//
// Signature formula: MD5(transaction_id + player_username + amount + operator_secret_key)
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
	// Validate transaction ID length
	if len(req.TransactionID) < constants.MinTransactionIDLength ||
		len(req.TransactionID) > constants.MaxTransactionIDLength {
		return nil, errors.ErrInvalidTransactionID
	}

	// Validate amount (minimum 10000 IDR)
	if req.Amount < constants.MinAmountIDR {
		return nil, errors.NewValidationError("amount", "minimum amount is 10000 IDR")
	}

	// Generate signature: transaction_id + player_username + amount + secret_key
	signatureData := fmt.Sprintf("%s%s%d%s",
		req.TransactionID,
		req.Username,
		req.Amount,
		s.client.SecretKey,
	)
	sig := signature.Generate(signatureData)

	// Build API request
	apiReq := idrAPIRequest{
		TransactionID: req.TransactionID,
		Username:      req.Username,
		Amount:        req.Amount,
		Signature:     sig,
	}

	// Add channel if specified
	if req.Channel != "" && constants.IsValidChannelIDR(req.Channel) {
		apiReq.Channel = string(req.Channel)
	}

	endpoint := fmt.Sprintf("/v2/integrations/operators/%s/idr/payment", s.client.AuthKey)
	resp, err := s.client.Post(ctx, endpoint, apiReq)
	if err != nil {
		return nil, err
	}

	result, err := client.ParseData[IDRResponse](resp.Data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetStatus retrieves the current status of an IDR payment order.
func (s *IDRService) GetStatus(ctx context.Context, transactionID string) (*IDRStatusResponse, error) {
	endpoint := fmt.Sprintf("/v2/integrations/operators/%s/idr/getpayment", s.client.AuthKey)
	resp, err := s.client.Get(ctx, endpoint, map[string]string{
		"transaction_id": transactionID,
	})
	if err != nil {
		return nil, err
	}

	result, err := client.ParseData[IDRStatusResponse](resp.Data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// VerifyCallback verifies the signature of an IDR payment callback.
//
// Callback Signature formula: MD5(idrpayment_id + amount + transaction_id + status + secret_key)
// Note: Amount in callback has 2 decimal places (e.g., "10000.00").
//
// This method only verifies the signature. To also verify the source IP,
// use [IDRService.VerifyCallbackWithIP] instead.
func (s *IDRService) VerifyCallback(callback *IDRCallback) error {
	return s.verifyCallbackSignature(callback)
}

// VerifyCallbackWithIP verifies both the signature and source IP of an IDR payment callback.
//
// The sourceIP parameter should be the IP address of the callback request,
// typically obtained from [http.Request.RemoteAddr] or the X-Forwarded-For header.
//
// If the client was configured with [WithCallbackIPWhitelist], this method will
// verify that the source IP is in the whitelist before verifying the signature.
// If no whitelist was configured, IP verification is skipped.
//
// Example:
//
//	func handleCallback(w http.ResponseWriter, r *http.Request) {
//	    sourceIP := r.RemoteAddr // or parse X-Forwarded-For
//	    if err := svc.VerifyCallbackWithIP(&callback, sourceIP); err != nil {
//	        // Handle error
//	    }
//	}
func (s *IDRService) VerifyCallbackWithIP(callback *IDRCallback, sourceIP string) error {
	// Verify IP first (fast fail)
	if err := s.client.VerifyCallbackIP(sourceIP); err != nil {
		return err
	}

	// Then verify signature
	return s.verifyCallbackSignature(callback)
}

// verifyCallbackSignature performs the actual signature verification.
func (s *IDRService) verifyCallbackSignature(callback *IDRCallback) error {
	// Check required fields
	if callback.IDRPaymentID == "" {
		return fmt.Errorf("%w: idrpayment_id", errors.ErrMissingCallbackField)
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
		callback.IDRPaymentID,
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
