# Use an official Golang runtime as the base image
FROM golang:1.20-alpine

# Install Redis server
RUN apk add --no-cache redis

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the source code to the container
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go server executable
RUN go build -o server .

# Expose the port for the Go server
EXPOSE 8080

# Start Redis server and the Go server
CMD ["sh", "-c", "redis-server --daemonize yes && ./server"]