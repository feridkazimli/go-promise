.PHONY: help test test-verbose test-coverage benchmark benchmark-all lint fmt vet clean install-tools

# Default target
.DEFAULT_GOAL := help

# Colored output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Go commands
GOCMD := go
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Project settings
PACKAGE := promise
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

## help: Show help menu
help:
	@echo "$(CYAN)═══════════════════════════════════════════════════$(NC)"
	@echo "$(GREEN)  Go Promise Library - Makefile Commands$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(YELLOW)Development Commands:$(NC)"
	@echo "  make test              - Run all tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Generate coverage report"
	@echo "  make benchmark         - Run benchmarks"
	@echo "  make benchmark-all     - Run all benchmarks (detailed)"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint              - Run linting"
	@echo "  make fmt               - Format code"
	@echo "  make vet               - Run Go vet analysis"
	@echo "  make check             - Run all checks"
	@echo ""
	@echo "$(YELLOW)Cleanup and Setup:$(NC)"
	@echo "  make clean             - Clean temporary files"
	@echo "  make install-tools     - Install development tools"
	@echo "  make deps              - Update dependencies"
	@echo ""

## test: Run all tests
test:
	@echo "$(CYAN)🧪 Running tests...$(NC)"
	@$(GOTEST) -v ./... -timeout 30s
	@echo "$(GREEN)✓ Tests completed successfully!$(NC)"

## test-verbose: Run tests with detailed output
test-verbose:
	@echo "$(CYAN)🧪 Running detailed tests...$(NC)"
	@$(GOTEST) -v -race ./... -timeout 30s
	@echo "$(GREEN)✓ Detailed tests completed successfully!$(NC)"

## test-coverage: Generate coverage report
test-coverage:
	@echo "$(CYAN)📊 Generating coverage report...$(NC)"
	@$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@$(GOCMD) tool cover -func=$(COVERAGE_FILE)
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_HTML)$(NC)"
	@echo "$(YELLOW)To open in browser: open $(COVERAGE_HTML)$(NC)"

## test-short: Run quick tests (skip long-running ones)
test-short:
	@echo "$(CYAN)⚡ Running short tests...$(NC)"
	@$(GOTEST) -short ./...
	@echo "$(GREEN)✓ Short tests completed successfully!$(NC)"

## benchmark: Run benchmarks
benchmark:
	@echo "$(CYAN)🚀 Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -run=^# ./...
	@echo "$(GREEN)✓ Benchmarks completed!$(NC)"

## benchmark-all: Run all benchmarks in detail
benchmark-all:
	@echo "$(CYAN)🚀 Running detailed benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -benchtime=5s -run=^# ./...
	@echo "$(GREEN)✓ Detailed benchmarks completed!$(NC)"

## benchmark-compare: Compare benchmark results (requires benchstat)
benchmark-compare:
	@echo "$(CYAN)📈 Comparing benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -run=^# ./... > bench_new.txt
	@if [ -f bench_old.txt ]; then \
		benchstat bench_old.txt bench_new.txt; \
	else \
		echo "$(YELLOW)⚠ bench_old.txt not found. Saving as first run...$(NC)"; \
		cp bench_new.txt bench_old.txt; \
	fi

## lint: Run linting (requires golangci-lint)
lint:
	@echo "$(CYAN)🔍 Running lint checks...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting completed!$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not installed!$(NC)"; \
		echo "$(YELLOW)To install: make install-tools$(NC)"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "$(CYAN)✨ Formatting code...$(NC)"
	@$(GOFMT) ./...
	@echo "$(GREEN)✓ Code formatted!$(NC)"

## vet: Run Go vet analysis
vet:
	@echo "$(CYAN)🔬 Running Go vet analysis...$(NC)"
	@$(GOVET) ./...
	@echo "$(GREEN)✓ Vet analysis completed!$(NC)"

## check: Run all checks
check: fmt vet lint test
	@echo "$(GREEN)✓ All checks passed successfully!$(NC)"

## clean: Clean temporary files
clean:
	@echo "$(CYAN)🧹 Cleaning up...$(NC)"
	@$(GOCLEAN)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -f bench_old.txt bench_new.txt
	@echo "$(GREEN)✓ Cleanup completed!$(NC)"

## deps: Update dependencies
deps:
	@echo "$(CYAN)📦 Updating dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy
	@$(GOMOD) verify
	@echo "$(GREEN)✓ Dependencies updated!$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(CYAN)🔧 Installing development tools...$(NC)"
	@echo "$(YELLOW)Installing golangci-lint...$(NC)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
	@echo "$(YELLOW)Installing benchstat...$(NC)"
	@go install golang.org/x/perf/cmd/benchstat@latest
	@echo "$(GREEN)✓ Tools installed successfully!$(NC)"

## watch: Run tests in watch mode (requires fswatch)
watch:
	@if command -v fswatch >/dev/null 2>&1; then \
		echo "$(CYAN)👀 Starting test watch mode...$(NC)"; \
		fswatch -o . | xargs -n1 -I{} make test; \
	else \
		echo "$(RED)✗ fswatch not installed!$(NC)"; \
		echo "$(YELLOW)macOS: brew install fswatch$(NC)"; \
		echo "$(YELLOW)Linux: apt-get install fswatch$(NC)"; \
		exit 1; \
	fi

## profile-cpu: Run CPU profiling
profile-cpu:
	@echo "$(CYAN)📊 Running CPU profiling...$(NC)"
	@$(GOTEST) -cpuprofile=cpu.prof -bench=. ./...
	@go tool pprof -http=:8080 cpu.prof

## profile-mem: Run memory profiling
profile-mem:
	@echo "$(CYAN)📊 Running memory profiling...$(NC)"
	@$(GOTEST) -memprofile=mem.prof -bench=. ./...
	@go tool pprof -http=:8080 mem.prof

## mod-graph: Show module dependency graph
mod-graph:
	@echo "$(CYAN)📊 Module dependency graph:$(NC)"
	@$(GOMOD) graph

## version: Show Go and module versions
version:
	@echo "$(CYAN)📌 Version Information:$(NC)"
	@echo "Go version: $$(go version)"
	@echo "Module: $$(go list -m)"
