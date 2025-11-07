// Package config loads the environment variables to the application
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the data structure for the env fields that should be loaded to the application
type Config struct {
	BaseAPIURL    string `envconfig:"BASE_API_URL" required:"true"`
	AuthSecretKey string `envconfig:"AUTH_SECRET_KEY" required:"true"`
	Port          string `envconfig:"PORT" default:"8089"`
}

// Load loads environment variables into Config and validates them.
func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}

	return &cfg
}
