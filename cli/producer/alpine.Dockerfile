FROM golang:1.24.6 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o producer ./cli/producer

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/producer .
CMD ["./producer"]

FROM golang:1.24.6-alpine3.22

# Required for confluent-kafka-go (via CGO)
RUN apk add --no-cache \
    gcc \
    musl-dev \
    librdkafka-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled
RUN go build -o app ./cli/producer

CMD ["./app"]