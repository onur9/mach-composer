package model

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/labd/mach-composer/internal/model/commercetools"
	"github.com/labd/mach-composer/internal/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type SiteAWS struct {
	AccountID string `yaml:"account_id"`
	Region    string `yaml:"region"`

	DeployRoleName string            `yaml:"deploy_role_name"`
	ExtraProviders []AWSProvider     `yaml:"extra_providers"`
	DefaultTags    map[string]string `yaml:"default_tags"`
}

// Site contains all configuration needed for a site.
type Site struct {
	Name         string
	Identifier   string
	RawEndpoints map[string]any `yaml:"endpoints"`
	Endpoints    []Endpoint     `yaml:"_endpoints"`

	Components []SiteComponent `yaml:"components"`

	AWS   *SiteAWS           `yaml:"aws,omitempty"`
	Azure *SiteAzureSettings `yaml:"azure,omitempty"`

	Commercetools *commercetools.Settings `yaml:"commercetools"`
	Amplience     *AmplienceConfig        `yaml:"amplience"`
	Sentry        *SentryConfig           `yaml:"sentry"`
}

func (s *Site) ResolveEndpoints() {
	for k, rv := range s.RawEndpoints {
		switch v := rv.(type) {
		case string:
			ep := Endpoint{
				Key: k,
				URL: v,
			}
			if err := defaults.Set(&ep); err != nil {
				panic(err)
			}
			s.Endpoints = append(s.Endpoints, ep)

		case map[string]any:
			// Do an extra serialize/deserialize step here. Simplest solution
			// for now.

			body, err := yaml.Marshal(v)
			if err != nil {
				panic(err)
			}

			ep := Endpoint{
				Key: k,
			}
			err = yaml.Unmarshal(body, &ep)
			if err != nil {
				panic(err)
			}

			if err := defaults.Set(&ep); err != nil {
				panic(err)
			}

			s.Endpoints = append(s.Endpoints, ep)
		default:
			panic("Missing")
		}
	}

	// Check if we need to add a default endpoint
	components := s.EndpointComponents()
	keys := make([]string, 0, len(s.Endpoints))
	for _, e := range s.Endpoints {
		keys = append(keys, e.Key)
	}

	// If one of the components has a 'default' endpoint defined,
	// we'll include it to our site endpoints.
	// A 'default' endpoint is one without a custom domain, so no further
	// Route53 or DNS zone settings required.
	componentKeys := []string{}
	for k := range components {
		componentKeys = append(componentKeys, k)
	}
	if stringContains(componentKeys, "default") && stringContains(keys, "default") {
		fmt.Println(
			"WARNING: 'default' endpoint used but not defined in the site endpoints.\n" +
				"MACH will create a default endpoint without any custom domain attached to it.\n" +
				"More info: https://docs.machcomposer.io/reference/syntax/sites.html#endpoints",
		)
		s.Endpoints = append(s.Endpoints, Endpoint{
			URL: "",
			Key: "default",
		})
	}
}

func (s *Site) EndpointComponents() map[string][]SiteComponent {
	// Check if we need to add a default endpoint
	endpoints := make(map[string][]SiteComponent)
	for _, c := range s.Components {
		for _, value := range c.Definition.Endpoints {
			endpoints[value] = append(endpoints[value], c)
		}
	}
	return endpoints

}

// UsedEndpoints returns only the endpoints that are actually used by the components.
func (s *Site) UsedEndpoints() []Endpoint {
	result := []Endpoint{}
	for _, ep := range s.Endpoints {
		if len(ep.Components) > 0 {
			result = append(result, ep)
		}
	}
	return result
}

// DNSZones returns the DNS zones of used endpoints.
func (s *Site) DNSZones() []string {
	result := []string{}
	endpoints := s.UsedEndpoints()
	for i := range endpoints {
		result = append(result, endpoints[i].Zone)
	}
	return utils.UniqueSlice(result)
}

// HasCDNEndpoint checks if there is an endpoint with a CDN enabled.
func (s *Site) HasCDNEndpoint() bool {
	endpoints := s.UsedEndpoints()
	for _, ep := range endpoints {
		if ep.AWS != nil && ep.AWS.EnableCDN {
			return true
		}
	}
	return false
}

// SiteAzureSettings Site-specific Azure settings
type SiteAzureSettings struct {
	Frontdoor  *AzureFrontdoorSettings `yaml:"frontdoor"`
	AlertGroup *AzureAlertGroup        `yaml:"alert_group"`

	// Can overwrite values from AzureConfig
	ResourceGroup  string
	TenantID       string `yaml:"tenant_id"`
	SubscriptionID string `yaml:"subscription_id"`

	Region           string
	ServiceObjectIds map[string]string           `yaml:"service_object_ids"`
	ServicePlans     map[string]AzureServicePlan `yaml:"service_plans"`
}

func (a *SiteAzureSettings) Merge(c *GlobalAzureConfig) {
	if a.Frontdoor == nil {
		a.Frontdoor = c.Frontdoor
	}
	if a.TenantID == "" {
		a.TenantID = c.TenantID
	}
	if a.SubscriptionID == "" {
		a.SubscriptionID = c.SubscriptionID
	}
	if a.Region == "" {
		a.Region = c.Region
	}

	if len(a.ServiceObjectIds) == 0 {
		a.ServiceObjectIds = c.ServiceObjectIds
	}

	for k, v := range c.ServicePlans {
		a.ServicePlans[k] = v
	}
}

func (a *SiteAzureSettings) ShortRegionName() string {
	if val, ok := azureRegionDisplayMapShort[a.Region]; ok {
		return val
	}
	logrus.Fatalf("No short name for region %s", a.Region)
	return ""
}

func (a *SiteAzureSettings) LongRegionName() string {
	if val, ok := azureRegionDisplayMapLon[a.Region]; ok {
		return val
	}
	logrus.Fatalf("No long name for region %s", a.Region)
	return ""
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
