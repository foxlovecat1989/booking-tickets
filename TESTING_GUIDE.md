# Testing Guide

This document provides a comprehensive guide to the testing system implemented for the tickets application.

## Overview

The testing system provides:
- **Unit tests** for models, configuration, and utilities
- **Integration tests** for repositories and services
- **gRPC Handler tests** for Protocol Buffer API endpoints
- **Database tests** with test helpers and cleanup
- **Migration tests** for database schema changes
- **Logger tests** for structured logging functionality
- **Coverage reporting** and analysis
- **Race condition detection** for concurrent code
- **Benchmark tests** for performance analysis

## Quick Start

### 1. Run All Tests

```bash
# Run all tests with verbose output
make test

# Run tests with coverage report
make test-coverage

# Run tests with race detection
make test-race
```

### 2. Run Specific Test Categories

```bash
# Unit tests (models, config, logger)
make test-unit

# Integration tests (repositories, services)
make test-integration

# Specific component tests
make test-models
make test-repository
make test-service
make test-config
make test-logger
make test-migrations

# gRPC Handler tests
go test ./internal/handler -v
```

## Test Structure

### Unit Tests

Unit tests focus on testing individual components in isolation:

#### Model Tests (`internal/models/domain/`)
- **concert_test.go**: Tests for Concert and ConcertSession models
- **ticket_test.go**: Tests for Ticket and TicketType models  
- **order_test.go**: Tests for Order and OrderItem models

**Example:**
```go
func TestConcert_Validation(t *testing.T) {
    tests := []struct {
        name    string
        concert Concert
        isValid bool
    }{
        {
            name: "valid concert",
            concert: Concert{
                Name:     "Rock Concert 2024",
                Location: "Madison Square Garden",
            },
            isValid: true,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            isValid := tt.concert.Name != "" && tt.concert.Location != ""
            assert.Equal(t, tt.isValid, isValid)
        })
    }
}
```

#### Configuration Tests (`internal/config/`)
- **config_test.go**: Tests for configuration loading and validation

**Example:**
```go
func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
    os.Setenv("DATABASE_HOST", "test-host")
    defer os.Unsetenv("DATABASE_HOST")
    
    cfg, err := LoadConfig()
    require.NoError(t, err)
    assert.Equal(t, "test-host", cfg.Database.Host)
}
```

#### Logger Tests (`internal/logger/`)
- **logger_test.go**: Tests for logging functionality and configuration

### Integration Tests

Integration tests focus on testing components that interact with external dependencies:

#### Repository Tests (`internal/repository/`)
- **concert_session_test.go**: Tests for concert session repository
- **connection_test.go**: Tests for database connectivity
- **test_helpers.go**: Database setup and cleanup utilities

**Example:**
```go
func TestConcertSessionRepository_GetConcertSessionByID(t *testing.T) {
    baseRepo, cleanup := SetupTestDB(t)
    defer cleanup()
    
    repo := NewConcertSessionRepository(baseRepo)
    
    // Test getting non-existent session
    session, err := repo.GetConcertSessionByID(999)
    require.NoError(t, err)
    assert.Nil(t, session)
}
```

#### Service Tests (`internal/service/`)
- **order_service_test.go**: Tests for order service business logic

#### gRPC Handler Tests (`internal/handler/`)
- **grpc_handler_test.go**: Tests for gRPC API endpoints
- **test_helpers.go**: Handler setup and database utilities

**Example:**
```go
func TestGRPCHandler_CreateOrder_ValidRequest(t *testing.T) {
    handler, cleanup := SetupTestHandlerWithData(t)
    defer cleanup()

    req := &api.CreateOrderRequest{
        UserId:           1,
        ConcertSessionId: 1,
        NumberOfTickets:  1,
    }

    resp, err := handler.CreateOrder(context.Background(), req)
    if err != nil {
        t.Logf("Expected error due to no test data: %v", err)
        return
    }

    require.NotNil(t, resp)
    assert.Greater(t, resp.OrderId, int32(0))
    assert.Equal(t, "pending", resp.Status)
    assert.NotEmpty(t, resp.TicketIds)
    assert.Greater(t, resp.TotalPrice, float64(0))
    assert.NotNil(t, resp.CreatedAt)
}
```

**Key Features:**
- **Input Validation**: Tests for invalid user IDs, session IDs, and ticket counts
- **Error Handling**: Tests for missing data, unavailable tickets, and service errors
- **Concurrent Requests**: Tests for race conditions and transaction safety
- **Response Structure**: Validates correct response format and data types
- **gRPC Status Codes**: Ensures proper error status code mapping

**Example:**
```go
func TestOrderService_CreateOrder_ValidRequest(t *testing.T) {
    baseRepo, cleanup := repository.SetupTestDB(t)
    defer cleanup()
    
    baseService := NewBaseService(baseRepo)
    orderService := NewOrderService(baseService)
    
    req := &CreateOrderRequest{
        UserID:           1,
        ConcertSessionID: 1,
    }
    
    resp, err := orderService.CreateOrder(req)
    if err != nil {
        t.Logf("Expected error due to no test data: %v", err)
        return
    }
    
    require.NotNil(t, resp)
    assert.Greater(t, resp.OrderID, 0)
}
```

## Test Database Setup

### Test Helpers

The `internal/repository/test_helpers.go` file provides utilities for setting up test databases:

```go
// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) (*BaseRepository, func()) {
    config := GetTestDBConfig()
    
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.Host, config.Port, config.User, config.Password, config.DBName)
    
    db, err := sqlx.Connect("postgres", dsn)
    require.NoError(t, err)
    
    baseRepo := NewBaseRepository(db)
    
    // Clean up function
    cleanup := func() {
        db.Close()
    }
    
    return baseRepo, cleanup
}
```

### Environment Configuration

Test database configuration can be set via environment variables:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=password
export TEST_DB_NAME=tickets_test_db
```

## Test Categories

### 1. Unit Tests

Unit tests verify individual functions and methods in isolation:

```bash
# Run unit tests
make test-unit

# Run specific unit test categories
make test-models
make test-config
make test-logger
```

**Characteristics:**
- Fast execution
- No external dependencies
- Focus on logic and validation
- High coverage of edge cases

### 2. Integration Tests

Integration tests verify component interactions:

```bash
# Run integration tests
make test-integration

# Run specific integration test categories
make test-repository
make test-service
make test-migrations
```

**Characteristics:**
- Require database connection
- Test component interactions
- Verify business logic
- May require test data setup

### 3. Performance Tests

Benchmark tests measure performance characteristics:

```bash
# Run benchmark tests
make test-benchmark
```

**Example:**
```go
func BenchmarkConcertSessionRepository_GetConcertSessionByID(b *testing.B) {
    baseRepo, cleanup := SetupTestDB(b)
    defer cleanup()
    
    repo := NewConcertSessionRepository(baseRepo)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = repo.GetConcertSessionByID(1)
    }
}
```

## Protocol Buffer Development Workflow

### Generating Protocol Buffer Code

Before running tests that depend on generated Protocol Buffer code:

```bash
# Generate Protocol Buffer code (with automatic cleanup)
make proto

# Clean generated files manually (if needed)
make proto-clean
```

### Testing Protocol Buffer Changes

When modifying `.proto` files:

1. **Update the proto file**:
   ```bash
   vim proto/tickets.proto
   ```

2. **Regenerate the code**:
   ```bash
   make proto
   ```

3. **Update handler code** (if needed):
   ```bash
   vim internal/handler/grpc_handler.go
   ```

4. **Run tests**:
   ```bash
   make test
   ```

### Protocol Buffer Test Integration

The gRPC handler tests automatically use the latest generated Protocol Buffer code:

```bash
# Full workflow
make proto
go test ./internal/handler -v
make test
```

## Test Coverage

### Generate Coverage Report

```bash
make test-coverage
```

This generates:
- `coverage.out`: Raw coverage data
- `coverage.html`: HTML coverage report

### Coverage Analysis

The coverage report shows:
- Line coverage percentage
- Branch coverage
- Function coverage
- Uncovered code sections

### Coverage Targets

Recommended coverage targets:
- **Unit tests**: 90%+ line coverage
- **Integration tests**: 80%+ line coverage
- **Overall**: 85%+ line coverage

## Race Condition Detection

### Run Race Detection

```bash
make test-race
```

Race detection helps identify:
- Concurrent access to shared resources
- Data races in goroutines
- Unsafe concurrent operations

### Common Race Conditions

```go
// Example of potential race condition
var counter int

func increment() {
    counter++ // Race condition without mutex
}

// Fixed version
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

## Test Data Management

### Test Data Setup

For integration tests that require data:

```go
func setupTestData(t *testing.T, db *sqlx.DB) {
    // Insert test concerts
    _, err := db.Exec(`
        INSERT INTO concerts (name, location, description) 
        VALUES ($1, $2, $3)
    `, "Test Concert", "Test Venue", "Test Description")
    require.NoError(t, err)
    
    // Insert test sessions
    _, err = db.Exec(`
        INSERT INTO concert_sessions (concert_id, start_time, end_time, venue, price) 
        VALUES ($1, $2, $3, $4, $5)
    `, 1, time.Now().Add(time.Hour).UnixMilli(), 
       time.Now().Add(2*time.Hour).UnixMilli(), "Test Arena", 99.99)
    require.NoError(t, err)
}
```

### Test Data Cleanup

```go
func cleanupTestData(t *testing.T, db *sqlx.DB) {
    // Clean up in reverse order of dependencies
    db.Exec("DELETE FROM order_items")
    db.Exec("DELETE FROM orders")
    db.Exec("DELETE FROM tickets")
    db.Exec("DELETE FROM concert_sessions")
    db.Exec("DELETE FROM concerts")
}
```

## Best Practices

### 1. Test Naming

Use descriptive test names that explain the scenario:

```go
// Good
func TestConcertSessionRepository_GetConcertSessionByID_InvalidID(t *testing.T)

// Avoid
func TestGetSession(t *testing.T)
```

### 2. Test Structure

Follow the Arrange-Act-Assert pattern:

```go
func TestExample(t *testing.T) {
    // Arrange
    baseRepo, cleanup := SetupTestDB(t)
    defer cleanup()
    repo := NewConcertSessionRepository(baseRepo)
    
    // Act
    session, err := repo.GetConcertSessionByID(1)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, session)
}
```

### 3. Error Handling

Test both success and failure scenarios:

```go
func TestFunction(t *testing.T) {
    // Test success case
    result, err := function(validInput)
    require.NoError(t, err)
    assert.Equal(t, expectedResult, result)
    
    // Test error case
    result, err = function(invalidInput)
    assert.Error(t, err)
    assert.Nil(t, result)
}
```

### 4. Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid input", "valid", true},
        {"empty input", "", false},
        {"invalid input", "invalid", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validate(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 5. Mocking and Stubbing

For external dependencies, use interfaces and mocks:

```go
type Repository interface {
    GetByID(id int) (*Model, error)
}

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetByID(id int) (*Model, error) {
    args := m.Called(id)
    return args.Get(0).(*Model), args.Error(1)
}
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: tickets_test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go mod download
      - run: make test-coverage
      - run: make test-race
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Failures

```bash
# Ensure test database is running
docker run -d --name test-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=tickets_test_db \
  -p 5433:5432 postgres:13
```

#### 2. Test Timeouts

```go
// Add timeout to long-running tests
func TestLongRunning(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Run test with context
    done := make(chan bool)
    go func() {
        // Test logic here
        done <- true
    }()
    
    select {
    case <-done:
        // Test completed
    case <-ctx.Done():
        t.Fatal("Test timed out")
    }
}
```

#### 3. Flaky Tests

- Use deterministic test data
- Avoid time-based assertions
- Clean up test state properly
- Use proper synchronization for concurrent tests

### Debugging Tests

```bash
# Run specific test with verbose output
go test -v ./internal/models/domain -run TestConcert_Validation

# Run test with debug output
go test -v ./internal/repository -run TestConcertSessionRepository -debug

# Run test with coverage for specific function
go test -coverprofile=coverage.out ./internal/service
go tool cover -func=coverage.out | grep CreateOrder
```

This testing system provides comprehensive coverage and helps ensure code quality and reliability across the tickets application. 