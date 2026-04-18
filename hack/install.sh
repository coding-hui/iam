#!/usr/bin/env bash
# Copyright (c) 2023 coding-hui. All rights reserved.
# Simplified IAM installation script: install, start, stop, status, restart, logs, uninstall

set -euo pipefail

# Paths
IAM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${IAM_ROOT}/_output"
IAM_DIR="${HOME}/.iam"
IAM_BIN="${IAM_DIR}/bin/apiserver"
IAM_CONF="${IAM_DIR}/conf/apiserver.yaml"
IAM_DATA="${HOME}/.iam/data"
IAM_PID="${IAM_DIR}/apiserver.pid"

# ==============================================================================
# Commands

cmd_install() {
  mkdir -p "${IAM_DIR}/bin" "${IAM_DIR}/conf" "${HOME}/.iam/data" "${IAM_DIR}/logs"

  # Build
  echo "[INFO] Building apiserver..."
  mkdir -p "${OUTPUT_DIR}/bin"
  CGO_ENABLED=1 go build -ldflags "-X github.com/coding-hui/common/version.GitVersion=$(git describe --tags --always 2>/dev/null || echo v0.0.0) -X github.com/coding-hui/common/version.GitCommit=$(git rev-parse HEAD 2>/dev/null)" \
    -o "${OUTPUT_DIR}/bin/apiserver" "${IAM_ROOT}/cmd/apiserver"

  # Config
  echo "[INFO] Generating configuration..."
  sed -e "s|{{IAM_DATA}}|${IAM_DATA}|g" \
      "${IAM_ROOT}/conf/apiserver.yaml" > "${IAM_CONF}"

  # Binary
  cp "${OUTPUT_DIR}/bin/apiserver" "${IAM_BIN}"
  chmod +x "${IAM_BIN}"

  echo "[INFO] Installed to ${IAM_BIN}"
  echo "[INFO] Config: ${IAM_CONF}"
  echo "[INFO] Run: ${IAM_BIN} --config ${IAM_CONF}"
}

cmd_start() {
  if [ -f "${IAM_PID}" ] && kill -0 "$(cat "${IAM_PID}")" 2>/dev/null; then
    echo "[WARN] Already running (PID: $(cat ${IAM_PID}))"
    return 0
  fi

  echo "[INFO] Starting apiserver..."
  nohup "${IAM_BIN}" --config "${IAM_CONF}" > "${IAM_DIR}/logs/apiserver.log" 2>&1 &
  echo $! > "${IAM_PID}"

  # Wait for startup
  local i=0
  while [ $((i++)) -lt 30 ] && ! nc -z -w1 127.0.0.1 8080 2>/dev/null; do
    sleep 1
  done

  if nc -z -w1 127.0.0.1 8080 2>/dev/null; then
    echo "[INFO] Started (PID: $(cat ${IAM_PID}))"
    echo "[INFO] HTTP: http://127.0.0.1:8080"
  else
    echo "[ERROR] Failed to start"
    return 1
  fi
}

cmd_stop() {
  if [ -f "${IAM_PID}" ]; then
    local pid=$(cat "${IAM_PID}")
    if kill -0 "${pid}" 2>/dev/null; then
      echo "[INFO] Stopping (PID: ${pid})..."
      kill "${pid}" 2>/dev/null || true
      sleep 1
      kill -9 "${pid}" 2>/dev/null || true
    fi
    rm -f "${IAM_PID}"
  fi
  echo "[INFO] Stopped"
}

cmd_status() {
  if [ -f "${IAM_PID}" ] && kill -0 "$(cat "${IAM_PID}")" 2>/dev/null; then
    if nc -z -w1 127.0.0.1 8080 2>/dev/null; then
      echo "[INFO] Running (PID: $(cat ${IAM_PID}), http://127.0.0.1:8080)"
    else
      echo "[WARN] Running but not responding (PID: $(cat ${IAM_PID}))"
    fi
  else
    echo "[WARN] Not running"
    return 1
  fi
}

cmd_restart() {
  cmd_stop; sleep 1; cmd_start
}

cmd_logs() {
  if [ -f "${IAM_DIR}/logs/apiserver.log" ]; then
    tail -50 "${IAM_DIR}/logs/apiserver.log"
  else
    echo "[WARN] No logs found"
  fi
}

cmd_uninstall() {
  echo "[INFO] Uninstalling..."
  cmd_stop
  rm -rf "${IAM_DIR}"
  echo "[INFO] Done"
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