package model

import "strings"

type ComponentAzureConfig struct {
	ServicePlan string `yaml:"service_plan"`
	ShortName   string `yaml:"short_name"`
}

type Component struct {
	Name         string
	Source       string
	Version      string `yaml:"version"`
	Branch       string
	Integrations []string
	Endpoints    map[string]string `yaml:"endpoints"`

	Azure *ComponentAzureConfig `yaml:"azure"`
}

// UseVersionReference indicates if the module should be referenced with the
// version.
// This will be mainly used for development purposes when referring to a local
// directory; versioning is not possible, but we should still be able to define
// a version in our component for the actual function deployment itself.
func (c *Component) UseVersionReference() bool {
	return strings.HasPrefix(c.Source, "git")
}
