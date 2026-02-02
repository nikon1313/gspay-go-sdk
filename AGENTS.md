# AGENTS.md - Guidelines for AI Coding Agents

This document provides comprehensive guidelines for AI agents working on the GSPAY Go SDK (Unofficial). Follow these instructions to maintain code quality, consistency, and correctness.

## 1. Build, Test, and Lint Commands

Execute these commands from the project root.

```bash
# Build all packages to verify compilation
go build ./...

# Run all tests (standard)
go test ./...

# Run tests with verbose output (recommended for debugging)
go test ./... -v

# Run tests with coverage analysis
go test ./... -cover

# Run a single test function (e.g., TestIDRService_Create in payment package)
go test ./src/payment -run TestIDRService_Create -v

# Run a specific subtest (e.g., "creates_payment_successfully")
go test ./src/payment -run "TestIDRService_Create/creates_payment_successfully" -v

# Run tests for a specific package only
go test ./src/client/...

# Static analysis (run before committing)
go vet ./...

# Format code (standard Go formatting)
go fmt ./...

# Tidy dependencies (run after adding/removing imports)
go mod tidy
```

## 2. Project Structure

Understand the layout before adding new files.

```
src/
├── balance/     # Balance query service
├── client/      # HTTP client, request handling, retry logic (with jitter)
├── constants/   # Enums (Banks, Channels, Status) and config constants
├── errors/      # Sentinel errors, APIError, ValidationError, LocalizedError
├── helper/      # Shared utility packages
│   ├── amount/  # Amount formatting (2 decimal places, i18n support)
│   └── gc/      # Garbage collection utilities (bytebufferpool wrapper)
├── i18n/        # Internationalization (Language, MessageKey, translations)
├── internal/    # Internal packages not exposed to users
│   └── signature/ # MD5 signature generation and verification
├── payment/     # Payment services (IDR, USDT)
└── payout/      # Payout/Withdrawal services (IDR)
```

## 3. Code Style & Conventions

### License Header
**REQUIRED** at the top of every `.go` file:
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

### Imports
Group imports: Standard Lib > Internal Project > External.
```go
import (
    "context"
    "fmt"

    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"

    "github.com/stretchr/testify/assert"
)
```

### Naming & Types
- **Packages**: Lowercase, single word (e.g., `payment`, `client`).
- **Exported Types**: PascalCase (e.g., `IDRRequest`, `PaymentStatus`).
- **Unexported Types**: camelCase (e.g., `idrAPIRequest`).
- **Constants**: PascalCase.
- **Errors**: Prefix with `Err` (e.g., `ErrInvalidAmount`).

### Error Handling
Use the `src/errors` package. Support i18n where applicable.
```go
// Sentinel error
return nil, errors.ErrInvalidTransactionID

// Validation error with i18n message
return nil, errors.NewValidationError("amount", 
    errors.GetMessage(s.client.Language, errors.KeyMinAmountIDR))

// Wrapping API errors
if apiErr := errors.GetAPIError(err); apiErr != nil { /* ... */ }
```

### Amount Formatting
**ALWAYS** use `src/helper/amount` for formatting amounts in signatures/requests.
```go
import amountfmt "github.com/H0llyW00dzZ/gspay-go-sdk/src/helper/amount"

// Format float64 for signature
formatted := amountfmt.FormatFloat(req.Amount)
```

## 4. Service Pattern Implementation

Follow this pattern for new services:

```go
type MyService struct {
    client *client.Client
}

func NewMyService(c *client.Client) *MyService {
    return &MyService{client: c}
}

// Context as first arg. Request struct pointer as second.
func (s *MyService) Create(ctx context.Context, req *MyRequest) (*MyResponse, error) {
    // 1. Validate inputs (use constants)
    // 2. Format data (e.g., amounts)
    // 3. Generate signature (s.client.GenerateSignature)
    // 4. Build internal API request struct
    // 5. Execute request (s.client.Post/Get)
    // 6. Parse response (client.ParseData)
}
```

## 5. Testing Guidelines

Use `testify` for assertions and table-driven tests.

```go
func TestService_Method(t *testing.T) {
    // Setup mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Assert request details
        assert.Equal(t, "/expected/endpoint", r.URL.Path)
        // Write mock response
    }))
    defer server.Close()

    client := client.New("key", "secret", client.WithBaseURL(server.URL))
    
    t.Run("success case", func(t *testing.T) {
        // Act & Assert
    })
}
```

## 6. API Signature Formulas

| Operation | Formula (MD5) |
|-----------|---------------|
| IDR Pay | `transaction_id + player_username + amount + secret_key` |
| IDR Payout | `transaction_id + player_username + amount + account_number + secret_key` |
| USDT Pay | `transaction_id + player_username + amount + secret_key` |
| Callbacks | `... + status + secret_key` (Verify specific order in code) |

**Note**: Amounts in signatures must be formatted to 2 decimal places (e.g., "10000.00").
