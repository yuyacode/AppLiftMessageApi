package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Env        string `env:"ENV"         envDefault:"dev"`
	Port       int    `env:"PORT"        envDefault:"8080"`
	DBHost     string `env:"DB_HOST"     envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT"     envDefault:"3306"`
	DBCompany  string `env:"DB_COMPANY"  envDefault:"company"`
	DBStudent  string `env:"DB_STUDENT"  envDefault:"student"`
	DBCommon   string `env:"DB_COMMON"   envDefault:"common"`
	DBUserName string `env:"DB_USERNAME" envDefault:"user3"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password3"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
