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
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]

# Builder stage
FROM base AS builder
COPY . .
RUN go build -o main ./cmd/api

# Production stage
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/main .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./main"]
