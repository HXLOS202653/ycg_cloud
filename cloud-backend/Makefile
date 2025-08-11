# Makefile for YCG Cloud Storage Backend
# Provides convenient commands for development workflow

.PHONY: help build run test clean lint format deps check-deps dev install-tools

# Default target
help: ## Show this help message
	@echo "YCG Cloud Storage Backend - Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build commands
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build completed: bin/server"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
	@echo "Linux build completed: bin/server-linux"

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o bin/server.exe cmd/server/main.go
	@echo "Windows build completed: bin/server.exe"

# Run commands
run: ## Run the application in development mode
	@echo "Starting development server..."
	@go run cmd/server/main.go

run-prod: ## Run with production config
	@echo "Starting production server..."
	@APP_ENV=production go run cmd/server/main.go

# Development commands
dev: format lint test ## Run all development checks (format, lint, test)

watch: ## Run with file watching (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to normal run..."; \
		make run; \
	fi

# Testing commands
test: ## Run all tests
	@echo "Running tests..."
	@APP_ENV=test go test -v ./...

test-short: ## Run tests in short mode (skip integration tests)
	@echo "Running tests in short mode..."
	@APP_ENV=test go test -short -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@APP_ENV=test go test -v -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

test-coverage-report: test-coverage ## Generate and open coverage report
	@echo "Opening coverage report..."
	@if command -v xdg-open > /dev/null; then \
		xdg-open coverage.html; \
	elif command -v open > /dev/null; then \
		open coverage.html; \
	elif command -v start > /dev/null; then \
		start coverage.html; \
	else \
		echo "Please open coverage.html manually"; \
	fi

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@APP_ENV=test go test -race -v ./...

test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	@APP_ENV=test go test -bench=. -benchmem ./...

test-unit: ## Run only unit tests (exclude integration tests)
	@echo "Running unit tests..."
	@APP_ENV=test go test -short -v ./internal/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@APP_ENV=test go test -v -tags=integration ./tests/...

test-clean: ## Clean test artifacts
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -rf ./test_uploads
	@go clean -testcache

# Code quality commands
lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run --config .golangci.yml

lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	@golangci-lint run --config .golangci.yml --fix

format: ## Format code with gofmt and goimports
	@echo "Formatting code..."
	@echo "Running gofmt..."
	@gofmt -s -w .
	@echo "Running goimports..."
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "Code formatting completed"

format-check: ## Check if code is properly formatted
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not properly formatted:"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "All files are properly formatted"; \
	fi

# Security commands
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Dependency commands
deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

deps-clean: ## Clean module cache
	@echo "Cleaning module cache..."
	@go clean -modcache

check-deps: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	@if command -v go-mod-outdated > /dev/null; then \
		go list -u -m -json all | go-mod-outdated -update -direct; \
	else \
		echo "go-mod-outdated not installed. Install with: go install github.com/psampaz/go-mod-outdated@latest"; \
	fi

# Database commands
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run cmd/migration/main.go

db-rollback: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@go run cmd/migration/main.go -rollback

db-seed: ## Seed database with test data
	@echo "Seeding database..."
	@go run cmd/migration/main.go -seed

# Installation commands
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing gosec..."
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Installing air (live reload)..."
	@go install github.com/cosmtrek/air@latest
	@echo "Installing go-mod-outdated..."
	@go install github.com/psampaz/go-mod-outdated@latest
	@echo "Installing swag (Swagger docs)..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "All tools installed successfully!"

# Documentation commands
docs: ## Generate API documentation
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o docs/swagger; \
	else \
		echo "swag not installed. Install with: make install-tools"; \
	fi

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t ycg-cloud-backend .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 ycg-cloud-backend

# Cleanup commands
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "Cleanup completed"

clean-all: clean deps-clean ## Clean everything including dependencies

# Health check commands
health: ## Check application health
	@echo "Checking application health..."
	@curl -f http://localhost:8080/health || echo "Application is not running or unhealthy"

# Git hooks setup
setup-hooks: ## Setup git hooks for code quality
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/sh\nmake format lint' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

# CI/CD simulation
ci: format-check lint test ## Run CI pipeline locally

# Quick development workflow
quick: format lint ## Quick format and lint check

# Full development workflow
full: clean deps format lint test build ## Full development workflow

# Production preparation
prod-ready: full security docs ## Prepare for production deployment