package v2

type DealID int64

type DealProductAttachmentID int64

type InstallmentID int64

type AdditionalDiscountID string

type PersonID int64

type OrganizationID int64

type ProductID int64

type ProductVariationID int64

type ProductImageID int64

type ActivityID int64

type ProjectID int64

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
