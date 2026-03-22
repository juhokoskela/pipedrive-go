package v1

type Activity struct {
	ID                ActivityID      `json:"id"`
	Subject           string          `json:"subject,omitempty"`
	Type              string          `json:"type,omitempty"`
	Done              bool            `json:"done,omitempty"`
	Busy              bool            `json:"busy_flag,omitempty"`
	Active            bool            `json:"active_flag,omitempty"`
	DueDate           string          `json:"due_date,omitempty"`
	DueTime           string          `json:"due_time,omitempty"`
	Duration          string          `json:"duration,omitempty"`
	UserID            *UserID         `json:"user_id,omitempty"`
	UpdateUserID      *UserID         `json:"update_user_id,omitempty"`
	DealID            *DealID         `json:"deal_id,omitempty"`
	PersonID          *PersonID       `json:"person_id,omitempty"`
	OrganizationID    *OrganizationID `json:"org_id,omitempty"`
	LeadID            *LeadID         `json:"lead_id,omitempty"`
	ProjectID         *ProjectID      `json:"project_id,omitempty"`
	PublicDescription string          `json:"public_description,omitempty"`
	SourceTimezone    string          `json:"source_timezone,omitempty"`
	AddTime           *DateTime       `json:"add_time,omitempty"`
	UpdateTime        *DateTime       `json:"update_time,omitempty"`
	MarkedAsDoneTime  *DateTime       `json:"marked_as_done_time,omitempty"`
}
