FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/app
RUN go build -o publisher ./cmd/publisher

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/publisher .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x /root/app /root/publisher

EXPOSE 8080

CMD ["./app"]