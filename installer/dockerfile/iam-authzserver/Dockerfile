# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

ARG BASE_IMAGE

# You can replace distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# Overwrite `BASE_IMAGE` by passing `--build-arg=BASE_IMAGE=gcr.io/distroless/static:nonroot`
FROM --platform=${BUILDPLATFORM:-linux/amd64} ${BASE_IMAGE:-alpine:latest}
LABEL maintainer="wecoding@yeah.net"

WORKDIR /opt/iam/

# This is required by daemon connecting with cri
# RUN apk add --no-cache ca-certificates bash expat

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY config.* /etc/iam/iam-authzserver.yaml
COPY iam-authzserver.* /usr/local/bin/iam-authzserver

EXPOSE 8010

ENTRYPOINT ["/usr/local/bin/iam-authzserver"]
CMD ["-c", "/etc/iam/iam-authzserver.yaml"]
