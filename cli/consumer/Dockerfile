# ---- Build stage (glibc base) ----
FROM golang:1.25.7-trixie AS consumer_builder

RUN apt-get update && apt-get install -y librdkafka-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# CGO must be enabled for confluent-kafka-go
RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cli/consumer

# ---- Final image (Alpine base) ----
FROM alpine:3.23

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=consumer_builder /app/app .

CMD ["./app"]