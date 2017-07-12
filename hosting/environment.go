package hosting

import (
	"os"
	"strings"
)

const (
	envVarName  = "GO_ENVIRONMENT"
	production  = "production"
	staging     = "staging"
	development = "development"
)

// EnvironmentName is  系统环境变量
var EnvironmentName string

func init() {
	EnvironmentName = os.Getenv(envVarName)
	if len(EnvironmentName) == 0 {
		EnvironmentName = production
	}
}

// IsProduction returns 当前环境是否为生产环境
func IsProduction() bool {
	return strings.ToLower(EnvironmentName) == production
}

// IsStaging returns 当前环境是否为测试环境
func IsStaging() bool {
	return strings.ToLower(EnvironmentName) == staging
}

// IsDevelopment returns 当前环境是否为开发环境
func IsDevelopment() bool {
	return strings.ToLower(EnvironmentName) == development
}
