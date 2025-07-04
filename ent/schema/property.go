package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Property struct {
	Base
}

func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("name"),
		field.String("property_type").Optional(),
		field.JSON("property_images", []string{}).Optional(), // 0 index logo image
		field.JSON("web_cards", WebCards{}),
		field.JSON("pricing_info", PropertyPricingInfo{}),
		field.JSON("meta_info", PropertyMetaInfo{}).Optional(),
		field.JSON("property_rera_info", PropertyReraInfo{}),
		field.JSON("search_context", []string{}).Optional(),
		field.Bool("is_deleted").Default(false),
		field.Bool("is_featured").Default(false),
		field.String("project_id").Optional(),
		field.String("developer_id").Optional(),
		field.String("location_id").Optional(),
	}
}

func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("properties").
			Unique().
			Field("project_id"),
		edge.To("leads", Leads.Type),
		edge.To("developer", Developer.Type).
			Unique().
			Field("developer_id"),
		edge.To("location", Location.Type).
			Unique().
			Field("location_id"),
	}
}

type PropertyReraInfo struct {
	ReraNumber string `json:"rera_number,omitempty"`
}

type WebCards struct {
	PropertyDetails PropertyDetails `json:"property_details,omitempty"`
	WhyChooseUs     struct {
		ImageUrls []string `json:"image_urls,omitempty"`
		USP_List  []string `json:"usp_list,omitempty"`
	} `json:"why_choose_us,omitempty"`
	PropertyFloorPlan PropertyFloorPlan `json:"property_floor_plan,omitempty"`
	KnowAbout         struct {
		Description string `json:"description,omitempty"`
	} `json:"know_about,omitempty"`
	VideoPresentation struct {
		Title    string `json:"title,omitempty"`
		VideoUrl string `json:"video_url,omitempty"`
	} `json:"video_presentation,omitempty"`
	LocationMap struct {
		Description   string `json:"description,omitempty"`
		GoogleMapLink string `json:"google_map_link,omitempty"`
	} `json:"location_map,omitempty"`
}

type PropertyFloorPlan struct {
	Title string              `json:"title,omitempty"`
	Plans []map[string]string `json:"plans,omitempty"`
}

// area details
type PropertyAreaDetails struct {
	CarpetArea       string `json:"carpet_area,omitempty"`         // in sq ft
	BuiltUpArea      string `json:"built_up_area,omitempty"`       // in sq ft
	SuperBuiltUpArea string `json:"super_built_up_area,omitempty"` // in sq ft
}

// pricing information
type PropertyPricingInfo struct {
	Price string `json:"price,omitempty"` // selling price
}

// property details
type PropertyDetails struct {
	BuiltUpArea struct {
		Value string `json:"value,omitempty"`
	} `json:"built_up_area,omitempty"`
	Sizes struct {
		Value string `json:"value,omitempty"`
	} `json:"size,omitempty"`
	FloorNumber struct {
		Value string `json:"value,omitempty"`
	} `json:"floor_number,omitempty"`
	Configuration struct {
		Value string `json:"value,omitempty"`
	} `json:"configuration,omitempty"`
	PossessionStatus struct {
		Value string `json:"value,omitempty"`
	} `json:"possession_status,omitempty"`
	Balconies struct {
		Value string `json:"value,omitempty"`
	} `json:"balconies,omitempty"`
	CoveredParking struct {
		Value string `json:"value,omitempty"`
	} `json:"covered_parking,omitempty"`
	Bedrooms struct {
		Value string `json:"value,omitempty"`
	} `json:"bedrooms,omitempty"`
	PropertyType struct {
		Value string `json:"value,omitempty"`
	} `json:"property_type,omitempty"`
	AgeOfProperty struct {
		Value string `json:"value,omitempty"`
	} `json:"age_of_property,omitempty"`
	FurnishingType struct {
		Value string `json:"value,omitempty"`
	} `json:"furnishing_type,omitempty"`
	Facing struct {
		Value string `json:"value,omitempty"`
	} `json:"facing,omitempty"`
	ReraNumber struct {
		Value string `json:"value,omitempty"`
	} `json:"rera_number,omitempty"`
	Bathrooms struct {
		Value string `json:"value,omitempty"`
	} `json:"bathrooms,omitempty"`
}

type PropertyMetaInfo struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Canonical   string `json:"canonical,omitempty"`
}
