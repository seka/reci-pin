# Default target
.DEFAULT_GOAL := help

# Server settings
PORT ?= 8080

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=recipin-backend

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

# Testing tool settings
TARGET ?= http://localhost:8080/api
STRESS_CONFIG ?= ./cmd/stress_test/config.json
STRESS_TARGET ?= http://localhost:8080/api

# ============================================================================
# Application targets
# ============================================================================

.PHONY: app.run
app.run: ## Run the API server with default settings
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

.PHONY: app.fmt
app.fmt: ## Format code and auto-fix issues
	@echo "Running golangci-lint fmt..."
	golangci-lint fmt ./...

.PHONY: app.test
app.test: ## Run tests
	go test -v ./...

.PHONY: app.lint
app.lint: ## Run linter
	@echo "Running golangci-lint ..."
	golangci-lint run ./...
	
.PHONY: app.lint.fix
app.lint.fix: ## Run linter and auto-fix
	@echo "Running golangci-lint run --fix ..."
	golangci-lint run --fix ./...

.PHONY: app.mock
app.mock: ## Generate mocks
	go generate ./...

.PHONY: app.clean
app.clean: ## Clean build artifacts
	go clean
	rm -f main

# ============================================================================
# Database seed targets
# ============================================================================

.PHONY: seed.run
seed.run: ## Run database seed
	@echo "Running database seed..."
	go run ./cmd/seed/main.go

.PHONY: seed.force
seed.force: ## Clean database and run seed
	@echo "Cleaning and seeding database..."
	go run ./cmd/seed/main.go -clean

# ============================================================================
# Integration test targets
# ============================================================================

.PHONY: integration-test.run
integration-test.run: ## Run integration tests (set TARGET to customize)
	@echo "Running integration tests against $(TARGET)..."
	go run ./cmd/integration_test/main.go -target=$(TARGET)

# ============================================================================
# Stress test targets
# ============================================================================

.PHONY: stress-test.run
stress-test.run: ## Run stress test (set STRESS_CONFIG and STRESS_TARGET to customize)
	@echo "Running stress test..."
	go run ./cmd/stress_test/main.go -config=$(STRESS_CONFIG) -target-url=$(STRESS_TARGET)

# ============================================================================
# Utility targets
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
