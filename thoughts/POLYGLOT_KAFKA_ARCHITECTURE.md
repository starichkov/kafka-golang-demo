# Polyglot Kafka Architecture - Cross-Language Communication Demo

**Created**: August 5, 2025  
**Purpose**: Design a multi-language Kafka communication system demonstrating enterprise-grade polyglot microservices architecture

## Executive Summary

This document outlines the architecture for a **polyglot microservices communication demo** using Kafka as the message bus between services built in **Java, Go, Python, and Node.js**. The goal is to demonstrate real-world, enterprise-level cross-language communication patterns using modern serialization formats.

**Key Insight**: This transforms the project from a simple Go+Kafka demo into a comprehensive **enterprise microservices architecture showcase** that demonstrates advanced distributed systems knowledge.

---

## Project Vision

### Multi-Language Service Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Service  │    │  Order Service  │    │Notification Svc │    │Analytics Service│
│   (Java/Spring) │    │      (Go)       │    │    (Python)     │    │    (Node.js)    │
│                 │    │                 │    │                 │    │                 │
│   REST API      │    │   REST API      │    │  Email/SMS      │    │   Dashboards    │
│   User mgmt     │    │   Order proc    │    │  Push notifs    │    │   Metrics       │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │                      │
          │ UserSignup           │ OrderCreated         │ EmailSent            │ AnalyticsUpdate
          │ UserUpdated          │ OrderShipped         │ SMSSent              │ MetricsCollected
          │                      │ PaymentProcessed     │                      │
          └──────────────────────┼──────────────────────┼──────────────────────┘
                                 │                      │
                        ┌────────▼──────────────────────▼────────┐
                        │              Apache Kafka             │
                        │                                       │
                        │  • Cross-language message bus        │
                        │  • Schema Registry (Avro/Protobuf)   │
                        │  • Event streaming & processing      │
                        │  • Reliable delivery guarantees      │
                        └───────────────────────────────────────┘
```

### Demo Flow Example

1. **Java Service** - User signup via REST API → publishes `UserSignupEvent`
2. **Go Service** - Consumes user events → creates welcome order → publishes `OrderCreatedEvent`  
3. **Python Service** - Consumes user/order events → sends email/SMS → publishes `NotificationSentEvent`
4. **Node.js Service** - Consumes all events → updates real-time analytics dashboard

**Result**: Single user action triggers a **cross-language event chain** demonstrating enterprise microservices communication.

---

## Serialization Format Analysis for Polyglot Systems

### Cross-Language Support Comparison

| Feature | JSON | Avro | Protocol Buffers |
|---------|------|------|------------------|
| **Java Support** | ✅ Native | ✅ Excellent | ✅ Excellent |
| **Go Support** | ✅ Native | ✅ Good | ✅ Excellent |
| **Python Support** | ✅ Native | ✅ Good | ✅ Excellent |
| **Node.js Support** | ✅ Native | ✅ Good | ✅ Excellent |
| **Type Safety** | ❌ Runtime | ⚠️ Runtime | ✅ Compile-time |
| **Schema Evolution** | ❌ Breaking | ✅ Excellent | ✅ Good |
| **Code Generation** | ❌ No | ❌ No | ✅ All languages |
| **IDE Support** | ⚠️ Basic | ⚠️ Limited | ✅ Full autocomplete |
| **Performance** | Slow | Fast | Fastest |

### Schema Management Approaches

#### **Avro: Runtime Schema Registry**
- **Central Schema Registry** - single source of truth for all languages
- **Dynamic schema resolution** - schemas fetched at runtime
- **Zero code generation** - schemas interpreted dynamically
- **Automatic compatibility checks** - registry validates schema evolution

**Pros for Polyglot:**
- ✅ No build coordination required
- ✅ Runtime schema evolution
- ✅ Confluent ecosystem integration

**Cons for Polyglot:**
- ❌ No compile-time validation
- ❌ Runtime schema resolution overhead
- ❌ Limited IDE support across languages

#### **Protocol Buffers: Build-time Code Generation**
- **Shared .proto files** - distributed across all projects  
- **Language-specific code generation** - compile schemas to native code
- **Strong typing** - compile-time validation in all languages
- **Manual coordination** - schema changes require rebuilds

**Pros for Polyglot:**
- ✅ Strong typing in all languages
- ✅ Compile-time validation prevents errors
- ✅ Excellent IDE support everywhere
- ✅ Best performance across all languages

**Cons for Polyglot:**
- ❌ Build coordination required
- ❌ Schema changes need synchronized deployments

---

## Recommended Architecture: Protocol Buffers First

### Why Protobuf for Polyglot Demo

**1. Superior Demo Experience**
- **Strong typing** shows consistency across all 4 languages
- **Compile-time validation** prevents runtime demo failures
- **IDE support** with auto-completion in Java, Go, Python, Node.js
- **Performance story** - measurable speed improvements

**2. Modern Microservices Pattern**
- **gRPC compatibility** - can extend demo to include gRPC services
- **Google ecosystem** - widely adopted in cloud-native applications
- **Developer familiarity** - more developers know Protobuf than Avro

**3. Easier Demo Setup**
- **No Schema Registry** required initially  
- **Build-time code generation** - simpler CI/CD
- **Self-contained** - each service has its schema code embedded

### Project Structure

```
kafka-polyglot-demo/
├── schemas/
│   └── proto/
│       ├── events/
│       │   ├── user_events.proto      # User lifecycle events
│       │   ├── order_events.proto     # Order processing events
│       │   ├── notification_events.proto # Communication events
│       │   └── analytics_events.proto # Metrics and tracking
│       ├── common/
│       │   ├── timestamps.proto       # Shared timestamp types
│       │   └── metadata.proto         # Common metadata structures
│       └── Makefile                   # Code generation for all languages
├── services/
│   ├── user-service-java/             # Spring Boot REST API
│   │   ├── src/main/java/
│   │   ├── src/main/proto/           # Generated protobuf classes
│   │   └── Dockerfile
│   ├── order-service-go/              # Go HTTP server  
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── proto/                     # Generated Go structs
│   │   └── Dockerfile
│   ├── notification-service-python/   # FastAPI service
│   │   ├── app/
│   │   ├── proto/                     # Generated Python classes
│   │   └── Dockerfile
│   ├── analytics-service-nodejs/      # Express + WebSocket server
│   │   ├── src/
│   │   ├── proto/                     # Generated JS classes
│   │   └── Dockerfile
│   └── docker-compose.yml
├── monitoring/
│   ├── grafana/
│   │   ├── dashboards/               # Cross-language metrics
│   │   └── provisioning/
│   ├── prometheus/
│   │   └── prometheus.yml
│   └── jaeger/                       # Distributed tracing
└── frontend/
    ├── dashboard/                    # Real-time demo dashboard
    └── admin/                        # Service management UI
```

---

## Schema Definitions

### Core Event Schemas

```protobuf
// schemas/proto/events/user_events.proto
syntax = "proto3";

package com.company.events;
option go_package = "kafka-polyglot-demo/proto/events";

import "google/protobuf/timestamp.proto";
import "common/metadata.proto";

message UserSignupEvent {
  string user_id = 1;
  string email = 2;
  string name = 3;
  google.protobuf.Timestamp signup_timestamp = 4;
  SignupSource source = 5;
  map<string, string> metadata = 6;
  common.EventMetadata event_metadata = 7;
}

message UserUpdatedEvent {
  string user_id = 1;
  string email = 2;
  string name = 3;
  google.protobuf.Timestamp updated_timestamp = 4;
  repeated string updated_fields = 5;
  common.EventMetadata event_metadata = 6;
}

enum SignupSource {
  SIGNUP_SOURCE_UNSPECIFIED = 0;
  SIGNUP_SOURCE_WEB = 1;
  SIGNUP_SOURCE_MOBILE = 2;
  SIGNUP_SOURCE_API = 3;
  SIGNUP_SOURCE_SOCIAL = 4;
}
```

```protobuf
// schemas/proto/events/order_events.proto
syntax = "proto3";

package com.company.events;
option go_package = "kafka-polyglot-demo/proto/events";

import "google/protobuf/timestamp.proto";
import "common/metadata.proto";

message OrderCreatedEvent {
  string order_id = 1;
  string user_id = 2;
  double total_amount = 3;
  string currency = 4;
  OrderStatus status = 5;
  repeated OrderItem items = 6;
  google.protobuf.Timestamp created_timestamp = 7;
  common.EventMetadata event_metadata = 8;
}

message OrderShippedEvent {
  string order_id = 1;
  string user_id = 2;
  string tracking_number = 3;
  string carrier = 4;
  google.protobuf.Timestamp shipped_timestamp = 5;
  common.EventMetadata event_metadata = 6;
}

message OrderItem {
  string product_id = 1;
  string product_name = 2;
  int32 quantity = 3;
  double unit_price = 4;
}

enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_PROCESSING = 3;
  ORDER_STATUS_SHIPPED = 4;
  ORDER_STATUS_DELIVERED = 5;
  ORDER_STATUS_CANCELLED = 6;
}
```

```protobuf
// schemas/proto/events/notification_events.proto
syntax = "proto3";

package com.company.events;
option go_package = "kafka-polyglot-demo/proto/events";

import "google/protobuf/timestamp.proto";
import "common/metadata.proto";

message EmailSentEvent {
  string user_id = 1;
  string email_address = 2;
  string template_id = 3;
  EmailType email_type = 4;
  map<string, string> template_variables = 5;
  google.protobuf.Timestamp sent_timestamp = 6;
  common.EventMetadata event_metadata = 7;
}

message SMSSentEvent {
  string user_id = 1;
  string phone_number = 2;
  string message_content = 3;
  SMSType sms_type = 4;
  google.protobuf.Timestamp sent_timestamp = 5;
  common.EventMetadata event_metadata = 6;
}

enum EmailType {
  EMAIL_TYPE_UNSPECIFIED = 0;
  EMAIL_TYPE_WELCOME = 1;
  EMAIL_TYPE_ORDER_CONFIRMATION = 2;
  EMAIL_TYPE_SHIPPING_NOTIFICATION = 3;
  EMAIL_TYPE_PROMOTIONAL = 4;
}

enum SMSType {
  SMS_TYPE_UNSPECIFIED = 0;
  SMS_TYPE_VERIFICATION = 1;
  SMS_TYPE_ORDER_UPDATE = 2;
  SMS_TYPE_PROMOTIONAL = 3;
}
```

```protobuf
// schemas/proto/common/metadata.proto
syntax = "proto3";

package common;
option go_package = "kafka-polyglot-demo/proto/common";

import "google/protobuf/timestamp.proto";

message EventMetadata {
  string event_id = 1;
  string correlation_id = 2;
  string causation_id = 3;
  string source_service = 4;
  string source_version = 5;
  google.protobuf.Timestamp timestamp = 6;
  map<string, string> headers = 7;
}
```

---

## Multi-Language Implementation Examples

### 1. Java Service (Spring Boot) - User Management

```java
// services/user-service-java/src/main/java/UserController.java
@RestController
@RequestMapping("/api/users")
public class UserController {
    
    @Autowired
    private UserEventProducer eventProducer;
    
    @PostMapping("/signup")
    public ResponseEntity<UserResponse> signup(@RequestBody SignupRequest request) {
        String userId = UUID.randomUUID().toString();
        
        // Create strongly-typed protobuf event
        UserSignupEvent event = UserSignupEvent.newBuilder()
            .setUserId(userId)
            .setEmail(request.getEmail())
            .setName(request.getName())
            .setSignupTimestamp(Timestamps.now())
            .setSource(SignupSource.SIGNUP_SOURCE_WEB)
            .putMetadata("ip_address", getClientIP())
            .putMetadata("user_agent", getUserAgent())
            .setEventMetadata(createEventMetadata("user-service", "1.0.0"))
            .build();
        
        // Publish to Kafka
        eventProducer.publishUserSignup(event);
        
        log.info("☕ Java: Published UserSignupEvent for user {}", userId);
        
        return ResponseEntity.ok(new UserResponse(userId, request.getEmail()));
    }
    
    private EventMetadata createEventMetadata(String service, String version) {
        return EventMetadata.newBuilder()
            .setEventId(UUID.randomUUID().toString())
            .setCorrelationId(MDC.get("correlation_id"))
            .setSourceService(service)
            .setSourceVersion(version)
            .setTimestamp(Timestamps.now())
            .build();
    }
}

@Component
public class UserEventProducer {
    
    @Autowired
    private KafkaTemplate<String, byte[]> kafkaTemplate;
    
    public void publishUserSignup(UserSignupEvent event) {
        byte[] serialized = event.toByteArray();
        
        kafkaTemplate.send("user-events", event.getUserId(), serialized)
            .addCallback(
                result -> log.info("✅ Published user signup: {} bytes", serialized.length),
                failure -> log.error("❌ Failed to publish user signup", failure)
            );
    }
}
```

### 2. Go Service (HTTP Server) - Order Processing

```go
// services/order-service-go/internal/handlers/order_handler.go
package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/google/uuid"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/timestamppb"
    
    pb "kafka-polyglot-demo/proto/events"
    "kafka-polyglot-demo/internal/kafka"
)

type OrderHandler struct {
    producer *kafka.EventProducer
}

type CreateOrderRequest struct {
    UserID string      `json:"user_id"`
    Items  []OrderItem `json:"items"`
}

type OrderItem struct {
    ProductID string  `json:"product_id"`
    Quantity  int32   `json:"quantity"`
    Price     float64 `json:"price"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    orderID := uuid.New().String()
    
    // Convert to protobuf items
    var pbItems []*pb.OrderItem
    var totalAmount float64
    
    for _, item := range req.Items {
        pbItems = append(pbItems, &pb.OrderItem{
            ProductId:   item.ProductID,
            ProductName: fmt.Sprintf("Product %s", item.ProductID), // In real app, fetch from DB
            Quantity:    item.Quantity,
            UnitPrice:   item.Price,
        })
        totalAmount += item.Price * float64(item.Quantity)
    }
    
    // Create strongly-typed protobuf event
    event := &pb.OrderCreatedEvent{
        OrderId:     orderID,
        UserId:      req.UserID,
        TotalAmount: totalAmount,
        Currency:    "USD",
        Status:      pb.OrderStatus_ORDER_STATUS_PENDING,
        Items:       pbItems,
        CreatedTimestamp: timestamppb.New(time.Now()),
        EventMetadata: &pb.EventMetadata{
            EventId:       uuid.New().String(),
            CorrelationId: getCorrelationID(r),
            SourceService: "order-service",
            SourceVersion: "1.0.0",
            Timestamp:     timestamppb.New(time.Now()),
        },
    }
    
    // Serialize and publish
    data, err := proto.Marshal(event)
    if err != nil {
        http.Error(w, "Serialization failed", http.StatusInternalServerError)
        return
    }
    
    if err := h.producer.PublishOrderEvent(r.Context(), "order-events", orderID, data); err != nil {
        http.Error(w, "Failed to publish event", http.StatusInternalServerError)
        return
    }
    
    log.Printf("🔥 Go: Published OrderCreatedEvent for order %s (%d bytes)", orderID, len(data))
    
    // Return response
    response := map[string]interface{}{
        "order_id": orderID,
        "status":   "created",
        "total":    totalAmount,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Consumer for user events (to create welcome orders)
func (h *OrderHandler) ProcessUserSignup(ctx context.Context, data []byte) error {
    var event pb.UserSignupEvent
    if err := proto.Unmarshal(data, &event); err != nil {
        return fmt.Errorf("failed to unmarshal user signup: %w", err)
    }
    
    log.Printf("🔥 Go: Processing signup for user %s, creating welcome order", event.UserId)
    
    // Create welcome order automatically
    welcomeOrder := &pb.OrderCreatedEvent{
        OrderId:     uuid.New().String(),
        UserId:      event.UserId,
        TotalAmount: 0.00, // Free welcome package
        Currency:    "USD",
        Status:      pb.OrderStatus_ORDER_STATUS_CONFIRMED,
        Items: []*pb.OrderItem{
            {
                ProductId:   "welcome-package",
                ProductName: "Welcome Package",
                Quantity:    1,
                UnitPrice:   0.00,
            },
        },
        CreatedTimestamp: timestamppb.New(time.Now()),
        EventMetadata: &pb.EventMetadata{
            EventId:       uuid.New().String(),
            CorrelationId: event.EventMetadata.CorrelationId, // Maintain correlation
            CausationId:   event.EventMetadata.EventId,       // Track causation
            SourceService: "order-service",
            SourceVersion: "1.0.0",
            Timestamp:     timestamppb.New(time.Now()),
        },
    }
    
    data, _ := proto.Marshal(welcomeOrder)
    return h.producer.PublishOrderEvent(ctx, "order-events", welcomeOrder.UserId, data)
}
```

### 3. Python Service (FastAPI) - Notification Service

```python
# services/notification-service-python/app/notification_service.py
from fastapi import FastAPI, BackgroundTasks
from confluent_kafka import Consumer, KafkaError
import asyncio
import logging
from typing import Dict, Any

# Generated protobuf classes
from proto.events import user_events_pb2, order_events_pb2, notification_events_pb2
from proto.common import metadata_pb2

class NotificationService:
    def __init__(self):
        self.app = FastAPI()
        self.consumer = Consumer({
            'bootstrap.servers': 'kafka:9092',
            'group.id': 'notification-service',
            'auto.offset.reset': 'earliest'
        })
        
        # Subscribe to multiple topics
        self.consumer.subscribe(['user-events', 'order-events'])
        
        self.email_client = EmailClient()
        self.sms_client = SMSClient()
        
    async def start_consuming(self):
        """Start consuming events from Kafka"""
        while True:
            msg = self.consumer.poll(timeout=1.0)
            
            if msg is None:
                continue
                
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    continue
                else:
                    logging.error(f"Consumer error: {msg.error()}")
                    continue
            
            # Route message based on topic
            topic = msg.topic()
            
            try:
                if topic == 'user-events':
                    await self.process_user_event(msg.value())
                elif topic == 'order-events':
                    await self.process_order_event(msg.value())
                    
            except Exception as e:
                logging.error(f"Error processing message from {topic}: {e}")
    
    async def process_user_event(self, message_data: bytes):
        """Process user lifecycle events"""
        # Try to parse as UserSignupEvent
        try:
            event = user_events_pb2.UserSignupEvent()
            event.ParseFromString(message_data)
            
            logging.info(f"🐍 Python: Processing UserSignupEvent for {event.user_id}")
            
            # Send welcome email
            await self.send_welcome_email(event)
            
        except Exception as e:
            # Try other user event types
            logging.debug(f"Not a UserSignupEvent: {e}")
    
    async def process_order_event(self, message_data: bytes):
        """Process order lifecycle events"""
        try:
            event = order_events_pb2.OrderCreatedEvent()
            event.ParseFromString(message_data)
            
            logging.info(f"🐍 Python: Processing OrderCreatedEvent {event.order_id} for user {event.user_id}")
            
            # Send order confirmation
            await self.send_order_confirmation(event)
            
        except Exception as e:
            # Try OrderShippedEvent
            try:
                shipped_event = order_events_pb2.OrderShippedEvent()
                shipped_event.ParseFromString(message_data)
                
                logging.info(f"🐍 Python: Processing OrderShippedEvent {shipped_event.order_id}")
                await self.send_shipping_notification(shipped_event)
                
            except Exception as e2:
                logging.debug(f"Unknown order event type: {e2}")
    
    async def send_welcome_email(self, event: user_events_pb2.UserSignupEvent):
        """Send welcome email to new user"""
        template_vars = {
            'user_name': event.name,
            'signup_source': event.source.name,
        }
        
        success = await self.email_client.send_email(
            to_email=event.email,
            template='welcome',
            variables=template_vars
        )
        
        if success:
            # Publish EmailSentEvent
            email_event = notification_events_pb2.EmailSentEvent(
                user_id=event.user_id,
                email_address=event.email,
                template_id='welcome',
                email_type=notification_events_pb2.EmailType.EMAIL_TYPE_WELCOME,
                template_variables=template_vars,
                sent_timestamp=self.current_timestamp(),
                event_metadata=self.create_event_metadata(
                    correlation_id=event.event_metadata.correlation_id
                )
            )
            
            await self.publish_notification_event(email_event)
            logging.info(f"📧 Sent welcome email to {event.email}")
    
    async def send_order_confirmation(self, event: order_events_pb2.OrderCreatedEvent):
        """Send order confirmation email"""
        # Calculate order summary
        item_summary = []
        for item in event.items:
            item_summary.append(f"{item.quantity}x {item.product_name} - ${item.unit_price:.2f}")
        
        template_vars = {
            'order_id': event.order_id,
            'total_amount': f"${event.total_amount:.2f}",
            'items': ', '.join(item_summary),
            'order_status': event.status.name
        }
        
        # For demo: get user email (in real app, query user service)
        user_email = f"user_{event.user_id}@example.com"
        
        success = await self.email_client.send_email(
            to_email=user_email,
            template='order_confirmation', 
            variables=template_vars
        )
        
        if success:
            email_event = notification_events_pb2.EmailSentEvent(
                user_id=event.user_id,
                email_address=user_email,
                template_id='order_confirmation',
                email_type=notification_events_pb2.EmailType.EMAIL_TYPE_ORDER_CONFIRMATION,
                template_variables=template_vars,
                sent_timestamp=self.current_timestamp(),
                event_metadata=self.create_event_metadata(
                    correlation_id=event.event_metadata.correlation_id
                )
            )
            
            await self.publish_notification_event(email_event)
            logging.info(f"📧 Sent order confirmation for {event.order_id}")
    
    def create_event_metadata(self, correlation_id: str = None):
        """Create event metadata for published events"""
        import uuid
        from google.protobuf.timestamp_pb2 import Timestamp
        
        metadata = metadata_pb2.EventMetadata()
        metadata.event_id = str(uuid.uuid4())
        metadata.correlation_id = correlation_id or str(uuid.uuid4())
        metadata.source_service = "notification-service"
        metadata.source_version = "1.0.0"
        metadata.timestamp.CopyFrom(self.current_timestamp())
        
        return metadata
    
    def current_timestamp(self):
        """Get current timestamp as protobuf Timestamp"""
        from google.protobuf.timestamp_pb2 import Timestamp
        timestamp = Timestamp()
        timestamp.GetCurrentTime()
        return timestamp

# Mock email client for demo
class EmailClient:
    async def send_email(self, to_email: str, template: str, variables: Dict[str, Any]) -> bool:
        # Simulate email sending
        await asyncio.sleep(0.1)  # Simulate network delay
        logging.info(f"📧 Mock: Sending {template} email to {to_email}")
        return True
```

### 4. Node.js Service (Express) - Analytics Dashboard

```javascript
// services/analytics-service-nodejs/src/analytics_service.js
const express = require('express');
const { Kafka } = require('kafkajs');
const WebSocket = require('ws');

// Generated protobuf classes
const { UserSignupEvent, UserUpdatedEvent } = require('./proto/user_events_pb');
const { OrderCreatedEvent, OrderShippedEvent } = require('./proto/order_events_pb');
const { EmailSentEvent, SMSSentEvent } = require('./proto/notification_events_pb');

class AnalyticsService {
    constructor() {
        this.app = express();
        this.wss = new WebSocket.Server({ port: 8080 });
        
        // Kafka setup
        this.kafka = Kafka({
            clientId: 'analytics-service',
            brokers: ['kafka:9092']
        });
        
        this.consumer = this.kafka.consumer({ groupId: 'analytics-service' });
        
        // In-memory analytics storage (use Redis/DB in production)
        this.metrics = {
            userSignups: 0,
            ordersCreated: 0,
            emailsSent: 0,
            totalRevenue: 0,
            signupsBySource: {},
            ordersByHour: {},
            recentEvents: []
        };
        
        this.setupRoutes();
        this.setupWebSocket();
    }
    
    setupRoutes() {
        this.app.get('/metrics', (req, res) => {
            res.json(this.metrics);
        });
        
        this.app.get('/dashboard', (req, res) => {
            res.sendFile(__dirname + '/dashboard.html');
        });
    }
    
    setupWebSocket() {
        this.wss.on('connection', (ws) => {
            console.log('🟨 Node.js: Dashboard client connected');
            
            // Send current metrics to new client
            ws.send(JSON.stringify({
                type: 'metrics_update',
                data: this.metrics
            }));
        });
    }
    
    async startConsuming() {
        await this.consumer.connect();
        await this.consumer.subscribe({ 
            topics: ['user-events', 'order-events', 'notification-events'] 
        });
        
        await this.consumer.run({
            eachMessage: async ({ topic, partition, message }) => {
                try {
                    await this.processMessage(topic, message.value);
                } catch (error) {
                    console.error(`🟨 Node.js: Error processing message from ${topic}:`, error);
                }
            }
        });
    }
    
    async processMessage(topic, messageData) {
        let event = null;
        let eventType = '';
        
        try {
            switch (topic) {
                case 'user-events':
                    event = this.parseUserEvent(messageData);
                    break;
                case 'order-events':
                    event = this.parseOrderEvent(messageData);
                    break;
                case 'notification-events':
                    event = this.parseNotificationEvent(messageData);
                    break;
            }
            
            if (event) {
                await this.updateMetrics(event);
                this.broadcastUpdate(event);
            }
            
        } catch (error) {
            console.error(`🟨 Node.js: Failed to parse ${topic} message:`, error);
        }
    }
    
    parseUserEvent(messageData) {
        // Try UserSignupEvent first
        try {
            const event = UserSignupEvent.deserializeBinary(messageData);
            console.log(`🟨 Node.js: Processing UserSignupEvent for ${event.getUserId()}`);
            
            return {
                type: 'user_signup',
                userId: event.getUserId(),
                email: event.getEmail(),
                source: event.getSource(),
                timestamp: event.getSignupTimestamp().toDate()
            };
        } catch (e) {
            // Try other user event types...
            return null;
        }
    }
    
    parseOrderEvent(messageData) {
        // Try OrderCreatedEvent
        try {
            const event = OrderCreatedEvent.deserializeBinary(messageData);
            console.log(`🟨 Node.js: Processing OrderCreatedEvent ${event.getOrderId()}`);
            
            return {
                type: 'order_created',
                orderId: event.getOrderId(),
                userId: event.getUserId(),
                totalAmount: event.getTotalAmount(),
                itemCount: event.getItemsList().length,
                timestamp: event.getCreatedTimestamp().toDate()
            };
        } catch (e) {
            // Try OrderShippedEvent
            try {
                const event = OrderShippedEvent.deserializeBinary(messageData);
                console.log(`🟨 Node.js: Processing OrderShippedEvent ${event.getOrderId()}`);
                
                return {
                    type: 'order_shipped',
                    orderId: event.getOrderId(),
                    trackingNumber: event.getTrackingNumber(),
                    timestamp: event.getShippedTimestamp().toDate()
                };
            } catch (e2) {
                return null;
            }
        }
    }
    
    parseNotificationEvent(messageData) {
        // Try EmailSentEvent
        try {
            const event = EmailSentEvent.deserializeBinary(messageData);
            console.log(`🟨 Node.js: Processing EmailSentEvent for ${event.getUserId()}`);
            
            return {
                type: 'email_sent',
                userId: event.getUserId(),
                emailType: event.getEmailType(),
                templateId: event.getTemplateId(),
                timestamp: event.getSentTimestamp().toDate()
            };
        } catch (e) {
            return null;
        }
    }
    
    async updateMetrics(event) {
        const hour = new Date().getHours();
        
        switch (event.type) {
            case 'user_signup':
                this.metrics.userSignups++;
                this.metrics.signupsBySource[event.source] = 
                    (this.metrics.signupsBySource[event.source] || 0) + 1;
                break;
                
            case 'order_created':
                this.metrics.ordersCreated++;
                this.metrics.totalRevenue += event.totalAmount;
                this.metrics.ordersByHour[hour] = 
                    (this.metrics.ordersByHour[hour] || 0) + 1;
                break;
                
            case 'email_sent':
                this.metrics.emailsSent++;
                break;
        }
        
        // Keep recent events for activity feed
        this.metrics.recentEvents.unshift({
            ...event,
            timestamp: new Date().toISOString()
        });
        
        // Keep only last 50 events
        if (this.metrics.recentEvents.length > 50) {
            this.metrics.recentEvents = this.metrics.recentEvents.slice(0, 50);
        }
    }
    
    broadcastUpdate(event) {
        const update = {
            type: 'event_update',
            event: event,
            metrics: this.metrics
        };
        
        // Broadcast to all connected WebSocket clients
        this.wss.clients.forEach(client => {
            if (client.readyState === WebSocket.OPEN) {
                client.send(JSON.stringify(update));
            }
        });
    }
    
    start() {
        const port = process.env.PORT || 3000;
        this.app.listen(port, () => {
            console.log(`🟨 Node.js: Analytics service running on port ${port}`);
        });
        
        this.startConsuming();
    }
}

module.exports = AnalyticsService;

// Start the service
if (require.main === module) {
    const service = new AnalyticsService();
    service.start();
}
```

---

## Code Generation & Build Process

### Makefile for Cross-Language Generation

```makefile
# schemas/Makefile
.PHONY: generate-all generate-java generate-go generate-python generate-nodejs clean

# Generate for all languages
generate-all: generate-java generate-go generate-python generate-nodejs

# Java generation
generate-java:
	@echo "Generating Java protobuf classes..."
	protoc --java_out=../services/user-service-java/src/main/java \
	       --proto_path=. \
	       proto/events/*.proto proto/common/*.proto

# Go generation  
generate-go:
	@echo "Generating Go protobuf structs..."
	protoc --go_out=../services/order-service-go \
	       --go_opt=paths=source_relative \
	       --proto_path=. \
	       proto/events/*.proto proto/common/*.proto

# Python generation
generate-python:
	@echo "Generating Python protobuf classes..."
	protoc --python_out=../services/notification-service-python \
	       --proto_path=. \
	       proto/events/*.proto proto/common/*.proto

# Node.js generation
generate-nodejs:
	@echo "Generating Node.js protobuf classes..."
	protoc --js_out=import_style=commonjs:../services/analytics-service-nodejs/src \
	       --proto_path=. \
	       proto/events/*.proto proto/common/*.proto

# Clean generated files
clean:
	find ../services -name "*_pb2.py" -delete
	find ../services -name "*_pb.js" -delete
	find ../services -name "*.pb.go" -delete
	find ../services -path "*/generated/*" -delete

# Validate schemas
validate:
	@echo "Validating protobuf schemas..."
	protoc --proto_path=. --descriptor_set_out=/dev/null proto/events/*.proto proto/common/*.proto
	@echo "✅ All schemas are valid"

# Generate documentation
docs:
	protoc --doc_out=../docs --doc_opt=html,index.html proto/events/*.proto proto/common/*.proto
```

### Docker Compose for Full Demo

```yaml
# services/docker-compose.yml
version: '3.8'

services:
  # Infrastructure
  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    depends_on: [zookeeper]
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"

  # Schema Registry (for Avro comparison)
  schema-registry:
    image: confluentinc/cp-schema-registry:7.4.0
    depends_on: [kafka]
    ports:
      - "8081:8081"
    environment:
      SCHEMA_REGISTRY_HOST_NAME: schema-registry
      SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:9092

  # Monitoring
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ../monitoring/prometheus:/etc/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    volumes:
      - ../monitoring/grafana:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

  # Microservices
  user-service:
    build: ./user-service-java
    ports:
      - "8080:8080"
    depends_on: [kafka]
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
      - SCHEMA_REGISTRY_URL=http://schema-registry:8081
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  order-service:
    build: ./order-service-go
    ports:
      - "8081:8081"
    depends_on: [kafka, user-service]
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
      - USER_SERVICE_URL=http://user-service:8080
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  notification-service:
    build: ./notification-service-python
    depends_on: [kafka]
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
      - EMAIL_API_KEY=demo_key
      - SMS_API_KEY=demo_key
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  analytics-service:
    build: ./analytics-service-nodejs
    ports:
      - "3000:3000"
      - "8080:8080"  # WebSocket port
    depends_on: [kafka]
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/metrics"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Demo dashboard
  demo-dashboard:
    build: ../frontend/dashboard
    ports:
      - "3002:80"
    depends_on: [analytics-service]
    environment:
      - ANALYTICS_WS_URL=ws://analytics-service:8080
```

---

## Demo Execution Flow

### 1. Setup and Start
```bash
# Generate protobuf code for all languages
cd schemas && make generate-all

# Start all services
cd services && docker-compose up --build

# Wait for services to be healthy
docker-compose ps
```

### 2. Demo Script
```bash
#!/bin/bash
# demo_script.sh

echo "🚀 Starting Polyglot Kafka Demo"
echo "================================="

echo "📊 Opening analytics dashboard..."
open http://localhost:3002/dashboard

echo "☕ Step 1: User signup via Java service"
curl -X POST http://localhost:8080/api/users/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "demo@example.com", "name": "Demo User"}'

echo ""
echo "⏱️  Waiting for event propagation..."
sleep 2

echo "🔥 Step 2: Manual order via Go service"
curl -X POST http://localhost:8081/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "items": [
      {"product_id": "prod_1", "quantity": 2, "price": 29.99},
      {"product_id": "prod_2", "quantity": 1, "price": 49.99}
    ]
  }'

echo ""
echo "⏱️  Waiting for notifications..."
sleep 3

echo "📈 Step 3: Check analytics"
curl http://localhost:3000/metrics | jq

echo ""
echo "✅ Demo completed! Check the dashboard for real-time updates"
echo "   Dashboard: http://localhost:3002"
echo "   Grafana:   http://localhost:3001 (admin/admin)"
echo "   Kafka UI:  http://localhost:9021"
```

### 3. Expected Output
```
🚀 Java: Published UserSignupEvent for user abc-123
🔥 Go: Processing signup for user abc-123, creating welcome order
🔥 Go: Published OrderCreatedEvent for order def-456
🐍 Python: Processing UserSignupEvent for abc-123
📧 Sent welcome email to demo@example.com
🐍 Python: Processing OrderCreatedEvent def-456
📧 Sent order confirmation for def-456
🟨 Node.js: Processing UserSignupEvent for abc-123
🟨 Node.js: Processing OrderCreatedEvent def-456
🟨 Node.js: Processing EmailSentEvent for abc-123

📊 Analytics Update:
- User Signups: 1
- Orders Created: 2 (1 welcome + 1 manual)
- Emails Sent: 2
- Total Revenue: $109.97
```

---

## Why This Approach Wins

### 1. Enterprise-Grade Demonstration
- **Real-world architecture** - polyglot microservices communicating via Kafka
- **Production patterns** - event sourcing, CQRS, saga patterns
- **Scalability showcase** - horizontal scaling, load balancing

### 2. Technical Depth
- **Cross-language serialization** - Protocol Buffers expertise
- **Event-driven architecture** - asynchronous communication patterns  
- **Distributed systems** - correlation IDs, distributed tracing
- **Monitoring & observability** - metrics, logs, traces

### 3. Professional Impact
- **Portfolio differentiator** - beyond typical CRUD applications
- **Architecture knowledge** - microservices, event streaming
- **Multi-language skills** - Java, Go, Python, Node.js proficiency
- **DevOps practices** - containerization, orchestration, monitoring

### 4. Extensibility
- **Add gRPC** - synchronous service-to-service communication
- **Add databases** - event sourcing with PostgreSQL/MongoDB
- **Add caching** - Redis for session management
- **Add API gateway** - Kong/Envoy for service mesh

This transforms your simple Kafka demo into a **comprehensive distributed systems showcase** that demonstrates enterprise-level architecture knowledge and cross-language communication expertise!