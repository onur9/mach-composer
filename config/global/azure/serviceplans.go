package azure

type ServicePlan struct {
	Kind                   string
	Tier                   string
	Size                   string
	Capacity               int  // TODO
	DedicatedResourceGroup bool `yaml:"dedicated_resource_group"`
	PerSiteScaling         bool `yaml:"per_site_scaling"`
}
