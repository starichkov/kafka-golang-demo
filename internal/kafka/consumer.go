package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	c     KafkaConsumerInterface
	topic string
}

func NewConsumer(brokers, groupID, topic string) (*Consumer, error) {
	c, err := ConsumerFactory.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{c: c, topic: topic}, nil
}

func (co *Consumer) Run(ctx context.Context) {
	co.RunWithChannel(ctx, nil)
}

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
				fmt.Printf("❌ Consumer error: %v\n", err)
				continue
			}
			fmt.Printf("📥 Received: %s from %s [%d] offset %d\n",
				string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
			if msgCh != nil {
				msgCh <- msg
			}
		}
	}
}

func (co *Consumer) Close() {
	err := co.c.Close()
	if err != nil {
		return
	}
}
