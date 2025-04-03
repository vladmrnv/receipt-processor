FROM golang:1.23.4-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o receipt-processor .

# Create a minimal image
FROM alpine:latest

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/receipt-processor /receipt-processor

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/receipt-processor"]