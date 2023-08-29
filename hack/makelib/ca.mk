# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for create CA files
#

.PHONY: ca.gen.%
ca.gen.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@${ROOT_DIR}/hack/gencerts.sh generate-iam-cert $(OUTPUT_DIR)/cert $(CA)

.PHONY: ca.gen
ca.gen: $(addprefix ca.gen., $(CERTIFICATES))
