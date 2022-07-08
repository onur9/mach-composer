package model

type AWSProvider struct {
	Name        string
	Region      string
	DefaultTags map[string]string `yaml:"default_tags"`
}

// AWSTFState is AWS S3 bucket state backend configuration.
type AWSTFState struct {
	Bucket    string `yaml:"bucket"`
	KeyPrefix string `yaml:"key_prefix"`
	Region    string `yaml:"region"`
	RoleARN   string `yaml:"role_arn"`
	LockTable string `yaml:"lock_table"`
	Encrypt   bool   `yaml:"encrypt" default:"true"`
}
