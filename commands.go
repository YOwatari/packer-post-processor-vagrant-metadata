package main

import (
	"github.com/YOwatari/packer-post-processor-vagrant-metadata/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"version-number": func() (cli.Command, error) {
			return &command.VersionNumberCommand{
				Meta:    *meta,
				Version: version,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
