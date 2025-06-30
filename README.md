# Golang Kafka Demo

[![Author](https://img.shields.io/badge/Author-Vadim%20Starichkov-blue?style=for-the-badge)](https://github.com/starichkov)
[![GitHub License](https://img.shields.io/github/license/starichkov/kafka-golang-demo?style=for-the-badge)](https://github.com/starichkov/kafka-golang-demo/blob/main/LICENSE.md)

This project demonstrates a basic Apache Kafka setup using Go (Golang) producers and consumers, fully containerized with
Docker Compose.

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

## 📦 Dependencies

- Go 1.24.4
- [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)
- [librdkafka](https://github.com/confluentinc/librdkafka)
- Apache Kafka official Docker image

## 📌 Notes

- Docker Compose waits for Kafka to be fully ready using health checks.
- Messages may show connection errors briefly until Kafka becomes available.
- Uses `libc6-compat` to run glibc-linked Go binaries on Alpine.

---

## 🧾 About TemplateTasks

TemplateTasks is a developer-focused initiative by Vadim Starichkov, currently operated as sole proprietorship in
Finland.  
All code is released under open-source licenses. Ownership may be transferred to a registered business entity in the
future.

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
