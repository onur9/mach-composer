package commercetools

type Store struct {
	Key                  string
	Name                 map[string]string
	Managed              bool `default:"true"`
	Languages            []string
	DistributionChannels []string `yaml:"distribution_channels"`
	SupplyChannels       []string `yaml:"supply_channels"`
}
