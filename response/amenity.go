package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

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

type StaticSiteDataResponse struct {
	ID                      string               `json:"id"`
	PropertyTypes           schema.PropertyTypes `json:"property_types"`
	CategoriesWithAmenities struct {
		Categories map[string][]Amenity `json:"categories"`
	} `json:"categories_with_amenities"`
	Testimonials []schema.Testimonials `json:"testimonials"`
	IsActive     bool                  `json:"is_active"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}

func GetStaticSiteDataFromEnt(data *ent.StaticSiteData) *StaticSiteDataResponse {
	// Convert categories to our response format
	categories := make(map[string][]Amenity)
	for category, amenities := range data.CategoriesWithAmenities.Categories {
		var responseAmenities []Amenity
		for _, amenity := range amenities {
			responseAmenities = append(responseAmenities, Amenity{
				Icon:  amenity.Icon,
				Value: amenity.Value,
			})
		}
		categories[category] = responseAmenities
	}

	return &StaticSiteDataResponse{
		ID:            data.ID,
		PropertyTypes: data.PropertyTypes,
		CategoriesWithAmenities: struct {
			Categories map[string][]Amenity `json:"categories"`
		}{
			Categories: categories,
		},
		Testimonials: data.Testimonials,
		IsActive:     data.IsActive,
		CreatedAt:    data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
