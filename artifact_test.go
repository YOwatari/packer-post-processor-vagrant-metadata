package main

import (
	"github.com/mitchellh/packer/packer"
	"strings"
	"testing"
)

func TestArtifact_ImplementsArtifact(t *testing.T) {
	var raw interface{}
	raw = &Artifact{}
	if _, ok := raw.(packer.Artifact); !ok {
		t.Errorf("Artifact should be a Artifact")
	}
}

func TestArtifact_URL(t *testing.T) {
	artifact := &Artifact{"https://www.packer.io/"}
	if !strings.Contains(artifact.String(), "vagrant metadata url: https://www.packer.io/") {
		t.Errorf("should return metadata info")
	}
}
