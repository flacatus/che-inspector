GO := GOFLAGS="-mod=vendor" go
CHEINSPECTOR := $(addprefix bin/, che-inspector)
PKG := github.com/flacatus/che-inspector
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
CLI_VERSION := $(or $(SOURCE_CLI_VERSION),$(shell git describe --always --tags HEAD))
TAGS := -tags "json1"

# -race is only supported on linux/amd64, linux/ppc64le, linux/arm64, freebsd/amd64, netbsd/amd64, darwin/amd64 and windows/amd64
ifeq ($(shell go env GOARCH),s390x)
TEST_RACE :=
else
TEST_RACE := -race
endif
$(CHEINSPECTOR): che-inspector_version_flags=-ldflags "-X '$(PKG)/pkg/cmd/version.gitCommit=$(GIT_COMMIT)' -X '$(PKG)/pkg/cmd/version.cliVersion=$(CLI_VERSION)' -X '$(PKG)/pkg/cmd/version.buildDate=$(BUILD_DATE)'"
$(CHEINSPECTOR):
	go mod vendor
	$(GO) build $(che-inspector_version_flags) $(extra_flags) $(TAGS) -o $@ main.go
.PHONY: build
build: clean $(CHEINSPECTOR)
.PHONY: cross
cross: che-inspector_version_flags=-ldflags "-X '$(PKG)/pkg/cmd/version.gitCommit=$(GIT_COMMIT)' -X '$(PKG)/pkg/cmd/version.cliVersion=$(CLI_VERSION)' -X '$(PKG)/pkg/cmd/version.buildDate=$(BUILD_DATE)'"
cross:
ifeq ($(shell go env GOARCH),amd64)
	go mod vendor
	GOOS=darwin CC=o64-clang CXX=o64-clang++ CGO_ENABLED=1 $(GO) -mod=readonly build $(che-inspector_version_flags) -o "bin/darwin-amd64-che-inspector" --ldflags "-extld=o64-clang" main.go
endif
.PHONY: clean
clean:
	@rm -rf ./bin
.PHONY: static
static: extra_flags=-ldflags '-w -extldflags "-static"' -tags "json1"
