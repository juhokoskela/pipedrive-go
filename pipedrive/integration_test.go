package pipedrive_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v1 "github.com/juhokoskela/pipedrive-go/pipedrive/v1"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

const integrationTimeout = 20 * time.Second

func integrationToken(t *testing.T) string {
	t.Helper()

	token := os.Getenv("PIPEDRIVE_API_TOKEN")
	if token == "" {
		t.Skip("set PIPEDRIVE_API_TOKEN to run integration tests")
	}
	return token
}

func integrationWriteEnabled(t *testing.T) {
	t.Helper()

	if os.Getenv("PIPEDRIVE_INTEGRATION_WRITE") != "1" {
		t.Skip("set PIPEDRIVE_INTEGRATION_WRITE=1 to run write integration tests")
	}
}

func newV1Client(t *testing.T, token string) *v1.Client {
	t.Helper()

	cfg := pipedrive.Config{Auth: pipedrive.APITokenAuth(token)}
	if baseURL := os.Getenv("PIPEDRIVE_BASE_URL_V1"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	client, err := v1.NewClient(cfg)
	if err != nil {
		t.Fatalf("v1.NewClient error: %v", err)
	}
	return client
}

func newV2Client(t *testing.T, token string) *v2.Client {
	t.Helper()

	cfg := pipedrive.Config{Auth: pipedrive.APITokenAuth(token)}
	if baseURL := os.Getenv("PIPEDRIVE_BASE_URL_V2"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	client, err := v2.NewClient(cfg)
	if err != nil {
		t.Fatalf("v2.NewClient error: %v", err)
	}
	return client
}

func TestIntegrationV1CurrentUser(t *testing.T) {
	token := integrationToken(t)
	client := newV1Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	user, err := client.Users.GetCurrent(ctx)
	if err != nil {
		t.Fatalf("Users.GetCurrent error: %v", err)
	}
	if user == nil || user.ID == 0 {
		t.Fatalf("unexpected current user: %#v", user)
	}
}

func TestIntegrationV1ListCurrencies(t *testing.T) {
	token := integrationToken(t)
	client := newV1Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	currencies, err := client.Currencies.List(ctx, v1.ListCurrenciesRequest{})
	if err != nil {
		t.Fatalf("Currencies.List error: %v", err)
	}
	if len(currencies) == 0 {
		t.Fatalf("expected currencies, got empty list")
	}
}

func TestIntegrationV2ListPipelines(t *testing.T) {
	token := integrationToken(t)
	client := newV2Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	pipelines, _, err := client.Pipelines.List(ctx)
	if err != nil {
		t.Fatalf("Pipelines.List error: %v", err)
	}
	if pipelines == nil {
		t.Fatalf("unexpected pipelines response: %#v", pipelines)
	}
}

func TestIntegrationV2ListOrganizations(t *testing.T) {
	token := integrationToken(t)
	client := newV2Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	organizations, _, err := client.Organizations.List(ctx, v2.WithOrganizationsPageSize(1))
	if err != nil {
		t.Fatalf("Organizations.List error: %v", err)
	}
	if organizations == nil {
		t.Fatalf("unexpected organizations response: %#v", organizations)
	}
}

func TestIntegrationV2ListPersons(t *testing.T) {
	token := integrationToken(t)
	client := newV2Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	persons, _, err := client.Persons.List(ctx, v2.WithPersonsPageSize(1))
	if err != nil {
		t.Fatalf("Persons.List error: %v", err)
	}
	if persons == nil {
		t.Fatalf("unexpected persons response: %#v", persons)
	}
}

func TestIntegrationV2OrganizationCreateDelete(t *testing.T) {
	token := integrationToken(t)
	integrationWriteEnabled(t)
	client := newV2Client(t, token)

	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()

	name := fmt.Sprintf("sdk-integration-%d", time.Now().UnixNano())
	org, err := client.Organizations.Create(ctx, v2.WithOrganizationName(name))
	if err != nil {
		t.Fatalf("Organizations.Create error: %v", err)
	}
	if org == nil || org.ID == 0 {
		t.Fatalf("unexpected organization response: %#v", org)
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), integrationTimeout)
		defer cleanupCancel()
		_, _ = client.Organizations.Delete(cleanupCtx, org.ID)
	})

	result, err := client.Organizations.Delete(ctx, org.ID)
	if err != nil {
		t.Fatalf("Organizations.Delete error: %v", err)
	}
	if result == nil || result.ID != org.ID {
		t.Fatalf("unexpected delete result: %#v", result)
	}
	deleted = true
}
