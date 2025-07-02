package request

type Amenities struct {
	Icon  string `json:"icon"`
	Value string `json:"value"`
}

type CreateAmenityRequest struct {
	Category map[string][]struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"category" validate:"required"`
}

type UpdateAmenityRequest struct {
	Category string `json:"category"`
	Icon     string `json:"icon"`
	Value    string `json:"value"`
}
