package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	for k, v := range map[string]string{
		"NEXUS_PORTAL_SERVER_PORT":            "8000",
		"NEXUS_PORTAL_SERVER_AUTH_PUBLIC_KEY": "some-key",
		"NEXUS_PORTAL_DB_HOST":                "some-db-host",
		"NEXUS_PORTAL_DB_NAME":                "some-db-name",
		"NEXUS_PORTAL_DB_USER":                "some-db-user",
		"NEXUS_PORTAL_DB_PASSWORD":            "some-db-password",
		"NEXUS_PORTAL_LOG_LEVEL":              "warn",
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
			Name:     "some-db-name",
			User:     "some-db-user",
			Password: "some-db-password",
		},
		LogLevel: slog.LevelWarn,
	}, c)
}
