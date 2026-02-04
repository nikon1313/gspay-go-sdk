---
trigger: always_on
---

# GSPAY Go SDK (Unofficial) Project Instructions

This document provides guidance for AI agents working on the GSPAY Go SDK project.

## Project Overview

This is an **unofficial** Go SDK for the GSPAY2 Payment Gateway API. It was independently developed to provide Go language compatibility for integrating with GSPAY2 services.

> **Note**: This SDK is not affiliated with, endorsed by, or officially supported by GSPAY. It is a community-driven implementation.

**Supported Features:**
- IDR (Indonesian Rupiah) payments and payouts
- USDT cryptocurrency payments
- Balance queries
- Webhook callback verification

## Project Structure

```
gspay-go-sdk/
├── go.mod                      # Module: github.com/H0llyW00dzZ/gspay-go-sdk
├── go.sum
├── README.md
├── opencode.json               # AI agent configuration
├── .github/instructions/       # AI agent instructions
├── src/
│   ├── client/                 # HTTP client and configuration
│   │   ├── client.go          # Client struct, options pattern
│   │   ├── request.go         # HTTP requests with retry logic
│   │   └── helpers.go         # Utility functions
│   ├── constants/              # Constants and enums
│   │   ├── constants.go       # Base constants (URLs, limits)
│   │   ├── status.go          # PaymentStatus type and methods
│   │   ├── banks.go           # Bank codes (IDR, MYR, THB)
│   │   └── channels.go        # Payment channels
│   ├── errors/                 # Error types
│   │   └── errors.go          # APIError, ValidationError
│   ├── payment/                # Payment services
│   │   ├── idr.go             # IDR payment service
│   │   └── usdt.go            # USDT payment service
│   ├── payout/                 # Payout services
│   │   └── idr.go             # IDR payout service
│   ├── balance/                # Balance service
│   │   └── balance.go
│   └── internal/               # Internal packages
│       └── signature/          # MD5 signature utilities
└── examples/                   # Usage examples
    ├── basic/main.go
    └── webhook/main.go
```

## Code Style Guidelines

### License Header

All Go files must include this Apache 2.0 license header:

```go
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
```

### Package Documentation

Each package should have a doc comment in one of its files:

```go
// Package client provides the HTTP client for the GSPAY2 API.
package client
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `client`, `payment`, `payout`)
- **Exported types**: PascalCase (e.g., `PaymentStatus`, `IDRRequest`)
- **Unexported types**: camelCase (e.g., `idrAPIRequest`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Errors**: Start with `Err` prefix (e.g., `ErrInvalidAmount`)

### Service Pattern

Services follow this pattern:

```go
// Service struct
type IDRService struct {
    client *client.Client
}

// Constructor
func NewIDRService(c *client.Client) *IDRService {
    return &IDRService{client: c}
}

// Methods
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
    // Implementation
}
```

### Error Handling

Use typed errors from the `errors` package:

```go
// Return sentinel errors
return nil, errors.ErrInvalidTransactionID

// Return validation errors
return nil, errors.NewValidationError("amount", "minimum amount is 10000 IDR")

// Return API errors
return nil, &errors.APIError{
    Code:    resp.StatusCode,
    Message: "HTTP Error",
}
```

When wrapping errors, use `errors.New` which automatically wraps causes with `%w`:

```go
// Wrap with context and localization
return errors.New(s.client.Language, errors.ErrRequestFailed, err)
```

## Testing Guidelines

### Test File Naming

Tests go in `*_test.go` files in the same package:

- `src/client/client_test.go`
- `src/payment/idr_test.go`

### Test Structure

Use table-driven tests and testify:

```go
func TestFunction(t *testing.T) {
    t.Run("description", func(t *testing.T) {
        // Arrange
        // Act
        // Assert using testify
        assert.Equal(t, expected, actual)
        require.NoError(t, err)
    })
}
```

### Mock HTTP Server

Use `httptest` for API tests:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Verify request
    assert.Equal(t, http.MethodPost, r.Method)
    
    // Return mock response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code":    200,
        "message": "success",
        "data":    `{"payment_url":"https://example.com"}`,
    })
}))
defer server.Close()

c := client.New("auth", "secret", client.WithBaseURL(server.URL))
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package
go test ./src/client/...
```

## Common Tasks

### Adding a New Payment Method

1. Create request/response structs in appropriate package
2. Add service method following existing patterns
3. Add signature generation if needed
4. Add callback verification if needed
5. Write tests with mock server
6. Update README with examples

### Adding a New Bank Code

1. Edit `src/constants/banks.go`
2. Add to appropriate map (`BanksIDR`, `BanksMYR`, `BanksTHB`)
3. Add test case in `src/constants/banks_test.go`

### Modifying API Endpoints

1. Edit `src/constants/endpoints.go` to add or modify the `EndpointKey` and path in the `endpoints` map.
2. Update the service implementation to use `constants.GetEndpoint()` instead of hardcoded strings.
3. Update signature generation if parameters change.
4. Update request/response structs.
5. Update tests to verify changes and ensure coverage for new endpoints.

## API Signature Formulas

### IDR Payment
```
MD5(transaction_id + player_username + amount + operator_secret_key)
```

### IDR Payment Callback
```
MD5(idrpayment_id + amount + transaction_id + status + secret_key)
Note: amount has 2 decimal places (e.g., "10000.00")
```

### IDR Payout
```
MD5(transaction_id + player_username + amount + account_number + operator_secret_key)
```

### IDR Payout Callback
```
MD5(idrpayout_id + account_number + amount + transaction_id + secret_key)
```

### USDT Payment
```
MD5(transaction_id + player_username + amount + operator_secret_key)
```

### USDT Payment Callback
```
MD5(cryptopayment_id + amount + transaction_id + status + secret_key)
```

## Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/stretchr/testify` | Testing assertions | v1.11.1 |

## Build & Verify

```bash
# Build all packages
go build ./...

# Run static analysis
go vet ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```
