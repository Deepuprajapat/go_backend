package domain

import "github.com/VI-IM/im_backend_go/ent/schema"

type Property struct {
	PropertyID       string
	Name             string
	PropertyType     string
	PropertyImages   []string
	WebCards         schema.WebCards
	PricingInfo      schema.PropertyPricingInfo
	PropertyReraInfo schema.PropertyReraInfo
	MetaInfo         schema.PropertyMetaInfo
	IsFeatured       bool
	IsDeleted        bool
	DeveloperID      string
	LocationID       string
	ProjectID        string
	CreatedByUserID  *string
}
