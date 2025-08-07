// Package kafka provides Kafka producer functionality with structured logging.
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka-golang-demo/internal/logging"
)

// Producer wraps a Kafka producer with additional functionality like structured logging
// and simplified message sending. It handles delivery reports in the background and
// provides a clean API for sending messages to a specific topic.
type Producer struct {
	p     KafkaProducerInterface
	topic string
}

// NewProducer creates a new Producer instance configured to send messages to the specified topic.
// It initializes the underlying Kafka producer with the provided broker addresses and sets up
// background delivery report logging. The producer is ready to send messages immediately after creation.
//
// Parameters:
//   - brokers: Comma-separated list of Kafka broker addresses (e.g., "localhost:9092")
//   - topic: The Kafka topic to send messages to
//
// Returns a configured Producer instance or an error if initialization fails.
func NewProducer(brokers, topic string) (*Producer, error) {
	p, err := ProducerFactory.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}

	// Background delivery report logger
	go func() {
		for e := range p.Events() {
			if ev, ok := e.(*kafka.Message); ok {
				if ev.TopicPartition.Error != nil {
					logging.Logger.Error("Delivery failed", "error", ev.TopicPartition.Error, "topic_partition", ev.TopicPartition)
				} else {
					logging.Logger.Info("Message delivered", "topic", *ev.TopicPartition.Topic, "partition", ev.TopicPartition.Partition, "offset", ev.TopicPartition.Offset)
				}
			}
		}
	}()

	return &Producer{p: p, topic: topic}, nil
}

// Send publishes a message to the configured Kafka topic.
// The message is sent asynchronously and delivery reports are handled in the background.
// Any partition assignment is handled automatically by Kafka.
//
// Parameters:
//   - message: The string message to send to Kafka
//
// Returns an error if the message cannot be queued for sending.
func (pr *Producer) Send(message string) error {
	return pr.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &pr.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
}

// Close gracefully shuts down the producer by flushing any pending messages
// and cleaning up resources. It waits up to 3 seconds for pending messages
// to be delivered before closing the connection.
func (pr *Producer) Close() {
	pr.p.Flush(3000)
	pr.p.Close()
}
