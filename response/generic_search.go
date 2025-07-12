package response

type GenericSearchData struct {
	Index        int               `json:"index,omitempty"`
	CanonicalURL string            `json:"canonical_url"`
	SearchTerm   string            `json:"title"`
	Filters      map[string]string `json:"filters"`
}
