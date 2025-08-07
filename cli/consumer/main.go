package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"kafka-golang-demo/internal/kafka"
)

func main() {
	brokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		topic = "demo-topic"
	}

	groupID := os.Getenv("GROUP_ID")
	if groupID == "" {
		groupID = "golang-demo-group"
	}

	consumer, err := kafka.NewConsumer(brokers, groupID, topic)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Handle SIGINT / SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	consumer.Run(ctx)
}

//package main
//
//import (
//	"context"
//	"fmt"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//)
//
//func main() {
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers": "localhost:9092",
//		"group.id":          "golang-demo-group",
//		"auto.offset.reset": "earliest",
//	})
//	if err != nil {
//		panic(fmt.Sprintf("Failed to create consumer: %s", err))
//	}
//	defer consumer.Close()
//
//	topic := "demo-topic"
//	if err := consumer.SubscribeTopics([]string{topic}, nil); err != nil {
//		panic(fmt.Sprintf("Failed to subscribe: %s", err))
//	}
//
//	fmt.Println("🟢 Listening for messages... (Ctrl+C to stop)")
//
//	// Set up cancelable context for graceful shutdown
//	ctx, cancel := context.WithCancel(context.Background())
//
//	// Signal handler for Ctrl+C
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//	go func() {
//		<-sigs
//		fmt.Println("🔴 Interrupt received, shutting down...")
//		cancel()
//	}()
//
//loop:
//	for {
//		select {
//		case <-ctx.Done():
//			break loop
//		default:
//			msg, err := consumer.ReadMessage(500 * time.Millisecond)
//			if err != nil {
//				// Ignore timeout errors, only print real failures
//				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
//					continue
//				}
//				fmt.Printf("❌ Consumer error: %v\n", err)
//				continue
//			}
//			if msg != nil {
//				fmt.Printf("📥 Received: %s from %s [%d] offset %d\n",
//					string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
//			}
//		}
//	}
//
//	fmt.Println("🛑 Consumer closed.")
//}
