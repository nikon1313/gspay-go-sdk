---
trigger: always_on
---

# DeepWiki MCP Server Instructions

This project has DeepWiki configured as an MCP server for accessing documentation and knowledge about Go packages and libraries.

## Available Tools

### 1. `mcp__deepwiki__ask_question`

Ask questions about Go packages, libraries, or programming concepts.

**Parameters:**
- `question` (string, required): The question to ask
- `repo_name` (string, optional): Repository name in format `owner/repo` (e.g., `stretchr/testify`)

**Example - General Question:**
```jsonc
{
  "question": "How do I use context.WithTimeout in Go?"
}
```

**Example - Repository-Specific Question:**
```jsonc
{
  "question": "How do I use assert.Equal?",
  "repo_name": "stretchr/testify"
}
```

### 2. `mcp__deepwiki__read_wiki_structure`

Get the documentation structure/outline for a repository.

**Parameters:**
- `repo_name` (string, required): Repository name in format `owner/repo`

**Example:**
```jsonc
{
  "repo_name": "stretchr/testify"
}
```

### 3. `mcp__deepwiki__read_wiki_contents`

Read specific documentation content from a repository's wiki.

**Parameters:**
- `repo_name` (string, required): Repository name in format `owner/repo`
- `path` (string, required): Path to the documentation page

**Example:**
```jsonc
{
  "repo_name": "stretchr/testify",
  "path": "assert"
}
```

## Use Cases for This Project

### Understanding Dependencies

This project uses the following dependencies that can be queried:

| Dependency | Repo Name | Common Questions |
|------------|-----------|------------------|
| testify | `stretchr/testify` | Assertions, mocking, test suites |

**Example Queries:**

```jsonc
// Learn about testify assertions
{
  "question": "What assertion methods are available in testify?",
  "repo_name": "stretchr/testify"
}

// Learn about mocking
{
  "question": "How do I create mocks with testify?",
  "repo_name": "stretchr/testify"
}
```

### Go Standard Library

Ask about Go standard library packages without specifying a repo:

```jsonc
// HTTP client
{
  "question": "How do I set timeout on http.Client in Go?"
}

// Context usage
{
  "question": "What is the difference between context.Background and context.TODO?"
}

// Crypto/MD5
{
  "question": "How do I generate MD5 hash in Go?"
}

// JSON handling
{
  "question": "How do I unmarshal JSON with unknown structure in Go?"
}
```

### Best Practices Queries

```jsonc
// Error handling
{
  "question": "What are Go best practices for error handling?"
}

// Testing
{
  "question": "How do I write table-driven tests in Go?"
}

// Project structure
{
  "question": "What is the recommended Go project layout?"
}
```

## When to Use DeepWiki

1. **Learning new packages**: Before using a new dependency
2. **Best practices**: When unsure about idiomatic Go patterns
3. **API reference**: When you need function signatures or usage examples
4. **Troubleshooting**: When encountering errors with third-party packages

## Combining with gopls

DeepWiki complements gopls:

- **gopls**: Navigates YOUR code (definitions, references, diagnostics)
- **DeepWiki**: Explains LIBRARY code and concepts

**Workflow Example:**
```
1. Use gopls `hover` to see a function signature
2. If it's from a library, use DeepWiki to understand usage
3. Use gopls `references` to see how it's used in your codebase
```
