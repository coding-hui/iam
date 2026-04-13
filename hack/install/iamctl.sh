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
  pushd ${IAM_ROOT}

  # 1. 生成并安装 CA 证书和私钥
  ./hack/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}/cert"
  cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${IAM_CONFIG_DIR}/cert

  ./hack/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert admin
  #iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/admin*pem ${IAM_CONFIG_DIR}/cert"
  cert_dir=$(dirname ${CONFIG_USER_CLIENT_CERTIFICATE})
  key_dir=$(dirname ${CONFIG_USER_CLIENT_KEY})
  mkdir -p ${cert_dir} ${key_dir}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin.pem ${CONFIG_USER_CLIENT_CERTIFICATE}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin-key.pem ${CONFIG_USER_CLIENT_KEY}

  # 2. 构建 iamctl
  CGO_ENABLED=1 go build -o ${LOCAL_OUTPUT_ROOT}/bin/iamctl github.com/coding-hui/iam/cmd/iamctl
  local bin_path=$(iam::common::get_bin_path)
  mkdir -p "${bin_path}"
  cp ${LOCAL_OUTPUT_ROOT}/bin/iamctl "${bin_path}/iamctl"

  # 3. 生成并安装 iamctl 的配置文件（iamctl.yaml）
  mkdir -p $HOME/.iam
  ./hack/genconfig.sh ${ENV_FILE} configs/iamctl-template.yaml > $HOME/.iam/iamctl.yaml 2>/dev/null || {
    # Fallback: create config manually if genconfig fails
    cat > $HOME/.iam/iamctl.yaml << 'EOFCONFIG'
apiVersion: v1
user:
  username: admin
  password: Admin@2021
  client-certificate: ${HOME}/.iam/cert/admin.pem
  client-key: ${HOME}/.iam/cert/admin-key.pem

server:
  address: https://127.0.0.1:8443
  timeout: 10s
  insecure-skip-tls-verify: true
EOFCONFIG
  }
  iam::iamctl::status || return 1
  iam::iamctl::info

  iam::log::info "install iamctl successfully"
  popd
}

# 卸载
function iam::iamctl::uninstall()
{
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)
  rm -f "${bin_path}/iamctl"
  rm -f $HOME/.iam/iamctl.yaml
  #iam::common::sudo "rm -f ${IAM_CONFIG_DIR}/cert/admin*pem"
  rm -f ${CONFIG_USER_CLIENT_CERTIFICATE}
  rm -f ${CONFIG_USER_CLIENT_KEY}
  set -o errexit

  iam::log::info "uninstall iamctl successfully"
}

# 状态检查
function iam::iamctl::status()
{
  iamctl user list | grep -q admin || {
   iam::log::error "cannot list user, iamctl maybe not installed properly"
   return 1
  }

 if echo | telnet ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>&1|grep refused &>/dev/null;then
   iam::log::error "cannot access insecure port, iamctl maybe not startup"
   return 1
 fi
}

if [[ "$*" =~ iam::iamctl:: ]];then
  :
fi
