# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for generate necessary files
#

PROTOC_INC_PATH=$(dir $(shell which protoc 2>/dev/null))/../include
API_DEPS=pkg/api/proto/apiserver/v1/cache.proto
API_DEPSRCS=$(API_DEPS:.proto=.pb.go)

.PHONY: gen.run
#gen.run: gen.errcode
gen.run: gen.clean gen.errcode

.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating iam error code go source files"
	@codegen -type=int ${ROOT_DIR}/pkg/code

.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation"
	@codegen -type=int -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/pkg/code

.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@${ROOT_DIR}/hack/gencerts.sh generate-iam-cert $(OUTPUT_DIR)/cert $(CA)

.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))

.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/hack/gen_default_config.sh

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name '*_generated.go' -delete

$(API_DEPSRCS): $(API_DEPS)
	@echo "===========> Generate protobuf files"
	@protoc --go_out=$(OUTPUT_DIR) --go-grpc_out=$(OUTPUT_DIR) $(@:.pb.go=.proto)
	@cp $(OUTPUT_DIR)/$(ROOT_PACKAGE)/$(dir $@)/*.go $(dir $@)
	@rm -rf $(OUTPUT_DIR)/$(ROOT_PACKAGE)

.PHONY: gen.proto
 gen.proto: $(API_DEPSRCS)
