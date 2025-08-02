FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o order-service

EXPOSE 9000

CMD ["./order-service"]
