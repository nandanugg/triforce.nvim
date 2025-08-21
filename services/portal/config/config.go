package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server Server   `envPrefix:"NEXUS_PORTAL_SERVER_"`
	DB     Database `envPrefix:"NEXUS_PORTAL_DB_"`

	LogLevel slog.Level `env:"NEXUS_PORTAL_LOG_LEVEL" envDefault:"info"`
}

type Server struct {
	Port          uint   `env:"PORT"`
	AuthPublicKey string `env:"AUTH_PUBLIC_KEY"`
}

type Database struct {
	Host     string `env:"HOST"`
	Name     string `env:"NAME"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

func Load() (Config, error) {
	var c Config
	err := env.Parse(&c)
	return c, err
}
