# Logging with Logrus

This document provides a comprehensive guide to the structured logging system implemented using Logrus for the tickets application.

## Overview

The logging system provides:
- **Structured logging** with JSON and text formats
- **Configurable log levels** (debug, info, warn, error, fatal, panic)
- **Multiple output destinations** (stdout, stderr, files)
- **Caller information** (file and line numbers)
- **Timestamp formatting** with ISO 8601 format
- **Field-based logging** for structured data
- **Error context** with automatic error field handling

## Quick Start

### 1. Basic Usage

```go
import "tickets/internal/logger"

// Initialize logger (usually done in main)
if err := logger.Init(&cfg.Logging); err != nil {
    logger.Fatalf("Failed to initialize logger: %v", err)
}

// Basic logging
logger.Info("Server started")
logger.Infof("Listening on port %d", 8080)
logger.Error("Database connection failed")
logger.Debug("Processing request")
```

### 2. Structured Logging

```go
// With fields
logger.WithField("user_id", 123).Info("User logged in")
logger.WithField("request_id", "abc-123").WithField("method", "GET").Info("Request processed")

// With multiple fields
fields := map[string]interface{}{
    "user_id": 123,
    "action": "login",
    "ip": "192.168.1.1",
}
logger.WithFields(fields).Info("User action")

// With errors
if err != nil {
    logger.WithError(err).Error("Database query failed")
}
```

## Configuration

### Configuration Structure

```yaml
logging:
  level: "info"              # debug, info, warn, error, fatal, panic
  format: "text"             # text, json
  output: "stdout"           # stdout, stderr, or file path
  include_caller: true       # Include file and line information
  include_timestamp: true    # Include timestamps
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `level` | string | `"info"` | Logging level (debug, info, warn, error, fatal, panic) |
| `format` | string | `"text"` | Output format (text, json) |
| `output` | string | `"stdout"` | Output destination (stdout, stderr, file path) |
| `include_caller` | bool | `true` | Include file and line information |
| `include_timestamp` | bool | `true` | Include timestamps in log entries |

### Environment-Specific Configurations

#### Development
```yaml
logging:
  level: "debug"
  format: "text"
  output: "stdout"
  include_caller: true
  include_timestamp: true
```

#### Production
```yaml
logging:
  level: "info"
  format: "json"
  output: "/var/log/tickets/app.log"
  include_caller: false
  include_timestamp: true
```

#### Testing
```yaml
logging:
  level: "error"
  format: "text"
  output: "stderr"
  include_caller: false
  include_timestamp: false
```

## Log Levels

### Debug Level
```go
logger.Debug("Processing request details")
logger.Debugf("Request body: %+v", requestBody)
```
Use for detailed debugging information.

### Info Level
```go
logger.Info("Server started successfully")
logger.Infof("User %s logged in", username)
```
Use for general application flow information.

### Warn Level
```go
logger.Warn("Database connection slow")
logger.Warnf("Rate limit exceeded for user %s", userID)
```
Use for potentially harmful situations.

### Error Level
```go
logger.Error("Failed to process payment")
logger.Errorf("Database query failed: %v", err)
logger.WithError(err).Error("API request failed")
```
Use for error conditions that don't stop the application.

### Fatal Level
```go
logger.Fatal("Cannot bind to port 8080")
logger.Fatalf("Configuration error: %v", err)
```
Use for errors that require application shutdown.

### Panic Level
```go
logger.Panic("Critical system failure")
logger.Panicf("Unrecoverable error: %v", err)
```
Use for unrecoverable errors that cause panic.

## Structured Logging Examples

### Request Logging
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    requestID := generateRequestID()
    
    logger.WithFields(logrus.Fields{
        "request_id": requestID,
        "method":     r.Method,
        "path":       r.URL.Path,
        "user_agent": r.UserAgent(),
        "ip":         r.RemoteAddr,
    }).Info("Request received")
    
    // Process request...
    
    logger.WithField("request_id", requestID).Info("Request completed")
}
```

### Database Operations
```go
func createUser(user *User) error {
    logger.WithField("email", user.Email).Info("Creating new user")
    
    if err := db.Create(user).Error; err != nil {
        logger.WithError(err).WithField("email", user.Email).Error("Failed to create user")
        return err
    }
    
    logger.WithFields(logrus.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    }).Info("User created successfully")
    
    return nil
}
```

### Error Handling
```go
func processPayment(payment *Payment) error {
    logger.WithFields(logrus.Fields{
        "payment_id": payment.ID,
        "amount":     payment.Amount,
        "currency":   payment.Currency,
    }).Info("Processing payment")
    
    if err := paymentGateway.Process(payment); err != nil {
        logger.WithError(err).WithFields(logrus.Fields{
            "payment_id": payment.ID,
            "amount":     payment.Amount,
        }).Error("Payment processing failed")
        return err
    }
    
    logger.WithField("payment_id", payment.ID).Info("Payment processed successfully")
    return nil
}
```

## Output Formats

### Text Format
```
time="2024-01-15T10:30:45.123Z" level=info msg="Server started" port=8080 caller=main.go:25
time="2024-01-15T10:30:45.124Z" level=info msg="Database connected" host=localhost port=5432 caller=main.go:30
```

### JSON Format
```json
{
  "caller": "main.go:25",
  "level": "info",
  "msg": "Server started",
  "port": 8080,
  "time": "2024-01-15T10:30:45.123Z"
}
{
  "caller": "main.go:30",
  "level": "info",
  "msg": "Database connected",
  "host": "localhost",
  "port": 5432,
  "time": "2024-01-15T10:30:45.124Z"
}
```

## Best Practices

### 1. Log Level Selection
- **Debug**: Detailed debugging information
- **Info**: General application flow
- **Warn**: Potentially harmful situations
- **Error**: Error conditions
- **Fatal**: Application shutdown errors
- **Panic**: Unrecoverable errors

### 2. Structured Data
```go
// Good: Structured logging with fields
logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "action": "login",
    "ip": clientIP,
}).Info("User logged in")

// Avoid: String concatenation
logger.Info("User " + user.ID + " logged in from " + clientIP)
```

### 3. Error Context
```go
// Good: Include error context
if err != nil {
    logger.WithError(err).WithField("user_id", userID).Error("Failed to update user")
}

// Avoid: Just logging the error
if err != nil {
    logger.Error(err.Error())
}
```

### 4. Sensitive Data
```go
// Good: Mask sensitive data
logger.WithField("email", maskEmail(user.Email)).Info("User registered")

// Avoid: Logging sensitive data directly
logger.WithField("password", user.Password).Info("User registered")
```

### 5. Performance Considerations
```go
// Good: Use structured logging for better performance
logger.WithField("user_id", userID).Info("User action")

// Avoid: String formatting in debug logs
if logger.GetLevel() >= logrus.DebugLevel {
    logger.Debugf("Processing user %s with data %+v", userID, userData)
}
```

## Integration Examples

### HTTP Middleware
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log request
        logger.WithFields(logrus.Fields{
            "method":     r.Method,
            "path":       r.URL.Path,
            "user_agent": r.UserAgent(),
            "ip":         r.RemoteAddr,
        }).Info("Request started")
        
        // Process request
        next.ServeHTTP(w, r)
        
        // Log response
        duration := time.Since(start)
        logger.WithFields(logrus.Fields{
            "method":   r.Method,
            "path":     r.URL.Path,
            "duration": duration.String(),
        }).Info("Request completed")
    })
}
```

### Database Operations
```go
func (r *UserRepository) Create(user *User) error {
    logger.WithField("email", user.Email).Debug("Creating user in database")
    
    if err := r.db.Create(user).Error; err != nil {
        logger.WithError(err).WithField("email", user.Email).Error("Failed to create user")
        return err
    }
    
    logger.WithFields(logrus.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    }).Info("User created successfully")
    
    return nil
}
```

### Background Jobs
```go
func processOrder(orderID string) {
    logger.WithField("order_id", orderID).Info("Starting order processing")
    
    // Process order...
    
    if err := processPayment(orderID); err != nil {
        logger.WithError(err).WithField("order_id", orderID).Error("Order processing failed")
        return
    }
    
    logger.WithField("order_id", orderID).Info("Order processed successfully")
}
```

## Testing

### Unit Tests
```go
func TestUserService(t *testing.T) {
    // Initialize logger for testing
    err := logger.Init(&logger.Config{
        Level:  logger.LogLevelError, // Only log errors during tests
        Format: logger.LogFormatText,
        Output: "stderr",
    })
    if err != nil {
        t.Fatalf("Failed to initialize logger: %v", err)
    }
    
    // Test code...
}
```

### Integration Tests
```go
func TestAPIEndpoint(t *testing.T) {
    // Initialize logger with test configuration
    err := logger.Init(&logger.Config{
        Level:  logger.LogLevelDebug,
        Format: logger.LogFormatJSON,
        Output: "stdout",
    })
    if err != nil {
        t.Fatalf("Failed to initialize logger: %v", err)
    }
    
    // Test API endpoints...
}
```

## Monitoring and Alerting

### Log Aggregation
- Use JSON format in production for better log aggregation
- Include structured fields for filtering and searching
- Use consistent field names across the application

### Metrics from Logs
```go
// Log metrics as structured data
logger.WithFields(logrus.Fields{
    "metric": "request_duration",
    "value":  duration.Milliseconds(),
    "unit":   "ms",
}).Info("Request performance")
```

### Error Tracking
```go
// Include error context for better debugging
if err != nil {
    logger.WithError(err).WithFields(logrus.Fields{
        "component": "payment_processor",
        "operation": "process_payment",
        "payment_id": paymentID,
    }).Error("Payment processing error")
}
```

## Troubleshooting

### Common Issues

#### Logger Not Initialized
```go
// Error: Logger not initialized
// Solution: Call logger.Init() before using logger functions
if err := logger.Init(&cfg.Logging); err != nil {
    log.Fatalf("Failed to initialize logger: %v", err)
}
```

#### Invalid Log Level
```go
// Error: Invalid log level
// Solution: Use valid log levels (debug, info, warn, error, fatal, panic)
config := &logger.Config{
    Level: logger.LogLevelInfo, // Valid level
}
```

#### File Output Permission Issues
```go
// Error: Cannot write to log file
// Solution: Ensure proper file permissions and directory exists
config := &logger.Config{
    Output: "/var/log/tickets/app.log", // Ensure directory exists and is writable
}
```

This logging system provides a robust, structured approach to application logging with full integration into your Go tickets application. 