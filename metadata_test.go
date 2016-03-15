package main

import (
	"testing"
)

func TestMetadataAdd(t *testing.T) {
	var m *Metadata

	m = &Metadata{
		Versions: []*Version{
			&Version{
				Version: "v1",
				Providers: []*Provider{
					&Provider{Name: "test"},
				},
			},
		},
	}

	if e := m.add("v2", &Provider{Name: "test"}); e != nil {
		t.Fatalf("should not happen error")
	}

	if len(m.Versions) != 2 && m.Versions[1].Version != "v2" {
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
					&Provider{Name: "test"},
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
