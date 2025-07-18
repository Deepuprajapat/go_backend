package migration_jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/VI-IM/im_backend_go/shared/logger"
)

// JSONDataLoader handles loading and querying JSON data
type JSONDataLoader struct {
	dataDir string
	mutex   sync.RWMutex

	// In-memory data structures
	cities                     map[int64]LCity
	localities                 map[int64]LLocality
	developers                 map[int64]LDeveloper
	projects                   map[int64]LProject
	properties                 map[int64]LProperty
	propertyConfigurations     map[int64]LPropertyConfiguration
	propertyConfigurationTypes map[int64]LPropertyConfigurationType
	projectImages              map[int64][]LProjectImage
	floorPlans                 map[int64][]LFloorPlan
	reras                      map[int64][]LRera
	amenities                  map[int64]LAmenity
	projectAmenities           map[int64][]LProjectAmenity
	paymentPlans               map[int64][]LPaymentPlan
	faqs                       map[int64][]LFAQ
	blogs                      map[int64]LBlog
	genericSearchData          map[int64]LGenericSearchData

	// Secondary indexes
	projectsByDeveloper map[int64][]int64
	projectsByLocality  map[int64][]int64
	propertiesByProject map[int64][]int64

	isLoaded bool
}

// TableData represents the structure of exported JSON files
type TableData struct {
	TableName   string                   `json:"table_name"`
	ExportedAt  string                   `json:"exported_at"`
	RowCount    int                      `json:"row_count"`
	Columns     []string                 `json:"columns"`
	ColumnTypes []string                 `json:"column_types"`
	Data        []map[string]interface{} `json:"data"`
}

// Global instance
var jsonLoader *JSONDataLoader

// InitializeJSONDataLoader initializes the JSON data loader with the specified directory
func InitializeJSONDataLoader(dataDir string) error {
	jsonLoader = &JSONDataLoader{
		dataDir:                    dataDir,
		cities:                     make(map[int64]LCity),
		localities:                 make(map[int64]LLocality),
		developers:                 make(map[int64]LDeveloper),
		projects:                   make(map[int64]LProject),
		properties:                 make(map[int64]LProperty),
		propertyConfigurations:     make(map[int64]LPropertyConfiguration),
		propertyConfigurationTypes: make(map[int64]LPropertyConfigurationType),
		projectImages:              make(map[int64][]LProjectImage),
		floorPlans:                 make(map[int64][]LFloorPlan),
		reras:                      make(map[int64][]LRera),
		amenities:                  make(map[int64]LAmenity),
		projectAmenities:           make(map[int64][]LProjectAmenity),
		paymentPlans:               make(map[int64][]LPaymentPlan),
		faqs:                       make(map[int64][]LFAQ),
		blogs:                      make(map[int64]LBlog),
		genericSearchData:          make(map[int64]LGenericSearchData),
		projectsByDeveloper:        make(map[int64][]int64),
		projectsByLocality:         make(map[int64][]int64),
		propertiesByProject:        make(map[int64][]int64),
	}

	return jsonLoader.LoadAllData()
}

// LoadAllData loads all JSON files into memory
func (j *JSONDataLoader) LoadAllData() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	logger.Get().Info().Msg("Loading JSON data into memory...")

	// Load each table
	tables := []string{
		"amenities", "blogs", "city", "developer", "faq", "floorplan",
		"generic_search", "locality", "payment_plan", "project",
		"project_amenities", "project_configuration", "project_configuration_type",
		"project_image", "property", "rera_info",
	}

	for _, tableName := range tables {
		if err := j.loadTable(tableName); err != nil {
			logger.Get().Error().Err(err).Msgf("Failed to load table %s", tableName)
			// Continue loading other tables even if one fails
		}
	}

	// Build secondary indexes
	j.buildIndexes()

	j.isLoaded = true
	logger.Get().Info().Msg("JSON data loaded successfully")

	return nil
}

// loadTable loads a specific table from JSON file
func (j *JSONDataLoader) loadTable(tableName string) error {
	filePath := filepath.Join(j.dataDir, tableName+".json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var tableData TableData
	if err := json.Unmarshal(data, &tableData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	// Process data based on table name
	switch tableName {
	case "city":
		return j.loadCities(tableData.Data)
	case "locality":
		return j.loadLocalities(tableData.Data)
	case "developer":
		return j.loadDevelopers(tableData.Data)
	case "project":
		return j.loadProjects(tableData.Data)
	case "property":
		return j.loadProperties(tableData.Data)
	case "project_configuration":
		return j.loadPropertyConfigurations(tableData.Data)
	case "project_configuration_type":
		return j.loadPropertyConfigurationTypes(tableData.Data)
	case "project_image":
		return j.loadProjectImages(tableData.Data)
	case "floorplan":
		return j.loadFloorPlans(tableData.Data)
	case "rera_info":
		return j.loadReras(tableData.Data)
	case "amenities":
		return j.loadAmenities(tableData.Data)
	case "project_amenities":
		return j.loadProjectAmenities(tableData.Data)
	case "payment_plan":
		return j.loadPaymentPlans(tableData.Data)
	case "faq":
		return j.loadFAQs(tableData.Data)
	case "blogs":
		return j.loadBlogs(tableData.Data)
	case "generic_search":
		return j.loadGenericSearchData(tableData.Data)
	default:
		logger.Get().Warn().Msgf("Unknown table: %s", tableName)
	}

	return nil
}

// Helper functions to load specific data types
func (j *JSONDataLoader) loadCities(data []map[string]interface{}) error {
	for _, row := range data {
		city := LCity{
			ID:          getInt64(row["id"]),
			Name:        getStringPtr(row["city_name"]),
			URL:         getStringPtr(row["city_url"]),
			CreatedDate: getInt64Ptr(row["created_date"]),
			IsActive:    getInt64Ptr(row["is_active"]),
			StateName:   getStringPtr(row["state_name"]),
			UpdatedDate: getInt64Ptr(row["updated_date"]),
			Phone:       getStringPtr(row["phone_number"]),
		}
		j.cities[city.ID] = city
	}
	return nil
}

func (j *JSONDataLoader) loadLocalities(data []map[string]interface{}) error {
	for _, row := range data {
		locality := LLocality{
			ID:          getInt64(row["id"]),
			CreatedDate: getInt64Ptr(row["created_date"]),
			Name:        getStringPtr(row["locality_name"]),
			URL:         getStringPtr(row["locality_url"]),
			UpdatedDate: getInt64Ptr(row["updated_date"]),
			CityID:      getInt64Ptr(row["city_id"]),
		}
		j.localities[locality.ID] = locality
	}
	return nil
}

func (j *JSONDataLoader) loadDevelopers(data []map[string]interface{}) error {
	for _, row := range data {
		developer := LDeveloper{
			ID:                 getInt64(row["id"]),
			About:              getStringPtr(row["about"]),
			AltDeveloperLogo:   getStringPtr(row["alt_developer_logo"]),
			CreatedDate:        getInt64Ptr(row["created_date"]),
			DeveloperAddress:   getStringPtr(row["developer_address"]),
			DeveloperLegalName: getStringPtr(row["developer_legal_name"]),
			DeveloperLogo:      getStringPtr(row["developer_logo"]),
			DeveloperName:      getStringPtr(row["developer_name"]),
			DeveloperURL:       getStringPtr(row["developer_url"]),
			Disclaimer:         getStringPtr(row["disclaimer"]),
			EstablishedYear:    getInt64Ptr(row["established_year"]),
			IsActive:           getBoolPtr(row["is_active"]),
			IsVerified:         getBoolPtr(row["is_verified"]),
			Overview:           getStringPtr(row["overview"]),
			ProjectDoneNo:      getStringPtr(row["project_done_no"]),
			UpdatedDate:        getInt64Ptr(row["updated_date"]),
			CityName:           getInt64Ptr(row["city_name"]),
			Phone:              getStringPtr(row["phone"]),
		}
		j.developers[developer.ID] = developer
	}
	return nil
}

func (j *JSONDataLoader) loadProjects(data []map[string]interface{}) error {
	for _, row := range data {
		project := LProject{
			ID:                    getInt64(row["id"]),
			AltProjectLogo:        getStringPtr(row["alt_project_logo"]),
			AltSitePlanImg:        getStringPtr(row["alt_site_plan_img"]),
			AmenitiesPara:         getStringPtr(row["amenities_para"]),
			AvailableUnit:         getStringPtr(row["available_unit"]),
			CoverPhoto:            getStringPtr(row["cover_photo"]),
			CreatedDate:           getInt64Ptr(row["created_date"]),
			FloorPara:             getStringPtr(row["floor_para"]),
			IsDeleted:             getBool(row["is_deleted"]),
			IsFeatured:            getBool(row["is_featured"]),
			IsPremium:             getBool(row["is_premium"]),
			IsPriority:            getBool(row["is_priority"]),
			LocationMap:           getStringPtr(row["location_map"]),
			LocationPara:          getStringPtr(row["location_para"]),
			MetaDescription:       getStringPtr(row["meta_description"]),
			MetaTitle:             getStringPtr(row["meta_title"]),
			OverviewPara:          getStringPtr(row["overview_para"]),
			PaymentPara:           getStringPtr(row["payment_para"]),
			PriceListPara:         getStringPtr(row["price_list_para"]),
			ProjectAbout:          getStringPtr(row["project_about"]),
			ProjectAddress:        getStringPtr(row["project_address"]),
			ProjectArea:           getStringPtr(row["project_area"]),
			ProjectBrochure:       getStringPtr(row["project_brochure"]),
			ProjectConfigurations: getStringPtr(row["project_configurations"]),
			ProjectDescription:    getStringPtr(row["project_description"]),
			ProjectLaunchDate:     getStringPtr(row["project_launch_date"]),
			ProjectLocationURL:    getStringPtr(row["project_location_url"]),
			ProjectLogo:           getStringPtr(row["project_logo"]),
			ProjectName:           getStringPtr(row["project_name"]),
			ProjectPossessionDate: getStringPtr(row["project_possession_date"]),
			ProjectRERA:           getStringPtr(row["project_rera"]),
			ProjectSchema:         getStringPtr(row["project_schema"]),
			ProjectUnits:          getStringPtr(row["project_units"]),
			ProjectURL:            getStringPtr(row["project_url"]),
			ProjectVideoCount:     getInt64Ptr(row["project_video_count"]),
			ProjectVideos:         getBytes(row["project_videos"]),
			ReraLink:              getStringPtr(row["rera_link"]),
			ShortAddress:          getStringPtr(row["short_address"]),
			SitePlanImg:           getStringPtr(row["siteplan_img"]),
			SitePlanPara:          getStringPtr(row["siteplan_para"]),
			Status:                getStringPtr(row["status"]),
			TotalFloor:            getStringPtr(row["total_floor"]),
			TotalTowers:           getStringPtr(row["total_towers"]),
			UpdatedDate:           getInt64Ptr(row["updated_date"]),
			USP:                   getStringPtr(row["usp"]),
			VideoPara:             getStringPtr(row["video_para"]),
			WhyPara:               getStringPtr(row["why_para"]),
			PropertyConfigTypeID:  getInt64Ptr(row["property_configuration_type_id"]),
			DeveloperID:           getInt64Ptr(row["developer_id"]),
			LocalityID:            getInt64Ptr(row["locality_id"]),
			UserID:                getInt64Ptr(row["user_id"]),
			MetaKeywords:          getStringPtr(row["meta_keywords"]),
		}
		j.projects[project.ID] = project
	}
	return nil
}

func (j *JSONDataLoader) loadProperties(data []map[string]interface{}) error {
	for _, row := range data {
		property := LProperty{
			ID:                getInt64(row["id"]),
			About:             getStringPtr(row["about"]),
			AgeOfProperty:     getStringPtr(row["age_of_property"]),
			AmenitiesPara:     getStringPtr(row["amenities_para"]),
			Balcony:           getStringPtr(row["balcony"]),
			Bathrooms:         getStringPtr(row["bathrooms"]),
			Bedrooms:          getStringPtr(row["bedrooms"]),
			BuiltupArea:       getStringPtr(row["builtup_area"]),
			CoverPhoto:        getStringPtr(row["cover_photo"]),
			CoveredParking:    getStringPtr(row["covered_parking"]),
			CreatedDate:       getInt64Ptr(row["created_date"]),
			Facing:            getStringPtr(row["facing"]),
			FloorImage2D:      getStringPtr(row["floor_image2d"]),
			FloorImage3D:      getStringPtr(row["floor_image3d"]),
			FloorPara:         getStringPtr(row["floor_para"]),
			Floors:            getStringPtr(row["floors"]),
			FurnishingType:    getStringPtr(row["furnishing_type"]),
			Images:            getStringPtr(row["images"]),
			IsDeleted:         getBool(row["is_deleted"]),
			IsFeatured:        getBool(row["is_featured"]),
			Latlong:           getStringPtr(row["latlong"]),
			ListingType:       getStringPtr(row["listing_type"]),
			LocationMap:       getStringPtr(row["location_map"]),
			LocationAdvantage: getStringPtr(row["location_advantage"]),
			LocationPara:      getStringPtr(row["location_para"]),
			LogoImage:         getStringPtr(row["logo_image"]),
			MetaDescription:   getStringPtr(row["meta_description"]),
			MetaKeywords:      getStringPtr(row["meta_keywords"]),
			MetaTitle:         getStringPtr(row["meta_title"]),
			OpenParking:       getStringPtr(row["open_parking"]),
			OverviewPara:      getStringPtr(row["overview_para"]),
			PossessionDate:    getStringPtr(row["possession_date"]),
			PossessionStatus:  getStringPtr(row["possession_status"]),
			Price:             getFloat64(row["price"]),
			ProductSchema:     getStringPtr(row["product_schema"]),
			PropertyAddress:   getStringPtr(row["property_address"]),
			PropertyName:      getStringPtr(row["property_name"]),
			PropertyURL:       getStringPtr(row["property_url"]),
			PropertyVideo:     getStringPtr(row["property_video"]),
			Rera:              getStringPtr(row["rera"]),
			Size:              getStringPtr(row["size"]),
			UpdatedByID:       getInt64Ptr(row["updated_by_id"]),
			UpdatedDate:       getInt64Ptr(row["updated_date"]),
			USP:               getStringPtr(row["usp"]),
			VideoPara:         getStringPtr(row["video_para"]),
			ConfigurationID:   getInt64Ptr(row["confiuration_id"]),
			DeveloperID:       getInt64Ptr(row["developer_id"]),
			LocalityID:        getInt64Ptr(row["locality_id"]),
			ProjectID:         getInt64Ptr(row["project_id"]),
			UserID:            getInt64Ptr(row["user_id"]),
			Highlights:        getStringPtr(row["highlights"]),
			LocaionMap:        getStringPtr(row["locaion_map"]),
		}
		j.properties[property.ID] = property
	}
	return nil
}

func (j *JSONDataLoader) loadPropertyConfigurations(data []map[string]interface{}) error {
	for _, row := range data {
		config := LPropertyConfiguration{
			ID:                       getInt64(row["id"]),
			CreatedDate:              getInt64Ptr(row["created_date"]),
			ProjectConfigurationName: getStringPtr(row["project_configuration_name"]),
			UpdatedDate:              getInt64Ptr(row["updated_date"]),
			ConfigurationTypeID:      getInt64Ptr(row["configuration_type_id"]),
		}
		j.propertyConfigurations[config.ID] = config
	}
	return nil
}

func (j *JSONDataLoader) loadPropertyConfigurationTypes(data []map[string]interface{}) error {
	for _, row := range data {
		configType := LPropertyConfigurationType{
			ID:                    getInt64(row["id"]),
			ConfigurationTypeName: getStringPtr(row["configuration_type_name"]),
			CreatedDate:           getInt64Ptr(row["created_date"]),
			PropertyType:          getStringPtr(row["property_type"]),
			UpdatedDate:           getInt64Ptr(row["updated_date"]),
		}
		j.propertyConfigurationTypes[configType.ID] = configType
	}
	return nil
}

func (j *JSONDataLoader) loadProjectImages(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		image := LProjectImage{
			ProjectID:    projectID,
			ImageAltName: getStringPtr(row["image_alt_name"]),
			ImageURL:     getString(row["image_url"]),
		}
		j.projectImages[projectID] = append(j.projectImages[projectID], image)
	}
	return nil
}

func (j *JSONDataLoader) loadFloorPlans(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		floorPlan := LFloorPlan{
			ID:              getInt64(row["id"]),
			CreatedDate:     getInt64Ptr(row["created_date"]),
			ImgURL:          getStringPtr(row["img_url"]),
			IsSoldOut:       getBool(row["is_sold_out"]),
			Price:           getFloat64(row["price"]),
			Size:            getInt64Ptr(row["size"]),
			Title:           getStringPtr(row["title"]),
			UpdatedDate:     getInt64Ptr(row["updated_date"]),
			ConfigurationID: getInt64Ptr(row["configuration_id"]),
			ProjectID:       getInt64Ptr(row["project_id"]),
			UserID:          getInt64Ptr(row["user_id"]),
		}
		j.floorPlans[projectID] = append(j.floorPlans[projectID], floorPlan)
	}
	return nil
}

func (j *JSONDataLoader) loadReras(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		rera := LRera{
			ID:              getInt64(row["id"]),
			CreatedDate:     getInt64Ptr(row["created_on"]),
			Phase:           getStringPtr(row["phase"]),
			ProjectReraName: getStringPtr(row["project_rera_name"]),
			QRImages:        getStringPtr(row["qr_images"]),
			ReraNumber:      getStringPtr(row["rera_number"]),
			Status:          getStringPtr(row["status"]),
			UpdatedDate:     getInt64Ptr(row["updated_on"]),
			ProjectID:       getInt64Ptr(row["project_id"]),
			UserID:          getInt64Ptr(row["user_id"]),
		}
		j.reras[projectID] = append(j.reras[projectID], rera)
	}
	return nil
}

func (j *JSONDataLoader) loadAmenities(data []map[string]interface{}) error {
	for _, row := range data {
		amenity := LAmenity{
			ID:                getInt64(row["id"]),
			AmenitiesCategory: getStringPtr(row["amenities_category"]),
			AmenitiesName:     getStringPtr(row["amenities_name"]),
			AmenitiesURL:      getStringPtr(row["amenities_url"]),
			CreatedDate:       getInt64Ptr(row["created_date"]),
			UpdatedDate:       getInt64Ptr(row["updated_date"]),
		}
		j.amenities[amenity.ID] = amenity
	}
	return nil
}

func (j *JSONDataLoader) loadProjectAmenities(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		amenityID := getInt64(row["amenities_id"])
		projectAmenity := LProjectAmenity{
			ProjectID: projectID,
			AmenityID: amenityID,
		}
		j.projectAmenities[projectID] = append(j.projectAmenities[projectID], projectAmenity)
	}
	return nil
}

func (j *JSONDataLoader) loadPaymentPlans(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		paymentPlan := LPaymentPlan{
			ID:               getInt64(row["id"]),
			PaymentPlanName:  getStringPtr(row["payment_plan_name"]),
			PaymentPlanValue: getStringPtr(row["payment_plan_value"]),
			ProjectID:        getInt64Ptr(row["project_id"]),
		}
		j.paymentPlans[projectID] = append(j.paymentPlans[projectID], paymentPlan)
	}
	return nil
}

func (j *JSONDataLoader) loadFAQs(data []map[string]interface{}) error {
	for _, row := range data {
		projectID := getInt64(row["project_id"])
		faq := LFAQ{
			ID:        getInt64(row["id"]),
			Question:  getStringPtr(row["question"]),
			Answer:    getStringPtr(row["answer"]),
			ProjectID: getInt64Ptr(row["project_id"]),
		}
		j.faqs[projectID] = append(j.faqs[projectID], faq)
	}
	return nil
}

func (j *JSONDataLoader) loadBlogs(data []map[string]interface{}) error {
	for _, row := range data {
		blog := LBlog{
			ID:           getInt64(row["id"]),
			Alt:          getStringPtr(row["alt"]),
			BlogSchema:   getStringPtr(row["blog_schema"]),
			BlogURL:      getStringPtr(row["blog_url"]),
			Canonical:    getStringPtr(row["canonical"]),
			CreatedDate:  getInt64Ptr(row["created_date"]),
			Description:  getStringPtr(row["description"]),
			Headings:     getStringPtr(row["headings"]),
			Images:       getStringPtr(row["images"]),
			IsPriority:   getBool(row["is_priority"]),
			SubHeadings:  getStringPtr(row["sub_headings"]),
			UpdatedDate:  getInt64Ptr(row["updated_date"]),
			UserID:       getInt64Ptr(row["user_id"]),
			MetaKeywords: getStringPtr(row["meta_keywords"]),
			MetaTitle:    getStringPtr(row["meta_title"]),
			IsDeleted:    getBool(row["is_deleted"]),
		}
		j.blogs[blog.ID] = blog
	}
	return nil
}

func (j *JSONDataLoader) loadGenericSearchData(data []map[string]interface{}) error {
	for _, row := range data {
		searchData := LGenericSearchData{
			ID:          getInt64(row["id"]),
			Path:        getStringPtr(row["path"]),
			SearchTerms: getStringPtr(row["search_terms"]),
			URL:         getStringPtr(row["url"]),
		}
		j.genericSearchData[searchData.ID] = searchData
	}
	return nil
}

// buildIndexes creates secondary indexes for efficient querying
func (j *JSONDataLoader) buildIndexes() {
	// Index projects by developer
	for _, project := range j.projects {
		if project.DeveloperID != nil {
			developerID := *project.DeveloperID
			j.projectsByDeveloper[developerID] = append(j.projectsByDeveloper[developerID], project.ID)
		}
	}

	// Index projects by locality
	for _, project := range j.projects {
		if project.LocalityID != nil {
			localityID := *project.LocalityID
			j.projectsByLocality[localityID] = append(j.projectsByLocality[localityID], project.ID)
		}
	}

	// Index properties by project
	for _, property := range j.properties {
		if property.ProjectID != nil {
			projectID := *property.ProjectID
			j.propertiesByProject[projectID] = append(j.propertiesByProject[projectID], property.ID)
		}
	}
}

// Helper functions for type conversion
func getInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func getInt64Ptr(v interface{}) *int64 {
	if v == nil {
		return nil
	}
	val := getInt64(v)
	return &val
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func getStringPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	val := getString(v)
	return &val
}

func getBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int64:
		return val != 0
	case int:
		return val != 0
	case float64:
		return val != 0
	case string:
		return strings.ToLower(val) == "true" || val == "1"
	}
	return false
}

func getBoolPtr(v interface{}) *bool {
	if v == nil {
		return nil
	}
	val := getBool(v)
	return &val
}

func getFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}

func getBytes(v interface{}) []byte {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return []byte(s)
	}
	return nil
}

// Query methods - these replace the database queries
func (j *JSONDataLoader) GetCityByID(id int64) (LCity, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return LCity{}, fmt.Errorf("data not loaded")
	}

	city, exists := j.cities[id]
	if !exists {
		return LCity{}, fmt.Errorf("no city found with ID %d", id)
	}

	return city, nil
}

func (j *JSONDataLoader) GetAllLocalities() ([]LLocality, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var localities []LLocality
	for _, locality := range j.localities {
		localities = append(localities, locality)
	}

	return localities, nil
}

func (j *JSONDataLoader) GetAllDevelopers() ([]LDeveloper, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var developers []LDeveloper
	for _, developer := range j.developers {
		developers = append(developers, developer)
	}

	return developers, nil
}

func (j *JSONDataLoader) GetAllProjects() ([]LProject, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var projects []LProject
	for _, project := range j.projects {
		projects = append(projects, project)
	}

	return projects, nil
}

func (j *JSONDataLoader) GetAllProperties() ([]LProperty, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var properties []LProperty
	for _, property := range j.properties {
		properties = append(properties, property)
	}

	return properties, nil
}

func (j *JSONDataLoader) GetProjectByID(id int64) (*LProject, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	project, exists := j.projects[id]
	if !exists {
		return nil, fmt.Errorf("no project found with ID %d", id)
	}

	return &project, nil
}

func (j *JSONDataLoader) GetDeveloperByID(id int64) (*LDeveloper, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	developer, exists := j.developers[id]
	if !exists {
		return nil, fmt.Errorf("no developer found with ID %d", id)
	}

	return &developer, nil
}

func (j *JSONDataLoader) GetLocalityByID(id int64) (*LLocality, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	locality, exists := j.localities[id]
	if !exists {
		return nil, fmt.Errorf("no locality found with ID %d", id)
	}

	return &locality, nil
}

func (j *JSONDataLoader) GetPropertyConfigurationByID(id int64) (*LPropertyConfiguration, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	config, exists := j.propertyConfigurations[id]
	if !exists {
		return nil, fmt.Errorf("no property configuration found with ID %d", id)
	}

	return &config, nil
}

func (j *JSONDataLoader) GetPropertyConfigurationTypeByID(id int64) (*LPropertyConfigurationType, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	configType, exists := j.propertyConfigurationTypes[id]
	if !exists {
		return nil, fmt.Errorf("no property configuration type found with ID %d", id)
	}

	return &configType, nil
}

func (j *JSONDataLoader) GetProjectImagesByProjectID(id int64) (*[]LProjectImage, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	images, exists := j.projectImages[id]
	if !exists {
		return &[]LProjectImage{}, nil // Return empty slice if no images found
	}

	return &images, nil
}

func (j *JSONDataLoader) GetFloorPlansByProjectID(id int64) (*[]LFloorPlan, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	floorPlans, exists := j.floorPlans[id]
	if !exists {
		return &[]LFloorPlan{}, nil
	}

	return &floorPlans, nil
}

func (j *JSONDataLoader) GetRerasByProjectID(id int64) ([]*LRera, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	reras, exists := j.reras[id]
	if !exists {
		return []*LRera{}, nil
	}

	// Convert to pointer slice
	var result []*LRera
	for i := range reras {
		result = append(result, &reras[i])
	}

	return result, nil
}

func (j *JSONDataLoader) GetAmenityByID(id int64) (*LAmenity, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	amenity, exists := j.amenities[id]
	if !exists {
		return nil, fmt.Errorf("no amenity found with ID %d", id)
	}

	return &amenity, nil
}

func (j *JSONDataLoader) GetProjectAmenitiesByProjectID(id int64) ([]*LAmenity, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	projectAmenities, exists := j.projectAmenities[id]
	if !exists {
		return []*LAmenity{}, nil
	}

	var amenities []*LAmenity
	for _, pa := range projectAmenities {
		if amenity, exists := j.amenities[pa.AmenityID]; exists {
			amenities = append(amenities, &amenity)
		}
	}

	return amenities, nil
}

func (j *JSONDataLoader) GetPaymentPlansByProjectID(id int64) ([]*LPaymentPlan, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	paymentPlans, exists := j.paymentPlans[id]
	if !exists {
		return []*LPaymentPlan{}, nil
	}

	var result []*LPaymentPlan
	for i := range paymentPlans {
		result = append(result, &paymentPlans[i])
	}

	return result, nil
}

func (j *JSONDataLoader) GetFaqsByProjectID(id int64) ([]*LFAQ, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	faqs, exists := j.faqs[id]
	if !exists {
		return []*LFAQ{}, nil
	}

	var result []*LFAQ
	for i := range faqs {
		result = append(result, &faqs[i])
	}

	return result, nil
}

func (j *JSONDataLoader) GetAllBlogs() ([]LBlog, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var blogs []LBlog
	for _, blog := range j.blogs {
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (j *JSONDataLoader) GetAllGenericSearchData() ([]LGenericSearchData, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if !j.isLoaded {
		return nil, fmt.Errorf("data not loaded")
	}

	var data []LGenericSearchData
	for _, item := range j.genericSearchData {
		data = append(data, item)
	}

	return data, nil
}

// GetJSONDataLoader returns the global JSON data loader instance
func GetJSONDataLoader() *JSONDataLoader {
	return jsonLoader
}
