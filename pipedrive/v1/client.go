package v1

import (
	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

const DefaultBaseURL = "https://api.pipedrive.com/v1"

type Client struct {
	Raw   *pipedrive.RawClient
	OAuth *OAuthService

	gen *genv1.ClientWithResponses

	Currencies                *CurrenciesService
	Activities                *ActivitiesService
	ActivityTypes             *ActivityTypesService
	CallLogs                  *CallLogsService
	Channels                  *ChannelsService
	Billing                   *BillingService
	LeadLabels                *LeadLabelsService
	LeadSources               *LeadSourcesService
	Leads                     *LeadsService
	LeadFields                *LeadFieldsService
	DealFields                *DealFieldsService
	Deals                     *DealsService
	PersonFields              *PersonFieldsService
	Persons                   *PersonsService
	OrganizationFields        *OrganizationFieldsService
	Organizations             *OrganizationsService
	ProductFields             *ProductFieldsService
	Products                  *ProductsService
	Files                     *FilesService
	NoteFields                *NoteFieldsService
	Notes                     *NotesService
	Stages                    *StagesService
	Filters                   *FiltersService
	Goals                     *GoalsService
	Mailbox                   *MailboxService
	Meetings                  *MeetingsService
	Projects                  *ProjectsService
	ProjectTemplates          *ProjectTemplatesService
	Roles                     *RolesService
	Teams                     *TeamsService
	Tasks                     *TasksService
	Webhooks                  *WebhooksService
	UserConnections           *UserConnectionsService
	UserSettings              *UserSettingsService
	PermissionSets            *PermissionSetsService
	Recents                   *RecentsService
	Pipelines                 *PipelinesService
	Users                     *UsersService
	OrganizationRelationships *OrganizationRelationshipsService
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

	oauthBaseURL := DefaultOAuthBaseURL
	oauthRaw, err := pipedrive.NewRawClient(oauthBaseURL, httpClient)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Raw:   raw,
		OAuth: &OAuthService{raw: oauthRaw, baseURL: oauthBaseURL},
		gen:   gen,
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
	c.LeadFields = &LeadFieldsService{client: c}
	c.DealFields = &DealFieldsService{client: c}
	c.Deals = &DealsService{client: c}
	c.PersonFields = &PersonFieldsService{client: c}
	c.Persons = &PersonsService{client: c}
	c.OrganizationFields = &OrganizationFieldsService{client: c}
	c.Organizations = &OrganizationsService{client: c}
	c.ProductFields = &ProductFieldsService{client: c}
	c.Products = &ProductsService{client: c}
	c.Files = &FilesService{client: c}
	c.NoteFields = &NoteFieldsService{client: c}
	c.Notes = &NotesService{client: c}
	c.Stages = &StagesService{client: c}
	c.Filters = &FiltersService{client: c}
	c.Goals = &GoalsService{client: c}
	c.Mailbox = &MailboxService{client: c}
	c.Meetings = &MeetingsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.ProjectTemplates = &ProjectTemplatesService{client: c}
	c.Roles = &RolesService{client: c}
	c.Teams = &TeamsService{client: c}
	c.Tasks = &TasksService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.UserConnections = &UserConnectionsService{client: c}
	c.UserSettings = &UserSettingsService{client: c}
	c.PermissionSets = &PermissionSetsService{client: c}
	c.Recents = &RecentsService{client: c}
	c.Pipelines = &PipelinesService{client: c}
	c.Users = &UsersService{client: c}
	c.OrganizationRelationships = &OrganizationRelationshipsService{client: c}
	return c, nil
}
