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

// Package main demonstrates webhook handling with the GSPAY Go SDK.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/payout"
)

var (
	paymentSvc *payment.IDRService
	usdtSvc    *payment.USDTService
	payoutSvc  *payout.IDRService
)

func main() {
	// Get credentials from environment variables
	authKey := os.Getenv("GSPAY_AUTH_KEY")
	secretKey := os.Getenv("GSPAY_SECRET_KEY")

	if authKey == "" || secretKey == "" {
		log.Fatal("Please set GSPAY_AUTH_KEY and GSPAY_SECRET_KEY environment variables")
	}

	// Create client and services
	c := client.New(authKey, secretKey)
	paymentSvc = payment.NewIDRService(c)
	usdtSvc = payment.NewUSDTService(c)
	payoutSvc = payout.NewIDRService(c)

	// Setup webhook handlers
	http.HandleFunc("/webhook/payment/idr", handlePaymentCallbackIDR)
	http.HandleFunc("/webhook/payout/idr", handlePayoutCallbackIDR)
	http.HandleFunc("/webhook/payment/usdt", handlePaymentCallbackUSDT)

	// Start server
	addr := ":8080"
	fmt.Printf("Starting webhook server on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /webhook/payment/idr  - IDR payment callbacks")
	fmt.Println("  POST /webhook/payout/idr   - IDR payout callbacks")
	fmt.Println("  POST /webhook/payment/usdt - USDT payment callbacks")

	log.Fatal(http.ListenAndServe(addr, nil))
}

// handlePaymentCallbackIDR handles IDR payment webhook callbacks.
func handlePaymentCallbackIDR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var callback payment.IDRCallback
	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		log.Printf("Failed to decode callback: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify signature
	if err := paymentSvc.VerifyCallback(&callback); err != nil {
		log.Printf("Invalid signature for transaction %s: %v", callback.TransactionID, err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Process the callback based on status
	log.Printf("Received IDR payment callback: txn=%s, payment_id=%s, amount=%s, status=%s",
		callback.TransactionID,
		callback.IDRPaymentID,
		callback.Amount,
		callback.Status.String(),
	)

	switch {
	case callback.Status.IsSuccess():
		// Payment successful - update order status, credit user account, etc.
		log.Printf("Payment successful: %s", callback.TransactionID)
		// TODO: Implement your business logic here

	case callback.Status.IsFailed():
		// Payment failed - notify user, cancel order, etc.
		log.Printf("Payment failed: %s", callback.TransactionID)
		// TODO: Implement your business logic here

	case callback.Status.IsPending():
		// Payment still pending
		log.Printf("Payment pending: %s", callback.TransactionID)
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handlePayoutCallbackIDR handles IDR payout webhook callbacks.
func handlePayoutCallbackIDR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var callback payout.IDRCallback
	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		log.Printf("Failed to decode callback: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify signature
	if err := payoutSvc.VerifyCallback(&callback); err != nil {
		log.Printf("Invalid signature for payout %s: %v", callback.TransactionID, err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Process the callback
	log.Printf("Received IDR payout callback: txn=%s, payout_id=%s, account=%s, amount=%s",
		callback.TransactionID,
		callback.IDRPayoutID,
		callback.AccountNumber,
		callback.Amount,
	)

	// TODO: Implement your business logic here
	// - Update withdrawal status
	// - Notify user

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handlePaymentCallbackUSDT handles USDT payment webhook callbacks.
func handlePaymentCallbackUSDT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var callback payment.USDTCallback
	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		log.Printf("Failed to decode callback: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify signature
	if err := usdtSvc.VerifyCallback(&callback); err != nil {
		log.Printf("Invalid signature for USDT payment %s: %v", callback.TransactionID, err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Process the callback
	log.Printf("Received USDT payment callback: txn=%s, payment_id=%s, amount=%s, status=%s",
		callback.TransactionID,
		callback.CryptoPaymentID,
		callback.Amount,
		callback.Status.String(),
	)

	if callback.Status.IsSuccess() {
		log.Printf("USDT Payment successful: %s", callback.TransactionID)
		// TODO: Implement your business logic here
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
