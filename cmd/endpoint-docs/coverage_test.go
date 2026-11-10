package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// This checks route references in handwritten wrappers, not runtime API behavior.
// Generated-client calls and Raw.Do calls must resolve to a spec operation. In v1,
// resolving against the derived legacy spec also rejects operations moved to v2.
func TestSDKOperationCoverage(t *testing.T) {
	for version, spec := range map[string]string{"v1": "derived/v1-legacy.yaml", "v2": "upstream/v2.yaml"} {
		t.Run(version, func(t *testing.T) {
			endpoints, err := parseSpec(filepath.Join("../../openapi", spec))
			if err != nil {
				t.Fatal(err)
			}
			files, err := filepath.Glob(filepath.Join("../../pipedrive", version, "*.go"))
			if err != nil {
				t.Fatal(err)
			}
			var sources []*ast.File
			for _, path := range files {
				if strings.HasSuffix(path, "_test.go") {
					continue
				}
				file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
				if err != nil {
					t.Fatal(err)
				}
				sources = append(sources, file)
			}
			for _, err := range checkCoverage(endpoints, sources) {
				t.Error(err)
			}
		})
	}
}

var pathParameter = regexp.MustCompile(`\{[^}]+\}|%[ds]`)

func routeKey(method, path string) string {
	return strings.ToUpper(method) + " " + pathParameter.ReplaceAllString(path, "{}")
}

func checkCoverage(endpoints []endpoint, files []*ast.File) []error {
	routes := make(map[string]string)
	generated := make(map[string]string)
	for _, e := range endpoints {
		routes[routeKey(e.Method, e.Path)] = e.OperationID
		name := strings.ToUpper(e.OperationID[:1]) + e.OperationID[1:]
		generated[name] = e.OperationID
	}
	covered := make(map[string]bool)
	var problems []error
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			values := make(map[string]ast.Expr)
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				assign, ok := node.(*ast.AssignStmt)
				if ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 {
					if name, ok := assign.Lhs[0].(*ast.Ident); ok {
						values[name.Name] = assign.Rhs[0]
					}
				}
				return true
			})
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				selector, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				receiver, ok := selector.X.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				var operation, reference string
				switch receiver.Sel.Name {
				case "gen":
					name := strings.TrimSuffix(selector.Sel.Name, "WithBody")
					operation, reference = generated[name], name
				case "Raw", "raw":
					if selector.Sel.Name != "Do" || len(call.Args) < 3 {
						return true
					}
					method, ok := call.Args[1].(*ast.SelectorExpr)
					if !ok {
						problems = append(problems, fmt.Errorf("%s: unresolved HTTP method", fn.Name))
						return true
					}
					path := routeLiteral(call.Args[2], values)
					reference = routeKey(strings.TrimPrefix(method.Sel.Name, "Method"), path)
					operation = routes[reference]
					// Both token grants use the same URL, despite separate upstream operations.
					if receiver.Sel.Name == "raw" && path == "/oauth/token" && method.Sel.Name == "MethodPost" {
						switch fn.Name.Name {
						case "GetTokens":
							operation = "get-tokens"
						case "RefreshTokens":
							operation = "refresh-tokens"
						}
					}
				default:
					return true
				}
				if operation == "" {
					problems = append(problems, fmt.Errorf("%s: route %s is absent from this version's spec", fn.Name, reference))
				} else {
					covered[operation] = true
				}
				return true
			})
			// Authorization returns a browser URL instead of making an HTTP request.
			if fn.Name.Name == "Authorize" && fn.Recv != nil {
				ast.Inspect(fn.Body, func(node ast.Node) bool {
					literal, ok := node.(*ast.BasicLit)
					if ok && literal.Kind == token.STRING && literal.Value == `"/oauth/authorize"` {
						covered["authorize"] = true
					}
					return true
				})
			}
		}
	}
	for _, e := range endpoints {
		if !covered[e.OperationID] {
			problems = append(problems, fmt.Errorf("%s %s (%s): missing handwritten wrapper", e.Method, e.Path, e.OperationID))
		}
	}
	return problems
}

func routeLiteral(expr ast.Expr, values map[string]ast.Expr) string {
	if name, ok := expr.(*ast.Ident); ok {
		expr = values[name.Name]
	}
	if call, ok := expr.(*ast.CallExpr); ok {
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Sprintf" || len(call.Args) == 0 {
			return ""
		}
		expr = call.Args[0]
	}
	if literal, ok := expr.(*ast.BasicLit); ok && literal.Kind == token.STRING {
		value, err := strconv.Unquote(literal.Value)
		if err == nil {
			return value
		}
	}
	return ""
}

func TestCoverageRejectsMissingAndOverlappingRoutes(t *testing.T) {
	endpoints := []endpoint{{Method: "GET", Path: "/legacy/{id}", OperationID: "getLegacy"}}
	for _, tc := range []struct {
		name, body string
		failures   int
	}{
		{"generated", `s.client.gen.GetLegacy(ctx, id)`, 0},
		{"raw", `path := fmt.Sprintf("/legacy/%d", id); s.client.Raw.Do(ctx, http.MethodGet, path, nil, nil, nil)`, 0},
		{"missing", `return`, 1},
		{"overlap", `s.client.Raw.Do(ctx, http.MethodGet, "/v2-only", nil, nil, nil)`, 2},
		{"unknown generated", `s.client.gen.GetRemoved(ctx)`, 2},
		{"unresolved path", `s.client.Raw.Do(ctx, http.MethodGet, dynamicPath(), nil, nil, nil)`, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file, err := parser.ParseFile(token.NewFileSet(), "fixture.go", "package sdk; func f() {"+tc.body+"}", 0)
			if err != nil {
				t.Fatal(err)
			}
			if got := checkCoverage(endpoints, []*ast.File{file}); len(got) != tc.failures {
				t.Fatalf("problems = %v, want %d", got, tc.failures)
			}
		})
	}
}
