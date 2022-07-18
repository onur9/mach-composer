package commercetools

type ProjectSettings struct {
	Languages       []string `yaml:"languages"`
	Currencies      []string `yaml:"currencies"`
	Countries       []string `yaml:"countries"`
	MessagesEnabled bool     `yaml:"messages_enabled" default:"true"`
}
