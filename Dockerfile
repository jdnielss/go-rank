# Use the official golang image as a base for the Go build stage
FROM golang:1.22-alpine as builder
RUN apk add --no-cache gcc g++ make

# Create the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o /code-runner

# Use the Eclipse Temurin image as a base for the runtime stage
FROM eclipse-temurin:22-jdk-alpine

# Install Node.js, Go, and Python
RUN apk add --no-cache nodejs npm go python3

# Copy the built Go application from the builder stage
COPY --from=builder /code-runner /code-runner

# Set the working directory
WORKDIR /app

# Set the PATH environment variable to include the directories of Node.js, Go, and Java binaries
ENV PATH="/usr/local/go/bin:/usr/local/bin:${PATH}"

# Run the application
ENTRYPOINT ["/code-runner"]