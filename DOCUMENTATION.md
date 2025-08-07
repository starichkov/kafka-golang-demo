# Go Documentation Guide

This guide explains how to view and generate documentation for the kafka-golang-demo project using Go's built-in documentation tools.

## 📚 Available Documentation Tools

### 1. `go doc` Command-Line Tool

The `go doc` command extracts and displays documentation from Go source code.

#### View Package Documentation

```bash
# View main package documentation
go doc .

# View internal kafka package documentation
go doc ./internal/kafka

# View logging package documentation  
go doc ./internal/logging
```

#### View Type Documentation

```bash
# View Producer type documentation
go doc ./internal/kafka.Producer

# View Consumer type documentation
go doc ./internal/kafka.Consumer

# View all interfaces in kafka package
go doc ./internal/kafka.KafkaProducerInterface
go doc ./internal/kafka.KafkaConsumerInterface
```

#### View Method Documentation

```bash
# View specific Producer methods
go doc ./internal/kafka.Producer.Send
go doc ./internal/kafka.Producer.Close

# View specific Consumer methods
go doc ./internal/kafka.Consumer.Run
go doc ./internal/kafka.Consumer.RunWithChannel
go doc ./internal/kafka.Consumer.Close

# View constructor functions
go doc ./internal/kafka.NewProducer
go doc ./internal/kafka.NewConsumer
```

#### View Factory Documentation

```bash
# View factory interfaces and implementations
go doc ./internal/kafka.ProducerFactoryInterface
go doc ./internal/kafka.ConsumerFactoryInterface
go doc ./internal/kafka.DefaultProducerFactory
go doc ./internal/kafka.DefaultConsumerFactory
```

### 2. `godoc` Web Server

The `godoc` tool provides a web interface for browsing Go documentation, similar to what you see on pkg.go.dev.

#### Installation

```bash
# Install godoc tool
go install golang.org/x/tools/cmd/godoc@latest
```

#### Running the Web Server

```bash
# Start local documentation server on port 6060
godoc -http=:6060

# Start on a different port
godoc -http=:8080

# Start and open browser automatically (if available)
godoc -http=:6060 -play
```

Once running, open your browser to:
- **Local documentation**: http://localhost:6060/
- **This project**: http://localhost:6060/pkg/kafka-golang-demo/

#### Advanced godoc Options

```bash
# Include internal packages in documentation
godoc -http=:6060 -internal

# Enable the Go Playground
godoc -http=:6060 -play

# Use custom template directory
godoc -http=:6060 -templates=/path/to/templates
```

## 📖 Documentation Structure

### Package-Level Documentation

- **Main package**: Comprehensive overview in `doc.go`
- **kafka package**: Core Kafka abstractions and implementations
- **logging package**: Structured JSON logging utilities

### Key Documented Components

#### Kafka Package (`internal/kafka`)

- **Interfaces**: `KafkaProducerInterface`, `KafkaConsumerInterface`
- **Implementations**: `Producer`, `Consumer`
- **Factories**: `ProducerFactoryInterface`, `ConsumerFactoryInterface`
- **Mocks**: Available for testing (see `mocks.go`)

#### Logging Package (`internal/logging`)

- **Global Logger**: Structured JSON logging with slog
- **Configuration**: JSON output to stdout with Info level

## 🔍 Documentation Features

### Code Examples

The documentation includes practical examples:

```bash
# View usage examples in main package documentation
go doc .
```

### Cross-References

Documentation automatically links related types and functions within the same package.

### Formatted Output

The documentation supports:
- **Headers**: Using `#` syntax in comments
- **Code blocks**: Indented code in comments  
- **Lists**: Bullet points and numbered lists
- **Links**: Automatic linking between Go identifiers

## 🌐 Online Documentation

If this project is published as a Go module, it will automatically appear on:
- **pkg.go.dev**: https://pkg.go.dev/kafka-golang-demo
- **go.dev**: Searchable in the Go package index

## 🛠 Integration with IDEs

Most Go-aware IDEs support godoc integration:

- **VS Code**: Go extension shows documentation on hover
- **GoLand/IntelliJ**: Built-in documentation viewer
- **Vim/Neovim**: Various Go plugins support doc viewing

## 📋 Best Practices

### Writing Good Documentation

1. **Start with the name**: Begin comments with the name being documented
2. **Complete sentences**: Use proper grammar and punctuation
3. **Be concise**: Focus on what, not how
4. **Include examples**: Show typical usage patterns
5. **Document parameters**: Explain inputs and outputs

### Maintaining Documentation

1. **Keep it current**: Update docs when code changes
2. **Test examples**: Ensure code examples actually work
3. **Be consistent**: Follow Go documentation conventions
4. **Review regularly**: Documentation should be as good as the code

## 🔧 Troubleshooting

### Common Issues

**Problem**: `godoc` command not found
```bash
# Solution: Install the tool
go install golang.org/x/tools/cmd/godoc@latest
```

**Problem**: Documentation not showing up
```bash
# Solution: Ensure comments are directly above declarations with no blank lines
# ✅ Correct:
// Producer wraps a Kafka producer
type Producer struct {

# ❌ Incorrect:
// Producer wraps a Kafka producer

type Producer struct {
```

**Problem**: HTML templates not loading
```bash
# Solution: Use the -templates flag to specify template directory
godoc -http=:6060 -templates=$GOROOT/lib/godoc
```

---

For more information about Go documentation, see:
- [Go Doc Comments](https://tip.golang.org/doc/comment)
- [Effective Go - Commentary](https://go.dev/doc/effective_go#commentary)
- [pkg.go.dev About](https://pkg.go.dev/about)