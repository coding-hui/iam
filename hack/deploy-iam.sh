#!/bin/bash
# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# Define deployment directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Set namespace variables (use default values if not set)
IAM_DEPEND_NAMESPACE="${IAM_DEPEND_NAMESPACE:-infra}"
IAM_SYSTEM_NAMESPACE="${IAM_SYSTEM_NAMESPACE:-iam-system}"

# Set installation flags, read from environment variables or default to false
INSTALL_MYSQL="${INSTALL_MYSQL:-false}"
INSTALL_REDIS="${INSTALL_REDIS:-false}"

# Function to wait for service startup
wait_for_service_start() {
  local namespace=$1
  local label_selector=$2
  local timeout=$3

  local start_time=$(date +%s)
  local elapsed_time=0

  while true; do
    local pod_status=$(kubectl get pods -n "$namespace" -l "$label_selector" -o jsonpath='{.items[0].status.phase}' 2>/dev/null)
    if [ "$pod_status" = "Running" ]; then
      echo "Service started successfully!"
      return 0
    fi

    local current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))

    echo "Waiting for service startup..."

    if [ "$elapsed_time" -gt "$timeout" ]; then
      echo "Timeout reached, service startup failed!"
      return 1
    fi

    sleep 5
  done
}

# Install MySQL
install_mysql() {
  if [ "$INSTALL_MYSQL" = true ]; then
    echo "Installing MySQL..."
    pushd "$IAM_ROOT/installer/helm/mysql" >/dev/null
    helm install mysql . -n "$IAM_DEPEND_NAMESPACE" --create-namespace
    popd >/dev/null

    wait_for_service_start "$IAM_DEPEND_NAMESPACE" "app.kubernetes.io/instance=mysql" 300
  else
    echo "Skipping MySQL installation"
  fi
}

# Install Redis
install_redis() {
  if [ "$INSTALL_REDIS" = true ]; then
    echo "Installing Redis..."
    pushd "$IAM_ROOT/installer/helm/redis" >/dev/null
    helm install redis . -n "$IAM_DEPEND_NAMESPACE" --create-namespace
    popd >/dev/null

    wait_for_service_start "$IAM_DEPEND_NAMESPACE" "app.kubernetes.io/instance=redis" 300
  else
    echo "Skipping Redis installation"
  fi
}

# Install IAM
install_iam() {
  echo "Installing IAM..."
  pushd "$IAM_ROOT/installer/helm/iam" >/dev/null
  helm install iam . -n "$IAM_SYSTEM_NAMESPACE" --create-namespace
  popd >/dev/null

  wait_for_service_start "$IAM_SYSTEM_NAMESPACE" "app.kubernetes.io/instance=iam" 300
}

# Execute installation steps
install_mysql
install_redis
install_iam

echo "IAM deployment completed!"
