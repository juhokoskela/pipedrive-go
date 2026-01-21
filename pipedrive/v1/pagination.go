package v1

type Pagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
	NextStart             int  `json:"next_start,omitempty"`
}

type CollectionPagination struct {
	NextCursor *string `json:"next_cursor,omitempty"`
}
