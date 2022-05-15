package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port     string `envconfig:"PORT"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	PGDSN    string `envconfig:"PG_DSN"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
