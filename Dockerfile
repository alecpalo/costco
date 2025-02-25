FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o costco cmd/costco/main.go

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/costco .

EXPOSE 8080

CMD ["./costco"]