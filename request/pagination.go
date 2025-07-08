package request

type GetAllAPIRequest struct {
	Page     int                    `json:"page" query:"page"`
	PageSize int                    `json:"page_size" query:"page_size"`
	Filters  map[string]interface{} `json:"filters,omitempty" query:"filters"`
}

func (p *GetAllAPIRequest) Validate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10 // Default page size
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Maximum page size
	}
}

func (p *GetAllAPIRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *GetAllAPIRequest) GetLimit() int {
	return p.PageSize
}
