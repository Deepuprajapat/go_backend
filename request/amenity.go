package request

import (
	"github.com/VI-IM/im_backend_go/ent/schema"
)

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

type AddAmenitiesToCategoryRequest struct {
	Category string      `json:"category" validate:"required"`
	Items    []Amenities `json:"items" validate:"required"`
}

type DeleteAmenitiesFromCategoryRequest struct {
	Category string   `json:"category" validate:"required"`
	Values   []string `json:"values" validate:"required"`
}

type DeleteCategoryWithAmenitiesRequest struct {
	Category string `json:"category" validate:"required"`
}

// Request type for updating static site data
type UpdateStaticSiteDataRequest struct {
	PropertyTypes           *schema.PropertyTypes `json:"property_types,omitempty"`
	CategoriesWithAmenities struct {
		Categories map[string][]struct {
			Icon  string `json:"icon"`
			Value string `json:"value"`
		} `json:"categories,omitempty"`
	} `json:"categories_with_amenities,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
}

type AddAmenityToCategoryRequest struct {
	CategoryName string `json:"category_name"`
	Amenities    []struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"amenities"`
}

type DeleteAmenityFromCategoryRequest struct {
	CategoryName string `json:"category_name"`
	Amenities    []struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"amenities"`
}
