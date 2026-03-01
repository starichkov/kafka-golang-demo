# Development Guidelines for Kafka-Golang Demo

This document provides project-specific information for developers working on the `kafka-golang-demo`.

## 1. Build/Configuration Instructions

### Prerequisites
- **Go 1.25.7** or later.
- **librdkafka**: The `confluent-kafka-go` client is a C wrapper and requires `librdkafka` to be installed on the host system for local builds with `CGO_ENABLED=1`.
- **Docker**: Required for running Kafka via Docker Compose and for integration tests (Testcontainers).

### Build for Local Execution
To build the producer or consumer locally:
```bash
CGO_ENABLED=1 go build -o producer ./cli/producer
CGO_ENABLED=1 go build -o consumer ./cli/consumer
```

### Run with Docker Compose
The easiest way to run the entire stack is:
```bash
docker compose up --build
```
This starts Kafka in KRaft mode and then runs the producer and consumer.

## 2. Testing Information

### Unit Tests
Unit tests use mocks for Kafka dependencies. They are located in `internal/kafka/` and can be run with:
```bash
go test -v ./internal/kafka
```

### Integration Tests
Integration tests use [Testcontainers](https://golang.testcontainers.org/) to spin up a real Kafka instance. **Docker must be running.**
```bash
go test -v integration_test.go
```
Note: Integration tests might take some time as they pull the Kafka image (`confluentinc/cp-kafka:7.9.3`) and wait for the container to start.

### Mocking Strategy
We use a factory pattern to inject mock producers and consumers. See `internal/kafka/interfaces.go` for the interface definitions and `internal/kafka/mocks.go` for the mock implementations.
To use mocks in a test:
1. Create a `MockProducerFactory` or `MockConsumerFactory`.
2. Replace the global `ProducerFactory` or `ConsumerFactory` with your mock.
3. Your code under test should use the global factory to create producers/consumers.

### Demo Test Example
Here is a simple test demonstration using mocks (similar to `internal/kafka/demo_test.go` which was used for verification):
```go
func TestSimpleMockProducer(t *testing.T) {
	// Setup mock factory
	mockFactory := NewMockProducerFactory()
	config := &kafka.ConfigMap{"bootstrap.servers": "localhost"}
	
	// Create producer using factory
	producer, err := mockFactory.NewProducer(config)
	if err != nil {
		t.Fatalf("Failed to create mock producer: %v", err)
	}
	
	mockProducer := producer.(*MockProducer)
	
	// Send a message
	topic := "test-topic"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("test-message"),
	}
	
	_ = mockProducer.Produce(msg, nil)
	
	// Verify
	if len(mockProducer.GetMessages()) != 1 {
		t.Errorf("Expected 1 message, got %d", len(mockProducer.GetMessages()))
	}
}
```

## 3. Additional Development Information

### Code Style
- Follow standard Go formatting (`go fmt`).
- Use structured logging via `log/slog` (available in `internal/logging/simple.go`).
- Errors should be handled explicitly and wrapped with context when appropriate.

### Kafka Implementation Details
- **KRaft Mode**: Kafka is configured to run without Zookeeper.
- **CGO**: `confluent-kafka-go` requires CGO. Ensure `CGO_ENABLED=1` when building or testing locally if you are using the real Kafka client.
- **Factory Pattern**: Always use `kafka.ProducerFactory` and `kafka.ConsumerFactory` to create instances to maintain testability.
