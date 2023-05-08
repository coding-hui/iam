#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# IAM 项目源码根目录
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${IAM_ROOT}/${OUT_DIR:-_output}"

# 设置统一的密码，方便记忆
readonly PASSWORD=${PASSWORD:-'iam59!z$'}

# Linux系统 going 用户
readonly LINUX_USERNAME=${LINUX_USERNAME:-going}
# Linux root & going 用户密码
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-${PASSWORD}}

# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${IAM_ROOT}/hack/install/environment.sh

# MariaDB 配置信息
readonly MARIADB_ADMIN_USERNAME=${MARIADB_ADMIN_USERNAME:-root}        # MariaDB root 用户
readonly MARIADB_ADMIN_PASSWORD=${MARIADB_ADMIN_PASSWORD:-${PASSWORD}} # MariaDB root 用户密码
readonly MARIADB_HOST=${MARIADB_HOST:-127.0.0.1:3306}                  # MariaDB 主机地址
readonly MARIADB_DATABASE=${MARIADB_DATABASE:-iam}                     # MariaDB iam 应用使用的数据库名
readonly MARIADB_USERNAME=${MARIADB_USERNAME:-iam}                     # iam 数据库用户名
readonly MARIADB_PASSWORD=${MARIADB_PASSWORD:-${PASSWORD}}             # iam 数据库密码

# Redis 配置信息
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1}           # Redis 主机地址
readonly REDIS_PORT=${REDIS_PORT:-6379}                # Redis 监听端口
readonly REDIS_USERNAME=${REDIS_USERNAME:-''}          # Redis 用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-${PASSWORD}} # Redis 密码

# MongoDB 配置
readonly MONGO_ADMIN_USERNAME=${MONGO_ADMIN_USERNAME:-root}        # MongoDB root 用户
readonly MONGO_ADMIN_PASSWORD=${MONGO_ADMIN_PASSWORD:-${PASSWORD}} # MongoDB root 用户密码
readonly MONGO_HOST=${MONGO_HOST:-127.0.0.1}                       # MongoDB 地址
readonly MONGO_PORT=${MONGO_PORT:-27017}                           # MongoDB 端口
readonly MONGO_USERNAME=${MONGO_USERNAME:-iam}                     # MongoDB 用户名
readonly MONGO_PASSWORD=${MONGO_PASSWORD:-${PASSWORD}}             # MongoDB 密码

# iam 配置
readonly IAM_DATA_DIR=${IAM_DATA_DIR:-/data/iam}           # iam 各组件数据目录
readonly IAM_INSTALL_DIR=${IAM_INSTALL_DIR:-/opt/iam}      # iam 安装文件存放目录
readonly IAM_CONFIG_DIR=${IAM_CONFIG_DIR:-/etc/iam}        # iam 配置文件存放目录
readonly IAM_LOG_DIR=${IAM_LOG_DIR:-/var/log/iam}          # iam 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${IAM_CONFIG_DIR}/cert/ca.pem} # CA

# iam-apiserver 配置
readonly IAM_APISERVER_HOST=${IAM_APISERVER_HOST:-127.0.0.1} # iam-apiserver 部署机器 IP 地址
readonly IAM_APISERVER_GRPC_BIND_ADDRESS=${IAM_APISERVER_GRPC_BIND_ADDRESS:-0.0.0.0}
readonly IAM_APISERVER_GRPC_BIND_PORT=${IAM_APISERVER_GRPC_BIND_PORT:-8081}
readonly IAM_APISERVER_INSECURE_BIND_ADDRESS=${IAM_APISERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly IAM_APISERVER_INSECURE_BIND_PORT=${IAM_APISERVER_INSECURE_BIND_PORT:-8080}
readonly IAM_APISERVER_SECURE_BIND_ADDRESS=${IAM_APISERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly IAM_APISERVER_SECURE_BIND_PORT=${IAM_APISERVER_SECURE_BIND_PORT:-8443}
readonly IAM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${IAM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${IAM_CONFIG_DIR}/cert/iam-apiserver.pem}
readonly IAM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${IAM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${IAM_CONFIG_DIR}/cert/iam-apiserver-key.pem}

# iam-authz-server 配置
readonly IAM_AUTHZ_SERVER_HOST=${IAM_AUTHZ_SERVER_HOST:-127.0.0.1} # iam-authz-server 部署机器 IP 地址
readonly IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=${IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly IAM_AUTHZ_SERVER_INSECURE_BIND_PORT=${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT:-9090}
readonly IAM_AUTHZ_SERVER_SECURE_BIND_ADDRESS=${IAM_AUTHZ_SERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly IAM_AUTHZ_SERVER_SECURE_BIND_PORT=${IAM_AUTHZ_SERVER_SECURE_BIND_PORT:-9443}
readonly IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${IAM_CONFIG_DIR}/cert/iam-authz-server.pem}
readonly IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${IAM_CONFIG_DIR}/cert/iam-authz-server-key.pem}
readonly IAM_AUTHZ_SERVER_CLIENT_CA_FILE=${IAM_AUTHZ_SERVER_CLIENT_CA_FILE:-${CA_FILE}}
readonly IAM_AUTHZ_SERVER_RPCSERVER=${IAM_AUTHZ_SERVER_RPCSERVER:-${IAM_APISERVER_HOST}:${IAM_APISERVER_GRPC_BIND_PORT}}

# iam-pump 配置
readonly IAM_PUMP_HOST=${IAM_PUMP_HOST:-127.0.0.1} # iam-pump 部署机器 IP 地址
readonly IAM_PUMP_COLLECTION_NAME=${IAM_PUMP_COLLECTION_NAME:-iam_analytics}
readonly IAM_PUMP_MONGO_URL=${IAM_PUMP_MONGO_URL:-mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${IAM_PUMP_COLLECTION_NAME}?authSource=${IAM_PUMP_COLLECTION_NAME}}

# iam-watcher配置
readonly IAM_WATCHER_HOST=${IAM_WATCHER_HOST:-127.0.0.1} # iam-watcher 部署机器 IP 地址

# iamctl 配置
readonly CONFIG_USER_USERNAME=${CONFIG_USER_USERNAME:-admin}
readonly CONFIG_USER_PASSWORD=${CONFIG_USER_PASSWORD:-Admin@2021}
readonly CONFIG_USER_CLIENT_CERTIFICATE=${CONFIG_USER_CLIENT_CERTIFICATE:-${HOME}/.iam/cert/admin.pem}
readonly CONFIG_USER_CLIENT_KEY=${CONFIG_USER_CLIENT_KEY:-${HOME}/.iam/cert/admin-key.pem}
readonly CONFIG_SERVER_ADDRESS=${CONFIG_SERVER_ADDRESS:-${IAM_APISERVER_HOST}:${IAM_APISERVER_SECURE_BIND_PORT}}
readonly CONFIG_SERVER_CERTIFICATE_AUTHORITY=${CONFIG_SERVER_CERTIFICATE_AUTHORITY:-${CA_FILE}}
