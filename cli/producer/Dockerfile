# ---- Build stage (glibc base) ----
FROM golang:1.24.5-bookworm AS producer_builder

RUN apt-get update && apt-get install -y librdkafka-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# CGO must be enabled for confluent-kafka-go
RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cli/producer

# ---- Final image (Alpine base) ----
FROM alpine:3.22.0

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=producer_builder /app/app .

CMD ["./app"]
