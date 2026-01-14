OPENAPI_V1_URL := https://developers.pipedrive.com/docs/api/v1/openapi.yaml
OPENAPI_V2_URL := https://developers.pipedrive.com/docs/api/v1/openapi-v2.yaml

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

