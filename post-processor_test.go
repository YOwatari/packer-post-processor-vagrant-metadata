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
		"name":       "test",
		"output":     "_tmp/metadata.json",
		"url_prefix": "files://",
		"box_dir":    "_tmp",
		"version":    "0.0.0",
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

func TestConfigure_RequiredConfigs(t *testing.T) {
	s := []string{"name", "output", "url_prefix", "box_dir", "version"}
	for _, v := range s {
		p := new(PostProcessor)
		c := testConfig()
		delete(c, v)
		if err := p.Configure(c); err == nil {
			t.Fatalf("should have error when missing %s", v)
		}
	}
}

func TestPostProcess_BadBuilderId(t *testing.T) {
	artifact := &packer.MockArtifact{
		BuilderIdValue: "invalid",
	}

	_, _, err := testPostProcessor(t).PostProcess(testUi(), artifact)
	if !strings.Contains(err.Error(), "Unknown artifact type") {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcess_BadFiles(t *testing.T) {
	name, err := testCreateFile("_tmp/invalid", "")
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
	if !strings.Contains(err.Error(), "Unknown files in artifact") {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcess_MissingBox(t *testing.T) {
	name, err := testCreateFile("_tmp/missing.box", "")
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
			t.Fatalf("err: %s", pathErr)
		}
	}
}

func TestPostProcess_BadMetadata(t *testing.T) {
	metadataName, err := testCreateFile("_tmp/invalid.json", `{name: invalid json}`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c := testConfig()
	c["output"] = metadataName

	boxName, err := testCreateFile("_tmp/vagrant.box", "")
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
		t.Fatalf("err: %s", err)
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestPostProcess_EofMetadata(t *testing.T) {
	metadataName, err := testCreateFile("_tmp/eof.json", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName, err := testCreateFile("_tmp/vagrant.box", "")
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
		t.Fatalf("err: %s", err)
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestPostProcess_Json(t *testing.T) {
	metadataName, err := testCreateFile("_tmp/metadata.json", "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(metadataName); err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName, err := testCreateFile("_tmp/vagrant.box", "")
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

	actual, err := ioutil.ReadFile("_tmp/metadata.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []byte(`{
    "name": "test",
    "discription": "",
    "versions": [
        {
            "version": "0.0.0",
            "providers": [
                {
                    "name": "id",
                    "url": "files:///_tmp/vagrant.box",
                    "checksum_type": "sha256",
                    "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
                }
            ]
        }
    ]
}`)

	if string(actual) != string(expected) {
		t.Fatalf("should be the same as expected json string")
	}

	for _, v := range []string{metadataName, boxName} {
		if err := os.Remove(v); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

// TODO: 追加json
