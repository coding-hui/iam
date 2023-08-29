# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for swagger
#

.PHONY: swag.run
swag.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	@swag i -g apiserver.go -dir ${ROOT_DIR}/internal/apiserver --parseDependency --parseInternal -o ${ROOT_DIR}/api/swagger

.PHONY: swag.fmt
swag.fmt: tools.verify.swagger
	@echo "===========> Format swag comments"
	@swag fmt -d ${ROOT_DIR}/

.PHONY: swag.serve
swag.serve: tools.verify.swagger
	@swag serve -F=redoc --no-open --port 36666 $(ROOT_DIR)/api/swagger/swagger.yaml
