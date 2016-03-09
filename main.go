package main

import (
	"github.com/mitchellh/packer/packer/plugin"
)

var Version string = "0.1"

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(&PostProcessor{})
	server.Serve()
}
