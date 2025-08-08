package kafka

import (
	"testing"
	"time"
	
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// setupMockProducer sets up a mock producer factory for testing
func setupMockProducer() (*MockProducerFactory, func()) {
	originalFactory := ProducerFactory
	mockFactory := NewMockProducerFactory()
	ProducerFactory = mockFactory

	return mockFactory, func() {
		ProducerFactory = originalFactory
	}
}

func TestNewProducer(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	tests := []struct {
		name         string
		brokers      string
		topic        string
		factoryError bool
		wantErr      bool
	}{
		{
			name:    "valid brokers and topic",
			brokers: "localhost:9092",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "empty brokers",
			brokers: "",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "empty topic",
			brokers: "localhost:9092",
			topic:   "",
			wantErr: false,
		},
		{
			name:    "invalid broker format",
			brokers: "invalid-broker-format",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:         "factory error",
			brokers:      "localhost:9092",
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

			producer, err := NewProducer(tt.brokers, tt.topic)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewProducer() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("NewProducer() unexpected error: %v", err)
				return
			}

			if producer == nil {
				t.Errorf("NewProducer() returned nil producer")
				return
			}

			if producer.topic != tt.topic {
				t.Errorf("NewProducer() topic = %v, want %v", producer.topic, tt.topic)
			}

			if producer.p == nil {
				t.Errorf("NewProducer() internal producer is nil")
			}

			// Clean up
			producer.Close()
		})
	}
}

func TestProducer_Send(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Get the mock producer to verify messages
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]

	tests := []struct {
		name        string
		message     string
		producerErr bool
		wantErr     bool
	}{
		{
			name:    "valid message",
			message: "test message",
			wantErr: false,
		},
		{
			name:    "empty message",
			message: "",
			wantErr: false,
		},
		{
			name:    "long message",
			message: string(make([]byte, 1000)),
			wantErr: false,
		},
		{
			name:    "unicode message",
			message: "测试消息 🚀",
			wantErr: false,
		},
		{
			name:        "producer error",
			message:     "test message",
			producerErr: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.producerErr {
				mockProducer.SetError(true, "mock producer error")
				defer mockProducer.SetError(false, "")
			}

			initialMsgCount := len(mockProducer.GetMessages())
			err := producer.Send(tt.message)

			if tt.wantErr && err == nil {
				t.Errorf("Producer.Send() expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Producer.Send() unexpected error: %v", err)
			}

			if !tt.wantErr && !tt.producerErr {
				// Verify message was stored in mock
				messages := mockProducer.GetMessages()
				if len(messages) != initialMsgCount+1 {
					t.Errorf("Expected %d messages, got %d", initialMsgCount+1, len(messages))
				} else {
					lastMsg := messages[len(messages)-1]
					if string(lastMsg.Value) != tt.message {
						t.Errorf("Expected message %q, got %q", tt.message, string(lastMsg.Value))
					}
				}
			}
		})
	}
}

func TestProducer_Close(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}

	// Get the mock producer to verify state
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]

	// Send a message before closing
	_ = producer.Send("test message before close")

	// Verify producer is not closed initially
	if mockProducer.IsClosed() {
		t.Errorf("Producer should not be closed initially")
	}

	// Test that Close doesn't panic and completes within reasonable time
	done := make(chan bool, 1)
	go func() {
		producer.Close()
		done <- true
	}()

	select {
	case <-done:
		// Close completed successfully
	case <-time.After(5 * time.Second):
		t.Errorf("Producer.Close() took too long to complete")
		return
	}

	// Verify producer is closed
	if !mockProducer.IsClosed() {
		t.Errorf("Producer should be closed after Close()")
	}

	// Test that calling Close multiple times doesn't panic
	producer.Close()
}

func TestProducer_StatisticsEventHandling(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Get the mock producer
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]


	// Send a statistics event to the producer
	statsEvent := &kafka.Stats{}
	// We can't easily set the internal JSON string in the Stats event,
	// so we'll test that the event processing doesn't panic
	mockProducer.events <- statsEvent

	// Give some time for event processing
	time.Sleep(100 * time.Millisecond)

	// Test should complete without panicking
	// The actual metrics processing is tested in metrics package tests
}

func TestProducer_MessageEventHandling(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Get the mock producer
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]

	// Send a message first
	err = producer.Send("test message")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Give time for background event processing
	time.Sleep(100 * time.Millisecond)

	// Verify the message was processed (stored in mock)
	messages := mockProducer.GetMessages()
	if len(messages) == 0 {
		t.Error("Expected at least one message to be stored")
	}

	// Test that message events are handled without panicking
	topic := "test-topic"
	partition := int32(0)
	offset := kafka.Offset(0)

	// Create a successful delivery message event
	successEvent := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
			Offset:    offset,
		},
		Value: []byte("test message"),
	}

	// Send the event
	mockProducer.events <- successEvent

	// Create a failed delivery message event
	failedEvent := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
			Error:     kafka.NewError(kafka.ErrMsgTimedOut, "Message timed out", false),
		},
		Value: []byte("failed message"),
	}

	// Send the failed event
	mockProducer.events <- failedEvent

	// Give time for event processing
	time.Sleep(100 * time.Millisecond)

	// Test should complete without panicking
	// The actual logging is tested implicitly (no assertion needed for logs)
}

func TestProducer_SendAfterClose(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}

	// Get the mock producer to verify state
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]

	producer.Close()

	// Verify producer is closed
	if !mockProducer.IsClosed() {
		t.Errorf("Producer should be closed after Close()")
	}

	// Sending after close should return an error
	err = producer.Send("message after close")
	if err == nil {
		t.Errorf("Producer.Send() after Close() should return an error")
	}
}

func TestProducer_ConcurrentSend(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Get the mock producer to verify messages
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Fatalf("Expected 1 mock producer, got %d", len(mockProducers))
	}
	mockProducer := mockProducers[0]

	// Test concurrent sends
	done := make(chan bool, 10)
	expectedMessages := 10 * 5 // 10 goroutines * 5 messages each

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 5; j++ {
				err := producer.Send("concurrent message")
				if err != nil {
					t.Errorf("Concurrent send failed: %v", err)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Errorf("Concurrent send test timed out")
			return
		}
	}

	// Verify all messages were sent
	messages := mockProducer.GetMessages()
	if len(messages) != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, len(messages))
	}
}

func TestProducer_StructFields(t *testing.T) {
	mockFactory, cleanup := setupMockProducer()
	defer cleanup()

	producer, err := NewProducer("localhost:9092", "test-topic")
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Test that the producer has the expected fields
	if producer.topic != "test-topic" {
		t.Errorf("Producer.topic = %v, want %v", producer.topic, "test-topic")
	}

	if producer.p == nil {
		t.Errorf("Producer.p should not be nil")
	}

	// Verify mock producer was created
	mockProducers := mockFactory.GetProducers()
	if len(mockProducers) != 1 {
		t.Errorf("Expected 1 mock producer, got %d", len(mockProducers))
	}
}
