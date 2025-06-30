package request

type CreateAmenityRequest struct {
	Category string `json:"category" validate:"required"`
	Icon     string `json:"icon" validate:"required"`
	Value    string `json:"value" validate:"required"`
}

type UpdateAmenityRequest struct {
	Category string `json:"category"`
	Icon     string `json:"icon"`
	Value    string `json:"value"`
}
