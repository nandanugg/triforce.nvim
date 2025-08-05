set dotenv-load

fmt:
	golangci-lint fmt

lint-go:
	golangci-lint run

lint-openapi:
	spectral lint services/sampleservice1/docs/openapi.yaml

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

db-create-schema-sampleservice1:
	psql "host=$NEXUS_SAMPLESERVICE1_DB_HOST user=$NEXUS_SAMPLESERVICE1_DB_USER password=$NEXUS_SAMPLESERVICE1_DB_PASSWORD dbname=$NEXUS_SAMPLESERVICE1_DB_NAME" \
		-c "create schema $NEXUS_SAMPLESERVICE1_DB_SCHEMA"

db-migrate-up-sampleservice1:
	migrate \
	-path services/sampleservice1/dbmigrations \
	-database "pgx://$NEXUS_SAMPLESERVICE1_DB_USER:$NEXUS_SAMPLESERVICE1_DB_PASSWORD@$NEXUS_SAMPLESERVICE1_DB_HOST/$NEXUS_SAMPLESERVICE1_DB_NAME?search_path=$NEXUS_SAMPLESERVICE1_DB_SCHEMA" \
	up

db-migrate-down-sampleservice1:
	migrate \
	-path services/sampleservice1/dbmigrations \
	-database "pgx://$NEXUS_SAMPLESERVICE1_DB_USER:$NEXUS_SAMPLESERVICE1_DB_PASSWORD@$NEXUS_SAMPLESERVICE1_DB_HOST/$NEXUS_SAMPLESERVICE1_DB_NAME?search_path=$NEXUS_SAMPLESERVICE1_DB_SCHEMA" \
	down

db-doc-sampleservice1:
	tbls doc "postgres://$NEXUS_SAMPLESERVICE1_DB_USER:$NEXUS_SAMPLESERVICE1_DB_PASSWORD@$NEXUS_SAMPLESERVICE1_DB_HOST/$NEXUS_SAMPLESERVICE1_DB_NAME?search_path=$NEXUS_SAMPLESERVICE1_DB_SCHEMA&sslmode=disable" services/sampleservice1/docs/db --rm-dist -t mermaid
	rm services/sampleservice1/docs/db/schema.json
