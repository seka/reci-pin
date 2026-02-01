# Default target
.DEFAULT_GOAL := help

# Server settings
PORT ?= 8080

# Database settings
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= recipin_dev
DB_SSLMODE ?= disable

# JWT settings
JWT_SECRET ?= change-me
JWT_EXPIRATION ?= 24

.PHONY: run
run: ## Run the API server with default settings
	go run ./cmd/api/main.go \
		-port $(PORT) \
		-db-host $(DB_HOST) \
		-db-port $(DB_PORT) \
		-db-user $(DB_USER) \
		-db-password $(DB_PASSWORD) \
		-db-name $(DB_NAME) \
		-db-sslmode $(DB_SSLMODE) \
		-jwt-secret $(JWT_SECRET) \
		-jwt-expiration $(JWT_EXPIRATION)

.PHONY: test 
test: ## Run tests
	go test -v ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

.PHONY: mock
mock: ## Generate mocks
	go generate ./...

.PHONY: clean
clean: ## Clean build artifacts
	go clean
	rm -f main

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
