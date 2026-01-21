package specdiff

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Report struct {
	V1Operations      int
	V2Operations      int
	RemovedOperations int
	LegacyOperations  int
}

type Options struct{}

func DeriveV1Legacy(v1, v2 []byte, _ *Options) ([]byte, Report, error) {
	var v1Doc yaml.Node
	if err := yaml.Unmarshal(v1, &v1Doc); err != nil {
		return nil, Report{}, fmt.Errorf("parse v1: %w", err)
	}
	var v2Doc yaml.Node
	if err := yaml.Unmarshal(v2, &v2Doc); err != nil {
		return nil, Report{}, fmt.Errorf("parse v2: %w", err)
	}

	v2IDs, v2Count, err := operationIDs(&v2Doc)
	if err != nil {
		return nil, Report{}, fmt.Errorf("collect v2 operations: %w", err)
	}
	v1Count, err := operationCount(&v1Doc)
	if err != nil {
		return nil, Report{}, fmt.Errorf("collect v1 operations: %w", err)
	}

	removed, remaining, err := removeOperations(&v1Doc, v2IDs)
	if err != nil {
		return nil, Report{}, fmt.Errorf("derive v1-legacy: %w", err)
	}

	out, err := yaml.Marshal(&v1Doc)
	if err != nil {
		return nil, Report{}, fmt.Errorf("encode derived yaml: %w", err)
	}
	out = bytes.TrimSpace(out)
	out = append(out, '\n')

	return out, Report{
		V1Operations:      v1Count,
		V2Operations:      v2Count,
		RemovedOperations: removed,
		LegacyOperations:  remaining,
	}, nil
}

var httpMethods = map[string]struct{}{
	"get":     {},
	"post":    {},
	"put":     {},
	"patch":   {},
	"delete":  {},
	"head":    {},
	"options": {},
	"trace":   {},
}

func operationIDs(doc *yaml.Node) (map[string]struct{}, int, error) {
	paths, err := pathsNode(doc)
	if err != nil {
		return nil, 0, err
	}

	ids := make(map[string]struct{})
	count := 0

	for i := 0; i+1 < len(paths.Content); i += 2 {
		pathItem := paths.Content[i+1]
		if pathItem.Kind != yaml.MappingNode {
			continue
		}
		for j := 0; j+1 < len(pathItem.Content); j += 2 {
			methodNode := pathItem.Content[j]
			if methodNode.Kind != yaml.ScalarNode {
				continue
			}
			if _, ok := httpMethods[methodNode.Value]; !ok {
				continue
			}
			opNode := pathItem.Content[j+1]
			opID := scalarMappingValue(opNode, "operationId")
			if opID == "" {
				continue
			}
			count++
			ids[opID] = struct{}{}
		}
	}

	return ids, count, nil
}

func operationCount(doc *yaml.Node) (int, error) {
	_, count, err := operationIDs(doc)
	return count, err
}

func removeOperations(v1Doc *yaml.Node, removeIDs map[string]struct{}) (removed int, remaining int, _ error) {
	paths, err := pathsNode(v1Doc)
	if err != nil {
		return 0, 0, err
	}

	for i := 0; i+1 < len(paths.Content); {
		pathKey := paths.Content[i]
		pathItem := paths.Content[i+1]
		if pathKey.Kind != yaml.ScalarNode || pathItem.Kind != yaml.MappingNode {
			i += 2
			continue
		}

		for j := 0; j+1 < len(pathItem.Content); {
			methodNode := pathItem.Content[j]
			if methodNode.Kind != yaml.ScalarNode {
				j += 2
				continue
			}
			if _, ok := httpMethods[methodNode.Value]; !ok {
				j += 2
				continue
			}

			opNode := pathItem.Content[j+1]
			opID := scalarMappingValue(opNode, "operationId")
			if opID == "" {
				j += 2
				continue
			}

			if _, ok := removeIDs[opID]; ok {
				pathItem.Content = append(pathItem.Content[:j], pathItem.Content[j+2:]...)
				removed++
				continue
			}

			j += 2
		}

		if !hasAnyHTTPMethod(pathItem) {
			paths.Content = append(paths.Content[:i], paths.Content[i+2:]...)
			continue
		}

		i += 2
	}

	remaining, err = operationCount(v1Doc)
	if err != nil {
		return 0, 0, err
	}
	return removed, remaining, nil
}

func pathsNode(doc *yaml.Node) (*yaml.Node, error) {
	if doc == nil {
		return nil, errors.New("nil document")
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil, errors.New("invalid yaml document")
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, errors.New("expected root mapping node")
	}

	paths := mappingValue(root, "paths")
	if paths == nil {
		return nil, errors.New("missing 'paths' in spec")
	}
	if paths.Kind != yaml.MappingNode {
		return nil, errors.New("'paths' is not a mapping")
	}
	return paths, nil
}

func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

func scalarMappingValue(mapping *yaml.Node, key string) string {
	node := mappingValue(mapping, key)
	if node == nil || node.Kind != yaml.ScalarNode {
		return ""
	}
	return node.Value
}

func hasAnyHTTPMethod(pathItem *yaml.Node) bool {
	if pathItem == nil || pathItem.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i+1 < len(pathItem.Content); i += 2 {
		methodNode := pathItem.Content[i]
		if methodNode.Kind != yaml.ScalarNode {
			continue
		}
		if _, ok := httpMethods[methodNode.Value]; ok {
			return true
		}
	}
	return false
}
