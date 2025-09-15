package main

import (
	"log/slog"
	"os"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/asesmenninebox"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/datapribadi"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/golongan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/hukumandisiplin"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jabatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenisjabatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jeniskenaikanpangkat"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/keluarga"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/kenaikangajiberkala"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/kepangkatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/kinerja"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pegawai"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pekerjaan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pelatihanfungsional"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pelatihanstruktural"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pelatihanteknis"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pendidikanformal"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/penghargaan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/penugasan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/sertifikasi"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/unitkerja"
)

func main() {
	c, err := config.Load()
	exitIfError("Error loading application config.", err)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LogLevel})))

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	exitIfError("Error creating new echo server.", err)

	keyfunc, err := api.NewAuthKeyfunc(c.Keycloak.Host, c.Keycloak.Realm, c.Keycloak.Audience)
	exitIfError("Error initializing auth keyfunc.", err)

	db, err := db.New(c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Schema)
	exitIfError("Error connecting to database with pgx.", err)

	dbRepository := repository.New(db)

	mwAuth := api.NewAuthMiddleware(config.Service, keyfunc)

	asesmenninebox.RegisterRoutes(e, db, mwAuth)
	golongan.RegisterRoutes(e, dbRepository, mwAuth)
	datapribadi.RegisterRoutes(e, db, mwAuth)
	hukumandisiplin.RegisterRoutes(e, dbRepository, mwAuth)
	jabatan.RegisterRoutes(e, dbRepository, mwAuth)
	jenisjabatan.RegisterRoutes(e, dbRepository, mwAuth)
	jeniskenaikanpangkat.RegisterRoutes(e, dbRepository, mwAuth)
	keluarga.RegisterRoutes(e, dbRepository, mwAuth)
	kenaikangajiberkala.RegisterRoutes(e, db, mwAuth)
	kepangkatan.RegisterRoutes(e, dbRepository, mwAuth)
	kinerja.RegisterRoutes(e, db, mwAuth)
	pegawai.RegisterRoutes(e, db, mwAuth)
	pekerjaan.RegisterRoutes(e, db, mwAuth)
	pelatihanfungsional.RegisterRoutes(e, db, mwAuth)
	pelatihanstruktural.RegisterRoutes(e, dbRepository, mwAuth)
	pelatihanteknis.RegisterRoutes(e, dbRepository, mwAuth)
	pendidikanformal.RegisterRoutes(e, dbRepository, mwAuth)
	penghargaan.RegisterRoutes(e, db, mwAuth)
	penugasan.RegisterRoutes(e, db, mwAuth)
	sertifikasi.RegisterRoutes(e, dbRepository, mwAuth)
	unitkerja.RegisterRoutes(e, dbRepository, mwAuth)
	err = api.StartEchoServer(e, c.Server.Port)
	exitIfError("Error starting server.", err)
}

func exitIfError(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
