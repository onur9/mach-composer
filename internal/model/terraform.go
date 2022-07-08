package model

// TerraformProvider version overwrites.
type TerraformProvider struct {
	AWS           string `yaml:"aws"`
	Azure         string `yaml:"azure"`
	Commercetools string `yaml:"commercetools"`
	Sentry        string `yaml:"sentry"`
	Contentful    string `yaml:"contentful"`
	Amplience     string `yaml:"amplience"`
}

type TerraformConfig struct {
	AzureRemoteState *AzureTFState     `yaml:"azure_remote_state"`
	AwsRemoteState   *AWSTFState       `yaml:"aws_remote_state"`
	Providers        TerraformProvider `yaml:"providers"`
}
