package main

import (
	"fmt"
	"golang-kafka-demo/internal/kafka"
	"os"
)

func main() {
	brokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	producer, err := kafka.NewProducer(brokers, "demo-topic")
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Message #%d", i)
		_ = producer.Send(msg)
	}
}

//package main
//
//import (
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//func main() {
//	// Create producer with bootstrap server
//	producer, err := kafka.NewProducer(&kafka.ConfigMap{
//		"bootstrap.servers": "localhost:9092",
//	})
//	if err != nil {
//		panic(fmt.Sprintf("Failed to create producer: %s", err))
//	}
//	defer producer.Close()
//
//	// Delivery report handler (runs in background goroutine)
//	go func() {
//		for e := range producer.Events() {
//			switch ev := e.(type) {
//			case *kafka.Message:
//				if ev.TopicPartition.Error != nil {
//					fmt.Printf("❌ Delivery failed: %v\n", ev.TopicPartition)
//				} else {
//					fmt.Printf("✅ Delivered to %v\n", ev.TopicPartition)
//				}
//			}
//		}
//	}()
//
//	// Handle Ctrl+C to allow graceful shutdown
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//	topic := "demo-topic"
//
//	for i := 0; i < 10; i++ {
//		value := fmt.Sprintf("Message #%d", i)
//		err := producer.Produce(&kafka.Message{
//			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
//			Value:          []byte(value),
//		}, nil)
//
//		if err != nil {
//			fmt.Printf("❌ Produce error: %v\n", err)
//		}
//	}
//
//	// Wait for message deliveries (flush)
//	fmt.Println("🔄 Flushing...")
//	producer.Flush(3000)
//
//	fmt.Println("✅ Done producing.")
//}
//
////import (
////	"fmt"
////	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
////)
////
////func main() {
////	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
////	if err != nil {
////		panic(err)
////	}
////	defer p.Close()
////
////	topic := "demo-topic"
////	for i := 0; i < 10; i++ {
////		msg := fmt.Sprintf("Message #%d", i)
////		err = p.Produce(&kafka.Message{
////			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
////			Value:          []byte(msg),
////		}, nil)
////		if err != nil {
////			fmt.Println("Produce error:", err)
////		}
////	}
////
////	// Wait for message deliveries before shutting down
////	p.Flush(1000)
////}
