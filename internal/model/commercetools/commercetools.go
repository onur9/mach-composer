package commercetools

import "github.com/creasty/defaults"

type Settings struct {
	ProjectKey      string           `yaml:"project_key"`
	ClientID        string           `yaml:"client_id"`
	ClientSecret    string           `yaml:"client_secret"`
	Scopes          string           `yaml:"scopes"`
	TokenURL        string           `yaml:"token_url" default:"https://auth.europe-west1.gcp.commercetools.com"`
	APIURL          string           `yaml:"api_url" default:"https://api.europe-west1.gcp.commercetools.com"`
	ProjectSettings *ProjectSettings `yaml:"project_settings"`

	Frontend *FrontendSettings `yaml:"frontend"`

	Channels      []Channel
	Taxes         []Tax
	TaxCategories []TaxCategory `yaml:"tax_categories"`
	Stores        []Store
	Zones         []Zone
}

func (s *Settings) SetDefaults() {
	if defaults.CanUpdate(s.Frontend) {
		s.Frontend = &FrontendSettings{
			CreateCredentials: true,
		}
		s.Frontend.SetDefaults()
	}
}

func (s *Settings) ManagedStores() []Store {
	managed := make([]Store, 0)
	for _, store := range s.Stores {
		if store.Managed {
			managed = append(managed, store)
		}
	}
	return managed
}

type ProjectSettings struct {
	Languages  []string `yaml:"languages"`
	Currencies []string `yaml:"currencies"`
	Countries  []string `yaml:"countries"`

	MessagesEnabled bool `yaml:"messages_enabled" default:"true"`
}

type FrontendSettings struct {
	CreateCredentials bool     `yaml:"create_credentials" default:"false"`
	PermissionScopes  []string `yaml:"permission_scopes"`
}

func (s *FrontendSettings) SetDefaults() {
	if defaults.CanUpdate(s.PermissionScopes) {
		s.PermissionScopes = []string{
			"create_anonymous_token",
			"manage_my_profile",
			"manage_my_orders",
			"manage_my_shopping_lists",
			"manage_my_payments",
			"view_products",
			"view_project_settings",
		}
	}
}

type Store struct {
	Key                  string
	Name                 map[string]string
	Managed              bool `default:"true"`
	Languages            []string
	DistributionChannels []string `yaml:"distribution_channels"`
	SupplyChannels       []string `yaml:"supply_channels"`
}

type Channel struct {
	Key         string
	Roles       []string
	Name        map[string]string
	Description map[string]string
}

type Tax struct {
	Country         string
	Amount          float64
	Name            string
	IncludedInPrice bool `yaml:"included_in_price" default:"true"`
}

type TaxCategory struct {
	Key   string
	Name  string
	Rates []Tax
}

type ZoneLocation struct {
	Country string
	State   string
}

type Zone struct {
	Name        string
	Description string
	Locations   []ZoneLocation
}
