package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	LogLevel string `env:"LOGLEVEL" envDefault:"info"`
	Server   string `env:"SERVER" envDefault:":8080"`
}

func Load() Config {
	cfg := Config{}
	env.Parse(&cfg)
	return cfg
}
