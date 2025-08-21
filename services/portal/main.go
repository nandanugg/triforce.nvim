package main

import (
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/pemberitahuan"
)

func main() {
	c, err := config.Load()
	exitIfError("Error loading application config.", err)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LogLevel})))

	db, err := db.New(c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name)
	exitIfError("Error connecting to database.", err)

	jwtPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(c.Server.AuthPublicKey))
	exitIfError("Error parsing auth public key.", err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	exitIfError("Error creating new echo server.", err)

	mwAuth := api.NewAuthMiddleware(jwtPublicKey)

	pemberitahuan.RegisterRoutes(e, db, mwAuth)

	err = api.StartEchoServer(e, c.Server.Port)
	exitIfError("Error starting server.", err)
}

func exitIfError(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
