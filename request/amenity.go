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

// New request type for adding amenities to a category
type AddAmenitiesToCategoryRequest struct {
	Category string      `json:"category" validate:"required"`
	Items    []Amenities `json:"items" validate:"required"`
}

// Request type for deleting amenities from a category
type DeleteAmenitiesFromCategoryRequest struct {
	Category string   `json:"category" validate:"required"`
	Values   []string `json:"values" validate:"required"` // List of amenity values to delete
}

// Request type for deleting a category with all its amenities
type DeleteCategoryRequest struct {
	Category string `json:"category" validate:"required"`
}
