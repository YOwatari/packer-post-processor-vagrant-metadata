.PHONY: all \
	deps build install \
	test/deps test \

VERSION:=$(shell git describe --tags --always --dirty) $(shell git name-rev --name-only HEAD)
BUILD_FLAGS:=-ldflags "-X main.Version \"$(VERSION)\""
ARTIFACTS_DIR:=$(CURDIR)/artifacts

GOOS:=linux darwin windows
GOARCH:=amd64


all: test build

deps:
	go get -d -v

build: deps
	go get github.com/mitchellh/gox
	go get github.com/cloudfoundry/gosigar
	gox -os="$(GOOS)" -arch="$(GOARCH)" -output="$(ARTIFACTS_DIR)/$(shell basename $(CURDIR))_{{.OS}}_{{.Arch}}" $(BUILD_FLAGS)

test/deps:
	go get -d -t -v

test: test/deps
	go test -v ./...

install: deps
	go install -v $(BUILD_FLAGS)
