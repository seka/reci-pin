.PHONY: run test lint mock clean help

# Default target
.DEFAULT_GOAL := help

# Server settings
PORT := 8080

run: ## Run the API server with default settings
	go run ./cmd/api/main.go -port $(PORT)

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run ./...

mock: ## Generate mocks
	go generate ./...

clean: ## Clean build artifacts
	go clean
	rm -f main

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
