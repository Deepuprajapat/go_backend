package request

type CheckURLExistsRequest struct {
	URL string `json:"url" validate:"required"`
}