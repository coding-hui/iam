BIN_DIR=_output/bin
RELEASE_DIR=_output/release
REPO_PATH=github.com/wecoding/iam
IMAGE_PREFIX=wecoding
CC ?= "gcc"
GOOS ?= linux
SUPPORT_PLUGINS ?= "no"
BUILDX_OUTPUT_TYPE ?= "docker"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Get OS architecture
OSARCH=$(shell uname -m)
ifeq ($(OSARCH),x86_64)
GOARCH?=amd64
else ifeq ($(OSARCH),x64)
GOARCH?=amd64
else ifeq ($(OSARCH),aarch64)
GOARCH?=arm64
else ifeq ($(OSARCH),aarch64_be)
GOARCH?=arm64
else ifeq ($(OSARCH),armv8b)
GOARCH?=arm64
else ifeq ($(OSARCH),armv8l)
GOARCH?=arm64
else ifeq ($(OSARCH),i386)
GOARCH?=x86
else ifeq ($(OSARCH),i686)
GOARCH?=x86
else ifeq ($(OSARCH),arm)
GOARCH?=arm
else
GOARCH?=$(OSARCH)
endif

# Run `make images DOCKER_PLATFORMS="linux/amd64,linux/arm64" BUILDX_OUTPUT_TYPE=registry IMAGE_PREFIX=[yourregistry]` to push multi-platform
DOCKER_PLATFORMS ?= "linux/${GOARCH}"

include Makefile.def

.EXPORT_ALL_VARIABLES:

all: apiserver

init:
	mkdir -p ${BIN_DIR}
	mkdir -p ${RELEASE_DIR}

apiserver: init
	if [ ${SUPPORT_PLUGINS} = "yes" ];then\
		CC=${CC} CGO_ENABLED=1 go build -ldflags ${LD_FLAGS} -o ${BIN_DIR}/apiserver ./cmd/apiserver;\
	else\
		CC=${CC} CGO_ENABLED=0 go build -ldflags ${LD_FLAGS} -o ${BIN_DIR}/apiserver ./cmd/apiserver;\
	fi;

apiserver-win: init
	if [ ${SUPPORT_PLUGINS} = "yes" ];then\
		env GOOS=windows CC=${CC} CGO_ENABLED=1 go build -ldflags ${LD_FLAGS} -o ${BIN_DIR}/apiserver.exe ./cmd/apiserver/main.go;\
	else\
		env GOOS=windows CC=${CC} CGO_ENABLED=0 go build -ldflags ${LD_FLAGS} -o ${BIN_DIR}/apiserver.exe ./cmd/apiserver/main.go;\
	fi;

image_bins: apiserver

apiserver-image:
	docker build -t "${IMAGE_PREFIX}/apiserver:$(RELEASE_VER)" . -f ./installer/dockerfile/apiserver/Dockerfile\
		--build-arg=VERSION=$(RELEASE_VER)\
		--build-arg=GITVERSION=$(GITVERSION)\
		--build-arg=BUILDPLATFORM=$(DOCKER_PLATFORMS)

unit-test:
	go clean -testcache
	go test -gcflags=all=-l -coverprofile=coverage.txt $(shell go list ./pkg/... ./cmd/...)

clean:
	rm -rf _output/
	rm -f *.log

build-swagger:
	go get -u github.com/swaggo/swag/cmd/swag
	swag i -g server.go -dir ./pkg/apiserver --parseDependency --parseInternal -o ./docs/apidoc

run-apiserver:
	go run ./cmd/apiserver/main.go
