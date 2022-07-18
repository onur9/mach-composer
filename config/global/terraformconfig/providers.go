package terraformconfig

type Providers struct {
	AWS           string `yaml:"aws"`
	Azure         string `yaml:"azure"`
	Commercetools string `yaml:"commercetools"`
	Sentry        string `yaml:"sentry"`
	Contentful    string `yaml:"contentful"`
	Amplience     string `yaml:"amplience"`
}
