package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestVersionNumberCommand_implement(t *testing.T) {
	var _ cli.Command = &VersionNumberCommand{}
}
