package response

type AmenityResponse struct {
	Categories map[string][]Amenity `json:"categories"`
}

type Amenity struct {
	Icon  string `json:"icon"`
	Value string `json:"value"`
}

type SingleAmenityResponse struct {
	Category string `json:"category"`
	Icon     string `json:"icon"`
	Value    string `json:"value"`
}
