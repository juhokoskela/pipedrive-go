package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/juhokoskela/pipedrive-go/internal/specdiff"
)

func main() {
	var (
		v1Path     = flag.String("v1", "openapi/upstream/v1.yaml", "path to OpenAPI v1 spec (yaml)")
		v2Path     = flag.String("v2", "openapi/upstream/v2.yaml", "path to OpenAPI v2 spec (yaml)")
		outPath    = flag.String("out", "openapi/derived/v1-legacy.yaml", "output path for derived v1-legacy spec (yaml)")
		reportPath = flag.String("report", "openapi/derived/v1-legacy.report.json", "output path for derivation report (json)")
	)
	flag.Parse()

	v1, err := os.ReadFile(*v1Path)
	fatalIf(err, "read v1 spec")
	v2, err := os.ReadFile(*v2Path)
	fatalIf(err, "read v2 spec")

	derived, report, err := specdiff.DeriveV1Legacy(v1, v2, nil)
	fatalIf(err, "derive v1-legacy spec")

	fatalIf(os.MkdirAll(filepath.Dir(*outPath), 0o755), "create output directory")
	fatalIf(os.WriteFile(*outPath, derived, 0o644), "write derived spec")

	reportJSON, err := json.MarshalIndent(report, "", "  ")
	fatalIf(err, "encode report")
	reportJSON = append(reportJSON, '\n')

	fatalIf(os.MkdirAll(filepath.Dir(*reportPath), 0o755), "create report directory")
	fatalIf(os.WriteFile(*reportPath, reportJSON, 0o644), "write report")

	fmt.Fprintf(os.Stderr, "derived %s (v1 ops=%d, v2 ops=%d, removed=%d, legacy=%d)\n",
		*outPath, report.V1Operations, report.V2Operations, report.RemovedOperations, report.LegacyOperations)
}

func fatalIf(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}

