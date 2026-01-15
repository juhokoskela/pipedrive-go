package v2

import (
	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

const DefaultBaseURL = "https://api.pipedrive.com/api/v2"

type Client struct {
	Raw *pipedrive.RawClient

	gen *genv2.ClientWithResponses

	Deals              *DealsService
	DealFields         *DealFieldsService
	Persons            *PersonsService
	PersonFields       *PersonFieldsService
	Organizations      *OrganizationsService
	OrganizationFields *OrganizationFieldsService
	Activities         *ActivitiesService
	ActivityFields     *ActivityFieldsService
	Pipelines          *PipelinesService
	Stages             *StagesService
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

	gen, err := genv2.NewClientWithResponses(baseURL, genv2.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Raw: raw,
		gen: gen,
	}
	c.Deals = &DealsService{client: c}
	c.DealFields = &DealFieldsService{client: c}
	c.Persons = &PersonsService{client: c}
	c.PersonFields = &PersonFieldsService{client: c}
	c.Organizations = &OrganizationsService{client: c}
	c.OrganizationFields = &OrganizationFieldsService{client: c}
	c.Activities = &ActivitiesService{client: c}
	c.ActivityFields = &ActivityFieldsService{client: c}
	c.Pipelines = &PipelinesService{client: c}
	c.Stages = &StagesService{client: c}
	return c, nil
}
