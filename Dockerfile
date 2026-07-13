# Build stage
FROM golang:1.26.3 AS builder

WORKDIR /app

# Copy go modules
COPY go.mod go.sum ./

COPY internal/vendor/onvif-go ./internal/vendor/onvif-go

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o onvif-controller ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install certificates
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/onvif-controller .

# Expose API port
EXPOSE 8090

# Start application
CMD ["./onvif-controller"]