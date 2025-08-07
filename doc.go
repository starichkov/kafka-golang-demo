// Package kafka-golang-demo is a comprehensive demonstration of Apache Kafka integration in Go.
//
// # Overview
//
// This project showcases best practices for building Kafka producers and consumers in Go,
// including dependency injection, structured logging, testing patterns, and Docker deployment.
// It's part of a multi-language Kafka demo series and serves as both a learning resource
// and a practical implementation reference.
//
// # Architecture
//
// The project follows a clean architecture pattern with the following key components:
//
//   - internal/kafka: Core Kafka abstractions and implementations
//   - internal/logging: Structured JSON logging using Go's slog package
//   - cli/producer: Command-line Kafka producer application
//   - cli/consumer: Command-line Kafka consumer application
//
// # Key Features
//
//   - Interface-based design for easy testing and dependency injection
//   - Structured JSON logging with contextual information
//   - Graceful shutdown and error handling
//   - Docker and Docker Compose support
//   - Comprehensive unit and integration tests
//   - Factory pattern for creating Kafka clients
//
// # Usage Examples
//
// ## Producer Example
//
//	producer, err := kafka.NewProducer("localhost:9092", "my-topic")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer producer.Close()
//
//	err = producer.Send("Hello, Kafka!")
//	if err != nil {
//	    log.Printf("Failed to send message: %v", err)
//	}
//
// ## Consumer Example
//
//	consumer, err := kafka.NewConsumer("localhost:9092", "my-group", "my-topic")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer consumer.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	consumer.Run(ctx)
//
// # Testing
//
// The project includes comprehensive test coverage with both unit tests and integration tests:
//
//   - Unit tests use mocks to isolate components
//   - Integration tests use Testcontainers to run real Kafka instances
//   - All tests follow Go testing conventions and can be run with go test
//
// # Documentation Generation
//
// This project uses Go's built-in documentation tools. To generate and view documentation:
//
//	# View package documentation
//	go doc ./internal/kafka
//
//	# View specific type documentation
//	go doc ./internal/kafka.Producer
//
//	# View specific method documentation
//	go doc ./internal/kafka.Producer.Send
//
// For more detailed information, see the individual package documentation and README.md.
package main
