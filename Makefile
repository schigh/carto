.DEFAULT_GOAL := help
_YELLOW	      =\033[0;33m
_NC           =\033[0m
PKG_DIRS      = $(shell go list -f '{{.Dir}}' ./...)
NOW           := $(shell date -u '+%Y%m%d_%I:%M:%S')
VERSION       = 0.0.1

.PHONY: help build build_release install fmt imports templates test_carto test_gen tests # generic commands
help: ## prints this help
	@grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

build: templates ## build carto
	@echo "${_YELLOW}building carto main...${_NC}"
	@go build -a -ldflags "-s -w -X \"main.Version=${NOW}\"" -o carto carto.go

build_release: templates ## build for release
	@echo "${_YELLOW}building carto main...${_NC}"
	@go build -a -ldflags "-s -w -X \"main.Version=${VERSION}\"" -o carto carto.go

install: build ## install dev build to $GOBIN
	@echo "${_YELLOW}installing carto...${_NC}"
	@mv carto ${GOBIN}/carto

install_release: build_release ## install release build to $GOBIN
	@echo "${_YELLOW}installing carto release...${_NC}"
	@mv carto ${GOBIN}/carto

fmt: ## apply the gofmt tool
	@echo "${_YELLOW}cleaning formatting...${_NC}"
	@gofmt -s -w .

imports: ## apply imports
	@echo "${_YELLOW}cleaning imports...${_NC}"
	@goimports -w -local github.com/schigh/carto .

templates: fmt imports ## generate templates
	@echo "${_YELLOW}generating templates...${_NC}"
	@go run tools/tmpl/main.go

test_carto: build ## test carto main
	@echo "${_YELLOW}testing carto main...${_NC}"
	@go test -v -cover github.com/schigh/carto

test_gen: ## test carto-generated code
	@echo "${_YELLOW}testing carto generated structs...${_NC}"
	@go test -race -v -cover github.com/schigh/carto/_tests

tests: build ## make tests for generated code
	@echo "${_YELLOW}generating test structs...${_NC}"
	@go run tools/tests/main.go
