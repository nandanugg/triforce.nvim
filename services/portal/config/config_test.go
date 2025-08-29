package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	for k, v := range map[string]string{
		"NEXUS_PORTAL_SERVER_PORT":                "8000",
		"NEXUS_PORTAL_DB_HOST":                    "some-db-host",
		"NEXUS_PORTAL_DB_PORT":                    "5432",
		"NEXUS_PORTAL_DB_SCHEMA":                  "portal",
		"NEXUS_PORTAL_DB_NAME":                    "some-db-name",
		"NEXUS_PORTAL_DB_USER":                    "some-db-user",
		"NEXUS_PORTAL_DB_PASSWORD":                "some-db-password",
		"NEXUS_PORTAL_LOG_LEVEL":                  "warn",
		"NEXUS_KEYCLOAK_HOST":                     "http://127.0.0.1:8080",
		"NEXUS_KEYCLOAK_REALM":                    "nexus",
		"NEXUS_KEYCLOAK_AUDIENCE":                 "portal",
		"NEXUS_KEYCLOAK_CLIENT_ID":                "portal",
		"NEXUS_KEYCLOAK_CLIENT_SECRET":            "some-secret",
		"NEXUS_KEYCLOAK_REDIRECT_URI":             "http://portal.local/callback",
		"NEXUS_KEYCLOAK_POST_LOGOUT_REDIRECT_URI": "http://portal.local/",
	} {
		t.Setenv(k, v)
	}

	c, err := Load()

	require.NoError(t, err)
	assert.Equal(t, Config{
		Server: Server{
			Port: 8000,
		},
		DB: Database{
			Host:     "some-db-host",
			Port:     5432,
			Schema:   "portal",
			Name:     "some-db-name",
			User:     "some-db-user",
			Password: "some-db-password",
		},
		Keycloak: Keycloak{
			Host:                  "http://127.0.0.1:8080",
			Realm:                 "nexus",
			Audience:              "portal",
			ClientID:              "portal",
			ClientSecret:          "some-secret",
			RedirectURI:           "http://portal.local/callback",
			PostLogoutRedirectURI: "http://portal.local/",
		},
		LogLevel: slog.LevelWarn,
	}, c)
}
