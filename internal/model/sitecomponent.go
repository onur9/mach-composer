package model

import "log"

type SiteComponent struct {
	Name      string
	Variables map[string]any
	Secrets   map[string]any

	Definition *Component
	Sentry     *SentryConfig `yaml:"sentry"`
}

func (sc SiteComponent) HasCloudIntegration() bool {
	if sc.Definition == nil {
		log.Fatalf("Component %s was not resolved properly (missing definition)", sc.Name)
	}
	for _, i := range sc.Definition.Integrations {
		if i == "aws" || i == "azure" {
			return true
		}
	}
	return false
}
