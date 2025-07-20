package response

type CheckURLExistsResponse struct {
	Exists     bool   `json:"exists"`
	EntityType string `json:"entity_type,omitempty"` // "blog", "project", "custom_search_page"
	EntityID   string `json:"entity_id,omitempty"`
}

type URLExistsResult struct {
	Exists     bool
	EntityType string
	EntityID   string
}