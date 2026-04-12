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

# 判断当前操作系统是否为 macOS
function iam::common::is_macos() {
  [[ "$(uname -s)" == "Darwin" ]]
}

# 判断当前 Linux 发行版是否为 Ubuntu/Debian
function iam::common::is_ubuntu() {
  command -v apt-get &>/dev/null
}

# 获取二进制文件安装路径
# macOS: ~/.local/bin (用户可写，无需 sudo)
# Linux: ${IAM_INSTALL_DIR}/bin
function iam::common::get_bin_path() {
  if iam::common::is_macos; then
    echo "${HOME}/.local/bin"
  else
    echo "${IAM_INSTALL_DIR}/bin"
  fi
}

