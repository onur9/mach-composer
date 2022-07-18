package global

type Global struct {
	Environment     Environment     `yaml:"environment"`
	Cloud           Cloud           `yaml:"cloud"`
	Azure           Azure           `yaml:"azure"`
	TerraformConfig TerraformConfig `yaml:"terraform_config"`
	Amplience       Amplience       `yaml:"amplience"`
	Sentry          Sentry          `yaml:"sentry"`
	Contentful      Contentful      `yaml:"contentful"`
}
