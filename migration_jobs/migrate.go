package migration_jobs

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
	"github.com/rs/zerolog/log"
)

// Use ent schema modal to migrate data from legacy database to new database

var (
	legacyToNewProjectIDMAP       = make(map[int64]string)
	legacyToNewDeveloperIDMAP     = make(map[int64]string)
	legacyToNewLocalityIDMAP      = make(map[int64]string)
	legacyToNewConfigurationIDMAP = make(map[string]string)
	legacyToNewConfigTypeIDMAP    = make(map[string]string)
	legacyToNewPropertyIDMAP      = make(map[int64]string)
)

func migrateProject(ctx context.Context, db *sql.DB) error {
	projects, err := FetchhAllProject(ctx, db)
	if err != nil {
		return err
	}

	for _, project := range projects {
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(project.ID, 10))))[:16]
		// append images in sequesnce

		projectRera, err := FetchReraByProjectID(ctx, db, project.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project RERA for project ID %d", project.ID)
			continue
		}

		reras := []schema.ReraListItem{}
		for _, rera := range projectRera {
			reras = append(reras, schema.ReraListItem{
				Phase:      *rera.Phase,
				ReraQR:     *rera.QRImages,
				ReraNumber: *rera.ReraNumber,
				Status:     *rera.Status,
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

		uspListNew := []schema.USPListItem{}
		for _, usp := range uspList {
			uspListNew = append(uspListNew, schema.USPListItem{
				Icon:        "",
				Description: usp,
			})
		}

		floorPlans, err := FetchFloorPlanByProjectID(ctx, db, project.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch floor plans for project ID %d", project.ID)
			continue
		}

		floorPlanItems := []schema.FloorPlanItem{}
		configurationProducts := []schema.ProductConfiguration{}
		for _, floorPlan := range *floorPlans {
			floorPlanItems = append(floorPlanItems, schema.FloorPlanItem{
				Title:        *floorPlan.Title,
				FlatType:     *floorPlan.Title,
				IsSoldOut:    floorPlan.IsSoldOut,
				Price:        strconv.FormatFloat(floorPlan.Price, 'f', -1, 64),
				BuildingArea: strconv.FormatInt(*floorPlan.Size, 10),
				Image:        *floorPlan.ImgURL,
			})
			configurationProducts = append(configurationProducts, schema.ProductConfiguration{
				ConfigurationName: *floorPlan.Title,
				Size:              strconv.FormatInt(*floorPlan.Size, 10),
				Price:             strconv.FormatFloat(floorPlan.Price, 'f', -1, 64),
			})
		}

		amenities, err := FetchProjectAmenitiesByProjectID(ctx, db, project.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch amenities for project ID %d", project.ID)
			continue
		}

		amenitiesMap := map[string][]schema.AmenityCategory{}
		for _, amenity := range amenities {
			amenitiesMap[*amenity.AmenitiesCategory] = append(amenitiesMap[*amenity.AmenitiesCategory], schema.AmenityCategory{
				Icon:  *amenity.AmenitiesURL,
				Value: *amenity.AmenitiesName,
			})
		}

		paymentPlans, err := FetchPaymentPlansByProjectID(ctx, db, project.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch payment plans for project ID %d", project.ID)
			continue
		}

		paymentPlansNew := []schema.Plan{}
		for _, paymentPlan := range paymentPlans {
			paymentPlansNew = append(paymentPlansNew, schema.Plan{
				Name:    *paymentPlan.PaymentPlanName,
				Details: *paymentPlan.PaymentPlanValue,
			})
		}

		developer, err := FetchDeveloperByID(ctx, db, *project.DeveloperID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch developer for project ID %d", project.ID)
			return err
		}

		faqs, err := FetchFaqsByProjectID(ctx, db, project.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch faqs for project ID %d", project.ID)
			return err
		}

		faqsNew := []schema.FAQ{}
		for _, faq := range faqs {
			faqsNew = append(faqsNew, schema.FAQ{
				Question: *faq.Question,
				Answer:   *faq.Answer,
			})
		}

		legacyToNewProjectIDMAP[project.ID] = id
		if err := newDB.Project.Create().
			SetID(id).
			SetName(*project.ProjectName).
			SetDescription(*project.ProjectDescription).
			SetStatus(enums.ProjectStatus(*project.Status)).
			SetTotalTowers(int(*project.TotalTowers)).
			SetTotalFloor(int(*project.TotalFloor)).
			SetTotalTowers(int(*project.TotalTowers)).
			SetTimelineInfo(schema.TimelineInfo{
				ProjectLaunchDate:     *project.ProjectLaunchDate,
				ProjectPossessionDate: *project.ProjectPossessionDate,
			}).
			SetMetaInfo(schema.SEOMeta{
				Title:         *project.MetaTitle,
				Description:   *project.MetaDescription,
				Keywords:      *project.MetaKeywords,
				Canonical:     *project.ProjectURL,
				ProjectSchema: *project.ProjectSchema,
			}).
			SetWebCards(schema.ProjectWebCards{
				Images: []string{*project.CoverPhoto},
				ReraInfo: schema.ReraInfo{
					WebsiteLink: *project.ReraLink,
					ReraList:    reras,
				},
				Details: schema.ProjectDetails{
					Area: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectArea,
					},
					Sizes: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectArea,
					},
					Units: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectArea,
					},
					Configuration: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectArea,
					},
					LaunchDate: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectLaunchDate,
					},
					PossessionDate: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectPossessionDate,
					},
					TotalTowers: struct {
						Value string `json:"value"`
					}{
						Value: strconv.FormatInt(*project.TotalTowers, 10),
					},
					TotalFloors: struct {
						Value string `json:"value"`
					}{
						Value: strconv.FormatInt(*project.TotalFloor, 10),
					},
					ProjectStatus: struct {
						Value string `json:"value"`
					}{
						Value: *project.Status,
					},
					Type: struct {
						Value string `json:"value"`
					}{
						Value: *project.ProjectConfigurations,
					},
				},
				WhyToChoose: schema.WhyToChoose{
					ImageUrls: []string{*project.SitePlanImg},
					USP_List:  uspListNew,
				},
				KnowAbout: schema.KnowAbout{
					Description:  *project.ProjectAbout,
					DownloadLink: *project.ProjectBrochure,
				},
				FloorPlan: schema.FloorPlan{
					Description: *project.PriceListPara,
					Products:    floorPlanItems,
				},
				PriceList: schema.PriceList{
					Description:          *project.PriceListPara,
					BHKOptionsWithPrices: configurationProducts,
				},
				Amenities: schema.Amenities{
					Description:             *project.AmenitiesPara,
					CategoriesWithAmenities: amenitiesMap,
				},
				VideoPresentation: schema.VideoPresentation{
					Description: *project.VideoPara,
					URL:         project.ProjectVideos,
				},
				PaymentPlans: schema.PaymentPlans{
					Description: *project.PaymentPara,
					Plans:       paymentPlansNew,
				},
				SitePlan: struct {
					Description string `json:"description"`
					Image       string `json:"image"`
				}{
					Description: *project.SitePlanPara,
					Image:       *project.SitePlanImg,
				},
				About: struct {
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
				}{
					Description:       *project.ProjectAbout,
					LogoURL:           *project.ProjectLogo,
					EstablishmentYear: strconv.FormatInt(*developer.EstablishedYear, 10),
					TotalProjects:     strconv.FormatInt(*developer.ProjectDoneNo, 10),
					ContactDetails: struct {
						Name           string `json:"name"`
						ProjectAddress string `json:"project_address"`
						Phone          string `json:"phone"`
						BookingLink    string `json:"booking_link"`
					}{
						Name:           *developer.DeveloperName,
						ProjectAddress: *developer.DeveloperAddress,
						Phone:          *developer.Phone,
						BookingLink:    *developer.DeveloperURL,
					},
				},
				Faqs: faqsNew,
			}).
			SetIsFeatured(project.IsFeatured).
			SetIsPremium(project.IsPremium).
			SetIsPriority(project.IsPriority).
			SetIsDeleted(project.IsDeleted).
			SetDeveloperID(legacyToNewDeveloperIDMAP[*project.DeveloperID]).
			SetLocationID(legacyToNewLocalityIDMAP[*project.LocalityID]).
			Exec(ctx); err != nil {
			return err
		}

	}

	// fetch all projects from legacy database

	// create project in iteration
	// map legacy project id to new project id
	// setDeveloperID(legacyToNewDeveloperIDMAP[lproject.DeveloperID])

	return nil
}

func MigrateDeveloper(ctx context.Context, db *sql.DB, newDB *ent.Client) error {
	log.Info().Msg("Migrating developers")
	ldeveloper, err := FetchAllDevelopers(ctx, db)
	if err != nil {
		return err
	}

	for _, developer := range ldeveloper {
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(developer.ID, 10))))[:16]
		legacyToNewDeveloperIDMAP[developer.ID] = id
		if err := newDB.Developer.Create().
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
			return err
		}

	}
	return nil
}

func MigrateLocality(ctx context.Context, db *sql.DB, newDB *ent.Client) error {
	//new location id will be generated
	llocality, err := FetchAllLocality(ctx, db)
	if err != nil {
		return err
	}

	for _, locality := range llocality {

		city, err := FetchCityByID(ctx, db, *locality.CityID)

		log.Info().Msgf("Fetched city %+v", city)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch city for locality ID %d", locality.ID)
			continue
		}

		if newDB == nil {
			return fmt.Errorf("newDB is nil — database connection not initialized")
		}

		phoneInt, err := parsePhoneJSONToString(city.Phone)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to convert phone for locality ID %d", locality.ID)
		}

		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(locality.ID, 10))))[:16]
		legacyToNewLocalityIDMAP[locality.ID] = id
		if err := newDB.Location.Create().
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

func migrateProperty(ctx context.Context, db *sql.DB) error {
	properties, err := fetchAllProperty(ctx, db)
	if err != nil {
		return err
	}
	for _, property := range properties {
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.FormatInt(property.ID, 10))))[:16]
		legacyToNewPropertyIDMAP[property.ID] = id
		project, err := FetchProjectByID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project for property ID %d", property.ID)
			continue
		}
		projectConfigurations, err := FetchProjectConfigurationsByID(ctx, db, *property.ConfigurationID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project configurations for property ID %d", property.ID)
			continue
		}
		locality, err := FetchLocalityByID(ctx, db, *property.LocalityID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch locality for property ID %d", property.ID)
			continue
		}
		developer, err := FetchDeveloperByID(ctx, db, *property.DeveloperID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch developer for property ID %d", property.ID)
		}

		projectImages, err := FetchProjectImagesByProjectID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch project images for property ID %d", property.ID)
			continue
		}

		propertyImages, err := parsePropertyImagesFromProjectImages(projectImages)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse property images for property ID %d", property.ID)
			continue
		}
		floorPlans, err := FetchFloorPlansByProjectID(ctx, db, *property.ProjectID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch floor plans for property ID %d", property.ID)
			continue
		}
		webCards, err := parseWebCardsFromProject(project)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse web cards for property ID %d", property.ID)
			continue
		}

		// property Type and configurationtype(Ground Floor, Apartment, etc.) name from projectConfigurations table

		if err := newDB.Property.Create().
			SetID(id).
			SetName(*property.PropertyName).
			SetPropertyImages(*propertyImages).
			SetProjectID(legacyToNewProjectIDMAP[*property.ProjectID]).
			SetWebCards(*webCards).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
