---
trigger: always_on
---

# AI Agent Workflows

This document describes common workflows for AI agents working on this project.

## Workflow 1: Bug Investigation

When investigating a bug or unexpected behavior:

```
1. Get diagnostics for the relevant file
   → Use: mcp_gopls_go_diagnostics
   
2. Get file context to understand dependencies
   → Use: mcp_gopls_go_file_context
   
3. Find all references to understand usage
   → Use: mcp_gopls_go_symbol_references
   
4. Check test files for expected behavior
   → Read: *_test.go files
   
5. Fix the bug and verify
   → Run: go test ./...
```

## Workflow 2: Adding New Feature

When implementing a new feature:

```
1. Understand existing patterns
   → Use: mcp_gopls_go_package_api on similar packages
   
2. Check for similar implementations
   → Use: mcp_gopls_go_search to find related symbols
   → Use: mcp_gopls_go_symbol_references to find usage patterns
   
3. Create the implementation following project conventions
   → Follow patterns in project.md
   
4. Write tests
   → Use table-driven tests with testify
   
5. Verify implementation
   → Run: go build ./...
   → Run: go test ./...
   → Run: go vet ./...
```

## Workflow 3: Code Review

When reviewing code changes:

```
1. Get package API overview
   → Use: mcp_gopls_go_package_api
   
2. Check for diagnostics (errors/warnings)
   → Use: mcp_gopls_go_diagnostics
   
3. Understand file dependencies
   → Use: mcp_gopls_go_file_context
   
4. Check test coverage
   → Run: go test ./... -cover
   
5. Verify code style
   → Check license headers
   → Check naming conventions
   → Run: go fmt ./...
```

## Workflow 4: Understanding Unfamiliar Code

When exploring code you don't understand:

```
1. Get package API summary
   → Use: mcp_gopls_go_package_api
   
2. Understand file dependencies
   → Use: mcp_gopls_go_file_context

3. Search for specific symbols
   → Use: mcp_gopls_go_search
   
4. For library functions:
   → Use: DeepWiki MCP for documentation
   
5. See how symbols are used
   → Use: mcp_gopls_go_symbol_references
```

## Workflow 5: Refactoring

When refactoring code:

```
1. Find all references to the target
   → Use: mcp_gopls_go_symbol_references
   
2. Understand the current implementation
   → Use: mcp_gopls_go_file_context
   → Use: mcp_gopls_go_package_api
   
3. Make changes carefully
   → Update all reference sites
   → Maintain backward compatibility if public API
   
4. Verify refactoring
   → Use: mcp_gopls_go_diagnostics
   → Run: go build ./...
   → Run: go test ./...
```

## Workflow 6: Writing Tests

When writing or improving tests:

```
1. Understand what needs testing
   → Use: mcp_gopls_go_package_api to list exported functions
   
2. Check existing test patterns
   → Read: existing *_test.go files
   
3. For testify usage questions
   → Use: DeepWiki MCP with repo_name: "stretchr/testify"
   
4. Write tests following table-driven pattern
   → Use mock HTTP servers for API tests
   → Use testify assertions
   
5. Verify test coverage
   → Run: go test ./... -cover -v
```

## Workflow 7: Debugging API Issues

When debugging API-related issues:

```
1. Check the request building
   → Find the API request struct
   → Verify signature generation
   
2. Check the response parsing
   → Use: mcp_gopls_go_file_context on request.go
   → Verify JSON field tags
   
3. Check error handling
   → Use: mcp_gopls_go_symbol_references on error types
   
4. Add test case reproducing the issue
   → Use httptest.NewServer for mocking
   
5. Fix and verify
   → Run tests with -v for verbose output
```

## Tool Selection Guide

| Task | Primary Tool | Secondary Tool |
|------|--------------|----------------|
| Find symbol usages | `mcp_gopls_go_symbol_references` | - |
| Search for symbols | `mcp_gopls_go_search` | - |
| Get package API | `mcp_gopls_go_package_api` | - |
| Check for errors | `mcp_gopls_go_diagnostics` | `go vet` |
| Understand file deps | `mcp_gopls_go_file_context` | - |
| Learn about libraries | DeepWiki MCP | - |
| Run tests | `go test` (bash) | - |
| Build verification | `go build` (bash) | `go vet` |

## Tips

1. **Always use absolute paths** when calling gopls tools
2. **Symbol names are case-insensitive** for fuzzy search
3. **Qualify symbols with package name** when using go_symbol_references if needed
4. **Run tests frequently** to catch issues early
5. **Check diagnostics** before committing changes
6. **Combine tools** for comprehensive understanding
