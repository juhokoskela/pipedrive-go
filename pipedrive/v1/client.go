package v1

import (
	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

const DefaultBaseURL = "https://api.pipedrive.com/v1"

type Client struct {
	Raw *pipedrive.RawClient

	gen *genv1.ClientWithResponses

	Currencies      *CurrenciesService
	Activities      *ActivitiesService
	ActivityTypes   *ActivityTypesService
	CallLogs        *CallLogsService
	Channels        *ChannelsService
	Billing         *BillingService
	LeadLabels      *LeadLabelsService
	LeadSources     *LeadSourcesService
	Leads           *LeadsService
	UserConnections *UserConnectionsService
	UserSettings    *UserSettingsService
}

func NewClient(cfg pipedrive.Config) (*Client, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := pipedrive.NewHTTPClient(cfg)

	raw, err := pipedrive.NewRawClient(baseURL, httpClient)
	if err != nil {
		return nil, err
	}

	gen, err := genv1.NewClientWithResponses(baseURL, genv1.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Raw: raw,
		gen: gen,
	}
	c.Currencies = &CurrenciesService{client: c}
	c.Activities = &ActivitiesService{client: c}
	c.ActivityTypes = &ActivityTypesService{client: c}
	c.CallLogs = &CallLogsService{client: c}
	c.Channels = &ChannelsService{client: c}
	c.Billing = &BillingService{client: c}
	c.LeadLabels = &LeadLabelsService{client: c}
	c.LeadSources = &LeadSourcesService{client: c}
	c.Leads = &LeadsService{client: c}
	c.UserConnections = &UserConnectionsService{client: c}
	c.UserSettings = &UserSettingsService{client: c}
	return c, nil
}
