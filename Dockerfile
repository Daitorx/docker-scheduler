# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o docker-scheduler .

# Runtime stage
FROM alpine:latest

# Install ca-certificates, docker-cli, and timezone data
RUN apk --no-cache add ca-certificates docker-cli tzdata

# Set timezone (can be overridden with -e TZ=...)
ENV TZ=Europe/Madrid

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/docker-scheduler .

# Copy static files
COPY --from=builder /app/static ./static

# Create data directory for SQLite
RUN mkdir -p /app/data

# Set Gin to release mode
ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Run the application
CMD ["./docker-scheduler"]
