.PHONY: all \
	setup clean \
	deps build install \
	test/deps test \

ARTIFACTS_DIR:=$(CURDIR)/artifacts

VERSION=$(shell gobump show | jq -r .version)
COMMIT=$(shell git rev-parse --verify HEAD)
BUILD_FLAGS=-ldflags "-X main.version=\"$(VERSION)\" -X main.GitCommit=\"$(COMMIT)\""
GOOS:=linux darwin windows
GOARCH:=amd64


all: test build

setup:
	@which jq >/dev/null 2>&1 || (echo "you need jq. https://stedolan.github.io/jq/"; exit 1;)
	go get github.com/motemen/gobump/cmd/gobump

deps:
	go get -d -v ./...

build: clean setup deps
	go get github.com/mitchellh/gox
	go get github.com/cloudfoundry/gosigar
	gox -os="$(GOOS)" -arch="$(GOARCH)" -output="$(ARTIFACTS_DIR)/$(shell basename $(CURDIR))_{{.OS}}_{{.Arch}}" $(BUILD_FLAGS)

install: setup deps
	go install -v $(BUILD_FLAGS)

clean:
	-find $(ARTIFACTS_DIR) -maxdepth 1 -mindepth 1 ! -name .gitkeep | xargs rm -rf

test/deps:
	go get -d -t -v ./...

test: test/deps
	go test -v ./...
