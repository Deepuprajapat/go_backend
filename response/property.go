package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

type Property struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	PropertyImages []string                   `json:"property_images"`
	WebCards       schema.WebCards            `json:"web_cards"`
	PricingInfo    schema.PropertyPricingInfo `json:"pricing_info"`
	PropertyRera   schema.PropertyReraInfo    `json:"property_rera_info"`
	MetaInfo       schema.PropertyMetaInfo    `json:"meta_info"`
	DeveloperID    string                     `json:"developer_id"`
	LocationID     string                     `json:"location_id"`
	ProjectID      string                     `json:"project_id,omitempty"`
}

func GetPropertyFromEnt(property *ent.Property) *Property {
	return &Property{
		ID:             property.ID,
		Name:           property.Name,
		PropertyImages: property.PropertyImages,
		WebCards:       property.WebCards,
		PricingInfo:    property.PricingInfo,
		PropertyRera:   property.PropertyReraInfo,
		MetaInfo:       property.MetaInfo,
		DeveloperID:    property.DeveloperID,
		LocationID:     property.LocationID,
		ProjectID:      property.ProjectID,
	}
}
