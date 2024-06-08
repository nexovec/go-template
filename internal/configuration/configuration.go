package configuration

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"errors"
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

const (
	EnumDeploymentDev   = "dev"
	EnumDeploymentDebug = "debug"
	EnumDeploymentProd  = "prod"
)

var (
	ErrUserExists           = "ERROR: User already exists (SQLSTATE P0001)"
	ErrNoRows               = "no rows in result set"
	ErrServerRefusedToStart = errors.New("server hasn't even attempted to start")
)

//go:generate envdoc -output ../../environments.md
type DeploymentConfiguration struct {
	Deployment     string `env:"DEPLOYMENT"`
	AppName        string `env:"APP_NAME"`
	AppPort        string `env:"APP_PORT" envDefault:"80"`
	AppHost        string `env:"APP_HOST" envDefault:"0.0.0.0"`
	MainConfigFile string `env:"MAIN_CONFIG_FILE"`
	DbDSNStringPgx string `env:"DATABASE_DSN_STRING_PGX"`
	isLoaded       bool
}

var deploymentConfiguration DeploymentConfiguration

func GetDeploymentConfig() (DeploymentConfiguration, error) {
	if deploymentConfiguration.isLoaded {
		return deploymentConfiguration, nil
	}
	appConfigurationInstance := DeploymentConfiguration{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&appConfigurationInstance, opts); err != nil {
		slog.Error(".env parsing error", "detail", err)
		return DeploymentConfiguration{}, err
	}

	// validation
	allowedVals := []string{EnumDeploymentDev, EnumDeploymentDebug, EnumDeploymentProd}
	if !slices.Contains(allowedVals, appConfigurationInstance.Deployment) {
		return DeploymentConfiguration{}, fmt.Errorf("DEPLOYMENT must be one of: %v", allowedVals)
	}

	appConfigurationInstance.isLoaded = true
	return appConfigurationInstance, nil
}

type GeneralConfiguration struct {
	DbConnectionPoolSize     int           `yaml:"db_connection_pool_size"`
	DbStatementCacheCapacity int           `yaml:"db_statement_cache_capacity"`
	MaxServerStartRetries    int           `yaml:"max_server_start_retries"`
	GracefulShutdownTimeout  time.Duration `yaml:"graceful_shutdown_timeout"`
	DbHealthcheckInterval    time.Duration `yaml:"db_health_check_interval"`
	isLoaded                 bool
}

func (conf GeneralConfiguration) Validate() error {
	var errList []error

	if conf.DbConnectionPoolSize <= 0 {
		errList = append(errList, errors.New("DbConnectionPoolSize must be greater than 0"))
	}

	if conf.DbStatementCacheCapacity <= 0 {
		errList = append(errList, errors.New("DbStatementCacheCapacity must be greater than 0"))
	}

	if conf.MaxServerStartRetries <= 0 {
		errList = append(errList, errors.New("MaxServerStartRetries must be greater than 0"))
	}

	if conf.GracefulShutdownTimeout <= 0 {
		errList = append(errList, errors.New("GracefulShutdownTimeout must be greater than 0"))
	}

	return errors.Join(errList...)
}

var generalConfiguration GeneralConfiguration

func GetGeneralConfig() (GeneralConfiguration, error) {
	if generalConfiguration.isLoaded {
		return generalConfiguration, nil
	}

	// load from file
	conf, err := GetDeploymentConfig()
	if err != nil {
		return GeneralConfiguration{}, err
	}
	filename := conf.MainConfigFile

	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return GeneralConfiguration{}, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&generalConfiguration)
	if err != nil {
		return GeneralConfiguration{}, err
	}
	err = generalConfiguration.Validate()
	if err != nil {
		return GeneralConfiguration{}, err
	}
	generalConfiguration.isLoaded = true
	return generalConfiguration, nil
}
