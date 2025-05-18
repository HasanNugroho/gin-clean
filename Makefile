.PHONY: all build run watch gen-docs \
        migration-create migration-up migration-down \
        docker-run docker-down \
        help

# Default target
all: build

## ---------- Build & Run ----------

setup:
	@echo "📦 Setting up project..."
	@go mod download & go mod tidy
	@cp .env.example .env || true

## build: Build the application binary
build:
	@echo "🔨 Building application..."
	@go build -o main cmd/api/main.go

## run: Run the application
run:
	@echo "🚀 Running application..."
	@go run ./cmd/api

## watch: Watch for file changes and auto-reload (requires air)
watch:
	@echo "👀 Watching for changes..."
	@air -c .air.toml


## ---------- Documentation ----------

## gen-docs: Generate Swagger API documentation
gen-docs:
	@echo "📖 Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs


## ---------- Migration ----------

include .env

## migration-create: Create a new DB migration. Usage: make migration-create desc=your_description
migration-create:
	@test -n "$(desc)" || (echo "❌ Missing desc param. Usage: make migration-create desc=your_description" && exit 1)
	@migrate create -ext=sql -dir=migrations -seq $(desc)

## migration-up: Apply all up migrations
migration-up:
	@echo "⬆️  Running DB migrations..."
	@migrate -source file://./migrations \
	         -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" up

## migration-down: Rollback all migrations
migration-down:
	@echo "⬇️  Reverting DB migrations..."
	@migrate -source file://./migrations \
	         -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" down


## ---------- Docker ----------

## docker-run: Start Docker containers (if defined)
docker-run:
	@echo "🐳 Starting Docker services..."
	@docker compose up -d

## docker-down: Stop Docker containers
docker-down:
	@echo "🛑 Stopping Docker services..."
	@docker compose down


## ---------- Help ----------

## help: Show this help message
help:
	@echo
	@echo "📦 Available Makefile commands:"
	@echo
	@grep -E '^##' Makefile | sed -e 's/^## //' | column -t -s ':' | sed -e 's/^/ /'
	@echo
