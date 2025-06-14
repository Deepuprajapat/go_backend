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
		field.String("id").Unique(),
		field.String("name"),
		field.Text("description"),
		field.Enum("status").
			Values("UNDER_CONSTRUCTION", "READY_TO_MOVE", "NEW_LAUNCH", "PRE_LAUNCH").
			Default("UNDER_CONSTRUCTION").Optional(),
		field.Enum("project_configurations").
			Values("1BHK", "2BHK", "3BHK", "4BHK", "5BHK", "6BHK", "7BHK", "8BHK").Optional(),
		field.Int("total_floor").Optional(),
		field.Int("total_towers").Optional(),
		field.Int("min_price").Default(0), // update on every add property
		field.Int("max_price").Default(0), // update on every add property
		field.String("price_unit").Default("cr"),
		field.JSON("timeline_info", TimelineInfo{}).Optional(),
		field.JSON("meta_info", SEOMeta{}).Optional(),
		field.JSON("web_cards", ProjectWebCards{}).Optional(),
		field.JSON("location_info", LocationInfo{}).Optional(),
		field.Bool("is_featured").Default(false).Optional(),
		field.Bool("is_premium").Default(false).Optional(),
		field.Bool("is_priority").Default(false).Optional(),
		field.Bool("is_deleted").Default(false).Optional(),
		field.JSON("search_context", []string{}).Optional(),
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
		Description string `json:"description"`
		Image       string `json:"image"`
	} `json:"site_plan"`
	About About `json:"about"`
	Faqs  []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"faqs"`
}

// project info
type ProjectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Area        string `json:"area"`
	LogoURL     string `json:"logo_url"`
	MinPrice    string `json:"min_price"`
	MaxPrice    string `json:"max_price"`
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
	Configuration struct {
		Value string `json:"value"`
	} `json:"configuration"`

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
	TotalProjects     string `json:"total_projects"`
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
}
