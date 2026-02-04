---
trigger: always_on
---

# AI Agent Workflows

This document describes common workflows for AI agents working on this project.

## Workflow 1: Bug Investigation

When investigating a bug or unexpected behavior:

```
1. Get diagnostics for the relevant file
   → Use: mcp__gopls__diagnostics
   
2. Find the function definition
   → Use: mcp__gopls__definition
   
3. Get hover info for context
   → Use: mcp__gopls__hover
   
4. Find all references to understand usage
   → Use: mcp__gopls__references
   
5. Check test files for expected behavior
   → Read: *_test.go files
   
6. Fix the bug and verify
   → Run: go test ./...
```

## Workflow 2: Adding New Feature

When implementing a new feature:

```
1. Understand existing patterns
   → Use: mcp__gopls__documentSymbol on similar files
   
2. Check for similar implementations
   → Use: mcp__gopls__references to find usage patterns
   
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
1. Get file overview
   → Use: mcp__gopls__documentSymbol
   
2. Check for diagnostics (errors/warnings)
   → Use: mcp__gopls__diagnostics
   
3. Verify function signatures
   → Use: mcp__gopls__hover
   
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
1. Get file structure
   → Use: mcp__gopls__documentSymbol
   
2. For each symbol of interest:
   → Use: mcp__gopls__hover for quick info
   → Use: mcp__gopls__definition to see implementation
   
3. For interfaces:
   → Use: mcp__gopls__implementation to find implementors
   
4. For library functions:
   → Use: mcp__deepwiki__ask_question for documentation
   
5. See how it's used
   → Use: mcp__gopls__references
```

## Workflow 5: Refactoring

When refactoring code:

```
1. Find all references to the target
   → Use: mcp__gopls__references
   
2. Understand the current implementation
   → Use: mcp__gopls__definition
   → Use: mcp__gopls__hover
   
3. Check for interface implementations
   → Use: mcp__gopls__implementation
   
4. Make changes carefully
   → Update all reference sites
   → Maintain backward compatibility if public API
   
5. Verify refactoring
   → Use: mcp__gopls__diagnostics
   → Run: go build ./...
   → Run: go test ./...
```

## Workflow 6: Writing Tests

When writing or improving tests:

```
1. Understand what needs testing
   → Use: mcp__gopls__documentSymbol to list functions
   
2. Check existing test patterns
   → Read: existing *_test.go files
   
3. For testify usage questions
   → Use: mcp__deepwiki__ask_question with repo_name: "stretchr/testify"
   
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
   → Use: mcp__gopls__definition on ParseData
   → Verify JSON field tags
   
3. Check error handling
   → Use: mcp__gopls__references on error types
   
4. Add test case reproducing the issue
   → Use httptest.NewServer for mocking
   
5. Fix and verify
   → Run tests with -v for verbose output
```

## Tool Selection Guide

| Task | Primary Tool | Secondary Tool |
|------|--------------|----------------|
| Find where something is defined | `mcp__gopls__definition` | - |
| Find all usages | `mcp__gopls__references` | - |
| Get quick documentation | `mcp__gopls__hover` | `mcp__deepwiki__ask_question` |
| Check for errors | `mcp__gopls__diagnostics` | `go vet` |
| List file contents | `mcp__gopls__documentSymbol` | - |
| Find interface implementations | `mcp__gopls__implementation` | - |
| Learn about libraries | `mcp__deepwiki__ask_question` | - |
| Run tests | `go test` (bash) | - |
| Build verification | `go build` (bash) | `go vet` |

## Tips

1. **Always use absolute paths** when calling gopls tools
2. **Line numbers are 1-based** (first line is 1)
3. **Column numbers are byte-based** (important for UTF-8)
4. **Run tests frequently** to catch issues early
5. **Check diagnostics** before committing changes
6. **Use hover** before definition for quick checks
7. **Combine tools** for comprehensive understanding
