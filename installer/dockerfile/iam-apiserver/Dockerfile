# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

FROM golang:1.19.4 AS builder
WORKDIR /go/src/wecoding/iam
COPY . .
RUN make build BINS="iam-apiserver"

ARG BASE_IMAGE

# You can replace distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# Overwrite `BASE_IMAGE` by passing `--build-arg=BASE_IMAGE=gcr.io/distroless/static:nonroot`
FROM ${BASE_IMAGE:-alpine:latest}
LABEL maintainer="wecoding@yeah.net"

WORKDIR /opt/iam/

ARG KUBE_VERSION="1.25.0"
ARG TARGETARCH
RUN apk add --update ca-certificates && \
    apk add --update openssl bash expat && \
    apk add --update -t deps curl && \
    curl -L https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/$TARGETARCH/kubectl -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    apk del --purge deps && \
    rm /var/cache/apk/* && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /go/src/wecoding/iam/_output/bin/iam-apiserver /usr/local/bin/iam-apiserver
COPY --from=builder /go/src/wecoding/iam/configs/iam-apiserver.yaml /etc/iam/iam-apiserver.yaml
COPY --from=builder /go/src/wecoding/iam/hack/gen-k8s-secret.sh /gen-k8s-secret.sh
COPY --from=builder /go/src/wecoding/iam/template/ /etc/iam/template
COPY --from=builder /go/src/wecoding/iam/template/ ./

RUN chmod +x /gen-k8s-secret.sh

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/iam-apiserver"]
CMD ["-c", "/etc/iam/iam-apiserver.yaml"]
