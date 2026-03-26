FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot_binary ./main.go

FROM alpine:3.19

RUN adduser -D appuser
USER appuser

WORKDIR /app
COPY --from=builder /app/bot_binary .

CMD ["./bot_binary"]