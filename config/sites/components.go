package sites

import (
	config "github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/config/sites/components"
)

type Component struct {
	Name      string
	Variables map[string]any
	Secrets   map[string]any

	Definition config.Component
	Sentry     components.Sentry `yaml:"sentry"`
}
