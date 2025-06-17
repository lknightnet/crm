package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	Config struct {
		App         `yaml:"app"`
		Database    `yaml:"database"`
		Http        `yaml:"http"`
		UserService `yaml:"user_service"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	Database struct {
		URI string `yaml:"uri"`
	}

	Http struct {
		Port            string        `yaml:"port"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}
	UserService struct {
		Api  string `yaml:"api"`
		Port string `yaml:"port"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	err = os.Setenv("app.name", cfg.App.Name)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
