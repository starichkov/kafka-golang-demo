# Structured Logging with Go's slog - Implementation Guide

**Created**: August 5, 2025  
**Purpose**: Replace basic fmt.Printf logging with production-grade structured logging using Go's standard library slog

## Overview

This document outlines the implementation of structured logging with correlation tracking using Go's official `slog` package (available since Go 1.21). This transforms basic console output into production-ready, machine-parseable logs that are essential for distributed systems and enterprise applications.

## Why slog for Kafka Applications

### Current State (Basic Logging)
```go
// Unstructured, hard to parse, no correlation
fmt.Printf("📤 Sent: Message #%d", i)
fmt.Printf("✅ Delivered to %v\n", ev.TopicPartition)
fmt.Printf("📥 Received: %s from %s [%d] offset %d\n", 
    string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
```

**Problems:**
- ❌ **Unstructured** - difficult to parse programmatically
- ❌ **No correlation** - can't trace messages across services
- ❌ **Missing context** - no timestamps, service info, log levels
- ❌ **Not production-ready** - can't integrate with ELK, Splunk, etc.

### Enhanced State (slog Structured Logging)
```json
{
  "time": "2025-08-05T15:30:45.123Z",
  "level": "INFO",
  "msg": "Message produced successfully",
  "correlation_id": "abc-123-def-456",
  "service": "kafka-producer",
  "kafka.topic": "demo-topic",
  "kafka.partition": 0,
  "kafka.offset": 42,
  "message.size_bytes": 156,
  "source": "internal/kafka/producer.go:67"
}
```

**Benefits:**
- ✅ **Structured JSON** - machine-parseable, searchable
- ✅ **Correlation tracking** - trace messages across services
- ✅ **Rich context** - service, topic, partition, offset information
- ✅ **Production-ready** - integrates with all log aggregation systems
- ✅ **Standard library** - no external dependencies, future-proof

---

## Implementation Architecture

### Package Structure
```
internal/
├── logging/
│   ├── slog.go           # slog configuration and helpers
│   ├── context.go        # Correlation ID context management
│   └── kafka_logger.go   # Kafka-specific logging functions
├── kafka/
│   ├── producer.go       # Enhanced with structured logging
│   └── consumer.go       # Enhanced with structured logging
└── headers/
    └── correlation.go    # Kafka header utilities
```

---

## Core Implementation

### 1. Logging Configuration and Setup

```go
// internal/logging/slog.go
package logging

import (
    "context"
    "log/slog"
    "os"
    "strings"
)

// Global logger instance - configured once, used everywhere
var Logger *slog.Logger

// Initialize logger on package import
func init() {
    // Create JSON handler for structured output
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: true, // Include source file and line number
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            // Customize attribute names for better readability
            switch a.Key {
            case slog.TimeKey:
                return slog.Attr{Key: "timestamp", Value: a.Value}
            case slog.MessageKey:
                return slog.Attr{Key: "message", Value: a.Value}
            case slog.SourceKey:
                // Shorten source paths
                if src, ok := a.Value.Any().(*slog.Source); ok {
                    // Convert /path/to/kafka-golang-demo/internal/kafka/producer.go:67
                    // to internal/kafka/producer.go:67
                    parts := strings.Split(src.File, "/")
                    for i, part := range parts {
                        if part == "kafka-golang-demo" && i+1 < len(parts) {
                            src.File = strings.Join(parts[i+1:], "/")
                            break
                        }
                    }
                }
            }
            return a
        },
    })
    
    Logger = slog.New(handler)
    
    // Log initialization
    Logger.Info("Structured logging initialized",
        "logger", "slog",
        "format", "json",
        "level", "info",
    )
}

// LogLevel configuration
type LogLevel string

const (
    LevelDebug LogLevel = "debug"
    LevelInfo  LogLevel = "info"
    LevelWarn  LogLevel = "warn"
    LevelError LogLevel = "error"
)

// SetLogLevel allows runtime log level changes
func SetLogLevel(level LogLevel) {
    var slogLevel slog.Level
    
    switch level {
    case LevelDebug:
        slogLevel = slog.LevelDebug
    case LevelInfo:
        slogLevel = slog.LevelInfo
    case LevelWarn:
        slogLevel = slog.LevelWarn
    case LevelError:
        slogLevel = slog.LevelError
    default:
        slogLevel = slog.LevelInfo
    }
    
    // Create new handler with updated level
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slogLevel,
        AddSource: true,
    })
    
    Logger = slog.New(handler)
    Logger.Info("Log level updated", "new_level", level)
}
```

### 2. Correlation ID Context Management

```go
// internal/logging/context.go
package logging

import (
    "context"
    "log/slog"
    
    "github.com/google/uuid"
)

// Context keys for correlation tracking
type contextKey string

const (
    CorrelationIDKey contextKey = "correlation_id"
    TraceIDKey       contextKey = "trace_id"
    ServiceNameKey   contextKey = "service_name"
    UserIDKey        contextKey = "user_id"
)

// Correlation ID management
func NewCorrelationID() string {
    return uuid.New().String()
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
    return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
    if id := ctx.Value(CorrelationIDKey); id != nil {
        if str, ok := id.(string); ok {
            return str
        }
    }
    return ""
}

// Trace ID management (for distributed tracing integration)
func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, TraceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
    if id := ctx.Value(TraceIDKey); id != nil {
        if str, ok := id.(string); ok {
            return str
        }
    }
    return ""
}

// Service name management
func WithServiceName(ctx context.Context, serviceName string) context.Context {
    return context.WithValue(ctx, ServiceNameKey, serviceName)
}

func GetServiceName(ctx context.Context) string {
    if name := ctx.Value(ServiceNameKey); name != nil {
        if str, ok := name.(string); ok {
            return str
        }
    }
    return "unknown-service"
}

// User ID management (for business context)
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) string {
    if id := ctx.Value(UserIDKey); id != nil {
        if str, ok := id.(string); ok {
            return str
        }
    }
    return ""
}

// Logger with context - extracts all context information
func LoggerFromContext(ctx context.Context) *slog.Logger {
    logger := Logger
    
    // Add correlation ID if present
    if correlationID := GetCorrelationID(ctx); correlationID != "" {
        logger = logger.With("correlation_id", correlationID)
    }
    
    // Add trace ID if present
    if traceID := GetTraceID(ctx); traceID != "" {
        logger = logger.With("trace_id", traceID)
    }
    
    // Add service name if present
    if serviceName := GetServiceName(ctx); serviceName != "" {
        logger = logger.With("service", serviceName)
    }
    
    // Add user ID if present (business context)
    if userID := GetUserID(ctx); userID != "" {
        logger = logger.With("user_id", userID)
    }
    
    return logger
}

// Create enriched context for new operations
func NewOperationContext(serviceName string) context.Context {
    ctx := context.Background()
    ctx = WithCorrelationID(ctx, NewCorrelationID())
    ctx = WithServiceName(ctx, serviceName)
    return ctx
}

// Create child context that inherits correlation but can add new context
func NewChildContext(parent context.Context, operation string) context.Context {
    // Inherit correlation ID from parent
    correlationID := GetCorrelationID(parent)
    if correlationID == "" {
        correlationID = NewCorrelationID()
    }
    
    ctx := WithCorrelationID(context.Background(), correlationID)
    
    // Inherit other context
    if serviceName := GetServiceName(parent); serviceName != "" {
        ctx = WithServiceName(ctx, serviceName)
    }
    if traceID := GetTraceID(parent); traceID != "" {
        ctx = WithTraceID(ctx, traceID)
    }
    if userID := GetUserID(parent); userID != "" {
        ctx = WithUserID(ctx, userID)
    }
    
    return ctx
}
```

### 3. Kafka-Specific Logging Functions

```go
// internal/logging/kafka_logger.go
package logging

import (
    "context"
    "log/slog"
    "time"
)

// Kafka event types for consistent logging
const (
    EventKafkaProducerStarted     = "kafka.producer.started"
    EventKafkaProducerStopped     = "kafka.producer.stopped"
    EventKafkaMessageProduced    = "kafka.message.produced"
    EventKafkaMessageDelivered   = "kafka.message.delivered"
    EventKafkaMessageFailed      = "kafka.message.failed"
    EventKafkaConsumerStarted    = "kafka.consumer.started"
    EventKafkaConsumerStopped    = "kafka.consumer.stopped"
    EventKafkaMessageConsumed    = "kafka.message.consumed"
    EventKafkaMessageProcessed   = "kafka.message.processed"
    EventKafkaConsumerRebalanced = "kafka.consumer.rebalanced"
    EventKafkaError              = "kafka.error"
)

// Producer lifecycle logging
func LogProducerStarted(ctx context.Context, brokers string, topic string) {
    LoggerFromContext(ctx).Info("Kafka producer started",
        "event_type", EventKafkaProducerStarted,
        "kafka.brokers", brokers,
        "kafka.topic", topic,
    )
}

func LogProducerStopped(ctx context.Context, topic string) {
    LoggerFromContext(ctx).Info("Kafka producer stopped",
        "event_type", EventKafkaProducerStopped,
        "kafka.topic", topic,
    )
}

// Message production logging
func LogMessageProduced(ctx context.Context, topic, key string, sizeBytes int) {
    LoggerFromContext(ctx).Info("Message published to Kafka",
        "event_type", EventKafkaMessageProduced,
        "kafka.topic", topic,
        "message.key", key,
        "message.size_bytes", sizeBytes,
    )
}

func LogMessageDelivered(ctx context.Context, topic string, partition int32, offset int64, latencyMs int64) {
    LoggerFromContext(ctx).Info("Message delivery confirmed",
        "event_type", EventKafkaMessageDelivered,
        "kafka.topic", topic,
        "kafka.partition", partition,
        "kafka.offset", offset,
        "delivery.latency_ms", latencyMs,
    )
}

func LogMessageDeliveryFailed(ctx context.Context, topic string, err error, retryAttempt int) {
    LoggerFromContext(ctx).Error("Message delivery failed",
        "event_type", EventKafkaMessageFailed,
        "kafka.topic", topic,
        "error", err.Error(),
        "retry_attempt", retryAttempt,
    )
}

// Consumer lifecycle logging
func LogConsumerStarted(ctx context.Context, brokers string, topic string, groupID string) {
    LoggerFromContext(ctx).Info("Kafka consumer started",
        "event_type", EventKafkaConsumerStarted,
        "kafka.brokers", brokers,
        "kafka.topic", topic,
        "kafka.consumer_group", groupID,
    )
}

func LogConsumerStopped(ctx context.Context, topic string, groupID string) {
    LoggerFromContext(ctx).Info("Kafka consumer stopped",
        "event_type", EventKafkaConsumerStopped,
        "kafka.topic", topic,
        "kafka.consumer_group", groupID,
    )
}

// Message consumption logging
func LogMessageConsumed(ctx context.Context, topic, key string, partition int32, offset int64, sizeBytes int) {
    LoggerFromContext(ctx).Info("Message consumed from Kafka",
        "event_type", EventKafkaMessageConsumed,
        "kafka.topic", topic,
        "message.key", key,
        "kafka.partition", partition,
        "kafka.offset", offset,
        "message.size_bytes", sizeBytes,
    )
}

func LogMessageProcessed(ctx context.Context, topic string, processingDurationMs int64, success bool) {
    level := slog.LevelInfo
    message := "Message processed successfully"
    
    if !success {
        level = slog.LevelError
        message = "Message processing failed"
    }
    
    LoggerFromContext(ctx).Log(ctx, level, message,
        "event_type", EventKafkaMessageProcessed,
        "kafka.topic", topic,
        "processing.duration_ms", processingDurationMs,
        "processing.success", success,
    )
}

// Consumer rebalancing
func LogConsumerRebalanced(ctx context.Context, groupID string, assignedPartitions []int32, revokedPartitions []int32) {
    LoggerFromContext(ctx).Info("Consumer group rebalanced",
        "event_type", EventKafkaConsumerRebalanced,
        "kafka.consumer_group", groupID,
        "partitions.assigned", assignedPartitions,
        "partitions.revoked", revokedPartitions,
    )
}

// Generic Kafka errors
func LogKafkaError(ctx context.Context, operation string, err error) {
    LoggerFromContext(ctx).Error("Kafka operation failed",
        "event_type", EventKafkaError,
        "kafka.operation", operation,
        "error", err.Error(),
    )
}

// Performance monitoring helpers
func LogOperationDuration(ctx context.Context, operation string, startTime time.Time, additionalFields ...slog.Attr) {
    duration := time.Since(startTime)
    
    attrs := []slog.Attr{
        slog.String("operation", operation),
        slog.Int64("duration_ms", duration.Milliseconds()),
        slog.Float64("duration_seconds", duration.Seconds()),
    }
    
    attrs = append(attrs, additionalFields...)
    
    LoggerFromContext(ctx).LogAttrs(ctx, slog.LevelInfo, "Operation completed", attrs...)
}
```

### 4. Kafka Header Utilities

```go
// internal/headers/correlation.go
package headers

import (
    "context"
    "time"
    
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "golang-kafka-demo/internal/logging"
)

// Header keys
const (
    HeaderCorrelationID = "correlation_id"
    HeaderTraceID       = "trace_id"
    HeaderSourceService = "source_service"
    HeaderTimestamp     = "timestamp"
    HeaderUserID        = "user_id"
    HeaderMessageType   = "message_type"
)

// Add correlation headers to Kafka message
func AddCorrelationHeaders(ctx context.Context, headers []kafka.Header) []kafka.Header {
    // Add correlation ID
    if correlationID := logging.GetCorrelationID(ctx); correlationID != "" {
        headers = append(headers, kafka.Header{
            Key:   HeaderCorrelationID,
            Value: []byte(correlationID),
        })
    }
    
    // Add trace ID
    if traceID := logging.GetTraceID(ctx); traceID != "" {
        headers = append(headers, kafka.Header{
            Key:   HeaderTraceID,
            Value: []byte(traceID),
        })
    }
    
    // Add source service
    if serviceName := logging.GetServiceName(ctx); serviceName != "" {
        headers = append(headers, kafka.Header{
            Key:   HeaderSourceService,
            Value: []byte(serviceName),
        })
    }
    
    // Add user ID if present
    if userID := logging.GetUserID(ctx); userID != "" {
        headers = append(headers, kafka.Header{
            Key:   HeaderUserID,
            Value: []byte(userID),
        })
    }
    
    // Add timestamp
    headers = append(headers, kafka.Header{
        Key:   HeaderTimestamp,
        Value: []byte(time.Now().Format(time.RFC3339Nano)),
    })
    
    return headers
}

// Extract correlation context from Kafka message headers
func ExtractCorrelationContext(headers []kafka.Header) context.Context {
    ctx := context.Background()
    
    for _, header := range headers {
        switch header.Key {
        case HeaderCorrelationID:
            ctx = logging.WithCorrelationID(ctx, string(header.Value))
        case HeaderTraceID:
            ctx = logging.WithTraceID(ctx, string(header.Value))
        case HeaderSourceService:
            // Don't override current service, but could be used for source tracking
        case HeaderUserID:
            ctx = logging.WithUserID(ctx, string(header.Value))
        }
    }
    
    return ctx
}

// Extract specific header value
func GetHeaderValue(headers []kafka.Header, key string) string {
    for _, header := range headers {
        if header.Key == key {
            return string(header.Value)
        }
    }
    return ""
}

// Check if header exists
func HasHeader(headers []kafka.Header, key string) bool {
    return GetHeaderValue(headers, key) != ""
}
```

---

## Enhanced Kafka Components

### 1. Enhanced Producer

```go
// internal/kafka/producer.go (enhanced with slog)
package kafka

import (
    "context"
    "fmt"
    "time"
    
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "golang-kafka-demo/internal/logging"
    "golang-kafka-demo/internal/headers"
)

type Producer struct {
    p     KafkaProducerInterface
    topic string
    ctx   context.Context // Service context
}

func NewProducer(brokers, topic string) (*Producer, error) {
    // Create service context
    ctx := logging.NewOperationContext("kafka-producer")
    
    // Log producer initialization
    logging.LoggerFromContext(ctx).Info("Initializing Kafka producer",
        "kafka.brokers", brokers,
        "kafka.topic", topic,
    )
    
    p, err := ProducerFactory.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": brokers,
        "client.id":         "golang-demo-producer",
        "acks":              "all",
        "retries":           3,
        "batch.size":        16384,
        "linger.ms":         5,
        "compression.type":  "snappy",
    })
    if err != nil {
        logging.LogKafkaError(ctx, "create_producer", err)
        return nil, err
    }

    producer := &Producer{
        p:     p,
        topic: topic,
        ctx:   ctx,
    }

    // Start delivery report handler
    go producer.handleDeliveryReports()

    // Log successful startup
    logging.LogProducerStarted(ctx, brokers, topic)

    return producer, nil
}

func (pr *Producer) Send(message string) error {
    return pr.SendWithContext(pr.ctx, "", message)
}

func (pr *Producer) SendWithContext(ctx context.Context, key, message string) error {
    start := time.Now()
    
    // Ensure we have correlation ID
    if logging.GetCorrelationID(ctx) == "" {
        ctx = logging.WithCorrelationID(ctx, logging.NewCorrelationID())
    }
    
    // Create Kafka message with correlation headers
    var kafkaHeaders []kafka.Header
    kafkaHeaders = headers.AddCorrelationHeaders(ctx, kafkaHeaders)
    
    kafkaMsg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &pr.topic, Partition: kafka.PartitionAny},
        Key:            []byte(key),
        Value:          []byte(message),
        Headers:        kafkaHeaders,
        Timestamp:      time.Now(),
    }
    
    // Produce message
    err := pr.p.Produce(kafkaMsg, nil)
    
    if err != nil {
        logging.LogKafkaError(ctx, "produce_message", err)
        return err
    }
    
    // Log successful production
    logging.LogMessageProduced(ctx, pr.topic, key, len(message))
    
    // Log operation timing
    logging.LogOperationDuration(ctx, "kafka_produce", start,
        slog.String("kafka.topic", pr.topic),
        slog.Int("message.size_bytes", len(message)),
    )
    
    return nil
}

func (pr *Producer) handleDeliveryReports() {
    for e := range pr.p.Events() {
        if ev, ok := e.(*kafka.Message); ok {
            // Extract correlation context from message headers
            msgCtx := headers.ExtractCorrelationContext(ev.Headers)
            
            // Ensure we have service context
            if logging.GetServiceName(msgCtx) == "" {
                msgCtx = logging.WithServiceName(msgCtx, "kafka-producer")
            }
            
            if ev.TopicPartition.Error != nil {
                logging.LogMessageDeliveryFailed(msgCtx, *ev.TopicPartition.Topic, ev.TopicPartition.Error, 0)
            } else {
                // Calculate delivery latency (simplified - in production you'd track send time)
                deliveryLatency := time.Since(ev.Timestamp).Milliseconds()
                
                logging.LogMessageDelivered(msgCtx, *ev.TopicPartition.Topic, 
                    ev.TopicPartition.Partition, int64(ev.TopicPartition.Offset), deliveryLatency)
            }
        }
    }
}

func (pr *Producer) Close() {
    logging.LoggerFromContext(pr.ctx).Info("Shutting down Kafka producer")
    
    // Flush pending messages
    unflushed := pr.p.Flush(30000) // 30 second timeout
    if unflushed > 0 {
        logging.LoggerFromContext(pr.ctx).Warn("Producer shutdown with unflushed messages",
            "unflushed_count", unflushed,
        )
    }
    
    pr.p.Close()
    logging.LogProducerStopped(pr.ctx, pr.topic)
}
```

### 2. Enhanced Consumer

```go
// internal/kafka/consumer.go (enhanced with slog)
package kafka

import (
    "context"
    "time"
    
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "golang-kafka-demo/internal/logging"
    "golang-kafka-demo/internal/headers"
)

type Consumer struct {
    c       KafkaConsumerInterface
    topic   string
    groupID string
    ctx     context.Context // Service context
}

func NewConsumer(brokers, groupID, topic string) (*Consumer, error) {
    // Create service context
    ctx := logging.NewOperationContext("kafka-consumer")
    
    logging.LoggerFromContext(ctx).Info("Initializing Kafka consumer",
        "kafka.brokers", brokers,
        "kafka.consumer_group", groupID,
        "kafka.topic", topic,
    )
    
    c, err := ConsumerFactory.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":        brokers,
        "group.id":                groupID,
        "auto.offset.reset":       "earliest",
        "enable.auto.commit":      false, // Manual commits for better control
        "session.timeout.ms":      6000,
        "heartbeat.interval.ms":   2000,
    })
    if err != nil {
        logging.LogKafkaError(ctx, "create_consumer", err)
        return nil, err
    }

    err = c.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        logging.LogKafkaError(ctx, "subscribe_topics", err)
        return nil, err
    }

    consumer := &Consumer{
        c:       c,
        topic:   topic,
        groupID: groupID,
        ctx:     ctx,
    }

    // Log successful startup
    logging.LogConsumerStarted(ctx, brokers, topic, groupID)

    return consumer, nil
}

func (co *Consumer) Run(ctx context.Context) {
    co.RunWithChannel(ctx, nil)
}

func (co *Consumer) RunWithChannel(ctx context.Context, msgCh chan<- *kafka.Message) {
    // Merge service context with operation context
    serviceCtx := logging.WithServiceName(ctx, "kafka-consumer")
    
    logging.LoggerFromContext(serviceCtx).Info("Starting message consumption loop",
        "kafka.topic", co.topic,
        "kafka.consumer_group", co.groupID,
    )

loop:
    for {
        select {
        case <-ctx.Done():
            logging.LoggerFromContext(serviceCtx).Info("Consumer shutdown requested",
                "reason", "context_cancelled",
            )
            break loop
        default:
            msg, err := co.c.ReadMessage(500 * time.Millisecond)
            if err != nil {
                if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
                    continue // Timeout is normal, not an error
                }
                logging.LogKafkaError(serviceCtx, "read_message", err)
                continue
            }
            
            // Extract correlation context from message
            msgCtx := headers.ExtractCorrelationContext(msg.Headers)
            // Ensure service context is preserved
            msgCtx = logging.WithServiceName(msgCtx, "kafka-consumer")
            
            // Process message with full context
            co.processMessage(msgCtx, msg, msgCh)
        }
    }
    
    logging.LogConsumerStopped(serviceCtx, co.topic, co.groupID)
}

func (co *Consumer) processMessage(ctx context.Context, msg *kafka.Message, msgCh chan<- *kafka.Message) {
    start := time.Now()
    
    // Log message consumption
    logging.LogMessageConsumed(ctx, *msg.TopicPartition.Topic, string(msg.Key),
        msg.TopicPartition.Partition, int64(msg.TopicPartition.Offset), len(msg.Value))
    
    // Extract message type for business context (if present in headers)
    messageType := headers.GetHeaderValue(msg.Headers, headers.HeaderMessageType)
    if messageType != "" {
        logging.LoggerFromContext(ctx).Info("Processing business message",
            "message.type", messageType,
            "kafka.offset", msg.TopicPartition.Offset,
        )
    }
    
    // Forward message if channel provided
    success := true
    if msgCh != nil {
        select {
        case msgCh <- msg:
            // Message forwarded successfully
        case <-time.After(5 * time.Second):
            logging.LoggerFromContext(ctx).Error("Message forwarding timeout",
                "timeout_seconds", 5,
            )
            success = false
        }
    }
    
    // Commit offset manually
    if success {
        if err := co.commitOffset(ctx, msg); err != nil {
            logging.LogKafkaError(ctx, "commit_offset", err)
            success = false
        }
    }
    
    // Log processing completion
    processingDuration := time.Since(start).Milliseconds()
    logging.LogMessageProcessed(ctx, *msg.TopicPartition.Topic, processingDuration, success)
}

func (co *Consumer) commitOffset(ctx context.Context, msg *kafka.Message) error {
    _, err := co.c.CommitMessage(msg)
    if err == nil {
        logging.LoggerFromContext(ctx).Debug("Offset committed",
            "kafka.topic", *msg.TopicPartition.Topic,
            "kafka.partition", msg.TopicPartition.Partition,
            "kafka.offset", msg.TopicPartition.Offset,
        )
    }
    return err
}

func (co *Consumer) Close() {
    logging.LoggerFromContext(co.ctx).Info("Shutting down Kafka consumer")
    
    if err := co.c.Close(); err != nil {
        logging.LogKafkaError(co.ctx, "close_consumer", err)
    }
}
```

### 3. Enhanced CLI Applications

```go
// cli/producer/main.go (enhanced with slog)
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang-kafka-demo/internal/kafka"
    "golang-kafka-demo/internal/logging"
)

func main() {
    // Create operation context
    ctx := logging.NewOperationContext("kafka-demo-producer")
    
    logging.LoggerFromContext(ctx).Info("Starting Kafka producer application")
    
    // Configuration
    brokers := getEnvOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
    topic := getEnvOrDefault("TOPIC", "demo-topic")
    messageCount := 10
    
    logging.LoggerFromContext(ctx).Info("Application configuration",
        "kafka.brokers", brokers,
        "kafka.topic", topic,
        "message.count", messageCount,
    )
    
    // Create producer
    producer, err := kafka.NewProducer(brokers, topic)
    if err != nil {
        logging.LogKafkaError(ctx, "create_producer", err)
        os.Exit(1)
    }
    defer producer.Close()
    
    // Graceful shutdown handling
    shutdownCtx, cancel := context.WithCancel(ctx)
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        sig := <-sigs
        logging.LoggerFromContext(ctx).Info("Shutdown signal received",
            "signal", sig.String(),
        )
        cancel()
    }()
    
    // Send messages
    for i := 0; i < messageCount; i++ {
        select {
        case <-shutdownCtx.Done():
            logging.LoggerFromContext(ctx).Info("Stopping message production due to shutdown")
            break
        default:
            messageCtx := logging.WithCorrelationID(ctx, logging.NewCorrelationID())
            
            msg := fmt.Sprintf("Message #%d produced at %s", i, time.Now().Format(time.RFC3339))
            key := fmt.Sprintf("key-%d", i)
            
            if err := producer.SendWithContext(messageCtx, key, msg); err != nil {
                logging.LogKafkaError(messageCtx, "send_message", err)
            } else {
                logging.LoggerFromContext(messageCtx).Info("Message sent successfully",
                    "message.number", i,
                    "message.key", key,
                )
            }
            
            // Small delay between messages for demo visibility
            time.Sleep(100 * time.Millisecond)
        }
    }
    
    logging.LoggerFromContext(ctx).Info("Message production completed",
        "total_messages", messageCount,
    )
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

```go  
// cli/consumer/main.go (enhanced with slog)
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "golang-kafka-demo/internal/kafka"
    "golang-kafka-demo/internal/logging"
)

func main() {
    // Create operation context
    ctx := logging.NewOperationContext("kafka-demo-consumer")
    
    logging.LoggerFromContext(ctx).Info("Starting Kafka consumer application")
    
    // Configuration
    brokers := getEnvOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
    topic := getEnvOrDefault("TOPIC", "demo-topic")
    groupID := getEnvOrDefault("GROUP_ID", "golang-demo-group")
    
    logging.LoggerFromContext(ctx).Info("Application configuration",
        "kafka.brokers", brokers,
        "kafka.topic", topic,
        "kafka.consumer_group", groupID,
    )
    
    // Create consumer
    consumer, err := kafka.NewConsumer(brokers, groupID, topic)
    if err != nil {
        logging.LogKafkaError(ctx, "create_consumer", err)
        os.Exit(1)
    }
    defer consumer.Close()

    // Graceful shutdown handling
    shutdownCtx, cancel := context.WithCancel(ctx)
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        sig := <-sigs
        logging.LoggerFromContext(ctx).Info("Shutdown signal received",
            "signal", sig.String(),
        )
        cancel()
    }()

    logging.LoggerFromContext(ctx).Info("Starting message consumption")
    
    // Start consuming
    consumer.Run(shutdownCtx)
    
    logging.LoggerFromContext(ctx).Info("Consumer application stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

---

## Dependencies

### Go Module Updates
```go
// go.mod additions
require github.com/google/uuid v1.4.0
```

**Note**: `log/slog` is part of the standard library since Go 1.21, so no additional logging dependencies are required!

---

## Example Output

### Before (Basic Printf Logging)
```
📤 Sent: Message #0
✅ Delivered to demo-topic [0] offset 42
📥 Received: Message #0 from demo-topic [0] offset 42
```

### After (Structured slog Logging)

**Producer Output:**
```json
{
  "timestamp": "2025-08-05T15:30:45.123Z",
  "level": "INFO",
  "message": "Structured logging initialized",
  "logger": "slog",
  "format": "json",
  "level": "info"
}

{
  "timestamp": "2025-08-05T15:30:45.125Z",
  "level": "INFO",
  "message": "Starting Kafka producer application",
  "correlation_id": "abc-123-def-456",
  "service": "kafka-demo-producer"
}

{
  "timestamp": "2025-08-05T15:30:45.130Z",
  "level": "INFO",
  "message": "Kafka producer started",
  "correlation_id": "abc-123-def-456",
  "service": "kafka-producer",
  "event_type": "kafka.producer.started",
  "kafka.brokers": "localhost:9092",
  "kafka.topic": "demo-topic"
}

{
  "timestamp": "2025-08-05T15:30:45.135Z",
  "level": "INFO",
  "message": "Message published to Kafka",
  "correlation_id": "def-456-ghi-789",
  "service": "kafka-producer",
  "event_type": "kafka.message.produced",
  "kafka.topic": "demo-topic",
  "message.key": "key-0",
  "message.size_bytes": 156
}

{
  "timestamp": "2025-08-05T15:30:45.140Z",
  "level": "INFO",
  "message": "Message delivery confirmed",
  "correlation_id": "def-456-ghi-789",
  "service": "kafka-producer",
  "event_type": "kafka.message.delivered",
  "kafka.topic": "demo-topic",
  "kafka.partition": 0,
  "kafka.offset": 42,
  "delivery.latency_ms": 5
}
```

**Consumer Output:**
```json
{
  "timestamp": "2025-08-05T15:30:45.145Z",
  "level": "INFO",
  "message": "Message consumed from Kafka",
  "correlation_id": "def-456-ghi-789",
  "service": "kafka-consumer", 
  "event_type": "kafka.message.consumed",
  "kafka.topic": "demo-topic",
  "message.key": "key-0",
  "kafka.partition": 0,
  "kafka.offset": 42,
  "message.size_bytes": 156
}

{
  "timestamp": "2025-08-05T15:30:45.147Z",
  "level": "INFO",
  "message": "Message processed successfully",
  "correlation_id": "def-456-ghi-789",
  "service": "kafka-consumer",
  "event_type": "kafka.message.processed",
  "kafka.topic": "demo-topic",
  "processing.duration_ms": 2,
  "processing.success": true
}
```

---

## Integration with Log Aggregation Systems

### ELK Stack Query Examples

**Find all messages for a correlation ID:**
```json
{
  "query": {
    "term": {
      "correlation_id": "def-456-ghi-789"
    }
  }
}
```

**Find producer errors:**
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"event_type": "kafka.message.failed"}},
        {"term": {"service": "kafka-producer"}}
      ]
    }
  }
}
```

**Find slow message processing:**
```json
{
  "query": {
    "range": {
      "processing.duration_ms": {
        "gte": 1000
      }
    }
  }
}
```

### Grafana Dashboard Queries

**Message throughput:**
```promql
rate(kafka_messages_produced_total[5m])
```

**Processing latency percentiles:**
```promql
histogram_quantile(0.95, rate(processing_duration_ms_bucket[5m]))
```

---

## Testing the Implementation

### Unit Tests Example

```go
// internal/logging/slog_test.go
package logging

import (
    "bytes"
    "context"
    "encoding/json"
    "log/slog"
    "testing"
    "time"
)

func TestStructuredLogging(t *testing.T) {
    // Capture log output
    var buf bytes.Buffer
    handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
    testLogger := slog.New(handler)
    
    // Temporarily replace global logger
    originalLogger := Logger
    Logger = testLogger
    defer func() { Logger = originalLogger }()
    
    // Test logging with correlation
    ctx := WithCorrelationID(context.Background(), "test-correlation-123")
    ctx = WithServiceName(ctx, "test-service")
    
    LogMessageProduced(ctx, "test-topic", "test-key", 100)
    
    // Parse logged JSON
    var logEntry map[string]interface{}
    if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
        t.Fatalf("Failed to parse log JSON: %v", err)
    }
    
    // Verify log structure
    expectedFields := map[string]interface{}{
        "correlation_id": "test-correlation-123",
        "service":        "test-service",
        "event_type":     EventKafkaMessageProduced,
        "kafka.topic":    "test-topic",
        "message.key":    "test-key",
        "message.size_bytes": 100,
    }
    
    for key, expectedValue := range expectedFields {
        if actualValue, exists := logEntry[key]; !exists {
            t.Errorf("Missing expected field: %s", key)
        } else if actualValue != expectedValue {
            t.Errorf("Field %s: expected %v, got %v", key, expectedValue, actualValue)
        }
    }
    
    // Verify timestamp format
    if timestamp, exists := logEntry["timestamp"]; exists {
        if _, err := time.Parse(time.RFC3339Nano, timestamp.(string)); err != nil {
            t.Errorf("Invalid timestamp format: %v", timestamp)
        }
    } else {
        t.Error("Missing timestamp field")
    }
}

func TestCorrelationContextPropagation(t *testing.T) {
    correlationID := "test-propagation-456"
    traceID := "trace-789"
    
    // Create parent context
    parentCtx := WithCorrelationID(context.Background(), correlationID)
    parentCtx = WithTraceID(parentCtx, traceID)
    
    // Create child context
    childCtx := NewChildContext(parentCtx, "child-operation")
    
    // Verify correlation ID is propagated
    if actualCorrelationID := GetCorrelationID(childCtx); actualCorrelationID != correlationID {
        t.Errorf("Correlation ID not propagated: expected %s, got %s", correlationID, actualCorrelationID)
    }
    
    // Verify trace ID is propagated
    if actualTraceID := GetTraceID(childCtx); actualTraceID != traceID {
        t.Errorf("Trace ID not propagated: expected %s, got %s", traceID, actualTraceID)
    }
}
```

### Integration Test Example

```go
// integration_test.go addition for structured logging
func TestStructuredLoggingIntegration(t *testing.T) {
    // Capture structured logs during integration test
    var logBuffer bytes.Buffer
    handler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })
    
    // Replace logger for test
    originalLogger := logging.Logger
    logging.Logger = slog.New(handler)
    defer func() { logging.Logger = originalLogger }()
    
    // Run integration test with structured logging
    ctx := context.Background()
    
    // ... existing integration test setup ...
    
    // Send test message with correlation tracking
    testCtx := logging.WithCorrelationID(ctx, "integration-test-correlation")
    err = producer.SendWithContext(testCtx, "test-key", "integration test message")
    if err != nil {
        t.Fatalf("Failed to send message: %v", err)
    }
    
    // Verify structured logs were generated
    logOutput := logBuffer.String()
    if logOutput == "" {
        t.Error("No structured logs generated during integration test")
    }
    
    // Verify correlation ID appears in logs
    if !strings.Contains(logOutput, "integration-test-correlation") {
        t.Error("Correlation ID not found in structured logs")
    }
    
    // Parse and verify log structure
    logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
    for _, line := range logLines {
        var logEntry map[string]interface{}
        if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
            t.Errorf("Invalid JSON log line: %s", line)
            continue
        }
        
        // Verify required fields are present
        requiredFields := []string{"timestamp", "level", "message"}
        for _, field := range requiredFields {
            if _, exists := logEntry[field]; !exists {
                t.Errorf("Missing required field '%s' in log entry: %s", field, line)
            }
        }
    }
}
```

---

## Migration Guide

### Step-by-Step Implementation

1. **Add Logging Package**
   ```bash
   mkdir -p internal/logging internal/headers
   # Copy the logging implementation files
   ```

2. **Update Dependencies**
   ```bash
   go get github.com/google/uuid@v1.4.0
   go mod tidy
   ```

3. **Replace Printf Statements**
   - Search for all `fmt.Printf` statements in kafka package
   - Replace with appropriate structured logging calls
   - Add correlation context where missing

4. **Update CLI Applications**
   - Add structured logging to main.go files
   - Implement graceful shutdown logging
   - Add configuration logging

5. **Test Integration**
   ```bash
   go test ./...
   docker-compose up --build
   ```

### Backward Compatibility

The implementation maintains backward compatibility:
- Existing `Send()` methods still work
- New `SendWithContext()` methods provide enhanced functionality
- Structured logs can be configured to output in different formats

---

## Benefits for Production

### 1. Operational Excellence
- **Structured logs** integrate with all modern log aggregation systems
- **Correlation tracking** enables distributed request tracing
- **Consistent formatting** across all services and languages

### 2. Debugging Efficiency
- **Find related logs** across producer/consumer/services using correlation ID
- **Filter by event types** to focus on specific operations
- **Performance analysis** with built-in timing and metrics

### 3. Monitoring Integration
- **Alert on error patterns** using structured fields
- **Dashboard creation** with consistent log structure
- **SLA monitoring** with processing duration tracking

### 4. Compliance & Auditing
- **Audit trails** with complete message lifecycle tracking
- **Data lineage** through correlation and trace IDs
- **Security monitoring** with user context tracking

This structured logging implementation transforms your Kafka demo from a basic example into a **production-ready, enterprise-grade application** that demonstrates modern observability practices essential for distributed systems.