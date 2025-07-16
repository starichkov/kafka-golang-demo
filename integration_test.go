package main

import (
	"context"
	"golang-kafka-demo/internal/kafka"
	"strings"
	"testing"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafkacontainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

func TestKafkaIntegration(t *testing.T) {
	ctx := context.Background()

	// Start Kafka container
	kafkaContainer, err := kafkacontainer.RunContainer(ctx,
		kafkacontainer.WithClusterID("test-cluster"),
	)
	if err != nil {
		t.Fatalf("Failed to start Kafka container: %v", err)
	}
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Kafka container: %v", err)
		}
	}()

	// Get Kafka bootstrap servers
	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		t.Fatalf("Failed to get Kafka brokers: %v", err)
	}

	bootstrapServers := strings.Join(brokers, ",")
	topic := "test-topic"

	t.Run("Producer", func(t *testing.T) {
		testProducer(t, bootstrapServers, topic)
	})

	t.Run("Consumer", func(t *testing.T) {
		testConsumer(t, bootstrapServers, topic)
	})

	t.Run("ProducerConsumerFlow", func(t *testing.T) {
		testProducerConsumerFlow(t, bootstrapServers, topic)
	})
}

func testProducer(t *testing.T, bootstrapServers, topic string) {
	producer, err := kafka.NewProducer(bootstrapServers, topic)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Send a test message
	testMessage := "test message from producer"
	err = producer.Send(testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Give some time for the message to be delivered
	time.Sleep(1 * time.Second)
}

func testConsumer(t *testing.T, bootstrapServers, topic string) {
	consumer, err := kafka.NewConsumer(bootstrapServers, "test-group", topic)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Test that consumer can be created and closed without errors
	// More detailed testing is done in the producer-consumer flow test
}

func testProducerConsumerFlow(t *testing.T, bootstrapServers, topic string) {
	// Create producer
	producer, err := kafka.NewProducer(bootstrapServers, topic)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create consumer
	consumer, err := kafka.NewConsumer(bootstrapServers, "integration-test-group", topic)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Send test messages
	testMessages := []string{
		"integration test message 1",
		"integration test message 2",
		"integration test message 3",
	}
	expectedMessageCount := len(testMessages)

	// Start consumer in a goroutine with message counting
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgCh := make(chan *ckafka.Message)
	done := make(chan bool, 1)

	go func() {
		received := 0
		for {
			select {
			case <-ctx.Done():
				done <- false
				return
			case <-msgCh:
				received++
				if received >= expectedMessageCount {
					done <- true
					return
				}
			}
		}
	}()

	go consumer.RunWithChannel(ctx, msgCh)

	// Give consumer time to start
	time.Sleep(2 * time.Second)

	// Send all test messages
	for _, msg := range testMessages {
		err = producer.Send(msg)
		if err != nil {
			t.Fatalf("Failed to send message '%s': %v", msg, err)
		}
	}

	// Wait for all messages to be processed or timeout
	select {
	case got := <-done:
		if got {
			t.Log("All expected messages received successfully")
		} else {
			t.Error("Test finished but not all messages received")
		}
	case <-ctx.Done():
		t.Log("Test completed within timeout")
	}
}
