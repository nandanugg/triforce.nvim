package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/auth"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/dokumenpendukung"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/pemberitahuan"
)

func main() {
	c, err := config.Load()
	exitIfError("Error loading application config.", err)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LogLevel})))

	db, err := db.New(c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name)
	exitIfError("Error connecting to database.", err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	exitIfError("Error creating new echo server.", err)

	keyfunc, err := api.NewAuthKeyfunc(c.Keycloak.Host, c.Keycloak.Realm, c.Keycloak.Audience)
	exitIfError("Error initializing auth keyfunc.", err)

	mwAuth := api.NewAuthMiddleware(keyfunc)
	client := newHTTPClient()

	dokumenpendukung.RegisterRoutes(e, db, mwAuth)
	pemberitahuan.RegisterRoutes(e, db, mwAuth)
	auth.RegisterRoutes(e, c.Keycloak, client)

	err = api.StartEchoServer(e, c.Server.Port)
	exitIfError("Error starting server.", err)
}

func exitIfError(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

func newHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}
