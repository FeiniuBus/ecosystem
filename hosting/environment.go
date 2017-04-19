package hosting

import (
	"os"
)

const envVarName = "GO_EnvironmentName"

var EnvironmentName string

func init() {
	EnvironmentName = os.Getenv(envVarName)
	if len(EnvironmentName) == 0 {
		EnvironmentName = "Production"
	}
}

func IsProduction() bool {
	return EnvironmentName == "Production"
}
func IsStaging() bool {
	return EnvironmentName == "Staging"
}
func IsDevelopment() bool {
	return EnvironmentName == "Development"
}
