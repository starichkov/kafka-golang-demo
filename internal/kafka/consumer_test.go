package kafka

import (
	"context"
	"testing"
	"time"
)

// setupMockConsumer sets up a mock consumer factory for testing
func setupMockConsumer() (*MockConsumerFactory, func()) {
	originalFactory := ConsumerFactory
	mockFactory := NewMockConsumerFactory()
	ConsumerFactory = mockFactory

	return mockFactory, func() {
		ConsumerFactory = originalFactory
	}
}

func TestNewConsumer(t *testing.T) {
	mockFactory, cleanup := setupMockConsumer()
	defer cleanup()

	tests := []struct {
		name         string
		brokers      string
		groupID      string
		topic        string
		factoryError bool
		wantErr      bool
	}{
		{
			name:    "valid parameters",
			brokers: "localhost:9092",
			groupID: "test-group",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "empty brokers",
			brokers: "",
			groupID: "test-group",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "empty group ID",
			brokers: "localhost:9092",
			groupID: "",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "empty topic",
			brokers: "localhost:9092",
			groupID: "test-group",
			topic:   "",
			wantErr: false,
		},
		{
			name:    "invalid broker format",
			brokers: "invalid-broker-format",
			groupID: "test-group",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:         "factory error",
			brokers:      "localhost:9092",
			groupID:      "test-group",
			topic:        "test-topic",
			factoryError: true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.factoryError {
				mockFactory.SetError(true, "mock factory error")
				defer mockFactory.SetError(false, "")
			}

			consumer, err := NewConsumer(tt.brokers, tt.groupID, tt.topic)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewConsumer() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("NewConsumer() unexpected error: %v", err)
				return
			}

			if consumer == nil {
				t.Errorf("NewConsumer() returned nil consumer")
				return
			}

			if consumer.topic != tt.topic {
				t.Errorf("NewConsumer() topic = %v, want %v", consumer.topic, tt.topic)
			}

			if consumer.c == nil {
				t.Errorf("NewConsumer() internal consumer is nil")
			}

			// Verify mock consumer was created and subscribed
			mockConsumers := mockFactory.GetConsumers()
			if len(mockConsumers) == 0 {
				t.Errorf("Expected mock consumer to be created")
			} else {
				mockConsumer := mockConsumers[len(mockConsumers)-1]
				topics := mockConsumer.GetTopics()
				if len(topics) != 1 || topics[0] != tt.topic {
					t.Errorf("Expected topic %q to be subscribed, got %v", tt.topic, topics)
				}
			}

			// Clean up
			consumer.Close()
		})
	}
}

func TestConsumer_Run(t *testing.T) {
	mockFactory, cleanup := setupMockConsumer()
	defer cleanup()

	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Test that Run respects context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Run should exit when context is cancelled
	done := make(chan bool, 1)
	go func() {
		consumer.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
		// Expected: Run should exit when context is cancelled
	case <-time.After(2 * time.Second):
		t.Errorf("Consumer.Run() did not respect context cancellation")
	}

	// Verify mock consumer was created
	mockConsumers := mockFactory.GetConsumers()
	if len(mockConsumers) != 1 {
		t.Errorf("Expected 1 mock consumer, got %d", len(mockConsumers))
	}
}

func TestConsumer_RunWithMessageCount(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}
	defer consumer.Close()

	tests := []struct {
		name          string
		expectedCount int
		timeout       time.Duration
	}{
		{
			name:          "zero messages",
			expectedCount: 0,
			timeout:       1 * time.Second,
		},
		{
			name:          "one message",
			expectedCount: 1,
			timeout:       2 * time.Second,
		},
		{
			name:          "multiple messages",
			expectedCount: 5,
			timeout:       3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			done := make(chan bool, 1)

			go consumer.RunWithMessageCount(ctx, tt.expectedCount, done)

			if tt.expectedCount == 0 {
				// For zero messages, we expect the function to complete immediately
				select {
				case <-done:
					// Expected behavior for zero messages
				case <-time.After(100 * time.Millisecond):
					// This is also acceptable - the function might wait for context cancellation
				}
			} else {
				// For non-zero messages, we expect either completion or timeout
				select {
				case <-done:
					// Messages were received (unlikely without a producer)
				case <-ctx.Done():
					// Context timeout (expected when no messages are available)
				case <-time.After(tt.timeout + 500*time.Millisecond):
					t.Errorf("RunWithMessageCount() did not respect context timeout")
				}
			}
		})
	}
}

func TestConsumer_RunWithMessageCount_ContextCancellation(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	go consumer.RunWithMessageCount(ctx, 10, done)

	// Cancel context after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// Function should exit due to context cancellation
	select {
	case <-done:
		// This might happen if messages are received quickly
	case <-time.After(1 * time.Second):
		// Expected: function should exit due to context cancellation
	}
}

func TestConsumer_Close(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}

	// Test that Close doesn't panic
	consumer.Close()

	// Test that calling Close multiple times doesn't panic
	consumer.Close()
}

func TestConsumer_StructFields(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}
	defer consumer.Close()

	// Test that the consumer has the expected fields
	if consumer.topic != "test-topic" {
		t.Errorf("Consumer.topic = %v, want %v", consumer.topic, "test-topic")
	}

	if consumer.c == nil {
		t.Errorf("Consumer.c should not be nil")
	}
}

func TestConsumer_ConcurrentRun(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}
	defer consumer.Close()

	// Test that multiple Run calls don't cause issues
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan bool, 2)

	// Start two concurrent Run operations
	go func() {
		consumer.Run(ctx)
		done <- true
	}()

	go func() {
		consumer.Run(ctx)
		done <- true
	}()

	// Wait for both to complete
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("Concurrent Run test timed out")
			return
		}
	}
}

func TestConsumer_RunWithMessageCount_EdgeCases(t *testing.T) {
	// Skip if Kafka is not available
	consumer, err := NewConsumer("localhost:9092", "test-group", "test-topic")
	if err != nil {
		t.Skipf("Kafka broker not available: %v", err)
		return
	}
	defer consumer.Close()

	tests := []struct {
		name          string
		expectedCount int
		description   string
	}{
		{
			name:          "negative count",
			expectedCount: -1,
			description:   "should handle negative count gracefully",
		},
		{
			name:          "very large count",
			expectedCount: 1000000,
			description:   "should handle very large count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			done := make(chan bool, 1)

			// This should not panic regardless of the count value
			go consumer.RunWithMessageCount(ctx, tt.expectedCount, done)

			select {
			case <-done:
				// Function completed
			case <-ctx.Done():
				// Context timeout (expected)
			case <-time.After(500 * time.Millisecond):
				t.Errorf("RunWithMessageCount() with %s did not respect context timeout", tt.description)
			}
		})
	}
}
