# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy go mod files for worker and shared
COPY modules/worker/go.mod modules/worker/go.sum ./
COPY shared/go.mod ../../shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY modules/worker/ .
COPY shared/ ../../shared/

# Build the worker application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/worker .

# Copy email templates
COPY --from=builder /app/templates ./templates

# Run the worker
CMD ["./worker"]
