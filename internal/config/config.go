package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type AuthKeyType string

const RunAddress = "http://localhost:8080"
const DatabaseURI = "postgresql://localhost:5432/postgres"
const AccrualSystemAddress = "http://localhost:8081"
const JwtSecret = "secret"

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
	AuthKey              AuthKeyType
}

func LoadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	flagsConfig := new(Config)
	flag.StringVar(&flagsConfig.RunAddress, "a", RunAddress, "server address")
	flag.StringVar(&flagsConfig.DataBaseURI, "d", DatabaseURI, "database uri")
	flag.StringVar(&flagsConfig.AccrualSystemAddress, "r", AccrualSystemAddress, "accrual system address")
	flag.StringVar(&flagsConfig.JWTSecret, "j", JwtSecret, "jwt secret")
	flag.Parse()

	if config.RunAddress == "" {
		config.RunAddress = flagsConfig.RunAddress
	}

	if config.DataBaseURI == "" {
		config.DataBaseURI = flagsConfig.DataBaseURI
	}

	if config.AccrualSystemAddress == "" {
		config.AccrualSystemAddress = flagsConfig.AccrualSystemAddress
	}

	if config.JWTSecret == "" {
		config.JWTSecret = flagsConfig.JWTSecret
	}

	return &config, nil
}
