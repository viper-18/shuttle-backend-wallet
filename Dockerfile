# Use the official Golang image to build the Go application
FROM golang:1.21.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image to run the Go application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install required packages
RUN apk add --no-cache ca-certificates

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
