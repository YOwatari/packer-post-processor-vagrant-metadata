package main

import (
	"testing"
)

func TestMetadataAdd(t *testing.T) {
	var expected, m *Metadata
	var p *Provider

	p = &Provider{Name: "test"}
	expected = &Metadata{
		Versions: []*Version{
			&Version{
				Version:   "v1",
				Providers: []*Provider{p},
			},
			&Version{
				Version:   "v2",
				Providers: []*Provider{p},
			},
		},
	}

	m = &Metadata{
		Versions: []*Version{
			&Version{
				Version:   "v1",
				Providers: []*Provider{p},
			},
		},
	}

	if e := m.add("v2", &Provider{Name: "test"}); e != nil {
		t.Errorf("should not happen error")
	}

	// TODO: fix
	if m.Versions[0].Version != expected.Versions[0].Version {
		t.Errorf("should be enable to add version to metadata")
	}
}

func TestAlreadyExsits(t *testing.T) {
	var m *Metadata

	m = &Metadata{
		Versions: []*Version{
			&Version{
				Version: "v1",
				Providers: []*Provider{
					&Provider{
						Name: "test",
					},
				},
			},
		},
	}

	if e := m.add("v1", &Provider{Name: "test"}); e != nil {
		if e.Error() != "test box for version v1 already exists in metadata" {
			t.Fatalf("should happen already exists error")
		}
	}
}
