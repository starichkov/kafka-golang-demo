package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"sync"
	"time"
)

// MockProducer implements KafkaProducerInterface for testing
type MockProducer struct {
	events       chan kafka.Event
	closed       bool
	mutex        sync.RWMutex
	messages     []*kafka.Message
	shouldError  bool
	errorMessage string
}

func NewMockProducer() *MockProducer {
	return &MockProducer{
		events:   make(chan kafka.Event, 100),
		messages: make([]*kafka.Message, 0),
	}
}

func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.closed {
		return fmt.Errorf("producer is closed")
	}

	if m.shouldError {
		return fmt.Errorf("%s", m.errorMessage)
	}

	// Store the message
	m.messages = append(m.messages, msg)

	// Send delivery report
	go func() {
		deliveryMsg := &kafka.Message{
			TopicPartition: msg.TopicPartition,
			Value:          msg.Value,
		}

		if m.shouldError {
			deliveryMsg.TopicPartition.Error = fmt.Errorf("%s", m.errorMessage)
		}

		// Hold the lock while checking closed status and sending to prevent race
		m.mutex.RLock()
		defer m.mutex.RUnlock()

		if !m.closed {
			select {
			case m.events <- deliveryMsg:
			default:
			}
		}
	}()

	return nil
}

func (m *MockProducer) Events() chan kafka.Event {
	return m.events
}

func (m *MockProducer) Flush(timeoutMs int) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.messages)
}

func (m *MockProducer) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.closed {
		m.closed = true
		close(m.events)
	}
}

func (m *MockProducer) SetError(shouldError bool, errorMessage string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldError = shouldError
	m.errorMessage = errorMessage
}

func (m *MockProducer) GetMessages() []*kafka.Message {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result := make([]*kafka.Message, len(m.messages))
	copy(result, m.messages)
	return result
}

func (m *MockProducer) IsClosed() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.closed
}

// MockConsumer implements KafkaConsumerInterface for testing
type MockConsumer struct {
	messages     chan *kafka.Message
	events       chan kafka.Event
	closed       bool
	mutex        sync.RWMutex
	topics       []string
	shouldError  bool
	errorMessage string
}

func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		messages: make(chan *kafka.Message, 100),
		events:   make(chan kafka.Event, 100),
		topics:   make([]string, 0),
	}
}

func (m *MockConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.closed {
		return fmt.Errorf("consumer is closed")
	}

	if m.shouldError {
		return fmt.Errorf("%s", m.errorMessage)
	}

	m.topics = make([]string, len(topics))
	copy(m.topics, topics)
	return nil
}

func (m *MockConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	m.mutex.RLock()
	if m.closed {
		m.mutex.RUnlock()
		return nil, fmt.Errorf("consumer is closed")
	}

	if m.shouldError {
		m.mutex.RUnlock()
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	m.mutex.RUnlock()

	select {
	case msg := <-m.messages:
		return msg, nil
	case <-time.After(timeout):
		// Return timeout error similar to real Kafka consumer
		return nil, kafka.NewError(kafka.ErrTimedOut, "Timed out", false)
	}
}

func (m *MockConsumer) Events() chan kafka.Event {
	return m.events
}

func (m *MockConsumer) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.closed {
		m.closed = true
		close(m.messages)
		close(m.events)
	}
	return nil
}

func (m *MockConsumer) SetError(shouldError bool, errorMessage string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldError = shouldError
	m.errorMessage = errorMessage
}

func (m *MockConsumer) AddMessage(topic string, partition int32, offset kafka.Offset, value []byte) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.closed {
		return
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
			Offset:    offset,
		},
		Value: value,
	}

	select {
	case m.messages <- msg:
	default:
		// Channel is full, drop message
	}
}

func (m *MockConsumer) GetTopics() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result := make([]string, len(m.topics))
	copy(result, m.topics)
	return result
}

func (m *MockConsumer) IsClosed() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.closed
}

// MockProducerFactory implements ProducerFactoryInterface for testing
type MockProducerFactory struct {
	producers    []*MockProducer
	shouldError  bool
	errorMessage string
	mutex        sync.Mutex
}

func NewMockProducerFactory() *MockProducerFactory {
	return &MockProducerFactory{
		producers: make([]*MockProducer, 0),
	}
}

func (f *MockProducerFactory) NewProducer(configMap *kafka.ConfigMap) (KafkaProducerInterface, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.shouldError {
		return nil, fmt.Errorf("%s", f.errorMessage)
	}

	producer := NewMockProducer()
	f.producers = append(f.producers, producer)
	return producer, nil
}

func (f *MockProducerFactory) SetError(shouldError bool, errorMessage string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.shouldError = shouldError
	f.errorMessage = errorMessage
}

func (f *MockProducerFactory) GetProducers() []*MockProducer {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	result := make([]*MockProducer, len(f.producers))
	copy(result, f.producers)
	return result
}

// MockConsumerFactory implements ConsumerFactoryInterface for testing
type MockConsumerFactory struct {
	consumers    []*MockConsumer
	shouldError  bool
	errorMessage string
	mutex        sync.Mutex
}

func NewMockConsumerFactory() *MockConsumerFactory {
	return &MockConsumerFactory{
		consumers: make([]*MockConsumer, 0),
	}
}

func (f *MockConsumerFactory) NewConsumer(configMap *kafka.ConfigMap) (KafkaConsumerInterface, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.shouldError {
		return nil, fmt.Errorf("%s", f.errorMessage)
	}

	consumer := NewMockConsumer()
	f.consumers = append(f.consumers, consumer)
	return consumer, nil
}

func (f *MockConsumerFactory) SetError(shouldError bool, errorMessage string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.shouldError = shouldError
	f.errorMessage = errorMessage
}

func (f *MockConsumerFactory) GetConsumers() []*MockConsumer {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	result := make([]*MockConsumer, len(f.consumers))
	copy(result, f.consumers)
	return result
}
