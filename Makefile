# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy gen add-copyright format lint build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/coding-hui/iam
VERSION_PACKAGE=github.com/coding-hui/common/version

# ==============================================================================
# Includes

include hack/makelib/common.mk # make sure include common.mk at the first include line
include hack/makelib/golang.mk
include hack/makelib/image.mk
include hack/makelib/deploy.mk
include hack/makelib/copyright.mk
include hack/makelib/gen.mk
include hack/makelib/ca.mk
include hack/makelib/release.mk
include hack/makelib/swagger.mk
include hack/makelib/dependencies.mk
include hack/makelib/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="iam-apiserver iam-authz-server"
  IMAGES           Backend images to make. Default is all of cmd.
                   This option is available when using: make image/image.multiarch/push/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-authz-server"
  REGISTRY_PREFIX  Docker registry prefix. Default is wecoding.
                   Example: make push REGISTRY_PREFIX=devops-wecoding-docker.pkg.coding.net/wecoding/docker-repo VERSION=v1.6.2
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-pump" PLATFORMS="linux/amd64,linux/arm64"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.build BUILDX_OUTPUT_TYPE=registry

## deploy: Deploy updated components to development env.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.k8s

## undeploy: Undeploy components.
.PHONY: undeploy
undeploy:
	@helm uninstall iam -n iam-system

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## release: Release iam
.PHONY: release
release:
	@$(MAKE) release.run

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports tools.verify.swag
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' ! -path 'web/*' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' ! -path 'web/*' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' ! -path 'web/*' | $(XARGS) golines -w --max-len=180 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt
	@swag fmt -d ${ROOT_DIR}/

## verify-copyright: Verify the boilerplate headers for all files.
.PHONY: verify-copyright
verify-copyright:
	@$(MAKE) copyright.verify

## add-copyright: Ensures source code files have copyright license headers.
.PHONY: add-copyright
add-copyright:
	@$(MAKE) copyright.add

## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## ca: Generate CA files for all iam components.
.PHONY: ca
ca:
	@$(MAKE) gen.ca

## swag: Generate swagger document.
.PHONY: swag
swag:
	@$(MAKE) swag.run

## serve-swagger: Serve swagger spec and docs.
.PHONY: swagger.serve
serve-swagger:
	@$(MAKE) swagger.serve

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: Install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

## tidy: Go mod tidy
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
