package global

import "github.com/labd/mach-composer/internal/config/sites/components"

// Sentry global Sentry configuration.
type Sentry struct {
	components.Sentry
	AuthToken    string `yaml:"auth_token"`
	BaseURL      string `yaml:"base_url"`
	Project      string `yaml:"project"`
	Organization string `yaml:"organization"`
}
