.PHONY: dev build build-frontend check format clean help tail-log

# Default target
help:
	@echo "Available targets:"
	@echo "  dev       - Start development servers (frontend + backend)"
	@echo "  build     - Build production binary with embedded frontend"
	@echo "  check     - Run linting and type checking"
	@echo "  format    - Format all code"
	@echo "  clean     - Clean build artifacts"
	@echo "  tail-log  - Show the last 100 lines of the log"

# Development mode - start both frontend and backend
dev:
	@ENV=development ./scripts/shoreman.sh

# Production build
build: build-frontend
	@echo "Building production binary..."
	@mkdir -p bin
	@go build -o bin/minibb ./cmd/minibb
	@echo "Production binary created at bin/minibb"

# Build frontend for production
build-frontend:
	@echo "Building frontend..."
	@cd web && npm install && npm run build

# Linting and type checking
check:
	@echo "Running Go checks..."
	@go vet ./...
	@go mod tidy
	@echo "Running frontend checks..."
	@cd web && npm run lint

# Format code
format:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Formatting frontend code..."
	@cd web && npm run format

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf web/dist/
	@rm -rf web/node_modules/

# Display the last 100 lines of development log with ANSI codes stripped
tail-log:
	@tail -100 ./dev.log | perl -pe 's/\e\[[0-9;]*m(?:\e\[K)?//g'