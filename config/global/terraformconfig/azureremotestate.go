package terraformconfig

// AzureRemoteState Azure storage account state backend configuration.
type AzureRemoteState struct {
	ResourceGroup  string `yaml:"resource_group"`
	StorageAccount string `yaml:"storage_account"`
	ContainerName  string `yaml:"container_name"`
	StateFolder    string `yaml:"state_folder"`
}
