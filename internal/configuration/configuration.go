package configuration

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/caarlos0/env/v11"
)

const (
	EnumDeploymentDev   = "dev"
	EnumDeploymentDebug = "debug"
	EnumDeploymentProd  = "prod"
)

//go:generate envdoc -output environments.md
type AppConfiguration struct {
	Deployment string `env:"DEPLOYMENT"`
	AppPort    string `env:"APP_PORT"`
	AppName    string `env:"APP_NAME"`
	isLoaded   bool
}

var appConfigurationInstance AppConfiguration

func GetAppConfiguration() (AppConfiguration, error) {
	if appConfigurationInstance.isLoaded {
		return appConfigurationInstance, nil
	}
	appConfigurationInstance := AppConfiguration{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&appConfigurationInstance, opts); err != nil {
		slog.Error(".env parsing error", "detail", err)
		return AppConfiguration{}, err
	}

	// validation
	allowedVals := []string{EnumDeploymentDev, EnumDeploymentDebug, EnumDeploymentProd}
	if !slices.Contains(allowedVals, appConfigurationInstance.Deployment) {
		return AppConfiguration{}, fmt.Errorf("DEPLOYMENT must be one of: %v", allowedVals)
	}

	appConfigurationInstance.isLoaded = true
	return appConfigurationInstance, nil
}
