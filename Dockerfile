# Use the official Go image as the builder
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the Go server
RUN go build -tags netgo -ldflags "-s -w" -o server ./cmd/server/server.go

# Use a lightweight image for the final container
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder
COPY --from=builder /app/server .

# Expose the port (Railway can map this dynamically)
EXPOSE 42069

# Run the server
CMD ["./server"]
