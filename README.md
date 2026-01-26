# GSPAY Go SDK (Unofficial)

[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/gspay-go-sdk.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/gspay-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/gspay-go-sdk)](https://goreportcard.com/report/github.com/H0llyW00dzZ/gspay-go-sdk)

An **unofficial** Go SDK for the GSPAY2 Payment Gateway API. This SDK provides a comprehensive, idiomatic Go interface for payment processing, payouts, and balance queries.

> **Disclaimer**: This is an unofficial SDK and is not affiliated with, endorsed by, or officially supported by GSPAY. It was independently developed to provide Go compatibility for integrating with the GSPAY2 Payment Gateway API. Use at your own discretion.

## Features

- **IDR Payments**: Create payments via QRIS, DANA, and bank virtual accounts
- **IDR Payouts**: Process withdrawals to Indonesian bank accounts and e-wallets
- **USDT Payments**: Accept cryptocurrency payments via TRC20 network
- **Balance Queries**: Check operator settlement balance
- **Callback Verification**: Secure signature verification for webhooks
- **Retry Logic**: Automatic retries with exponential backoff for transient failures
- **Context Support**: Full context.Context support for cancellation and timeouts

## Installation

```bash
go get github.com/H0llyW00dzZ/gspay-go-sdk
```

## Project Structure

```
gspay-go-sdk/
├── src/
│   ├── client/      # HTTP client and configuration
│   ├── constants/   # Bank codes, channels, status codes
│   ├── errors/      # Error types and helpers
│   ├── payment/     # Payment services (IDR, USDT)
│   ├── payout/      # Payout services (IDR)
│   ├── balance/     # Balance query service
│   └── internal/    # Internal utilities (signature)
└── examples/        # Usage examples
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

func main() {
    // Create a new client
    c := client.New("your-auth-key", "your-secret-key")

    // Create payment service
    paymentSvc := payment.NewIDRService(c)

    ctx := context.Background()

    // Create an IDR payment
    resp, err := paymentSvc.Create(ctx, &payment.IDRRequest{
        TransactionID:  client.GenerateTransactionID("TXN"),
        Username:       "user123",
        Amount:         50000, // 50,000 IDR
        Channel:        constants.ChannelQRIS,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment URL: %s\n", resp.PaymentURL)
    fmt.Printf("Payment ID: %s\n", resp.IDRPaymentID)
    fmt.Printf("Expires: %s\n", resp.ExpireDate)
}
```

## Configuration Options

The client supports various configuration options using functional options pattern:

```go
c := client.New(
    "auth-key",
    "secret-key",
    client.WithBaseURL("https://custom-api.example.com"),
    client.WithTimeout(60 * time.Second),
    client.WithRetries(5),
    client.WithRetryWait(500*time.Millisecond, 5*time.Second),
    client.WithHTTPClient(customHTTPClient),
)
```

| Option | Description | Default |
|--------|-------------|---------|
| `WithBaseURL` | Set custom API base URL | `https://api.thegspay.com` |
| `WithTimeout` | Set request timeout | `30s` |
| `WithRetries` | Set number of retry attempts | `3` |
| `WithRetryWait` | Set min/max wait between retries | `500ms` / `2s` |
| `WithHTTPClient` | Use custom HTTP client | Default `http.Client` |

## Usage Examples

### Create IDR Payment

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
paymentSvc := payment.NewIDRService(c)

    resp, err := paymentSvc.Create(ctx, &payment.IDRRequest{
        TransactionID:  client.GenerateTransactionID("TXN"),
        Username:       "user123",
        Amount:         50000,
        Channel:        constants.ChannelQRIS, // Optional: QRIS, DANA, or BNI
    })
if err != nil {
    log.Fatal(err)
}

// Redirect user to payment page
fmt.Printf("Redirect to: %s\n", resp.PaymentURL)

// Optionally add return URL
redirectURL := client.BuildReturnURL(resp.PaymentURL, "https://mysite.com/complete")
```

### Check Payment Status

```go
status, err := paymentSvc.GetStatus(ctx, "TXN20260126143022123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status.String())
if status.Status.IsSuccess() {
    fmt.Println("Payment completed!")
}
```

### Create IDR Payout

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payout"
)

c := client.New("auth-key", "secret-key")
payoutSvc := payout.NewIDRService(c)

    resp, err := payoutSvc.Create(ctx, &payout.IDRRequest{
        TransactionID:  client.GenerateTransactionID("PAY"),
        Username:       "user123",
        AccountName:    "John Doe",
        AccountNumber:  "1234567890",
        Amount:         50000,
        BankCode:       "BCA",
        Description:    "Withdrawal request",
    })
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Payout ID: %s\n", resp.IDRPayoutID)
```

### Create USDT Payment

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
usdtSvc := payment.NewUSDTService(c)

    resp, err := usdtSvc.Create(ctx, &payment.USDTRequest{
        TransactionID:  client.GenerateTransactionID("USD"),
        Username:       "user123",
        Amount:         10.50, // 10.50 USDT
    })
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Payment URL: %s\n", resp.PaymentURL)
fmt.Printf("Crypto Payment ID: %s\n", resp.CryptoPaymentID)
```

### Check Balance

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/balance"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
)

c := client.New("auth-key", "secret-key")
balanceSvc := balance.NewService(c)

resp, err := balanceSvc.Get(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Balance: %s\n", resp.Balance)
```

### Verify Payment Callback

Handle webhooks from GSPAY2 securely:

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
paymentSvc := payment.NewIDRService(c)

func handleCallback(w http.ResponseWriter, r *http.Request) {
    var callback payment.IDRCallback
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Verify signature
    if err := paymentSvc.VerifyCallback(&callback); err != nil {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Process the callback
    if callback.Status.IsSuccess() {
        // Payment successful, update order status
        fmt.Printf("Payment %s completed for transaction %s\n",
            callback.IDRPaymentID, callback.TransactionID)
    }

    w.WriteHeader(http.StatusOK)
}
```

### Verify Payout Callback

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payout"
)

c := client.New("auth-key", "secret-key")
payoutSvc := payout.NewIDRService(c)

func handlePayoutCallback(w http.ResponseWriter, r *http.Request) {
    var callback payout.IDRCallback
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if err := payoutSvc.VerifyCallback(&callback); err != nil {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Process successful payout
    fmt.Printf("Payout %s completed\n", callback.IDRPayoutID)
    w.WriteHeader(http.StatusOK)
}
```

## Error Handling

The SDK provides typed errors for easy handling:

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
)

resp, err := paymentSvc.Create(ctx, req)
if err != nil {
    // Check for API errors
    if apiErr := errors.GetAPIError(err); apiErr != nil {
        log.Printf("API Error %d: %s", apiErr.Code, apiErr.Message)
        return
    }

    // Check for specific validation errors
    if errors.Is(err, errors.ErrInvalidTransactionID) {
        log.Println("Invalid transaction ID")
        return
    }

    if errors.Is(err, errors.ErrInvalidAmount) {
        log.Println("Invalid amount")
        return
    }

    // Handle other errors
    log.Printf("Error: %v", err)
}
```

## Supported Banks

### Indonesia (IDR)

| Code | Bank Name |
|------|-----------|
| `BCA` | Bank BCA |
| `BRI` | Bank BRI |
| `MANDIRI` | Bank Mandiri |
| `BNI` | Bank BNI |
| `CIMB` | Bank CIMB Niaga |
| `PERMATA` | Bank Permata |
| `DANAMON` | Bank Danamon Indonesia |
| `DANA` | DANA |
| `OVO` | OVO |

### Malaysia (MYR)

| Code | Bank Name |
|------|-----------|
| `MBB` | Maybank |
| `CIMB` | CIMB |
| `PBB` | Public Bank |
| `HLB` | Hong Leong Bank |
| `RHB` | RHB |
| `TNG` | Touch n Go eWallet |
| ... | [See full list](src/constants/banks.go) |

### Thailand (THB)

| Code | Bank Name |
|------|-----------|
| `BBL` | Bangkok Bank |
| `KBANK` | Kasikornbank |
| `KTB` | Krung Thai Bank |
| `SCB` | Siam Commercial Bank |
| ... | [See full list](src/constants/banks.go) |

## Payment Channels (IDR)

| Channel | Description |
|---------|-------------|
| `constants.ChannelQRIS` | QRIS QR Payment |
| `constants.ChannelDANA` | DANA E-Wallet |
| `constants.ChannelBNI` | BNI Virtual Account |

## Payment Status

| Status | Value | Description |
|--------|-------|-------------|
| `constants.StatusPending` | 0 | Payment pending or expired |
| `constants.StatusSuccess` | 1 | Payment successful |
| `constants.StatusFailed` | 2 | Payment failed |
| `constants.StatusTimeout` | 4 | Payment timed out |

```go
// Check status using helper methods
if status.IsSuccess() {
    // Payment completed
}

if status.IsFailed() {
    // Payment failed or timed out
}

if status.IsPending() {
    // Payment still pending
}

// Get human-readable label
fmt.Println(status.String()) // "Success", "Pending/Expired", etc.
```

## Helper Functions

### Generate Transaction ID

```go
// Generate unique transaction ID (max 20 chars)
txnID := client.GenerateTransactionID("TXN")
// Result: "TXN20260126143022123"
```

### Build Return URL

```go
// Add return URL to payment URL
fullURL := client.BuildReturnURL(paymentURL, "https://mysite.com/complete")
```

### Format Currency

```go
// Format IDR amount
formatted := client.FormatAmountIDR(50000)
// Result: "Rp 50.000"

// Format USDT amount
formatted := client.FormatAmountUSDT(10.50)
// Result: "10.50 USDT"
```

### Bank Utilities

```go
// Check if bank code is valid
if constants.IsValidBankIDR("BCA") {
    // Valid Indonesian bank
}

// Get bank name
name := constants.GetBankName("BCA", constants.CurrencyIDR)
// Result: "Bank BCA"

// Get all bank codes for a currency
codes := constants.GetBankCodes(constants.CurrencyIDR)
```

## Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover
```

## Disclaimer

This is an **unofficial** SDK. It is not affiliated with, endorsed by, or officially supported by GSPAY or its parent company. This SDK was independently developed by the community to provide Go language compatibility for integrating with the GSPAY2 Payment Gateway API.

The authors of this SDK are not responsible for any issues arising from its use. Please ensure you understand the GSPAY2 API terms of service before using this SDK in production.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
