.DEFAULT_GOAL := ci
_YELLOW	      =\033[0;33m
_NC           =\033[0m
PKG_DIRS      = $(shell go list -f '{{.Dir}}' ./...)
NOW           := $(shell date -u '+%Y%m%d_%I:%M:%S')
VERSION       = 0.0.1

.PHONY: help
help: ## prints this help
	@grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

.PHONY: setup
setup: # install deps
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@go mod tidy

.PHONY: lint
lint: ## runs the code linter
	@echo "${_YELLOW}linting...${_NC}"
	@golangci-lint run

.PHONY: build
build: templates ## build carto
	@echo "${_YELLOW}building carto main...${_NC}"
	@go build -a -ldflags "-s -w -X \"main.Version=${NOW}\"" -o carto carto.go

.PHONY: build_release
build_release: templates ## build for release
	@echo "${_YELLOW}building carto main...${_NC}"
	@go build -a -ldflags "-s -w -X \"main.Version=${VERSION}\"" -o carto carto.go

.PHONY: install
install: build ## install dev build to $GOBIN
	@echo "${_YELLOW}installing carto...${_NC}"
	@mv carto ${GOBIN}/carto

.PHONY: install_release
install_release: build_release ## install release build to $GOBIN
	@echo "${_YELLOW}installing carto release...${_NC}"
	@mv carto ${GOBIN}/carto

.PHONY: fmt
fmt: ## apply the gofmt tool
	@echo "${_YELLOW}cleaning formatting...${_NC}"
	@gofmt -s -w .

.PHONY: imports
imports: ## apply imports
	@echo "${_YELLOW}cleaning imports...${_NC}"
	@goimports -w -local github.com/schigh/carto .

.PHONY: templates
templates:  ## generate templates
	@echo "${_YELLOW}generating templates...${_NC}"
	@go run tools/tmpl/main.go
	@gofmt -s -w .
	@goimports -w -local github.com/schigh/carto .

.PHONY: test_carto
test_carto: build ## test carto main
	@echo "${_YELLOW}testing carto main...${_NC}"
	@go test -race -v -cover github.com/schigh/carto

.PHONY: test_gen
test_gen: tests ## test carto-generated code
	@echo "${_YELLOW}testing carto generated structs...${_NC}"
	@go test -race -v -cover github.com/schigh/carto/_tests

.PHONY: tests
tests: templates build ## make tests for generated code
	@echo "${_YELLOW}generating test structs...${_NC}"
	@go run tools/tests/main.go

.PHONY: ci
ci: test_carto test_gen ## builds for ci
	@echo "${_YELLOW}ci build complete${_NC}"
