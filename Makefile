# ==============================================================================
# CRM Party Service Makefile
# ==============================================================================

SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Application Configuration
APP_NAME       ?= party-service
BIN_DIR        ?= bin
BINARY_NAME    ?= $(BIN_DIR)/$(APP_NAME)
MAIN_SRC       ?= ./cmd/service
COVERAGE_OUT   ?= coverage.out
COVERAGE_HTML  ?= coverage.html

# Tools
GO             ?= go
SQLC           ?= $(shell which sqlc 2>/dev/null || echo $(HOME)/go/bin/sqlc)
DOCKER         ?= docker
DOCKER_COMPOSE ?= docker compose

# Environment Defaults for Local Run
PORT           ?= 8080
DATABASE_URL   ?= postgres://postgres:postgres@localhost:5432/party?sslmode=disable
ENV            ?= development

# Colors for Terminal Output
CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RESET  := \033[0m

# ==============================================================================
# Development & Build Targets
# ==============================================================================

.PHONY: help
help: ## Show this help message
	@echo -e "$(CYAN)Usage:$(RESET) make $(GREEN)<target>$(RESET)"
	@echo ""
	@echo -e "$(CYAN)Available Targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the statically-linked service binary
	@echo -e "$(CYAN)[*] Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 $(GO) build -ldflags="-w -s" -o $(BINARY_NAME) $(MAIN_SRC)
	@echo -e "$(GREEN)[✓] Build complete: $(BINARY_NAME)$(RESET)"

.PHONY: run
run: ## Run service locally using environment variables
	@echo -e "$(CYAN)[*] Starting $(APP_NAME) locally...$(RESET)"
	@PORT=$(PORT) DATABASE_URL=$(DATABASE_URL) ENV=$(ENV) $(GO) run $(MAIN_SRC)

.PHONY: clean
clean: ## Remove build artifacts, binaries, and test coverage outputs
	@echo -e "$(CYAN)[*] Cleaning up build artifacts...$(RESET)"
	@rm -rf $(BIN_DIR) $(COVERAGE_OUT) $(COVERAGE_HTML)
	@echo -e "$(GREEN)[✓] Clean complete$(RESET)"

# ==============================================================================
# Testing & Code Quality Targets
# ==============================================================================

.PHONY: test
test: ## Run all unit and integration tests
	@echo -e "$(CYAN)[*] Running all tests...$(RESET)"
	@$(GO) test -v -count=1 ./...

.PHONY: test-unit
test-unit: ## Run internal unit tests only
	@echo -e "$(CYAN)[*] Running unit tests...$(RESET)"
	@$(GO) test -v -count=1 ./internal/...

.PHONY: test-e2e
test-e2e: ## Run E2E API integration tests only
	@echo -e "$(CYAN)[*] Running E2E API integration tests...$(RESET)"
	@$(GO) test -v -count=1 ./test/api/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage and generate HTML report
	@echo -e "$(CYAN)[*] Generating test coverage report...$(RESET)"
	@$(GO) test -count=1 -coverprofile=$(COVERAGE_OUT) ./internal/... ./test/api/...
	@$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo -e "$(GREEN)[✓] Coverage report generated: $(COVERAGE_HTML)$(RESET)"

.PHONY: fmt
fmt: ## Format Go source files using gofmt
	@echo -e "$(CYAN)[*] Formatting Go code...$(RESET)"
	@gofmt -s -w .
	@echo -e "$(GREEN)[✓] Code formatted$(RESET)"

.PHONY: vet
vet: ## Examine Go source code for suspicious constructs
	@echo -e "$(CYAN)[*] Running go vet...$(RESET)"
	@$(GO) vet ./...
	@echo -e "$(GREEN)[✓] Vet passed$(RESET)"

.PHONY: tidy
tidy: ## Tidy and verify Go module dependencies
	@echo -e "$(CYAN)[*] Tidying Go modules...$(RESET)"
	@$(GO) mod tidy
	@$(GO) mod verify
	@echo -e "$(GREEN)[✓] Modules verified$(RESET)"

# ==============================================================================
# Code Generation
# ==============================================================================

.PHONY: sqlc
sqlc: ## Regenerate Go database code using sqlc
	@echo -e "$(CYAN)[*] Running sqlc generate...$(RESET)"
	@$(SQLC) generate
	@echo -e "$(GREEN)[✓] SQL code generation complete$(RESET)"

# ==============================================================================
# Docker & Container Targets
# ==============================================================================

.PHONY: docker-build
docker-build: ## Build Docker container image for party-service
	@echo -e "$(CYAN)[*] Building Docker image $(APP_NAME):latest...$(RESET)"
	@$(DOCKER) build -t $(APP_NAME):latest .
	@echo -e "$(GREEN)[✓] Docker build complete$(RESET)"

.PHONY: docker-up
docker-up: ## Boot Postgres and party-service in background via Docker Compose
	@echo -e "$(CYAN)[*] Starting containers with Docker Compose...$(RESET)"
	@$(DOCKER_COMPOSE) up --build -d
	@echo -e "$(GREEN)[✓] Containers started$(RESET)"

.PHONY: docker-down
docker-down: ## Stop and remove all Docker Compose containers and networks
	@echo -e "$(CYAN)[*] Stopping containers...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@echo -e "$(GREEN)[✓] Containers stopped$(RESET)"

.PHONY: docker-db
docker-db: ## Start only PostgreSQL database container
	@echo -e "$(CYAN)[*] Starting PostgreSQL container...$(RESET)"
	@$(DOCKER_COMPOSE) up -d postgres
	@echo -e "$(GREEN)[✓] PostgreSQL is running on port 5432$(RESET)"

.PHONY: docker-monitoring
docker-monitoring: ## Start Prometheus and Grafana monitoring stack
	@echo -e "$(CYAN)[*] Starting Prometheus (9090) and Grafana (3000)...$(RESET)"
	@$(DOCKER_COMPOSE) up -d prometheus grafana
	@echo -e "$(GREEN)[✓] Prometheus available at: http://localhost:9090$(RESET)"
	@echo -e "$(GREEN)[✓] Grafana available at:    http://localhost:3000 (admin/admin)$(RESET)"

.PHONY: docker-logs
docker-logs: ## Follow live logs from Docker Compose containers
	@$(DOCKER_COMPOSE) logs -f

.PHONY: docker-ps
docker-ps: ## List status of running Docker Compose services
	@$(DOCKER_COMPOSE) ps
