package logger

import (
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	// Logger is the global logger instance
	Logger *logrus.Logger
)

// LogLevel represents the logging level
type LogLevel string

const (
	// LogLevelDebug represents debug level logging
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo represents info level logging
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn represents warn level logging
	LogLevelWarn LogLevel = "warn"
	// LogLevelError represents error level logging
	LogLevelError LogLevel = "error"
	// LogLevelFatal represents fatal level logging
	LogLevelFatal LogLevel = "fatal"
	// LogLevelPanic represents panic level logging
	LogLevelPanic LogLevel = "panic"
)

// LogFormat represents the logging format
type LogFormat string

const (
	// LogFormatJSON represents JSON format logging
	LogFormatJSON LogFormat = "json"
	// LogFormatText represents text format logging
	LogFormatText LogFormat = "text"
)

// Config holds the logger configuration
type Config struct {
	Level  LogLevel  `json:"level" yaml:"level"`
	Format LogFormat `json:"format" yaml:"format"`
	// Output specifies the output destination (stdout, stderr, or file path)
	Output string `json:"output" yaml:"output"`
	// IncludeCaller adds file and line information to log entries
	IncludeCaller bool `json:"include_caller" yaml:"include_caller"`
	// IncludeTimestamp adds timestamp to log entries
	IncludeTimestamp bool `json:"include_timestamp" yaml:"include_timestamp"`
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:            LogLevelInfo,
		Format:           LogFormatText,
		Output:           "stdout",
		IncludeCaller:    true,
		IncludeTimestamp: true,
	}
}

// Init initializes the global logger with the given configuration
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	Logger = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		return err
	}
	Logger.SetLevel(level)

	// Set log format
	switch config.Format {
	case LogFormatJSON:
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := f.File
				if idx := strings.LastIndex(filename, "/"); idx != -1 {
					filename = filename[idx+1:]
				}
				return "", filename + ":" + string(rune(f.Line))
			},
		})
	case LogFormatText:
		fallthrough
	default:
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   config.IncludeTimestamp,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := f.File
				if idx := strings.LastIndex(filename, "/"); idx != -1 {
					filename = filename[idx+1:]
				}
				return "", filename + ":" + string(rune(f.Line))
			},
		})
	}

	// Set output
	switch config.Output {
	case "stdout":
		Logger.SetOutput(os.Stdout)
	case "stderr":
		Logger.SetOutput(os.Stderr)
	default:
		// Try to open file
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Logger.SetOutput(file)
	}

	// Set caller reporting
	if config.IncludeCaller {
		Logger.SetReportCaller(true)
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		// Initialize with default config if not already initialized
		if err := Init(DefaultConfig()); err != nil {
			// If initialization fails, create a basic logger
			Logger = logrus.New()
			Logger.SetLevel(logrus.InfoLevel)
		}
	}
	return Logger
}

// WithField adds a field to the logger
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields adds multiple fields to the logger
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithError adds an error field to the logger
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Panic logs a panic message and panics
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Panicf logs a formatted panic message and panics
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}
