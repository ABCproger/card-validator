package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Env string `yaml:"env" env:"APP_ENV" env-default:"development"`
	} `yaml:"app"`
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	} `yaml:"http"`
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"log"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("config/config.yml", cfg); err != nil {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("config: %w", err)
		}
	}
	return cfg, nil
}
