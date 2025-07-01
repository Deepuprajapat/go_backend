package request

import (
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type AddProjectRequest struct {
	ProjectName string `json:"project_name" validate:"required"`
	ProjectURL  string `json:"project_url" validate:"required"`
	ProjectType string `json:"project_type" validate:"required"`
	Locality    string `json:"locality" validate:"required"`
	ProjectCity string `json:"project_city" validate:"required"`
	DeveloperID string `json:"developer_id" validate:"required"`
}

type UpdateProjectRequest struct {
	ProjectID    string                 `json:"project_id" validate:"required"`
	ProjectName  string                 `json:"project_name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Status       enums.ProjectStatus    `json:"status,omitempty"`
	PriceUnit    string                 `json:"price_unit,omitempty"`
	TimelineInfo schema.TimelineInfo    `json:"timeline_info,omitempty"`
	MetaInfo     schema.SEOMeta         `json:"meta_info,omitempty"`
	WebCards     schema.ProjectWebCards `json:"web_cards,omitempty"`
	LocationInfo schema.LocationInfo    `json:"location_info,omitempty"`
	IsFeatured   bool                   `json:"is_featured,omitempty"`
	IsPremium    bool                   `json:"is_premium,omitempty"`
	IsPriority   bool                   `json:"is_priority,omitempty"`
	IsDeleted    bool                   `json:"is_deleted,omitempty"`
}

type UpdatePropertyRequest struct {
	PropertyID       string                     `json:"property_id"`
	Name             string                     `json:"name"`
	PropertyImages   []string                   `json:"property_images"`
	WebCards         schema.WebCards            `json:"web_cards"`
	PricingInfo      schema.PropertyPricingInfo `json:"pricing_info"`
	PropertyReraInfo schema.PropertyReraInfo    `json:"property_rera_info"`
	MetaInfo         schema.PropertyMetaInfo    `json:"meta_info"`
	IsFeatured       bool                       `json:"is_featured"`
	IsDeleted        bool                       `json:"is_deleted"`
	DeveloperID      string                     `json:"developer_id"`
	LocationID       string                     `json:"location_id"`
	ProjectID        string                     `json:"project_id"`
}

type AddPropertyRequest struct {
	ProjectID      string `json:"project_id"`
	Name           string `json:"name"`
	PropertyType   string `json:"property_type"`
	AgeOfProperty  string `json:"age_of_property"`
	FloorNumber    string `json:"floor_number"`
	Facing         string `json:"facing"`
	Furnishing     string `json:"furnishing"`
	BalconyCount   string `json:"balcony_count"`
	BedroomsCount  string `json:"bedrooms_count"`
	CoveredParking string `json:"covered_parking"`
}

type ProjectFilterRequest struct {
	Configurations []string `json:"configurations"` // List of configurations like "1BHK", "2BHK", etc.
	IsPremium      bool     `json:"is_premium"`
	IsPriority     bool     `json:"is_priority"`
	IsFeatured     bool     `json:"is_featured"`
	LocationID     string   `json:"location_id"`
	DeveloperID    string   `json:"developer_id"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
}

func (r *ProjectFilterRequest) ToMap() map[string]interface{} {
	filters := make(map[string]interface{})
	if len(r.Configurations) > 0 {
		filters["configurations"] = r.Configurations
	}
	if r.IsPremium {
		filters["is_premium"] = r.IsPremium
	}
	if r.IsPriority {
		filters["is_priority"] = r.IsPriority
	}
	if r.IsFeatured {
		filters["is_featured"] = r.IsFeatured
	}
	if r.LocationID != "" {
		filters["location_id"] = r.LocationID
	}
	if r.DeveloperID != "" {
		filters["developer_id"] = r.DeveloperID
	}
	if r.Name != "" {
		filters["name"] = r.Name
	}
	if r.Type != "" {
		filters["type"] = r.Type
	}
	return filters
}
