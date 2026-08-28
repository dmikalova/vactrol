# Vactrol developer tasks. Run `make help` to list targets.

.PHONY: help run build test cover vet fmt fmt-check check tidy gen generate-comments generate-rules

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-14s %s\n", $$1, $$2}'

run: ## Launch the vactrol TUI (card explorer / play a hotseat game)
	go run ./cmd/tui

build: ## Build all packages
	go build ./...

test: ## Run all tests
	go test ./...

cover: ## Run tests and report engine coverage (kept at 100%)
	go test ./internal/engine/ -coverprofile=coverage.out
	@go tool cover -func=coverage.out | tail -1

vet: ## Run go vet
	go vet ./...

fmt: ## Format all Go files
	gofmt -w .

fmt-check: ## Fail if any Go file is not gofmt-clean (does not modify files)
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then \
		echo "gofmt needed:"; echo "$$out"; exit 1; fi

check: fmt-check build vet test cover ## Full green gate: fmt-check, build, vet, test, coverage
	@echo "ALL GREEN"

tidy: ## Tidy module dependencies
	go mod tidy

generate-comments: ## Rewrite each card's doc comment from its definition
	go run ./cmd/gencomments

generate-rules: ## Generate docs/rulebook.md from engine doc comments
	go run ./cmd/genrules

gen: generate-comments generate-rules ## Regenerate card comments and the rulebook
