# Simple Makefile for a Go project

include .env
export
# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Migrate the database
migrate-up:
	@echo "Running migrations..."
	@echo "Connection string: postgres://$(DB_USERNAME):***@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"
	@goose -dir "internal/database/migrations" postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up

# Rollback the database
migrate-down:
	@echo "Rolling back migrations..."
	@goose -dir "internal/database/migrations" postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
	    air; \
	    Write-Output 'Watching...'; \
	} else { \
	    Write-Output 'Installing air...'; \
	    go install github.com/air-verse/air@latest; \
	    air; \
	    Write-Output 'Watching...'; \
	}"


.PHONY: all build run test clean watch docker-run docker-down itest migrate-up migrate-down