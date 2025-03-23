# Stage 1: Build the Go application
FROM golang:1.21 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o service-availability main.go

# Stage 2: Run the application using a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/service-availability .

# Expose the port your service runs on
EXPOSE 8080

# Pass environment variables and start the service
ENTRYPOINT ["./service-availability"]
