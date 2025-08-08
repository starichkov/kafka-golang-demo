package main

import (
	"context"
	"fmt"
	"kafka-golang-demo/internal/kafka"
	"kafka-golang-demo/internal/logging"
	"kafka-golang-demo/internal/metrics"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	metricsAddr := os.Getenv("METRICS_ADDR")
	if metricsAddr == "" {
		metricsAddr = ":8080"
	}

	// Start metrics server
	metricsServer := metrics.NewServer(metricsAddr)
	go func() {
		if err := metricsServer.Start(); err != nil {
			logging.Logger.Error("Metrics server failed", "error", err)
		}
	}()

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logging.Logger.Info("Received shutdown signal")
		cancel()
	}()

	// Create producer
	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// Send messages
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			msg := fmt.Sprintf("Message #%d", i)
			_ = producer.Send(msg)
			logging.Logger.Info("Message sent", "number", i, "message", msg)
			time.Sleep(1 * time.Second) // Allow time for metrics to be collected
		}
	}

	logging.Logger.Info("Producer finished, metrics available", "metrics_endpoint", fmt.Sprintf("http://localhost%s/metrics", metricsAddr))
	
	// Keep running to allow metrics collection
	<-ctx.Done()
	
	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	
	if err := metricsServer.Stop(shutdownCtx); err != nil {
		logging.Logger.Error("Error stopping metrics server", "error", err)
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
