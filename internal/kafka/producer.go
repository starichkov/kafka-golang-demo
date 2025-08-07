package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka-golang-demo/internal/logging"
)

type Producer struct {
	p     KafkaProducerInterface
	topic string
}

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

func (pr *Producer) Send(message string) error {
	return pr.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &pr.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
}

func (pr *Producer) Close() {
	pr.p.Flush(3000)
	pr.p.Close()
}
