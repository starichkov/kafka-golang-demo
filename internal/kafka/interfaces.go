// Package kafka provides abstractions and implementations for Apache Kafka producers and consumers.
// It includes interfaces for dependency injection and factory patterns to enable testability.
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

// KafkaProducerInterface abstracts the kafka.Producer functionality for dependency injection.
// This interface allows for easy mocking in unit tests and provides a clean abstraction
// over the Confluent Kafka Go client's producer operations.
type KafkaProducerInterface interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	Close()
}

// KafkaConsumerInterface abstracts the kafka.Consumer functionality for dependency injection.
// This interface allows for easy mocking in unit tests and provides a clean abstraction
// over the Confluent Kafka Go client's consumer operations.
type KafkaConsumerInterface interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Events() chan kafka.Event
	Close() error
}

// ProducerFactoryInterface abstracts the creation of Kafka producers.
// This factory pattern enables dependency injection and makes it easy to switch
// between real Kafka producers and mocks for testing.
type ProducerFactoryInterface interface {
	NewProducer(configMap *kafka.ConfigMap) (KafkaProducerInterface, error)
}

// ConsumerFactoryInterface abstracts the creation of Kafka consumers.
// This factory pattern enables dependency injection and makes it easy to switch
// between real Kafka consumers and mocks for testing.
type ConsumerFactoryInterface interface {
	NewConsumer(configMap *kafka.ConfigMap) (KafkaConsumerInterface, error)
}

// DefaultProducerFactory implements ProducerFactoryInterface using real Kafka producers.
// This is the production implementation that creates actual Kafka producer instances.
type DefaultProducerFactory struct{}

// NewProducer creates a new Kafka producer with the provided configuration.
// It returns a KafkaProducerInterface that wraps the actual Kafka producer.
func (f *DefaultProducerFactory) NewProducer(configMap *kafka.ConfigMap) (KafkaProducerInterface, error) {
	return kafka.NewProducer(configMap)
}

// DefaultConsumerFactory implements ConsumerFactoryInterface using real Kafka consumers.
// This is the production implementation that creates actual Kafka consumer instances.
type DefaultConsumerFactory struct{}

// NewConsumer creates a new Kafka consumer with the provided configuration.
// It returns a KafkaConsumerInterface that wraps the actual Kafka consumer.
func (f *DefaultConsumerFactory) NewConsumer(configMap *kafka.ConfigMap) (KafkaConsumerInterface, error) {
	return kafka.NewConsumer(configMap)
}

// Global factory instances that can be replaced with mocks in tests.
// These variables provide the default implementations but can be overridden
// for testing purposes to inject mock implementations.
var (
	ProducerFactory ProducerFactoryInterface = &DefaultProducerFactory{}
	ConsumerFactory ConsumerFactoryInterface = &DefaultConsumerFactory{}
)
