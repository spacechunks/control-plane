.PHONY: api-spec
api-spec:
	scripts/yaml-merge.sh $(FILE)
	npx --yes @redocly/cli bundle $(shell dirname $(FILE))/openapi.bundled.yaml -o $(shell dirname $(FILE))/openapi.bundled.yaml
