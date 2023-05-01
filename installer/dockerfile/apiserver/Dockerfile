# Copyright (c) 2023 coding-hui. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

ARG BASE_IMAGE
# Build the manager binary
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.19.4 as builder

WORKDIR /go/src/wecoding/
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# It's a proxy for CN developer, please unblock it if you have network issue
ARG GOPROXY
ENV GOPROXY=${GOPROXY:-https://goproxy.cn}

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source for building core
COPY cmd/apiserver/ cmd/apiserver/
COPY docs/ docs/
COPY pkg/ pkg/
COPY version/ version/
COPY Makefile Makefile
COPY Makefile.def Makefile.def

# Build
ARG VERSION
ARG GITVERSION
RUN make apiserver GOOS=linux RELEASE_VER=${VERSION} GITVERSION=${GITVERSION}

# Use alpine as base image due to the discussion in issue #1448
# You can replace distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# Overwrite `BASE_IMAGE` by passing `--build-arg=BASE_IMAGE=gcr.io/distroless/static:nonroot`
FROM ${BASE_IMAGE:-alpine@sha256:e2e16842c9b54d985bf1ef9242a313f36b856181f188de21313820e177002501}
# This is required by daemon connecting with cri
RUN apk add --no-cache ca-certificates bash expat

WORKDIR /go/src/wecoding/

COPY --from=builder /go/src/wecoding/_output/bin/apiserver /usr/local/bin/apiserver

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/apiserver"]
