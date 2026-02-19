# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Install git (required for fetching Go dependencies)
RUN apk add --no-cache git

WORKDIR /app

# 1. Copy everything currently in the folder (like main.go)
COPY . .

# 2. Initialize Go modules and fetch dependencies automatically
# This replaces the need for you to have 'go' installed on Windows
RUN go mod init app && \
    go get github.com/lib/pq && \
    go get github.com/streadway/amqp && \
    go mod tidy

# 3. Build the actual executable
RUN go build -o main .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Install curl (used for the healthcheck in docker-compose)
RUN apk add --no-cache curl

WORKDIR /root/

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Start the service
CMD ["./main"]