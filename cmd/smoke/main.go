package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v1 "github.com/juhokoskela/pipedrive-go/pipedrive/v1"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

const defaultTimeout = 20 * time.Second

func main() {
	token := strings.TrimSpace(os.Getenv("PIPEDRIVE_API_TOKEN"))
	if token == "" {
		fatalf("PIPEDRIVE_API_TOKEN is required")
	}

	timeout := defaultTimeout
	if envTimeout := strings.TrimSpace(os.Getenv("PIPEDRIVE_SMOKE_TIMEOUT")); envTimeout != "" {
		parsed, err := time.ParseDuration(envTimeout)
		if err != nil {
			fatalf("invalid PIPEDRIVE_SMOKE_TIMEOUT: %v", err)
		}
		if parsed > 0 {
			timeout = parsed
		}
	}

	cfgV1 := pipedrive.Config{Auth: pipedrive.APITokenAuth(token)}
	if baseURL := strings.TrimSpace(os.Getenv("PIPEDRIVE_BASE_URL_V1")); baseURL != "" {
		cfgV1.BaseURL = baseURL
	}

	v1Client, err := v1.NewClient(cfgV1)
	if err != nil {
		fatalf("v1.NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	user, err := v1Client.Users.GetCurrent(ctx)
	cancel()
	if err != nil {
		fatalf("v1 Users.GetCurrent: %v", err)
	}
	fmt.Printf("v1 current user: id=%d name=%s\n", user.ID, user.Name)

	cfgV2 := pipedrive.Config{Auth: pipedrive.APITokenAuth(token)}
	if baseURL := strings.TrimSpace(os.Getenv("PIPEDRIVE_BASE_URL_V2")); baseURL != "" {
		cfgV2.BaseURL = baseURL
	}

	v2Client, err := v2.NewClient(cfgV2)
	if err != nil {
		fatalf("v2.NewClient: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	pipelines, next, err := v2Client.Pipelines.List(ctx)
	cancel()
	if err != nil {
		fatalf("v2 Pipelines.List: %v", err)
	}

	fmt.Printf("v2 pipelines: count=%d", len(pipelines))
	if next != nil {
		fmt.Printf(" next_cursor=%s", *next)
	}
	fmt.Println()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
