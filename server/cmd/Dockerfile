FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/pharmacy-server ./cmd/main.go


# Этап 2: Финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pharmacy-server .

RUN chmod +x /app/pharmacy-server

EXPOSE 8080

CMD ["./pharmacy-server"]

LABEL org.opencontainers.image.source=https://github.com/hnnsly/library-console
LABEL org.opencontainers.image.description="Backend API контейнер для консоли управления аптекой \"PharmShark\""
