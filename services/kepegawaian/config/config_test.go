package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	for k, v := range map[string]string{
		"NEXUS_KEPEGAWAIAN_SERVER_PORT":            "8000",
		"NEXUS_KEPEGAWAIAN_SERVER_AUTH_PUBLIC_KEY": "some-key",
		"NEXUS_KEPEGAWAIAN_DB_HOST":                "some-db-host",
		"NEXUS_KEPEGAWAIAN_DB_PORT":                "5432",
		"NEXUS_KEPEGAWAIAN_DB_SCHEMA":              "kepegawaian",
		"NEXUS_KEPEGAWAIAN_DB_NAME":                "some-db-name",
		"NEXUS_KEPEGAWAIAN_DB_USER":                "some-db-user",
		"NEXUS_KEPEGAWAIAN_DB_PASSWORD":            "some-db-password",
		"NEXUS_KEPEGAWAIAN_LOG_LEVEL":              "warn",
		"NEXUS_KEYCLOAK_HOST":                      "http://127.0.0.1:8080",
		"NEXUS_KEYCLOAK_REALM":                     "nexus",
		"NEXUS_KEYCLOAK_AUDIENCE":                  "portal",
	} {
		t.Setenv(k, v)
	}

	c, err := Load()

	require.NoError(t, err)
	assert.Equal(t, Config{
		Server: Server{
			Port:          8000,
			AuthPublicKey: "some-key",
		},
		DB: Database{
			Host:     "some-db-host",
			Port:     5432,
			Schema:   "kepegawaian",
			Name:     "some-db-name",
			User:     "some-db-user",
			Password: "some-db-password",
		},
		Keycloak: Keycloak{
			Host:     "http://127.0.0.1:8080",
			Realm:    "nexus",
			Audience: "portal",
		},
		LogLevel: slog.LevelWarn,
	}, c)
}
