package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Project struct {
	ent.Schema
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.JSON("basic_info", BasicInfo{}),
		field.JSON("timeline_info", TimelineInfo{}),
		field.JSON("seo_meta", SEOMeta{}),
		field.JSON("website_cards", WebsiteCards{}),
		field.Bool("is_featured").Default(false),
		field.Bool("is_premium").Default(false),
		field.Bool("is_priority").Default(false),
		field.Bool("is_deleted").Default(false),
		field.JSON("search_context", []string{}),
		field.Time("updated_at"),
		field.Time("created_at").Default(time.Now()),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("properties", Property.Type),
		edge.From("location", Location.Type).Ref("projects").Unique(),
		edge.From("developer", Developer.Type).Ref("projects").Unique(),
	}
}

// website cards
type WebsiteCards struct {
	Images            []string          `json:"images"`
	ReraInfo          ReraInfo          `json:"rera_info"`
	Details           Details           `json:"details"`
	WhyToChoose       WhyToChoose       `json:"why_to_choose"`
	KnowAbout         KnowAbout         `json:"know_about"`
	FloorPlan         FloorPlan         `json:"floor_plan"`
	PriceList         PriceList         `json:"price_list"`
	Amenities         Amenities         `json:"amenities"`
	VideoPresentation VideoPresentation `json:"video_presentation"`
	PaymentPlans      PaymentPlans      `json:"payment_plans"`
	SitePlan          SitePlan          `json:"site_plan"`
	About             About             `json:"about"`
	Faqs              Faqs              `json:"faqs"`
}

// project details
type Details struct {
	ProjectArea struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"project_area"`
	Sizes struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"sizes"`
	ProjectUnits struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"project_units"`
	LaunchDate struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"launch_date"`
	PossessionDate struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"possession_date"`
	TotalTowers struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"total_towers"`
	TotalFloors struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"total_floors"`
	ProjectStatus struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"project_status"`
	PropertyType struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	} `json:"property_type"`
}

// Rera info
type ReraInfo struct {
	CreatedOn       time.Time `json:"created_on"`
	Phase           string    `json:"phase"`
	ProjectReraName string    `json:"project_rera_name"`
	QRImages        struct {
		Url string `json:"url"`
	} `json:"qr_images"`
	ReraNumber string    `json:"rera_number"`
	Status     string    `json:"status"`
	UpdatedOn  time.Time `json:"updated_on"`
}

// why to choose
type WhyToChoose struct {
	Images []struct {
		Order int    `json:"order"`
		Url   string `json:"url"`
	} `json:"images"`
	UspList []struct {
		Icon        string `json:"icon"`
		HtmlContent string `json:"html_content"`
	} `json:"usp_list"`
}

// know about
type KnowAbout struct {
	HtmlText     string `json:"html_text"`
	DownloadLink string `json:"download_link"`
}

// floor plan
type FloorPlan struct {
	Discription string `json:"discription"`
	Products    []struct {
		Title          string `json:"title"`
		Congfiguration string `json:"flat_type"`
		Price          string `json:"price"`
		BuildingArea   string `json:"building_area"`
		Image          string `json:"image"`
		ExpertLink     string `json:"expert_link"`
		BrochureLink   string `json:"brochure_link"`
	} `json:"plans"`
}

type PriceList struct {
	Discription          string `json:"title"`
	BHKOptionsWithPrices []struct {
		BHKOption string `json:"bhk_option"`
		Size      string `json:"size"`
		Price     string `json:"price"`
	} `json:"bhk_options_with_prices"`
}

// amenities
type Amenities struct {
	Discription   string `json:"title"`
	AmenitiesList []struct {
		AmenityType []struct {
			Icon string `json:"icon"`
			Text string `json:"value"`
		} `json:"amenity_type"`
	} `json:"amenities_list"`
}

// video presentation
type VideoPresentation struct {
	Discription string `json:"title"`
	URL         string `json:"video_url"`
}

// about

type About struct {
	// Title             string `json:"title"`
	LogoURL           string `json:"logo_url"`
	EstablishmentYear string `json:"establishment_year"`
	TotalProperties   string `json:"total_properties"`
	HTMLContent       string `json:"html_content"`
	ContactDetails    struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		BookingLink string `json:"booking_link"`
	} `json:"contact_details"`
}

// site plan
type SitePlan struct {
	Discription string `json:"discription"`
	Image       string `json:"image"`
}

// payment plans
type PaymentPlans struct {
	Discription string `json:"discription"`
	Plans       []struct {
		Name    string `json:"name"`
		Details string `json:"details"`
	} `json:"plans"`
}

// faqs
type Faqs struct {
	Faqs []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"faqs"`
}

// basic project information
type BasicInfo struct {
	ProjectDescription    string `json:"project_description"`
	ProjectArea           string `json:"project_area"`
	ProjectLogoURL        string `json:"project_logo_url"`
	ProjectUnits          string `json:"project_units"`
	ProjectConfigurations string `json:"project_configurations"`
	AvailableUnit         string `json:"available_unit"`
	TotalFloor            string `json:"total_floor"`
	TotalTowers           string `json:"total_towers"`
	Status                string `json:"status"`
}

// timeline information
type TimelineInfo struct {
	ProjectLaunchDate     string `json:"project_launch_date"`
	ProjectPossessionDate string `json:"project_possession_date"`
}

// SEO and meta information
type SEOMeta struct {
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	ProjectUrl      string `json:"project_url"`
}
