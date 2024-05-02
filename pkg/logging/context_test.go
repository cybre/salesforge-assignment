package logging_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/cybre/salesforge-assignment/pkg/logging"
)

func TestWithLogger(t *testing.T) {
	logger := &slog.Logger{} // Create a mock logger

	// Create a context with the logger
	ctx := logging.WithLogger(context.Background(), logger)

	// Retrieve the logger from the context
	retrievedLogger := ctx.Value("loggerKey").(*slog.Logger)

	// Check if the retrieved logger is the same as the original logger
	if retrievedLogger != logger {
		t.Errorf("Expected logger to be %v, but got %v", logger, retrievedLogger)
	}
}

func TestFromContext(t *testing.T) {
	logger := &slog.Logger{} // Create a mock logger

	// Create a context with the logger
	ctx := context.WithValue(context.Background(), "loggerKey", logger)

	// Retrieve the logger from the context using FromContext
	retrievedLogger := logging.FromContext(ctx)

	// Check if the retrieved logger is the same as the original logger
	if retrievedLogger != logger {
		t.Errorf("Expected logger to be %v, but got %v", logger, retrievedLogger)
	}
}

func TestFromContext_DefaultLogger(t *testing.T) {
	// Create a context without a logger
	ctx := context.Background()

	// Retrieve the logger from the context using FromContext
	retrievedLogger := logging.FromContext(ctx)

	// Check if the retrieved logger is the default logger
	if retrievedLogger != slog.Default() {
		t.Errorf("Expected logger to be the default logger, but got %v", retrievedLogger)
	}
}
