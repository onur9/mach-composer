package aws

type ExtraProvider struct {
	Name        string
	Region      string
	DefaultTags map[string]string `yaml:"default_tags"`
}
