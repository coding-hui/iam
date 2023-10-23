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

.PHONY: deploy.run.all
deploy.run.all:
	@echo "===========> Deploying all"
	@$(MAKE) deploy.run

.PHONY: deploy.run
deploy.run: $(addprefix deploy.run., $(DEPLOYS))

.PHONY: deploy.run.%
deploy.run.%:
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Deploying $* $(VERSION)-$(ARCH)"
	@${ROOT_DIR}/hack/deploy-iam.sh
