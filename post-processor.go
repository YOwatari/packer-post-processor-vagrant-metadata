package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Name         string `mapstructure:"name"`
	MetadataPath string `mapstructure:"output"`
	UrlPrefix    string `mapstructure:"url_prefix"`
	BoxDir       string `mapstructure:"box_dir"`
	Version      string `mapstructure:"version"`

	Description string `mapstructure:"description"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"description",
			},
		},
	}, raws...)

	if err != nil {
		return err
	}

	errs := new(packer.MultiError)

	// required configuration
	templates := map[string]*string{
		"name":       &p.config.Name,
		"output":     &p.config.MetadataPath,
		"url_prefix": &p.config.UrlPrefix,
		"box_dir":    &p.config.BoxDir,
		"version":    &p.config.Version,
	}

	for key, ptr := range templates {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", key))
		}
	}

	for key, ptr := range templates {
		if err = interpolate.Validate(*ptr, &p.config.ctx); err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Error parsing %s template: %s", key, err))
		}
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	if artifact.BuilderId() != "mitchellh.post-processor.vagrant" {
		return nil, false, fmt.Errorf(
			"Unknown artifact type, requires box from vagrant post-processor: %s", artifact.BuilderId())
	}

	box := artifact.Files()[0]
	if !strings.HasSuffix(box, ".box") {
		return nil, false, fmt.Errorf(
			"Unknown files in artifact from vagrant post-processor: %s", artifact.Files())
	}

	provider := providerFromBuilderName(artifact.Id())

	file, err := os.Open(box)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, false, err
	}
	size := info.Size()
	ui.Message(fmt.Sprintf("Box size: %s (%d bytes)", box, size))

	metadata, err := p.getMetadata()
	if err != nil {
		return nil, false, err
	}

	ui.Message("Generating checksum")
	checksum, err := sum256(file)
	if err != nil {
		return nil, false, err
	}
	ui.Message(fmt.Sprintf("Checksum is %s", checksum))

	ui.Message(fmt.Sprintf("Adding %s %s box to metadata", provider, p.config.Version))
	if err := metadata.add(p.config.Version, &Provider{
		Name:         provider,
		Url:          fmt.Sprintf("%s/%s/%s", p.config.UrlPrefix, p.config.BoxDir, path.Base(box)),
		ChecksumType: "sha256",
		Checksum:     checksum,
	}); err != nil {
		return nil, false, err
	}

	ui.Message(fmt.Sprintf("Saving the metadata: %s", p.config.MetadataPath))
	if err := p.putMetadata(metadata); err != nil {
		return nil, false, err
	}

	return &Artifact{fmt.Sprintf("%s/%s", p.config.UrlPrefix, p.config.MetadataPath)}, true, nil
}

func (p *PostProcessor) getMetadata() (*Metadata, error) {
	body, err := os.Open(p.config.MetadataPath)
	if err != nil {
		return &Metadata{Name: p.config.Name}, nil
	}
	defer body.Close()

	metadata := &Metadata{}
	if err := json.NewDecoder(body).Decode(metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (p *PostProcessor) putMetadata(metadata *Metadata) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(metadata); err != nil {
		return err
	}

	if err := ioutil.WriteFile(p.config.MetadataPath, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func sum256(file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func providerFromBuilderName(name string) string {
	switch name {
	case "aws":
		return "aws"
	case "digitalocean":
		return "digitalocean"
	case "virtualbox":
		return "virtualbox"
	case "vmware":
		return "vmware_desktop"
	case "parallels":
		return "parallels"
	default:
		return name
	}
}
