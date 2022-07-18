package global

import "github.com/labd/mach-composer/internal/config/global/azure"

type Azure struct {
	TenantID         string `yaml:"tenant_id"`
	SubscriptionID   string `yaml:"subscription_id"`
	Region           string
	Frontdoor        azure.Frontdoor              `yaml:"frontdoor"`
	ResourcesPrefix  string                       `yaml:"resources_prefix"`
	ServiceObjectIds map[string]string            `yaml:"service_object_ids"`
	ServicePlans     map[string]azure.ServicePlan `yaml:"service_plans"`
}
