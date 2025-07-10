package migration_jobs

type LProject struct {
	ID                    int64   `json:"id"`
	AltProjectLogo        *string `json:"alt_project_logo"`
	AltSitePlanImg        *string `json:"alt_site_plan_img"`
	AmenitiesPara         *string `json:"amenities_para"`
	AvailableUnit         *string `json:"available_unit"`
	CoverPhoto            *string `json:"cover_photo"`
	CreatedDate           *int64  `json:"created_date"`
	FloorPara             *string `json:"floor_para"`
	IsDeleted             bool    `json:"is_deleted"`
	IsFeatured            bool    `json:"is_featured"`
	IsPremium             bool    `json:"is_premium"`
	IsPriority            bool    `json:"is_priority"`
	LocationMap           *string `json:"location_map"`
	LocationPara          *string `json:"location_para"`
	MetaDescription       *string `json:"meta_description"`
	MetaTitle             *string `json:"meta_title"`
	OverviewPara          *string `json:"overview_para"`
	PaymentPara           *string `json:"payment_para"`
	PriceListPara         *string `json:"price_list_para"`
	ProjectAbout          *string `json:"project_about"`
	ProjectAddress        *string `json:"project_address"`
	ProjectArea           *string `json:"project_area"`
	ProjectBrochure       *string `json:"project_brochure"`
	ProjectConfigurations *string `json:"project_configurations"`
	ProjectDescription    *string `json:"project_description"`
	ProjectLaunchDate     *string `json:"project_launch_date"`
	ProjectLocationURL    *string `json:"project_location_url"`
	ProjectLogo           *string `json:"project_logo"`
	ProjectName           *string `json:"project_name"`
	ProjectPossessionDate *string `json:"project_possession_date"`
	ProjectRERA           *string `json:"project_rera"`
	ProjectSchema         *string `json:"project_schema"`
	ProjectUnits          *string `json:"project_units"`
	ProjectURL            *string `json:"project_url"`
	ProjectVideoCount     *int64  `json:"project_video_count"`
	ProjectVideos         []byte  `json:"project_videos"`
	ReraLink              *string `json:"rera_link"`
	ShortAddress          *string `json:"short_address"`
	SitePlanImg           *string `json:"siteplan_img"`
	SitePlanPara          *string `json:"siteplan_para"`
	Status                *string `json:"status"`
	TotalFloor            *string `json:"total_floor"`
	TotalTowers           *string `json:"total_towers"`
	UpdatedDate           *int64  `json:"updated_date"`
	USP                   *string `json:"usp"`
	VideoPara             *string `json:"video_para"`
	WhyPara               *string `json:"why_para"`
	PropertyConfigTypeID  *int64  `json:"property_configuration_type_id"`
	DeveloperID           *int64  `json:"developer_id"`
	LocalityID            *int64  `json:"locality_id"`
	UserID                *int64  `json:"user_id"`
	MetaKeywords          *string `json:"meta_keywords"`
}

type LProperty struct {
	ID                int64   `json:"id"`
	About             *string `json:"about"`
	AgeOfProperty     *string `json:"age_of_property"`
	AmenitiesPara     *string `json:"amenities_para"`
	Balcony           *string `json:"balcony"`
	Bathrooms         *string `json:"bathrooms"`
	Bedrooms          *string `json:"bedrooms"`
	BuiltupArea       *string `json:"builtup_area"`
	CoverPhoto        *string `json:"cover_photo"`
	CoveredParking    *string `json:"covered_parking"`
	CreatedDate       *int64  `json:"created_date"`
	Facing            *string `json:"facing"`
	FloorImage2D      *string `json:"floor_image2d"`
	FloorImage3D      *string `json:"floor_image3d"`
	FloorPara         *string `json:"floor_para"`
	Floors            *string `json:"floors"`
	FurnishingType    *string `json:"furnishing_type"`
	Images            *string `json:"images"`
	IsDeleted         bool    `json:"is_deleted"`
	IsFeatured        bool    `json:"is_featured"`
	Latlong           *string `json:"latlong"`
	ListingType       *string `json:"listing_type"`
	LocationMap       *string `json:"location_map"`
	LocationAdvantage *string `json:"location_advantage"`
	LocationPara      *string `json:"location_para"`
	LogoImage         *string `json:"logo_image"`
	MetaDescription   *string `json:"meta_description"`
	MetaKeywords      *string `json:"meta_keywords"`
	MetaTitle         *string `json:"meta_title"`
	OpenParking       *string `json:"open_parking"`
	OverviewPara      *string `json:"overview_para"`
	PossessionDate    *string `json:"possession_date"`
	PossessionStatus  *string `json:"possession_status"`
	Price             float64 `json:"price"`
	ProductSchema     *string `json:"product_schema"`
	PropertyAddress   *string `json:"property_address"`
	PropertyName      *string `json:"property_name"`
	PropertyURL       *string `json:"property_url"`
	PropertyVideo     *string `json:"property_video"`
	Rera              *string `json:"rera"`
	Size              *string `json:"size"`
	UpdatedByID       *int64  `json:"updated_by_id"`
	UpdatedDate       *int64  `json:"updated_date"`
	USP               *string `json:"usp"`
	VideoPara         *string `json:"video_para"`
	ConfigurationID   *int64  `json:"confiuration_id"`
	DeveloperID       *int64  `json:"developer_id"`
	LocalityID        *int64  `json:"locality_id"`
	ProjectID         *int64  `json:"project_id"`
	UserID            *int64  `json:"user_id"`
	Highlights        *string `json:"highlights"`
	LocaionMap        *string `json:"locaion_map"`
}

type LDeveloper struct {
	ID                 int64   `json:"id"`
	About              *string `json:"about"`
	AltDeveloperLogo   *string `json:"alt_developer_logo"`
	CreatedDate        *int64  `json:"created_date"`
	DeveloperAddress   *string `json:"developer_address"`
	DeveloperLegalName *string `json:"developer_legal_name"`
	DeveloperLogo      *string `json:"developer_logo"`
	DeveloperName      *string `json:"developer_name"`
	DeveloperURL       *string `json:"developer_url"`
	Disclaimer         *string `json:"disclaimer"`
	EstablishedYear    *int64  `json:"established_year"`
	IsActive           *bool   `json:"is_active"`
	IsVerified         *bool   `json:"is_verified"`
	Overview           *string `json:"overview"`
	ProjectDoneNo      *string `json:"project_done_no"`
	UpdatedDate        *int64  `json:"updated_date"`
	CityName           *int64  `json:"city_name"`
	Phone              *string `json:"phone"`
}

type LLocality struct {
	ID          int64   `json:"id"`
	CreatedDate *int64  `json:"created_date"`
	Name        *string `json:"locality_name"`
	URL         *string `json:"locality_url"`
	UpdatedDate *int64  `json:"updated_date"`
	CityID      *int64  `json:"city_id"`
}

type LCity struct {
	ID          int64   `json:"id"`
	Name        *string `json:"city_name"`
	URL         *string `json:"city_url"`
	CreatedDate *int64  `json:"created_date"`
	IsActive    *int64  `json:"is_active"`
	StateName   *string `json:"state_name"`
	UpdatedDate *int64  `json:"updated_date"`
	Phone       *string `json:"phone_number"`
}

type LPropertyConfiguration struct {
	ID                       int64   `json:"id"`
	CreatedDate              *int64  `json:"created_date"`
	ProjectConfigurationName *string `json:"project_configuration_name"`
	UpdatedDate              *int64  `json:"updated_date"`
	ConfigurationTypeID      *int64  `json:"configuration_type_id"`
}
type LProjectImage struct {
	ProjectID    int64   `json:"project_id"`
	ImageAltName *string `json:"image_alt_name"`
	ImageURL     string  `json:"image_url"`
}

// Floor Plan table structure - matches MySQL floorplan table
type LFloorPlan struct {
	ID              int64   `json:"id"`
	CreatedDate     *int64  `json:"created_date"`
	ImgURL          *string `json:"img_url"`
	IsSoldOut       bool    `json:"is_sold_out"`
	Price           float64 `json:"price"`
	Size            *int64  `json:"size"`
	Title           *string `json:"title"`
	UpdatedDate     *int64  `json:"updated_date"`
	ConfigurationID *int64  `json:"configuration_id"`
	ProjectID       *int64  `json:"project_id"`
	UserID          *int64  `json:"user_id"`
}
type LRera struct {
	ID              int64   `json:"id"`
	CreatedDate     *int64  `json:"created_on"`
	Phase           *string `json:"phase"`
	ProjectReraName *string `json:"project_rera_name"`
	QRImages        *string `json:"qr_images"`
	ReraNumber      *string `json:"rera_number"`
	Status          *string `json:"status"`
	UpdatedDate     *int64  `json:"updated_on"`
	ProjectID       *int64  `json:"project_id"`
	UserID          *int64  `json:"user_id"`
}
type LAmenity struct {
	ID                int64   `json:"id"`
	AmenitiesCategory *string `json:"amenities_category"`
	AmenitiesName     *string `json:"amenities_name"`
	AmenitiesURL      *string `json:"amenities_url"`
	CreatedDate       *int64  `json:"created_date"`
	UpdatedDate       *int64  `json:"updated_date"`
}

type LPaymentPlan struct {
	ID               int64   `json:"id"`
	PaymentPlanName  *string `json:"payment_plan_name"`
	PaymentPlanValue *string `json:"payment_plan_value"`
	ProjectID        *int64  `json:"project_id"`
}

type LFAQ struct {
	ID        int64   `json:"id"`
	Question  *string `json:"question"`
	Answer    *string `json:"answer"`
	ProjectID *int64  `json:"project_id"`
}

type LPropertyConfigurationType struct {
	ID                    int64   `json:"id"`
	ConfigurationTypeName *string `json:"configuration_type_name"`
	CreatedDate           *int64  `json:"created_date"`
	PropertyType          *string `json:"property_type"`
	UpdatedDate           *int64  `json:"updated_date"`
}

type LBlog struct {
	ID           int64    `json:"id"`
	Alt          *string  `json:"alt"`
	BlogSchema   []string `json:"blog_schema"`
	BlogURL      *string  `json:"blog_url"`
	Canonical    *string  `json:"canonical"`
	CreatedDate  *int64   `json:"created_date"`
	Description  *string  `json:"description"`
	Headings     *string  `json:"headings"`
	Images       *string  `json:"images"`
	IsPriority   bool     `json:"is_priority"`
	SubHeadings  *string  `json:"sub_headings"`
	UpdatedDate  *int64   `json:"updated_date"`
	UserID       *int64   `json:"user_id"`
	MetaKeywords *string  `json:"meta_keywords"`
	MetaTitle    *string  `json:"meta_title"`
	IsDeleted    bool     `json:"is_deleted"`
}

type JavaProject struct {
	ID        string `json:"id"`
	VideoURL  string `json:"videoUrl"`
	ProjectID string `json:"projectId"`
}

type JavaGetProjectByIDResponse struct {
	Content []struct {
		ID     int64    `json:"id"`
		Videos []string `json:"videos"`
	} `json:"content"`
}

const (
	// Update these with actual Java API endpoints
	javaAPIBaseURL     = "https://api.investmango.com"
	getAllProjectsPath = "/project/get/all?size=400&isDeleted=false"
)
