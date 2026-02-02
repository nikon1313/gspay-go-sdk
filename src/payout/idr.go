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

// Package payout provides payout-related functionality for the GSPAY2 SDK.
package payout

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	amountfmt "github.com/H0llyW00dzZ/gspay-go-sdk/src/helper/amount"
)

// IDRRequest represents a request to create an IDR payout (withdrawal).
type IDRRequest struct {
	// TransactionID is a unique transaction ID.
	TransactionID string `json:"transaction_id"`
	// Username is the customer ID or username.
	Username string `json:"player_username"`
	// AccountName is the recipient's bank account name.
	AccountName string `json:"account_name"`
	// AccountNumber is the recipient's bank account number.
	AccountNumber string `json:"account_number"`
	// Amount is the payout amount in IDR (no decimals).
	Amount int64 `json:"amount"`
	// BankCode is the target bank code (see constants.BanksIDR).
	BankCode string `json:"bank_target"`
	// Description is an optional transaction description.
	Description string `json:"trx_description,omitempty"`
}

// idrAPIRequest is the internal API request structure.
type idrAPIRequest struct {
	TransactionID string `json:"transaction_id"`
	Username      string `json:"player_username"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Amount        int64  `json:"amount"`
	BankTarget    string `json:"bank_target"`
	Signature     string `json:"signature"`
	Description   string `json:"trx_description,omitempty"`
}

// IDRResponse represents the response from creating an IDR payout.
type IDRResponse struct {
	// IDRPayoutID is the unique payout ID assigned by GSPAY2.
	IDRPayoutID json.Number `json:"idrpayout_id"`
	// Status is the initial payout status.
	Status constants.PaymentStatus `json:"status"`
}

// IDRStatusResponse represents the response from querying IDR payout status.
type IDRStatusResponse struct {
	// IDRPayoutID is the unique payout ID.
	IDRPayoutID json.Number `json:"idrpayout_id"`
	// TransactionID is the transaction ID.
	TransactionID string `json:"transaction_id"`
	// AccountName is the recipient's account name.
	AccountName string `json:"account_name"`
	// AccountNumber is the recipient's account number.
	AccountNumber string `json:"account_number"`
	// Amount is the payout amount.
	Amount json.Number `json:"amount"`
	// Status is the current payout status.
	Status constants.PaymentStatus `json:"status"`
	// Completed indicates if the payout has been completed.
	Completed bool `json:"completed"`
	// PayoutSuccess indicates if the payout was successful.
	PayoutSuccess bool `json:"payout_success"`
	// Remark contains additional information about the payout.
	Remark string `json:"remark"`
	// Signature is the response signature.
	Signature string `json:"signature"`
}

// IDRCallback represents the callback data received from GSPAY2 for IDR payouts.
//
// According to GSPAY2 documentation, the callback contains:
//   - idrpayout_id: Unique payout ID (bigint)
//   - transaction_id: Unique transaction ID submitted
//   - account_name: Bank account name submitted
//   - account_number: Bank account number submitted
//   - amount: Amount submitted (decimal, 2 decimal places)
//   - completed: Payout completion stage (boolean)
//   - payout_success: Success status (boolean)
//   - remark: Bank transaction reference/status or error message
//   - signature: MD5 hash verification
//
// Signature formula: idrpayout_id + account_number + amount + transaction_id + operator_secret_key
type IDRCallback struct {
	// IDRPayoutID is the unique payout ID (bigint from GSPAY2).
	IDRPayoutID json.Number `json:"idrpayout_id"`
	// TransactionID is the original transaction ID.
	TransactionID string `json:"transaction_id"`
	// AccountName is the bank account name submitted.
	AccountName string `json:"account_name"`
	// AccountNumber is the recipient's account number.
	AccountNumber string `json:"account_number"`
	// Amount is the payout amount (decimal from GSPAY2).
	Amount json.Number `json:"amount"`
	// Completed indicates the payout completion stage.
	Completed bool `json:"completed"`
	// PayoutSuccess indicates if the payout was successful.
	PayoutSuccess bool `json:"payout_success"`
	// Remark indicates the bank transaction reference/status or error message.
	Remark string `json:"remark"`
	// Signature is the callback signature for verification.
	Signature string `json:"signature"`
}

// IDRService handles IDR payout operations.
type IDRService struct{ client *client.Client }

// NewIDRService creates a new IDR payout service.
func NewIDRService(c *client.Client) *IDRService { return &IDRService{client: c} }

// Create creates a new IDR payout (withdrawal) to an Indonesian bank account or e-wallet.
//
// Amount is deducted immediately from settlement balance.
//
// Signature formula: MD5(transaction_id + player_username + amount + account_number + operator_secret_key)
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
	// Validate bank code
	bankCode := strings.ToUpper(req.BankCode)
	if !constants.IsValidBankIDR(bankCode) {
		return nil, fmt.Errorf("%w: %s", errors.New(s.client.Language, errors.ErrInvalidBankCode), bankCode)
	}

	// Validate amount (minimum 10000 IDR)
	if req.Amount < constants.MinAmountIDR {
		return nil, errors.NewValidationError("amount", errors.GetMessage(s.client.Language, errors.KeyMinPayoutAmountIDR))
	}

	// Generate signature: transaction_id + player_username + amount + account_number + secret_key
	signatureData := fmt.Sprintf("%s%s%d%s%s",
		req.TransactionID,
		req.Username,
		req.Amount,
		req.AccountNumber,
		s.client.SecretKey,
	)
	sig := s.client.GenerateSignature(signatureData)

	// Build API request
	apiReq := idrAPIRequest{
		TransactionID: req.TransactionID,
		Username:      req.Username,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		BankTarget:    bankCode,
		Signature:     sig,
	}

	if req.Description != "" {
		apiReq.Description = req.Description
	}

	endpoint := fmt.Sprintf(constants.GetEndpoint(constants.EndpointPayoutIDRCreate), s.client.AuthKey)
	resp, err := s.client.Post(ctx, endpoint, apiReq)
	if err != nil {
		return nil, err
	}

	result, err := client.ParseData[IDRResponse](resp.Data, s.client.Language)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetStatus retrieves the current status of an IDR payout.
func (s *IDRService) GetStatus(ctx context.Context, transactionID string) (*IDRStatusResponse, error) {
	endpoint := fmt.Sprintf(constants.GetEndpoint(constants.EndpointPayoutIDRStatus), s.client.AuthKey)
	resp, err := s.client.Get(ctx, endpoint, map[string]string{
		"transaction_id": transactionID,
	})
	if err != nil {
		return nil, err
	}

	result, err := client.ParseData[IDRStatusResponse](resp.Data, s.client.Language)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// VerifySignature verifies a signature for IDR payout operations.
//
// This is a generic method that can be used to verify signatures from any GSPAY2 API response
// that includes signature verification (callbacks, etc.).
//
// Formula: MD5(id + account_number + amount + transaction_id + operator_secret_key)
// Note: Amount should be formatted with 2 decimal places (e.g., "10000.00").
func (s *IDRService) VerifySignature(id, accountNumber, amount, transactionID, receivedSignature string) error {
	lang := errors.Language(s.client.Language)

	// Check required fields
	if id == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "id")
	}
	if accountNumber == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "account_number")
	}
	if amount == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "amount")
	}
	if transactionID == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "transaction_id")
	}
	if receivedSignature == "" {
		return errors.New(lang, errors.ErrMissingCallbackField, "signature")
	}

	// Format amount with 2 decimal places
	formattedAmount, err := amountfmt.Format(amount, s.client.Language)
	if err != nil {
		return err
	}

	// Generate expected signature
	// Formula: MD5(id + account_number + amount + transaction_id + operator_secret_key)
	signatureData := fmt.Sprintf("%s%s%s%s%s",
		id,
		accountNumber,
		formattedAmount,
		transactionID,
		s.client.SecretKey,
	)
	expectedSignature := s.client.GenerateSignature(signatureData)

	// Constant-time comparison to prevent timing attacks
	if !s.client.VerifySignature(expectedSignature, receivedSignature) {
		return errors.New(lang, errors.ErrInvalidSignature)
	}

	return nil
}

// VerifyCallback verifies the signature of an IDR payout callback.
//
// Callback Signature formula: MD5(idrpayout_id + account_number + amount + transaction_id + operator_secret_key)
// Note: Amount in callback has 2 decimal places (e.g., "10000.00").
//
// This method only verifies the signature. To also verify the source IP,
// use [IDRService.VerifyCallbackWithIP] instead.
func (s *IDRService) VerifyCallback(callback *IDRCallback) error {
	return s.VerifySignature(
		string(callback.IDRPayoutID),
		callback.AccountNumber,
		string(callback.Amount),
		callback.TransactionID,
		callback.Signature,
	)
}

// VerifyCallbackWithIP verifies both the signature and source IP of an IDR payout callback.
//
// The sourceIP parameter should be the IP address of the callback request,
// typically obtained from [http.Request.RemoteAddr] or the X-Forwarded-For header.
//
// If the client was configured with [WithCallbackIPWhitelist], this method will
// verify that the source IP is in the whitelist before verifying the signature.
// If no whitelist was configured, IP verification is skipped.
func (s *IDRService) VerifyCallbackWithIP(callback *IDRCallback, sourceIP string) error {
	// Verify IP first (fast fail)
	if err := s.client.VerifyCallbackIP(sourceIP); err != nil {
		return err
	}

	// Then verify signature
	return s.VerifyCallback(callback)
}
