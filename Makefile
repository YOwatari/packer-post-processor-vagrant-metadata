.PHONY: all echo depends build

NAME:=$(shell basename $(CURDIR))
VERSION:=$(shell git log --pretty=format:"%h (%ad)" --date=short -1)
GOOS:=linux darwin windows
GOARCH:=amd64


all: build

echo:
	@echo "   name: $(NAME)"
	@echo "version: $(VERSION)"
	@echo " GOROOT: $$GOROOT"
	@echo " GOPATH: $$GOPATH"
	@go version

depends:
	go get github.com/mattn/gom
	go get github.com/mitchellh/gox
	go get github.com/cloudfoundry/gosigar
	gom install

build: depends echo
	gox -os="$(GOOS)" -arch="$(GOARCH)" -output="artifacts/{{.OS}}-{{.Arch}}/$(NAME)" -ldflags "-X main.version '$(VERSION)'"
