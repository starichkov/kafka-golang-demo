# Golang Kafka Demo

[![Author](https://img.shields.io/badge/Author-Vadim%20Starichkov-blue?style=for-the-badge)](https://github.com/starichkov)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/starichkov/kafka-golang-demo?style=for-the-badge)
[![GitHub License](https://img.shields.io/github/license/starichkov/kafka-golang-demo?style=for-the-badge)](https://github.com/starichkov/kafka-golang-demo/blob/main/LICENSE.md)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/starichkov/kafka-golang-demo/build.yml?style=for-the-badge)](https://github.com/starichkov/kafka-golang-demo/actions/workflows/build.yml)
[![Codecov](https://img.shields.io/codecov/c/github/starichkov/kafka-golang-demo?style=for-the-badge)](https://codecov.io/gh/starichkov/kafka-golang-demo)

This project demonstrates a basic Apache Kafka setup using Go (Golang) producers and consumers, fully containerized with Docker Compose. It uses the `confluent-kafka-go` client and runs Kafka in **KRaft mode**.

*This project has been created with assistance from JetBrains Junie to evaluate its capabilities and to speed up the development process.*

## 🚀 Key Features

- **Kafka 3.9.1 in KRaft mode** (no Zookeeper required).
- **Go Producer & Consumer** using `confluent-kafka-go` (v2).
- **Structured Logging** using `log/slog`.
- **Fully Containerized** with multi-stage Alpine-based Docker builds.
- **Resilient**: Kafka health checks with delayed startup for clients and auto-restart policies.
- **Automated Setup**: Auto topic creation via Kafka config and a setup container.
- **Testing**: Includes unit tests with mocks and integration tests using [Testcontainers](https://golang.testcontainers.org/).

## 📋 Requirements

- **Go**: 1.25.6 or later.
- **librdkafka**: Required for `confluent-kafka-go`. Must be installed on the host system for local builds.
- **Docker & Docker Compose**: Required for running the full stack and integration tests.

## 🏁 Getting Started

### Run with Docker Compose

The easiest way to run the entire stack:

```bash
docker compose up --build
```

This will:
1. Start Kafka in KRaft mode.
2. Wait until Kafka is healthy.
3. Create the required topics.
4. Run the Go producer (sends 10 messages and exits).
5. Run the Go consumer (continuously prints incoming messages).

### Local Execution

To build and run the components locally, ensure `librdkafka` is installed and Kafka is accessible.

**Build:**
```bash
CGO_ENABLED=1 go build -o producer ./cli/producer
CGO_ENABLED=1 go build -o consumer ./cli/consumer
```

**Run:**
```bash
# Ensure environment variables are set or use defaults
./producer
./consumer
```

## 🛠 Configuration

Configuration is handled via environment variables.

| Variable | Description | Default |
|----------|-------------|---------|
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka bootstrap address | `localhost:9092` |
| `TOPIC` | Kafka topic name | `demo-topic` |
| `GROUP_ID` | Consumer group ID | `golang-demo-group` |

## 🧪 Testing

### Unit Tests
Unit tests use a factory pattern and mocks for Kafka dependencies.

```bash
go test -v ./internal/kafka
```

### Integration Tests
Integration tests spin up a real Kafka instance in a Docker container.

```bash
go test -v integration_test.go
```
*Note: Docker must be running. The first run might take longer to pull the Kafka image (`confluentinc/cp-kafka:7.9.3`).*

For more testing options and coverage reports, see [TESTING.md](TESTING.md).

## 📚 Documentation

- [DOCUMENTATION.md](DOCUMENTATION.md): Guide on using Go's documentation tools (`go doc`, `godoc`).
- [TESTING.md](TESTING.md): Testing and coverage commands cheat sheet.

## 🧱 Project Structure

```
.
├── cli/
│   ├── producer/       # Go producer CLI entry point
│   └── consumer/       # Go consumer CLI entry point
├── internal/
│   ├── kafka/          # Kafka wrapper logic, interfaces, and mocks
│   └── logging/        # Structured JSON logging utilities
├── docker-compose.yml  # Kafka and applications orchestration
├── go.mod / go.sum     # Go modules and dependencies
├── integration_test.go # End-to-end integration tests
└── README.md
```

## 📦 Dependencies

- [confluent-kafka-go/v2](https://github.com/confluentinc/confluent-kafka-go): Kafka client for Go.
- [testcontainers-go](https://github.com/testcontainers/testcontainers-go): For integration testing.
- [Apache Kafka](https://kafka.apache.org/): Distributed event streaming platform.

## 📌 Notes

- **Alpine & CGO**: Uses `libc6-compat` in Docker images to support the C-linked `confluent-kafka-go`.
- **Startup**: Messages may show brief connection errors until the Kafka broker is fully ready.

---

## 🧾 About TemplateTasks

TemplateTasks is a personal software development initiative by Vadim Starichkov, focused on sharing open-source libraries, services, and technical demos. It operates independently and outside the scope of any employment.

## 📜 License & Attribution

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE.md) file for details.

### Using This Project?

If you use this code in your own projects, attribution is required under the MIT License:

```
Based on kafka-golang-demo by Vadim Starichkov, TemplateTasks
https://github.com/starichkov/kafka-golang-demo
```

**Copyright © 2026 Vadim Starichkov, TemplateTasks**
