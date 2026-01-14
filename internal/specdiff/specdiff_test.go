package specdiff

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDeriveV1Legacy_RemovesOperationsPresentInV2ByOperationID(t *testing.T) {
	t.Parallel()

	v1 := `
openapi: 3.0.1
info: { title: Pipedrive API v1, version: 1.0.0 }
paths:
  /foo:
    get:
      operationId: getFoo
      responses: { "200": { description: ok } }
  /bar:
    get:
      operationId: getBar
      responses: { "200": { description: ok } }
`
	v2 := `
openapi: 3.0.1
info: { title: Pipedrive API v2, version: 2.0.0 }
paths:
  /foo:
    get:
      operationId: getFoo
      responses: { "200": { description: ok } }
`

	out, report, err := DeriveV1Legacy([]byte(v1), []byte(v2), nil)
	if err != nil {
		t.Fatalf("DeriveV1Legacy returned error: %v", err)
	}

	if report.V1Operations != 2 {
		t.Fatalf("expected V1Operations=2, got %d", report.V1Operations)
	}
	if report.V2Operations != 1 {
		t.Fatalf("expected V2Operations=1, got %d", report.V2Operations)
	}
	if report.RemovedOperations != 1 {
		t.Fatalf("expected RemovedOperations=1, got %d", report.RemovedOperations)
	}
	if report.LegacyOperations != 1 {
		t.Fatalf("expected LegacyOperations=1, got %d", report.LegacyOperations)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(out, &root); err != nil {
		t.Fatalf("failed to parse output YAML: %v", err)
	}
	paths := mustMappingValue(t, root.Content[0], "paths")
	if hasMappingKey(paths, "/foo") {
		t.Fatalf("expected /foo to be removed from derived spec")
	}
	if !hasMappingKey(paths, "/bar") {
		t.Fatalf("expected /bar to remain in derived spec")
	}

	if !strings.Contains(string(out), "getBar") {
		t.Fatalf("expected output to contain remaining operationId getBar")
	}
}

func TestDeriveV1Legacy_RemovesPathWhenNoMethodsRemain(t *testing.T) {
	t.Parallel()

	v1 := `
openapi: 3.0.1
paths:
  /only:
    get:
      operationId: onlyOp
      responses: { "200": { description: ok } }
`
	v2 := `
openapi: 3.0.1
paths:
  /only:
    get:
      operationId: onlyOp
      responses: { "200": { description: ok } }
`

	out, report, err := DeriveV1Legacy([]byte(v1), []byte(v2), nil)
	if err != nil {
		t.Fatalf("DeriveV1Legacy returned error: %v", err)
	}
	if report.LegacyOperations != 0 {
		t.Fatalf("expected LegacyOperations=0, got %d", report.LegacyOperations)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(out, &root); err != nil {
		t.Fatalf("failed to parse output YAML: %v", err)
	}
	paths := mustMappingValue(t, root.Content[0], "paths")
	if hasMappingKey(paths, "/only") {
		t.Fatalf("expected /only to be removed from derived spec")
	}
}

func mustMappingValue(t *testing.T, mapping *yaml.Node, key string) *yaml.Node {
	t.Helper()

	if mapping.Kind != yaml.MappingNode {
		t.Fatalf("expected mapping node, got kind=%d", mapping.Kind)
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	t.Fatalf("missing mapping key %q", key)
	return nil
}

func hasMappingKey(mapping *yaml.Node, key string) bool {
	if mapping.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return true
		}
	}
	return false
}

