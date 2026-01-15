package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type DealStatus string

const (
	DealStatusOpen    DealStatus = "open"
	DealStatusWon     DealStatus = "won"
	DealStatusLost    DealStatus = "lost"
	DealStatusDeleted DealStatus = "deleted"
)

type DealSortField string

const (
	DealSortByAddTime    DealSortField = "add_time"
	DealSortByID         DealSortField = "id"
	DealSortByUpdateTime DealSortField = "update_time"
)

type DealIncludeField string

const (
	DealIncludeFieldActivitiesCount       DealIncludeField = "activities_count"
	DealIncludeFieldDoneActivitiesCount   DealIncludeField = "done_activities_count"
	DealIncludeFieldEmailMessagesCount    DealIncludeField = "email_messages_count"
	DealIncludeFieldFilesCount            DealIncludeField = "files_count"
	DealIncludeFieldFirstWonTime          DealIncludeField = "first_won_time"
	DealIncludeFieldFollowersCount        DealIncludeField = "followers_count"
	DealIncludeFieldLastActivityID        DealIncludeField = "last_activity_id"
	DealIncludeFieldLastIncomingMailTime  DealIncludeField = "last_incoming_mail_time"
	DealIncludeFieldLastOutgoingMailTime  DealIncludeField = "last_outgoing_mail_time"
	DealIncludeFieldNextActivityID        DealIncludeField = "next_activity_id"
	DealIncludeFieldNotesCount            DealIncludeField = "notes_count"
	DealIncludeFieldParticipantsCount     DealIncludeField = "participants_count"
	DealIncludeFieldProductsCount         DealIncludeField = "products_count"
	DealIncludeFieldSmartBccEmail         DealIncludeField = "smart_bcc_email"
	DealIncludeFieldUndoneActivitiesCount DealIncludeField = "undone_activities_count"
)

type DealSearchField string

const (
	DealSearchFieldCustomFields DealSearchField = "custom_fields"
	DealSearchFieldNotes        DealSearchField = "notes"
	DealSearchFieldTitle        DealSearchField = "title"
)

type DealSearchStatus string

const (
	DealSearchStatusOpen DealSearchStatus = "open"
	DealSearchStatusWon  DealSearchStatus = "won"
	DealSearchStatusLost DealSearchStatus = "lost"
)

type DealSearchIncludeField string

const (
	DealSearchIncludeFieldDealCCEmail DealSearchIncludeField = "deal.cc_email"
)

type DealProductSortField string

const (
	DealProductSortByAddTime    DealProductSortField = "add_time"
	DealProductSortByDealID     DealProductSortField = "deal_id"
	DealProductSortByID         DealProductSortField = "id"
	DealProductSortByOrderNr    DealProductSortField = "order_nr"
	DealProductSortByUpdateTime DealProductSortField = "update_time"
)

type DealProductDiscountType string

const (
	DealProductDiscountTypeAmount     DealProductDiscountType = "amount"
	DealProductDiscountTypePercentage DealProductDiscountType = "percentage"
)

type DealProductTaxMethod string

const (
	DealProductTaxMethodExclusive DealProductTaxMethod = "exclusive"
	DealProductTaxMethodInclusive DealProductTaxMethod = "inclusive"
	DealProductTaxMethodNone      DealProductTaxMethod = "none"
)

type AdditionalDiscountType string

const (
	AdditionalDiscountTypeAmount     AdditionalDiscountType = "amount"
	AdditionalDiscountTypePercentage AdditionalDiscountType = "percentage"
)

type InstallmentSortField string

const (
	InstallmentSortByBillingDate InstallmentSortField = "billing_date"
	InstallmentSortByDealID      InstallmentSortField = "deal_id"
	InstallmentSortByID          InstallmentSortField = "id"
)

type Deal struct {
	ID                DealID                 `json:"id"`
	Title             string                 `json:"title,omitempty"`
	Value             *float64               `json:"value,omitempty"`
	Currency          string                 `json:"currency,omitempty"`
	Status            DealStatus             `json:"status,omitempty"`
	OwnerID           *UserID                `json:"owner_id,omitempty"`
	PersonID          *PersonID              `json:"person_id,omitempty"`
	OrgID             *OrganizationID        `json:"org_id,omitempty"`
	StageID           *StageID               `json:"stage_id,omitempty"`
	PipelineID        *PipelineID            `json:"pipeline_id,omitempty"`
	AddTime           *time.Time             `json:"add_time,omitempty"`
	UpdateTime        *time.Time             `json:"update_time,omitempty"`
	ExpectedCloseDate *string                `json:"expected_close_date,omitempty"`
	LostReason        *string                `json:"lost_reason,omitempty"`
	WonTime           *time.Time             `json:"won_time,omitempty"`
	LostTime          *time.Time             `json:"lost_time,omitempty"`
	IsArchived        bool                   `json:"is_archived,omitempty"`
	IsDeleted         bool                   `json:"is_deleted,omitempty"`
	VisibleTo         *int                   `json:"visible_to,omitempty"`
	LabelIDs          []int                  `json:"label_ids,omitempty"`
	CustomFields      map[string]interface{} `json:"custom_fields,omitempty"`
}

type DealSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type DealSearchResults struct {
	Items []DealSearchItem `json:"items,omitempty"`
}

type DealDeleteResult struct {
	ID DealID `json:"id"`
}

type DealConversionJob struct {
	ConversionID ConversionID `json:"conversion_id"`
}

type DealConversionStatus struct {
	LeadID       *LeadID          `json:"lead_id,omitempty"`
	DealID       *DealID          `json:"deal_id,omitempty"`
	ConversionID ConversionID     `json:"conversion_id,omitempty"`
	Status       ConversionStatus `json:"status,omitempty"`
}

type DealProduct struct {
	ID                     DealProductAttachmentID `json:"id"`
	DealID                 *DealID                 `json:"deal_id,omitempty"`
	ProductID              *ProductID              `json:"product_id,omitempty"`
	ProductVariationID     *ProductVariationID     `json:"product_variation_id,omitempty"`
	Name                   string                  `json:"name,omitempty"`
	ItemPrice              *float64                `json:"item_price,omitempty"`
	Quantity               *float64                `json:"quantity,omitempty"`
	Discount               *float64                `json:"discount,omitempty"`
	DiscountType           DealProductDiscountType `json:"discount_type,omitempty"`
	Tax                    *float64                `json:"tax,omitempty"`
	TaxMethod              DealProductTaxMethod    `json:"tax_method,omitempty"`
	BillingFrequency       BillingFrequency        `json:"billing_frequency,omitempty"`
	BillingFrequencyCycles *int                    `json:"billing_frequency_cycles,omitempty"`
	BillingStartDate       *string                 `json:"billing_start_date,omitempty"`
	Currency               string                  `json:"currency,omitempty"`
	Comments               string                  `json:"comments,omitempty"`
	IsEnabled              *bool                   `json:"is_enabled,omitempty"`
	OrderNr                *int                    `json:"order_nr,omitempty"`
	Sum                    *float64                `json:"sum,omitempty"`
	AddTime                *time.Time              `json:"add_time,omitempty"`
	UpdateTime             *time.Time              `json:"update_time,omitempty"`
}

type DealProductInput struct {
	payload dealProductPayload
}

type DealProductDeleteResult struct {
	ID DealProductAttachmentID `json:"id"`
}

type DealProductsDeleteResult struct {
	IDs                   []DealProductAttachmentID `json:"ids"`
	MoreItemsInCollection *bool                     `json:"more_items_in_collection,omitempty"`
}

type AdditionalDiscount struct {
	ID          AdditionalDiscountID   `json:"id"`
	DealID      *DealID                `json:"deal_id,omitempty"`
	Amount      *float64               `json:"amount,omitempty"`
	Type        AdditionalDiscountType `json:"type,omitempty"`
	Description *string                `json:"description,omitempty"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
	CreatedBy   *UserID                `json:"created_by,omitempty"`
	UpdatedBy   *UserID                `json:"updated_by,omitempty"`
}

type AdditionalDiscountDeleteResult struct {
	ID AdditionalDiscountID `json:"id"`
}

type Installment struct {
	ID          InstallmentID `json:"id"`
	DealID      *DealID       `json:"deal_id,omitempty"`
	Amount      *float64      `json:"amount,omitempty"`
	BillingDate *string       `json:"billing_date,omitempty"`
	Description *string       `json:"description,omitempty"`
}

type InstallmentDeleteResult struct {
	ID InstallmentID `json:"id"`
}

type DealsService struct {
	client *Client
}

type GetDealOption interface {
	applyGetDeal(*getDealOptions)
}

type ListDealsOption interface {
	applyListDeals(*listDealsOptions)
}

type ListArchivedDealsOption interface {
	applyListArchivedDeals(*listArchivedDealsOptions)
}

type CreateDealOption interface {
	applyCreateDeal(*createDealOptions)
}

type UpdateDealOption interface {
	applyUpdateDeal(*updateDealOptions)
}

type DeleteDealOption interface {
	applyDeleteDeal(*deleteDealOptions)
}

type SearchDealsOption interface {
	applySearchDeals(*searchDealsOptions)
}

type ConvertDealOption interface {
	applyConvertDeal(*convertDealOptions)
}

type GetDealConversionStatusOption interface {
	applyGetDealConversionStatus(*getDealConversionStatusOptions)
}

type GetDealFollowersOption interface {
	applyGetDealFollowers(*getDealFollowersOptions)
}

type AddDealFollowerOption interface {
	applyAddDealFollower(*addDealFollowerOptions)
}

type DeleteDealFollowerOption interface {
	applyDeleteDealFollower(*deleteDealFollowerOptions)
}

type GetDealFollowersChangelogOption interface {
	applyGetDealFollowersChangelog(*getDealFollowersChangelogOptions)
}

type ListDealProductsOption interface {
	applyListDealProducts(*listDealProductsOptions)
}

type ListDealsProductsOption interface {
	applyListDealsProducts(*listDealsProductsOptions)
}

type AddDealProductOption interface {
	applyAddDealProduct(*addDealProductOptions)
}

type AddManyDealProductsOption interface {
	applyAddManyDealProducts(*addManyDealProductsOptions)
}

type UpdateDealProductOption interface {
	applyUpdateDealProduct(*updateDealProductOptions)
}

type DeleteDealProductOption interface {
	applyDeleteDealProduct(*deleteDealProductOptions)
}

type DeleteDealProductsOption interface {
	applyDeleteDealProducts(*deleteDealProductsOptions)
}

type ListAdditionalDiscountsOption interface {
	applyListAdditionalDiscounts(*listAdditionalDiscountsOptions)
}

type AddAdditionalDiscountOption interface {
	applyAddAdditionalDiscount(*addAdditionalDiscountOptions)
}

type UpdateAdditionalDiscountOption interface {
	applyUpdateAdditionalDiscount(*updateAdditionalDiscountOptions)
}

type DeleteAdditionalDiscountOption interface {
	applyDeleteAdditionalDiscount(*deleteAdditionalDiscountOptions)
}

type ListInstallmentsOption interface {
	applyListInstallments(*listInstallmentsOptions)
}

type AddInstallmentOption interface {
	applyAddInstallment(*addInstallmentOptions)
}

type UpdateInstallmentOption interface {
	applyUpdateInstallment(*updateInstallmentOptions)
}

type DeleteInstallmentOption interface {
	applyDeleteInstallment(*deleteInstallmentOptions)
}

type DealRequestOption interface {
	GetDealOption
	ListDealsOption
	ListArchivedDealsOption
	CreateDealOption
	UpdateDealOption
	DeleteDealOption
	SearchDealsOption
	ConvertDealOption
	GetDealConversionStatusOption
	GetDealFollowersOption
	AddDealFollowerOption
	DeleteDealFollowerOption
	GetDealFollowersChangelogOption
	ListDealProductsOption
	ListDealsProductsOption
	AddDealProductOption
	AddManyDealProductsOption
	UpdateDealProductOption
	DeleteDealProductOption
	DeleteDealProductsOption
	ListAdditionalDiscountsOption
	AddAdditionalDiscountOption
	UpdateAdditionalDiscountOption
	DeleteAdditionalDiscountOption
	ListInstallmentsOption
	AddInstallmentOption
	UpdateInstallmentOption
	DeleteInstallmentOption
}

type DealOption interface {
	CreateDealOption
	UpdateDealOption
}

type DealProductOption interface {
	AddDealProductOption
	UpdateDealProductOption
}

type AdditionalDiscountOption interface {
	AddAdditionalDiscountOption
	UpdateAdditionalDiscountOption
}

type InstallmentOption interface {
	AddInstallmentOption
	UpdateInstallmentOption
}

type getDealOptions struct {
	params         genv2.GetDealParams
	requestOptions []pipedrive.RequestOption
}

type listDealsOptions struct {
	params         genv2.GetDealsParams
	requestOptions []pipedrive.RequestOption
}

type listArchivedDealsOptions struct {
	params         genv2.GetArchivedDealsParams
	requestOptions []pipedrive.RequestOption
}

type createDealOptions struct {
	payload        dealPayload
	requestOptions []pipedrive.RequestOption
}

type updateDealOptions struct {
	payload        dealPayload
	requestOptions []pipedrive.RequestOption
}

type deleteDealOptions struct {
	requestOptions []pipedrive.RequestOption
}

type searchDealsOptions struct {
	params         genv2.SearchDealsParams
	requestOptions []pipedrive.RequestOption
}

type convertDealOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getDealConversionStatusOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getDealFollowersOptions struct {
	params         genv2.GetDealFollowersParams
	requestOptions []pipedrive.RequestOption
}

type addDealFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteDealFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getDealFollowersChangelogOptions struct {
	params         genv2.GetDealFollowersChangelogParams
	requestOptions []pipedrive.RequestOption
}

type listDealProductsOptions struct {
	params         genv2.GetDealProductsParams
	requestOptions []pipedrive.RequestOption
}

type listDealsProductsOptions struct {
	params         genv2.GetDealsProductsParams
	requestOptions []pipedrive.RequestOption
}

type addDealProductOptions struct {
	payload        dealProductPayload
	requestOptions []pipedrive.RequestOption
}

type addManyDealProductsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updateDealProductOptions struct {
	payload        dealProductPayload
	requestOptions []pipedrive.RequestOption
}

type deleteDealProductOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteDealProductsOptions struct {
	params         genv2.DeleteManyDealProductsParams
	requestOptions []pipedrive.RequestOption
}

type listAdditionalDiscountsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addAdditionalDiscountOptions struct {
	payload        additionalDiscountPayload
	requestOptions []pipedrive.RequestOption
}

type updateAdditionalDiscountOptions struct {
	payload        additionalDiscountPayload
	requestOptions []pipedrive.RequestOption
}

type deleteAdditionalDiscountOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listInstallmentsOptions struct {
	params         genv2.GetInstallmentsParams
	requestOptions []pipedrive.RequestOption
}

type addInstallmentOptions struct {
	payload        installmentPayload
	requestOptions []pipedrive.RequestOption
}

type updateInstallmentOptions struct {
	payload        installmentPayload
	requestOptions []pipedrive.RequestOption
}

type deleteInstallmentOptions struct {
	requestOptions []pipedrive.RequestOption
}

type dealPayload struct {
	title             *string
	value             *float64
	currency          *string
	ownerID           *UserID
	personID          *PersonID
	orgID             *OrganizationID
	stageID           *StageID
	pipelineID        *PipelineID
	status            *DealStatus
	expectedCloseDate *string
	probability       *float64
	lostReason        *string
	visibleTo         *int
	labelIDs          []int
	customFields      map[string]interface{}
	isArchived        *bool
	isDeleted         *bool
	archiveTime       *string
	closeTime         *string
	lostTime          *string
	wonTime           *string
}

type dealProductPayload struct {
	productID              *ProductID
	productVariationID     *ProductVariationID
	itemPrice              *float64
	quantity               *float64
	discount               *float64
	discountType           *DealProductDiscountType
	comments               *string
	tax                    *float64
	taxMethod              *DealProductTaxMethod
	isEnabled              *bool
	billingFrequency       *BillingFrequency
	billingFrequencyCycles *int
	billingStartDate       *string
}

type additionalDiscountPayload struct {
	amount       *float64
	discountType *AdditionalDiscountType
	description  *string
}

type installmentPayload struct {
	amount      *float64
	billingDate *string
	description *string
}

type dealRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o dealRequestOptions) applyGetDeal(cfg *getDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListDeals(cfg *listDealsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListArchivedDeals(cfg *listArchivedDealsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyCreateDeal(cfg *createDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyUpdateDeal(cfg *updateDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteDeal(cfg *deleteDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applySearchDeals(cfg *searchDealsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyConvertDeal(cfg *convertDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyGetDealConversionStatus(cfg *getDealConversionStatusOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyGetDealFollowers(cfg *getDealFollowersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyAddDealFollower(cfg *addDealFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteDealFollower(cfg *deleteDealFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyGetDealFollowersChangelog(cfg *getDealFollowersChangelogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListDealProducts(cfg *listDealProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListDealsProducts(cfg *listDealsProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyAddDealProduct(cfg *addDealProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyAddManyDealProducts(cfg *addManyDealProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyUpdateDealProduct(cfg *updateDealProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteDealProduct(cfg *deleteDealProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteDealProducts(cfg *deleteDealProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListAdditionalDiscounts(cfg *listAdditionalDiscountsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyAddAdditionalDiscount(cfg *addAdditionalDiscountOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyUpdateAdditionalDiscount(cfg *updateAdditionalDiscountOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteAdditionalDiscount(cfg *deleteAdditionalDiscountOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListInstallments(cfg *listInstallmentsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyAddInstallment(cfg *addInstallmentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyUpdateInstallment(cfg *updateInstallmentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyDeleteInstallment(cfg *deleteInstallmentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getDealOptionFunc func(*getDealOptions)

func (f getDealOptionFunc) applyGetDeal(cfg *getDealOptions) {
	f(cfg)
}

type listDealsOptionFunc func(*listDealsOptions)

func (f listDealsOptionFunc) applyListDeals(cfg *listDealsOptions) {
	f(cfg)
}

type listArchivedDealsOptionFunc func(*listArchivedDealsOptions)

func (f listArchivedDealsOptionFunc) applyListArchivedDeals(cfg *listArchivedDealsOptions) {
	f(cfg)
}

type searchDealsOptionFunc func(*searchDealsOptions)

func (f searchDealsOptionFunc) applySearchDeals(cfg *searchDealsOptions) {
	f(cfg)
}

type getDealFollowersOptionFunc func(*getDealFollowersOptions)

func (f getDealFollowersOptionFunc) applyGetDealFollowers(cfg *getDealFollowersOptions) {
	f(cfg)
}

type getDealFollowersChangelogOptionFunc func(*getDealFollowersChangelogOptions)

func (f getDealFollowersChangelogOptionFunc) applyGetDealFollowersChangelog(cfg *getDealFollowersChangelogOptions) {
	f(cfg)
}

type listDealProductsOptionFunc func(*listDealProductsOptions)

func (f listDealProductsOptionFunc) applyListDealProducts(cfg *listDealProductsOptions) {
	f(cfg)
}

type listDealsProductsOptionFunc func(*listDealsProductsOptions)

func (f listDealsProductsOptionFunc) applyListDealsProducts(cfg *listDealsProductsOptions) {
	f(cfg)
}

type listInstallmentsOptionFunc func(*listInstallmentsOptions)

func (f listInstallmentsOptionFunc) applyListInstallments(cfg *listInstallmentsOptions) {
	f(cfg)
}

type deleteDealProductsOptionFunc func(*deleteDealProductsOptions)

func (f deleteDealProductsOptionFunc) applyDeleteDealProducts(cfg *deleteDealProductsOptions) {
	f(cfg)
}

type dealFieldOption func(*dealPayload)

func (f dealFieldOption) applyCreateDeal(cfg *createDealOptions) {
	f(&cfg.payload)
}

func (f dealFieldOption) applyUpdateDeal(cfg *updateDealOptions) {
	f(&cfg.payload)
}

type dealProductFieldOption func(*dealProductPayload)

func (f dealProductFieldOption) applyAddDealProduct(cfg *addDealProductOptions) {
	f(&cfg.payload)
}

func (f dealProductFieldOption) applyUpdateDealProduct(cfg *updateDealProductOptions) {
	f(&cfg.payload)
}

type additionalDiscountFieldOption func(*additionalDiscountPayload)

func (f additionalDiscountFieldOption) applyAddAdditionalDiscount(cfg *addAdditionalDiscountOptions) {
	f(&cfg.payload)
}

func (f additionalDiscountFieldOption) applyUpdateAdditionalDiscount(cfg *updateAdditionalDiscountOptions) {
	f(&cfg.payload)
}

type installmentFieldOption func(*installmentPayload)

func (f installmentFieldOption) applyAddInstallment(cfg *addInstallmentOptions) {
	f(&cfg.payload)
}

func (f installmentFieldOption) applyUpdateInstallment(cfg *updateInstallmentOptions) {
	f(&cfg.payload)
}

func WithDealRequestOptions(opts ...pipedrive.RequestOption) DealRequestOption {
	return dealRequestOptions{requestOptions: opts}
}

func WithDealsFilterID(id int) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithDealsIDs(ids ...DealID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithDealsOwnerID(id UserID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := int(id)
		cfg.params.OwnerId = &value
	})
}

func WithDealsPersonID(id PersonID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithDealsOrganizationID(id OrganizationID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := int(id)
		cfg.params.OrgId = &value
	})
}

func WithDealsPipelineID(id PipelineID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := int(id)
		cfg.params.PipelineId = &value
	})
}

func WithDealsStageID(id StageID) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := int(id)
		cfg.params.StageId = &value
	})
}

func WithDealsStatus(statuses ...DealStatus) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		csv := joinCSV(statuses)
		if csv == "" {
			return
		}
		value := genv2.GetDealsParamsStatus(csv)
		cfg.params.Status = &value
	})
}

func WithDealsUpdatedSince(t time.Time) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedSince = &value
	})
}

func WithDealsUpdatedUntil(t time.Time) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedUntil = &value
	})
}

func WithDealsSortBy(field DealSortField) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := genv2.GetDealsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithDealsSortDirection(direction SortDirection) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		value := genv2.GetDealsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithDealsIncludeFields(fields ...DealIncludeField) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetDealsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithDealsCustomFields(fields ...string) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithDealsPageSize(limit int) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealsCursor(cursor string) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithArchivedDealsFilterID(id int) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithArchivedDealsIDs(ids ...DealID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithArchivedDealsOwnerID(id UserID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := int(id)
		cfg.params.OwnerId = &value
	})
}

func WithArchivedDealsPersonID(id PersonID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithArchivedDealsOrganizationID(id OrganizationID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := int(id)
		cfg.params.OrgId = &value
	})
}

func WithArchivedDealsPipelineID(id PipelineID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := int(id)
		cfg.params.PipelineId = &value
	})
}

func WithArchivedDealsStageID(id StageID) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := int(id)
		cfg.params.StageId = &value
	})
}

func WithArchivedDealsStatus(statuses ...DealStatus) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		csv := joinCSV(statuses)
		if csv == "" {
			return
		}
		value := genv2.GetArchivedDealsParamsStatus(csv)
		cfg.params.Status = &value
	})
}

func WithArchivedDealsUpdatedSince(t time.Time) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedSince = &value
	})
}

func WithArchivedDealsUpdatedUntil(t time.Time) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedUntil = &value
	})
}

func WithArchivedDealsSortBy(field DealSortField) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := genv2.GetArchivedDealsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithArchivedDealsSortDirection(direction SortDirection) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		value := genv2.GetArchivedDealsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithArchivedDealsIncludeFields(fields ...DealIncludeField) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetArchivedDealsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithArchivedDealsCustomFields(fields ...string) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithArchivedDealsPageSize(limit int) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithArchivedDealsCursor(cursor string) ListArchivedDealsOption {
	return listArchivedDealsOptionFunc(func(cfg *listArchivedDealsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealIncludeFields(fields ...DealIncludeField) GetDealOption {
	return getDealOptionFunc(func(cfg *getDealOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetDealParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithDealCustomFields(fields ...string) GetDealOption {
	return getDealOptionFunc(func(cfg *getDealOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithDealTitle(title string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.title = &title
	})
}

func WithDealValue(value float64) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.value = &value
	})
}

func WithDealCurrency(currency string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.currency = &currency
	})
}

func WithDealOwnerID(id UserID) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.ownerID = &id
	})
}

func WithDealPersonID(id PersonID) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.personID = &id
	})
}

func WithDealOrganizationID(id OrganizationID) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.orgID = &id
	})
}

func WithDealStageID(id StageID) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.stageID = &id
	})
}

func WithDealPipelineID(id PipelineID) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.pipelineID = &id
	})
}

func WithDealStatus(status DealStatus) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.status = &status
	})
}

func WithDealExpectedCloseDate(date string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if date == "" {
			return
		}
		payload.expectedCloseDate = &date
	})
}

func WithDealProbability(probability float64) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.probability = &probability
	})
}

func WithDealLostReason(reason string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if reason == "" {
			return
		}
		payload.lostReason = &reason
	})
}

func WithDealVisibleTo(visibleTo int) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.visibleTo = &visibleTo
	})
}

func WithDealLabelIDs(ids ...int) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if len(ids) == 0 {
			return
		}
		payload.labelIDs = append(payload.labelIDs, ids...)
	})
}

func WithDealCustomFieldsMap(fields map[string]interface{}) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if len(fields) == 0 {
			return
		}
		payload.customFields = fields
	})
}

func WithDealArchived(isArchived bool) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.isArchived = &isArchived
	})
}

func WithDealDeleted(isDeleted bool) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		payload.isDeleted = &isDeleted
	})
}

func WithDealArchiveTime(value string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if value == "" {
			return
		}
		payload.archiveTime = &value
	})
}

func WithDealCloseTime(value string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if value == "" {
			return
		}
		payload.closeTime = &value
	})
}

func WithDealLostTime(value string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if value == "" {
			return
		}
		payload.lostTime = &value
	})
}

func WithDealWonTime(value string) DealOption {
	return dealFieldOption(func(payload *dealPayload) {
		if value == "" {
			return
		}
		payload.wonTime = &value
	})
}

func WithDealSearchFields(fields ...DealSearchField) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchDealsParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithDealSearchExactMatch(enabled bool) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithDealSearchPersonID(id PersonID) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithDealSearchOrganizationID(id OrganizationID) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		value := int(id)
		cfg.params.OrganizationId = &value
	})
}

func WithDealSearchStatus(status DealSearchStatus) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		value := genv2.SearchDealsParamsStatus(status)
		cfg.params.Status = &value
	})
}

func WithDealSearchIncludeFields(fields ...DealSearchIncludeField) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchDealsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithDealSearchPageSize(limit int) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealSearchCursor(cursor string) SearchDealsOption {
	return searchDealsOptionFunc(func(cfg *searchDealsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealFollowersPageSize(limit int) GetDealFollowersOption {
	return getDealFollowersOptionFunc(func(cfg *getDealFollowersOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealFollowersCursor(cursor string) GetDealFollowersOption {
	return getDealFollowersOptionFunc(func(cfg *getDealFollowersOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealFollowersChangelogPageSize(limit int) GetDealFollowersChangelogOption {
	return getDealFollowersChangelogOptionFunc(func(cfg *getDealFollowersChangelogOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealFollowersChangelogCursor(cursor string) GetDealFollowersChangelogOption {
	return getDealFollowersChangelogOptionFunc(func(cfg *getDealFollowersChangelogOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealProductsPageSize(limit int) ListDealProductsOption {
	return listDealProductsOptionFunc(func(cfg *listDealProductsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealProductsCursor(cursor string) ListDealProductsOption {
	return listDealProductsOptionFunc(func(cfg *listDealProductsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealProductsSortBy(field DealProductSortField) ListDealProductsOption {
	return listDealProductsOptionFunc(func(cfg *listDealProductsOptions) {
		value := genv2.GetDealProductsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithDealProductsSortDirection(direction SortDirection) ListDealProductsOption {
	return listDealProductsOptionFunc(func(cfg *listDealProductsOptions) {
		value := genv2.GetDealProductsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithDealsProductsPageSize(limit int) ListDealsProductsOption {
	return listDealsProductsOptionFunc(func(cfg *listDealsProductsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealsProductsCursor(cursor string) ListDealsProductsOption {
	return listDealsProductsOptionFunc(func(cfg *listDealsProductsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealsProductsSortBy(field DealProductSortField) ListDealsProductsOption {
	return listDealsProductsOptionFunc(func(cfg *listDealsProductsOptions) {
		value := genv2.GetDealsProductsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithDealsProductsSortDirection(direction SortDirection) ListDealsProductsOption {
	return listDealsProductsOptionFunc(func(cfg *listDealsProductsOptions) {
		value := genv2.GetDealsProductsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithDealProductProductID(id ProductID) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.productID = &id
	})
}

func WithDealProductVariationID(id ProductVariationID) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.productVariationID = &id
	})
}

func WithDealProductItemPrice(price float64) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.itemPrice = &price
	})
}

func WithDealProductQuantity(quantity float64) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.quantity = &quantity
	})
}

func WithDealProductDiscount(discount float64) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.discount = &discount
	})
}

func WithDealProductDiscountType(discountType DealProductDiscountType) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.discountType = &discountType
	})
}

func WithDealProductComments(comments string) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		if comments == "" {
			return
		}
		payload.comments = &comments
	})
}

func WithDealProductTax(tax float64) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.tax = &tax
	})
}

func WithDealProductTaxMethod(method DealProductTaxMethod) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.taxMethod = &method
	})
}

func WithDealProductIsEnabled(enabled bool) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.isEnabled = &enabled
	})
}

func WithDealProductBillingFrequency(frequency BillingFrequency) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.billingFrequency = &frequency
	})
}

func WithDealProductBillingFrequencyCycles(cycles int) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		payload.billingFrequencyCycles = &cycles
	})
}

func WithDealProductBillingStartDate(date string) DealProductOption {
	return dealProductFieldOption(func(payload *dealProductPayload) {
		if date == "" {
			return
		}
		payload.billingStartDate = &date
	})
}

func WithDealProductAttachmentIDs(ids ...DealProductAttachmentID) DeleteDealProductsOption {
	return deleteDealProductsOptionFunc(func(cfg *deleteDealProductsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithAdditionalDiscountAmount(amount float64) AdditionalDiscountOption {
	return additionalDiscountFieldOption(func(payload *additionalDiscountPayload) {
		payload.amount = &amount
	})
}

func WithAdditionalDiscountType(discountType AdditionalDiscountType) AdditionalDiscountOption {
	return additionalDiscountFieldOption(func(payload *additionalDiscountPayload) {
		payload.discountType = &discountType
	})
}

func WithAdditionalDiscountDescription(description string) AdditionalDiscountOption {
	return additionalDiscountFieldOption(func(payload *additionalDiscountPayload) {
		if description == "" {
			return
		}
		payload.description = &description
	})
}

func WithInstallmentAmount(amount float64) InstallmentOption {
	return installmentFieldOption(func(payload *installmentPayload) {
		payload.amount = &amount
	})
}

func WithInstallmentBillingDate(date string) InstallmentOption {
	return installmentFieldOption(func(payload *installmentPayload) {
		if date == "" {
			return
		}
		payload.billingDate = &date
	})
}

func WithInstallmentDescription(description string) InstallmentOption {
	return installmentFieldOption(func(payload *installmentPayload) {
		if description == "" {
			return
		}
		payload.description = &description
	})
}

func WithInstallmentsPageSize(limit int) ListInstallmentsOption {
	return listInstallmentsOptionFunc(func(cfg *listInstallmentsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithInstallmentsCursor(cursor string) ListInstallmentsOption {
	return listInstallmentsOptionFunc(func(cfg *listInstallmentsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithInstallmentsSortBy(field InstallmentSortField) ListInstallmentsOption {
	return listInstallmentsOptionFunc(func(cfg *listInstallmentsOptions) {
		value := genv2.GetInstallmentsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithInstallmentsSortDirection(direction SortDirection) ListInstallmentsOption {
	return listInstallmentsOptionFunc(func(cfg *listInstallmentsOptions) {
		value := genv2.GetInstallmentsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func NewDealProductInput(opts ...DealProductOption) DealProductInput {
	var cfg addDealProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddDealProduct(&cfg)
	}
	return DealProductInput{payload: cfg.payload}
}

func newGetDealOptions(opts []GetDealOption) getDealOptions {
	var cfg getDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDeal(&cfg)
	}
	return cfg
}

func newListDealsOptions(opts []ListDealsOption) listDealsOptions {
	var cfg listDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListDeals(&cfg)
	}
	return cfg
}

func newListArchivedDealsOptions(opts []ListArchivedDealsOption) listArchivedDealsOptions {
	var cfg listArchivedDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListArchivedDeals(&cfg)
	}
	return cfg
}

func newCreateDealOptions(opts []CreateDealOption) createDealOptions {
	var cfg createDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateDeal(&cfg)
	}
	return cfg
}

func newUpdateDealOptions(opts []UpdateDealOption) updateDealOptions {
	var cfg updateDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateDeal(&cfg)
	}
	return cfg
}

func newDeleteDealOptions(opts []DeleteDealOption) deleteDealOptions {
	var cfg deleteDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDeal(&cfg)
	}
	return cfg
}

func newSearchDealsOptions(opts []SearchDealsOption) searchDealsOptions {
	var cfg searchDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchDeals(&cfg)
	}
	return cfg
}

func newConvertDealOptions(opts []ConvertDealOption) convertDealOptions {
	var cfg convertDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyConvertDeal(&cfg)
	}
	return cfg
}

func newGetDealConversionStatusOptions(opts []GetDealConversionStatusOption) getDealConversionStatusOptions {
	var cfg getDealConversionStatusOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDealConversionStatus(&cfg)
	}
	return cfg
}

func newGetDealFollowersOptions(opts []GetDealFollowersOption) getDealFollowersOptions {
	var cfg getDealFollowersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDealFollowers(&cfg)
	}
	return cfg
}

func newAddDealFollowerOptions(opts []AddDealFollowerOption) addDealFollowerOptions {
	var cfg addDealFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddDealFollower(&cfg)
	}
	return cfg
}

func newDeleteDealFollowerOptions(opts []DeleteDealFollowerOption) deleteDealFollowerOptions {
	var cfg deleteDealFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealFollower(&cfg)
	}
	return cfg
}

func newGetDealFollowersChangelogOptions(opts []GetDealFollowersChangelogOption) getDealFollowersChangelogOptions {
	var cfg getDealFollowersChangelogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDealFollowersChangelog(&cfg)
	}
	return cfg
}

func newListDealProductsOptions(opts []ListDealProductsOption) listDealProductsOptions {
	var cfg listDealProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListDealProducts(&cfg)
	}
	return cfg
}

func newListDealsProductsOptions(opts []ListDealsProductsOption) listDealsProductsOptions {
	var cfg listDealsProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListDealsProducts(&cfg)
	}
	return cfg
}

func newAddDealProductOptions(opts []AddDealProductOption) addDealProductOptions {
	var cfg addDealProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddDealProduct(&cfg)
	}
	return cfg
}

func newAddManyDealProductsOptions(opts []AddManyDealProductsOption) addManyDealProductsOptions {
	var cfg addManyDealProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddManyDealProducts(&cfg)
	}
	return cfg
}

func newUpdateDealProductOptions(opts []UpdateDealProductOption) updateDealProductOptions {
	var cfg updateDealProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateDealProduct(&cfg)
	}
	return cfg
}

func newDeleteDealProductOptions(opts []DeleteDealProductOption) deleteDealProductOptions {
	var cfg deleteDealProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealProduct(&cfg)
	}
	return cfg
}

func newDeleteDealProductsOptions(opts []DeleteDealProductsOption) deleteDealProductsOptions {
	var cfg deleteDealProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealProducts(&cfg)
	}
	return cfg
}

func newListAdditionalDiscountsOptions(opts []ListAdditionalDiscountsOption) listAdditionalDiscountsOptions {
	var cfg listAdditionalDiscountsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListAdditionalDiscounts(&cfg)
	}
	return cfg
}

func newAddAdditionalDiscountOptions(opts []AddAdditionalDiscountOption) addAdditionalDiscountOptions {
	var cfg addAdditionalDiscountOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddAdditionalDiscount(&cfg)
	}
	return cfg
}

func newUpdateAdditionalDiscountOptions(opts []UpdateAdditionalDiscountOption) updateAdditionalDiscountOptions {
	var cfg updateAdditionalDiscountOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateAdditionalDiscount(&cfg)
	}
	return cfg
}

func newDeleteAdditionalDiscountOptions(opts []DeleteAdditionalDiscountOption) deleteAdditionalDiscountOptions {
	var cfg deleteAdditionalDiscountOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteAdditionalDiscount(&cfg)
	}
	return cfg
}

func newListInstallmentsOptions(opts []ListInstallmentsOption) listInstallmentsOptions {
	var cfg listInstallmentsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListInstallments(&cfg)
	}
	return cfg
}

func newAddInstallmentOptions(opts []AddInstallmentOption) addInstallmentOptions {
	var cfg addInstallmentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddInstallment(&cfg)
	}
	return cfg
}

func newUpdateInstallmentOptions(opts []UpdateInstallmentOption) updateInstallmentOptions {
	var cfg updateInstallmentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateInstallment(&cfg)
	}
	return cfg
}

func newDeleteInstallmentOptions(opts []DeleteInstallmentOption) deleteInstallmentOptions {
	var cfg deleteInstallmentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteInstallment(&cfg)
	}
	return cfg
}

func (s *DealsService) Get(ctx context.Context, id DealID, opts ...GetDealOption) (*Deal, error) {
	cfg := newGetDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetDealWithResponse(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) List(ctx context.Context, opts ...ListDealsOption) ([]Deal, *string, error) {
	cfg := newListDealsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListPager(opts ...ListDealsOption) *pipedrive.CursorPager[Deal] {
	cfg := newListDealsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Deal, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEach(ctx context.Context, fn func(Deal) error, opts ...ListDealsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *DealsService) ListArchived(ctx context.Context, opts ...ListArchivedDealsOption) ([]Deal, *string, error) {
	cfg := newListArchivedDealsOptions(opts)
	return s.listArchived(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListArchivedPager(opts ...ListArchivedDealsOption) *pipedrive.CursorPager[Deal] {
	cfg := newListArchivedDealsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Deal, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listArchived(ctx, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachArchived(ctx context.Context, fn func(Deal) error, opts ...ListArchivedDealsOption) error {
	return s.ListArchivedPager(opts...).ForEach(ctx, fn)
}

func (s *DealsService) Create(ctx context.Context, opts ...CreateDealOption) (*Deal, error) {
	cfg := newCreateDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddDealWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Update(ctx context.Context, id DealID, opts ...UpdateDealOption) (*Deal, error) {
	cfg := newUpdateDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateDealWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Delete(ctx context.Context, id DealID, opts ...DeleteDealOption) (*DealDeleteResult, error) {
	cfg := newDeleteDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteDealWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal delete data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Search(ctx context.Context, term string, opts ...SearchDealsOption) (*DealSearchResults, *string, error) {
	cfg := newSearchDealsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchDealsWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *DealSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing deal search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) ConvertToLead(ctx context.Context, id DealID, opts ...ConvertDealOption) (*DealConversionJob, error) {
	cfg := newConvertDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.ConvertDealToLeadWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealConversionJob `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal conversion data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) ConversionStatus(ctx context.Context, id DealID, conversionID ConversionID, opts ...GetDealConversionStatusOption) (*DealConversionStatus, error) {
	cfg := newGetDealConversionStatusOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	convUUID, err := parseUUID(string(conversionID), "conversion id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.GetDealConversionStatusWithResponse(ctx, int(id), convUUID, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealConversionStatus `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal conversion status data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) ListFollowers(ctx context.Context, id DealID, opts ...GetDealFollowersOption) ([]Follower, *string, error) {
	cfg := newGetDealFollowersOptions(opts)
	return s.listFollowers(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListFollowersPager(id DealID, opts ...GetDealFollowersOption) *pipedrive.CursorPager[Follower] {
	cfg := newGetDealFollowersOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Follower, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listFollowers(ctx, id, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachFollowers(ctx context.Context, id DealID, fn func(Follower) error, opts ...GetDealFollowersOption) error {
	return s.ListFollowersPager(id, opts...).ForEach(ctx, fn)
}

func (s *DealsService) AddFollower(ctx context.Context, id DealID, userID UserID, opts ...AddDealFollowerOption) (*Follower, error) {
	cfg := newAddDealFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(map[string]interface{}{
		"user_id": int(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddDealFollowerWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Follower `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing follower data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteFollower(ctx context.Context, id DealID, followerID UserID, opts ...DeleteDealFollowerOption) (*FollowerDeleteResult, error) {
	cfg := newDeleteDealFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteDealFollowerWithResponse(ctx, int(id), int(followerID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *FollowerDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete follower data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) FollowersChangelog(ctx context.Context, id DealID, opts ...GetDealFollowersChangelogOption) ([]FollowerChangelog, *string, error) {
	cfg := newGetDealFollowersChangelogOptions(opts)
	return s.followersChangelog(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *DealsService) FollowersChangelogPager(id DealID, opts ...GetDealFollowersChangelogOption) *pipedrive.CursorPager[FollowerChangelog] {
	cfg := newGetDealFollowersChangelogOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]FollowerChangelog, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.followersChangelog(ctx, id, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachFollowersChangelog(ctx context.Context, id DealID, fn func(FollowerChangelog) error, opts ...GetDealFollowersChangelogOption) error {
	return s.FollowersChangelogPager(id, opts...).ForEach(ctx, fn)
}

func (s *DealsService) ListProducts(ctx context.Context, id DealID, opts ...ListDealProductsOption) ([]DealProduct, *string, error) {
	cfg := newListDealProductsOptions(opts)
	return s.listDealProducts(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListProductsPager(id DealID, opts ...ListDealProductsOption) *pipedrive.CursorPager[DealProduct] {
	cfg := newListDealProductsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]DealProduct, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listDealProducts(ctx, id, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachProducts(ctx context.Context, id DealID, fn func(DealProduct) error, opts ...ListDealProductsOption) error {
	return s.ListProductsPager(id, opts...).ForEach(ctx, fn)
}

func (s *DealsService) ListProductsAcrossDeals(ctx context.Context, dealIDs []DealID, opts ...ListDealsProductsOption) ([]DealProduct, *string, error) {
	if len(dealIDs) == 0 {
		return nil, nil, fmt.Errorf("deal IDs are required")
	}
	cfg := newListDealsProductsOptions(opts)
	cfg.params.DealIds = make([]int, 0, len(dealIDs))
	for _, id := range dealIDs {
		cfg.params.DealIds = append(cfg.params.DealIds, int(id))
	}
	return s.listDealsProducts(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListProductsAcrossDealsPager(dealIDs []DealID, opts ...ListDealsProductsOption) *pipedrive.CursorPager[DealProduct] {
	cfg := newListDealsProductsOptions(opts)
	cfg.params.DealIds = make([]int, 0, len(dealIDs))
	for _, id := range dealIDs {
		cfg.params.DealIds = append(cfg.params.DealIds, int(id))
	}
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]DealProduct, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listDealsProducts(ctx, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachProductsAcrossDeals(ctx context.Context, dealIDs []DealID, fn func(DealProduct) error, opts ...ListDealsProductsOption) error {
	return s.ListProductsAcrossDealsPager(dealIDs, opts...).ForEach(ctx, fn)
}

func (s *DealsService) AddProduct(ctx context.Context, id DealID, opts ...AddDealProductOption) (*DealProduct, error) {
	cfg := newAddDealProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddDealProductWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealProduct `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal product data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) AddProducts(ctx context.Context, id DealID, products []DealProductInput, opts ...AddManyDealProductsOption) ([]DealProduct, error) {
	cfg := newAddManyDealProductsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	items := make([]map[string]interface{}, 0, len(products))
	for _, product := range products {
		items = append(items, product.payload.toMap())
	}
	body, err := json.Marshal(map[string]interface{}{
		"data": items,
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddManyDealProductsWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data []DealProduct `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *DealsService) UpdateProduct(ctx context.Context, id DealID, attachmentID DealProductAttachmentID, opts ...UpdateDealProductOption) (*DealProduct, error) {
	cfg := newUpdateDealProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateDealProductWithBodyWithResponse(ctx, int(id), int(attachmentID), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealProduct `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal product data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteProduct(ctx context.Context, id DealID, attachmentID DealProductAttachmentID, opts ...DeleteDealProductOption) (*DealProductDeleteResult, error) {
	cfg := newDeleteDealProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteDealProductWithResponse(ctx, int(id), int(attachmentID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *DealProductDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal product delete data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteProducts(ctx context.Context, id DealID, opts ...DeleteDealProductsOption) (*DealProductsDeleteResult, error) {
	cfg := newDeleteDealProductsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteManyDealProductsWithResponse(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *struct {
			IDs []DealProductAttachmentID `json:"ids"`
		} `json:"data"`
		AdditionalData *struct {
			MoreItemsInCollection *bool `json:"more_items_in_collection"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal products delete data in response")
	}

	result := DealProductsDeleteResult{IDs: payload.Data.IDs}
	if payload.AdditionalData != nil {
		result.MoreItemsInCollection = payload.AdditionalData.MoreItemsInCollection
	}
	return &result, nil
}

func (s *DealsService) ListAdditionalDiscounts(ctx context.Context, id DealID, opts ...ListAdditionalDiscountsOption) ([]AdditionalDiscount, error) {
	cfg := newListAdditionalDiscountsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetAdditionalDiscountsWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data []AdditionalDiscount `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *DealsService) AddAdditionalDiscount(ctx context.Context, id DealID, opts ...AddAdditionalDiscountOption) (*AdditionalDiscount, error) {
	cfg := newAddAdditionalDiscountOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.PostAdditionalDiscountWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *AdditionalDiscount `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing additional discount data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) UpdateAdditionalDiscount(ctx context.Context, id DealID, discountID AdditionalDiscountID, opts ...UpdateAdditionalDiscountOption) (*AdditionalDiscount, error) {
	cfg := newUpdateAdditionalDiscountOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	discountUUID, err := parseUUID(string(discountID), "discount id")
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateAdditionalDiscountWithBodyWithResponse(ctx, int(id), discountUUID, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *AdditionalDiscount `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing additional discount data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteAdditionalDiscount(ctx context.Context, id DealID, discountID AdditionalDiscountID, opts ...DeleteAdditionalDiscountOption) (*AdditionalDiscountDeleteResult, error) {
	cfg := newDeleteAdditionalDiscountOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	discountUUID, err := parseUUID(string(discountID), "discount id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.DeleteAdditionalDiscount(ctx, int(id), discountUUID, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errorFromResponse(resp, body)
	}

	var payload struct {
		Data *AdditionalDiscountDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing additional discount delete data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) ListInstallments(ctx context.Context, dealIDs []DealID, opts ...ListInstallmentsOption) ([]Installment, *string, error) {
	if len(dealIDs) == 0 {
		return nil, nil, fmt.Errorf("deal IDs are required")
	}
	cfg := newListInstallmentsOptions(opts)
	cfg.params.DealIds = make([]int, 0, len(dealIDs))
	for _, id := range dealIDs {
		cfg.params.DealIds = append(cfg.params.DealIds, int(id))
	}
	return s.listInstallments(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListInstallmentsPager(dealIDs []DealID, opts ...ListInstallmentsOption) *pipedrive.CursorPager[Installment] {
	cfg := newListInstallmentsOptions(opts)
	cfg.params.DealIds = make([]int, 0, len(dealIDs))
	for _, id := range dealIDs {
		cfg.params.DealIds = append(cfg.params.DealIds, int(id))
	}
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Installment, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listInstallments(ctx, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEachInstallments(ctx context.Context, dealIDs []DealID, fn func(Installment) error, opts ...ListInstallmentsOption) error {
	return s.ListInstallmentsPager(dealIDs, opts...).ForEach(ctx, fn)
}

func (s *DealsService) AddInstallment(ctx context.Context, id DealID, opts ...AddInstallmentOption) (*Installment, error) {
	cfg := newAddInstallmentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.PostInstallmentWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Installment `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing installment data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) UpdateInstallment(ctx context.Context, id DealID, installmentID InstallmentID, opts ...UpdateInstallmentOption) (*Installment, error) {
	cfg := newUpdateInstallmentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateInstallmentWithBodyWithResponse(ctx, int(id), int(installmentID), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Installment `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing installment data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteInstallment(ctx context.Context, id DealID, installmentID InstallmentID, opts ...DeleteInstallmentOption) (*InstallmentDeleteResult, error) {
	cfg := newDeleteInstallmentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteInstallmentWithResponse(ctx, int(id), int(installmentID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *InstallmentDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing installment delete data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) list(ctx context.Context, params genv2.GetDealsParams, requestOptions []pipedrive.RequestOption) ([]Deal, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Deal `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) listArchived(ctx context.Context, params genv2.GetArchivedDealsParams, requestOptions []pipedrive.RequestOption) ([]Deal, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetArchivedDealsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Deal `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) listFollowers(ctx context.Context, id DealID, params genv2.GetDealFollowersParams, requestOptions []pipedrive.RequestOption) ([]Follower, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealFollowersWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Follower `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) followersChangelog(ctx context.Context, id DealID, params genv2.GetDealFollowersChangelogParams, requestOptions []pipedrive.RequestOption) ([]FollowerChangelog, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealFollowersChangelogWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []FollowerChangelog `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) listDealProducts(ctx context.Context, id DealID, params genv2.GetDealProductsParams, requestOptions []pipedrive.RequestOption) ([]DealProduct, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealProductsWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []DealProduct `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) listDealsProducts(ctx context.Context, params genv2.GetDealsProductsParams, requestOptions []pipedrive.RequestOption) ([]DealProduct, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealsProductsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []DealProduct `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *DealsService) listInstallments(ctx context.Context, params genv2.GetInstallmentsParams, requestOptions []pipedrive.RequestOption) ([]Installment, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetInstallmentsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Installment `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (p dealPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.title != nil {
		body["title"] = *p.title
	}
	if p.value != nil {
		body["value"] = *p.value
	}
	if p.currency != nil {
		body["currency"] = *p.currency
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.personID != nil {
		body["person_id"] = int(*p.personID)
	}
	if p.orgID != nil {
		body["org_id"] = int(*p.orgID)
	}
	if p.stageID != nil {
		body["stage_id"] = int(*p.stageID)
	}
	if p.pipelineID != nil {
		body["pipeline_id"] = int(*p.pipelineID)
	}
	if p.status != nil {
		body["status"] = string(*p.status)
	}
	if p.expectedCloseDate != nil {
		body["expected_close_date"] = *p.expectedCloseDate
	}
	if p.probability != nil {
		body["probability"] = *p.probability
	}
	if p.lostReason != nil {
		body["lost_reason"] = *p.lostReason
	}
	if p.visibleTo != nil {
		body["visible_to"] = *p.visibleTo
	}
	if len(p.labelIDs) > 0 {
		body["label_ids"] = p.labelIDs
	}
	if p.customFields != nil {
		body["custom_fields"] = p.customFields
	}
	if p.isArchived != nil {
		body["is_archived"] = *p.isArchived
	}
	if p.isDeleted != nil {
		body["is_deleted"] = *p.isDeleted
	}
	if p.archiveTime != nil {
		body["archive_time"] = *p.archiveTime
	}
	if p.closeTime != nil {
		body["close_time"] = *p.closeTime
	}
	if p.lostTime != nil {
		body["lost_time"] = *p.lostTime
	}
	if p.wonTime != nil {
		body["won_time"] = *p.wonTime
	}
	return body
}

func (p dealProductPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.productID != nil {
		body["product_id"] = int(*p.productID)
	}
	if p.productVariationID != nil {
		body["product_variation_id"] = int(*p.productVariationID)
	}
	if p.itemPrice != nil {
		body["item_price"] = *p.itemPrice
	}
	if p.quantity != nil {
		body["quantity"] = *p.quantity
	}
	if p.discount != nil {
		body["discount"] = *p.discount
	}
	if p.discountType != nil {
		body["discount_type"] = string(*p.discountType)
	}
	if p.comments != nil {
		body["comments"] = *p.comments
	}
	if p.tax != nil {
		body["tax"] = *p.tax
	}
	if p.taxMethod != nil {
		body["tax_method"] = string(*p.taxMethod)
	}
	if p.isEnabled != nil {
		body["is_enabled"] = *p.isEnabled
	}
	if p.billingFrequency != nil {
		body["billing_frequency"] = string(*p.billingFrequency)
	}
	if p.billingFrequencyCycles != nil {
		body["billing_frequency_cycles"] = *p.billingFrequencyCycles
	}
	if p.billingStartDate != nil {
		body["billing_start_date"] = *p.billingStartDate
	}
	return body
}

func (p additionalDiscountPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.amount != nil {
		body["amount"] = *p.amount
	}
	if p.discountType != nil {
		body["type"] = string(*p.discountType)
	}
	if p.description != nil {
		body["description"] = *p.description
	}
	return body
}

func (p installmentPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.amount != nil {
		body["amount"] = *p.amount
	}
	if p.billingDate != nil {
		body["billing_date"] = *p.billingDate
	}
	if p.description != nil {
		body["description"] = *p.description
	}
	return body
}
