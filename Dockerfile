# ===== Build Stage =====
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o server ./cmd

# ===== Runtime Stage =====
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache openssl bash

COPY --from=builder /app/server /app/server

RUN mkdir -p /app/keys && touch /app/.env

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/server /app/entrypoint.sh

ENV PRIVATE_KEY_PATH=/app/keys/private.key
ENV PUBLIC_KEY_PATH=/app/keys/public.key
ENV ENV_PATH=/app/.env

ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["/app/server"]
