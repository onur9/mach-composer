package sites

import (
	"github.com/labd/mach-composer/internal/config/global"
	"github.com/labd/mach-composer/internal/config/sites/components"
)

type Site struct {
	Name       string
	Identifier string
	Endpoints  map[string]any `yaml:"endpoints"`
	// Endpoints    []Endpoint     `yaml:"_endpoints"` // TODO

	Components []Component `yaml:"components"`

	AWS   *AWS   `yaml:"aws,omitempty"`
	Azure *Azure `yaml:"azure,omitempty"`

	Commercetools Commercetools     `yaml:"commercetools"`
	Amplience     global.Amplience  `yaml:"amplience"`
	Sentry        components.Sentry `yaml:"sentry"`
}
