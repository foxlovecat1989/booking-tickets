package logger

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != LogLevelInfo {
		t.Errorf("Expected default level to be info, got %s", config.Level)
	}

	if config.Format != LogFormatText {
		t.Errorf("Expected default format to be text, got %s", config.Format)
	}

	if config.Output != "stdout" {
		t.Errorf("Expected default output to be stdout, got %s", config.Output)
	}

	if !config.IncludeCaller {
		t.Error("Expected default include_caller to be true")
	}

	if !config.IncludeTimestamp {
		t.Error("Expected default include_timestamp to be true")
	}
}

func TestInit(t *testing.T) {
	// Test with nil config (should use default)
	err := Init(nil)
	if err != nil {
		t.Errorf("Failed to initialize logger with nil config: %v", err)
	}

	if Logger == nil {
		t.Error("Logger should not be nil after initialization")
	}

	// Test with custom config
	config := &Config{
		Level:            LogLevelDebug,
		Format:           LogFormatJSON,
		Output:           "stdout",
		IncludeCaller:    false,
		IncludeTimestamp: true,
	}

	err = Init(config)
	if err != nil {
		t.Errorf("Failed to initialize logger with custom config: %v", err)
	}

	if Logger == nil {
		t.Error("Logger should not be nil after initialization")
	}
}

func TestInitWithInvalidLevel(t *testing.T) {
	config := &Config{
		Level: "invalid_level",
	}

	err := Init(config)
	if err == nil {
		t.Error("Expected error when using invalid log level")
	}
}

func TestInitWithFileOutput(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	config := &Config{
		Level:  LogLevelInfo,
		Format: LogFormatText,
		Output: tempFile.Name(),
	}

	err = Init(config)
	if err != nil {
		t.Errorf("Failed to initialize logger with file output: %v", err)
	}

	if Logger == nil {
		t.Error("Logger should not be nil after initialization")
	}
}

func TestGetLogger(t *testing.T) {
	// Reset logger to nil
	Logger = nil

	// Get logger should initialize with default config if not already initialized
	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger should return a logger instance")
	}

	if Logger == nil {
		t.Error("Global logger should be set after GetLogger")
	}
}

func TestLoggingFunctions(t *testing.T) {
	// Initialize logger for testing
	err := Init(&Config{
		Level:  LogLevelDebug,
		Format: LogFormatText,
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Test basic logging functions
	Debug("debug message")
	Debugf("debug message with format: %s", "test")

	Info("info message")
	Infof("info message with format: %s", "test")

	Warn("warning message")
	Warnf("warning message with format: %s", "test")

	Error("error message")
	Errorf("error message with format: %s", "test")
}

func TestWithField(t *testing.T) {
	err := Init(&Config{
		Level:  LogLevelInfo,
		Format: LogFormatText,
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	entry := WithField("key", "value")
	if entry == nil {
		t.Error("WithField should return a logger entry")
	}

	entry.Info("message with field")
}

func TestWithFields(t *testing.T) {
	err := Init(&Config{
		Level:  LogLevelInfo,
		Format: LogFormatText,
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	entry := WithFields(fields)
	if entry == nil {
		t.Error("WithFields should return a logger entry")
	}

	entry.Info("message with fields")
}

func TestWithError(t *testing.T) {
	err := Init(&Config{
		Level:  LogLevelInfo,
		Format: LogFormatText,
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	testError := &os.PathError{Op: "test", Path: "test", Err: os.ErrNotExist}
	entry := WithError(testError)
	if entry == nil {
		t.Error("WithError should return a logger entry")
	}

	entry.Error("message with error")
}
