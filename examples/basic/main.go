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

// Package main demonstrates basic usage of the GSPAY Go SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/H0llyW00dzZ/gspay-go-sdk/src/balance"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
	"github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

func main() {
	// Get credentials from environment variables
	authKey := os.Getenv("GSPAY_AUTH_KEY")
	secretKey := os.Getenv("GSPAY_SECRET_KEY")

	if authKey == "" || secretKey == "" {
		log.Fatal("Please set GSPAY_AUTH_KEY and GSPAY_SECRET_KEY environment variables")
	}

	// Create client with custom options
	c := client.New(
		authKey,
		secretKey,
		client.WithTimeout(60*time.Second),
		client.WithRetries(3),
	)

	ctx := context.Background()

	// Create services
	paymentSvc := payment.NewIDRService(c)
	usdtSvc := payment.NewUSDTService(c)
	balanceSvc := balance.NewService(c)

	// Example 1: Create IDR Payment
	fmt.Println("=== Creating IDR Payment ===")
	paymentResp, err := paymentSvc.Create(ctx, &payment.IDRRequest{
		TransactionID: client.GenerateTransactionID("TXN"),
		Username:      "demo_user",
		Amount:        50000,
		Channel:       constants.ChannelQRIS,
	})
	if err != nil {
		if apiErr := errors.GetAPIError(err); apiErr != nil {
			log.Printf("API Error: %d - %s", apiErr.Code, apiErr.Message)
		} else {
			log.Printf("Error: %v", err)
		}
	} else {
		fmt.Printf("Payment URL: %s\n", paymentResp.PaymentURL)
		fmt.Printf("Payment ID: %s\n", paymentResp.IDRPaymentID)
		fmt.Printf("Expires: %s\n", paymentResp.ExpireDate)
		fmt.Printf("Amount: %s\n", client.FormatAmountIDR(50000))
		if paymentResp.QR != "" {
			fmt.Printf("QR Code Data: %s\n", paymentResp.QR)
		}
	}

	fmt.Println()

	// Example 2: Create USDT Payment
	fmt.Println("=== Creating USDT Payment ===")
	usdtResp, err := usdtSvc.Create(ctx, &payment.USDTRequest{
		TransactionID: client.GenerateTransactionID("USD"),
		Username:      "demo_user",
		Amount:        10.50,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Payment URL: %s\n", usdtResp.PaymentURL)
		fmt.Printf("Crypto Payment ID: %s\n", usdtResp.CryptoPaymentID)
		fmt.Printf("Amount: %s\n", client.FormatAmountUSDT(10.50))
	}

	fmt.Println()

	// Example 3: Check Balance
	fmt.Println("=== Checking Balance ===")
	balanceResp, err := balanceSvc.Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Balance: %s\n", balanceResp)
	}

	fmt.Println()

	// Example 4: List supported banks
	fmt.Println("=== Supported Indonesian Banks ===")
	for code, name := range constants.BanksIDR {
		fmt.Printf("  %s: %s\n", code, name)
	}
}
