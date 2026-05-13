package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HTTPPort    string `envconfig:"HTTP_PORT" default:"8080"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
}

func Load() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, err
	}
	return c, nil
}
