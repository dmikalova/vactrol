# Vactrol developer tasks. Run `make help` to list targets.

.PHONY: help run build test cover vet fmt tidy generate-comments

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

run: ## Launch the vactrol TUI (card explorer / play a hotseat game)
	go run ./cmd/tui

build: ## Build all packages
	go build ./...

test: ## Run all tests
	go test ./...

cover: ## Run tests and report total coverage
	go test ./internal/engine/ -coverprofile=coverage.out
	@go tool cover -func=coverage.out | tail -1

vet: ## Run go vet
	go vet ./...

fmt: ## Format all Go files
	gofmt -w .

tidy: ## Tidy module dependencies
	go mod tidy

generate-comments: ## Rewrite each card's doc comment from its definition
	go run ./cmd/gencomments
