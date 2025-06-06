package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Project struct {
	Base
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.JSON("basic_info", BasicInfo{}),
		field.JSON("timeline_info", TimelineInfo{}),
		field.JSON("meta_info", SEOMeta{}),
		field.JSON("web_cards", ProjectWebCards{}),
		field.JSON("location_info", LocationInfo{}),
		field.Bool("is_featured").Default(false),
		field.Bool("is_premium").Default(false),
		field.Bool("is_priority").Default(false),
		field.Bool("is_deleted").Default(false),
		field.JSON("search_context", []string{}),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("properties", Property.Type),
		edge.From("location", Location.Type).Ref("projects").Unique(),
		edge.From("developer", Developer.Type).Ref("projects").Unique(),
	}
}

// web cards
type ProjectWebCards struct {
	Images            []string          `json:"images"`
	ReraInfo          ReraInfo          `json:"rera_info"`
	Details           ProjectDetails    `json:"project_details"`
	WhyToChoose       WhyToChoose       `json:"why_to_choose"`
	KnowAbout         KnowAbout         `json:"know_about"`
	FloorPlan         FloorPlan         `json:"floor_plan"`
	PriceList         PriceList         `json:"price_list"`
	Amenities         Amenities         `json:"amenities"`
	VideoPresentation VideoPresentation `json:"video_presentation"`
	PaymentPlans      PaymentPlans      `json:"payment_plans"`
	SitePlan          struct {
		HTMLContent string `json:"html_content"`
		Image       string `json:"image"`
	} `json:"site_plan"`
	About About `json:"about"`
	Faqs  []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"faqs"`
}

// project details
type ProjectDetails struct {
	Area struct {
		Value string `json:"value"`
	} `json:"area"`
	Sizes struct {
		Value string `json:"value"`
	} `json:"sizes"`
	Units struct {
		Value string `json:"value"`
	} `json:"units"`
	LaunchDate struct {
		Value string `json:"value"`
	} `json:"launch_date"`
	PossessionDate struct {
		Value string `json:"value"`
	} `json:"possession_date"`
	TotalTowers struct {
		Value string `json:"value"`
	} `json:"total_towers"`
	TotalFloors struct {
		Value string `json:"value"`
	} `json:"total_floors"`
	ProjectStatus struct {
		Value string `json:"value"`
	} `json:"project_status"`
	Type struct {
		Value string `json:"value"`
	} `json:"type"`
}

// Rera info
type ReraInfo struct {
	WebsiteLink string `json:"website_link"`
	ReraList    []struct {
		Phase      string `json:"phase"`
		ReraQR     string `json:"rera_qr"`
		ReraNumber string `json:"rera_number"`
		Status     string `json:"status"`
	} `json:"rera_list"`
}

// why to choose
type WhyToChoose struct {
	ImageUrls []string `json:"image_urls"`
	USP_List  []struct {
		Icon        string `json:"icon"`
		Description string `json:"description"`
	} `json:"usp_list"`
}

// know about
type KnowAbout struct {
	Description  string `json:"description"`
	DownloadLink string `json:"download_link"`
}

// floor plan
type FloorPlan struct {
	Title    string `json:"title"`
	Products []struct {
		Title        string `json:"title"`
		FlatType     string `json:"flat_type"`
		Price        string `json:"price"`
		BuildingArea string `json:"building_area"`
		Image        string `json:"image"`
	} `json:"products"`
}

type PriceList struct {
	Description          string `json:"description"`
	BHKOptionsWithPrices []struct {
		BHKOption string `json:"bhk_option"`
		Size      string `json:"size"`
		Price     string `json:"price"`
	} `json:"bhk_options_with_prices"`
}

// amenities
type Amenities struct {
	Description             string `json:"description"`
	CategoriesWithAmenities map[string][]struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"categories_with_amenities"`
}

// video presentation
type VideoPresentation struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// about

type About struct {
	Description       string `json:"description"`
	LogoURL           string `json:"logo_url"`
	EstablishmentYear string `json:"establishment_year"`
	TotalProperties   string `json:"total_properties"`
	ContactDetails    struct {
		Name           string `json:"name"`
		ProjectAddress string `json:"project_address"`
		Phone          string `json:"phone"`
		BookingLink    string `json:"booking_link"`
	} `json:"contact_details"`
}

// payment plans
type PaymentPlans struct {
	Description string `json:"description"`
	Plans       []struct {
		Name    string `json:"name"`
		Details string `json:"details"`
	} `json:"plans"`
}

// basic project information
type BasicInfo struct {
	ProjectName           string `json:"project_name"`
	ProjectDescription    string `json:"project_description"`
	ProjectArea           string `json:"project_area"`
	ProjectUnits          string `json:"project_units"`
	ProjectConfigurations string `json:"project_configurations"`
	AvailableUnit         string `json:"available_unit"`
	TotalFloor            string `json:"total_floor"`
	TotalTowers           string `json:"total_towers"`
	ProjectType           string `json:"project_type"`
	MinPrice              string `json:"min_price"`
	MaxPrice              string `json:"max_price"`
	Status                string `json:"status"`
}

// timeline information
type TimelineInfo struct {
	ProjectLaunchDate     string `json:"project_launch_date"`
	ProjectPossessionDate string `json:"project_possession_date"`
}

// SEO and meta information
type SEOMeta struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Keywords      string `json:"keywords"`
	Canonical     string `json:"canonical"`
	ProjectSchema string `json:"project_schema"` //[ "<script type=\"application/ld+json\">\n{\n  \"@context\": \"https://schema.org/\",\n  \"@type\": \"Product\",\n  \"name\": \"ACE Divino\",\n  \"image\": \"https://image.investmango.com/images/img/ace-divino/ace-divino-greater-noida-west.webp\",\n  \"description\": \"ACE Divino Sector 1, Noida Extension: Explore prices, floor plans, payment options, location, photos, videos, and more. Download the project brochure now!\",\n  \"brand\": {\n    \"@type\": \"Brand\",\n    \"name\": \"Ace Group of India\"\n  },\n  \"offers\": {\n    \"@type\": \"AggregateOffer\",\n    \"url\": \"https://www.investmango.com/ace-divino\",\n    \"priceCurrency\": \"INR\",\n    \"lowPrice\": \"18800000\",\n    \"highPrice\": \"22500000\"\n  }\n}\n</script>" ]
}

type LocationInfo struct {
	ShortAddress  string `json:"short_address"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	GoogleMapLink string `json:"google_map_link"`
	City          string `json:"city"`
}
