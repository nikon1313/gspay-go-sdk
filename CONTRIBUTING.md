# Contributing to GSPAY Go SDK

Thank you for your interest in contributing to the GSPAY Go SDK! This document provides guidelines and information for contributors.

## 🚀 Ways to Contribute

- **🐛 Bug Reports**: Report bugs via [GitHub Issues](https://github.com/H0llyW00dzZ/gspay-go-sdk/issues)
- **💡 Feature Requests**: Suggest new features or improvements
- **📝 Documentation**: Improve documentation, examples, or guides
- **💻 Code Contributions**: Submit pull requests for new features or bug fixes
- **🧪 Testing**: Add tests or improve test coverage

## 📋 Development Setup

### Prerequisites

- Go 1.25.6 or later
- Git
- Basic understanding of Go modules and testing

### Setup Steps

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gspay-go-sdk.git
   cd gspay-go-sdk
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Run tests** to ensure everything works:
   ```bash
   go test ./...
   ```

5. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 🏗️ Project Structure

```
gspay-go-sdk/
├── src/
│   ├── balance/     # Balance query service
│   ├── client/      # HTTP client and core functionality
│   ├── constants/   # Bank codes, payment statuses, channels
│   ├── errors/      # Error types and handling
│   ├── helper/      # Helper utilities
│   ├── internal/    # Internal utilities (signature generation)
│   ├── payment/     # Payment services (IDR, future THB/MYR)
│   └── payout/      # Payout services (IDR)
├── examples/        # Usage examples
├── go.mod           # Go module definition
└── README.md        # Main documentation
```

## 💻 Code Standards

### Go Code Style

- Follow standard Go formatting: `go fmt`
- Use `gofmt -s` for additional simplifications
- Run `go vet` and fix all warnings
- Ensure `golint` passes (if available)

### Naming Conventions

```go
// Types
type PaymentRequest struct { ... }    // PascalCase for exported types
type paymentAPIRequest struct { ... } // camelCase for internal types

// Functions
func CreatePayment(...) (...)         // PascalCase for exported functions
func createAPIRequest(...) (...)      // camelCase for internal functions

// Variables
var PaymentStatusPending = 0          // PascalCase for exported constants
var defaultTimeout = 30 * time.Second // camelCase for internal variables
```

### Error Handling

- Return typed errors from the `errors` package
- Use `fmt.Errorf` for wrapping: `return fmt.Errorf("%w: %s", errors.ErrInvalidAmount, amount)`
- Include context in error messages

### Documentation

- Add doc comments for all exported functions, types, and methods
- Use proper Go doc format
- Include usage examples where helpful

## 🧪 Testing

### Test Requirements

- **100% coverage** for new code
- Use table-driven tests for multiple scenarios
- Mock HTTP responses using `httptest`
- Test both success and error cases
- Test edge cases and input validation

### Test Structure

```go
func TestPaymentService_Create(t *testing.T) {
    t.Run("successful payment creation", func(t *testing.T) {
        // Setup mock server
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mock API response
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]any{
                "code": 200,
                "data": `{"payment_url":"https://...","id":"123"}`,
            })
        }))
        defer server.Close()

        // Test your code
        client := New("auth", "secret", WithBaseURL(server.URL))
        svc := payment.NewIDRService(client)

        resp, err := svc.Create(context.Background(), &payment.IDRRequest{
            TransactionID: "TXN123",
            Username:      "user123",
            Amount:        50000,
        })

        assert.NoError(t, err)
        assert.NotEmpty(t, resp.PaymentURL)
    })
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package
go test ./src/payment

# Run with verbose output
go test ./... -v
```

## 🔧 Adding New Payment Methods

### For New Currencies (e.g., THB, MYR)

1. **Add constants** in `src/constants/`:
   ```go
   // Currency codes
   const CurrencyTHB Currency = "THB"

   // Bank codes
   var BanksTHB = map[string]string{
       "BBL": "Bangkok Bank",
       // ... add more banks
   }

   // Payment channels
   var ChannelsTHB = []string{"QRIS", "BANK_TRANSFER"}
   ```

2. **Create payment service** in `src/payment/`:
   ```go
   // src/payment/thb.go
   type THBService struct { client *client.Client }

   func NewTHBService(c *client.Client) *THBService {
       return &THBService{client: c}
   }

   func (s *THBService) Create(ctx context.Context, req *THBRequest) (*THBResponse, error) {
       // Implementation following IDR service pattern
   }
   ```

3. **Add callback verification**:
   ```go
   func (s *THBService) VerifyCallback(callback *THBCallback) error {
       // MD5 signature verification
   }
   ```

4. **Update client** to support new service:
   ```go
   // Add THB service constructor
   func NewTHBService(c *Client) *THBService { ... }
   ```

5. **Add comprehensive tests** following the existing patterns

6. **Update documentation** in README.md and examples

### Implementation Checklist

- [ ] Constants added for currency, banks, channels
- [ ] Payment service implemented with proper error handling
- [ ] Callback verification implemented
- [ ] Unit tests with 100% coverage
- [ ] Integration tests with mock API
- [ ] Documentation updated
- [ ] Examples added
- [ ] Changelog updated

## 📝 Pull Request Process

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/add-thb-support
   ```

2. **Make your changes** following the code standards

3. **Run tests** and ensure they pass:
   ```bash
   go test ./... -v
   go vet ./...
   ```

4. **Update documentation** if needed

5. **Commit your changes** with clear messages:
   ```bash
   git add .
   git commit -m "feat: add THB payment support

   - Implement THB payment service
   - Add callback verification
   - Add comprehensive tests

   Closes #123"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/add-thb-support
   ```

7. **Create a Pull Request** on GitHub:
   - Use a clear title and description
   - Reference any related issues
   - Include screenshots/demo for UI changes
   - Request review from maintainers

### PR Title Format

```
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
Examples:
- feat(thb): add THB payment support
- fix(callback): resolve signature verification bug
- docs(readme): update installation instructions
```

## 🐛 Bug Reports

When reporting bugs, please include:

- **Go version**: `go version`
- **SDK version**: Git commit hash or tag
- **Expected behavior**
- **Actual behavior**
- **Steps to reproduce**
- **Error messages/logs**
- **Code sample** that demonstrates the issue

## 💡 Feature Requests

Feature requests should include:

- **Use case**: What problem does this solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other approaches?
- **Additional context**: Screenshots, examples, etc.

## 📜 Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors:

- Be respectful and inclusive
- Focus on constructive feedback
- Accept responsibility for mistakes
- Show empathy towards other contributors
- Help create a positive community

## 📞 Getting Help

- **Documentation**: Check README.md and examples first
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions
- **Community**: Join relevant Go communities for general questions

## 🎉 Recognition

Contributors will be recognized:
- In the CHANGELOG for significant contributions
- As co-authors on releases
- In the project's contributor list
- Through GitHub's contributor insights

Thank you for contributing to the GSPAY Go SDK! 🚀
