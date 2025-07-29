package main

import (
	"errors"
	"tickets/internal/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger with default configuration
	if err := logger.Init(logger.DefaultConfig()); err != nil {
		panic(err)
	}

	// Example logging usage
	logger.Info("Application started")
	logger.Debug("Debug information")
	logger.Warn("Warning message")
	logger.Error("Error occurred")

	// Logging with fields
	logger.WithField("user_id", 123).Info("User logged in")
	logger.WithFields(logrus.Fields{
		"request_id": "abc123",
		"method":     "GET",
		"path":       "/api/tickets",
	}).Info("API request received")

	// Logging with error
	err := errors.New("database connection failed")
	logger.WithError(err).Error("Database connection failed")
}
