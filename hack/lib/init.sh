#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
IAM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${IAM_ROOT}/hack/lib/util.sh"
source "${IAM_ROOT}/hack/lib/logging.sh"
source "${IAM_ROOT}/hack/lib/color.sh"

iam::log::install_errexit

source "${IAM_ROOT}/hack/lib/version.sh"
source "${IAM_ROOT}/hack/lib/golang.sh"
