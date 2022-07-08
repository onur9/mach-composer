package config

import (
	"context"
	"os"

	"github.com/labd/mach-composer/internal/model"
	"github.com/labd/mach-composer/internal/utils"
	"gopkg.in/yaml.v3"
)

type Global struct {
	Environment     string                    `yaml:"environment"`
	Cloud           string                    `yaml:"cloud"`
	Azure           *model.GlobalAzureConfig  `yaml:"azure"`
	TerraformConfig model.TerraformConfig     `yaml:"terraform_config"`
	AmplienceConfig *model.AmplienceConfig    `yaml:"amplience"`
	SentryConfig    *model.GlobalSentryConfig `yaml:"sentry"`
}

type Config struct {
	Filename     string
	MACHComposer MACHComposer      `yaml:"mach_composer"`
	Global       Global            `yaml:"global"`
	Sites        []model.Site      `yaml:"sites"`
	Components   []model.Component `yaml:"components"`

	Variables *Variables

	IsEncrypted bool
}

func (c *Config) HasSite(id string) bool {
	for i := range c.Sites {
		if c.Sites[i].Identifier == id {
			return true
		}
	}
	return false
}

type RawConfig struct {
	Filename     string
	MACHComposer MACHComposer `yaml:"mach_composer"`
	Global       Global       `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`
	Sops         yaml.Node    `yaml:"sops"`
}

type MACHComposer struct {
	Version       string
	VariablesFile string `yaml:"variables_file"`
}

// decryptYAML takes a filename and returns the decrypted YAML.
// This command directly calls the sops binary instead of using the
// go.mozilla.org/sops/v3/decrypt package since that adds numerous dependencies
// and adds ~19mb to the generated binary.
func decryptYAML(ctx context.Context, filename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return utils.RunSops(ctx, wd, "-d", filename, "--output-type=yaml")
}
