package response

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int `json:"current_page"`
		PageSize    int `json:"page_size"`
		TotalItems  int `json:"total_items"`
		TotalPages  int `json:"total_pages"`
	} `json:"pagination"`
}

func NewPaginatedResponse(data interface{}, page, pageSize, totalItems int) *PaginatedResponse {
	resp := &PaginatedResponse{
		Data: data,
	}
	resp.Pagination.CurrentPage = page
	resp.Pagination.PageSize = pageSize
	resp.Pagination.TotalItems = totalItems
	resp.Pagination.TotalPages = (totalItems + pageSize - 1) / pageSize
	return resp
}
