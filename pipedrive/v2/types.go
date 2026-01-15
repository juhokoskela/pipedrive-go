package v2

type DealID int64

type PersonID int64

type OrganizationID int64

type ActivityID int64

type LeadID string

type ConversionID string

type UserID int64

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

type ConversionStatus string

const (
	ConversionStatusNotStarted ConversionStatus = "not_started"
	ConversionStatusRunning    ConversionStatus = "running"
	ConversionStatusCompleted  ConversionStatus = "completed"
	ConversionStatusFailed     ConversionStatus = "failed"
	ConversionStatusRejected   ConversionStatus = "rejected"
)
