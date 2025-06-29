package request

type PaginationRequest struct {
	Page     int `json:"page" query:"page"`           // Current page number (1-based)
	PageSize int `json:"page_size" query:"page_size"` // Number of items per page
}

func (p *PaginationRequest) Validate() {
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

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}
