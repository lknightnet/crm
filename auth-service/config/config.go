package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	Config struct {
		App        `yaml:"app"`
		Database   `yaml:"database"`
		JWT        `yaml:"jwt"`
		Http       `yaml:"http"`
		JWTService `yaml:"jwt_service"`
		Auth       `yaml:"auth"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	Database struct {
		URI string `yaml:"uri"`
	}

	JWT struct {
		SignKey string `yaml:"sign_key"`
	}

	Http struct {
		Port            string        `yaml:"port"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	JWTService struct {
		Api  string `yaml:"api"`
		Port string `yaml:"port"`
	}

	Auth struct {
		Signature string `yaml:"signature"`
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
