OPENAPI_V1_URL := https://developers.pipedrive.com/docs/api/v1/openapi.yaml
OPENAPI_V2_URL := https://developers.pipedrive.com/docs/api/v1/openapi-v2.yaml
OAPI_CODEGEN_VERSION := v2.5.1
OAPI_CODEGEN := go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

.PHONY: test
test:
	go test ./...

.PHONY: update-specs
update-specs:
	mkdir -p openapi/upstream
	curl -L -o openapi/upstream/v1.yaml $(OPENAPI_V1_URL)
	curl -L -o openapi/upstream/v2.yaml $(OPENAPI_V2_URL)

.PHONY: derive-v1-legacy
derive-v1-legacy:
	go run ./cmd/derive-v1-legacy -v1 openapi/upstream/v1.yaml -v2 openapi/upstream/v2.yaml -out openapi/derived/v1-legacy.yaml -report openapi/derived/v1-legacy.report.json

.PHONY: openapi
openapi: update-specs derive-v1-legacy

.PHONY: generate
generate: derive-v1-legacy
	mkdir -p internal/gen/v1 internal/gen/v2
	$(OAPI_CODEGEN) -package v2 -generate types,client -o internal/gen/v2/openapi.gen.go openapi/upstream/v2.yaml
	$(OAPI_CODEGEN) -package v1 -generate types,client -o internal/gen/v1/openapi.gen.go openapi/derived/v1-legacy.yaml
	gofmt -w internal/gen/v1/openapi.gen.go internal/gen/v2/openapi.gen.go
