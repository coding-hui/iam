# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper for installing IAM services
#
# Usage:
#   make install                    # Install with SQLite + Redis (local development)
#   make install INSTALL_MODE=local   # Same as above
#   make install INSTALL_MODE=docker # Install with Docker
#   make install INSTALL_MODE=k8s    # Install with Kubernetes
#   make install INSTALL_MODE=all     # Install all components
#   make stop                          # Stop all IAM services
#   make clean.install                 # Clean up installed files

# ==============================================================================
# Configuration

INSTALL_MODE ?= local

# ==============================================================================
# Install targets

## install.local: Install iam-apiserver with SQLite + Redis for local development
.PHONY: install.local
install.local: prerequisites.local
	@echo "===========> Starting Redis..."
	@if command -v brew &> /dev/null; then \
		brew services start redis 2>/dev/null || redis-server --daemonize yes; \
	elif command -v redis-server &> /dev/null; then \
		redis-server --daemonize yes 2>/dev/null || true; \
	else \
		echo "Redis not found. Please install Redis first."; \
	fi
	@echo "===========> Building iam-apiserver..."
	@$(MAKE) go.build BINS=iam-apiserver
	@echo "===========> Starting iam-apiserver..."
	@$(LOCAL_OUTPUT_ROOT)/platforms/$(GOOS)/$(GOARCH)/iam-apiserver -c $(ROOT_DIR)/configs/iam-apiserver-sqlite.yaml &> $(LOCAL_OUTPUT_ROOT)/iam-apiserver.log &
	@sleep 2
	@echo "===========> iam-apiserver started successfully"
	@echo "Log file: $(LOCAL_OUTPUT_ROOT)/iam-apiserver.log"
	@echo "API Server: http://127.0.0.1:8000"

## install.docker: Install IAM with Docker
.PHONY: install.docker
install.docker: prerequisites.docker
	@$(ROOT_DIR)/hack/install/install.sh iam::install::install_storage
	@$(ROOT_DIR)/hack/install/install.sh iam::install::prepare_iam
	@$(ROOT_DIR)/hack/install/install.sh iam::apiserver::install

## install.k8s: Install IAM with Kubernetes
.PHONY: install.k8s
install.k8s: prerequisites.k8s
	@echo "===========> Installing IAM with Kubernetes..."
	@echo "This feature is not implemented yet."
	@echo "Please use 'make install INSTALL_MODE=docker' or 'make install INSTALL_MODE=local' instead."

## install.all: Install all IAM components
.PHONY: install.all
install.all: prerequisites.all
	@$(ROOT_DIR)/hack/install/install.sh iam::install::install_iam

## stop: Stop all IAM services
.PHONY: stop
stop:
	@echo "===========> Stopping IAM services..."
	@pkill -f iam-apiserver 2>/dev/null || true
	@pkill -f iam-authz-server 2>/dev/null || true
	@echo "===========> IAM services stopped"

## clean.install: Clean up installed files
.PHONY: clean.install
clean.install:
	@echo "===========> Cleaning up IAM installation..."
	@rm -rf $(LOCAL_OUTPUT_ROOT)/iam-apiserver.log
	@rm -rf $(LOCAL_OUTPUT_ROOT)/platforms/iam-apiserver
	@$(ROOT_DIR)/hack/install/install.sh iam::install::uninstall_iam 2>/dev/null || true
	@echo "===========> Cleanup completed"

## prerequisites.local: Check prerequisites for local development
.PHONY: prerequisites.local
prerequisites.local:
	@echo "===========> Checking prerequisites for local installation..."
	@command -v go &> /dev/null || { echo "Error: Go is not installed."; exit 1; }
	@command -v redis-server &> /dev/null || { echo "Error: Redis is not installed."; exit 1; }
	@echo "===========> Prerequisites check passed"

## prerequisites.docker: Check prerequisites for Docker installation
.PHONY: prerequisites.docker
prerequisites.docker:
	@echo "===========> Checking prerequisites for Docker installation..."
	@command -v docker &> /dev/null || { echo "Error: Docker is not installed."; exit 1; }
	@echo "===========> Prerequisites check passed"

## prerequisites.k8s: Check prerequisites for Kubernetes installation
.PHONY: prerequisites.k8s
prerequisites.k8s:
	@echo "===========> Checking prerequisites for Kubernetes installation..."
	@command -v kubectl &> /dev/null || { echo "Error: kubectl is not installed."; exit 1; }
	@command -v helm &> /dev/null || { echo "Error: helm is not installed."; exit 1; }
	@echo "===========> Prerequisites check passed"

## prerequisites.all: Check prerequisites for full installation
.PHONY: prerequisites.all
prerequisites.all:
	@echo "===========> Checking prerequisites for full installation..."
	@command -v docker &> /dev/null || { echo "Error: Docker is not installed."; exit 1; }
	@echo "===========> Prerequisites check passed"
