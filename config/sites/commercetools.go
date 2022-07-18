package sites

import "github.com/labd/mach-composer/internal/config/sites/commercetools"

type Commercetools struct {
	ProjectKey      string                        `yaml:"project_key"`
	ClientID        string                        `yaml:"client_id"`
	ClientSecret    string                        `yaml:"client_secret"`
	Scopes          string                        `yaml:"scopes"`
	TokenURL        string                        `yaml:"token_url" default:"https://auth.europe-west1.gcp.commercetools.com"`
	APIURL          string                        `yaml:"api_url" default:"https://api.europe-west1.gcp.commercetools.com"`
	ProjectSettings commercetools.ProjectSettings `yaml:"project_settings"`
	Frontend        commercetools.Frontend        `yaml:"frontend"`
	TaxCategories   []commercetools.TaxCategory   `yaml:"tax_categories"`
	Channels        []commercetools.Channel
	Taxes           []commercetools.Tax
	Stores          []commercetools.Store
	Zones           []commercetools.Zone
}
