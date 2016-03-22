package main

import (
	"github.com/mitchellh/packer/packer/plugin"
)

const version = "0.1.1"

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(&PostProcessor{})
	server.Serve()
}
