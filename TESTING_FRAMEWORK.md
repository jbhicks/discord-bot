# Discord Bot Testing Framework

This document outlines the comprehensive testing strategy for the Discord bot, covering all tool calls and services with a focus on temporary file handling, mocking, and integration testing.

## Overview

The Discord bot consists of multiple interconnected services that need thorough testing to ensure reliability and proper functionality. This framework provides guidelines for testing each component while handling temporary file operations appropriately.

## Services and Tool Calls

### 1. LLM Service (`internal/llm/client.go`)
- **Tool Calls**: HTTP requests to llama.cpp server
- **Testing Focus**: 
  - JSON marshaling/unmarshaling
  - HTTP request/response handling
  - Error handling (timeouts, invalid responses)
  - Mock API responses

### 2. Image Generation Service (`internal/imagegen/client.go`)
- **Tool Calls**: HTTP requests to Stable Diffusion WebUI
- **Testing Focus**:
  - Request construction and validation
  - Base64 image decoding
  - API response parsing
  - Mock image generation

### 3. Stock News Services
- **MarketAux API** (`internal/stocknews/marketaux.go`)
- **AlphaVantage API** (`internal/stocknews/alphavantage.go`)
- **Fallback Mechanism** (`internal/stocknews/client.go`)
- **Testing Focus**:
  - API request construction
  - Response parsing and data conversion
  - Fallback behavior testing
  - Mock API responses

### 4. Sentiment Analysis Services
- **Reddit Scraper** (`internal/sentiment/reddit.go`)
- **X.com Scraper** (`internal/sentiment/x.go`)
- **Aggregator** (`internal/sentiment/aggregator.go`)
- **Testing Focus**:
  - Web scraping logic
  - Sentiment analysis algorithms
  - Fallback behavior
  - Mock scraping responses

### 5. Office Generation Services (`internal/officegen/`)
- **Document Generator**
- **Spreadsheet Generator** 
- **Presentation Generator**
- **Testing Focus**:
  - JSON parsing and data handling
  - File I/O operations (temporary files)
  - PDF generation workflow
  - Filename sanitization

### 6. Discord Command Implementations
- **Stock Command** (`pkg/commands/stock.go`)
- **Imagine Command** (`pkg/commands/imagine.go`)
- **PDF Command** (`pkg/commands/pdf.go`)
- **Testing Focus**:
  - Command data structures
  - Option parsing
  - Discord interaction handling
  - Error responses

## Temporary File Handling Strategy

### Test Directory Management
All temporary files created during testing will be:
1. Created in a dedicated test directory under `/tmp` or OS temp directory
2. Automatically cleaned up after each test run
3. Named with unique identifiers to prevent conflicts

### Implementation Details
```go
// Create test directory
testDir, err := os.MkdirTemp("", "discord-bot-test-*")
if err != nil {
    t.Fatalf("Failed to create test directory: %v", err)
}
defer os.RemoveAll(testDir) // Clean up after test

// Use test directory for temporary files
tempFile, err := os.Create(filepath.Join(testDir, "test-file.txt"))
```

### File Operations Testing
- Test file creation and deletion
- Test file content writing and reading
- Test file permissions and access
- Test cleanup procedures

## Mock Implementation Strategy

### Interface-Based Mocking
All external services will be implemented through interfaces to enable easy mocking:

```go
// Example interface for LLM service
type LLMClient interface {
    Chat(prompt string) (string, error)
    HealthCheck() error
}

// Real implementation
type llm.Client struct { /* ... */ }

// Mock implementation for testing
type MockLLMClient struct { /* ... */ }
func (m *MockLLMClient) Chat(prompt string) (string, error) { /* ... */ }
```

### Mock Service Implementations
- **Mock HTTP Clients**: Return predefined responses for testing
- **Mock API Responses**: Simulate various API scenarios
- **Mock File Operations**: Avoid actual disk I/O in unit tests
- **Mock External Dependencies**: Simulate service failures and edge cases

## Testing Approaches

### Unit Testing
- Test individual functions and methods in isolation
- Use mocks for external dependencies
- Test edge cases and error conditions
- Verify data structures and parsing logic

### Integration Testing
- Test interactions between multiple services
- Test full command execution flows
- Test database operations
- Test error handling and fallback mechanisms

### End-to-End Testing
- Test complete command execution from Discord interaction to response
- Test with various input parameters
- Test user feedback and error messages
- Test concurrent command execution

## Test Infrastructure

### Test Utilities
Create a `testutils` package with:
- Temporary file management functions
- Mock service factories
- Test data fixtures
- Helper functions for common testing patterns

### Test Structure
```
test/
├── unit/
│   ├── llm/
│   ├── imagegen/
│   ├── stocknews/
│   ├── sentiment/
│   └── officegen/
├── integration/
│   ├── command/
│   └── service/
└── fixtures/
    ├── news/
    ├── images/
    └── llm-responses/
```

### Test Execution
- Run unit tests with `go test ./...`
- Run integration tests with `go test -tags=integration ./...`
- Use `go test -v` for verbose output
- Use `go test -race` for race condition detection

## Testing Best Practices

### 1. Isolation
- Each test should be independent
- Use `t.Cleanup()` for automatic resource cleanup
- Avoid shared state between tests

### 2. Mocking Strategy
- Mock external dependencies completely
- Test both success and failure scenarios
- Use realistic mock data that matches real API responses

### 3. Temporary File Management
- Always clean up temporary files after tests
- Use unique identifiers for test files
- Test file operations without affecting production files

### 4. Error Handling
- Test all error paths
- Verify proper error messages are returned
- Test fallback mechanisms

## Example Test Structure

```go
func TestLLMClient_Chat(t *testing.T) {
    // Setup
    mockHTTPClient := &MockHTTPClient{
        Response: &http.Response{
            StatusCode: 200,
            Body:       io.NopCloser(strings.NewReader(`{"choices": [{"message": {"content": "test response"}}]}`)),
        },
    }
    
    client := llm.NewClient("http://localhost:8081")
    client.HTTPClient = mockHTTPClient
    
    // Execute
    response, err := client.Chat("test prompt")
    
    // Verify
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if response != "test response" {
        t.Errorf("Expected 'test response', got %q", response)
    }
}
```

## Running Tests

### Local Development
```bash
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/llm

# Run integration tests
go test -tags=integration ./...
```

### CI/CD Pipeline
- Run unit tests on every commit
- Run integration tests on scheduled basis
- Run performance tests with load testing
- Validate test coverage requirements

## Test Coverage Goals

### Unit Test Coverage
- Core business logic: 90%
- Data parsing and validation: 100%
- Error handling: 100%

### Integration Test Coverage
- Service interactions: 80%
- Command execution flows: 90%
- External API integrations: 80%

### End-to-End Test Coverage
- Full command execution: 70%
- User interaction scenarios: 70%
- Error scenarios: 80%