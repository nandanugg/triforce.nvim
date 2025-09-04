set dotenv-load

fmt:
	golangci-lint fmt

lint-go:
	golangci-lint run

lint-openapi:
	spectral lint \
		services/kepegawaian/docs/openapi.yaml \
		services/portal/docs/openapi.yaml

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
	psql "host=$NEXUS_KEPEGAWAIAN_DB_HOST port=$NEXUS_KEPEGAWAIAN_DB_PORT user=$NEXUS_KEPEGAWAIAN_DB_USER password=$NEXUS_KEPEGAWAIAN_DB_PASSWORD dbname=$NEXUS_KEPEGAWAIAN_DB_NAME" \
		-c "create schema $NEXUS_KEPEGAWAIAN_DB_SCHEMA"

db-migrate-up-kepegawaian:
	migrate \
	-path services/kepegawaian/db/migrations \
	-database "pgx://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST:$NEXUS_KEPEGAWAIAN_DB_PORT/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=$NEXUS_KEPEGAWAIAN_DB_SCHEMA" \
	up

db-migrate-down-kepegawaian:
	migrate \
	-path services/kepegawaian/db/migrations \
	-database "pgx://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST:$NEXUS_KEPEGAWAIAN_DB_PORT/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=$NEXUS_KEPEGAWAIAN_DB_SCHEMA" \
	down

db-doc-kepegawaian:
	tbls doc "postgres://$NEXUS_KEPEGAWAIAN_DB_USER:$NEXUS_KEPEGAWAIAN_DB_PASSWORD@$NEXUS_KEPEGAWAIAN_DB_HOST:$NEXUS_KEPEGAWAIAN_DB_PORT/$NEXUS_KEPEGAWAIAN_DB_NAME?search_path=$NEXUS_KEPEGAWAIAN_DB_SCHEMA&sslmode=disable" services/kepegawaian/docs/db --rm-dist -t mermaid
	rm services/kepegawaian/docs/db/schema.json

db-create-schema-portal:
	psql "host=$NEXUS_PORTAL_DB_HOST port=$NEXUS_PORTAL_DB_PORT user=$NEXUS_PORTAL_DB_USER password=$NEXUS_PORTAL_DB_PASSWORD dbname=$NEXUS_PORTAL_DB_NAME" \
		-c "create schema $NEXUS_PORTAL_DB_SCHEMA"

db-migrate-up-portal:
	migrate \
	-path services/portal/dbmigrations \
	-database "pgx://$NEXUS_PORTAL_DB_USER:$NEXUS_PORTAL_DB_PASSWORD@$NEXUS_PORTAL_DB_HOST:$NEXUS_PORTAL_DB_PORT/$NEXUS_PORTAL_DB_NAME?search_path=$NEXUS_PORTAL_DB_SCHEMA" \
	up

db-migrate-down-portal:
	migrate \
	-path services/portal/dbmigrations \
	-database "pgx://$NEXUS_PORTAL_DB_USER:$NEXUS_PORTAL_DB_PASSWORD@$NEXUS_PORTAL_DB_HOST:$NEXUS_PORTAL_DB_PORT/$NEXUS_PORTAL_DB_NAME?search_path=$NEXUS_PORTAL_DB_SCHEMA" \
	down

db-doc-portal:
	tbls doc "postgres://$NEXUS_PORTAL_DB_USER:$NEXUS_PORTAL_DB_PASSWORD@$NEXUS_PORTAL_DB_HOST:$NEXUS_PORTAL_DB_PORT/$NEXUS_PORTAL_DB_NAME?search_path=$NEXUS_PORTAL_DB_SCHEMA&sslmode=disable" services/portal/docs/db --rm-dist -t mermaid
	rm services/portal/docs/db/schema.json
