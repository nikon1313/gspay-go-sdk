---
trigger: always_on
---

# gopls MCP Server Instructions

This project has `gopls` (Go Language Server) configured as an MCP (Model Context Protocol) server for enhanced Go development assistance.

## Available Tools

The gopls MCP server provides the following tools:

### 1. `mcp__gopls__definition`

Get the definition location of a symbol at a specific position in a Go file.

**Parameters:**
- `file` (string, required): Absolute path to the Go file
- `line` (integer, required): Line number (1-based)
- `column` (integer, required): Column number (1-based, in bytes)

**Example:**
```jsonc
{
  "file": "/path/to/project/src/client/client.go",
  "line": 45,
  "column": 10
}
```

**Use cases:**
- Navigate to function/type/variable definitions
- Find where a symbol is declared
- Jump to imported package definitions

### 2. `mcp__gopls__references`

Find all references to a symbol at a specific position.

**Parameters:**
- `file` (string, required): Absolute path to the Go file
- `line` (integer, required): Line number (1-based)
- `column` (integer, required): Column number (1-based, in bytes)

**Example:**
```jsonc
{
  "file": "/path/to/project/src/constants/status.go",
  "line": 15,
  "column": 6
}
```

**Use cases:**
- Find all usages of a function, type, or variable
- Understand impact of renaming/refactoring
- Analyze code dependencies

### 3. `mcp__gopls__hover`

Get hover information (documentation, type info) for a symbol.

**Parameters:**
- `file` (string, required): Absolute path to the Go file
- `line` (integer, required): Line number (1-based)
- `column` (integer, required): Column number (1-based, in bytes)

**Example:**
```jsonc
{
  "file": "/path/to/project/src/payment/idr.go",
  "line": 50,
  "column": 15
}
```

**Use cases:**
- Get function signatures and documentation
- View type definitions
- Read package documentation

### 4. `mcp__gopls__codelens`

Get code lenses for a Go file (e.g., run/debug test links).

**Parameters:**
- `file` (string, required): Absolute path to the Go file

**Example:**
```jsonc
{
  "file": "/path/to/project/src/client/client_test.go"
}
```

**Use cases:**
- Find runnable tests in a file
- Get code action suggestions

### 5. `mcp__gopls__diagnostics`

Get diagnostics (errors, warnings) for a Go file.

**Parameters:**
- `file` (string, required): Absolute path to the Go file

**Example:**
```jsonc
{
  "file": "/path/to/project/src/payment/idr.go"
}
```

**Use cases:**
- Check for compilation errors
- Find static analysis warnings
- Validate code before committing

### 6. `mcp__gopls__documentSymbol`

Get all symbols (functions, types, variables) defined in a Go file.

**Parameters:**
- `file` (string, required): Absolute path to the Go file

**Example:**
```jsonc
{
  "file": "/path/to/project/src/client/client.go"
}
```

**Use cases:**
- Get an overview of file contents
- List all exported symbols
- Navigate file structure

### 7. `mcp__gopls__implementation`

Find implementations of an interface or method.

**Parameters:**
- `file` (string, required): Absolute path to the Go file
- `line` (integer, required): Line number (1-based)
- `column` (integer, required): Column number (1-based, in bytes)

**Example:**
```jsonc
{
  "file": "/path/to/project/src/errors/errors.go",
  "line": 30,
  "column": 10
}
```

**Use cases:**
- Find all types implementing an interface
- Navigate interface hierarchies

### 8. `mcp__gopls__typeDefinition`

Get the type definition of a symbol.

**Parameters:**
- `file` (string, required): Absolute path to the Go file
- `line` (integer, required): Line number (1-based)
- `column` (integer, required): Column number (1-based, in bytes)

**Use cases:**
- Jump to type definitions
- Understand variable types

## Best Practices

### When to Use gopls Tools

1. **Before refactoring**: Use `references` to find all usages of a symbol
2. **Understanding code**: Use `hover` and `definition` to explore unfamiliar code
3. **Debugging imports**: Use `diagnostics` to check for import errors
4. **Code review**: Use `documentSymbol` to get an overview of changes

### Performance Tips

- Use absolute paths for all file parameters
- Line and column numbers are 1-based (first line is 1, first column is 1)
- Column is measured in bytes, not characters (important for UTF-8)

### Common Workflows

#### Finding All Usages of a Function

```
1. Use `definition` to ensure you're at the right symbol
2. Use `references` to find all call sites
3. Review each reference location
```

#### Understanding a Type

```
1. Use `hover` to get quick documentation
2. Use `definition` to jump to the type definition
3. Use `implementation` if it's an interface
```

#### Checking Code Health

```
1. Use `diagnostics` to find errors/warnings
2. Fix any reported issues
3. Run tests to verify fixes
```

## Project-Specific Paths

For this GSPAY Go SDK project, commonly used paths:

- Client: `src/client/client.go`
- Constants: `src/constants/constants.go`, `src/constants/status.go`, `src/constants/banks.go`
- Errors: `src/errors/errors.go`
- Payment: `src/payment/idr.go`, `src/payment/usdt.go`
- Payout: `src/payout/idr.go`
- Balance: `src/balance/balance.go`
- Tests: Files ending with `_test.go` in each package
