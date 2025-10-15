package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server   Server   `envPrefix:"NEXUS_KEPEGAWAIAN_SERVER_"`
	DB       Database `envPrefix:"NEXUS_KEPEGAWAIAN_DB_"`
	Keycloak Keycloak `envPrefix:"NEXUS_KEYCLOAK_"`

	LogLevel slog.Level `env:"NEXUS_KEPEGAWAIAN_LOG_LEVEL"`
}

type Server struct {
	Port uint `env:"PORT"`
}

type Database struct {
	Host     string `env:"HOST"`
	Port     uint   `env:"PORT"`
	Schema   string `env:"SCHEMA"`
	Name     string `env:"NAME"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

type Keycloak struct {
	Host     string `env:"HOST"`
	Realm    string `env:"REALM"`
	Audience string `env:"AUDIENCE"`
}

func Load() (Config, error) {
	var c Config
	err := env.Parse(&c)
	return c, err
}
