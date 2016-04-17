package main

import (
	"bufio"
	"bytes"
	"github.com/mitchellh/packer/packer"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"
)

const validBuilderId = "mitchellh.post-processor.vagrant"

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":        "test",
		"description": "test",
		"output":      "_test/valid.json",
		"url_prefix":  "files://",
		"box_dir":     "_test",
		"version":     "0.0.0",
	}
}

func testPostProcessor(t *testing.T) *PostProcessor {
	var p PostProcessor
	if err := p.Configure(testConfig()); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &p
}

func testUi() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func testCreateFile(path, body string) (string, error) {
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.WriteString(body)
	writer.Flush()

	return f.Name(), nil
}

func TestImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestConfigure_requiredConfigures(t *testing.T) {
	s := []string{"name", "output", "url_prefix", "box_dir", "version"}
	for _, v := range s {
		p := new(PostProcessor)
		c := testConfig()
		delete(c, v)
		if err := p.Configure(c); err == nil {
			t.Errorf("should happen error when missing %s", v)
		}
	}
}

func TestPostProcess_badBuilderId(t *testing.T) {
	artifact := &packer.MockArtifact{
		BuilderIdValue: "invalid",
	}

	_, _, err := testPostProcessor(t).PostProcess(testUi(), artifact)
	if err != nil && !strings.Contains(err.Error(), "Unknown artifact type") {
		t.Errorf("should happen error about unknown artifact.\nerr: %s", err)
	}
}

func TestPostProcess_badFiles(t *testing.T) {
	name, err := testCreateFile("_test/invalid", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(name); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: validBuilderId,
		FilesValue:     []string{name},
	}

	_, _, err = testPostProcessor(t).PostProcess(testUi(), artifact)
	if err != nil && !strings.Contains(err.Error(), "Unknown files in artifact") {
		t.Errorf("should happen error about box files.\nerr: %s", err)
	}
}

func TestPostProcess_missingBox(t *testing.T) {
	name, err := testCreateFile("_test/missing.box", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(name); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: validBuilderId,
		FilesValue:     []string{string(name)},
	}

	_, _, err = testPostProcessor(t).PostProcess(testUi(), artifact)
	if err != nil {
		pathErr := err.(*os.PathError)
		if pathErr.Err.(syscall.Errno) != syscall.ENOENT {
			t.Errorf("should happen error about missing box file.\nerr: %s", pathErr)
		}
	}
}

func TestPostProcess_badMetadata(t *testing.T) {
	metadataName, err := testCreateFile("_test/invalid.json", `{name: invalid json}`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c := testConfig()
	c["output"] = metadataName

	boxName, err := testCreateFile("_test/vagrant.box", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var p PostProcessor
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: validBuilderId,
		FilesValue:     []string{boxName},
	}

	if _, _, err := p.PostProcess(testUi(), artifact); err == nil {
		t.Errorf("should happen error about invalid json format.\nerr: %s", err)
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestPostProcess_eofMetadata(t *testing.T) {
	metadataName, err := testCreateFile("_test/eof.json", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName, err := testCreateFile("_test/vagrant.box", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c := testConfig()
	c["output"] = metadataName

	var p PostProcessor
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: validBuilderId,
		FilesValue:     []string{boxName},
	}

	if _, _, err := p.PostProcess(testUi(), artifact); err == nil {
		t.Errorf("should happen error about EOF.\nerr: %s", err)
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestPostProcess_expectedJson(t *testing.T) {
	metadataName, err := testCreateFile("_test/valid.json", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(metadataName); err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName, err := testCreateFile("_test/vagrant.box", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c := testConfig()
	c["output"] = metadataName

	var p PostProcessor
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: validBuilderId,
		FilesValue:     []string{boxName},
	}

	if _, _, err := p.PostProcess(testUi(), artifact); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := ioutil.ReadFile("_test/valid.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []byte(`{
    "name": "test",
    "description": "test",
    "versions": [
        {
            "version": "0.0.0",
            "providers": [
                {
                    "name": "id",
                    "url": "files:///_test/vagrant.box",
                    "checksum_type": "sha256",
                    "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
                }
            ]
        }
    ]
}`)

	if string(actual) != string(expected) {
		t.Errorf("should be the same as expected json string")
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}
