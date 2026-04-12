#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::pump::info() {
cat << EOF
iam-pump listen on: ${IAM_PUMP_HOST}
EOF
}

# 创建 launchd plist 文件（仅 macOS，用户级 Agent，无需 sudo）
function iam::pump::create_plist() {
  local bin_path=$(iam::common::get_bin_path)
  cat > "${IAM_ROOT}/io.github.coding-hui.iam-pump.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.github.coding-hui.iam-pump</string>
    <key>ProgramArguments</key>
    <array>
        <string>${bin_path}/iam-pump</string>
        <string>--config</string>
        <string>${IAM_CONFIG_DIR}/iam-pump.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${IAM_LOG_DIR}/iam-pump.log</string>
    <key>StandardErrorPath</key>
    <string>${IAM_LOG_DIR}/iam-pump.error.log</string>
    <key>WorkingDirectory</key>
    <string>${IAM_INSTALL_DIR}</string>
</dict>
</plist>
EOF
}

# 安装
function iam::pump::install()
{
  pushd ${IAM_ROOT}

  # 0. 确保目录存在
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}"
  mkdir -p "${IAM_LOG_DIR}"
  mkdir -p "$(iam::common::get_bin_path)"

  # 1. 构建 iam-pump
  make build BINS=iam-pump
  local bin_path=$(iam::common::get_bin_path)
  cp "${LOCAL_OUTPUT_ROOT}/bin/iam-pump" "${bin_path}/iam-pump"

  # 2. 生成并安装 iam-pump 的配置文件（iam-pump.yaml）
  ./hack/genconfig.sh "${ENV_FILE}" configs/iam-pump.yaml > "${IAM_CONFIG_DIR}/iam-pump.yaml"

  if iam::common::is_macos; then
    # 3. 创建并安装 iam-pump launchd plist 文件（macOS 用户级 Agent）
    iam::pump::create_plist
    cp "${IAM_ROOT}/io.github.coding-hui.iam-pump.plist" "${HOME}/Library/LaunchAgents/"

    # 4. 启动 iam-pump 服务（macOS 用户级 Agent，无需 sudo）
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-pump.plist" 2>/dev/null || true
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-pump.plist"

    # 等待服务启动
    sleep 3
  else
    # 3. 创建并安装 iam-pump systemd unit 文件（Linux）
    echo "${LINUX_PASSWORD}" | sudo -S bash -c \
      "./hack/genconfig.sh ${ENV_FILE} init/iam-pump.service > /etc/systemd/system/iam-pump.service"

    # 4. 启动 iam-pump 服务（Linux）
    iam::common::sudo "systemctl daemon-reload"
    iam::common::sudo "systemctl restart iam-pump"
    iam::common::sudo "systemctl enable iam-pump"
  fi

  iam::pump::status || return 1
  iam::pump::info

  iam::log::info "install iam-pump successfully"
  popd
}

# 卸载
function iam::pump::uninstall()
{
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)

  if iam::common::is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-pump.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-pump.plist"
  else
    iam::common::sudo "systemctl stop iam-pump"
    iam::common::sudo "systemctl disable iam-pump"
    iam::common::sudo "rm -f /etc/systemd/system/iam-pump.service"
  fi

  rm -f "${bin_path}/iam-pump"
  rm -f "${IAM_CONFIG_DIR}/iam-pump.yaml"
  set -o errexit
  iam::log::info "uninstall iam-pump successfully"
}

# 状态检查
function iam::pump::status()
{
  if iam::common::is_macos; then
    launchctl list | grep -q 'io.github.coding-hui.iam-pump' || {
      iam::log::error "iam-pump failed to start, maybe not installed properly"
      return 1
    }

    pgrep -x iam-pump &>/dev/null || {
      iam::log::error "iam-pump process not running"
      return 1
    }

    # 监听端口在配置文件中是 hardcode (7070)
    if echo | nc -z -w 2 127.0.0.1 7070 2>&1 | grep -q "refused\|failed"; then
      iam::log::error "cannot access health check port, iam-pump maybe not startup"
      return 1
    fi
  else
    systemctl status iam-pump | grep -q 'active' || {
      iam::log::error "iam-pump failed to start, maybe not installed properly"
      return 1
    }

    # 监听端口在配置文件中是 hardcode (7070)
    if echo | telnet 127.0.0.1 7070 2>&1 | grep refused &>/dev/null; then
      iam::log::error "cannot access health check port, iam-pump maybe not startup"
      return 1
    fi
  fi

  iam::log::info "iam-pump status active"
}

if [[ "$*" =~ iam::pump:: ]];then
  eval $*
fi
