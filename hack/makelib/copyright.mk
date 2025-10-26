# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for copyright
#

.PHONY: copyright.verify
copyright.verify: tools.verify.addlicense
	@echo "===========> Verifying the boilerplate headers for all files"
	@addlicense --check -f $(ROOT_DIR)/hack/boilerplate.txt --skip-dirs "api/*" --skip-dirs "installer/*" --skip-dirs "web/node_modules/*" --skip-dirs "_output/*" $(ROOT_DIR)

.PHONY: copyright.add
copyright.add: tools.verify.addlicense
	@addlicense -v -f $(ROOT_DIR)/hack/boilerplate.txt --skip-dirs "api/*" --skip-dirs "web/*" --skip-dirs "installer/*" --skip-dirs "_output/*" $(ROOT_DIR)
