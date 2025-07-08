package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type Project struct {
	Base
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("name"),
		field.Text("description"),
		field.String("status").GoType(enums.ProjectStatus("")),
		field.String("min_price").Default("0").Optional(),
		field.String("max_price").Default("0").Optional(),
		field.JSON("timeline_info", TimelineInfo{}).Optional(),
		field.Enum("project_type").Values("RESIDENTIAL", "COMMERCIAL"),
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
	Images            []string          `json:"images,omitempty"`
	ReraInfo          ReraInfo          `json:"rera_info,omitempty"`
	Details           ProjectDetails    `json:"project_details,omitempty"`
	WhyToChoose       WhyToChoose       `json:"why_to_choose,omitempty"`
	KnowAbout         KnowAbout         `json:"know_about,omitempty"`
	FloorPlan         FloorPlan         `json:"floor_plan,omitempty"`
	PriceList         PriceList         `json:"price_list,omitempty"`
	Amenities         Amenities         `json:"amenities,omitempty"`
	VideoPresentation VideoPresentation `json:"video_presentation,omitempty"`
	PaymentPlans      PaymentPlans      `json:"payment_plans,omitempty"`
	SitePlan          struct {
		Description string `json:"description,omitempty"`
		Image       string `json:"image,omitempty"`
	} `json:"site_plan,omitempty"`
	About struct {
		Description       string `json:"description,omitempty"`
		LogoURL           string `json:"logo_url,omitempty"`
		EstablishmentYear string `json:"establishment_year,omitempty"`
		TotalProjects     string `json:"total_projects,omitempty"`
		ContactDetails    struct {
			Name           string `json:"name,omitempty"`
			ProjectAddress string `json:"project_address,omitempty"`
			Phone          string `json:"phone,omitempty"`
			BookingLink    string `json:"booking_link,omitempty"`
		} `json:"contact_details,omitempty"`
	} `json:"about,omitempty"`
	Faqs []FAQ `json:"faqs,omitempty"`
}

type FAQ struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

// project details
type ProjectDetails struct {
	Area struct {
		Value string `json:"value,omitempty"`
	} `json:"area,omitempty"`
	Sizes struct {
		Value string `json:"value,omitempty"`
	} `json:"sizes,omitempty"`
	Units struct {
		Value string `json:"value,omitempty"`
	} `json:"units,omitempty"`
	Configuration struct {
		Value string `json:"value,omitempty"`
	} `json:"configuration,omitempty"`
	TotalFloor struct {
		Value string `json:"value,omitempty"`
	} `json:"total_floor,omitempty"`
	TotalTowers struct {
		Value string `json:"value,omitempty"`
	} `json:"total_towers,omitempty"`
	LaunchDate struct {
		Value string `json:"value,omitempty"`
	} `json:"launch_date,omitempty"`
	PossessionDate struct {
		Value string `json:"value,omitempty"`
	} `json:"possession_date,omitempty"`
	Type struct {
		Value string `json:"value,omitempty"`
	} `json:"type,omitempty"`
}

// Rera info
type ReraInfo struct {
	WebsiteLink string         `json:"website_link,omitempty"`
	ReraList    []ReraListItem `json:"rera_list,omitempty"`
}

type ReraListItem struct {
	Phase      string `json:"phase,omitempty"`
	ReraQR     string `json:"rera_qr,omitempty"`
	ReraNumber string `json:"rera_number,omitempty"`
	Status     string `json:"status,omitempty"`
}

// why to choose
type WhyToChoose struct {
	ImageUrls []string `json:"image_urls,omitempty"`
	USP_List  []string `json:"usp_list,omitempty"`
}

// know about
type KnowAbout struct {
	Description  string `json:"description,omitempty"`
	DownloadLink string `json:"download_link,omitempty"`
}

// floor plan
type FloorPlan struct {
	Description string          `json:"description,omitempty"`
	Products    []FloorPlanItem `json:"products,omitempty"`
}

type FloorPlanItem struct {
	Title        string `json:"title,omitempty"`
	FlatType     string `json:"flat_type,omitempty"`
	Price        string `json:"price,omitempty"`
	IsSoldOut    bool   `json:"is_sold_out,omitempty"`
	BuildingArea string `json:"building_area,omitempty"`
	Image        string `json:"image,omitempty"`
}

type PriceList struct {
	Description          string                 `json:"description,omitempty"`
	BHKOptionsWithPrices []ProductConfiguration `json:"product_configurations,omitempty"`
}

type ProductConfiguration struct {
	ConfigurationName string `json:"configuration_name,omitempty"`
	Size              string `json:"size,omitempty"`
	Price             string `json:"price,omitempty"`
}

// amenities
type Amenities struct {
	Description             string                       `json:"description,omitempty"`
	CategoriesWithAmenities map[string][]AmenityCategory `json:"categories_with_amenities,omitempty"`
}

type AmenityCategory struct {
	Icon  string `json:"icon,omitempty"`
	Value string `json:"value,omitempty"`
}

// video presentation
type VideoPresentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// payment plans
type PaymentPlans struct {
	Description string `json:"description,omitempty"`
	Plans       []Plan `json:"plans,omitempty"`
}

type Plan struct {
	Name    string `json:"name,omitempty"`
	Details string `json:"details,omitempty"`
}

// timeline information
type TimelineInfo struct {
	ProjectLaunchDate     string `json:"project_launch_date,omitempty"`
	ProjectPossessionDate string `json:"project_possession_date,omitempty"`
}

// SEO and meta information
type SEOMeta struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Keywords      string `json:"keywords,omitempty"`
	Canonical     string `json:"canonical,omitempty"`
	ProjectSchema string `json:"project_schema,omitempty"` //[ "<script type=\"application/ld+json\">\n{\n  \"@context\": \"https://schema.org/\",\n  \"@type\": \"Product\",\n  \"name\": \"ACE Divino\",\n  \"image\": \"https://image.investmango.com/images/img/ace-divino/ace-divino-greater-noida-west.webp\",\n  \"description\": \"ACE Divino Sector 1, Noida Extension: Explore prices, floor plans, payment options, location, photos, videos, and more. Download the project brochure now!\",\n  \"brand\": {\n    \"@type\": \"Brand\",\n    \"name\": \"Ace Group of India\"\n  },\n  \"offers\": {\n    \"@type\": \"AggregateOffer\",\n    \"url\": \"https://www.investmango.com/ace-divino\",\n    \"priceCurrency\": \"INR\",\n    \"lowPrice\": \"18800000\",\n    \"highPrice\": \"22500000\"\n  }\n}\n</script>" ]
}

type LocationInfo struct {
	ShortAddress  string `json:"short_address,omitempty"`
	Longitude     string `json:"longitude,omitempty"`
	Latitude      string `json:"latitude,omitempty"`
	GoogleMapLink string `json:"google_map_link,omitempty"`
}

type Configurations struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"` // apartment, villa, penthouse, studio
}
