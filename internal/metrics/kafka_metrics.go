// Package metrics provides Prometheus metrics collection for Kafka operations.
package metrics

import (
	"encoding/json"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// KafkaMetrics holds all Kafka-related Prometheus metrics
type KafkaMetrics struct {
	// Message counters
	MessagesProduced *prometheus.CounterVec
	MessagesConsumed *prometheus.CounterVec
	
	// Byte counters
	BytesProduced *prometheus.CounterVec
	BytesConsumed *prometheus.CounterVec
	
	// Error counters
	ProducerErrors *prometheus.CounterVec
	ConsumerErrors *prometheus.CounterVec
	
	// Consumer lag gauge
	ConsumerLag *prometheus.GaugeVec
	
	// Request latency histogram
	RequestLatency *prometheus.HistogramVec
}

var (
	kafkaMetrics *KafkaMetrics
	once         sync.Once
)

// GetKafkaMetrics returns the singleton KafkaMetrics instance
func GetKafkaMetrics() *KafkaMetrics {
	once.Do(func() {
		kafkaMetrics = &KafkaMetrics{
			MessagesProduced: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_messages_produced_total",
					Help: "Total number of messages produced to Kafka topics",
				},
				[]string{"topic", "partition"},
			),
			MessagesConsumed: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_messages_consumed_total",
					Help: "Total number of messages consumed from Kafka topics",
				},
				[]string{"topic", "partition", "consumer_group"},
			),
			BytesProduced: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_bytes_produced_total",
					Help: "Total number of bytes produced to Kafka topics",
				},
				[]string{"topic", "partition"},
			),
			BytesConsumed: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_bytes_consumed_total",
					Help: "Total number of bytes consumed from Kafka topics",
				},
				[]string{"topic", "partition", "consumer_group"},
			),
			ProducerErrors: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_producer_errors_total",
					Help: "Total number of producer errors",
				},
				[]string{"topic", "error_code"},
			),
			ConsumerErrors: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "kafka_consumer_errors_total",
					Help: "Total number of consumer errors",
				},
				[]string{"topic", "error_code"},
			),
			ConsumerLag: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "kafka_consumer_lag",
					Help: "Current consumer lag in messages",
				},
				[]string{"topic", "partition", "consumer_group"},
			),
			RequestLatency: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "kafka_request_latency_seconds",
					Help:    "Request latency in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"operation", "broker"},
			),
		}
	})
	return kafkaMetrics
}

// LibrdkafkaStats represents the structure of librdkafka statistics JSON
type LibrdkafkaStats struct {
	Name     string                     `json:"name"`
	Type     string                     `json:"type"`
	TxMsgs   int64                      `json:"txmsgs"`
	TxBytes  int64                      `json:"txbytes"`
	RxMsgs   int64                      `json:"rxmsgs"`
	RxBytes  int64                      `json:"rxbytes"`
	TxErr    int64                      `json:"txerrs"`
	RxErr    int64                      `json:"rxerrs"`
	Brokers  map[string]BrokerStats     `json:"brokers"`
	Topics   map[string]TopicStats      `json:"topics"`
	CGroups  map[string]ConsumerGroup   `json:"cgrp,omitempty"`
}

// BrokerStats contains broker-specific statistics
type BrokerStats struct {
	Name     string  `json:"name"`
	NodeID   int32   `json:"nodeid"`
	State    string  `json:"state"`
	Rtt      Rtt     `json:"rtt"`
}

// Rtt contains round-trip time statistics
type Rtt struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
	Avg int64 `json:"avg"`
}

// TopicStats contains topic-specific statistics
type TopicStats struct {
	Topic      string                     `json:"topic"`
	Partitions map[string]PartitionStats  `json:"partitions"`
}

// PartitionStats contains partition-specific statistics
type PartitionStats struct {
	Partition  int32 `json:"partition"`
	Leader     int32 `json:"leader"`
	TxMsgs     int64 `json:"txmsgs"`
	TxBytes    int64 `json:"txbytes"`
	RxMsgs     int64 `json:"rxmsgs"`
	RxBytes    int64 `json:"rxbytes"`
	ConsumerLag int64 `json:"consumer_lag"`
}

// ConsumerGroup contains consumer group statistics
type ConsumerGroup struct {
	State       string `json:"state"`
	StateAge    int64  `json:"stateage"`
	JoinState   string `json:"join_state"`
	RebalanceAge int64 `json:"rebalance_age"`
}

// ProcessLibrdkafkaStats processes librdkafka statistics JSON and updates Prometheus metrics
func ProcessLibrdkafkaStats(statsJSON string) error {
	var stats LibrdkafkaStats
	if err := json.Unmarshal([]byte(statsJSON), &stats); err != nil {
		return err
	}

	metrics := GetKafkaMetrics()
	
	// Process topic-level metrics
	for topicName, topicStats := range stats.Topics {
		for partitionID, partitionStats := range topicStats.Partitions {
			// Producer metrics (if this is a producer)
			if stats.Type == "producer" {
				metrics.MessagesProduced.WithLabelValues(topicName, partitionID).Add(float64(partitionStats.TxMsgs))
				metrics.BytesProduced.WithLabelValues(topicName, partitionID).Add(float64(partitionStats.TxBytes))
			}
			
			// Consumer metrics (if this is a consumer)
			if stats.Type == "consumer" {
				// We need the consumer group from cgrp, use "unknown" if not available
				consumerGroup := "unknown"
				for cgName := range stats.CGroups {
					consumerGroup = cgName
					break // Use first consumer group found
				}
				
				metrics.MessagesConsumed.WithLabelValues(topicName, partitionID, consumerGroup).Add(float64(partitionStats.RxMsgs))
				metrics.BytesConsumed.WithLabelValues(topicName, partitionID, consumerGroup).Add(float64(partitionStats.RxBytes))
				
				// Set consumer lag if available
				if partitionStats.ConsumerLag >= 0 {
					metrics.ConsumerLag.WithLabelValues(topicName, partitionID, consumerGroup).Set(float64(partitionStats.ConsumerLag))
				}
			}
		}
	}
	
	// Process broker-level metrics
	for brokerName, brokerStats := range stats.Brokers {
		// Update request latency histogram with RTT data
		if brokerStats.Rtt.Avg > 0 {
			// Convert microseconds to seconds for Prometheus convention
			latencySeconds := float64(brokerStats.Rtt.Avg) / 1000000.0
			operation := "rtt"
			if stats.Type == "producer" {
				operation = "produce"
			} else if stats.Type == "consumer" {
				operation = "fetch"
			}
			metrics.RequestLatency.WithLabelValues(operation, brokerName).Observe(latencySeconds)
		}
	}
	
	// Process global error counters
	if stats.TxErr > 0 {
		metrics.ProducerErrors.WithLabelValues("global", "tx_error").Add(float64(stats.TxErr))
	}
	if stats.RxErr > 0 {
		metrics.ConsumerErrors.WithLabelValues("global", "rx_error").Add(float64(stats.RxErr))
	}
	
	return nil
}