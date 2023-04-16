# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: clean

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/coding-hui/iam
VERSION_PACKAGE=github.com/coding-hui/iam/version

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

## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## ca: Generate CA files for all iam components.
.PHONY: ca
ca:
	@$(MAKE) gen.ca

## swagger: Generate swagger document.
.PHONY: swagger
swagger:
	@$(MAKE) swagger.run

## serve-swagger: Serve swagger spec and docs.
.PHONY: swagger.serve
serve-swagger:
	@$(MAKE) swagger.serve

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: clean
clean:
	rm -rf _output/
	rm -f *.log
