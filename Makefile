# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

.DEFAULT_GOAL := all

# ==============================================================================
# Variables

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
OUTPUT_DIR := $(ROOT_DIR)/_output
TOOLS_DIR := $(OUTPUT_DIR)/tools
ROOT_PACKAGE := github.com/coding-hui/iam
VERSION_PACKAGE := github.com/coding-hui/common/version

GO := go
CGO_ENABLED ?= 1
PLATFORM ?= $(GOOS)_$(GOARCH)
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
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

# Platforms for multi-arch build
PLATFORMS ?= linux/amd64 linux/arm64

# ==============================================================================
# Usage

define USAGE_OPTIONS
Options:
  BINS            Binaries to build. Default is all cmd/*.
  PLATFORMS       Multi-platform build. Default: linux/amd64 linux/arm64.
  VERSION         Version info compiled into binaries.
  COVERAGE        Minimum test coverage. Default: 60.
  CGO_ENABLED     Enable CGO. Default: 1.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## all: Full pipeline (tidy, gen, copyright, format, lint, build)
.PHONY: all
all: tidy gen add-copyright format lint build

## build: Build binaries for host platform
.PHONY: build
build: go.build

## build.multiarch: Build binaries for multiple platforms
.PHONY: build.multiarch
build.multiarch: go.build.multiarch

## test: Run unit tests with race detection and coverage
.PHONY: test
test: go.test

## cover: Run tests and generate HTML coverage report
.PHONY: cover
cover: go.test.cover

## lint: Run golangci-lint
.PHONY: lint
lint: go.lint

## format: Format Go source files
.PHONY: format
format: tools.verify.goimports tools.verify.golines tools.verify.swag
	@echo "==> Formatting codes"
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' | \
		xargs -r gofmt -s -w
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' | \
		xargs -r goimports -w -local $(ROOT_PACKAGE)
	@find $(ROOT_DIR) -type f -name '*.go' ! -path '$(ROOT_DIR)/web/*' ! -path '$(ROOT_DIR)/vendor/*' | \
		xargs -r golines -w --max-len=180 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt
	@swag fmt -d $(ROOT_DIR)/

## tidy: Go mod tidy
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## gen: Generate code (error codes)
.PHONY: gen
gen: gen.errcode

## gen.errcode: Generate error code files
.PHONY: gen.errcode
gen.errcode: tools.verify.codegen
	@echo "==> Generating error codes"
	@codegen -type=int $(ROOT_DIR)/pkg/code

## swag: Generate Swagger docs
.PHONY: swag
swag: tools.verify.swag
	@echo "==> Generating swagger docs"
	@swag i -g apiserver.go -dir $(ROOT_DIR)/internal/apiserver --parseDependency --parseInternal -o $(ROOT_DIR)/api/swagger

## serve-swagger: Serve Swagger UI
.PHONY: serve-swagger
serve-swagger: tools.verify.swag
	@swag serve -F=redoc --no-open --port 36666 $(ROOT_DIR)/api/swagger/swagger.yaml

## verify-copyright: Verify license headers
.PHONY: verify-copyright
verify-copyright: tools.verify.addlicense
	@echo "==> Verifying license headers"
	@addlicense --check -f $(ROOT_DIR)/hack/boilerplate.txt \
		--skip-dirs "api/*" --skip-dirs "installer/*" --skip-dirs "web/node_modules/*" \
		--skip-dirs "_output/*" $(ROOT_DIR)

## add-copyright: Add license headers
.PHONY: add-copyright
add-copyright: tools.verify.addlicense
	@echo "==> Adding license headers"
	@addlicense -v -f $(ROOT_DIR)/hack/boilerplate.txt \
		--skip-dirs "api/*" --skip-dirs "web/*" --skip-dirs "installer/*" \
		--skip-dirs "_output/*" $(ROOT_DIR)

## clean: Remove build output
.PHONY: clean
clean:
	@echo "==> Cleaning build output"
	@rm -vrf $(OUTPUT_DIR)

## image: Build Docker image for host platform
.PHONY: image
image: tools.verify.swag
	@echo "==> Building Docker image"
	@docker build -t $(REGISTRY_PREFIX)/iam-apiserver:$(VERSION) \
		-f $(ROOT_DIR)/installer/dockerfile/iam-apiserver/Dockerfile $(ROOT_DIR)

## image.multiarch: Build multi-arch Docker image
.PHONY: image.multiarch
image.multiarch:
	@echo "==> Building multi-arch Docker image (requires docker buildx)"
	@docker buildx create --use 2>/dev/null || true
	@docker buildx build --platform $(PLATFORMS) \
		-t $(REGISTRY_PREFIX)/iam-apiserver:$(VERSION) \
		--output=type=registry \
		-f $(ROOT_DIR)/installer/dockerfile/iam-apiserver/Dockerfile $(ROOT_DIR)

## deploy: Deploy to Kubernetes dev environment
.PHONY: deploy
deploy:
	@$(ROOT_DIR)/hack/deploy-iam.sh

## undeploy: Remove Kubernetes deployment
.PHONY: undeploy
undeploy:
	@helm uninstall iam -n iam-system 2>/dev/null || true

## install: Install IAM services locally (use INSTALL_MODE=local/docker/k8s)
.PHONY: install
install:
	@$(ROOT_DIR)/hack/install.sh install

## start: Start IAM services
.PHONY: start
start:
	@$(ROOT_DIR)/hack/install.sh start

## stop: Stop IAM services
.PHONY: stop
stop:
	@$(ROOT_DIR)/hack/install.sh stop

## status: Check IAM services status
.PHONY: status
status:
	@$(ROOT_DIR)/hack/install.sh status

## restart: Restart IAM services
.PHONY: restart
restart:
	@$(ROOT_DIR)/hack/install.sh restart

## logs: Show IAM logs
.PHONY: logs
logs:
	@$(ROOT_DIR)/hack/install.sh logs

## uninstall: Uninstall IAM services
.PHONY: uninstall
uninstall:
	@$(ROOT_DIR)/hack/install.sh uninstall

## check-updates: Check outdated Go dependencies
.PHONY: check-updates
check-updates:
	@$(GO) list -u -m -json all | \
		$(TOOLS_DIR)/go-mod-outdated@latest -update -direct 2>/dev/null || \
		echo "Install go-mod-outdated: go install github.com/psampaz/go-mod-outdated@latest"

## tools: Install build tools
.PHONY: tools
tools: tools.verify.golangci-lint tools.verify.addlicense tools.verify.goimports \
	tools.verify.golines tools.verify.swag tools.verify.go-junit-report tools.verify.codegen

## help: Show this help
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"

# ==============================================================================
# Go build targets

.PHONY: go.build go.build.multiarch go.lint go.test go.test.cover

go.build: $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "==> Building $(COMMAND) $(VERSION) for $(OS)/$(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/bin
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(OS) GOARCH=$(ARCH) $(GO) build \
		$(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/bin/$(COMMAND)$(GO_OUT_EXT) \
		$(ROOT_PACKAGE)/cmd/$(COMMAND)

go.build.multiarch: $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

go.lint: tools.verify.golangci-lint
	@echo "==> Linting"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

go.test:
	@echo "==> Running tests"
	@mkdir -p $(OUTPUT_DIR)
	@set -o pipefail; $(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short -v ./... 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml) || true
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html 2>/dev/null || true

go.test.cover: go.test
	@echo "==> Coverage report"
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/hack/coverage.awk

# ==============================================================================
# Tool verification and installation (inline, no makelib needed)

TOOLS_DIR := $(OUTPUT_DIR)/tools

.PHONY: tools.verify.golangci-lint tools.verify.addlicense tools.verify.goimports \
	tools.verify.golines tools.verify.swag tools.verify.go-junit-report tools.verify.codegen

tools.verify.%:
	@if ! which $* &>/dev/null; then \
		$(MAKE) install.$*; \
	fi

.PHONY: install.golangci-lint install.addlicense install.goimports install.golines \
	install.swag install.go-junit-report install.codegen

install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

install.addlicense:
	@$(GO) install github.com/coding-hui/addlicense@latest

install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

install.golines:
	@$(GO) install github.com/segmentio/golines@latest

install.swag:
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest

install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

install.codegen:
	@$(GO) install $(ROOT_DIR)/tools/codegen/codegen.go
