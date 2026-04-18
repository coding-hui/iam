# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

.DEFAULT_GOAL := all

# ==============================================================================
# Variables

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
OUTPUT_DIR := $(ROOT_DIR)/_output
ROOT_PACKAGE := github.com/coding-hui/iam
VERSION_PACKAGE := github.com/coding-hui/common/version

GO := go
GO_OUT_EXT :=
ifeq ($(GOOS), windows)
	GO_OUT_EXT := .exe
endif

# Version info from git
VERSION := $(shell git describe --tags --always --match='v*' 2>/dev/null || echo "v0.0.0")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_TREE_STATE := $(shell git status --porcelain 2>/dev/null || echo "")
ifeq ($(GIT_TREE_STATE),)
	GIT_TREE_STATE := clean
else
	GIT_TREE_STATE := dirty
endif

GO_LDFLAGS := -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_BUILD_FLAGS := -ldflags "$(GO_LDFLAGS)"

# Binaries to build
COMMANDS := $(filter-out %.md, $(wildcard $(ROOT_DIR)/cmd/*))
BINS := $(foreach cmd,$(COMMANDS),$(notdir $(cmd)))

# Minimum coverage
COVERAGE ?= 60

# ==============================================================================
# Targets

## all: Full pipeline (tidy, gen, format, lint, build)
.PHONY: all
all: tidy gen format lint build

## build: Build binaries for host platform
.PHONY: build
build:
	@echo "==> Building binaries"
	@mkdir -p $(OUTPUT_DIR)/bin
	@$(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/bin/apiserver$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/apiserver

## test: Run unit tests with race detection and coverage
.PHONY: test
test:
	@echo "==> Running tests"
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out -timeout=10m -shuffle=on ./... 2>&1 || true
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html 2>/dev/null || true

## lint: Run golangci-lint
.PHONY: lint
lint:
	@$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 2>/dev/null || true
	@echo "==> Linting"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

## format: Format Go source files
.PHONY: format
format:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest 2>/dev/null || true
	@$(GO) install github.com/segmentio/golines@latest 2>/dev/null || true
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest 2>/dev/null || true
	@echo "==> Formatting codes"
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' ! -path '$(ROOT_DIR)/_output/*' | \
		xargs -r gofmt -s -w
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' ! -path '$(ROOT_DIR)/_output/*' | \
		xargs -r goimports -w -local $(ROOT_PACKAGE)
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' ! -path '$(ROOT_DIR)/_output/*' | \
		xargs -r golines -w --max-len=180 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt
	@swag fmt -d $(ROOT_DIR)/

## tidy: Go mod tidy
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## gen: Generate code (error codes)
.PHONY: gen
gen:
	@$(GO) install $(ROOT_DIR)/tools/codegen/codegen.go 2>/dev/null || true
	@echo "==> Generating error codes"
	@codegen -type=int $(ROOT_DIR)/pkg/code

## verify-copyright: Verify license headers
.PHONY: verify-copyright
verify-copyright:
	@$(GO) install github.com/coding-hui/addlicense@latest 2>/dev/null || true
	@echo "==> Verifying license headers"
	@addlicense --check -f $(ROOT_DIR)/hack/boilerplate.txt \
		--skip-dirs "api/*" --skip-dirs "installer/*" --skip-dirs "web/*" \
		--skip-dirs "_output/*" --skip-dirs ".idea/*" $(ROOT_DIR)

## add-copyright: Add license headers
.PHONY: add-copyright
add-copyright:
	@$(GO) install github.com/coding-hui/addlicense@latest 2>/dev/null || true
	@echo "==> Adding license headers"
	@addlicense -v -f $(ROOT_DIR)/hack/boilerplate.txt \
		--skip-dirs "api/*" --skip-dirs "installer/*" --skip-dirs "web/*" \
		--skip-dirs "_output/*" --skip-dirs ".idea/*" $(ROOT_DIR)

## clean: Remove build output
.PHONY: clean
clean:
	@echo "==> Cleaning build output"
	@rm -vrf $(OUTPUT_DIR)

## image: Build Docker image
.PHONY: image
image:
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest 2>/dev/null || true
	@echo "==> Building Docker image"
	@docker build -t $(REGISTRY_PREFIX)/apiserver:$(VERSION) \
		-f $(ROOT_DIR)/installer/dockerfile/apiserver/Dockerfile $(ROOT_DIR)

## help: Show this help
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'