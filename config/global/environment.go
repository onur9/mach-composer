package global

type Environment string

var (
	EnvironmentProduction  Environment = "production"
	EnvironmentDevelopment Environment = "development"
	EnvironmentTest        Environment = "test"
)
