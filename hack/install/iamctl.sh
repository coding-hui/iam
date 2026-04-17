#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::iamctl::info() {
cat << EOF
iamctl test command: iamctl user list
EOF
}

# 安装
function iam::iamctl::install()
{
  pushd ${IAM_ROOT} > /dev/null

  iam::log::section "Installing iamctl"

  # 1. 生成并安装 CA 证书和私钥
  iam::log::substep "Generating certificates..."
  ./hack/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert
  mkdir -p "${IAM_INSTALL_DIR}" 2>/dev/null
  mkdir -p "${IAM_CONFIG_DIR}/cert" 2>/dev/null
  cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${IAM_CONFIG_DIR}/cert 2>/dev/null

  ./hack/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert admin
  cert_dir=$(dirname ${CONFIG_USER_CLIENT_CERTIFICATE})
  key_dir=$(dirname ${CONFIG_USER_CLIENT_KEY})
  mkdir -p ${cert_dir} ${key_dir} 2>/dev/null
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin.pem ${CONFIG_USER_CLIENT_CERTIFICATE} 2>/dev/null
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin-key.pem ${CONFIG_USER_CLIENT_KEY} 2>/dev/null

  # 2. 构建 iamctl (包含版本信息)
  iam::log::substep "Building iamctl..."
  local version=$(git describe --tags --always --match='v*' 2>/dev/null || echo "v0.0.0")
  local git_commit=$(git rev-parse HEAD 2>/dev/null || echo "")
  local git_tree_state="dirty"
  if [ -z "$(git status --porcelain 2>/dev/null)" ]; then
    git_tree_state="clean"
  fi
  local build_date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
  CGO_ENABLED=1 go build -ldflags "\
    -X github.com/coding-hui/common/version.GitVersion=${version} \
    -X github.com/coding-hui/common/version.GitCommit=${git_commit} \
    -X github.com/coding-hui/common/version.GitTreeState=${git_tree_state} \
    -X github.com/coding-hui/common/version.BuildDate=${build_date}" \
    -o ${LOCAL_OUTPUT_ROOT}/bin/iamctl github.com/coding-hui/iam/cmd/iamctl
  local bin_path=$(iam::common::get_bin_path)
  mkdir -p "${bin_path}" 2>/dev/null
  cp ${LOCAL_OUTPUT_ROOT}/bin/iamctl "${bin_path}/iamctl" 2>/dev/null

  # 3. 生成并安装 iamctl 的配置文件（iamctl.yaml）
  mkdir -p $HOME/.iam
  envsubst < ${IAM_ROOT}/configs/iamctl-template.yaml > $HOME/.iam/iamctl.yaml

  iam::iamctl::status || return 1

  popd >/dev/null 2>&1
}

# 卸载
function iam::iamctl::uninstall()
{
  iam::log::section "Uninstalling iamctl"
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)
  iam::log::substep "Removing binary..."
  rm -f "${bin_path}/iamctl"
  iam::log::substep "Removing config..."
  rm -f $HOME/.iam/iamctl.yaml
  #iam::common::sudo "rm -f ${IAM_CONFIG_DIR}/cert/admin*pem"
  rm -f ${CONFIG_USER_CLIENT_CERTIFICATE}
  rm -f ${CONFIG_USER_CLIENT_KEY}
  set -o errexit

  iam::log::info "Uninstall iamctl successfully"
}

# 状态检查
function iam::iamctl::status()
{
  local bin_path=$(iam::common::get_bin_path)
  "${bin_path}/iamctl" user list --iamconfig=$HOME/.iam/iamctl.yaml 2>/dev/null | grep -q ADMIN || {
   iam::log::error "cannot list user, iamctl maybe not installed properly"
   return 1
  }
}

if [[ "$*" =~ iam::iamctl:: ]];then
  :
fi
