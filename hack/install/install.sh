#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# The root of the build/dist directory
IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)

# Enable CGO for macOS to support sqlite driver (requires CGO)
if [[ "$(uname -s)" == "Darwin" ]]; then
  export CGO_ENABLED=1
fi

source "${IAM_ROOT}/hack/install/common.sh"

source ${IAM_ROOT}/hack/install/mariadb.sh
source ${IAM_ROOT}/hack/install/redis.sh
source ${IAM_ROOT}/hack/install/mongodb.sh
source ${IAM_ROOT}/hack/install/iam-apiserver.sh
source ${IAM_ROOT}/hack/install/iam-authz-server.sh
source ${IAM_ROOT}/hack/install/iam-pump.sh
source ${IAM_ROOT}/hack/install/iam-watcher.sh
source ${IAM_ROOT}/hack/install/iamctl.sh
source ${IAM_ROOT}/hack/install/man.sh
source ${IAM_ROOT}/hack/install/test.sh

# 如果是通过脚本安装，需要先尝试获取安装脚本指定的 Tag，Tag 记录在 version 文件中
function iam::install::obtain_branch_flag() {
  if [ -f "${IAM_ROOT}"/version ]; then
    echo $(cat "${IAM_ROOT}"/version)
  fi
}

function iam::install::prepare_iam() {
  rm -rf $WORKSPACE/golang/src/github.com/coding-hui/iam # clean up

  # 1. 下载 iam 项目代码，先强制删除 iam 目录，确保 iam 源码都是最新的指定版本
  mkdir -p $WORKSPACE/golang/src/github.com/coding-hui && cd $WORKSPACE/golang/src/github.com/coding-hui
  git clone -b $(iam::install::obtain_branch_flag) --depth=1 https://github.com/coding-hui/iam
  go work use ./iam

  # NOTICE: 因为切换编译路径，所以这里要重新赋值 IAM_ROOT 和 LOCAL_OUTPUT_ROOT
  IAM_ROOT=$WORKSPACE/golang/src/github.com/coding-hui/iam
  LOCAL_OUTPUT_ROOT="${IAM_ROOT}/${OUT_DIR:-_output}"

  pushd ${IAM_ROOT} > /dev/null

  # 2. 配置 $HOME/.bashrc 添加一些便捷入口
  if ! grep -q 'Alias for quick access' $HOME/.bashrc; then
    cat <<'EOF' >>$HOME/.bashrc
# Alias for quick access
export GOSRC="$WORKSPACE/golang/src"
export IAM_ROOT="$GOSRC/github.com/coding-hui/iam"
alias ch="cd $GOSRC/github.com/coding-hui"
alias i="cd $GOSRC/github.com/coding-hui/iam"
EOF
  fi

  # 3. 初始化 MariaDB 数据库，创建 iam 数据库

  # 3.1 登录数据库并创建 iam 用户
  mysql -h127.0.0.1 -P3306 -u"${MARIADB_ADMIN_USERNAME}" -p"${MARIADB_ADMIN_PASSWORD}" <<EOF
grant all on iam.* TO ${MARIADB_USERNAME}@127.0.0.1 identified by "${MARIADB_PASSWORD}";
flush privileges;
EOF

  # 3.2 用 iam 用户登录 mysql，执行 iam.sql 文件，创建 iam 数据库
  mysql -h127.0.0.1 -P3306 -u${MARIADB_USERNAME} -p"${MARIADB_PASSWORD}" <<EOF
source configs/iam.sql;
show databases;
EOF

  # 4. 创建必要的目录
  echo ${LINUX_PASSWORD} | sudo -S mkdir -p ${IAM_DATA_DIR}/{iam-apiserver,iam-authz-server,iam-pump,iam-watcher}
  iam::common::sudo "mkdir -p ${IAM_INSTALL_DIR}/bin"
  iam::common::sudo "mkdir -p ${IAM_CONFIG_DIR}/cert"
  iam::common::sudo "mkdir -p ${IAM_LOG_DIR}"

  # 5. 安装 cfssl 工具集
  ! command -v cfssl &>/dev/null || ! command -v cfssl-certinfo &>/dev/null || ! command -v cfssljson &>/dev/null && {
    iam::install::install_cfssl || return 1
  }

  # 6. 配置 hosts
  if ! egrep -q 'iam.*coding-hui.com' /etc/hosts; then
    echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' >> /etc/hosts
    127.0.0.1 iam.api.coding-hui.com
    127.0.0.1 iam.authz.coding-hui.com
    EOF"
  fi

  iam::log::info "prepare for iam installation successfully"
  popd
}

function iam::install::unprepare_iam() {
  pushd ${IAM_ROOT} > /dev/null

  # 1. 删除 iam 数据库和用户 (仅 Linux，macOS 使用 SQLite)
  if ! iam::common::is_macos && command -v mysql &>/dev/null; then
    mysql -h127.0.0.1 -P3306 -u"${MARIADB_ADMIN_USERNAME}" -p"${MARIADB_ADMIN_PASSWORD}" <<EOF
drop database iam;
drop user ${MARIADB_USERNAME}@127.0.0.1
EOF
  fi

  # 2. 删除创建的目录
  iam::common::sudo "rm -rf ${IAM_DATA_DIR}"
  iam::common::sudo "rm -rf ${IAM_INSTALL_DIR}"
  iam::common::sudo "rm -rf ${IAM_CONFIG_DIR}"
  iam::common::sudo "rm -rf ${IAM_LOG_DIR}"

  # 3. 删除配置 hosts
  local sed_i_flag
  sed_i_flag=$(iam::common::get_sed_i_flag)
  echo ${LINUX_PASSWORD} | sudo -S sed ${sed_i_flag} '/iam.api.coding-hui.com/d' /etc/hosts
  echo ${LINUX_PASSWORD} | sudo -S sed ${sed_i_flag} '/iam.authz.coding-hui.com/d' /etc/hosts

  iam::log::info "unprepare for iam installation successfully"
  popd
}

function iam::install::install_cfssl() {
  mkdir -p $HOME/bin/
  curl -fsSL https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssl_1.6.1_linux_amd64 -O $HOME/bin/cfssl
  curl -fsSL https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssljson_1.6.1_linux_amd64 -O $HOME/bin/cfssljson
  curl -fsSL https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssl-certinfo_1.6.1_linux_amd64 -O $HOME/bin/cfssl-certinfo
  chmod +x $HOME/bin/{cfssl,cfssljson,cfssl-certinfo}
  iam::log::info "install cfssl tools successfully"
}

function iam::install::install_storage() {
  iam::redis::install || return 1
}

function iam::install::uninstall_storage() {
  iam::redis::uninstall || return 1
}

# 安装 IAM 应用
function iam::install::install_iam() {
  # 1. 安装存储服务 (Redis)
  iam::install::install_storage || return 1

  # 2. 安装 iam-apiserver 服务
  iam::apiserver::install || return 1

  # 3. 安装 iam-authz-server 服务
  iam::authzserver::install || return 1

  # 4. 安装 iamctl 客户端工具
  iam::iamctl::install || return 1
}

function iam::install::uninstall_iam() {
  iam::iamctl::uninstall || return 1
  iam::apiserver::uninstall || return 1
  iam::install::uninstall_storage || return 1
}

function iam::install::uninstall() {
  iam::install::uninstall_iam || return 1
  iam::log::info "uninstall iam application successfully"
}

# ==============================================================================
# Main dispatch - supports both direct sourcing and CLI execution
# Usage:
#   ./install.sh install        # Install IAM
#   ./install.sh uninstall     # Uninstall IAM
#   ./install.sh start         # Start services
#   ./install.sh stop          # Stop services
#   ./install.sh status        # Check status
#   ./install.sh restart       # Restart services
#   ./install.sh logs          # Show logs

if [[ -n "${1}" ]]; then
  case "${1}" in
    install)
      iam::install::install_iam
      ;;
    uninstall)
      iam::install::uninstall
      ;;
    start)
      iam::apiserver::start
      ;;
    stop)
      iam::apiserver::stop
      ;;
    status)
      iam::apiserver::status
      ;;
    restart)
      iam::apiserver::stop
      iam::apiserver::start
      ;;
    logs)
      iam::apiserver::logs
      ;;
    *)
      echo "Usage: $0 {install|uninstall|start|stop|status|restart|logs}"
      exit 1
      ;;
  esac
fi
