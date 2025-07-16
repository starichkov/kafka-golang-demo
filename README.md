# Golang Kafka Demo

[![Author](https://img.shields.io/badge/Author-Vadim%20Starichkov-blue?style=for-the-badge)](https://github.com/starichkov)
[![GitHub License](https://img.shields.io/github/license/starichkov/kafka-golang-demo?style=for-the-badge)](https://github.com/starichkov/kafka-golang-demo/blob/main/LICENSE.md)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/starichkov/kafka-golang-demo/build.yml?style=for-the-badge)](https://github.com/starichkov/kafka-golang-demo/actions/workflows/build.yml)
[![Codecov](https://img.shields.io/codecov/c/github/starichkov/kafka-golang-demo?style=for-the-badge)](https://codecov.io/gh/starichkov/kafka-golang-demo)

This project demonstrates a basic Apache Kafka setup using Go (Golang) producers and consumers, fully containerized with
Docker Compose.

*This project has been created with assistance from JetBrains Junie to evaluate its capabilities and to speed up development process.*

## 👨‍💻 Author

**Vadim Starichkov** | [GitHub](https://github.com/starichkov) | [LinkedIn](https://www.linkedin.com/in/vadim-starichkov/)

*Developed to demonstrate modern Golang development practices and patterns, with Apache Kafka integration using containerized and testable architecture.*

## 🧱 Project Structure

```
.
├── cli/
│   ├── producer/       # Go CLI producer
│   └── consumer/       # Go CLI consumer
├── internal/kafka/     # Shared Kafka wrapper logic
├── docker-compose.yml  # Kafka + apps
├── go.mod / go.sum     # Go modules
└── README.md
```

## 🚀 Features

- Kafka 3.9.1 in **KRaft mode** (no Zookeeper)
- Go producer and consumer using `confluent-kafka-go`
- Multi-stage Alpine-based Docker builds
- Kafka health checks with delayed startup for clients
- Auto topic creation via Kafka config

## 🧪 Usage

### 1. Build and Run

```bash
docker compose up --build
```

This will:

- Start Kafka in KRaft mode
- Wait until Kafka is healthy
- Run the Go producer once (sends 10 messages)
- Run the Go consumer (prints incoming messages)

### 2. Kafka Admin (optional)

To list topics manually (inside the Kafka container):

```bash
docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## 🛠 Configuration

Set in `docker-compose.yml` as environment variables:

| Variable                  | Description                    |
|---------------------------|--------------------------------|
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka bootstrap address (host) |
| `TOPIC`                   | Kafka topic name               |
| `GROUP_ID`                | Consumer group ID (optional)   |

## 🧪 Integration Tests

This project includes comprehensive integration tests using [Testcontainers](https://golang.testcontainers.org/) that spin up a real Kafka instance to test the producer and consumer functionality.

### Running Integration Tests

```bash
go test -v ./integration_test.go
```

### What the Tests Cover

The integration tests include:

1. **Producer Test**: Verifies that the producer can successfully create and send messages to Kafka
2. **Consumer Test**: Verifies that the consumer can be created and configured properly
3. **Producer-Consumer Flow Test**: End-to-end test that:
   - Starts a real Kafka container using Testcontainers
   - Creates both producer and consumer
   - Sends multiple test messages
   - Verifies that messages are received by the consumer

### Test Configuration

- Uses `confluentinc/confluent-local:7.5.0` Kafka image (compatible with Testcontainers)
- Creates isolated test topics for each test run
- Automatically handles container lifecycle (start/stop/cleanup)
- Tests run with a 30-second timeout to ensure reliability

### Prerequisites for Testing

- Docker must be running (Testcontainers requires Docker)
- Go 1.24.5 or later
- Internet connection (to pull Kafka container image on first run)

## 📦 Dependencies

- Go 1.24.5
- [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)
- [librdkafka](https://github.com/confluentinc/librdkafka)
- Apache Kafka official Docker image
- [testcontainers-go](https://github.com/testcontainers/testcontainers-go) (for integration tests)

## 📌 Notes

- Docker Compose waits for Kafka to be fully ready using health checks.
- Messages may show connection errors briefly until Kafka becomes available.
- Uses `libc6-compat` to run glibc-linked Go binaries on Alpine.

---

## 🧾 About TemplateTasks

TemplateTasks is a personal software development initiative by Vadim Starichkov, focused on sharing open-source libraries, services, and technical demos.

It operates independently and outside the scope of any employment.

All code is released under permissive open-source licenses. The legal structure may evolve as the project grows.

## 📜 License & Attribution

This project is licensed under the **MIT License** - see
the [LICENSE](https://github.com/starichkov/kafka-golang-demo/blob/main/LICENSE.md) file for details.

### Using This Project?

If you use this code in your own projects, attribution is required under the MIT License:

```
Based on kafka-golang-demo by Vadim Starichkov, TemplateTasks

https://github.com/starichkov/kafka-golang-demo
```

**Copyright © 2025 Vadim Starichkov, TemplateTasks**
