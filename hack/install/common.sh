#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Sourced flag
COMMON_SOURCED=true

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${IAM_ROOT}/hack/lib/init.sh"
source "${IAM_ROOT}/hack/install/environment.sh"

# 不输入密码执行需要 root 权限的命令
function iam::common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}

# ==============================================================================
# Aliases for backward compatibility
# These functions are now in hack/lib/util.sh
# ==============================================================================

function iam::common::is_macos() {
  iam::util::is_macos
}

function iam::common::is_ubuntu() {
  iam::util::is_ubuntu
}

function iam::common::get_bin_path() {
  iam::util::get_bin_path "$@"
}

function iam::common::get_sed_i_flag() {
  iam::util::get_sed_i_flag
}

function iam::common::launchd_load() {
  iam::util::launchd_load "$@"
}

function iam::common::launchd_unload() {
  iam::util::launchd_unload "$@"
}

function iam::common::launchd_is_running() {
  iam::util::launchd_is_running "$@"
}

function iam::common::brew_service_start() {
  iam::util::brew_service_start "$@"
}

function iam::common::brew_service_stop() {
  iam::util::brew_service_stop "$@"
}

function iam::common::pkill() {
  iam::util::pkill "$@"
}
