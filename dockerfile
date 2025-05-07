# Use Ubuntu as the base image
FROM ubuntu:22.04

# Install Go 1.24 (official binary) and dependencies
RUN apt-get update && apt-get install -y curl git ca-certificates && \
    curl -OL https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"

# Create and set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go API
RUN go build -o main .

# Expose port (change as needed)
EXPOSE 8080

# Set environment file (runtime)
# Use docker run --env-file to inject it at runtime
CMD ["./main"]