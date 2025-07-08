package migration_jobs

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/VI-IM/im_backend_go/ent"
	entproject "github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	batchSize = 200 // Number of items to process in each batch
)

var (
	// Shared maps protected by mutex
	mapMutex                  sync.RWMutex
	legacyToNewProjectIDMAP   = make(map[int64]string)
	legacyToNewDeveloperIDMAP = make(map[int64]string)
	legacyToNewLocalityIDMAP  = make(map[int64]string)
	projectIDToVideoURLMAP    = make(map[int64][]string)
	staticSiteWebCardsID      string
)

// Helper functions to safely access maps
func setProjectIDMapping(legacyID int64, newID string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	legacyToNewProjectIDMAP[legacyID] = newID
}

func getProjectIDMapping(legacyID int64) (*string, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	id, ok := legacyToNewProjectIDMAP[legacyID]
	return &id, ok
}

func setDeveloperIDMapping(legacyID int64, newID string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	legacyToNewDeveloperIDMAP[legacyID] = newID
}

func getDeveloperIDMapping(legacyID int64) (*string, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	id, ok := legacyToNewDeveloperIDMAP[legacyID]
	return &id, ok
}

func setLocalityIDMapping(legacyID int64, newID string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	legacyToNewLocalityIDMAP[legacyID] = newID
}

func getLocalityIDMapping(legacyID int64) (*string, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	id, ok := legacyToNewLocalityIDMAP[legacyID]
	return &id, ok
}

// processBatch processes a batch of items with the given processor function
func processBatch[T any](ctx context.Context, items []T, batchSize int, processor func(context.Context, []T) error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, (len(items)+batchSize-1)/batchSize)

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		wg.Add(1)
		go func(batch []T) {
			defer wg.Done()
			if err := processor(ctx, batch); err != nil {
				errChan <- err
			}
		}(batch)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Use ent schema modal to migrate data from legacy database to new database

func MigrateProject(ctx context.Context, txn *ent.Tx) error {
	projects, err := FetchhAllProject(ctx)
	if err != nil {
		return err
	}

	projectFetchFromAPI, err := fetchAllProjectIDs(&http.Client{})
	if err != nil {
		return err
	}

	for _, project := range projectFetchFromAPI.Content {
		projectIDToVideoURLMAP[project.ID] = project.Videos
	}

	log.Info().Msg("Fetched all projects --------->>>> success")
	processProjectBatch := func(ctx context.Context, batch []LProject) error {
		for _, project := range batch {
			id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(project.ID, 10))))[:16]

			projectRera, err := FetchReraByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch project RERA for project ID %d", project.ID)
				continue
			}

			reras := []schema.ReraListItem{}
			for _, rera := range projectRera {
				reras = append(reras, schema.ReraListItem{
					Phase:      safeStr(rera.Phase),
					ReraQR:     safeStr(rera.QRImages),
					ReraNumber: safeStr(rera.ReraNumber),
					Status:     safeStr(rera.Status),
				})
			}

			uspList := []string{}
			if project.USP != nil {
				uspText := strings.Trim(*project.USP, "[]")
				uspItems := strings.Split(uspText, "\",")
				for _, item := range uspItems {
					item = strings.Trim(item, "\" ")
					if item != "" {
						uspList = append(uspList, item)
					}
				}
			}

			uspListNew := []string{}
			for _, usp := range uspList {
				uspListNew = append(uspListNew, safeStr(&usp))
			}

			floorPlans, err := FetchFloorPlanByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch floor plans for project ID %d", project.ID)
				continue
			}

			var minPrice float64 = -1
			var maxPrice float64 = -1
			var minSize int64 = -1
			var maxSize int64 = -1

			floorPlanItems := []schema.FloorPlanItem{}
			configurationProducts := []schema.ProductConfiguration{}
			for _, floorPlan := range *floorPlans {

				price := floorPlan.Price
				if price > 0 {
					if minPrice == -1 || price < minPrice {
						minPrice = price
					}
					if maxPrice == -1 || price > maxPrice {
						maxPrice = price
					}
				}

				if floorPlan.Size != nil && *floorPlan.Size > 0 {
					size := *floorPlan.Size
					if minSize == -1 || size < minSize {
						minSize = size
					}
					if maxSize == -1 || size > maxSize {
						maxSize = size
					}
				}

				floorPlanItems = append(floorPlanItems, schema.FloorPlanItem{
					Title:        safeStr(floorPlan.Title),
					FlatType:     safeStr(floorPlan.Title),
					IsSoldOut:    floorPlan.IsSoldOut,
					Price:        strconv.FormatFloat(floorPlan.Price, 'f', -1, 64),
					BuildingArea: strconv.FormatInt(*floorPlan.Size, 10),
					Image:        safeStr(floorPlan.ImgURL),
				})
				configurationProducts = append(configurationProducts, schema.ProductConfiguration{
					ConfigurationName: safeStr(floorPlan.Title),
					Size:              strconv.FormatInt(*floorPlan.Size, 10),
					Price:             strconv.FormatFloat(floorPlan.Price, 'f', -1, 64),
				})
			}

			var sizeRange string
			if minSize != -1 && maxSize != -1 {
				if minSize == maxSize {
					sizeRange = fmt.Sprintf("%d sq.ft.", minSize)
				} else {
					sizeRange = fmt.Sprintf("%d - %d sq.ft.", minSize, maxSize)
				}
			} else {
				sizeRange = "N/A"
			}

			// Convert float prices to integers (calculated but not used in current migration)

			amenities, err := FetchProjectAmenitiesByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch amenities for project ID %d", project.ID)
				continue
			}

			amenitiesMap := map[string][]schema.AmenityCategory{}
			for _, amenity := range amenities {
				amenitiesMap[*amenity.AmenitiesCategory] = append(amenitiesMap[*amenity.AmenitiesCategory], schema.AmenityCategory{
					Icon:  safeStr(amenity.AmenitiesURL),
					Value: safeStr(amenity.AmenitiesName),
				})
			}

			paymentPlans, err := FetchPaymentPlansByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch payment plans for project ID %d", project.ID)
				continue
			}

			paymentPlansNew := []schema.Plan{}
			for _, paymentPlan := range paymentPlans {
				paymentPlansNew = append(paymentPlansNew, schema.Plan{
					Name:    safeStr(paymentPlan.PaymentPlanName),
					Details: safeStr(paymentPlan.PaymentPlanValue),
				})
			}

			if project.DeveloperID == nil {
				log.Error().Msgf("Developer ID is nil for project ID %d", project.ID)
				continue
			}

			developer, err := FetchDeveloperByID(ctx, *project.DeveloperID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch developer for project ID %d", project.ID)
				continue
			}

			faqs, err := FetchFaqsByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch faqs for project ID %d", project.ID)
				continue
			}

			projectImages, err := FetchProjectImagesByProjectID(ctx, project.ID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch project images for project ID %d", project.ID)
				continue
			}

			projectImagesNew := []LProjectImage{}
			var isLogoSeeded bool
			for _, image := range *projectImages {
				if !isLogoSeeded {
					projectImagesNew = append(projectImagesNew, LProjectImage{
						ImageURL:     *project.ProjectLogo,
						ImageAltName: project.AltProjectLogo,
					})
					isLogoSeeded = true
				}
				projectImagesNew = append(projectImagesNew, LProjectImage{
					ImageURL:     image.ImageURL,
					ImageAltName: image.ImageAltName,
				})
			}

			imageURLs := make([]string, len(projectImagesNew))
			for i, img := range projectImagesNew {
				imageURLs[i] = img.ImageURL
			}

			faqsNew := []schema.FAQ{}
			for _, faq := range faqs {
				faqsNew = append(faqsNew, schema.FAQ{
					Question: safeStr(faq.Question),
					Answer:   safeStr(faq.Answer),
				})
			}

			setProjectIDMapping(project.ID, id)

			developerID, ok := getDeveloperIDMapping(*project.DeveloperID)
			if !ok {
				log.Error().Msgf("Developer ID mapping not found for project ID %d", project.ID)
				continue
			}

			localityID, ok := getLocalityIDMapping(*project.LocalityID)
			if !ok {
				log.Error().Msgf("Locality ID mapping not found for project ID %d", project.ID)
				continue
			}
			projectConfigurationType, err := FetchProjectConfigurationByID(ctx, *project.PropertyConfigTypeID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch project configuration type for project ID %d", project.ID)
				continue
			}

			videoUrl := projectIDToVideoURLMAP[project.ID]

			if err := txn.Project.Create().
				SetID(id).
				SetName(safeStr(project.ProjectName)).
				SetMinPrice(strconv.FormatFloat(minPrice, 'f', -1, 64)).
				SetMaxPrice(strconv.FormatFloat(maxPrice, 'f', -1, 64)).
				SetDescription(safeStr(project.ProjectDescription)).
				SetStatus(enums.ProjectStatus(*project.Status)).
				SetTimelineInfo(schema.TimelineInfo{
					ProjectLaunchDate:     safeStr(project.ProjectLaunchDate),
					ProjectPossessionDate: safeStr(project.ProjectPossessionDate),
				}).
				SetMetaInfo(schema.SEOMeta{
					Title:         safeStr(project.MetaTitle),
					Description:   safeStr(project.MetaDescription),
					Keywords:      safeStr(project.MetaKeywords),
					Canonical:     safeStr(project.ProjectURL),
					ProjectSchema: project.ProjectSchema,
				}).
				SetWebCards(schema.ProjectWebCards{
					Images: imageURLs,
					ReraInfo: schema.ReraInfo{
						WebsiteLink: safeStr(project.ReraLink),
						ReraList:    reras,
					},
					Details: schema.ProjectDetails{
						Area: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.ProjectArea),
						},
						Sizes: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: sizeRange,
						},
						Units: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.ProjectUnits),
						},
						TotalFloor: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.TotalFloor),
						},
						TotalTowers: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.TotalTowers),
						},
						Configuration: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.ProjectConfigurations),
						},
						LaunchDate: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.ProjectLaunchDate),
						},
						PossessionDate: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(project.ProjectPossessionDate),
						},
						Type: struct {
							Value string `json:"value,omitempty"`
						}{
							Value: safeStr(projectConfigurationType.PropertyType),
						},
					},
					WhyToChoose: schema.WhyToChoose{
						ImageUrls: imageURLs[1:],
						USP_List:  uspListNew,
					},
					KnowAbout: schema.KnowAbout{
						Description:  *project.ProjectAbout,
						DownloadLink: safeStr(project.ProjectBrochure),
					},
					FloorPlan: schema.FloorPlan{
						Description: safeStr(project.PriceListPara),
						Products:    floorPlanItems,
					},
					PriceList: schema.PriceList{
						Description:          safeStr(project.PriceListPara),
						BHKOptionsWithPrices: configurationProducts,
					},
					Amenities: schema.Amenities{
						Description:             safeStr(project.AmenitiesPara),
						CategoriesWithAmenities: amenitiesMap,
					},
					VideoPresentation: schema.VideoPresentation{
						Description: safeStr(project.VideoPara),
						URLs:        videoUrl,
					},
					PaymentPlans: schema.PaymentPlans{
						Description: safeStr(project.PaymentPara),
						Plans:       paymentPlansNew,
					},
					SitePlan: struct {
						Description string `json:"description,omitempty"`
						Image       string `json:"image,omitempty"`
					}{
						Description: safeStr(project.SitePlanPara),
						Image:       safeStr(project.SitePlanImg),
					},
					About: struct {
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
					}{
						Description:       safeStr(developer.About),
						LogoURL:           safeStr(developer.DeveloperLogo),
						EstablishmentYear: strconv.FormatInt(*developer.EstablishedYear, 10),
						TotalProjects:     safeStr(developer.ProjectDoneNo),
						ContactDetails: struct {
							Name           string `json:"name,omitempty"`
							ProjectAddress string `json:"project_address,omitempty"`
							Phone          string `json:"phone,omitempty"`
							BookingLink    string `json:"booking_link,omitempty"`
						}{
							Name:           safeStr(developer.DeveloperName),
							ProjectAddress: safeStr(developer.DeveloperAddress),
							Phone:          safeStr(developer.Phone),
							BookingLink:    safeStr(project.ProjectBrochure),
						},
					},
					Faqs: faqsNew,
				}).
				SetLocationInfo(schema.LocationInfo{
					ShortAddress:  safeStr(project.ShortAddress),
					Latitude:      "",
					Longitude:     "",
					GoogleMapLink: safeStr(project.ProjectLocationURL),
				}).
				SetProjectType(entproject.ProjectType(*projectConfigurationType.PropertyType)).
				SetIsFeatured(project.IsFeatured).
				SetIsPremium(project.IsPremium).
				SetIsPriority(project.IsPriority).
				SetIsDeleted(project.IsDeleted).
				SetDeveloperID(*developerID).
				SetLocationID(*localityID).
				Exec(ctx); err != nil {
				log.Error().Err(err).Msgf("Failed to insert project ID %d", project.ID)
				continue
			}
		}
		return nil
	}

	if err := processBatch(ctx, projects, batchSize, processProjectBatch); err != nil {
		return err
	}

	log.Info().Msg("Projects migrated successfully --------->>>> success")
	return nil
}

func MigrateDeveloper(ctx context.Context, txn *ent.Tx) error {
	log.Info().Msg("fetching developers")
	developers, err := FetchAllDevelopers(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("fetched developers --------->>>> success")

	processDeveloperBatch := func(ctx context.Context, batch []LDeveloper) error {
		for _, developer := range batch {
			id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(developer.ID, 10))))[:16]
			setDeveloperIDMapping(developer.ID, id)
			if err := txn.Developer.Create().
				SetID(id).
				SetName(safeStr(developer.DeveloperName)).
				SetLegalName(safeStr(developer.DeveloperLegalName)).
				SetIdentifier(safeStr(developer.DeveloperName)).
				SetEstablishedYear(safeInt(developer.EstablishedYear)).
				SetMediaContent(schema.DeveloperMediaContent{
					DeveloperAddress: safeStr(developer.DeveloperAddress),
					Phone:            safeStr(developer.Phone),
					DeveloperLogo:    safeStr(developer.DeveloperLogo),
					AltDeveloperLogo: safeStr(developer.AltDeveloperLogo),
					About:            safeStr(developer.About),
					Overview:         safeStr(developer.Overview),
					Disclaimer:       safeStr(developer.Disclaimer),
				}).
				SetIsVerified(developer.IsVerified != nil && *developer.IsVerified).
				Exec(ctx); err != nil {
				log.Error().Err(err).Msgf("Failed to insert developer ID %d", developer.ID)
				continue
			}
		}
		return nil
	}

	if err := processBatch(ctx, developers, batchSize, processDeveloperBatch); err != nil {
		return err
	}

	log.Info().Msg("Developers migrated successfully --------->>>> success")
	return nil
}

func MigrateLocality(ctx context.Context, txn *ent.Tx) error {
	log.Info().Msg("fetching localities")
	localities, err := FetchAllLocality(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("fetched localities --------->>>> success")

	processLocalityBatch := func(ctx context.Context, batch []LLocality) error {
		for _, locality := range batch {
			city, err := FetchCityByID(ctx, *locality.CityID)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch city for locality ID %d", locality.ID)
				continue
			}

			phoneInt, err := parsePhoneJSONToString(city.Phone)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to convert phone for locality ID %d", locality.ID)
				continue
			}

			id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(locality.ID, 10))))[:16]
			setLocalityIDMapping(locality.ID, id)

			if err := txn.Location.Create().
				SetID(id).
				SetLocalityName(safeStr(locality.Name)).
				SetCity(safeStr(city.Name)).
				SetState(safeStr(city.StateName)).
				SetPhoneNumber(*phoneInt).
				SetCountry("India").
				SetPincode("112222").
				SetIsActive(true).
				Exec(ctx); err != nil {
				log.Error().Err(err).Msgf("Failed to insert locality ID %d", locality.ID)
				continue
			}
		}
		return nil
	}

	if err := processBatch(ctx, localities, batchSize, processLocalityBatch); err != nil {
		return err
	}

	log.Info().Msg("Localities migrated successfully --------->>>> success")
	return nil
}

func MigrateProperty(ctx context.Context, txn *ent.Tx) error {
	log.Info().Msg("fetching properties --------->>>> success")
	properties, err := fetchAllProperty(ctx)
	if err != nil {
		return err
	}
	log.Info().Msgf("Fetched all properties --------->>>> success")

	processPropertyBatch := func(ctx context.Context, batch []LProperty) error {
		for _, property := range batch {
			id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(property.ID, 10))))[:16]

			if property.PropertyName == nil {
				log.Error().Msgf("Property name is nil for property ID %d", property.ID)
				continue
			}

			var propertyConfiguration *LPropertyConfiguration
			if property.ConfigurationID != nil {
				propertyConfiguration, err = FetchPropertyConfigurationByID(ctx, *property.ConfigurationID)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to fetch project configurations for property ID %d", property.ID)
					continue
				}
			}

			if propertyConfiguration == nil {
				propertyConfiguration = &LPropertyConfiguration{
					ProjectConfigurationName: strPtr("Not Available"),
				}
			}

			uspList := []string{}
			if property.USP != nil {
				uspText := strings.Trim(*property.USP, "[]")
				uspItems := strings.Split(uspText, "\",")
				for _, item := range uspItems {
					item = strings.Trim(item, "\" ")
					if item != "" {
						uspList = append(uspList, item)
					}
				}
			}

			parsedImages, err := parsePropertyImagesFromPropertyImages(property.Images)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to parse property images for property ID %d", property.ID)
				continue
			}

			var propertyType *LPropertyConfigurationType
			if property.ConfigurationID != nil {
				propertyType, err = FetchPropertyConfigurationTypeByID(ctx, *property.ConfigurationID)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to fetch property type for property ID %d", property.ID)
					continue
				}
			}

			if propertyType == nil {
				propertyType = &LPropertyConfigurationType{
					PropertyType: strPtr("Not Available"),
				}
			}

			ImageUrlWithType := make(map[string]string)
			if property.FloorPara != nil && *property.FloorPara != "" {
				ImageUrlWithType["2D"] = safeStr(property.FloorImage2D)
				ImageUrlWithType["3D"] = safeStr(property.FloorImage3D)
			}

			webCard := schema.WebCards{
				PropertyDetails: schema.PropertyDetails{
					BuiltUpArea: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.BuiltupArea),
					},
					Sizes: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Size),
					},
					FloorNumber: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Floors),
					},
					Configuration: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(propertyConfiguration.ProjectConfigurationName),
					},
					PossessionStatus: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.PossessionStatus),
					},
					Balconies: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Balcony),
					},
					CoveredParking: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.CoveredParking),
					},
					Bedrooms: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Bedrooms),
					},
					PropertyType: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(propertyType.ConfigurationTypeName),
					},
					AgeOfProperty: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.AgeOfProperty),
					},
					FurnishingType: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.FurnishingType),
					},
					ReraNumber: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Rera),
					},
					Facing: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Facing),
					},
					Bathrooms: struct {
						Value string `json:"value,omitempty"`
					}{
						Value: safeStr(property.Bathrooms),
					},
				},
				PropertyFloorPlan: schema.PropertyFloorPlan{
					Title: safeStr(property.FloorPara),
					Plans: []map[string]string{
						{
							"2D": ImageUrlWithType["2D"],
							"3D": ImageUrlWithType["3D"],
						},
					},
				},
				KnowAbout: struct {
					Description string `json:"description,omitempty"`
				}{
					Description: safeStr(property.About),
				},
				LocationMap: struct {
					Description   string `json:"description,omitempty"`
					GoogleMapLink string `json:"google_map_link,omitempty"`
				}{
					Description:   safeStr(property.LocationPara),
					GoogleMapLink: safeStr(property.LocaionMap),
				},
			}

			var projectID *string
			var ok bool
			if property.ProjectID != nil {
				projectID, ok = getProjectIDMapping(*property.ProjectID)
				if !ok {
					log.Error().Msgf("Project ID mapping not found for property ID %d", property.ID)
					return fmt.Errorf("project ID mapping not found for property ID %d", property.ID)
				}
			}

			var developerID *string
			if property.DeveloperID != nil {
				developerID, ok = getDeveloperIDMapping(*property.DeveloperID)
				if !ok {
					log.Error().Msgf("Developer ID mapping not found for property ID %d", property.ID)
					return fmt.Errorf("developer ID mapping not found for property ID %d", property.ID)
				}
			}

			var localityID *string
			if property.LocalityID != nil {
				localityID, ok = getLocalityIDMapping(*property.LocalityID)
				if !ok {
					log.Error().Msgf("Locality ID mapping not found for property ID %d", property.ID)
					return fmt.Errorf("locality ID mapping not found for property ID %d", property.ID)
				}
			}

			entProperty := txn.Property.Create().
				SetID(id).
				SetName(*property.PropertyName).
				SetPropertyImages(parsedImages).
				SetPropertyType(safeStr(propertyType.PropertyType)).
				SetWebCards(webCard).
				SetPricingInfo(schema.PropertyPricingInfo{
					Price: strconv.FormatFloat(property.Price, 'f', -1, 64),
				}).
				SetPropertyReraInfo(schema.PropertyReraInfo{
					ReraNumber: safeStr(property.Rera),
				}).
				SetMetaInfo(schema.PropertyMetaInfo{
					Title:       safeStr(property.MetaTitle),
					Description: safeStr(property.MetaDescription),
					Keywords:    safeStr(property.MetaKeywords),
					Canonical:   safeStr(property.PropertyURL),
				}).
				SetIsFeatured(property.IsFeatured).
				SetIsDeleted(property.IsDeleted).
				SetDeveloperID(*developerID).
				SetLocationID(*localityID)

			if projectID != nil {
				entProperty = entProperty.SetProjectID(*projectID)
			}

			if err := entProperty.Exec(ctx); err != nil {
				log.Error().Err(err).Msgf("Failed to insert property ID %d", property.ID)
				return err
			}

		}
		return nil
	}

	if err := processBatch(ctx, properties, batchSize, processPropertyBatch); err != nil {
		return err
	}

	log.Info().Msg("Properties migrated successfully --------->>>> success")
	return nil
}

func safeStrToInt(s *string) int {
	if s != nil {
		if n, err := strconv.Atoi(*s); err == nil {
			return n
		}
	}
	return 0
}

func strPtr(s string) *string {
	return &s
}

func MigrateStaticSiteData(ctx context.Context, txn *ent.Tx) error {
	log.Info().Msg("fetching property configurations --------->>>> success")
	residentialConfigs := []string{
		"1BHK", "2BHK", "2.5BHK", "3BHK", "3.5BHK", "4BHK", "4.5BHK",
		"5BHK", "5.5BHK", "6BHK", "6.5BHK", "7BHK", "7.5BHK", "8BHK", "8.5BHK",
		"VILLAS", "PENTHOUSE", "DUPLEX", "SIMPLEX", "STUDIO APARTMENT", "PLOTS",
		"INDEPENDENT FLOOR",
	}

	commercialConfigs := []string{
		"SHOPS", "SUITS", "OFFICE", "RETAIL SHOP", "RENTAL SPACES",
		"LEASE SPACES", "FOODCOURT", "ANCHOR SPACES", "CO-WORKING SPACES",
		"VIRTUAL SPACES",
	}

	// Create residential configurations
	staticSiteWebCardsID = uuid.New().String()
	propertyTypes := schema.PropertyTypes{
		Commercial:  commercialConfigs,
		Residential: residentialConfigs,
	}
	if err := txn.StaticSiteData.Create().
		SetID(staticSiteWebCardsID).
		SetPropertyTypes(propertyTypes).
		Exec(ctx); err != nil {
		log.Error().Err(err).Msgf("Failed to insert residential configuration: %s", staticSiteWebCardsID)
		return err
	}

	log.Info().Msg("Property configurations migrated successfully --------->>>> success")
	return nil
}

func MigrateBlogs(ctx context.Context, txn *ent.Tx) error {
	if txn == nil {
		log.Fatal().Msg("Transaction (txn) is nil")
	}

	log.Info().Msg("Fetching blogs")
	blogs, err := FetchAllBlogs(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("Fetched blogs --------->>>> success")

	processBlogBatch := func(ctx context.Context, batch []LBlog) error {
		for _, blog := range batch {

			id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(blog.ID, 10))))[:16]
			// Parse images from JSON string
			var images []string
			if blog.Images != nil {
				if err := json.Unmarshal([]byte(*blog.Images), &images); err != nil {
					log.Error().Err(err).Msgf("Failed to parse images for blog ID %d", blog.ID)
				}
			}
			// Clean up blog schema format
			//"[\"1\",\"2\",\"3\"]"

			blogContent := schema.BlogContent{
				Title:       safeStr(blog.Headings),
				Description: safeStr(blog.Description),
			}
			if len(images) > 0 {
				blogContent.Image = images[0]
				blogContent.ImageAlt = safeStr(blog.Alt)
			}

			// Create SEO meta info
			seoMetaInfo := schema.SEOMetaInfo{
				BlogSchema: blog.BlogSchema,
				Canonical:  safeStr(blog.Canonical),
				Title:      safeStr(blog.SubHeadings),
				Keywords:   safeStr(blog.MetaKeywords),
			}

			// Create blog entry
			if err := txn.Blogs.Create().
				SetID(id).
				SetBlogURL(safeStr(blog.BlogURL)).
				SetBlogContent(blogContent).
				SetSeoMetaInfo(seoMetaInfo).
				SetIsPriority(blog.IsPriority).
				SetCreatedAt(*blog.CreatedDate).
				SetUpdatedAt(int64(safeInt(blog.UpdatedDate))).
				SetIsDeleted(blog.IsDeleted).
				Exec(ctx); err != nil {
				log.Error().Err(err).Msgf("Failed to insert blog ID %d", blog.ID)
				continue
			}
		}
		return nil
	}

	if err := processBatch(ctx, blogs, batchSize, processBlogBatch); err != nil {
		return err
	}

	log.Info().Msg("Blogs migrated successfully --------->>>> success")
	return nil
}
