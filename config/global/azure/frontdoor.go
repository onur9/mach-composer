package azure

type SSLKeyVault struct {
	Name          string // TODO
	ResourceGroup string `yaml:"resource_group"`
	SecretName    string `yaml:"secret_name"`
}

type Frontdoor struct {
	DNSResourceGroup string      `yaml:"dns_resource_group"`
	SSLKeyVault      SSLKeyVault `yaml:"ssl_key_vault"`
	// SuppressChanges is an undocumented option to work around some tenacious issues
	// with using Frontdoor in the Azure Terraform provider
	SuppressChanges bool `yaml:"suppress_changes"`
}
