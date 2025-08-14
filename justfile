set dotenv-load

fmt:
	golangci-lint fmt

lint-go:
	golangci-lint run

lint-openapi:
	spectral lint services/kepegawaian/docs/openapi.yaml

lint: lint-go lint-openapi

test:
	go test ./...

test-nocache:
	go test -count=1 ./...

run service:
	go run ./services/{{service}}

build service:
	go build -trimpath -ldflags="-s -w" ./services/{{service}}

db-migrate-new service migration_name:
	migrate create -ext sql -dir services/{{service}}/dbmigrations -seq {{migration_name}}

db-create-schema-kepegawaian:
	psql "host=$NEXUS_KEPEGAWAIAN_DB_HOST user=$NEXUS_KEPEGAWAIAN_DB_USER password=$NEXUS_KEPEGAWAIAN_DB_PASSWORD dbname=$NEXUS_KEPEGAWAIAN_DB_NAME" \
		-c "create schema kepegawaian"

db-migrate-up-kepegawaian:
	migrate \
	-path services/kepegawaian/dbmigrations \
	-database "pgx://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=kepegawaian" \
	up

db-migrate-down-kepegawaian:
	migrate \
	-path services/kepegawaian/dbmigrations \
	-database "pgx://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=kepegawaian" \
	down

db-doc-kepegawaian:
	tbls doc "postgres://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=kepegawian&sslmode=disable" services/kepegawaian/docs/db --rm-dist -t mermaid
	rm services/kepegawaian/docs/db/schema.json
