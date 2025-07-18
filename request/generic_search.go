package request

type CustomSearchPage struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Slug        string                 `json:"slug,omitempty"`
	Filters     map[string]interface{} `json:"filters"`
	MetaInfo    *MetaInfo              `json:"meta_info"`
	SearchTerm  string                 `json:"search_term"`
}

type MetaInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}
