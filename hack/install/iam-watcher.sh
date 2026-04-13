#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::watcher::info() {
cat << EOF
iam-watcher listen on: ${IAM_WATCHER_HOST}
EOF
}

# 创建 launchd plist 文件（仅 macOS，用户级 Agent，无需 sudo）
function iam::watcher::create_plist() {
  local bin_path=$(iam::common::get_bin_path)
  cat > "${IAM_ROOT}/io.github.coding-hui.iam-watcher.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.github.coding-hui.iam-watcher</string>
    <key>ProgramArguments</key>
    <array>
        <string>${bin_path}/iam-watcher</string>
        <string>--config</string>
        <string>${IAM_CONFIG_DIR}/iam-watcher.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${IAM_LOG_DIR}/iam-watcher.log</string>
    <key>StandardErrorPath</key>
    <string>${IAM_LOG_DIR}/iam-watcher.error.log</string>
    <key>WorkingDirectory</key>
    <string>${IAM_INSTALL_DIR}</string>
</dict>
</plist>
EOF
}

# 安装
function iam::watcher::install()
{
  pushd ${IAM_ROOT} > /dev/null

  # 0. 确保目录存在
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}"
  mkdir -p "${IAM_LOG_DIR}"
  mkdir -p "$(iam::common::get_bin_path)"

  # 1. 构建 iam-watcher
  make build BINS=iam-watcher
  local bin_path=$(iam::common::get_bin_path)
  cp "${LOCAL_OUTPUT_ROOT}/bin/iam-watcher" "${bin_path}/iam-watcher"

  # 2. 生成并安装 iam-watcher 的配置文件（iam-watcher.yaml）
  ./hack/genconfig.sh "${ENV_FILE}" configs/iam-watcher.yaml > "${IAM_CONFIG_DIR}/iam-watcher.yaml"

  if iam::common::is_macos; then
    # 3. 创建并安装 iam-watcher launchd plist 文件（macOS 用户级 Agent）
    iam::watcher::create_plist
    cp "${IAM_ROOT}/io.github.coding-hui.iam-watcher.plist" "${HOME}/Library/LaunchAgents/"

    # 4. 启动 iam-watcher 服务（macOS 用户级 Agent，无需 sudo）
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-watcher.plist" 2>/dev/null || true
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-watcher.plist"

    # 等待服务启动
    sleep 3
  else
    # 3. 创建并安装 iam-watcher systemd unit 文件（Linux）
    echo "${LINUX_PASSWORD}" | sudo -S bash -c \
      "./hack/genconfig.sh ${ENV_FILE} init/iam-watcher.service > /etc/systemd/system/iam-watcher.service"

    # 4. 启动 iam-watcher 服务（Linux）
    iam::common::sudo "systemctl daemon-reload"
    iam::common::sudo "systemctl restart iam-watcher"
    iam::common::sudo "systemctl enable iam-watcher"
  fi

  iam::watcher::status || return 1
  iam::watcher::info

  iam::log::info "install iam-watcher successfully"
  popd
}

# 卸载
function iam::watcher::uninstall()
{
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)

  if iam::common::is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-watcher.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-watcher.plist"
  else
    iam::common::sudo "systemctl stop iam-watcher"
    iam::common::sudo "systemctl disable iam-watcher"
    iam::common::sudo "rm -f /etc/systemd/system/iam-watcher.service"
  fi

  rm -f "${bin_path}/iam-watcher"
  rm -f "${IAM_CONFIG_DIR}/iam-watcher.yaml"
  set -o errexit
  iam::log::info "uninstall iam-watcher successfully"
}

# 状态检查
function iam::watcher::status()
{
  if iam::common::is_macos; then
    launchctl list | grep -q 'io.github.coding-hui.iam-watcher' || {
      iam::log::error "iam-watcher failed to start, maybe not installed properly"
      return 1
    }

    pgrep -x iam-watcher &>/dev/null || {
      iam::log::error "iam-watcher process not running"
      return 1
    }

    # 监听端口在配置文件中是 hardcode (5050)
    if echo | nc -z -w 2 127.0.0.1 5050 2>&1 | grep -q "refused\|failed"; then
      iam::log::error "cannot access health check port, iam-watcher maybe not startup"
      return 1
    fi
  else
    systemctl status iam-watcher | grep -q 'active' || {
      iam::log::error "iam-watcher failed to start, maybe not installed properly"
      return 1
    }

    # 监听端口在配置文件中是 hardcode (5050)
    if echo | telnet 127.0.0.1 5050 2>&1 | grep refused &>/dev/null; then
      iam::log::error "cannot access health check port, iam-watcher maybe not startup"
      return 1
    fi
  fi

  iam::log::info "iam-watcher status active"
}

if [[ "$*" =~ iam::watcher:: ]];then
  :
fi
