package commercetools

type Frontend struct {
	CreateCredentials bool     `yaml:"create_credentials" default:"false"`
	PermissionScopes  []string `yaml:"permission_scopes"`
}
