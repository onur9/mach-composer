package components

// Sentry is the base sentry config.
type Sentry struct {
	DSN             string `yaml:"dsn"`
	RateLimitWindow int    `yaml:"rate_limit_window"`
	RateLimitCount  int    `yaml:"rate_limit_count"`
}
