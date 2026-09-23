# Etapa 1: compilar
FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

RUN go build -o servidor main.go

# Etapa 2: imagen final liviana
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/servidor .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./servidor"]