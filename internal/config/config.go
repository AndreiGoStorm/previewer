package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Loading `yaml:"loading"`
		Log     `yaml:"logger"`
		Cache   `yaml:"cache"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host"`
		Port int    `env-required:"true" yaml:"port"`
	}

	Loading struct {
		Protocol string `env-required:"true" yaml:"protocol"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level"`
	}

	Cache struct {
		Capacity int `env-required:"true" yaml:"capacity"`
	}
)

func New(path string) *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("config new readconfig: %w", err))
	}

	return cfg
}
