package config

import (
	"errors"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort         string `envconfig:"HTTP_PORT" default:"8080"`
	PostgresURL      string `envconfig:"POSTGRES_URL" required:"true"`
	OpenAIAPIKey     string `envconfig:"OPENAI_API_KEY" required:"true"`
	OpenAIModel      string `envconfig:"OPENAI_MODEL" default:"gpt-4o-mini"`
	OpenAIImageAPIKey string `envconfig:"GPT_IMAGE_1"`
	OpenAIImageModel string `envconfig:"OPENAI_IMAGE_MODEL" default:"gpt-image-1"`
}

func Load() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, err
	}
	// envconfig's `required` only checks presence, not emptiness;
	// an `FOO=` line counts as "set". Treat empty as missing.
	if strings.TrimSpace(c.PostgresURL) == "" {
		return Config{}, errors.New("config: POSTGRES_URL is required and must be non-empty")
	}
	if strings.TrimSpace(c.OpenAIAPIKey) == "" {
		return Config{}, errors.New("config: OPENAI_API_KEY is required and must be non-empty")
	}
	if strings.TrimSpace(c.OpenAIModel) == "" {
		c.OpenAIModel = "gpt-4o-mini"
	}
	if strings.TrimSpace(c.OpenAIImageModel) == "" {
		c.OpenAIImageModel = "gpt-image-1"
	}
	return c, nil
}
