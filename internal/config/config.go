package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

const RunAddress = "http://localhost:8080"
const DatabaseUri = "postgresql://localhost:5432/postgres"
const AccrualSystemAddress = "http://localhost:8081"

type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	flagsConfig := new(Config)
	flag.StringVar(&flagsConfig.ServerAddress, "a", RunAddress, "server address")
	flag.StringVar(&flagsConfig.DataBaseURI, "d", DatabaseUri, "database uri")
	flag.StringVar(&flagsConfig.AccrualSystemAddress, "r", AccrualSystemAddress, "accrual system address")
	flag.Parse()

	if config.ServerAddress == "" {
		config.ServerAddress = flagsConfig.ServerAddress
	}

	if config.DataBaseURI == "" {
		config.DataBaseURI = flagsConfig.DataBaseURI
	}

	if config.AccrualSystemAddress == "" {
		config.AccrualSystemAddress = flagsConfig.AccrualSystemAddress
	}

	return &config, nil
}
