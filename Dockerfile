# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/config ./config

# Copy locales directory
COPY --from=builder /app/locales ./locales

# Copy migrations directory
COPY --from=builder /app/migrations ./migrations

# Copy scripts directory
COPY --from=builder /app/scripts ./scripts

# Install PostgreSQL client for migrations
RUN apk --no-cache add postgresql-client

# Make the migration script executable
RUN chmod +x ./scripts/run-migrations.sh

# Expose the application port
EXPOSE 8003

# Run the application
CMD ["./main"]
