package sites

import "github.com/labd/mach-composer/internal/config/sites/aws"

type AWS struct {
	AccountID      string              `yaml:"account_id"`
	Region         string              `yaml:"region"`
	DeployRoleName string              `yaml:"deploy_role_name"`
	ExtraProviders []aws.ExtraProvider `yaml:"extra_providers"`
	DefaultTags    map[string]string   `yaml:"default_tags"`
}
