package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	for k, v := range map[string]string{
		"NEXUS_SAMPLESERVICE1_SERVER_PORT":            "8000",
		"NEXUS_SAMPLESERVICE1_SERVER_AUTH_PUBLIC_KEY": "some-key",
		"NEXUS_SAMPLESERVICE1_DB_HOST":                "some-db-host",
		"NEXUS_SAMPLESERVICE1_DB_NAME":                "some-db-name",
		"NEXUS_SAMPLESERVICE1_DB_USER":                "some-db-user",
		"NEXUS_SAMPLESERVICE1_DB_PASSWORD":            "some-db-password",
		"NEXUS_SAMPLESERVICE1_DB_SCHEMA":              "some-db-schema",
		"NEXUS_SAMPLESERVICE1_LOG_LEVEL":              "warn",
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
			Schema:   "some-db-schema",
		},
		LogLevel: slog.LevelWarn,
	}, c)
}
