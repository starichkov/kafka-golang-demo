#FROM golang:1.24.4 AS builder
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -o producer ./cli/producer
#
#FROM alpine:3.22.0
#WORKDIR /app
#COPY --from=builder /app/producer .
#CMD ["./producer"]

#FROM golang:1.24.4-alpine3.22
#
## Required for confluent-kafka-go (via CGO)
#RUN apk add --no-cache \
#    gcc \
#    musl-dev \
#    librdkafka-dev
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#RUN go mod download
#
#COPY . .
#
## Build with CGO enabled
#RUN go build -o app ./cli/producer
#
#CMD ["./app"]

# ---- Build stage (glibc base) ----
FROM golang:1.24.4-bookworm AS producer_builder

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
