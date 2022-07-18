package terraformconfig

// AWSRemoteState is AWS S3 bucket state backend configuration.
type AWSRemoteState struct {
	Bucket    string `yaml:"bucket"`
	KeyPrefix string `yaml:"key_prefix"`
	Region    string `yaml:"region"`
	RoleARN   string `yaml:"role_arn"`
	LockTable string `yaml:"lock_table"`
	Encrypt   bool   `yaml:"encrypt" default:"true"`
}
