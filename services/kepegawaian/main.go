package main

import (
	"log/slog"
	"os"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/datapribadi"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/golongan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jabatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenisdiklat"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenishukuman"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenisjabatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jeniskenaikanpangkat"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenispegawai"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenispenghargaan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/keluarga"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/pegawai"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatasesmenninebox"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayathukumandisiplin"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatjabatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatkenaikangajiberkala"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatkepangkatan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatkinerja"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpelatihanfungsional"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpelatihansiasn"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpelatihanteknis"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpendidikan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpenghargaan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatpenugasan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/riwayatsertifikasi"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/statuspernikahan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/suratkeputusan"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/tingkatpendidikan"
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

	datapribadi.RegisterRoutes(e, dbRepository, mwAuth)
	golongan.RegisterRoutes(e, dbRepository, mwAuth)
	jabatan.RegisterRoutes(e, dbRepository, mwAuth)
	jenisdiklat.RegisterRoutes(e, dbRepository, mwAuth)
	jenishukuman.RegisterRoutes(e, dbRepository, mwAuth)
	jenisjabatan.RegisterRoutes(e, dbRepository, mwAuth)
	jeniskenaikanpangkat.RegisterRoutes(e, dbRepository, mwAuth)
	jenispegawai.RegisterRoutes(e, dbRepository, mwAuth)
	jenispenghargaan.RegisterRoutes(e, dbRepository, mwAuth)
	keluarga.RegisterRoutes(e, dbRepository, mwAuth)
	pegawai.RegisterRoutes(e, db, mwAuth)
	riwayatasesmenninebox.RegisterRoutes(e, dbRepository, mwAuth)
	riwayathukumandisiplin.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatjabatan.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatkenaikangajiberkala.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatkepangkatan.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatkinerja.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpelatihanfungsional.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpelatihansiasn.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpelatihanteknis.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpendidikan.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpenghargaan.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatpenugasan.RegisterRoutes(e, dbRepository, mwAuth)
	riwayatsertifikasi.RegisterRoutes(e, dbRepository, mwAuth)
	suratkeputusan.RegisterRoutes(e, dbRepository, mwAuth)
	statuspernikahan.RegisterRoutes(e, dbRepository, mwAuth)
	tingkatpendidikan.RegisterRoutes(e, dbRepository, mwAuth)
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
