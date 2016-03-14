.PHONY: all \
	setup clean \
	deps build install \
	test/deps test \

ARTIFACTS_DIR:=$(CURDIR)/artifacts
TMP_DIR:=$(CURDIR)/_tmp

VERSION=$(shell gobump show | $(TMP_DIR)/bin/jq -r .version)
BUILD_FLAGS=-ldflags "-X main.Version \"v$(VERSION)\""
GOOS:=linux darwin windows
GOARCH:=amd64


all: test build

setup: clean tool

ifndef JQ_URL
  ifeq ($(shell expr substr $(shell uname -s) 1 5), Linux)
    JQ_URL=http://stedolan.github.io/jq/download/linux64/jq
  endif
  ifeq ($(shell uname), Darwin)
    JQ_URL=https://github.com/stedolan/jq/releases/download/jq-1.5/jq-osx-amd64
  endif
endif
tool:
	go get github.com/motemen/gobump/cmd/gobump
	curl -L $(JQ_URL) -o $(TMP_DIR)/bin/jq
	chmod +x $(TMP_DIR)/bin/jq

deps:
	go get -d -v ./...

build: setup deps
	go get github.com/mitchellh/gox
	go get github.com/cloudfoundry/gosigar
	gox -os="$(GOOS)" -arch="$(GOARCH)" -output="$(ARTIFACTS_DIR)/$(shell basename $(CURDIR))_{{.OS}}_{{.Arch}}" $(BUILD_FLAGS)

install: setup deps
	go install -v $(BUILD_FLAGS)

clean:
	-find $(TMP_DIR) -maxdepth 1 -mindepth 1 ! -name .gitkeep | xargs rm -rf
	-find $(ARTIFACTS_DIR) -maxdepth 1 -mindepth 1 ! -name .gitkeep | xargs rm -rf

test/setup: clean tool

test/deps:
	go get -d -t -v ./...

test: setup test/deps
	go test -v ./...
