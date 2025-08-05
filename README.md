# nexus-be

Nexus Backend Server

## Development

* Bahasa pemrograman: [Go](https://go.dev/learn/)
* Database: [PostgreSQL](https://www.postgresql.org/)

Tools tambahan (opsional):

* [Docker + Docker Compsoe](https://docs.docker.com/compose/gettingstarted/) -
  untuk menyediakan PostgreSQL di environment local.
* [golangci-lint](https://golangci-lint.run/welcome/install/) - linter &
  formatter untuk kode Go.
* [spectral](https://docs.stoplight.io/docs/spectral/b8391e051b7d8-installation) -
  linter untuk dokumen OpenAPI.
* [migrate](https://github.com/golang-migrate/migrate) - untuk menjalankan
  perubahan schema database.
* [tbls](https://github.com/k1LoW/tbls) - untuk membuat dokumen ERD.
* [just](https://github.com/casey/just) - untuk menjalankan shell command yang
  berkaitan dengan pengembangan.

### Quick Start

1. Buat file .env di root project directory dengan menggunakan file .env.sample
   sebagai referensi. Edit sesuai dengan konfigurasi di local.
1. Jalankan test dengan perintah `just test` (atau `source .env; go test ./...`).

### Struktur Direktori

```
/
├─services
│ └─X
│   ├─config
│   ├─dbmigrations
│   ├─docs
│   └─modules
│     └─Y
└─lib
```

**services/X**

Seluruh kode untuk service X. Entrypoint service (`func main()`) ada di sini.

**services/X/config**

Definisi konfigurasi untuk service X.

**services/X/dbmigrations**

Daftar database migration (DDL) untuk service X.

**services/X/docs**

File-file dokumentasi untuk service X (skema API, ERD, dll).

**services/X/modules/Y**

Seluruh kode untuk module Y pada service X.

**lib**

Fungsi-fungsi tambahan yang bisa digunakan di service manapun.
