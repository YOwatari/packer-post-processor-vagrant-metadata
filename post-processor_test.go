package main

import (
	"bufio"
	"bytes"
	"github.com/mitchellh/packer/packer"
	"os"
	"strings"
	"syscall"
	"testing"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":       "test",
		"output":     "_tmp/metadata.json",
		"url_prefix": "files://",
		"box_dir":    "_tmp/",
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
	artifact := &packer.MockArtifact{
		BuilderIdValue: "mitchellh.post-processor.vagrant",
		FilesValue:     []string{"invalid"},
	}

	_, _, err := testPostProcessor(t).PostProcess(testUi(), artifact)
	if !strings.Contains(err.Error(), "Unknown files in artifact") {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcess_MissingBox(t *testing.T) {
	f, err := os.Create("_tmp/missing.box")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	name := f.Name()

	if err := f.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(name); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: "mitchellh.post-processor.vagrant",
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
	f, err := os.Create("_tmp/invalid.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	writer := bufio.NewWriter(f)
	writer.WriteString(`{name: invalid json}`)
	writer.Flush()

	metadataName := f.Name()
	c := testConfig()
	c["output"] = metadataName

	if err := f.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if f, err = os.Create("_tmp/vagrant.box"); err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName := f.Name()

	if err := f.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	var p PostProcessor
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: "mitchellh.post-processor.vagrant",
		FilesValue:     []string{boxName},
	}

	if _, _, err := p.PostProcess(testUi(), artifact); err == nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(metadataName); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(boxName); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcess_Hoge(t *testing.T) {
	f, err := os.Create("_tmp/valid.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	metadataName := f.Name()
	c := testConfig()
	c["output"] = metadataName

	if err := f.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(metadataName); err != nil {
		t.Fatalf("err: %s", err)
	}

	if f, err = os.Create("_tmp/vagrant.box"); err != nil {
		t.Fatalf("err: %s", err)
	}

	boxName := f.Name()

	if err := f.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	var p PostProcessor
	if err := p.Configure(c); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: "mitchellh.post-processor.vagrant",
		FilesValue:     []string{boxName},
	}

	if _, _, err := p.PostProcess(testUi(), artifact); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Remove(boxName); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TODO: sha256確認
// TODO: 初期生成json
// TODO: 追加json
