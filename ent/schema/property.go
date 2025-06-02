package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Property struct {
	ent.Schema
}

func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name"),
		field.Text("description"),
		field.JSON("basic_info", PropertyBasicInfo{}),
		field.JSON("area_details", PropertyAreaDetails{}),
		field.JSON("location_details", PropertyLocationDetails{}),
		field.JSON("pricing_info", PropertyPricingInfo{}),
		field.JSON("features", PropertyFeatures{}),
		field.JSON("status_info", PropertyStatusInfo{}),
		field.JSON("property_images", PropertyImages{}),
		field.JSON("property_details", PropertyDetails{}),
		field.JSON("property_amenities", PropertyAmenities{}),
		field.JSON("property_video_presentation", PropertyVideoPresentation{}),
		field.JSON("property_know_about", PropertyKnowAbout{}),
		field.JSON("property_specifications", PropertySpecifications{}),
		field.JSON("pricing_details", PropertyPricingDetails{}).Optional(),
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

// property details - specific to individual property
type PropertyDetails struct {
	UnitNumber struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"unit_number"`
	FloorPlan struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"floor_plan"`
	Direction struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"direction"`
	VentilationType struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"ventilation_type"`
	PowerBackup struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"power_backup"`
	WaterSupply struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"water_supply"`
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

// property video presentation
type PropertyVideoPresentation struct {
	Title          string `json:"title"`
	VirtualTourUrl string `json:"virtual_tour_url"`
	VideoUrl       string `json:"video_url"`
	ThreeDViewUrl  string `json:"3d_view_url"`
}

// property know about - detailed information
type PropertyKnowAbout struct {
	HtmlText          string `json:"html_text"`
	FloorPlanPdf      string `json:"floor_plan_pdf"`
	SpecificationsPdf string `json:"specifications_pdf"`
	LegalDocuments    string `json:"legal_documents"`
}

// property specifications
type PropertySpecifications struct {
	Flooring struct {
		LivingRoom string `json:"living_room"`
		Bedrooms   string `json:"bedrooms"`
		Kitchen    string `json:"kitchen"`
		Bathrooms  string `json:"bathrooms"`
	} `json:"flooring"`
	Doors struct {
		Main     string `json:"main"`
		Internal string `json:"internal"`
	} `json:"doors"`
	Windows struct {
		Type     string `json:"type"`
		Material string `json:"material"`
	} `json:"windows"`
	Kitchen struct {
		Platform string `json:"platform"`
		Sink     string `json:"sink"`
		Chimney  string `json:"chimney"`
	} `json:"kitchen"`
	Bathroom struct {
		Fixtures string `json:"fixtures"`
		Fittings string `json:"fittings"`
	} `json:"bathroom"`
	Electrical struct {
		Wiring   string `json:"wiring"`
		Points   string `json:"points"`
		Switches string `json:"switches"`
	} `json:"electrical"`
}

// property pricing details
type PropertyPricingDetails struct {
	BasePrice                string `json:"base_price"`
	ParkingCharges           string `json:"parking_charges"`
	ClubMembership           string `json:"club_membership"`
	PreferentialFloorCharges string `json:"preferential_floor_charges"`
	MaintenanceDeposit       string `json:"maintenance_deposit"`
	UtilityDeposit           string `json:"utility_deposit"`
	TotalPrice               string `json:"total_price"`
	PaymentSchedule          []struct {
		Milestone  string `json:"milestone"`
		Percentage string `json:"percentage"`
		Amount     string `json:"amount"`
	} `json:"payment_schedule"`
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
	Price              string `json:"price"` // selling price
	PricePerSqft       string `json:"price_per_sqft"`
	MaintenanceCharges string `json:"maintenance_charges"`
	BookingAmount      string `json:"booking_amount"`
	StampDuty          string `json:"stamp_duty"`
	RegistrationFee    string `json:"registration_fee"`
}

// property features
type PropertyFeatures struct {
	FurnishingStatus string `json:"furnishing_status"` // Furnished, Semi-furnished, Unfurnished
	ParkingSpaces    string `json:"parking_spaces"`
	Balconies        string `json:"balconies"`
	Study            bool   `json:"study"`
	ServantRoom      bool   `json:"servant_room"`
	Store            bool   `json:"store"`
	Terrace          bool   `json:"terrace"`
}

// status information
type PropertyStatusInfo struct {
	AvailabilityStatus string `json:"availability_status"` // Available, Sold, Reserved, Under Construction
	ReraID             string `json:"rera_id"`
	LoanAvailable      bool   `json:"loan_available"`
	ReadyToMove        bool   `json:"ready_to_move"`
	UnderConstruction  bool   `json:"under_construction"`
}
