# Build Stage
FROM golang:1.26.1-alpine AS builder

# Install build dependencies for CGO (required for sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for the sqlite driver we're using
RUN CGO_ENABLED=1 GOOS=linux go build -o server-watcher-app main.go

# Runner Stage
FROM alpine:latest

# Install runtime dependencies for sqlite, nsenter and health checks
RUN apk add --no-cache ca-certificates tzdata util-linux

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server-watcher-app .

# Ensure the database file exists or is accessible
# We'll handle the DB file as a volume in docker-compose for persistence

EXPOSE 5001

CMD ["./server-watcher-app"]
