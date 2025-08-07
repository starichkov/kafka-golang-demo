// Package logging provides structured JSON logging using Go's standard slog package.
// It exposes a global Logger instance that outputs JSON-formatted logs to stdout.
package logging

import (
	"log/slog"
	"os"
)

// Logger is the global structured logger instance using slog with JSON output.
var Logger *slog.Logger

// init initializes the global Logger with JSON output to stdout
func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
