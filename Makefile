all: setup

help: ## Show this help
	@perl -nle 'BEGIN {printf "Usage:\n  make \033[33m<target>\033[0m\n\nTargets:\n"} printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 if /^([a-zA-Z_-].+):.*\s+## (.*)/' $(MAKEFILE_LIST)
.PHONY: help

setup: ## Install development tools
	go list -tags=tools -f='{{ join .Imports "\n" }}' ./tools.go | tr -d [ | tr -d ] | xargs -I{} go install {}
.PHONY: setup

mock-generate:
	go generate ./mock/github
.PHONY: mock

lint: ## Run golangci-lint
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.64.5 golangci-lint run --modules-download-mode=readonly /app/...
.PHONY: lint

lint-fix: ## Run golangci-lint with --fix
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.64.5 golangci-lint run --fix --modules-download-mode=readonly /app/...
.PHONY: lint-fix
