package main

import (
	"strings"
	"testing"
)

func TestAdd_initial(t *testing.T) {
	m := &Metadata{}

	if err := m.Add("v1", &Provider{Name: "test"}); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestAdd_appendProvider(t *testing.T) {
	m := &Metadata{}

	if err := m.Add("v1", &Provider{Name: "test1"}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := m.Add("v1", &Provider{Name: "test2"}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(m.Versions[0].Providers) != 2 {
		t.Errorf("should be enale to append the other provider")
	}
}

func TestAdd_appendVersion(t *testing.T) {
	m := &Metadata{}

	if err := m.Add("v1", &Provider{Name: "test"}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := m.Add("v2", &Provider{Name: "test"}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(m.Versions) != 2 {
		t.Errorf("should be enable to add version to metadata")
	}

}

func TestAdd_alreadyExsits(t *testing.T) {
	m := &Metadata{}

	if err := m.Add("v1", &Provider{Name: "test"}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := m.Add("v1", &Provider{Name: "test"}); err != nil {
		if !strings.Contains(err.Error(), "already exists in metadata") {
			t.Errorf("err: %s", err)
		}
	}
}
