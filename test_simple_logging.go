package main

import (
	"fmt"

	"kafka-golang-demo/internal/logging"
)

func main() {
	fmt.Println("Testing simple structured logging...")

	// Test basic logging
	logging.Logger.Info("This is a structured log message",
		"service", "test",
		"operation", "demo",
		"count", 123,
	)

	logging.Logger.Error("This is an error message",
		"error", "something went wrong",
		"code", 500,
	)

	fmt.Println("✅ Simple structured logging test completed!")
}
