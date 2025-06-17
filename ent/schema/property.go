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
		field.JSON("property_images", PropertyImages{}),
		field.JSON("web_cards", WebCards{}),
		field.JSON("basic_info", PropertyBasicInfo{}),
		field.JSON("location_details", PropertyLocationDetails{}),
		field.JSON("pricing_info", PropertyPricingInfo{}),
		field.JSON("property_rera_info", PropertyReraInfo{}),
		field.JSON("search_context", []string{}),
	}
}

func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).Ref("properties").Unique(),
		edge.To("leads", Leads.Type),
	}
}

// property images
type PropertyImages struct {
	Images []struct {
		Order int    `json:"order"`
		Url   string `json:"url"`
		Type  string `json:"type"` // floor_plan, interior, exterior, amenities
	} `json:"images"`
}

type PropertyReraInfo struct {
	Phase      string `json:"phase"`
	Status     string `json:"status"`
	ReraNumber string `json:"rera_number"`
	ReraQR     string `json:"rera_qr"`
}

type WebCards struct {
	PropertyDetails   PropertyDetails `json:"property_details"`
	PropertyFloorPlan []struct {
		Title string `json:"title"`
		Plans []struct {
			Title        string `json:"title"`
			FlatType     string `json:"flat_type"`
			Price        string `json:"price"`
			BuildingArea string `json:"building_area"`
			Image        string `json:"image"`
			ExpertLink   string `json:"expert_link"`
			BrochureLink string `json:"brochure_link"`
		} `json:"plans"`
	} `json:"property_floor_plan"`
	KnowAbout struct {
		HtmlText string `json:"html_text"`
	} `json:"know_about"`
	VideoPresentation struct {
		Title    string `json:"title"`
		VideoUrl string `json:"video_url"`
	} `json:"video_presentation"`
	GoogleMapLink string `json:"google_map_link"`
}

// property amenities - specific to this property unit
type PropertyAmenities struct {
	UnitAmenities []struct {
		Icon string `json:"icon"`
		Name string `json:"name"`
	} `json:"unit_amenities"`
	FloorAmenities []struct {
		Icon string `json:"icon"`
		Name string `json:"name"`
	} `json:"floor_amenities"`
}

// basic property information
type PropertyBasicInfo struct {
	PropertyType string `json:"property_type"` // apartment, villa, penthouse, studio
	BHKType      string `json:"bhk_type"`      // 1BHK, 2BHK, 3BHK, etc.
	Bedrooms     int    `json:"bedrooms"`
	Bathrooms    int    `json:"bathrooms"`
}

// area details
type PropertyAreaDetails struct {
	CarpetArea       string `json:"carpet_area"`         // in sq ft
	BuiltUpArea      string `json:"built_up_area"`       // in sq ft
	SuperBuiltUpArea string `json:"super_built_up_area"` // in sq ft
}

// location details
type PropertyLocationDetails struct {
	FloorNumber int    `json:"floor_number"`
	Facing      string `json:"facing"` // North, South, East, West, etc.
	Tower       string `json:"tower"`  // Tower name/number if applicable
	Wing        string `json:"wing"`   // Wing name if applicable
}

// pricing information
type PropertyPricingInfo struct {
	StartingPrice      string `json:"starting_price"`
	Price              string `json:"price"` // selling price
	PricePerSqft       string `json:"price_per_sqft"`
	MaintenanceCharges string `json:"maintenance_charges"`
	BookingAmount      string `json:"booking_amount"`
	StampDuty          string `json:"stamp_duty"`
	RegistrationFee    string `json:"registration_fee"`
}

// property details
type PropertyDetails struct {
	PropertyType      string `json:"property_type"`
	FurnishingType    string `json:"furnishing_type"`
	ListingType       string `json:"listing_type"`
	PossessionStatus  string `json:"possession_status"`
	AgeOfProperty     string `json:"age_of_property"`
	FloorPara         string `json:"floor_para"`
	LocationPara      string `json:"location_para"`
	LocationAdvantage string `json:"location_advantage"`
	OverviewPara      string `json:"overview_para"`
	Floors            string `json:"floors"`
	Images            string `json:"images"`
	Latlong           string `json:"latlong"`
}
