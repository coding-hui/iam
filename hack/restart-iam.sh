#!/bin/bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

set -o errexit
set -o nounset
set -o pipefail

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${IAM_ROOT}/hack/lib/init.sh"

TAG="${TAG:-latest}"
NAMESPACE="${NAMESPACE:-iam-system}" # 设置默认的命名空间
REGISTRY="${REPO:-devops-wecoding-docker.pkg.coding.net/wecoding/images}"
DEPLOYS="${DEPLOYS:-iam-apiserver}" # 部署列表，使用逗号分隔

function wait_for_installation_finish() {
  iam::log::info "waiting for IAM pod ready..."
  kubectl -n "$NAMESPACE" wait --timeout=180s --for=condition=Ready pods -l app.kubernetes.io/instance=iam

  start_time=$(date +%s)
  timeout_seconds=900 # 超时时间，单位为秒

  while :; do
    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))
    if [ "$elapsed_time" -ge "$timeout_seconds" ]; then
      iam::log::error "Timeout: IAM Pod did not become ready within $timeout_seconds seconds."
      exit 1
    fi

    iam::log::info "waiting for IAM ready..."

    pod_status=$(kubectl -n "$NAMESPACE" get pod -l app.kubernetes.io/instance=iam -o jsonpath="{.items[0].status.containerStatuses[0].ready}")
    if [ "$pod_status" == "true" ]; then
      break
    fi

    sleep 1
  done

  iam::log::info "IAM is ready!"
}

if [[ "$DEPLOYS" == *iam-apiserver* ]]; then
  # Update iam-apiserver image
  kubectl -n "$NAMESPACE" set image deployment/iam-apiserver apiserver="${REGISTRY}/iam-apiserver:${TAG}"
  # Restart iam-apiserver deployment
  kubectl -n "$NAMESPACE" rollout restart deployment/iam-apiserver
fi

wait_for_installation_finish
