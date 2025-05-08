# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Stage 2: Run the binary in a minimal image
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Expose the port your app runs on (e.g., 8080)
EXPOSE 8080

# Run the binary
CMD ["./server"]
