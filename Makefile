DATABASE_URL := postgres://postgres:test@localhost:5432/postgres?sslmode=disable

define start_db
	@docker run --name testdb --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=test postgres:17.2
	@dbmate \
		--migrations-dir $(MIGRATIONS_DIR) \
		--schema-file $(SCHEMA_FILE) \
		--wait migrate \
		|| docker stop testdb
endef


dbgen-server: MIGRATIONS_DIR = ./internal/service/server/postgres/migrations
dbgen-server: SCHEMA_FILE = ./internal/service/server/postgres/schema.sql
dbgen-server: dbschema sqlc

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: dbschema
dbschema: export DATABASE_URL := $(DATABASE_URL)
dbschema:
	$(call start_db)
	@docker stop testdb

.PHONY: job-api
job-api: FILE = ./api/job/v1alpha1/oapi-merge.yaml
job-api: api-spec
	$(MAKE) oapi-codegen OAPI_GEN_FILE=./api/job/v1alpha1/oapi-gen.yaml  OAPI_SPEC=./api/job/v1alpha1/spec/openapi.bundled.yaml

.PHONY: server-api
server-api: FILE = ./api/server/v1alpha1/oapi-merge.yaml
server-api: api-spec
	$(MAKE) oapi-codegen OAPI_GEN_FILE=./api/server/v1alpha1/oapi-gen.yaml  OAPI_SPEC=./api/server/v1alpha1/spec/openapi.bundled.yaml

.PHONY: oapi-codegen
oapi-codegen:
	go tool oapi-codegen -config $(OAPI_GEN_FILE) $(OAPI_SPEC)

.PHONY: api-spec
api-spec:
	scripts/yaml-merge.sh $(FILE)
	npx --yes @redocly/cli bundle $(shell dirname $(FILE))/spec/openapi.bundled.yaml -o $(shell dirname $(FILE))/spec/openapi.bundled.yaml

.PHONY: dev
stop-dev:
	docker compose -f dev/docker-compose.yaml down

.PHONY: dev
start-dev:
	docker compose -f dev/docker-compose.yaml build
	docker compose -f dev/docker-compose.yaml up db -d
	dbmate --url postgres://app:secret@localhost:5432/server_service?sslmode=disable wait

	# SERVER SERVICE DB INIT
	docker compose -f dev/docker-compose.yaml exec -T db psql -U app -d postgres -c 'CREATE DATABASE server_service'
	dbmate --no-dump-schema --url postgres://app:secret@localhost:5432/server_service?sslmode=disable --migrations-dir ./internal/service/server/postgres/migrations migrate
	docker compose -f dev/docker-compose.yaml exec -T db psql -v ON_ERROR_STOP=1 -U app -d server_service < ./dev/server_service_init.sql

	SERVER_CLUSTERS="$$(jq -c . ./dev/clusterset_testdata.json)" docker compose -f dev/docker-compose.yaml up app