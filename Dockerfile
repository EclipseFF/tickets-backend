# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /build

# Copy go.mod and go.sum files to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o /main ./cmd/api/

# Stage 2: Create a lightweight container
FROM alpine:3

# Copy the Go binary from the builder stage
COPY --from=builder /main /bin/main

# Define the entry point for the container
ENTRYPOINT ["/bin/main"]
