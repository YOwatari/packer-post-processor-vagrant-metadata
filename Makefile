.PHONY: all \
	setup/tool \
	deps build install clean \
	test/deps test \

VERSION=$(shell gobump show | $(CURDIR)/_tmp/bin/jq -r .version)
BUILD_FLAGS=-ldflags "-X main.Version \"v$(VERSION)\""
ARTIFACTS_DIR:=$(CURDIR)/artifacts

GOOS:=linux darwin windows
GOARCH:=amd64


all: test build

ifndef JQ_URL
  ifeq ($(shell expr substr $(shell uname -s) 1 5), Linux)
    JQ_URL=http://stedolan.github.io/jq/download/linux64/jq
  endif
  ifeq ($(shell uname), Darwin)
    JQ_URL=https://github.com/stedolan/jq/releases/download/jq-1.5/jq-osx-amd64
  endif
endif
setup/tool:
	go get github.com/motemen/gobump/cmd/gobump
	curl -L $(JQ_URL) -o $(CURDIR)/_tmp/bin/jq
	chmod +x $(CURDIR)/_tmp/bin/jq

deps: clean setup/tool
	go get -d -v ./...

build: deps
	go get github.com/mitchellh/gox
	go get github.com/cloudfoundry/gosigar
	gox -os="$(GOOS)" -arch="$(GOARCH)" -output="$(ARTIFACTS_DIR)/$(shell basename $(CURDIR))_{{.OS}}_{{.Arch}}" $(BUILD_FLAGS)

install: deps
	go install -v $(BUILD_FLAGS)

clean:
	-find $(CURDIR)/_tmp -maxdepth 1 -mindepth 1 ! -name .gitkeep | xargs rm -rf
	-find $(CURDIR)/artifacts -maxdepth 1 -mindepth 1 ! -name .gitkeep | xargs rm -rf

test/deps: clean
	go get -d -t -v ./...

test: test/deps
	go test -v ./...
