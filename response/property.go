package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

// SimpleDeveloper contains only the essential developer information
type SimpleDeveloper struct {
	Name             string `json:"name"`
	DeveloperLogo    string `json:"developer_logo"`
	DeveloperAddress string `json:"developer_address"`
}

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
	Developer      *SimpleDeveloper           `json:"developer,omitempty"`
}

func GetPropertyFromEnt(property *ent.Property) *Property {
	var developer *SimpleDeveloper
	if property.Edges.Developer != nil {
		developer = &SimpleDeveloper{
			Name:             property.Edges.Developer.Name,
			DeveloperLogo:    property.Edges.Developer.MediaContent.DeveloperLogo,
			DeveloperAddress: property.Edges.Project.WebCards.About.ContactDetails.ProjectAddress,
		}
	}

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
		Developer:      developer,
	}
}

type AddPropertyResponse struct {
	PropertyID string `json:"property_id"`
}

type PropertyListResponse struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	PossessionStatus string   `json:"possession_status"`
	BuiltUpArea      string   `json:"built_up_area"`
	Facing           string   `json:"facing"`
	FloorNumber      string   `json:"floor_number"`
	Images           []string `json:"images"`
	Location         string   `json:"location"`
	DeveloperName    string   `json:"developer_name"`
}

func GetPropertyListResponse(property *ent.Property, developerName string, location string) *PropertyListResponse {
	return &PropertyListResponse{
		ID:               property.ID,
		Name:             property.Name,
		PossessionStatus: property.WebCards.PropertyDetails.PossessionStatus.Value,
		BuiltUpArea:      property.WebCards.PropertyDetails.BuiltUpArea.Value,
		Facing:           property.WebCards.PropertyDetails.Facing.Value,
		FloorNumber:      property.WebCards.PropertyDetails.FloorNumber.Value,
		Images:           property.PropertyImages,
		Location:         location,
		DeveloperName:    developerName,
	}
}
