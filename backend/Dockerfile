# Base stage
FROM golang:1.26-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Development stage
FROM base AS dev
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/air-verse/air@latest
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
    -jwt-expiration=${JWT_EXPIRATION_HOURS:-24} \
    -storage-bucket=${STORAGE_BUCKET:-recipin-bucket} \
    -storage-endpoint=${STORAGE_ENDPOINT} \
    -storage-region=${AWS_REGION:-us-east-1} \
    -storage-access-key=${AWS_ACCESS_KEY_ID} \
    -storage-secret-key=${AWS_SECRET_ACCESS_KEY} \
    -storage-public-url=${STORAGE_PUBLIC_URL} \
    -es-addresses=${ELASTICSEARCH_ADDRESS:-http://elasticsearch:9200} \
    -mail-host=${MAIL_HOST:-localhost} \
    -mail-port=${MAIL_PORT:-1025} \
    -mail-from=${MAIL_FROM:-no-reply@reci-pin.com}

# Builder stage
FROM base AS builder
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o main ./cmd/api

# Production stage
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ./main \
    -port=${BACKEND_PORT} \
    -db-host=${DB_HOST} \
    -db-port=${DB_PORT} \
    -db-user=${DB_USER} \
    -db-password=${DB_PASSWORD} \
    -db-name=${DB_NAME} \
    -db-sslmode=${DB_SSLMODE} \
    -jwt-secret=${JWT_SECRET} \
    -jwt-expiration=${JWT_EXPIRATION_HOURS} \
    -storage-bucket=${STORAGE_BUCKET} \
    -storage-endpoint=${STORAGE_ENDPOINT} \
    -storage-region=${AWS_REGION} \
    -storage-access-key=${AWS_ACCESS_KEY_ID} \
    -storage-secret-key=${AWS_SECRET_ACCESS_KEY} \
    -storage-public-url=${STORAGE_PUBLIC_URL} \
    -es-addresses=${ELASTICSEARCH_ADDRESS} \
    -mail-host=${MAIL_HOST} \
    -mail-port=${MAIL_PORT} \
    -mail-from=${MAIL_FROM}
