# Use official Go image
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN go build -o auth-service

# Expose the port
EXPOSE 8000

# Set environment variables file
ENV GO_ENV=production

# Run the executable
CMD ["./auth-service"]
