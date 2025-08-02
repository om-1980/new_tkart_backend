FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o product-service

EXPOSE 7000

CMD ["./product-service"]
