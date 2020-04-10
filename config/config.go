package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	LogLevel string `env:"LOGLEVEL" envDefault:"info"`
	Server   string `env:"SERVER" envDefault:":8080"`
}

func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
