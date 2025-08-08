// Package kafka provides Kafka consumer functionality with structured logging and context support.
package kafka

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka-golang-demo/internal/logging"
	"kafka-golang-demo/internal/metrics"
)

// Consumer wraps a Kafka consumer with additional functionality like structured logging,
// context-based cancellation, and optional message channel forwarding. It provides a
// clean API for consuming messages from a specific topic with proper error handling.
type Consumer struct {
	c     KafkaConsumerInterface
	topic string
}

// NewConsumer creates a new Consumer instance configured to consume messages from the specified topic.
// It initializes the underlying Kafka consumer with the provided configuration and subscribes
// to the topic. The consumer is ready to start consuming messages immediately after creation.
//
// Parameters:
//   - brokers: Comma-separated list of Kafka broker addresses (e.g., "localhost:9092")
//   - groupID: The consumer group ID for coordinated consumption
//   - topic: The Kafka topic to consume messages from
//
// Returns a configured Consumer instance or an error if initialization or subscription fails.
func NewConsumer(brokers, groupID, topic string) (*Consumer, error) {
	c, err := ConsumerFactory.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":      brokers,
		"group.id":               groupID,
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 5000, // Enable statistics collection every 5 seconds
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	// Start background statistics processor
	go func() {
		for e := range c.Events() {
			if statsEvent, ok := e.(*kafka.Stats); ok {
				// Process librdkafka statistics for Prometheus metrics
				if err := metrics.ProcessLibrdkafkaStats(statsEvent.String()); err != nil {
					logging.Logger.Error("Failed to process consumer Kafka statistics", "error", err)
				}
			}
		}
	}()

	return &Consumer{c: c, topic: topic}, nil
}

// Run starts the consumer loop that continuously polls for messages until the context is cancelled.
// Messages are logged but not forwarded to any channel. This is a convenience method that calls
// RunWithChannel with a nil channel.
//
// The consumer will respect context cancellation and exit gracefully when ctx.Done() is signaled.
func (co *Consumer) Run(ctx context.Context) {
	co.RunWithChannel(ctx, nil)
}

// RunWithChannel starts the consumer loop that continuously polls for messages until the context is cancelled.
// If msgCh is provided, received messages are forwarded to the channel in addition to being logged.
// This method provides flexibility for applications that need to process messages programmatically.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - msgCh: Optional channel to forward received messages (can be nil)
//
// The consumer handles timeouts gracefully and logs errors appropriately. It will exit when
// the context is cancelled.
func (co *Consumer) RunWithChannel(ctx context.Context, msgCh chan<- *kafka.Message) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			msg, err := co.c.ReadMessage(500 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					continue
				}
				logging.Logger.Error("Consumer error", "error", err)
				continue
			}
			logging.Logger.Info("Message received",
				"message", string(msg.Value),
				"topic", *msg.TopicPartition.Topic,
				"partition", msg.TopicPartition.Partition,
				"offset", msg.TopicPartition.Offset)
			if msgCh != nil {
				msgCh <- msg
			}
		}
	}
}

// Close gracefully shuts down the consumer and cleans up resources.
// Any errors during closure are silently ignored to prevent issues during shutdown.
func (co *Consumer) Close() {
	err := co.c.Close()
	if err != nil {
		return
	}
}
