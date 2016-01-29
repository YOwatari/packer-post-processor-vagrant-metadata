package main

import (
	"fmt"
)

const BuilderId = "yowatari.post-processor.vagrant-extend"

type Artifact struct {
	Url string
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Files() []string {
	return []string{a.Url}
}

func (a *Artifact) Id() string {
	return ""
}

func (a *Artifact) String() string {
	return fmt.Sprintf("vagrant metadata url: %s", a.Url)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	return nil
}
