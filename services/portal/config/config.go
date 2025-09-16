package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

const Service = "portal"

type Config struct {
	Server   Server   `envPrefix:"NEXUS_PORTAL_SERVER_"`
	DB       Database `envPrefix:"NEXUS_PORTAL_DB_"`
	Keycloak Keycloak `envPrefix:"NEXUS_KEYCLOAK_"`

	LogLevel slog.Level `env:"NEXUS_PORTAL_LOG_LEVEL" envDefault:"info"`
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
	Host         string `env:"HOST"`
	Realm        string `env:"REALM"`
	Audience     string `env:"AUDIENCE"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	PublicHost   string `env:"PUBLIC_HOST"`
	PrivateKey   string `env:"PRIVATE_KEY"`
	KID          string `env:"KID"`
}

func Load() (Config, error) {
	var c Config
	err := env.Parse(&c)
	return c, err
}
