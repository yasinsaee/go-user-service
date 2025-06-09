# ===== Build Stage =====
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

# ===== Runtime Stage =====
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN chmod +x /app/server

ENTRYPOINT ["/app/server"]
