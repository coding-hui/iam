# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for deploy to developer env
#

KUBECTL := kubectl
NAMESPACE ?= iam-system
CONTEXT ?= wecoding.dev

DEPLOYS=iam-apiserver

.PHONY: deploy.k8s.all
deploy.k8s.all:
	@echo "===========> Deploying all"
	@$(MAKE) deploy.k8s

.PHONY: deploy.restart
deploy.restart: $(addprefix deploy.restart., $(DEPLOYS))

.PHONY: deploy.restart.%
deploy.restart.%:
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Restarting $* $(VERSION)-$(ARCH)"
	@${ROOT_DIR}/hack/restart-iam.sh

.PHONY: deploy.k8s
deploy.k8s: $(addprefix deploy.k8s., $(DEPLOYS))

.PHONY: deploy.k8s.%
deploy.k8s.%:
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Deploying $* $(VERSION)-$(ARCH) to kubernetes"
	@${ROOT_DIR}/hack/deploy-iam.sh
