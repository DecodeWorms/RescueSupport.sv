# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage

# Create a work directory inside Docker
WORKDIR /app

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all Go source files and other necessary files
COPY . .

# Set environment variables for Go build
ENV CGO_ENABLED=0
ENV GOOS=linux

# Build the Go application
RUN go build -o rescue-support.sv .

# Run the tests in the container
FROM build-stage AS run-test-stage

# Install docker-cli to access docker within the container
RUN apt-get update && apt-get install -y docker.io

# Bind-mount the Docker socket (when running the container)
CMD ["go", "test", "-v", "./..."]
