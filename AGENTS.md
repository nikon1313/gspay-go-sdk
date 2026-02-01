# AGENTS.md - Guidelines for AI Coding Agents

This document provides guidelines for AI agents working on the GSPAY Go SDK (Unofficial).

## Build, Test, and Lint Commands

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run a single test function
go test ./src/payment -run TestIDRService_Create -v

# Run a specific subtest
go test ./src/payment -run "TestIDRService_Create/creates_payment_successfully" -v

# Run tests for a specific package
go test ./src/client/...

# Static analysis
go vet ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

## Project Structure

```
src/
├── client/      # HTTP client, options, helpers
├── constants/   # Bank codes, channels, status codes
├── errors/      # Sentinel errors, APIError, ValidationError, LocalizedError
├── i18n/        # Internationalization (Language, MessageKey, translations)
├── payment/     # IDR and USDT payment services
├── payout/      # IDR payout service
├── balance/     # Balance query service
├── helper/      # Helper utilities
│   └── gc/      # Buffer pool management (bytebufferpool wrapper)
└── internal/    # Internal packages (signature)
```

## Code Style Guidelines

### License Header (Required on all .go files)

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

### Import Organization

Group imports in this order with blank lines between groups:
1. Standard library
2. Project internal packages
3. External dependencies

```go
import (
    "context"
    "fmt"

    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"

    "github.com/stretchr/testify/assert"
)
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Packages | lowercase, single word | `client`, `payment`, `payout` |
| Exported types | PascalCase | `PaymentStatus`, `IDRRequest` |
| Unexported types | camelCase | `idrAPIRequest` |
| Constants | PascalCase (exported) | `StatusSuccess`, `DefaultTimeout` |
| Sentinel errors | `Err` prefix | `ErrInvalidAmount`, `ErrInvalidSignature` |
| Constructors | `New` prefix | `NewIDRService`, `NewValidationError` |
| Options | `With` prefix | `WithTimeout`, `WithRetries` |

### Type Definitions

```go
// Exported request struct with JSON tags and doc comments
type IDRRequest struct {
    // TransactionID is a unique transaction ID (5-20 characters).
    TransactionID string `json:"transaction_id"`
    // Amount is the payment amount in IDR (no decimals).
    Amount int64 `json:"amount"`
}

// Unexported internal API struct
type idrAPIRequest struct {
    TransactionID string `json:"transaction_id"`
    Signature     string `json:"signature"`
}
```

### Service Pattern

```go
// Service struct
type IDRService struct {
    client *client.Client
}

// Constructor
func NewIDRService(c *client.Client) *IDRService {
    return &IDRService{client: c}
}

// Methods accept context.Context as first parameter
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
    // Implementation
}
```

### Error Handling

```go
// Return sentinel errors for known conditions
return nil, errors.ErrInvalidTransactionID

// Return validation errors with field context
return nil, errors.NewValidationError("amount", "minimum amount is 10000 IDR")

// Wrap errors with context using fmt.Errorf
return nil, fmt.Errorf("%w: %s", errors.ErrInvalidBankCode, bankCode)

// Check errors using errors.Is or type extraction
if apiErr := errors.GetAPIError(err); apiErr != nil {
    // Handle API error
}

// Localized error messages
return nil, errors.NewValidationError("amount", 
    errors.GetMessage(s.client.Language, errors.KeyMinAmountIDR))
```

### Internationalization (i18n)

```go
import "github.com/H0llyW00dzZ/gspay-go-sdk/src/i18n"

// Client with Indonesian error messages
c := client.New("auth", "secret", client.WithLanguage(i18n.Indonesian))

// Get translated message
msg := i18n.Get(i18n.Indonesian, i18n.MsgMinAmountIDR)
// Result: "jumlah minimum adalah 10000 IDR"

// Supported languages: i18n.English (default), i18n.Indonesian
```

### Functional Options Pattern

```go
type Option func(*Client)

func WithTimeout(timeout time.Duration) Option {
    return func(c *Client) {
        if timeout >= 5*time.Second {
            c.Timeout = timeout
        }
    }
}
```

## Testing Guidelines

### Test Structure (table-driven with testify)

```go
func TestFunction(t *testing.T) {
    t.Run("description of test case", func(t *testing.T) {
        // Arrange
        // Act
        // Assert
        assert.Equal(t, expected, actual)
        require.NoError(t, err)
    })
}
```

### Mock HTTP Server Pattern

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    assert.Equal(t, http.MethodPost, r.Method)
    
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

### Test Assertions

- Use `require` for critical checks that should stop the test
- Use `assert` for non-critical checks that can continue
- Use `assert.ErrorIs` for sentinel error checking
- Use `require.NotNil` before accessing pointer fields

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/stretchr/testify` | Test assertions (assert, require) |
| `github.com/valyala/bytebufferpool` | Buffer pool for efficient memory reuse |

## API Signature Formulas

| Operation | Formula |
|-----------|---------|
| IDR Payment | `MD5(transaction_id + player_username + amount + secret_key)` |
| IDR Payment Callback | `MD5(idrpayment_id + amount + transaction_id + status + secret_key)` |
| IDR Payout | `MD5(transaction_id + player_username + amount + account_number + secret_key)` |
| IDR Payout Callback | `MD5(idrpayout_id + account_number + amount + transaction_id + secret_key)` |
| USDT Payment | `MD5(transaction_id + player_username + amount + secret_key)` |

Note: Callback amounts have 2 decimal places (e.g., "10000.00").
