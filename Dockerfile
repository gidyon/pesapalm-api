# Stage 1: Build the Go binary
FROM golang:1.23 as builder

# Set environment variables for Go cross-compilation
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GO111MODULE=on

# Create an app directory
WORKDIR /app

# Cache dependencies by copying go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary, specifying the main file in cmd/
RUN go build -v -o service -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" ./cmd/app/*.go

# Stage 2: Run the binary in a minimal container
FROM alpine:3.18

# Set up working directory
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/service .

# Set the app environment variable path if needed
ENV CONFIG_FILE_PATH=/root/app.env

# Ensure the binary has execution permission
RUN chmod +x ./service

# Expose the port (update if your app uses a different port)
EXPOSE 8080
EXPOSE 9090

# Command to run the service
ENTRYPOINT ["./service"]
CMD ["--config-file","/root/app.env"]
