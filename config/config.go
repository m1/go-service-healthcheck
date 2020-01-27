package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	// EnvironmentProduction ...
	EnvironmentProduction = "production"

	// EnvironmentStaging ...
	EnvironmentStaging = "staging"

	// EnvironmentDevelopment ...
	EnvironmentDevelopment = "development"
)

var (
	// GitCommit ...
	GitCommit string

	// BuildDate ...
	BuildDate string

	// Version ...
	Version string
)

// Config ...
type Config struct {
	Debug       bool   `env:"DEBUG" envDefault:"false"`
	Environment string `env:"ENV" envDefault:"production"`

	ServiceConfig
	APIConfig
	RunnerConfig
	DBConfig
}

// ServiceConfig ...
type ServiceConfig struct {
	Name      string `env:"SERVICE_NAME" envDefault:"name"`
	GitCommit string
	BuildDate string
	Version   string
}

// APIConfig ...
type APIConfig struct {
	Port   string `env:"API_PORT" envDefault:"8080"`
	Domain string `env:"API_DOMAIN" envDefault:"http://localhost"`
}

// DBConfig is the config struct for DBs ...
type DBConfig struct {
	File string `env:"DB_FILE" envDefault:"db.db"`
}

// DBConfig is the config struct for DBs ...
type RunnerConfig struct {
	ScrapeIntervalSeconds int `env:"RUNNER_SCRAPE_INTERVAL_SECONDS" envDefault:"15"`
	HTTPTimeoutSeconds    int `env:"RUNNER_HTTP_TIMEOUT_SECONDS" envDefault:"15"`
}

// LoadConfig ...
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.ServiceConfig)
	if err != nil {
		return nil, err
	}
	config.ServiceConfig.BuildDate = BuildDate
	config.ServiceConfig.GitCommit = GitCommit
	config.ServiceConfig.Version = Version

	err = env.Parse(&config.APIConfig)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.RunnerConfig)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.DBConfig)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
