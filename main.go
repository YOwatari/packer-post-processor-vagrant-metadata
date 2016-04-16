package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"os"
)

func main() {
	if len(os.Args[1:]) != 0 {
		os.Exit(Run(os.Args[1:]))
	}

	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(&PostProcessor{})
	server.Serve()
}
