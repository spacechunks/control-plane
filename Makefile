.PHONY: check-river-migrations
check-river-migrations:
	@./scripts/check-river-migrations.sh

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
