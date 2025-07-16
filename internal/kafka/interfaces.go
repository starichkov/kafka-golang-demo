package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

// KafkaProducerInterface abstracts the kafka.Producer functionality
type KafkaProducerInterface interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	Close()
}

// KafkaConsumerInterface abstracts the kafka.Consumer functionality
type KafkaConsumerInterface interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

// ProducerFactoryInterface abstracts the creation of Kafka producers
type ProducerFactoryInterface interface {
	NewProducer(configMap *kafka.ConfigMap) (KafkaProducerInterface, error)
}

// ConsumerFactoryInterface abstracts the creation of Kafka consumers
type ConsumerFactoryInterface interface {
	NewConsumer(configMap *kafka.ConfigMap) (KafkaConsumerInterface, error)
}

// DefaultProducerFactory implements ProducerFactoryInterface using real Kafka
type DefaultProducerFactory struct{}

func (f *DefaultProducerFactory) NewProducer(configMap *kafka.ConfigMap) (KafkaProducerInterface, error) {
	return kafka.NewProducer(configMap)
}

// DefaultConsumerFactory implements ConsumerFactoryInterface using real Kafka
type DefaultConsumerFactory struct{}

func (f *DefaultConsumerFactory) NewConsumer(configMap *kafka.ConfigMap) (KafkaConsumerInterface, error) {
	return kafka.NewConsumer(configMap)
}

// Global factory instances (can be replaced with mocks in tests)
var (
	ProducerFactory ProducerFactoryInterface = &DefaultProducerFactory{}
	ConsumerFactory ConsumerFactoryInterface = &DefaultConsumerFactory{}
)
