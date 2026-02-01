# Base stage
FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# Development stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
EXPOSE 8080
CMD air -c .air.toml -- \
    -port=${BACKEND_PORT:-8080} \
    -db-host=${DB_HOST:-postgres} \
    -db-port=${DB_PORT:-5432} \
    -db-user=${DB_USER:-postgres} \
    -db-password=${DB_PASSWORD:-postgres} \
    -db-name=${DB_NAME:-recipin_dev} \
    -db-sslmode=${DB_SSLMODE:-disable} \
    -jwt-secret=${JWT_SECRET:-change-me} \
    -jwt-expiration=${JWT_EXPIRATION_HOURS:-24}

# Builder stage
FROM base AS builder
COPY . .
RUN go build -o main ./cmd/api

# Production stage
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ./main \
    -port=${BACKEND_PORT:-8080} \
    -db-host=${DB_HOST:-postgres} \
    -db-port=${DB_PORT:-5432} \
    -db-user=${DB_USER:-postgres} \
    -db-password=${DB_PASSWORD:-postgres} \
    -db-name=${DB_NAME:-recipin_dev} \
    -db-sslmode=${DB_SSLMODE:-disable} \
    -jwt-secret=${JWT_SECRET:-change-me} \
    -jwt-expiration=${JWT_EXPIRATION_HOURS:-24}
