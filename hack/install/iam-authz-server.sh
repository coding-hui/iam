#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::authzserver::info() {
cat << EOF
iam-authz-server insecure listen on: ${IAM_AUTHZ_SERVER_HOST}:${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT}
iam-authz-server secure listen on: ${IAM_AUTHZ_SERVER_HOST}:${IAM_AUTHZ_SERVER_SECURE_BIND_PORT}
EOF
}

# 创建 launchd plist 文件（仅 macOS，用户级 Agent，无需 sudo）
function iam::authzserver::create_plist() {
  local bin_path=$(iam::common::get_bin_path)
  cat > "${IAM_ROOT}/io.github.coding-hui.iam-authz-server.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.github.coding-hui.iam-authz-server</string>
    <key>ProgramArguments</key>
    <array>
        <string>${bin_path}/iam-authz-server</string>
        <string>--config</string>
        <string>${IAM_CONFIG_DIR}/iam-authz-server.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${IAM_LOG_DIR}/iam-authz-server.log</string>
    <key>StandardErrorPath</key>
    <string>${IAM_LOG_DIR}/iam-authz-server.error.log</string>
    <key>WorkingDirectory</key>
    <string>${IAM_INSTALL_DIR}</string>
</dict>
</plist>
EOF
}

# 安装
function iam::authzserver::install()
{
  pushd ${IAM_ROOT}

  # 0. 确保目录存在
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}/cert"
  mkdir -p "${IAM_LOG_DIR}"
  mkdir -p "$(iam::common::get_bin_path)"

  # 1. 生成 CA 证书和私钥
  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}/cert"
  cp "${LOCAL_OUTPUT_ROOT}/cert/ca"* "${IAM_CONFIG_DIR}/cert"

  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}/cert" iam-authz-server
  cp "${LOCAL_OUTPUT_ROOT}/cert/iam-authz-server"*.pem "${IAM_CONFIG_DIR}/cert"

  # 2. 构建 iam-authz-server
  make build BINS=iam-authz-server
  local bin_path=$(iam::common::get_bin_path)
  cp "${LOCAL_OUTPUT_ROOT}/bin/iam-authz-server" "${bin_path}/iam-authz-server"

  # 3. 生成并安装 iam-authz-server 的配置文件（iam-authz-server.yaml）
  ./hack/genconfig.sh "${ENV_FILE}" configs/iam-authz-server.yaml > "${IAM_CONFIG_DIR}/iam-authz-server.yaml"

  if iam::common::is_macos; then
    # 4. 创建并安装 iam-authz-server launchd plist 文件（macOS 用户级 Agent）
    iam::authzserver::create_plist
    cp "${IAM_ROOT}/io.github.coding-hui.iam-authz-server.plist" "${HOME}/Library/LaunchAgents/"

    # 5. 启动 iam-authz-server 服务（macOS 用户级 Agent，无需 sudo）
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-authz-server.plist" 2>/dev/null || true
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-authz-server.plist"

    # 等待服务启动
    sleep 3
  else
    # 4. 创建并安装 iam-authz-server systemd unit 文件（Linux）
    echo "${LINUX_PASSWORD}" | sudo -S bash -c \
      "./hack/genconfig.sh ${ENV_FILE} init/iam-authz-server.service > /etc/systemd/system/iam-authz-server.service"

    # 5. 启动 iam-authz-server 服务（Linux）
    iam::common::sudo "systemctl daemon-reload"
    iam::common::sudo "systemctl restart iam-authz-server"
    iam::common::sudo "systemctl enable iam-authz-server"
  fi

  iam::authzserver::status || return 1
  iam::authzserver::info

  iam::log::info "install iam-authz-server successfully"
  popd
}

# 卸载
function iam::authzserver::uninstall()
{
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)

  if iam::common::is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-authz-server.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-authz-server.plist"
  else
    iam::common::sudo "systemctl stop iam-authz-server"
    iam::common::sudo "systemctl disable iam-authz-server"
    iam::common::sudo "rm -f /etc/systemd/system/iam-authz-server.service"
  fi

  rm -f "${bin_path}/iam-authz-server"
  rm -f "${IAM_CONFIG_DIR}/iam-authz-server.yaml"
  rm -f "${IAM_CONFIG_DIR}/cert/iam-authz-server"*.pem
  set -o errexit
  iam::log::info "uninstall iam-authz-server successfully"
}

# 状态检查
function iam::authzserver::status()
{
  if iam::common::is_macos; then
    launchctl list | grep -q 'io.github.coding-hui.iam-authz-server' || {
      iam::log::error "iam-authz-server failed to start, maybe not installed properly"
      return 1
    }

    pgrep -x iam-authz-server &>/dev/null || {
      iam::log::error "iam-authz-server process not running"
      return 1
    }

    if echo | nc -z -w 2 ${IAM_AUTHZ_SERVER_HOST} ${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT} 2>&1 | grep -q "refused\|failed"; then
      iam::log::error "cannot access insecure port, iam-authz-server maybe not startup"
      return 1
    fi
  else
    systemctl status iam-authz-server | grep -q 'active' || {
      iam::log::error "iam-authz-server failed to start, maybe not installed properly"
      return 1
    }

    if echo | telnet ${IAM_AUTHZ_SERVER_HOST} ${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT} 2>&1 | grep refused &>/dev/null; then
      iam::log::error "cannot access insecure port, iam-authz-server maybe not startup"
      return 1
    fi
  fi

  iam::log::info "iam-authz-server status active"
}

if [[ "$*" =~ iam::authzserver:: ]];then
  :
fi
