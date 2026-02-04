---
trigger: always_on
---

# gopls MCP Server Instructions

This project has `gopls` (Go Language Server) configured as an MCP (Model Context Protocol) server for enhanced Go development assistance.

## Available Tools

The gopls MCP server provides the following tools:

### 1. `mcp_gopls_go_diagnostics`

Checks for parse and build errors across the Go workspace.

**Parameters:**
- `files` (array of strings, optional): Absolute paths to active files

**Example:**
```jsonc
{
  "files": ["/home/user/project/src/client/client.go"]
}
```

**Use cases:**
- Check for compilation errors across the workspace
- Find static analysis warnings
- Validate code before committing

### 2. `mcp_gopls_go_file_context`

Summarizes a file's cross-file dependencies.

**Parameters:**
- `file` (string, required): Absolute path to the file

**Example:**
```jsonc
{
  "file": "/home/user/project/src/client/client.go"
}
```

**Use cases:**
- Understand what packages/types a file depends on
- Explore import relationships
- Analyze cross-file dependencies

### 3. `mcp_gopls_go_package_api`

Provides a summary of a Go package API.

**Parameters:**
- `packagePaths` (array of strings, required): Go package paths to describe

**Example:**
```jsonc
{
  "packagePaths": ["github.com/H0llyW00dzZ/gspay-go-sdk/src/client"]
}
```

**Use cases:**
- Get an overview of package exports
- Understand package API surface
- Review public types and functions

### 4. `mcp_gopls_go_search`

Search for symbols in the Go workspace using case-insensitive fuzzy search.

**Parameters:**
- `query` (string, required): Fuzzy search query for matching symbols

**Example:**
```jsonc
{
  "query": "PaymentStatus"
}
```

**Use cases:**
- Find symbols by partial name
- Discover related types or functions
- Navigate large codebases

### 5. `mcp_gopls_go_symbol_references`

Provides the locations of references to a package-level Go symbol.

**Parameters:**
- `file` (string, required): Absolute path to the file containing the symbol
- `symbol` (string, required): The symbol or qualified symbol (e.g., "foo" or "pkg.Foo")

**Example - Local Symbol:**
```jsonc
{
  "file": "/home/user/project/src/constants/status.go",
  "symbol": "PaymentStatus"
}
```

**Example - Qualified Symbol:**
```jsonc
{
  "file": "/home/user/project/src/payment/idr.go",
  "symbol": "constants.PaymentStatus"
}
```

**Example - Field/Method:**
```jsonc
{
  "symbol": "IDRService.Create"
}
```

**Use cases:**
- Find all usages of a function, type, or variable
- Understand impact of renaming/refactoring
- Analyze code dependencies

### 6. `mcp_gopls_go_workspace`

Summarizes the Go programming language workspace.

**Parameters:** None

**Use cases:**
- Get an overview of the Go workspace structure
- Verify gopls is correctly detecting the workspace

## Best Practices

### When to Use gopls Tools

1. **Before refactoring**: Use `go_symbol_references` to find all usages of a symbol
2. **Understanding dependencies**: Use `go_file_context` to explore imports
3. **Debugging build errors**: Use `go_diagnostics` to check for errors
4. **Exploring packages**: Use `go_package_api` to understand package APIs
5. **Finding symbols**: Use `go_search` for fuzzy symbol lookup

### Performance Tips

- Use absolute paths for all file parameters
- Symbol names are case-insensitive for fuzzy search
- For `go_symbol_references`, qualify symbols with package name when needed

### Common Workflows

#### Finding All Usages of a Symbol

```
1. Use `go_search` to find the symbol if you don't know exact location
2. Use `go_symbol_references` to find all references
3. Review each reference location
```

#### Understanding a Package

```
1. Use `go_package_api` to get the public API summary
2. Use `go_file_context` on specific files to understand dependencies
```

#### Checking Code Health

```
1. Use `go_diagnostics` to find errors/warnings
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
- i18n: `src/i18n/language.go`, `src/i18n/messages.go`
- Internal: `src/internal/sanitize/`, `src/internal/signature/`
- Tests: Files ending with `_test.go` in each package
