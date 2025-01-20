# Use the official Go image as the base image (for building the app)
FROM golang:1.23.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o server .

# Use a smaller base image for the final image
FROM alpine:latest

# Install necessary dependencies (ca-certificates)
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go binary from the builder image
COPY --from=builder /app/server .

# Copy the .env file into the container
COPY .env .env

# Expose the port the server runs on
EXPOSE 8080

# Command to run the Go application (reads port from the environment)
CMD ["sh", "-c", "./server --port=${PORT:-8080}"]
