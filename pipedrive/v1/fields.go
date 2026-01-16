package v1

type FieldType string

type Field struct {
	ID                  *FieldID                 `json:"id,omitempty"`
	Name                string                   `json:"name,omitempty"`
	Key                 string                   `json:"key,omitempty"`
	FieldType           FieldType                `json:"field_type,omitempty"`
	Active              bool                     `json:"active_flag,omitempty"`
	AddTime             *DateTime                `json:"add_time,omitempty"`
	AddVisible          bool                     `json:"add_visible_flag,omitempty"`
	BulkEditAllowed     bool                     `json:"bulk_edit_allowed,omitempty"`
	CreatedByUserID     *UserID                  `json:"created_by_user_id,omitempty"`
	DetailsVisible      bool                     `json:"details_visible_flag,omitempty"`
	Edit                bool                     `json:"edit_flag,omitempty"`
	FilteringAllowed    bool                     `json:"filtering_allowed,omitempty"`
	Important           bool                     `json:"important_flag,omitempty"`
	IndexVisible        bool                     `json:"index_visible_flag,omitempty"`
	IsSubfield          bool                     `json:"is_subfield,omitempty"`
	LastUpdatedByUserID *UserID                  `json:"last_updated_by_user_id,omitempty"`
	Mandatory           bool                     `json:"mandatory_flag,omitempty"`
	Options             []map[string]interface{} `json:"options,omitempty"`
	OptionsDeleted      []map[string]interface{} `json:"options_deleted,omitempty"`
	OrderNr             int                      `json:"order_nr,omitempty"`
	Searchable          bool                     `json:"searchable_flag,omitempty"`
	Sortable            bool                     `json:"sortable_flag,omitempty"`
	Subfields           []map[string]interface{} `json:"subfields,omitempty"`
	UpdateTime          *DateTime                `json:"update_time,omitempty"`
}

type FieldPagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
}
