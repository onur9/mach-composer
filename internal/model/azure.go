package model

import "strings"

var azureRegionDisplayMapShort = map[string]string{
	"eastasia":           "ea",
	"southeastasia":      "sea",
	"centralus":          "cus",
	"eastus":             "eus",
	"eastus2":            "eus2",
	"westus":             "wus",
	"northcentralus":     "ncus",
	"southcentralus":     "scus",
	"northeurope":        "ne",
	"westeurope":         "we",
	"japanwest":          "jw",
	"japaneast":          "je",
	"brazilsouth":        "bs",
	"australiaeast":      "ae",
	"australiasoutheast": "ase",
	"southindia":         "si",
	"centralindia":       "ci",
	"westindia":          "wi",
	"canadacentral":      "cc",
	"canadaeast":         "ce",
	"uksouth":            "us",
	"ukwest":             "uw",
	"westcentralus":      "wc",
	"westus2":            "wus2",
	"koreacentral":       "kc",
	"koreasouth":         "ks",
	"francecentral":      "fc",
	"francesouth":        "fs",
	"australiacentral":   "ac",
	"australiacentral2":  "ac2",
	"southafricanorth":   "san",
	"southafricawest":    "saw",
}

var azureRegionDisplayMapLon = map[string]string{
	"eastasia":           "East Asia",
	"southeastasia":      "Southeast Asia",
	"centralus":          "Central US",
	"eastus":             "East US",
	"eastus2":            "East US 2",
	"westus":             "West US",
	"northcentralus":     "North Central US",
	"southcentralus":     "South Central US",
	"northeurope":        "North Europe",
	"westeurope":         "West Europe",
	"japanwest":          "Japan West",
	"japaneast":          "Japan East",
	"brazilsouth":        "Brazil South",
	"australiaeast":      "Australia East",
	"australiasoutheast": "Australia Southeast",
	"southindia":         "South India",
	"centralindia":       "Central India",
	"westindia":          "West India",
	"canadacentral":      "Canada Central",
	"canadaeast":         "Canada East",
	"uksouth":            "UK South",
	"ukwest":             "UK West",
	"westcentralus":      "West Central US",
	"westus2":            "West US 2",
	"koreacentral":       "Korea Central",
	"koreasouth":         "Korea South",
	"francecentral":      "France Central",
	"francesouth":        "France South",
	"australiacentral":   "Australia Central",
	"australiacentral2":  "Australia Central 2",
	"southafricanorth":   "South Africa North",
	"southafricawest":    "South Africa West",
}

// AzureTFState Azure storage account state backend configuration.
type AzureTFState struct {
	ResourceGroup  string `yaml:"resource_group"`
	StorageAccount string `yaml:"storage_account"`
	ContainerName  string `yaml:"container_name"`
	StateFolder    string `yaml:"state_folder"`
}

type GlobalAzureConfig struct {
	TenantID         string `yaml:"tenant_id"`
	SubscriptionID   string `yaml:"subscription_id"`
	Region           string
	Frontdoor        *AzureFrontdoorSettings     `yaml:"frontdoor"`
	ResourcesPrefix  string                      `yaml:"resources_prefix"`
	ServiceObjectIds map[string]string           `yaml:"service_object_ids"`
	ServicePlans     map[string]AzureServicePlan `yaml:"service_plans"`
}

type AzureAlertGroup struct {
	Name        string
	AlertEmails []string `yaml:"alert_emals"`
	WebhookURL  string   `yaml:"webhook_url"`
	LogicApp    string   `yaml:"logic_app"`
}

func (a *AzureAlertGroup) LogicAppName() string {
	if a.LogicApp != "" {
		parts := strings.Split(a.LogicApp, ".")
		return parts[1]
	}
	return ""
}

func (a *AzureAlertGroup) LogicAppResourceGroup() string {
	if a.LogicApp != "" {
		parts := strings.Split(a.LogicApp, ".")
		return parts[0]
	}
	return ""
}

type AzureFrontdoorSSLConfig struct {
	Name          string
	ResourceGroup string `yaml:"resource_group"`
	SecretName    string `yaml:"secret_name"`
}

type AzureServicePlan struct {
	Kind                   string
	Tier                   string
	Size                   string
	Capacity               int
	DedicatedResourceGroup bool `yaml:"dedicated_resource_group"`
	PerSiteScaling         bool `yaml:"per_site_scaling"`
}

type AzureFrontdoorSettings struct {
	DNSResourceGroup string                   `yaml:"dns_resource_group"`
	SslKeyVault      *AzureFrontdoorSSLConfig `yaml:"ssl_key_vault"`
	// undocumented option to work around some tenacious issues
	// with using Frontdoor in the Azure Terraform provider
	SuppressChanges bool `yaml:"suppress_changes"`
}
