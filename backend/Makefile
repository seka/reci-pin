# Default target
.DEFAULT_GOAL := help

# Server settings
PORT ?= 8080

# Go cache and module cache confinement
export GOCACHE=$(shell pwd)/.gocache
export GOMODCACHE=$(shell pwd)/.gomodcache

# ============================================================================
# Development targets
# ============================================================================

.PHONY: bootstrap
bootstrap: ## Install all development tools
	@echo "Installing development tools..."
	go install go.uber.org/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	@echo "All tools installed successfully!"

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

# Elasticsearch settings
ES_HOST ?= localhost
ES_PORT ?= 9200

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
	@echo "Running go fix and golangci-lint run --fix ..."
	go fix ./...
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
seed.run: ## Run database seed and sync to Elasticsearch
	@echo "Running database seed..."
	go run ./cmd/seed/main.go
	$(MAKE) seed.image
	$(MAKE) sync-es.run

.PHONY: seed.force
seed.force: ## Clean database, clear ES index and run seed
	@echo "Cleaning and seeding database..."
	go run ./cmd/seed/main.go -clean
	$(MAKE) seed.image
	$(MAKE) sync-es.clear
	$(MAKE) sync-es.run

.PHONY: seed.image
seed.image: ## Upload seed images to LocalStack
	@echo "Uploading seed images to LocalStack..."
	docker cp internal/infrastructure/storage/mock/seed_placeholder.png reci-pin-localstack:/tmp/seed_placeholder.png
	docker exec reci-pin-localstack awslocal s3 cp /tmp/seed_placeholder.png s3://recipin-bucket/recipes/1/seed_1.jpg --content-type image/png
	docker exec reci-pin-localstack awslocal s3 cp /tmp/seed_placeholder.png s3://recipin-bucket/recipes/2/seed_1.jpg --content-type image/png
	docker exec reci-pin-localstack awslocal s3 cp /tmp/seed_placeholder.png s3://recipin-bucket/recipes/2/seed_2.jpg --content-type image/png

.PHONY: sync-es.run
sync-es.run: ## Sync data from DB to Elasticsearch
	@echo "Syncing data to Elasticsearch from host..."
	go run ./cmd/sync_es/main.go \
		-db-host $(DB_HOST) \
		-db-port $(DB_PORT) \
		-db-user $(DB_USER) \
		-db-password $(DB_PASSWORD) \
		-db-name $(DB_NAME) \
		-db-sslmode $(DB_SSLMODE) \
		-es-host $(ES_HOST) \
		-es-port $(ES_PORT)

.PHONY: sync-es.clear
sync-es.clear: ## Clear Elasticsearch index
	@echo "Clearing Elasticsearch index..."
	curl -f -X DELETE "http://$(ES_HOST):$(ES_PORT)/recipes" || true

.PHONY: sync-es.docker
sync-es.docker: ## Sync data from DB to Elasticsearch (run inside Docker)
	@echo "Syncing data to Elasticsearch inside Docker container..."
	docker-compose -f ../docker-compose.yml exec backend go run ./cmd/sync_es/main.go \
		-db-host postgres \
		-es-host elasticsearch

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
