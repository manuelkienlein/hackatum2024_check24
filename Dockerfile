# Use the official Golang image as a base image
FROM golang:1.22.2

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /app

# Copy Go modules and application source code
COPY . .
RUN go mod download

# Build the Go application
RUN go build -o main ./cmd

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./main"]