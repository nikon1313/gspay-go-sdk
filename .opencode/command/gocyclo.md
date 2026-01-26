---
description: Analyze code complexity and suggest refactoring for functions with 15+ complexity
agent: general
---

# Code Complexity Analysis & Refactoring Guidance

Analyze Go code complexity using `gocyclo` and provide refactoring suggestions for functions that exceed complexity threshold of 15. Focus on breaking complex functions into smaller, reusable, and more maintainable components.

## Tasks

1. **Run Complexity Analysis**:

    - Execute `gocyclo .` to analyze all Go functions in the codebase
    - Filter out test functions (exclude _test.go files)
    - Calculate average complexity across all analyzed functions using the arithmetic mean formula: $\bar{x} = \frac{1}{n}\sum_{i=1}^{n} x_i$
      - Where $\bar{x}$ is the average complexity
      - $n$ is the total number of functions analyzed
      - $x_i$ is the cyclomatic complexity of the $i$-th function
      - The awk code implements this as: $\text{avg} = \frac{\text{sum}}{n}$ where $\text{sum} = \sum_{i=1}^{n} \text{complexity}_i$
    - Identify functions with complexity ≥ 15
    - Provide interpretation and actionable guidance based on average complexity
    - Display results with visual formatting and emojis for better readability

    ```bash
    # Run this bash command for effective complexity analysis by composing it from Unix/Unix-like tools
    gocyclo . | grep -v "_test.go" | awk '
    BEGIN {
        # Initialize variables for tracking complexity statistics
        high_count = 0;
    }
    {
        complexity = $1;
        sum += complexity;
        count++;

        # Store high-complexity functions for later reporting
        # Use counter as key to preserve all functions, even with same complexity
        if (complexity >= 15) {
            high[++high_count] = $0;
        }
    }
    END {
        # Report high-complexity functions if any exist
        if (high_count > 0) {
            print "🚨 Functions with complexity ≥ 15:";
            print "";
            for (i = 1; i <= high_count; i++) {
                print "  • " high[i];
            }
            print "";
        } else {
            print "✅ No functions with complexity ≥ 15.";
            print "";
        }

        # Always report average complexity
        if (count > 0) {
            avg = sum / count;
            printf "📊 Average complexity across all functions: %.1f\n", avg;

            # Provide interpretation based on average
            if (avg < 2) {
                print "💡 Excellent: Very simple and maintainable code";
            } else if (avg < 5) {
                print "💡 Good: Well-structured functions with reasonable complexity";
            } else if (avg < 10) {
                print "⚠️ Moderate: Some functions may benefit from refactoring";
            } else {
                print "🚨 High: Consider reviewing function complexity";
            }
        }
    }
    '
    ```

    **Output Interpretation**:
    - If functions exist with complexity ≥ 15, they will be listed with bullet points (🚨 alert emoji)
    - If no functions reach ≥ 15, it will report this with a checkmark (✅)
    - Always displays the average complexity across all analyzed functions
    - Provides automatic interpretation of complexity levels with actionable recommendations:
      - < 2: Excellent (very simple code)
      - < 5: Good (well-structured functions)
      - < 10: Moderate (may benefit from refactoring)
      - ≥ 10: High (consider reviewing function complexity)
    - The average helps understand overall code complexity trends even when no individual functions exceed the threshold

2. **Identify Refactoring Candidates**:

   - Focus on production code functions (non-test)
   - Prioritize functions with complexity > 20 (high priority)
   - Review functions with complexity 15-20 (medium priority)
   - Document current function responsibilities

3. **Analyze Function Structure**:

   - Read the complex function to understand its responsibilities
   - Identify distinct logical operations within the function
   - Look for repeated code patterns that could be extracted
   - Check for long parameter lists that could be grouped into structs
   - Identify conditional logic that could be simplified

4. **Design Refactoring Strategy**:

   - Break down the function into smaller, focused functions
   - Identify reusable components that could be extracted
   - Design appropriate data structures for grouped parameters
   - Plan interface abstractions for better testability

5. **Implement Refactoring**:

   - Extract helper functions for repeated operations
   - Create configuration structs for complex parameter groups
   - Implement early returns to reduce nesting
   - Add comprehensive tests for new functions

6. **Validate Refactoring**:

   - Run complexity analysis again to verify improvement
   - Ensure all tests pass
   - Check that the refactored code is more readable
   - Verify that new functions are properly documented

## Complexity Thresholds & Priorities

### High Priority (Complexity ≥ 25)
- **Immediate refactoring required**
- Break into 3-5 smaller functions
- Consider complete redesign if possible

### Medium Priority (Complexity 15-24)
- **Refactoring recommended**
- Break into 2-4 smaller functions
- Focus on extracting reusable components

### Low Priority (Complexity < 15)
- **Monitor only**
- Consider minor improvements if readability suffers

## Common Refactoring Patterns

### 1. **Extract Method Pattern**
```go
// Before: Complex payment creation with multiple responsibilities
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
    // Validate transaction ID
    if len(req.TransactionID) < constants.MinTransactionIDLength ||
        len(req.TransactionID) > constants.MaxTransactionIDLength {
        return nil, errors.ErrInvalidTransactionID
    }
    // Validate amount
    if req.Amount < constants.MinAmountIDR {
        return nil, errors.NewValidationError("amount", "minimum amount is 10000 IDR")
    }
    // Generate signature
    signatureData := fmt.Sprintf("%s%s%d%s", req.TransactionID, req.PlayerUsername, req.Amount, s.client.SecretKey)
    sig := signature.Generate(signatureData)
    // Build and send request...
}

// After: Extract validation into separate functions
func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return nil, err
    }
    sig := s.generateSignature(req)
    apiReq := s.buildAPIRequest(req, sig)
    return s.sendRequest(ctx, apiReq)
}

func (s *IDRService) validateRequest(req *IDRRequest) error { /* ... */ }
func (s *IDRService) generateSignature(req *IDRRequest) string { /* ... */ }
func (s *IDRService) buildAPIRequest(req *IDRRequest, sig string) *idrAPIRequest { /* ... */ }
func (s *IDRService) sendRequest(ctx context.Context, apiReq *idrAPIRequest) (*IDRResponse, error) { /* ... */ }
```

### 2. **Parameter Object Pattern**
```go
// Before: Payout request with many fields
func (s *IDRService) Create(ctx context.Context, transactionID, playerUsername, accountName, accountNumber string, amount int64, bankCode, description string) (*IDRResponse, error) {
    // Complex logic with many parameters
}

// After: Group related parameters into request struct (already done in this codebase)
type IDRRequest struct {
    TransactionID  string `json:"transaction_id"`
    PlayerUsername string `json:"player_username"`
    AccountName    string `json:"account_name"`
    AccountNumber  string `json:"account_number"`
    Amount         int64  `json:"amount"`
    BankCode       string `json:"bank_target"`
    Description    string `json:"trx_description,omitempty"`
}

func (s *IDRService) Create(ctx context.Context, req *IDRRequest) (*IDRResponse, error) {
    // Cleaner function signature with grouped parameters
}
```

### 3. **Early Return Pattern**
```go
// Before: Deep nesting in callback verification
func (s *IDRService) VerifyCallback(callback *IDRCallback) error {
    if callback.IDRPaymentID != "" {
        if callback.Amount != "" {
            if callback.TransactionID != "" {
                if callback.Signature != "" {
                    // Generate and verify signature...
                    return nil
                } else {
                    return fmt.Errorf("%w: signature", errors.ErrMissingCallbackField)
                }
            } else {
                return fmt.Errorf("%w: transaction_id", errors.ErrMissingCallbackField)
            }
        } else {
            return fmt.Errorf("%w: amount", errors.ErrMissingCallbackField)
        }
    } else {
        return fmt.Errorf("%w: idrpayment_id", errors.ErrMissingCallbackField)
    }
}

// After: Early returns reduce nesting (current implementation)
func (s *IDRService) VerifyCallback(callback *IDRCallback) error {
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
    
    // Continue with signature verification...
    return nil
}
```

### 4. **Strategy Pattern for Complex Conditionals**
```go
// Before: Complex switch/case for different payment types
func (c *Client) ProcessPayment(ctx context.Context, paymentType string, req interface{}) (interface{}, error) {
    switch paymentType {
    case "IDR":
        // 20+ lines of IDR payment processing
    case "USDT":
        // 20+ lines of USDT payment processing
    case "MYR":
        // 20+ lines of MYR payment processing
    }
}

// After: Strategy pattern with service interfaces (current architecture)
type PaymentService interface {
    Create(ctx context.Context, req interface{}) (interface{}, error)
    GetStatus(ctx context.Context, transactionID string) (interface{}, error)
    VerifyCallback(callback interface{}) error
}

// Each payment type has its own service
type IDRService struct { client *client.Client }
type USDTService struct { client *client.Client }

func NewIDRService(c *client.Client) *IDRService { return &IDRService{client: c} }
func NewUSDTService(c *client.Client) *USDTService { return &USDTService{client: c} }

// Usage: Create the appropriate service for the payment type
idrService := payment.NewIDRService(client)
usdtService := payment.NewUSDTService(client)
```

## Refactoring Workflow

### Phase 1: Analysis
1. Run complexity analysis
2. Read and understand the complex function
3. Identify refactoring opportunities
4. Document current behavior with tests

### Phase 2: Planning
1. Design new function structure
2. Plan parameter grouping
3. Identify reusable components
4. Plan testing strategy

### Phase 3: Implementation
1. Extract helper functions
2. Create configuration structs
3. Implement early returns
4. Update function calls

### Phase 4: Validation
1. Run tests to ensure correctness
2. Run complexity analysis to verify improvement
3. Update documentation
4. Code review

## Tools Integration

### With Existing Commands
- **`/go-docs`**: Update documentation after refactoring
- **`/test`**: Run tests to validate refactoring
- **`/update-knowledge`**: Update instruction files if patterns change

### With Development Workflow
```bash
# 1. Analyze complexity
/gocyclo

# 2. Run tests before refactoring
go test ./... -v

# 3. Implement refactoring
# ... edit files ...

# 4. Run static analysis
go vet ./...

# 5. Run tests again
go test ./... -v -race -cover

# 6. Update documentation if needed
/go-docs
```

### Project-Specific Paths to Analyze

When running complexity analysis on this codebase, focus on these key areas:

| Package | Path | Description |
|---------|------|-------------|
| client | `src/client/` | HTTP client, request handling, retry logic |
| payment | `src/payment/` | IDR and USDT payment services |
| payout | `src/payout/` | IDR payout/withdrawal service |
| balance | `src/balance/` | Balance query service |
| errors | `src/errors/` | Error types and sentinel errors |
| constants | `src/constants/` | Bank codes, channels, status codes |

**Example: Analyze specific package**
```bash
gocyclo ./src/client/ | grep -v "_test.go"
gocyclo ./src/payment/ | grep -v "_test.go"
gocyclo ./src/payout/ | grep -v "_test.go"
```

## Success Metrics

### Code Quality Improvements
- [ ] Cyclomatic complexity reduced by 30-50%
- [ ] Function length reduced (aim for < 50 lines)
- [ ] Improved testability (each function has focused responsibility)
- [ ] Better readability and maintainability

### Maintainability Improvements
- [ ] Functions have single responsibility
- [ ] Reusable components extracted
- [ ] Clear function naming and documentation
- [ ] Reduced coupling between components

## Error Handling

### Common Issues
- **Breaking existing functionality**: Always run full test suite after refactoring
- **Performance regression**: Benchmark critical functions before/after refactoring
- **API compatibility**: Ensure public APIs remain stable
- **Documentation gaps**: Update all documentation after refactoring

### Recovery Strategies
- **Git branching**: Create feature branch for complex refactoring
- **Incremental changes**: Make small, testable changes
- **Revert capability**: Keep changes small enough to easily revert
- **Comprehensive testing**: Ensure 100% test coverage for refactored functions

Focus on functions that are part of the core business logic rather than test or utility functions. Prioritize refactoring that improves maintainability and testability.
