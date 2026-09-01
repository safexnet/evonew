.PHONY: help dev run build test clean swagger deps docker-build docker-run install setup migrate-up migrate-down logs

# Configuration
APP_NAME=evolution-go
MAIN_PATH=cmd/evolution-go/main.go
BUILD_DIR=build
GO=go
VERSION=$(shell grep -om1 "v[0-9].*" CHANGELOG.md)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
GOFLAGS=-v

# Output colors
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

##@ Help

help: ## Displays this help message
	@echo "$(GREEN)Evolution GO - Makefile$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

dev: ## Runs application in development mode
	@echo "$(GREEN)🚀 Running Evolution GO in development mode...$(NC)"
	$(GO) run $(LDFLAGS) $(MAIN_PATH) -dev

run: ## Runs application in production mode
	@echo "$(GREEN)🚀 Running Evolution GO...$(NC)"
	$(GO) run $(MAIN_PATH)

watch: ## Runs application with hot reload (requires air)
	@if command -v air > /dev/null; then \
		echo "$(GREEN)🔥 Running with hot reload...$(NC)"; \
		air; \
	else \
		echo "$(RED)❌ Air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		exit 1; \
	fi

##@ Build

build: ## Compiles the application
	@echo "$(GREEN)🔨 Compiling $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

build-linux: ## Compiles for Linux
	@echo "$(GREEN)🔨 Compiling for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Linux build complete$(NC)"

build-windows: ## Compiles for Windows
	@echo "$(GREEN)🔨 Compiling for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✅ Windows build complete$(NC)"

build-all: build build-linux build-windows ## Compiles for all platforms
	@echo "$(GREEN)✅ All builds complete$(NC)"

install: build ## Compiles and installs to GOPATH
	@echo "$(GREEN)📦 Installing $(APP_NAME)...$(NC)"
	$(GO) install $(MAIN_PATH)
	@echo "$(GREEN)✅ Installed successfully$(NC)"

##@ Tests

test: ## Runs all tests
	@echo "$(GREEN)🧪 Running tests...$(NC)"
	$(GO) test -v ./...

test-coverage: ## Runs tests with coverage
	@echo "$(GREEN)🧪 Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage generated: coverage.html$(NC)"

test-race: ## Runs tests checking race conditions
	@echo "$(GREEN)🧪 Running tests with race detector...$(NC)"
	$(GO) test -race -v ./...

bench: ## Runs benchmarks
	@echo "$(GREEN)⚡ Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./...

##@ Dependencies

deps: ## Installs dependencies
	@echo "$(GREEN)📦 Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod verify
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

deps-update: ## Updates dependencies
	@echo "$(GREEN)📦 Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

deps-clean: ## Cleans unused dependencies
	@echo "$(GREEN)🧹 Cleaning dependencies...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependencies cleaned$(NC)"

deps-reset: ## Cleans cache and reinstalls dependencies
	@echo "$(GREEN)🔄 Resetting dependencies and cache...$(NC)"
	@echo "$(YELLOW)Cleaning cache and modules...$(NC)"
	$(GO) clean -cache -modcache -i -r
	@echo "$(YELLOW)Downloading modules...$(NC)"
	$(GO) mod download
	@echo "$(YELLOW)Organizing modules...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependencies reset and updated$(NC)"

##@ Documentation

swagger: ## Generates Swagger documentation
	@echo "$(GREEN)📚 Generating Swagger documentation...$(NC)"
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH) -o ./docs; \
		echo "$(GREEN)✅ Swagger generated successfully$(NC)"; \
	else \
		echo "$(RED)❌ Swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest$(NC)"; \
		exit 1; \
	fi

docs: ## Opens local documentation
	@echo "$(GREEN)📖 Opening documentation...$(NC)"
	@if [ -f "docs/wiki/README.md" ]; then \
		echo "Documentation available at: docs/wiki/README.md"; \
	else \
		echo "$(RED)❌ Documentation not found$(NC)"; \
	fi

##@ Database

migrate-up: ## Runs database migrations
	@echo "$(GREEN)🗃️  Running migrations...$(NC)"
	@if [ -d "migrations" ]; then \
		$(GO) run $(MAIN_PATH) migrate up; \
	else \
		echo "$(YELLOW)⚠️  migrations directory not found$(NC)"; \
	fi

migrate-down: ## Reverts database migrations
	@echo "$(YELLOW)⚠️  Reverting migrations...$(NC)"
	@if [ -d "migrations" ]; then \
		$(GO) run $(MAIN_PATH) migrate down; \
	else \
		echo "$(YELLOW)⚠️  migrations directory not found$(NC)"; \
	fi

##@ Docker

docker-build: ## Build Docker image
	@echo "$(GREEN)🐳 Building Docker image...$(NC)"
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-run: ## Run Docker container
	@echo "$(GREEN)🐳 Starting container...$(NC)"
	docker run -p 4000:4000 --env-file .env $(APP_NAME):latest

docker-compose-up: ## Start all services with docker-compose
	@echo "$(GREEN)🐳 Starting services with docker-compose...$(NC)"
	docker-compose up -d

docker-compose-down: ## Stop all docker-compose services
	@echo "$(YELLOW)🐳 Stopping services...$(NC)"
	docker-compose down

docker-compose-logs: ## Show docker-compose logs
	docker-compose logs -f

##@ Formatting & Linting

fmt: ## Formats code
	@echo "$(GREEN)✨ Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

lint: ## Runs linter (requires golangci-lint)
	@echo "$(GREEN)🔍 Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✅ Lint complete$(NC)"; \
	else \
		echo "$(RED)❌ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		exit 1; \
	fi

vet: ## Runs go vet
	@echo "$(GREEN)🔍 Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✅ Vet complete$(NC)"

check: fmt vet lint test ## Runs all checks

##@ Cleanup

clean: ## Removes build files
	@echo "$(YELLOW)🧹 Cleaning build files...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Clean complete$(NC)"

clean-all: clean ## Removes build files and cache
	@echo "$(YELLOW)🧹 Complete clean (including cache)...$(NC)"
	$(GO) clean -cache -testcache -modcache
	@echo "$(GREEN)✅ Clean complete$(NC)"

##@ Utilities

setup: deps swagger ## Complete setup for development environment
	@echo "$(GREEN)🎉 Setup complete!$(NC)"
	@echo ""
	@echo "To start development, run:"
	@echo "  $(YELLOW)make dev$(NC)"
	@echo ""
	@echo "Other useful commands:"
	@echo "  $(YELLOW)make help$(NC)       - View all commands"
	@echo "  $(YELLOW)make test$(NC)       - Run tests"
	@echo "  $(YELLOW)make build$(NC)      - Compile application"

logs: ## Displays application logs (if running)
	@echo "$(GREEN)📋 Displaying logs...$(NC)"
	@if [ -f "logs/app.log" ]; then \
		tail -f logs/app.log; \
	else \
		echo "$(YELLOW)⚠️  Log file not found$(NC)"; \
	fi

version: ## Displays Go version and dependencies
	@echo "$(GREEN)📌 Versions:$(NC)"
	@$(GO) version
	@echo ""
	@echo "$(GREEN)Main dependencies:$(NC)"
	@$(GO) list -m all | grep -E '(whatsmeow|postgres|minio)'

status: ## Checks application status
	@echo "$(GREEN)🔍 Checking status...$(NC)"
	@curl -s http://localhost:4000/health || echo "$(RED)❌ Application is not running$(NC)"

##@ Desenvolvimento Avançado

profile-cpu: ## Profile de CPU (requer aplicação rodando)
	@echo "$(GREEN)📊 Capturando profile de CPU...$(NC)"
	curl http://localhost:4000/debug/pprof/profile?seconds=30 > cpu.prof
	$(GO) tool pprof -http=:8080 cpu.prof

profile-mem: ## Profile de memória (requer aplicação rodando)
	@echo "$(GREEN)📊 Capturando profile de memória...$(NC)"
	curl http://localhost:4000/debug/pprof/heap > mem.prof
	$(GO) tool pprof -http=:8080 mem.prof

generate: ## Roda go generate
	@echo "$(GREEN)⚙️  Executando go generate...$(NC)"
	$(GO) generate ./...

mod-graph: ## Exibe gráfico de dependências
	@echo "$(GREEN)📊 Gráfico de dependências:$(NC)"
	$(GO) mod graph
