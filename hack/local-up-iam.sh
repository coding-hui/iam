# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
DOCKER_COMPOSE_DIR=${IAM_ROOT}/installer/docker-compose

source "${IAM_ROOT}/hack/lib/init.sh"

# prepare docker images
function prepare {
  iam::log::info "Preparing..."
  iam::log::info "Building docker images"
  make image
}

function install_iam {
  iam::log::info "Install IAM..."
  docker-compose -f ${DOCKER_COMPOSE_DIR}/docker-compose.yml up -d
}

function uninstall_iam {
  iam::log::info "Uninstall IAM..."
  docker-compose -f ${DOCKER_COMPOSE_DIR}/docker-compose.yml down
}

echo $* | grep -E -q "\-\-quit|\-q"
if [[ $? -eq 0 ]]; then
  uninstall_iam
  exit 0
fi

prepare
install_iam
