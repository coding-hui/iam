#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/hack/install/common.sh

# 安装后打印必要的信息
function iam::redis::info() {
cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# ---- macOS ----------------------------------------------------------------

function iam::redis::_install_macos() {
  if ! command -v brew &>/dev/null; then
    iam::log::error_exit "Homebrew is not installed. Please install it first: https://brew.sh"
    return 1
  fi

  iam::log::section "Installing Redis"

  # 找到配置文件（Apple Silicon 与 Intel 路径不同）
  local redis_conf=""
  if [[ -f "/opt/homebrew/etc/redis.conf" ]]; then
    redis_conf="/opt/homebrew/etc/redis.conf"
  elif [[ -f "/usr/local/etc/redis.conf" ]]; then
    redis_conf="/usr/local/etc/redis.conf"
  else
    iam::log::error_exit "Redis configuration file not found"
    return 1
  fi

  # 检查 Redis 是否已安装，未安装则安装
  if ! brew list redis &>/dev/null 2>&1; then
    iam::log::substep "Installing redis..."
    HOMEBREW_NO_INSTALL_CLEANUP=1 HOMEBREW_NO_ENV_HINTS=1 brew install redis
  fi

  # 允许后台运行
  sed -i '' 's/^daemonize no/daemonize yes/' ${redis_conf}
  # 允许外网连接
  sed -i '' 's/^bind 127.0.0.1/# bind 127.0.0.1/' ${redis_conf}
  # 设置密码
  sed -i '' 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' ${redis_conf}
  sed -i '' 's/^requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' ${redis_conf}
  # 关闭保护模式
  sed -i '' 's/^protected-mode yes/protected-mode no/' ${redis_conf}

  iam::log::substep "Starting redis service..."
  iam::common::brew_service_stop redis
  iam::common::brew_service_start redis
  sleep 2
}

function iam::redis::_uninstall_macos() {
  iam::common::brew_service_stop redis
  iam::common::pkill redis-server
  brew uninstall redis 2>/dev/null || true
}

function iam::redis::_status_macos() {
  if ! pgrep -f redis-server &>/dev/null; then
    iam::log::error_exit "Redis not running"
    return 1
  fi

  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" PING 2>/dev/null \
    | grep -q "PONG" || {
    iam::log::error "cannot connect to Redis"
    return 1
  }
}

# ---- Ubuntu/Debian --------------------------------------------------------

function iam::redis::_install_ubuntu() {
  iam::common::sudo "apt-get -y install redis-server"

  local redis_conf="/etc/redis/redis.conf"
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' ${redis_conf}

  iam::common::sudo "ufw disable"
  iam::common::sudo "ufw status"

  iam::common::sudo "redis-server ${redis_conf}"
}

function iam::redis::_uninstall_ubuntu() {
  iam::common::sudo "/etc/init.d/redis-server stop"
  iam::common::sudo "apt-get -y remove redis-server"
  iam::common::sudo "rm -rf /var/lib/redis"
}

function iam::redis::_status_ubuntu() {
  if [[ -z "$(pgrep redis-server)" ]]; then
    iam::log::error_exit "Redis not running, maybe not installed properly"
    return 1
  fi

  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
    iam::log::error "cannot connect to Redis, maybe not initialized properly"
    return 1
  }
}

# ---- RHEL/CentOS ----------------------------------------------------------

function iam::redis::_install_linux() {
  iam::common::sudo "yum -y install redis"

  local redis_conf="/etc/redis.conf"
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' ${redis_conf}
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' ${redis_conf}

  iam::common::sudo "systemctl stop firewalld.service"
  iam::common::sudo "systemctl disable firewalld.service"

  iam::common::sudo "redis-server ${redis_conf}"
}

function iam::redis::_uninstall_linux() {
  iam::common::sudo "killall redis-server"
  iam::common::sudo "yum -y remove redis"
  iam::common::sudo "rm -rf /var/lib/redis"
}

function iam::redis::_status_linux() {
  if [[ -z "$(pgrep redis-server)" ]]; then
    iam::log::error_exit "Redis not running, maybe not installed properly"
    return 1
  fi

  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
    iam::log::error "cannot connect to Redis, maybe not initialized properly"
    return 1
  }
}

# ---- 公共入口 ---------------------------------------------------------------

# 安装
function iam::redis::install()
{
  if iam::common::is_macos; then
    iam::redis::_install_macos
  elif iam::common::is_ubuntu; then
    iam::redis::_install_ubuntu
  else
    iam::redis::_install_linux
  fi

  iam::redis::status || return 1
}

# 卸载
function iam::redis::uninstall()
{
  iam::log::section "Uninstalling Redis"
  set +o errexit
  if iam::common::is_macos; then
    iam::log::substep "Stopping redis service..."
    iam::redis::_uninstall_macos
  elif iam::common::is_ubuntu; then
    iam::log::substep "Stopping redis service..."
    iam::redis::_uninstall_ubuntu
  else
    iam::log::substep "Stopping redis service..."
    iam::redis::_uninstall_linux
  fi
  set -o errexit
  iam::log::info "Uninstall Redis successfully"
}

# 状态检查
function iam::redis::status()
{
  if iam::common::is_macos; then
    iam::redis::_status_macos
  elif iam::common::is_ubuntu; then
    iam::redis::_status_ubuntu
  else
    iam::redis::_status_linux
  fi
}

if [[ "$*" =~ iam::redis:: ]];then
: # no-op
fi
