# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper for installing IAM services
#
# Usage:
#   make install                    # Install with SQLite + Redis (local development)
#   make install INSTALL_MODE=local # Same as above
#   make install INSTALL_MODE=docker # Install with Docker
#   make install INSTALL_MODE=k8s   # Install with Kubernetes
#   make install INSTALL_MODE=all   # Install all components
#   make start                       # Start IAM services
#   make stop                        # Stop IAM services
#   make status                      # Check service status
#   make restart                     # Restart IAM services
#   make logs                        # Show IAM logs
#   make clean.install               # Clean up installed files

# ==============================================================================
# Configuration

INSTALL_MODE ?= local

# ==============================================================================
# Install targets

## install.local: Install iam-apiserver with SQLite + Redis for local development
.PHONY: install.local
install.local: prerequisites.local
	@$(ROOT_DIR)/hack/install/install.sh install
	@echo "========================================"
	@echo "  IAM installed successfully!"
	@echo ""
	@echo "  Services"
	@echo "  --------"
	@echo "  API Server   http://127.0.0.1:8080"
	@echo "  Secure API   https://127.0.0.1:8443"
	@echo "  gRPC Server  127.0.0.1:8081"
	@echo ""
	@echo "  Directories"
	@echo "  ----------"
	@echo "  ~/.iam/conf/"
	@echo "  ~/.iam/data/"
	@echo "  ~/.iam/logs/"
	@echo ""
	@echo "  Quick Start"
	@echo "  ----------"
	@echo "  ~/.local/bin/iamctl user list"
	@echo "  make logs"
	@echo "========================================"

## install.k8s: Install IAM with Kubernetes
.PHONY: install.k8s
install.k8s: prerequisites.k8s
	@echo "==> Installing IAM with Kubernetes..."
	@echo "This feature is not implemented yet."
	@echo "Please use 'make install INSTALL_MODE=docker' or 'make install INSTALL_MODE=local' instead."

## install.docker: Install IAM with Docker
.PHONY: install.docker
install.docker: prerequisites.docker
	@echo "==> Installing IAM with Docker..."
	@echo "This feature is not implemented yet."
	@echo "Please use 'make install INSTALL_MODE=local' instead."

# ==============================================================================
# Service management targets

## start: Start IAM services
.PHONY: start
start:
	@$(ROOT_DIR)/hack/install/install.sh start

## stop: Stop IAM services
.PHONY: stop
stop:
	@$(ROOT_DIR)/hack/install/install.sh stop

## status: Check IAM services status
.PHONY: status
status:
	@$(ROOT_DIR)/hack/install/install.sh status

## restart: Restart IAM services
.PHONY: restart
restart:
	@$(ROOT_DIR)/hack/install/install.sh restart

## logs: Show IAM logs
.PHONY: logs
logs:
	@$(ROOT_DIR)/hack/install/install.sh logs

# ==============================================================================
# Uninstall target

## uninstall: Uninstall IAM services
.PHONY: uninstall
uninstall: prerequisites.local
	@echo "Uninstalling IAM services..."
	@$(ROOT_DIR)/hack/install/install.sh uninstall
	@echo "Uninstall completed"

## clean.install: Clean up installed files (但不卸载服务)
.PHONY: clean.install
clean.install:
	@echo "Cleaning installed files..."
	@rm -rf $(shell bash -c 'source $(ROOT_DIR)/hack/install/common.sh 2>/dev/null; echo $${IAM_INSTALL_DIR:-$$HOME/.iam}')
	@echo "Clean completed"

# ==============================================================================
# Prerequisites

## prerequisites.local: Check prerequisites for local development
.PHONY: prerequisites.local
prerequisites.local:
	@echo "==> Checking prerequisites for local installation..."
	@command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed."; exit 1; }
	@echo "==> Prerequisites check passed"

## prerequisites.docker: Check prerequisites for Docker installation
.PHONY: prerequisites.docker
prerequisites.docker:
	@echo "==> Checking prerequisites for Docker installation..."
	@command -v docker >/dev/null 2>&1 || { echo "Error: Docker is not installed."; exit 1; }
	@echo "==> Prerequisites check passed"

## prerequisites.k8s: Check prerequisites for Kubernetes installation
.PHONY: prerequisites.k8s
prerequisites.k8s:
	@echo "==> Checking prerequisites for Kubernetes installation..."
	@command -v kubectl >/dev/null 2>&1 || { echo "Error: kubectl is not installed."; exit 1; }
	@command -v helm >/dev/null 2>&1 || { echo "Error: helm is not installed."; exit 1; }
	@echo "==> Prerequisites check passed"
