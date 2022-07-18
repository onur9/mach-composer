package global

import "github.com/labd/mach-composer/internal/config/global/terraformconfig"

type TerraformConfig struct {
	AzureRemoteState terraformconfig.AzureRemoteState `yaml:"azure_remote_state"`
	AWSRemoteState   terraformconfig.AWSRemoteState   `yaml:"aws_remote_state"`
	Providers        terraformconfig.Providers        `yaml:"providers"`
}
