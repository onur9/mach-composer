package components

type Component struct {
	Name         string
	Source       string
	Branch       string
	Integrations []string

	Version   string            `yaml:"version"`
	Endpoints map[string]string `yaml:"endpoints"`
	Azure     Azure             `yaml:"azure"`
}
