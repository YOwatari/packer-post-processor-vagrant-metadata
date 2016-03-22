package main

import (
	"fmt"
)

type Metadata struct {
	Name        string     `json:"name"`
	Description string     `json:"discription"`
	Versions    []*Version `json:"versions"`
}

type Version struct {
	Version   string      `json:"version"`
	Providers []*Provider `json:"providers"`
}

type Provider struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	ChecksumType string `json:"checksum_type"`
	Checksum     string `json:"checksum"`
}

func (m *Metadata) Add(version string, provider *Provider) error {
	for _, vs := range m.Versions {
		if vs.Version == version {
			for _, p := range vs.Providers {
				if p.Name == provider.Name {
					return fmt.Errorf("%s box for version %s already exists in metadata", p.Name, version)
				}
			}
			vs.Providers = append(vs.Providers, provider)
			return nil
		}
	}
	m.Versions = append(m.Versions, &Version{
		Version:   version,
		Providers: []*Provider{provider},
	})
	return nil
}
