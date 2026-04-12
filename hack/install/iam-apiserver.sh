#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::apiserver::info() {
cat << EOF
iam-apserver insecure listen on: ${IAM_APISERVER_HOST}:${IAM_APISERVER_INSECURE_BIND_PORT}
iam-apserver secure listen on: ${IAM_APISERVER_HOST}:${IAM_APISERVER_SECURE_BIND_PORT}
EOF
}

# 创建 launchd plist 文件（仅 macOS，用户级 Agent，无需 sudo）
function iam::apiserver::create_plist() {
  local bin_path=$(iam::common::get_bin_path)
  cat > "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.github.coding-hui.iam-apiserver</string>
    <key>ProgramArguments</key>
    <array>
        <string>${bin_path}/iam-apiserver</string>
        <string>--config</string>
        <string>${IAM_CONFIG_DIR}/iam-apiserver.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${IAM_LOG_DIR}/iam-apiserver.log</string>
    <key>StandardErrorPath</key>
    <string>${IAM_LOG_DIR}/iam-apiserver.error.log</string>
    <key>WorkingDirectory</key>
    <string>${IAM_INSTALL_DIR}</string>
</dict>
</plist>
EOF
}

# 安装
function iam::apiserver::install()
{
  pushd "${IAM_ROOT}"

  # 0. 确保目录存在
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}/cert"
  mkdir -p "${IAM_LOG_DIR}"
  mkdir -p "${IAM_DATA_DIR}"
  mkdir -p "$(iam::common::get_bin_path)"

  chmod +x ./hack/gencerts.sh
  chmod +x ./hack/genconfig.sh

  # 1. 生成 CA 证书和私钥
  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}"/cert
  cp "${LOCAL_OUTPUT_ROOT}/cert/ca"* "${IAM_CONFIG_DIR}/cert"

  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}"/cert iam-apiserver
  cp "${LOCAL_OUTPUT_ROOT}/cert/iam-apiserver"*.pem "${IAM_CONFIG_DIR}/cert"

  # 2. 构建 iam-apiserver
  make build BINS=iam-apiserver
  local bin_path=$(iam::common::get_bin_path)
  cp "${LOCAL_OUTPUT_ROOT}/bin/iam-apiserver" "${bin_path}/iam-apiserver"

  # 3. 生成并安装 iam-apiserver 的配置文件（iam-apiserver.yaml）
  ./hack/genconfig.sh "${ENV_FILE}" configs/iam-apiserver-template.yaml > "${IAM_CONFIG_DIR}/iam-apiserver.yaml"

  if iam::common::is_macos; then
    # 4. 创建并安装 iam-apiserver launchd plist 文件
    iam::apiserver::create_plist
    cp "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" "${HOME}/Library/LaunchAgents/"

    # 5. 启动 iam-apiserver 服务
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist" 2>/dev/null || true
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"

    # 等待服务启动
    sleep 3
  else
    # 4. 创建并安装 iam-apiserver systemd unit 文件（Linux）
    echo "${LINUX_PASSWORD}" | sudo -S bash -c \
      "./hack/genconfig.sh ${ENV_FILE} init/iam-apiserver.service > /etc/systemd/system/iam-apiserver.service"

    # 5. 启动 iam-apiserver 服务（Linux）
    iam::common::sudo "systemctl daemon-reload"
    iam::common::sudo "systemctl restart iam-apiserver"
    iam::common::sudo "systemctl enable iam-apiserver"
  fi

  iam::apiserver::status || return 1
  iam::apiserver::info

  iam::log::info "install iam-apiserver successfully"
  popd
}

# 卸载
function iam::apiserver::uninstall()
{
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)

  if iam::common::is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    iam::common::sudo "systemctl stop iam-apiserver"
    iam::common::sudo "systemctl disable iam-apiserver"
    iam::common::sudo "rm -f /etc/systemd/system/iam-apiserver.service"
  fi

  rm -f "${bin_path}/iam-apiserver"
  rm -f "${IAM_CONFIG_DIR}/iam-apiserver.yaml"
  rm -f "${IAM_CONFIG_DIR}/cert/iam-apiserver"*.pem
  set -o errexit
  iam::log::info "uninstall iam-apiserver successfully"
}

# 状态检查
function iam::apiserver::status()
{
  if iam::common::is_macos; then
    launchctl list | grep -q 'io.github.coding-hui.iam-apiserver' || {
      iam::log::error "iam-apiserver failed to start, maybe not installed properly"
      return 1
    }

    pgrep -x iam-apiserver &>/dev/null || {
      iam::log::error "iam-apiserver process not running"
      return 1
    }

    if echo | nc -z -w 2 ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>&1 | grep -q "refused\|failed"; then
      iam::log::error "cannot access insecure port, iam-apiserver maybe not startup"
      return 1
    fi
  else
    systemctl status iam-apiserver | grep -q 'active' || {
      iam::log::error "iam-apiserver failed to start, maybe not installed properly"
      return 1
    }

    if echo | telnet ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>&1 | grep refused &>/dev/null; then
      iam::log::error "cannot access insecure port, iam-apiserver maybe not startup"
      return 1
    fi
  fi

  iam::log::info "iam-apiserver status active"
}

if [[ "$*" =~ iam::apiserver:: ]];then
  eval $*
fi
