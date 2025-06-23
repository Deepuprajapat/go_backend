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
		field.Text("description"),
		field.JSON("property_images", []string{}), // 0 index logo image
		field.JSON("web_cards", WebCards{}),
		field.JSON("pricing_info", PropertyPricingInfo{}),
		field.JSON("property_rera_info", PropertyReraInfo{}),
		field.JSON("search_context", []string{}).Optional(),
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
	ReraNumber string `json:"rera_number"`
}

type WebCards struct {
	PropertyDetails PropertyDetails `json:"property_details"`
	WhyChooseUs     struct {
		ImageUrls []string `json:"image_urls"`
		USP_List  []string `json:"usp_list"`
	} `json:"why_choose_us"`
	PropertyFloorPlan PropertyFloorPlan `json:"property_floor_plan"`
	KnowAbout         struct {
		Description string `json:"description"`
	} `json:"know_about"`
	VideoPresentation struct {
		Title    string `json:"title"`
		VideoUrl string `json:"video_url"`
	} `json:"video_presentation"`
	LocationMap struct {
		Description   string `json:"description"`
		GoogleMapLink string `json:"google_map_link"`
	} `json:"location_map"`
}

type PropertyUSPListItem struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type PropertyFloorPlan struct {
	Title string         `json:"title"`
	Plans []PropertyPlan `json:"plans"`
}

type PropertyPlan struct {
	Title        string `json:"title"`
	FlatType     string `json:"flat_type"`
	Price        string `json:"price"`
	BuildingArea string `json:"building_area"`
	Image        string `json:"image"`
}

// property amenities - specific to this property unit
type PropertyAmenities struct {
	Description             string                       `json:"description"`
	CategoriesWithAmenities map[string][]AmenityCategory `json:"categories_with_amenities"`
}

// area details
type PropertyAreaDetails struct {
	CarpetArea       string `json:"carpet_area"`         // in sq ft
	BuiltUpArea      string `json:"built_up_area"`       // in sq ft
	SuperBuiltUpArea string `json:"super_built_up_area"` // in sq ft
}

// pricing information
type PropertyPricingInfo struct {
	Price string `json:"price"` // selling price
}

// property details
type PropertyDetails struct {
	BuiltUpArea struct {
		Value string `json:"value"`
	} `json:"built_up_area"`
	Sizes struct {
		Value string `json:"value"`
	} `json:"size"`
	FloorNumber struct {
		Value string `json:"value"`
	} `json:"floor_number"`
	Configuration struct {
		Value string `json:"value"`
	} `json:"configuration"`
	PossessionStatus struct {
		Value string `json:"value"`
	} `json:"possession_status"`
	Balconies struct {
		Value string `json:"value"`
	} `json:"balconies"`
	CoveredParking struct {
		Value string `json:"value"`
	} `json:"covered_parking"`
	Bedrooms struct {
		Value string `json:"value"`
	} `json:"bedrooms"`
	PropertyType struct {
		Value string `json:"value"`
	} `json:"property_type"`
	AgeOfProperty struct {
		Value string `json:"value"`
	} `json:"age_of_property"`
	FurnishingType struct {
		Value string `json:"value"`
	} `json:"furnishing_type"`
	Facing struct {
		Value string `json:"value"`
	} `json:"facing"`
	ReraNumber struct {
		Value string `json:"value"`
	} `json:"rera_number"`
	Bathrooms struct {
		Value string `json:"value"`
	} `json:"bathrooms"`
}
