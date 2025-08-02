# Stage 1 - Build
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o delivery-service ./main.go

# Stage 2 - Run
FROM debian:bullseye-slim
WORKDIR /
COPY --from=builder /app/delivery-service .

EXPOSE 8003
ENTRYPOINT ["./delivery-service"]
