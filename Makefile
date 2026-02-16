# --------------------------------------
# Configuración general
# --------------------------------------

SHELL           := /usr/bin/env bash -o pipefail
.SHELLFLAGS     := -ec
.DEFAULT_GOAL   := help

GO              ?= go
CURRENT_PATH    := $(shell pwd)
OUTPUT_PATH     := $(CURRENT_PATH)/bin
SERVICE_NAME    := bike-rental

export CGO_ENABLED := 0
export GO111MODULE := on

# --------------------------------------
# Comandos útiles
# --------------------------------------

.PHONY: help migrate migrate-fresh migrate-no-seed run test test-coverage build-linux clean swagger docker-build docker-run docker-stop docker-clean docker-logs

help:
	@echo "Available commands:"
	@echo ""
	@echo "Local Development:"
	@echo "  make migrate         - Run database migrations with seed data"
	@echo "  make migrate-fresh   - Drop tables and run fresh migration"
	@echo "  make migrate-no-seed - Run migration without seed data"
	@echo "  make run             - Run the API server"
	@echo "  make build-linux     - Build static binary for Linux"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make clean           - Remove build artifacts and database"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build    - Build distroless Docker image (~10MB)"
	@echo "  make docker-run      - Run container from image"
	@echo "  make docker-stop     - Stop running container"
	@echo "  make docker-clean    - Remove container and images"
	@echo "  make docker-logs     - View container logs"

migrate:
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh

migrate-fresh:
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh fresh

migrate-no-seed:
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh no-seed

##@ Build

build-linux:
	@echo ">> Compilando binario para Linux"
	@bash $(CURRENT_PATH)/build.sh $(GO) $(OUTPUT_PATH)/$(SERVICE_NAME)-linux-amd64 cmd/api/main.go linux amd64

run:
	@go run cmd/api/main.go

test:
	@go test ./... -v

test-coverage:
	@echo ">> Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out 2>&1 | tee /tmp/test_output.txt
	@go tool cover -html=coverage.out -o coverage.html
	@echo "\n>> Updating coverage badge..."
	@bash -c ' \
		get_color() { \
			local val=$$1; \
			local pct=$$(echo "$$val" | cut -d. -f1); \
			if [ "$$pct" -ge 80 ] 2>/dev/null; then echo "brightgreen"; \
			elif [ "$$pct" -ge 60 ] 2>/dev/null; then echo "yellow"; \
			elif [ "$$pct" -ge 40 ] 2>/dev/null; then echo "orange"; \
			else echo "red"; fi; \
		}; \
		handlers=$$(grep "internal/handlers" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		services=$$(grep "internal/services" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		repos=$$(grep "internal/repositories" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		utils=$$(grep "internal/utils" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		types=$$(grep "internal/types" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		config=$$(grep "internal/config" /tmp/test_output.txt | grep -o "[0-9.]*%" | tr -d "%" || echo "0.0"); \
		total=$$(echo "scale=1; ($$handlers + $$services + $$repos + $$utils + $$types + $$config) / 6" | bc); \
		color=$$(get_color $$total); \
		badge="![Coverage](https://img.shields.io/badge/coverage-$${total}%25-$$color)."; \
		sed -i.bak "3s|^.*Coverage.*$$|$$badge|" README.md && rm -f README.md.bak; \
		rm -f /tmp/test_output.txt; \
		echo "Coverage: $$total% (handlers=$$handlers% services=$$services% repos=$$repos% utils=$$utils% types=$$types% config=$$config%)"; \
		echo "✓ Coverage badge updated in README.md"; \
	'

swagger: ## Genera documentación Swagger
	@if ! command -v swag &> /dev/null; then \
		echo ">> swag not found, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo ">> Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs
	@echo ">> Swagger docs generated in /docs"
	@echo ">> Visit http://localhost:8080/swagger/index.html"

##@ Utilidades

clean:
	@echo ">> Limpiando archivos de compilación"
	@rm -rf $(OUTPUT_PATH)
	@rm -f coverage.out coverage.html
	@rm -f data/bike_rental.db
	@$(GO) clean -x -i ./... > /dev/null 2>&1 || true

##@ Docker

docker-build: swagger build-linux
	@echo ">> Building distroless Docker image..."
	@ls -la bin/  # Debug: verificar que el archivo existe
	@docker build -f Dockerfile -t bike-rental-service:latest .
	@echo ">> Docker image built successfully"
	@docker images bike-rental-service:latest

docker-run:
	@echo "Running Docker container with .env file..."
	@docker run -d \
		--name bike-rental-service \
		--env-file .env \
		-e SQLITE_PATH=/app/data/bike_rental.db \
		-p 8080:8080 \
		-v $(PWD)/data:/app/data \
		bike-rental-service:latest
	@echo "Container started with environment variables from .env"
	@echo "API: http://localhost:8080"
	@echo "Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "View logs: make docker-logs"

docker-stop:
	@echo "Stopping Docker container..."
	@docker stop bike-rental-service || true
	@docker rm bike-rental-service || true
	@echo "Container stopped and removed"

docker-clean: docker-stop
	@echo "Cleaning Docker images..."
	@docker rmi bike-rental-service:latest || true
	@echo "Docker images cleaned"

docker-logs:
	@docker logs -f bike-rental-service