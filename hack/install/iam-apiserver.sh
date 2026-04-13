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
IAM apiserver listen on:
  - http://127.0.0.1:${IAM_APISERVER_INSECURE_BIND_PORT}
  - https://127.0.0.1:${IAM_APISERVER_SECURE_BIND_PORT}
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
  pushd "${IAM_ROOT}" > /dev/null

  iam::log::section "Installing IAM apiserver"

  # 0. 确保目录存在
  mkdir -p "${IAM_INSTALL_DIR}"
  mkdir -p "${IAM_CONFIG_DIR}/cert"
  mkdir -p "${IAM_LOG_DIR}"
  mkdir -p "${IAM_DATA_DIR}"
  mkdir -p "$(iam::common::get_bin_path)"

  chmod +x ./hack/gencerts.sh
  chmod +x ./hack/genconfig.sh

  # 1. 生成 CA 证书和私钥
  iam::log::substep "Generating certificates..."
  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}"/cert
  cp "${LOCAL_OUTPUT_ROOT}/cert/ca"* "${IAM_CONFIG_DIR}/cert" 2>/dev/null

  ./hack/gencerts.sh generate-iam-cert "${LOCAL_OUTPUT_ROOT}"/cert iam-apiserver
  cp "${LOCAL_OUTPUT_ROOT}/cert/iam-apiserver"*.pem "${IAM_CONFIG_DIR}/cert" 2>/dev/null

  # 2. 构建 iam-apiserver
  iam::log::substep "Building iam-apiserver..."
  CGO_ENABLED=1 go build -o "${LOCAL_OUTPUT_ROOT}/bin/iam-apiserver" github.com/coding-hui/iam/cmd/iam-apiserver
  local bin_path=$(iam::common::get_bin_path)
  cp "${LOCAL_OUTPUT_ROOT}/bin/iam-apiserver" "${bin_path}/iam-apiserver" 2>/dev/null

  # 3. 生成并安装 iam-apiserver 的配置文件（iam-apiserver.yaml）
  ./hack/genconfig.sh "${ENV_FILE}" configs/iam-apiserver-template.yaml > "${IAM_CONFIG_DIR}/iam-apiserver.yaml" 2>/dev/null

  # 4. 启动服务
  iam::log::substep "Starting services..."

  if iam::common::is_macos; then
    # 创建并安装 launchd plist
    iam::apiserver::create_plist
    mkdir -p "${HOME}/Library/LaunchAgents" 2>/dev/null
    cp "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" "${HOME}/Library/LaunchAgents/" 2>/dev/null
    iam::common::launchd_load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    # 创建并安装 iam-apiserver systemd unit 文件（Linux）
    echo "${LINUX_PASSWORD}" | sudo -S bash -c \
      "./hack/genconfig.sh ${ENV_FILE} init/iam-apiserver.service > /etc/systemd/system/iam-apiserver.service" 2>/dev/null

    iam::common::sudo "systemctl daemon-reload"
    iam::common::sudo "systemctl restart iam-apiserver"
    iam::common::sudo "systemctl enable iam-apiserver"
  fi

  # 等待服务启动（每1s检查一次，最多60s）
  iam::apiserver::wait_for_start || return 1

  iam::apiserver::status || return 1

  popd >/dev/null 2>&1
}

# 卸载
function iam::apiserver::uninstall()
{
  iam::log::section "Uninstalling IAM apiserver"
  set +o errexit
  local bin_path=$(iam::common::get_bin_path)

  if iam::common::is_macos; then
    iam::log::substep "Stopping launchd service..."
    iam::common::launchd_unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    iam::log::substep "Stopping systemd service..."
    iam::common::sudo "systemctl stop iam-apiserver"
    iam::common::sudo "systemctl disable iam-apiserver"
    iam::common::sudo "rm -f /etc/systemd/system/iam-apiserver.service"
  fi

  iam::log::substep "Removing binary..."
  rm -f "${bin_path}/iam-apiserver"
  rm -f "${IAM_CONFIG_DIR}/iam-apiserver.yaml"
  rm -f "${IAM_CONFIG_DIR}/cert/iam-apiserver"*.pem
  set -o errexit
  iam::log::info "Uninstall IAM apiserver successfully"
}

# 等待服务启动（每1s检查一次，最多60s）
function iam::apiserver::wait_for_start() {
  local max_attempts=60
  local attempt=0
  echo -n "    Waiting for IAM apiserver"
  while true; do
    if nc -z -w 1 ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>/dev/null; then
      echo " ✔"
      return 0
    fi
    sleep 1
    echo -n "."
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
      echo ""
      iam::log::error "IAM apiserver failed to start within 60 seconds"
      return 1
    fi
  done
}

# 状态检查
function iam::apiserver::status()
{
  if iam::common::is_macos; then
    nc -z -w 1 ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>/dev/null || {
      iam::log::error "IAM apiserver is not running"
      return 1
    }
    iam::log::substep "✔ IAM apiserver is running"
  else
    systemctl status iam-apiserver | grep -q 'active' || {
      iam::log::error "IAM apiserver is not running"
      return 1
    }
    iam::log::substep "✔ IAM apiserver is running"
  fi
}

# 启动服务
function iam::apiserver::start()
{
  if iam::common::is_macos; then
    local plist="${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
    if [ ! -f "${plist}" ]; then
      iam::log::error "IAM apiserver plist not found, please run 'make install' first"
      return 1
    fi
    iam::common::launchd_load "${plist}"
    iam::log::info "✔ IAM apiserver started"
  else
    iam::common::sudo "systemctl start iam-apiserver"
    iam::log::info "✔ IAM apiserver started"
  fi
}

# 停止服务
function iam::apiserver::stop()
{
  if iam::common::is_macos; then
    iam::common::launchd_unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
    iam::log::info "✔ IAM apiserver stopped"
  else
    iam::common::sudo "systemctl stop iam-apiserver"
    iam::log::info "✔ IAM apiserver stopped"
  fi
}

# 查看日志
function iam::apiserver::logs()
{
  if [ -f "${IAM_LOG_DIR}/iam-apiserver.log" ]; then
    tail -50 "${IAM_LOG_DIR}/iam-apiserver.log"
  else
    iam::log::error "log file not found: ${IAM_LOG_DIR}/iam-apiserver.log"
    return 1
  fi
}

if [[ "$*" =~ iam::apiserver:: ]];then
  :
fi
