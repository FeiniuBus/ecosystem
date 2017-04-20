package hosting

import (
	"os"
)

const (
	envVarName  = "GO_ENVIRONMENT"
	production  = "production"
	staging     = "staging"
	development = "development"
)

var EnvironmentName string

func init() {
	EnvironmentName = os.Getenv(envVarName)
	if len(EnvironmentName) == 0 {
		EnvironmentName = production
	}
}

func IsProduction() bool {
	return EnvironmentName == production
}
func IsStaging() bool {
	return EnvironmentName == staging
}
func IsDevelopment() bool {
	return EnvironmentName == development
}
