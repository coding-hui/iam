#!/bin/bash

# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

workdir=$(cd $(dirname -- $0) && cd ../ && pwd -P)
cd ${workdir}

make build BINS="iam-apiserver" PLATFORMS="linux,amd64"

echo "You need to wait for the file to be synchronized to the remote k8s cluster and then start the application."

#replace_count=$(cat .nocalhost/config.yaml | grep replace_to_your_name | wc -l)
#if [ ${replace_count} -eq 1 ]; then
#  current_user=$(whoami)
#  sed -i '' "s/replace_to_your_name/${current_user}/g" .nocalhost/config.yaml
#  echo "CURRENT_USER is set to ${current_user} in config.yaml, please re-enter nocalhost DevMode to make it take effect."
#fi
