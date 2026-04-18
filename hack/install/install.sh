#!/usr/bin/env bash
# Copyright (c) 2023 coding-hui. All rights reserved.
# Simplified IAM installation script: install, start, stop, status, restart, logs, uninstall

set -euo pipefail

# Paths
IAM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUTPUT_DIR="${IAM_ROOT}/_output"
IAM_DIR="${HOME}/.iam"
IAM_BIN="${IAM_DIR}/bin/iam-apiserver"
IAM_CONF="${IAM_DIR}/conf/iam-apiserver.yaml"
IAM_CERT="${IAM_DIR}/conf/cert"
IAM_DATA="${HOME}/.iam/data"
IAM_LOG="${IAM_DIR}/logs"

# ==============================================================================
# Commands

cmd_install() {
  mkdir -p "${IAM_DIR}/bin" "${IAM_DIR}/conf/cert" "${HOME}/.iam/data" "${IAM_DIR}/logs"

  # Build
  echo "[INFO] Building iam-apiserver..."
  mkdir -p "${OUTPUT_DIR}/bin"
  CGO_ENABLED=1 go build -ldflags "-X github.com/coding-hui/common/version.GitVersion=$(git describe --tags --always 2>/dev/null || echo v0.0.0) -X github.com/coding-hui/common/version.GitCommit=$(git rev-parse HEAD 2>/dev/null)" \
    -o "${OUTPUT_DIR}/bin/iam-apiserver" "${IAM_ROOT}/cmd/apiserver"

  # CA cert
  echo "[INFO] Generating certificates..."
  openssl genrsa -out "${IAM_CERT}/ca-key.pem" 2048 2>/dev/null
  openssl req -x509 -new -nodes -key "${IAM_CERT}/ca-key.pem" -sha256 -days 3650 -out "${IAM_CERT}/ca.pem" -subj "/CN=IAM CA" 2>/dev/null
  openssl genrsa -out "${IAM_CERT}/iam-apiserver-key.pem" 2048 2>/dev/null
  openssl req -new -key "${IAM_CERT}/iam-apiserver-key.pem" -out "${IAM_CERT}/iam-apiserver.csr" -subj "/CN=IAM APIServer" 2>/dev/null
  openssl x509 -req -in "${IAM_CERT}/iam-apiserver.csr" -CA "${IAM_CERT}/ca.pem" -CAkey "${IAM_CERT}/ca-key.pem" -CAcreateserial \
    -out "${IAM_CERT}/iam-apiserver.pem" -days 825 -sha256 2>/dev/null
  rm -f "${IAM_CERT}/iam-apiserver.csr" "${IAM_CERT}/ca-key.pem"

  # Config
  echo "[INFO] Generating configuration..."
  cat > "${IAM_CONF}" << EOF
server:
  insecure:
    bind_address: 127.0.0.1
    bind_port: 8080
  secure:
    bind_address: 0.0.0.0
    bind_port: 8443
    tls:
      cert_file: ${IAM_CERT}/iam-apiserver.pem
      key_file: ${IAM_CERT}/iam-apiserver-key.pem
  grpc:
    bind_address: 0.0.0.0
    bind_port: 8081
database:
  type: sqlite
  dsn: ${IAM_DATA}/iam.db
cache:
  type: memory
  redis:
    host: 127.0.0.1
    port: 6379
log:
  level: info
  path: ${IAM_LOG}
EOF

  # Binary
  cp "${OUTPUT_DIR}/bin/iam-apiserver" "${IAM_BIN}"
  chmod +x "${IAM_BIN}"

  # Service
  if [[ "$(uname -s)" == "Darwin" ]]; then
    cat > "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" << PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key><string>io.github.coding-hui.iam-apiserver</string>
    <key>ProgramArguments</key><array><string>${IAM_BIN}</string><string>--config</string><string>${IAM_CONF}</string></array>
    <key>RunAtLoad</key><true/>
    <key>StandardOutPath</key><string>${IAM_LOG}/iam-apiserver.log</string>
    <key>StandardErrorPath</key><string>${IAM_LOG}/iam-apiserver.error.log</string>
</dict>
</plist>
PLIST
    mkdir -p "${HOME}/Library/LaunchAgents"
    cp "${IAM_ROOT}/io.github.coding-hui.iam-apiserver.plist" "${HOME}/Library/LaunchAgents/"
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    cat > /tmp/iam-apiserver.service << EOF
[Unit] Description=IAM API Server After=network.target
[Service] Type=simple ExecStart=${IAM_BIN} --config ${IAM_CONF} WorkingDirectory=${IAM_DIR} Restart=on-failure
[Install] WantedBy=multi-user.target
EOF
    sudo cp /tmp/iam-apiserver.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable iam-apiserver
    sudo systemctl start iam-apiserver
  fi

  # Wait
  echo "[INFO] Waiting for service..."
  local i=0
  while [ $((i++)) -lt 30 ] && ! nc -z -w1 127.0.0.1 8080 2>/dev/null; do
    sleep 1
  done

  echo "[INFO] IAM apiserver installed"
  echo "  HTTP:  http://127.0.0.1:8080"
  echo "  HTTPS: https://127.0.0.1:8443"
}

cmd_uninstall() {
  echo "[INFO] Uninstalling..."
  if [[ "$(uname -s)" == "Darwin" ]]; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist" 2>/dev/null || true
    rm -f "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl stop iam-apiserver 2>/dev/null || true
    sudo systemctl disable iam-apiserver 2>/dev/null || true
    sudo rm -f /etc/systemd/system/iam-apiserver.service
  fi
  rm -rf "${IAM_DIR}"
  echo "[INFO] Done"
}

cmd_start() {
  if [[ "$(uname -s)" == "Darwin" ]]; then
    launchctl load "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl start iam-apiserver
  fi
  echo "[INFO] Started"
}

cmd_stop() {
  if [[ "$(uname -s)" == "Darwin" ]]; then
    launchctl unload "${HOME}/Library/LaunchAgents/io.github.coding-hui.iam-apiserver.plist"
  else
    sudo systemctl stop iam-apiserver
  fi
  echo "[INFO] Stopped"
}

cmd_status() {
  if nc -z -w1 127.0.0.1 8080 2>/dev/null; then
    echo "[INFO] Running (http://127.0.0.1:8080)"
  else
    echo "[WARN] Not running"
    return 1
  fi
}

cmd_restart() {
  cmd_stop; sleep 1; cmd_start
}

cmd_logs() {
  tail -50 "${IAM_LOG}/iam-apiserver.log" 2>/dev/null || echo "[WARN] No logs found"
}

# ==============================================================================
# Dispatch

case "${1:-}" in
  install)   cmd_install ;;
  uninstall) cmd_uninstall ;;
  start)     cmd_start ;;
  stop)      cmd_stop ;;
  status)    cmd_status ;;
  restart)   cmd_restart ;;
  logs)      cmd_logs ;;
  *)         echo "Usage: $0 {install|uninstall|start|stop|status|restart|logs}" ;;
esac
