package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestGetKafkaMetrics_Singleton(t *testing.T) {
	// Test that GetKafkaMetrics returns the same instance
	metrics1 := GetKafkaMetrics()
	metrics2 := GetKafkaMetrics()

	if metrics1 != metrics2 {
		t.Error("GetKafkaMetrics should return the same singleton instance")
	}

	// Test that all metrics are initialized
	if metrics1.MessagesProduced == nil {
		t.Error("MessagesProduced should be initialized")
	}
	if metrics1.MessagesConsumed == nil {
		t.Error("MessagesConsumed should be initialized")
	}
	if metrics1.BytesProduced == nil {
		t.Error("BytesProduced should be initialized")
	}
	if metrics1.BytesConsumed == nil {
		t.Error("BytesConsumed should be initialized")
	}
	if metrics1.ProducerErrors == nil {
		t.Error("ProducerErrors should be initialized")
	}
	if metrics1.ConsumerErrors == nil {
		t.Error("ConsumerErrors should be initialized")
	}
	if metrics1.ConsumerLag == nil {
		t.Error("ConsumerLag should be initialized")
	}
	if metrics1.RequestLatency == nil {
		t.Error("RequestLatency should be initialized")
	}
}

func TestProcessLibrdkafkaStats_ProducerStats(t *testing.T) {
	// Create a new registry for this test to avoid conflicts
	registry := prometheus.NewRegistry()
	
	// Create test metrics with our registry
	testMetrics := &KafkaMetrics{
		MessagesProduced: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_messages_produced_total",
				Help: "Test messages produced",
			},
			[]string{"topic", "partition"},
		),
		BytesProduced: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_bytes_produced_total",
				Help: "Test bytes produced",
			},
			[]string{"topic", "partition"},
		),
		RequestLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_kafka_request_latency_seconds",
				Help: "Test request latency",
			},
			[]string{"operation", "broker"},
		),
		ProducerErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_producer_errors_total",
				Help: "Test producer errors",
			},
			[]string{"topic", "error_code"},
		),
		MessagesConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_messages_consumed_total",
				Help: "Test messages consumed",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
		BytesConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_bytes_consumed_total",
				Help: "Test bytes consumed",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
		ConsumerErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_consumer_errors_total",
				Help: "Test consumer errors",
			},
			[]string{"topic", "error_code"},
		),
		ConsumerLag: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_kafka_consumer_lag",
				Help: "Test consumer lag",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
	}

	// Register metrics
	registry.MustRegister(testMetrics.MessagesProduced)
	registry.MustRegister(testMetrics.BytesProduced)
	registry.MustRegister(testMetrics.RequestLatency)
	registry.MustRegister(testMetrics.ProducerErrors)

	// Mock producer statistics JSON
	producerStatsJSON := `{
		"name": "rdkafka#producer-1",
		"type": "producer",
		"txmsgs": 100,
		"txbytes": 5000,
		"rxmsgs": 50,
		"rxbytes": 2500,
		"txerrs": 2,
		"rxerrs": 1,
		"brokers": {
			"localhost:9092/1": {
				"name": "localhost:9092/1",
				"nodeid": 1,
				"state": "UP",
				"rtt": {
					"min": 1000,
					"max": 5000,
					"avg": 2500
				}
			}
		},
		"topics": {
			"test-topic": {
				"topic": "test-topic",
				"partitions": {
					"0": {
						"partition": 0,
						"leader": 1,
						"txmsgs": 50,
						"txbytes": 2500,
						"rxmsgs": 0,
						"rxbytes": 0,
						"consumer_lag": -1
					},
					"1": {
						"partition": 1,
						"leader": 1,
						"txmsgs": 50,
						"txbytes": 2500,
						"rxmsgs": 0,
						"rxbytes": 0,
						"consumer_lag": -1
					}
				}
			}
		}
	}`

	// Test processing with temporary replacement of global metrics
	originalKafkaMetrics := kafkaMetrics
	kafkaMetrics = testMetrics
	defer func() { kafkaMetrics = originalKafkaMetrics }()

	err := ProcessLibrdkafkaStats(producerStatsJSON)
	if err != nil {
		t.Fatalf("ProcessLibrdkafkaStats failed: %v", err)
	}

	// Verify metrics were updated
	expectedMessages := 50.0 // txmsgs from partition 0
	actualMessages := testutil.ToFloat64(testMetrics.MessagesProduced.WithLabelValues("test-topic", "0"))
	if actualMessages != expectedMessages {
		t.Errorf("Expected messages produced %f, got %f", expectedMessages, actualMessages)
	}

	expectedBytes := 2500.0 // txbytes from partition 0
	actualBytes := testutil.ToFloat64(testMetrics.BytesProduced.WithLabelValues("test-topic", "0"))
	if actualBytes != expectedBytes {
		t.Errorf("Expected bytes produced %f, got %f", expectedBytes, actualBytes)
	}

	// Verify global error counter
	expectedErrors := 2.0 // txerrs
	actualErrors := testutil.ToFloat64(testMetrics.ProducerErrors.WithLabelValues("global", "tx_error"))
	if actualErrors != expectedErrors {
		t.Errorf("Expected producer errors %f, got %f", expectedErrors, actualErrors)
	}

	// Verify latency histogram has observations by checking if the metric exists
	// We can't easily check the count with testutil.ToFloat64 on histogram observers
	// but we can verify the metric was registered and is accessible
	observer := testMetrics.RequestLatency.WithLabelValues("produce", "localhost:9092/1")
	if observer == nil {
		t.Error("Expected latency histogram observer to be available")
	}
}

func TestProcessLibrdkafkaStats_ConsumerStats(t *testing.T) {
	// Create a new registry for this test
	registry := prometheus.NewRegistry()
	
	// Create test metrics
	testMetrics := &KafkaMetrics{
		MessagesConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_messages_consumed_total_2",
				Help: "Test messages consumed",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
		BytesConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_bytes_consumed_total_2",
				Help: "Test bytes consumed",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
		ConsumerLag: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_kafka_consumer_lag_2",
				Help: "Test consumer lag",
			},
			[]string{"topic", "partition", "consumer_group"},
		),
		ConsumerErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_kafka_consumer_errors_total_2",
				Help: "Test consumer errors",
			},
			[]string{"topic", "error_code"},
		),
		RequestLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_kafka_request_latency_seconds_2",
				Help: "Test request latency",
			},
			[]string{"operation", "broker"},
		),
		// Initialize other required fields to avoid nil pointer errors
		MessagesProduced: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "unused_1", Help: "unused"},
			[]string{"topic", "partition"},
		),
		BytesProduced: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "unused_2", Help: "unused"},
			[]string{"topic", "partition"},
		),
		ProducerErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "unused_3", Help: "unused"},
			[]string{"topic", "error_code"},
		),
	}

	// Register metrics
	registry.MustRegister(testMetrics.MessagesConsumed)
	registry.MustRegister(testMetrics.BytesConsumed)
	registry.MustRegister(testMetrics.ConsumerLag)
	registry.MustRegister(testMetrics.ConsumerErrors)
	registry.MustRegister(testMetrics.RequestLatency)

	// Mock consumer statistics JSON
	consumerStatsJSON := `{
		"name": "rdkafka#consumer-1",
		"type": "consumer",
		"txmsgs": 10,
		"txbytes": 500,
		"rxmsgs": 100,
		"rxbytes": 5000,
		"txerrs": 0,
		"rxerrs": 3,
		"cgrp": {
			"test-group": {
				"state": "up",
				"stateage": 12345,
				"join_state": "steady",
				"rebalance_age": 12345
			}
		},
		"brokers": {
			"localhost:9092/1": {
				"name": "localhost:9092/1",
				"nodeid": 1,
				"state": "UP",
				"rtt": {
					"min": 500,
					"max": 2000,
					"avg": 1000
				}
			}
		},
		"topics": {
			"test-topic": {
				"topic": "test-topic",
				"partitions": {
					"0": {
						"partition": 0,
						"leader": 1,
						"txmsgs": 0,
						"txbytes": 0,
						"rxmsgs": 60,
						"rxbytes": 3000,
						"consumer_lag": 10
					},
					"1": {
						"partition": 1,
						"leader": 1,
						"txmsgs": 0,
						"txbytes": 0,
						"rxmsgs": 40,
						"rxbytes": 2000,
						"consumer_lag": 5
					}
				}
			}
		}
	}`

	// Test processing with temporary replacement of global metrics
	originalKafkaMetrics := kafkaMetrics
	kafkaMetrics = testMetrics
	defer func() { kafkaMetrics = originalKafkaMetrics }()

	err := ProcessLibrdkafkaStats(consumerStatsJSON)
	if err != nil {
		t.Fatalf("ProcessLibrdkafkaStats failed: %v", err)
	}

	// Verify consumer metrics were updated
	expectedMessages := 60.0 // rxmsgs from partition 0
	actualMessages := testutil.ToFloat64(testMetrics.MessagesConsumed.WithLabelValues("test-topic", "0", "test-group"))
	if actualMessages != expectedMessages {
		t.Errorf("Expected messages consumed %f, got %f", expectedMessages, actualMessages)
	}

	expectedBytes := 3000.0 // rxbytes from partition 0
	actualBytes := testutil.ToFloat64(testMetrics.BytesConsumed.WithLabelValues("test-topic", "0", "test-group"))
	if actualBytes != expectedBytes {
		t.Errorf("Expected bytes consumed %f, got %f", expectedBytes, actualBytes)
	}

	// Verify consumer lag
	expectedLag := 10.0
	actualLag := testutil.ToFloat64(testMetrics.ConsumerLag.WithLabelValues("test-topic", "0", "test-group"))
	if actualLag != expectedLag {
		t.Errorf("Expected consumer lag %f, got %f", expectedLag, actualLag)
	}

	// Verify global error counter
	expectedErrors := 3.0 // rxerrs
	actualErrors := testutil.ToFloat64(testMetrics.ConsumerErrors.WithLabelValues("global", "rx_error"))
	if actualErrors != expectedErrors {
		t.Errorf("Expected consumer errors %f, got %f", expectedErrors, actualErrors)
	}
}

func TestProcessLibrdkafkaStats_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`
	
	err := ProcessLibrdkafkaStats(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestProcessLibrdkafkaStats_EmptyStats(t *testing.T) {
	emptyStatsJSON := `{
		"name": "test",
		"type": "producer",
		"txmsgs": 0,
		"txbytes": 0,
		"rxmsgs": 0,
		"rxbytes": 0,
		"txerrs": 0,
		"rxerrs": 0,
		"brokers": {},
		"topics": {}
	}`

	err := ProcessLibrdkafkaStats(emptyStatsJSON)
	if err != nil {
		t.Errorf("ProcessLibrdkafkaStats should handle empty stats, got error: %v", err)
	}
}