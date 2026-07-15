# =============================================================================
# HelixKnowledge Skill Graph System - Makefile
# =============================================================================
# Build automation for the Skill Graph System
# =============================================================================

# ---------------------------------------------------------------------------
# Variables
# ---------------------------------------------------------------------------
VERSION     ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME  ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Binary names
BINARY_SERVER = server
BINARY_WORKER = worker
BINARY_CLI    = skillctl
BINARY_TUI    = skill-tui

# Directories
CMD_DIR     = ./cmd
BUILD_DIR   = ./bin
DIST_DIR    = ./dist
MIGRATIONS_DIR = ./migrations

# Go settings
GO          = go
GOFLAGS     = -mod=vendor
LDFLAGS     = -s -w \
              -X main.Version=$(VERSION) \
              -X main.Commit=$(COMMIT) \
              -X main.BuildTime=$(BUILD_TIME) \
              -extldflags '-static'

# Container settings
CONTAINER_RUNTIME ?= docker
COMPOSE_CMD       ?= $(CONTAINER_RUNTIME) compose
IMAGE_NAME        ?= skill-system
IMAGE_TAG         ?= $(VERSION)

# Database settings (for local dev)
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_NAME     ?= skilldb
DB_USER     ?= skilluser
DB_PASSWORD ?= skillpassword
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ---------------------------------------------------------------------------
# Colors for output
# ---------------------------------------------------------------------------
BLUE        = \033[0;34m
GREEN       = \033[0;32m
YELLOW      = \033[1;33m
RED         = \033[0;31m
NC          = \033[0m

# ---------------------------------------------------------------------------
# Default target
# ---------------------------------------------------------------------------
.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Help
# ---------------------------------------------------------------------------
.PHONY: help
help: ## Show this help message
	@echo -e "$(BLUE)HelixKnowledge Skill Graph System - Make Targets$(NC)"
	@echo "═══════════════════════════════════════════════════════════════"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo -e "$(BLUE)Variables:$(NC)"
	@echo "  VERSION=$(VERSION)"
	@echo "  COMMIT=$(COMMIT)"
	@echo "  IMAGE_NAME=$(IMAGE_NAME)"
	@echo "  IMAGE_TAG=$(IMAGE_TAG)"

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------
.PHONY: build
build: build-server build-worker build-cli build-tui ## Build all binaries

.PHONY: build-server
build-server: ## Build API server binary
	@echo -e "$(BLUE)Building server...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_SERVER) $(CMD_DIR)/server
	@echo -e "$(GREEN)Built: $(BUILD_DIR)/$(BINARY_SERVER)$(NC)"

.PHONY: build-worker
build-worker: ## Build background worker binary
	@echo -e "$(BLUE)Building worker...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_WORKER) $(CMD_DIR)/worker
	@echo -e "$(GREEN)Built: $(BUILD_DIR)/$(BINARY_WORKER)$(NC)"

.PHONY: build-cli
build-cli: ## Build CLI binary
	@echo -e "$(BLUE)Building CLI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_CLI) $(CMD_DIR)/cli
	@echo -e "$(GREEN)Built: $(BUILD_DIR)/$(BINARY_CLI)$(NC)"

.PHONY: build-tui
build-tui: ## Build TUI binary
	@echo -e "$(BLUE)Building TUI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_TUI) $(CMD_DIR)/tui
	@echo -e "$(GREEN)Built: $(BUILD_DIR)/$(BINARY_TUI)$(NC)"

# ---------------------------------------------------------------------------
# Development
# ---------------------------------------------------------------------------
.PHONY: run
run: build-server ## Run server in development mode
	@echo -e "$(BLUE)Starting server...$(NC)"
	$(BUILD_DIR)/$(BINARY_SERVER) --config ./config/config.toml

.PHONY: run-worker
run-worker: build-worker ## Run worker in development mode
	@echo -e "$(BLUE)Starting worker...$(NC)"
	$(BUILD_DIR)/$(BINARY_WORKER) --config ./config/config.toml

.PHONY: dev
 dev: ## Start development stack (docker compose)
	@echo -e "$(BLUE)Starting development stack...$(NC)"
	$(COMPOSE_CMD) -f docker-compose.yml up -d db
	@echo -e "$(GREEN)Database ready. Run 'make run' to start server.$(NC)"

.PHONY: dev-down
dev-down: ## Stop development stack
	@echo -e "$(BLUE)Stopping development stack...$(NC)"
	$(COMPOSE_CMD) -f docker-compose.yml down

# ---------------------------------------------------------------------------
# Testing
# ---------------------------------------------------------------------------
.PHONY: test
test: ## Run all tests
	@echo -e "$(BLUE)Running tests...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo -e "$(BLUE)Running unit tests...$(NC)"
	$(GO) test -v -short ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo -e "$(BLUE)Running integration tests...$(NC)"
	$(GO) test -v -run Integration ./...

.PHONY: coverage
coverage: test ## Show test coverage report
	@echo -e "$(BLUE)Coverage report:$(NC)"
	$(GO) tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: test ## Generate HTML coverage report
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo -e "$(GREEN)Report: coverage.html$(NC)"

# ---------------------------------------------------------------------------
# Linting & Quality
# ---------------------------------------------------------------------------
.PHONY: lint
lint: ## Run all linters
	@echo -e "$(BLUE)Running linters...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo -e "$(YELLOW)golangci-lint not installed, running go vet...$(NC)"; \
		$(GO) vet ./...; \
	fi

.PHONY: fmt
fmt: ## Format Go source code
	@echo -e "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo -e "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo -e "$(BLUE)Tidying modules...$(NC)"
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: generate
generate: ## Run go generate
	@echo -e "$(BLUE)Running go generate...$(NC)"
	$(GO) generate ./...

# ---------------------------------------------------------------------------
# Database
# ---------------------------------------------------------------------------
.PHONY: migrate-up
migrate-up: ## Apply database migrations
	@echo -e "$(BLUE)Applying migrations...$(NC)"
	./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down: ## Rollback one migration
	@echo -e "$(BLUE)Rolling back migration...$(NC)"
	./scripts/migrate.sh down

.PHONY: migrate-status
migrate-status: ## Show migration status
	./scripts/migrate.sh status

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=add_skills)
	@if [ -z "$(name)" ]; then \
		echo -e "$(RED)Error: name required. Usage: make migrate-create name=add_skills$(NC)"; \
		exit 1; \
	fi
	./scripts/migrate.sh create $(name)

.PHONY: db-reset
db-reset: ## Reset database (drop and recreate)
	@echo -e "$(YELLOW)WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(COMPOSE_CMD) exec db psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"; \
		$(COMPOSE_CMD) exec db psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"; \
		$(MAKE) migrate-up; \
		echo -e "$(GREEN)Database reset complete$(NC)"; \
	else \
		echo "Cancelled"; \
	fi

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo -e "$(BLUE)Building Docker image...$(NC)"
	$(COMPOSE_CMD) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME)
	@echo -e "$(GREEN)Image: $(IMAGE_NAME):$(IMAGE_TAG)$(NC)"

.PHONY: docker-push
docker-push: docker-build ## Push Docker image to registry
	@echo -e "$(BLUE)Pushing image...$(NC)"
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

.PHONY: docker-up
docker-up: ## Start services with docker compose
	@echo -e "$(BLUE)Starting services...$(NC)"
	$(COMPOSE_CMD) up -d

.PHONY: docker-down
docker-down: ## Stop services
	@echo -e "$(BLUE)Stopping services...$(NC)"
	$(COMPOSE_CMD) down

.PHONY: docker-logs
docker-logs: ## Show service logs
	$(COMPOSE_CMD) logs -f

.PHONY: docker-ps
docker-ps: ## Show running containers
	$(COMPOSE_CMD) ps

# ---------------------------------------------------------------------------
# Backup & Restore
# ---------------------------------------------------------------------------
.PHONY: backup
backup: ## Create backup
	./scripts/backup.sh

.PHONY: restore
restore: ## Restore from backup (usage: make restore file=backup.tar.gz)
	@if [ -z "$(file)" ]; then \
		echo -e "$(RED)Error: file required. Usage: make restore file=backup.tar.gz$(NC)"; \
		exit 1; \
	fi
	./scripts/restore.sh $(file)

# ---------------------------------------------------------------------------
# Packaging
# ---------------------------------------------------------------------------
.PHONY: package
package: ## Create distributable package
	@echo -e "$(BLUE)Creating package...$(NC)"
	./scripts/package.sh --version $(VERSION) --output $(DIST_DIR)

.PHONY: package-no-source
package-no-source: ## Create package without source code
	@echo -e "$(BLUE)Creating binary-only package...$(NC)"
	./scripts/package.sh --version $(VERSION) --output $(DIST_DIR) --no-source

# ---------------------------------------------------------------------------
# Cleaning
# ---------------------------------------------------------------------------
.PHONY: clean
clean: ## Remove build artifacts
	@echo -e "$(BLUE)Cleaning...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean -cache

.PHONY: clean-all
clean-all: docker-down clean ## Full cleanup including containers and volumes
	@echo -e "$(BLUE)Full cleanup...$(NC)"
	$(COMPOSE_CMD) down -v --remove-orphans
	rm -rf $(BUILD_DIR) $(DIST_DIR)

# ---------------------------------------------------------------------------
# Utilities
# ---------------------------------------------------------------------------
.PHONY: version
version: ## Show version info
	@echo "Version:  $(VERSION)"
	@echo "Commit:   $(COMMIT)"
	@echo "Built:    $(BUILD_TIME)"

.PHONY: install
install: ## Run installation script
	@echo -e "$(BLUE)Running installer...$(NC)"
	./scripts/install.sh

.PHONY: status
status: ## Show stack status
	./scripts/status.sh

.PHONY: docs
docs: ## Generate documentation
	@echo -e "$(BLUE)Documentation is in the docs/ directory$(NC)"
	@ls -1 docs/

# ---------------------------------------------------------------------------
# CI/CD Helpers
# ---------------------------------------------------------------------------
.PHONY: ci-build
ci-build: fmt vet lint test build ## CI pipeline: format, lint, test, build

.PHONY: ci-docker
ci-docker: docker-build test-integration ## CI pipeline: build image, run integration tests

.PHONY: release
release: clean ci-build package ## Full release build
	@echo -e "$(GREEN)Release $(VERSION) ready in $(DIST_DIR)/$(NC)"
