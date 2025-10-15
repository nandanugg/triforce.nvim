package main

import (
	"log/slog"
	"os"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/auth"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/dokumenpendukung"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/pemberitahuan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/resourcepermission"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/modules/role"
)

func main() {
	c, err := config.Load()
	exitIfError("Error loading application config.", err)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LogLevel})))

	db, err := db.New(c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Schema)
	exitIfError("Error connecting to database.", err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	exitIfError("Error creating new echo server.", err)

	keyfunc, err := api.NewAuthKeyfunc(c.Keycloak.Host, c.Keycloak.Realm, c.Keycloak.Audience)
	exitIfError("Error initializing auth keyfunc.", err)

	client := api.NewHTTPClient()
	privateKey, err := api.LoadRSAPrivateKey(c.Keycloak.PrivateKey)
	exitIfError("Error loading rsa private key.", err)

	dbRepository := repository.New(db)
	mwAuth := api.NewAuthMiddleware(api.NewAuthService(db), keyfunc)

	auth.RegisterRoutes(e, dbRepository, c.Keycloak, client, privateKey, keyfunc.Keyfunc)
	dokumenpendukung.RegisterRoutes(e, db, mwAuth)
	pemberitahuan.RegisterRoutes(e, db, mwAuth)
	resourcepermission.RegisterRoutes(e, dbRepository, mwAuth)
	role.RegisterRoutes(e, db, dbRepository, mwAuth)

	port := uint(c.Server.Port)
	exitIfError("Error parsing server port.", err)
	err = api.StartEchoServer(e, uint(port))
	exitIfError("Error starting server.", err)
}

func exitIfError(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
