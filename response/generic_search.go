package response

type CustomSearchPage struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Slug        string                 `json:"slug,omitempty"`
	Projects    []*ProjectListResponse `json:"projects,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	SearchTerm  string                 `json:"search_term,omitempty"`
	MetaInfo    *MetaInfo              `json:"meta_info,omitempty"`
}

type MetaInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

type Link struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
