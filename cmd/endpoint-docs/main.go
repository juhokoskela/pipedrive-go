package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type endpoint struct {
	Method      string
	Path        string
	Summary     string
	OperationID string
}

var methodOrder = map[string]int{
	"GET":     0,
	"POST":    1,
	"PUT":     2,
	"PATCH":   3,
	"DELETE":  4,
	"OPTIONS": 5,
	"HEAD":    6,
	"TRACE":   7,
}

func main() {
	v1Path := flag.String("v1", "openapi/derived/v1-legacy.yaml", "path to v1 legacy OpenAPI spec (yaml)")
	v2Path := flag.String("v2", "openapi/upstream/v2.yaml", "path to v2 OpenAPI spec (yaml)")
	outDir := flag.String("out-dir", "docs", "output directory for markdown tables")
	flag.Parse()

	if err := writeDoc(*v2Path, filepath.Join(*outDir, "endpoints-v2.md"), "Pipedrive API v2 endpoints"); err != nil {
		fatal(err)
	}
	if err := writeDoc(*v1Path, filepath.Join(*outDir, "endpoints-v1-legacy.md"), "Pipedrive API v1 legacy endpoints"); err != nil {
		fatal(err)
	}
}

func writeDoc(specPath, outPath, title string) error {
	endpoints, err := parseSpec(specPath)
	if err != nil {
		return err
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Path != endpoints[j].Path {
			return endpoints[i].Path < endpoints[j].Path
		}
		li := methodRank(endpoints[i].Method)
		lj := methodRank(endpoints[j].Method)
		if li != lj {
			return li < lj
		}
		return endpoints[i].Method < endpoints[j].Method
	})

	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	b.WriteString("Generated from `" + specPath + "` by `cmd/endpoint-docs`. Do not edit manually.\n\n")
	b.WriteString(fmt.Sprintf("Total operations: %d\n\n", len(endpoints)))
	b.WriteString("| Method | Path | Summary | Operation ID |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, e := range endpoints {
		b.WriteString("| ")
		b.WriteString(e.Method)
		b.WriteString(" | `")
		b.WriteString(e.Path)
		b.WriteString("` | ")
		b.WriteString(sanitizeCell(e.Summary))
		b.WriteString(" | ")
		b.WriteString(formatCodeCell(e.OperationID))
		b.WriteString(" |\n")
	}
	b.WriteString("\n")

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	return nil
}

func parseSpec(path string) ([]endpoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	paths := findMapValue(&root, "paths")
	if paths == nil || paths.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("paths not found in spec")
	}

	var out []endpoint
	for i := 0; i < len(paths.Content); i += 2 {
		pathKey := paths.Content[i]
		pathValue := paths.Content[i+1]
		if pathKey == nil || pathValue == nil || pathValue.Kind != yaml.MappingNode {
			continue
		}
		for j := 0; j < len(pathValue.Content); j += 2 {
			methodNode := pathValue.Content[j]
			opNode := pathValue.Content[j+1]
			if methodNode == nil || opNode == nil {
				continue
			}
			method := strings.ToUpper(methodNode.Value)
			if !isHTTPMethod(method) {
				continue
			}
			out = append(out, endpoint{
				Method:      method,
				Path:        pathKey.Value,
				Summary:     scalarValue(opNode, "summary"),
				OperationID: scalarValue(opNode, "operationId"),
			})
		}
	}

	return out, nil
}

func findMapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]
		if k != nil && k.Value == key {
			return v
		}
	}
	return nil
}

func scalarValue(node *yaml.Node, key string) string {
	if node == nil || node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]
		if k != nil && v != nil && k.Value == key && v.Kind == yaml.ScalarNode {
			return v.Value
		}
	}
	return ""
}

func isHTTPMethod(method string) bool {
	_, ok := methodOrder[method]
	return ok
}

func methodRank(method string) int {
	if rank, ok := methodOrder[method]; ok {
		return rank
	}
	return 100
}

func sanitizeCell(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	return value
}

func formatCodeCell(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	value = strings.ReplaceAll(value, "`", "\\`")
	return "`" + value + "`"
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
