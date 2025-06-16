package migration

import (
	"database/sql"
	"time"
)

// PropertyLegacyData represents the property information from legacy database

// PropertyAmenitiesLegacyData represents the property amenities mapping from legacy database

// ProjectAmenitiesLegacyData represents the project amenities mapping from legacy database
type ProjectAmenitiesLegacyData struct {
	ProjectID   uint64 `json:"project_id"`
	AmenitiesID uint64 `json:"amenities_id"`
}

// PaymentPlanData represents the payment plan data from legacy database
type PaymentPlanData struct {
	ID               int64  `json:"id"`
	PaymentPlanName  string `json:"payment_plan_name"`
	PaymentPlanValue string `json:"payment_plan_value"`
	ProjectID        int64  `json:"project_id"`
}

// FloorPlanData represents the floor plan data from legacy database
type FloorPlanData struct {
	ID              int64   `json:"id"`
	CreatedDate     int64   `json:"created_date"`
	ImgURL          string  `json:"img_url"`
	IsSoldOut       bool    `json:"is_sold_out"`
	Price           float64 `json:"price"`
	Size            int64   `json:"size"`
	Title           string  `json:"title"`
	UpdatedDate     int64   `json:"updated_date"`
	ConfigurationID int64   `json:"configuration_id"`
	ProjectID       int64   `json:"project_id"`
	UserID          int64   `json:"user_id"`
}

// ProjectImageLegacyData represents the project images from legacy database
type ProjectImageLegacyData struct {
	ProjectID    int64  `json:"project_id"`
	ImageAltName string `json:"image_alt_name"`
	ImageURL     string `json:"image_url"`
}

// ProjectConfigTypeData represents the project configuration type data from legacy database
type ProjectConfigTypeData struct {
	ID                    int64  `json:"id"`
	ConfigurationTypeName string `json:"configuration_type_name"`
	PropertyType          string `json:"property_type"`
	CreatedDate           int64  `json:"created_date"`
	UpdatedDate           int64  `json:"updated_date"`
}

// BHKOptionPrice represents a BHK option with its price
type BHKOptionPrice struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

// PriceList represents the price list for a project
type PriceList struct {
	BHKOptionPrices []BHKOptionPrice `json:"bhk_option_prices"`
}

// ProjectVideo represents a video in the project
type ProjectVideo struct {
	VideoURL string `json:"video_url"`
}

// ProjectLegacyData represents the complete project data with its relationships
type ProjectLegacyData struct {
	// Project main fields
	ID                    int64          `json:"id"`
	ProjectName           string         `json:"project_name"`
	ProjectURL            string         `json:"project_url"`
	ProjectDescription    string         `json:"project_description"`
	ProjectAbout          string         `json:"project_about"`
	ProjectAddress        string         `json:"project_address"`
	ProjectArea           string         `json:"project_area"`
	ProjectBrochure       string         `json:"project_brochure"`
	ProjectConfigurations string         `json:"project_configurations"` //
	ProjectLaunchDate     string         `json:"project_launch_date"`
	ProjectLocationURL    string         `json:"project_location_url"`
	ProjectLogo           string         `json:"project_logo"`
	ProjectPossessionDate string         `json:"project_possession_date"`
	ProjectRERA           string         `json:"project_rera"`
	ProjectSchema         string         `json:"project_schema"`
	ProjectUnits          string         `json:"project_units"`
	ProjectVideoCount     int64          `json:"project_video_count"`
	ProjectVideos         []ProjectVideo `json:"project_videos"`
	RERALink              string         `json:"rera_link"`
	ShortAddress          string         `json:"short_address"`
	Status                string         `json:"status"`
	TotalFloor            string         `json:"total_floor"`
	TotalTowers           string         `json:"total_towers"`
	USP                   string         `json:"usp"`
	UserID                int64          `json:"user_id"`
	AvailableUnit         string         `json:"available_unit"`

	// Media and SEO fields
	AltProjectLogo  string   `json:"alt_project_logo"`
	AltSitePlanImg  string   `json:"alt_site_plan_img"`
	CoverPhoto      string   `json:"cover_photo"`
	LocationMap     string   `json:"location_map"`
	SitePlanImg     string   `json:"siteplan_img"`
	MetaDescription string   `json:"meta_description"`
	MetaTitle       string   `json:"meta_title"`
	MetaKeywords    []string `json:"meta_keywords"`

	// Content paragraphs
	AmenitiesPara string `json:"amenities_para"`
	FloorPara     string `json:"floor_para"`
	LocationPara  string `json:"location_para"`
	OverviewPara  string `json:"overview_para"`
	PaymentPara   string `json:"payment_para"`
	PriceListPara string `json:"price_list_para"`
	SitePlanPara  string `json:"siteplan_para"`
	VideoPara     string `json:"video_para"`
	WhyPara       string `json:"why_para"`

	// Boolean flags
	IsDeleted  bool `json:"is_deleted"`
	IsFeatured bool `json:"is_featured"`
	IsPremium  bool `json:"is_premium"`
	IsPriority bool `json:"is_priority"`

	// Timestamps
	CreatedDate sql.NullInt64 `json:"created_date"`
	UpdatedDate sql.NullInt64 `json:"updated_date"`

	// Related data
	Developer     *DeveloperData               `json:"developer"`
	Locality      *LocalityData                `json:"locality"`
	Configuration *ConfigurationData           `json:"configuration"`
	ConfigType    *ProjectConfigTypeData       `json:"config_type"` // Add ConfigType field
	Properties    []PropertyLegacyData         `json:"properties"`  // List of properties
	Amenities     []ProjectAmenitiesLegacyData `json:"amenities"`   // List of project amenities
	FAQs          []FAQLegacyData              `json:"faqs"`        // List of project FAQs
	AmenitiesData []AmenitiesData              `json:"amenities_data"`
	PaymentPlans  []PaymentPlanData            `json:"payment_plans"` // Add PaymentPlans field
	FloorPlans    []FloorPlanData              `json:"floor_plans"`   // Add FloorPlans field

	// Add ReraInfo field
	ReraInfo []ReraInfoData `json:"rera_info"`

	// Add ProjectImages field
	ProjectImages []ProjectImageLegacyData `json:"project_images"`

	// Add PriceList field
	PriceList PriceList `json:"price_list"`
}

// FAQLegacyData represents the FAQ information from legacy database
type FAQLegacyData struct {
	ID        int64  `json:"id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	ProjectID int64  `json:"project_id"`
}

type PropertyLegacyData struct {
	ID                uint64                        `json:"id"`
	About             *string                       `json:"about"`
	AgeOfProperty     *string                       `json:"age_of_property"`
	AmenitiesPara     *string                       `json:"amenities_para"`
	Balcony           *string                       `json:"balcony"`
	Bathrooms         *string                       `json:"bathrooms"`
	Bedrooms          *string                       `json:"bedrooms"`
	BuiltupArea       *string                       `json:"builtup_area"`
	CoverPhoto        *string                       `json:"cover_photo"`
	CoveredParking    *string                       `json:"covered_parking"`
	CreatedDate       *int64                        `json:"created_date"`
	Facing            *string                       `json:"facing"`
	FloorImage2D      *string                       `json:"floor_image2d"`
	FloorImage3D      *string                       `json:"floor_image3d"`
	FloorPara         *string                       `json:"floor_para"`
	Floors            *string                       `json:"floors"`
	FurnishingType    *string                       `json:"furnishing_type"`
	Images            *string                       `json:"images"`
	IsDeleted         bool                          `json:"is_deleted"`
	IsFeatured        bool                          `json:"is_featured"`
	Latlong           *string                       `json:"latlong"`
	ListingType       *string                       `json:"listing_type"`
	LocationMap       *string                       `json:"location_map"`
	LocationAdvantage *string                       `json:"location_advantage"`
	LocationPara      *string                       `json:"location_para"`
	LogoImage         *string                       `json:"logo_image"`
	MetaDescription   *string                       `json:"meta_description"`
	MetaKeywords      *string                       `json:"meta_keywords"`
	MetaTitle         *string                       `json:"meta_title"`
	OpenParking       *string                       `json:"open_parking"`
	OverviewPara      *string                       `json:"overview_para"`
	PossessionDate    *string                       `json:"possession_date"`
	PossessionStatus  *string                       `json:"possession_status"`
	Price             float64                       `json:"price"`
	ProductSchema     *string                       `json:"product_schema"`
	PropertyAddress   *string                       `json:"property_address"`
	PropertyName      *string                       `json:"property_name"`
	PropertyURL       *string                       `json:"property_url"`
	PropertyVideo     *string                       `json:"property_video"`
	Rera              *string                       `json:"rera"`
	Size              *string                       `json:"size"`
	UpdatedByID       *uint64                       `json:"updated_by_id"`
	UpdatedDate       *int64                        `json:"updated_date"`
	USP               *string                       `json:"usp"`
	VideoPara         *string                       `json:"video_para"`
	ConfiurationID    *uint64                       `json:"confiuration_id"`
	DeveloperID       *uint64                       `json:"developer_id"`
	LocalityID        *uint64                       `json:"locality_id"`
	ProjectID         *uint64                       `json:"project_id"`
	UserID            *uint64                       `json:"user_id"`
	Highlights        *string                       `json:"highlights"`
	Amenities         []PropertyAmenitiesLegacyData `json:"amenities"`
}

type PropertyAmenitiesLegacyData struct {
	PropertyID  uint64 `json:"property_id"`
	AmenitiesID uint64 `json:"amenities_id"`
}

// DeveloperData represents the developer information
type DeveloperData struct {
	ID                 int64         `json:"id"`
	DeveloperName      string        `json:"developer_name"`
	DeveloperLegalName string        `json:"developer_legal_name"`
	DeveloperURL       string        `json:"developer_url"`
	DeveloperAddress   string        `json:"developer_address"`
	DeveloperLogo      string        `json:"developer_logo"`
	AltDeveloperLogo   string        `json:"alt_developer_logo"`
	About              string        `json:"about"`
	Overview           string        `json:"overview"`
	Disclaimer         string        `json:"disclaimer"`
	EstablishedYear    int64         `json:"established_year"`
	ProjectDoneNo      int64         `json:"project_done_no"`
	Phone              string        `json:"phone"`
	IsActive           bool          `json:"is_active"`
	IsVerified         bool          `json:"is_verified"`
	CityName           int64         `json:"city_name"`
	CreatedDate        sql.NullInt64 `json:"created_date"`
	UpdatedDate        sql.NullInt64 `json:"updated_date"`
}

// LocalityData represents the locality information
type LocalityData struct {
	ID           int64         `json:"id"`
	LocalityName string        `json:"locality_name"`
	LocalityURL  string        `json:"locality_url"`
	CityID       int64         `json:"city_id"`
	CreatedDate  sql.NullInt64 `json:"created_date"`
	UpdatedDate  sql.NullInt64 `json:"updated_date"`
}

// ConfigurationData represents the project configuration information
type ConfigurationData struct {
	ID                       int64         `json:"id"`
	ProjectConfigurationName string        `json:"project_configuration_name"`
	ConfigurationTypeID      int64         `json:"configuration_type_id"`
	CreatedDate              sql.NullInt64 `json:"created_date"`
	UpdatedDate              sql.NullInt64 `json:"updated_date"`
}

// AmenitiesData represents the amenities data from legacy database
type AmenitiesData struct {
	ID                int64  `json:"id"`
	AmenitiesCategory string `json:"amenities_category"`
	AmenitiesName     string `json:"amenities_name"`
	AmenitiesURL      string `json:"amenities_url"`
	CreatedDate       int64  `json:"created_date"`
	UpdatedDate       int64  `json:"updated_date"`
}

// ReraInfoData represents the RERA info data from legacy database
type ReraInfoData struct {
	ID              int64    `json:"id"`
	CreatedOn       int64    `json:"created_on"`
	Phase           []string `json:"phase"`
	ProjectReraName string   `json:"project_rera_name"`
	QRImages        []string `json:"qr_images"`
	ReraNumber      []string `json:"rera_number"`
	Status          []string `json:"status"`
	UpdatedOn       int64    `json:"updated_on"`
	ProjectID       int64    `json:"project_id"`
	UserID          int64    `json:"user_id"`
}

// Helper function to convert Unix timestamp to time.Time
func unixToTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Time{}
	}
	return time.Unix(timestamp, 0)
}
