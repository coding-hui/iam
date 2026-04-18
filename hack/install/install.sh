#!/usr/bin/env bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Simplified IAM installation script for local development.
# Supports: install, start, stop, status, restart, logs, uninstall

set -o errexit
set -o nounset
set -o pipefail

IAM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
OUTPUT_DIR="${IAM_ROOT}/_output"

# Detect platform
is_macos() { [[ "$(uname -s)" == "Darwin" ]]; }
is_linux() { [[ "$(uname -s)" == "Linux" ]]; }

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Paths
if is_macos; then
  IAM_INSTALL_DIR="${HOME}/.iam"
else
  IAM_INSTALL_DIR="/opt/iam"
fi
IAM_BIN_DIR="${IAM_INSTALL_DIR}/bin"
IAM_CONFIG_DIR="${IAM_INSTALL_DIR}/conf"
IAM_DATA_DIR="${HOME}/.iam/data"
IAM_LOG_DIR="${IAM_INSTALL_DIR}/logs"
IAM_CERT_DIR="${IAM_CONFIG_DIR}/cert"

# ==============================================================================
# CA Certificate Generation (simplified, self-contained)

generate_ca() {
  mkdir -p "${IAM_CERT_DIR}"

  # Generate CA private key
  openssl genrsa -out "${IAM_CERT_DIR}/ca-key.pem" 2048 2>/dev/null || {
    log_error "Failed to generate CA key"
    return 1
  }

  # Generate CA certificate
  openssl req -x509 -new -nodes -key "${IAM_CERT_DIR}/ca-key.pem" \
    -sha256 -days 3650 -out "${IAM_CERT_DIR}/ca.pem" \
    -subj "/CN=IAM CA" 2>/dev/null || {
    log_error "Failed to generate CA certificate"
    return 1
  }

  # Generate server certificate
  openssl genrsa -out "${IAM_CERT_DIR}/iam-apiserver-key.pem" 2048 2>/dev/null || {
    log_error "Failed to generate server key"
    return 1
  }

  cat > "${IAM_CERT_DIR}/iam-apiserver.cnf" << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
[req_distinguished_name]
[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
IP.1 = 127.0.0.1
IP.2 = ::1
DNS.1 = localhost
EOF

  openssl req -new -key "${IAM_CERT_DIR}/iam-apiserver-key.pem" \
    -out "${IAM_CERT_DIR}/iam-apiserver.csr" \
    -subj "/CN=IAM APIServer" \
    -config "${IAM_CERT_DIR}/iam-apiserver.cnf" 2>/dev/null || {
    log_error "Failed to generate CSR"
    return 1
  }

  openssl x509 -req -in "${IAM_CERT_DIR}/iam-apiserver.csr" \
    -CA "${IAM_CERT_DIR}/ca.pem" \
    -CAkey "${IAM_CERT_DIR}/ca-key.pem" \
    -CAcreateserial \
    -out "${IAM_CERT_DIR}/iam-apiserver.pem" \
    -days 825 -sha256 \
    -extensions v3_req \
    -extfile "${IAM_CERT_DIR}/iam-apiserver.cnf" 2>/dev/null || {
    log_error "Failed to generate server certificate"
    return 1
  }

  # Cleanup
  rm -f "${IAM_CERT_DIR}/iam-apiserver.csr" "${IAM_CERT_DIR}/iam-apiserver.cnf" "${IAM_CERT_DIR}/ca-key.pem"

  log_info "CA certificates generated"
}

# ==============================================================================
# Build iam-apiserver

build_apiserver() {
  log_info "Building iam-apiserver..."

  local version=$(git describe --tags --always --match='v*' 2>/dev/null || echo "v0.0.0")
  local git_commit=$(git rev-parse HEAD 2>/dev/null || echo "")
  local git_tree_state="dirty"
  if [ -z "$(git status --porcelain 2>/dev/null)" ]; then
    git_tree_state="clean"
  fi
  local build_date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

  mkdir -p "${OUTPUT_DIR}/bin"

  CGO_ENABLED=1 go build -ldflags "\
    -X github.com/coding-hui/common/version.GitVersion=${version} \
    -X github.com/coding-hui/common/version.GitCommit=${git_commit} \
    -X github.com/coding-hui/common/version.GitTreeState=${git_tree_state} \
    -X github.com/coding-hui/common/version.BuildDate=${build_date}" \
    -o "${OUTPUT_DIR}/bin/iam-apiserver" "${IAM_ROOT}/cmd/apiserver" || {
    log_error "Failed to build iam-apiserver"
    return 1
  }

  log_info "iam-apiserver built successfully"
}

# ==============================================================================
# Generate config file

generate_config() {
  log_info "Generating configuration..."

  # Use template if exists, otherwise create minimal config
  local template="${IAM_ROOT}/configs/iam-apiserver-template.yaml"
  local config="${IAM_CONFIG_DIR}/iam-apiserver.yaml"

  mkdir -p "${IAM_CONFIG_DIR}"

  if [ -f "${template}" ]; then
    # Simple substitution for common variables
    sed -e "s|\${IAM_APISERVER_HOST}|127.0.0.1|g" \
        -e "s|\${IAM_APISERVER_INSECURE_BIND_PORT}|8080|g" \
        -e "s|\${IAM_APISERVER_SECURE_BIND_PORT}|8443|g" \
        -e "s|\${IAM_APISERVER_GRPC_BIND_PORT}|8081|g" \
        -e "s|\${IAM_CONFIG_DIR}|${IAM_CONFIG_DIR}|g" \
        -e "s|\${IAM_CERT_DIR}|${IAM_CERT_DIR}|g" \
        -e "s|\${IAM_DATA_DIR}|${IAM_DATA_DIR}|g" \
        -e "s|\${IAM_LOG_DIR}|${IAM_LOG_DIR}|g" \
        -e "s|\${REDIS_HOST}|127.0.0.1|g" \
        -e "s|\${REDIS_PORT}|6379|g" \
        "${template}" > "${config}" 2>/dev/null || true
  else
    # Create minimal SQLite config
    cat > "${config}" << EOF
server:
  insecure:
    bind_address: 127.0.0.1
    bind_port: 8080
  secure:
    bind_address: 0.0.0.0
    bind_port: 8443
    tls:
      cert_file: ${IAM_CERT_DIR}/iam-apiserver.pem
      key_file: ${IAM_CERT_DIR}/iam-apiserver-key.pem
  grpc:
    bind_address: 0.0.0.0
    bind_port: 8081

database:
  type: sqlite
  dsn: ${IAM_DATA_DIR}/iam.db

cache:
  type: memory
  redis:
    host: 127.0.0.1
    port: 6379

log:
  level: info
  path: ${IAM_LOG_DIR}
EOF
  fi

  log_info "Configuration generated"
}

# ==============================================================================
# Service management

install_service() {
  log_info "Installing IAM apiserver..."

  mkdir -p "${IAM_INSTALL_DIR}" "${IAM_BIN_DIR}" "${IAM_CONFIG_DIR}" "${IAM_DATA_DIR}" "${IAM_LOG_DIR}" "${IAM_CERT_DIR}"

  # Build
  build_apiserver || return 1

  # Generate CA and certs
  generate_ca || return 1

  # Copy binary
  cp "${OUTPUT_DIR}/bin/iam-apiserver" "${IAM_BIN_DIR}/iam-apiserver"
  chmod +x "${IAM_BIN_DIR}/iam-apiserver"

  # Generate config
  generate_config || return 1

  if is_macos; then
    # Create launchd plist
    cat > "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.github.coding-hui.iam-apiserver</string>
    <key>ProgramArguments</key>
    <array>
        <string>${IAM_BIN_DIR}/iam-apiserver</string>
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
    mkdir -p "${HOME}/Library/LaunchAgents"
    cp "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" "${HOME}/Library/LaunchAgents/"
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    # Create systemd unit
    cat > /tmp/iam-apiserver.service << EOF
[Unit]
Description=IAM API Server
After=network.target

[Service]
Type=simple
ExecStart=${IAM_BIN_DIR}/iam-apiserver --config ${IAM_CONFIG_DIR}/iam-apiserver.yaml
WorkingDirectory=${IAM_INSTALL_DIR}
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
    sudo cp /tmp/iam-apiserver.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable iam-apiserver
    sudo systemctl start iam-apiserver
  fi

  # Wait for startup
  log_info "Waiting for service to start..."
  local max_attempts=30
  local attempt=0
  while ! nc -z -w1 127.0.0.1 8080 2>/dev/null; do
    sleep 1
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
      log_error "Service failed to start within ${max_attempts} seconds"
      return 1
    fi
  done

  log_info "IAM apiserver installed successfully"
  log_info "  - HTTP: http://127.0.0.1:8080"
  log_info "  - HTTPS: https://127.0.0.1:8443"
  log_info "  - Config: ${IAM_CONFIG_DIR}/iam-apiserver.yaml"
}

uninstall_service() {
  log_info "Uninstalling IAM apiserver..."

  if is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl stop iam-apiserver 2>/dev/null || true
    sudo systemctl disable iam-apiserver 2>/dev/null || true
    sudo rm -f /etc/systemd/system/iam-apiserver.service
  fi

  rm -rf "${IAM_INSTALL_DIR}"
  log_info "IAM apiserver uninstalled"
}

start_service() {
  if is_macos; then
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl start iam-apiserver
  fi
  log_info "IAM apiserver started"
}

stop_service() {
  if is_macos; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl stop iam-apiserver
  fi
  log_info "IAM apiserver stopped"
}

restart_service() {
  stop_service
  sleep 1
  start_service
  log_info "IAM apiserver restarted"
}

check_status() {
  if nc -z -w1 127.0.0.1 8080 2>/dev/null; then
    log_info "IAM apiserver is running (http://127.0.0.1:8080)"
    return 0
  else
    log_warn "IAM apiserver is not running"
    return 1
  fi
}

show_logs() {
  local log_file="${IAM_LOG_DIR}/iam-apiserver.log"
  if [ -f "${log_file}" ]; then
    tail -50 "${log_file}"
  else
    log_warn "Log file not found: ${log_file}"
  fi
}

# ==============================================================================
# Main dispatch

case "${1:-}" in
  install)
    install_service
    ;;
  uninstall)
    uninstall_service
    ;;
  start)
    start_service
    ;;
  stop)
    stop_service
    ;;
  restart)
    restart_service
    ;;
  status)
    check_status
    ;;
  logs)
    show_logs
    ;;
  *)
    echo "Usage: $0 {install|uninstall|start|stop|restart|status|logs}"
    exit 1
    ;;
esac
