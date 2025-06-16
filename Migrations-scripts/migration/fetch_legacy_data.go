package migration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Helper function to convert MySQL bit(1) to bool
func bitToBool(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return b[0] == 1
}

func FetchLegacyProjectData(ctx context.Context, db *sql.DB) ([]ProjectLegacyData, error) {
	// Query to get all projects with their relationships
	query := `
		SELECT 
			p.id, p.project_name, p.project_url, p.project_description,
			p.project_about, p.project_address, p.project_area,
			p.project_brochure, p.project_configurations, p.project_launch_date,
			p.project_location_url, p.project_logo, p.project_possession_date,
			p.project_rera, p.project_schema, p.project_units,
			p.project_video_count, p.project_videos, p.rera_link,
			p.short_address, p.status, p.total_floor, p.total_towers,
			p.usp, p.user_id, p.available_unit,
			p.alt_project_logo, p.alt_site_plan_img, p.cover_photo,
			p.location_map, p.siteplan_img, p.meta_desciption,
			p.meta_title, p.meta_keywords, p.amenities_para,
			p.floor_para, p.location_para, p.overview_para,
			p.payment_para, p.price_list_para, p.siteplan_para,
			p.video_para, p.why_para, p.is_deleted,
			p.is_featured, p.is_premium, p.is_priority,
			p.created_date, p.updated_date,

			d.id as dev_id, d.developer_name, d.developer_legal_name, d.developer_url, 
			d.developer_address, d.developer_logo, d.alt_developer_logo, d.about,
			d.overview, d.disclaimer, d.established_year, d.project_done_no,
			d.phone, d.is_active, d.is_verified, d.city_name, d.created_date as dev_created,
			d.updated_date as dev_updated,
			
			l.id as loc_id, l.locality_name, l.locality_url, l.city_id,
			l.created_date as loc_created, l.updated_date as loc_updated,
			
			pc.id as config_id, pc.project_configuration_name,
			pc.configuration_type_id, pc.created_date as config_created,
			pc.updated_date as config_updated,

			pct.id as config_type_id, pct.configuration_type_name,
			COALESCE(pct.property_type, '') as property_type, pct.created_date as config_type_created,
			pct.updated_date as config_type_updated
		FROM 
			project p
			LEFT JOIN developer d ON p.developer_id = d.id
			LEFT JOIN locality l ON p.locality_id = l.id
			LEFT JOIN city c ON l.city_id = c.id
			LEFT JOIN project_configuration pc ON p.property_configuration_type_id = pc.id
			LEFT JOIN project_configuration_type pct ON pc.configuration_type_id = pct.id
		WHERE 
			p.is_deleted = 0
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying projects: %v", err)
	}
	defer rows.Close()

	var projects []ProjectLegacyData

	for rows.Next() {
		var p ProjectLegacyData
		p.Developer = &DeveloperData{}
		p.Locality = &LocalityData{}
		p.Configuration = &ConfigurationData{}
		p.ConfigType = &ProjectConfigTypeData{}

		// Temporary variables for scanning nullable fields
		var (
			altProjectLogo, altSitePlanImg, amenitiesPara                   sql.NullString
			coverPhoto, floorPara, locationMap, locationPara                sql.NullString
			metaDescription, metaTitle, overviewPara, paymentPara           sql.NullString
			priceListPara, projectAbout, projectAddress, projectArea        sql.NullString
			projectBrochure, projectConfigs, projectDesc, projectLaunchDate sql.NullString
			projectLocURL, projectLogo, projectName, projectPossessionDate  sql.NullString
			projectRERA, projectSchema, projectUnits, projectURL            sql.NullString
			reraLink, shortAddress, sitePlanImg, sitePlanPara               sql.NullString
			status, totalFloor, totalTowers, usp, availableUnit             sql.NullString
			videoPara, whyPara                                              sql.NullString
			metaKeywordsJSON                                                sql.NullString // Changed to sql.NullString since it's a TEXT field
			projectVideoCount                                               sql.NullInt64
			projectVideos                                                   []byte
			isDeletedBit, isFeaturedBit, isPremiumBit, isPriorityBit        []byte
			userID                                                          sql.NullInt64

			// Developer nullable fields
			devName, devLegalName, devURL, devAddr, devLogo, devAltLogo sql.NullString
			devAbout, devOverview, devDisclaimer, devPhone              sql.NullString
			devEstYear, devProjectDone, devCityName                     sql.NullInt64
			devIsActiveBit, devIsVerifiedBit                            []byte

			// Locality nullable fields
			locName, locURL sql.NullString
			locCityID       sql.NullInt64

			// Configuration nullable fields
			configName   sql.NullString
			configTypeID sql.NullInt64

			// Project configuration type nullable fields
			configTypeCreatedDate, configTypeUpdatedDate sql.NullInt64
			propertyType                                 sql.NullString
		)

		err := rows.Scan(
			// Project fields (52 columns)
			&p.ID, &projectName, &projectURL, &projectDesc,
			&projectAbout, &projectAddress, &projectArea,
			&projectBrochure, &projectConfigs, &projectLaunchDate,
			&projectLocURL, &projectLogo, &projectPossessionDate,
			&projectRERA, &projectSchema, &projectUnits,
			&projectVideoCount, &projectVideos, &reraLink,
			&shortAddress, &status, &totalFloor, &totalTowers,
			&usp, &userID, &availableUnit,
			&altProjectLogo, &altSitePlanImg, &coverPhoto,
			&locationMap, &sitePlanImg, &metaDescription,
			&metaTitle, &metaKeywordsJSON, &amenitiesPara,
			&floorPara, &locationPara, &overviewPara,
			&paymentPara, &priceListPara, &sitePlanPara,
			&videoPara, &whyPara, &isDeletedBit,
			&isFeaturedBit, &isPremiumBit, &isPriorityBit,
			&p.CreatedDate, &p.UpdatedDate,

			// Developer fields (18 columns)
			&p.Developer.ID, &devName, &devLegalName, &devURL, &devAddr,
			&devLogo, &devAltLogo, &devAbout, &devOverview, &devDisclaimer,
			&devEstYear, &devProjectDone, &devPhone, &devIsActiveBit, &devIsVerifiedBit,
			&devCityName, &p.Developer.CreatedDate, &p.Developer.UpdatedDate,

			// Locality fields (6 columns)
			&p.Locality.ID, &locName, &locURL, &locCityID,
			&p.Locality.CreatedDate, &p.Locality.UpdatedDate,

			// Configuration fields (5 columns)
			&p.Configuration.ID, &configName, &configTypeID,
			&p.Configuration.CreatedDate, &p.Configuration.UpdatedDate,

			// Project configuration type fields (5 columns)
			&p.ConfigType.ID, &p.ConfigType.ConfigurationTypeName,
			&propertyType, &configTypeCreatedDate,
			&configTypeUpdatedDate,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}

		// Handle nullable fields
		if projectName.Valid {
			p.ProjectName = projectName.String
		}
		if projectURL.Valid {
			p.ProjectURL = projectURL.String
		}
		if projectDesc.Valid {
			p.ProjectDescription = projectDesc.String
		}
		if projectAbout.Valid {
			p.ProjectAbout = projectAbout.String
		}
		if projectAddress.Valid {
			p.ProjectAddress = projectAddress.String
		}
		if projectArea.Valid {
			p.ProjectArea = projectArea.String
		}
		if projectBrochure.Valid {
			p.ProjectBrochure = projectBrochure.String
		}
		if projectConfigs.Valid {
			p.ProjectConfigurations = projectConfigs.String
		}
		if projectLaunchDate.Valid {
			p.ProjectLaunchDate = projectLaunchDate.String
		}
		if projectLocURL.Valid {
			p.ProjectLocationURL = projectLocURL.String
		}
		if projectLogo.Valid {
			p.ProjectLogo = projectLogo.String
		}
		if projectPossessionDate.Valid {
			p.ProjectPossessionDate = projectPossessionDate.String
		}
		if projectRERA.Valid {
			p.ProjectRERA = projectRERA.String
		}
		if projectSchema.Valid {
			p.ProjectSchema = projectSchema.String
		}
		if projectUnits.Valid {
			p.ProjectUnits = projectUnits.String
		}
		if reraLink.Valid {
			p.RERALink = reraLink.String
		}
		if shortAddress.Valid {
			p.ShortAddress = shortAddress.String
		}
		if status.Valid {
			p.Status = status.String
		}
		if totalFloor.Valid {
			p.TotalFloor = totalFloor.String
		}
		if totalTowers.Valid {
			p.TotalTowers = totalTowers.String
		}
		if usp.Valid {
			p.USP = usp.String
		}
		if availableUnit.Valid {
			p.AvailableUnit = availableUnit.String
		}
		if altProjectLogo.Valid {
			p.AltProjectLogo = altProjectLogo.String
		}
		if altSitePlanImg.Valid {
			p.AltSitePlanImg = altSitePlanImg.String
		}
		if coverPhoto.Valid {
			p.CoverPhoto = coverPhoto.String
		}
		if locationMap.Valid {
			p.LocationMap = locationMap.String
		}
		if sitePlanImg.Valid {
			p.SitePlanImg = sitePlanImg.String
		}
		if metaDescription.Valid {
			p.MetaDescription = metaDescription.String
		}
		if metaTitle.Valid {
			p.MetaTitle = metaTitle.String
		}
		if amenitiesPara.Valid {
			p.AmenitiesPara = amenitiesPara.String
		}
		if floorPara.Valid {
			p.FloorPara = floorPara.String
		}
		if locationPara.Valid {
			p.LocationPara = locationPara.String
		}
		if overviewPara.Valid {
			p.OverviewPara = overviewPara.String
		}
		if paymentPara.Valid {
			p.PaymentPara = paymentPara.String
		}
		if priceListPara.Valid {
			p.PriceListPara = priceListPara.String
		}
		if sitePlanPara.Valid {
			p.SitePlanPara = sitePlanPara.String
		}
		if videoPara.Valid {
			p.VideoPara = videoPara.String
		}
		if whyPara.Valid {
			p.WhyPara = whyPara.String
		}
		if projectVideoCount.Valid {
			p.ProjectVideoCount = projectVideoCount.Int64
		}
		if userID.Valid {
			p.UserID = userID.Int64
		}

		// Handle bit fields
		p.IsDeleted = bitToBool(isDeletedBit)
		p.IsFeatured = bitToBool(isFeaturedBit)
		p.IsPremium = bitToBool(isPremiumBit)
		p.IsPriority = bitToBool(isPriorityBit)

		// Handle Developer fields
		if devName.Valid {
			p.Developer.DeveloperName = devName.String
		}
		if devLegalName.Valid {
			p.Developer.DeveloperLegalName = devLegalName.String
		}
		if devURL.Valid {
			p.Developer.DeveloperURL = devURL.String
		}
		if devAddr.Valid {
			p.Developer.DeveloperAddress = devAddr.String
		}
		if devLogo.Valid {
			p.Developer.DeveloperLogo = devLogo.String
		}
		if devAltLogo.Valid {
			p.Developer.AltDeveloperLogo = devAltLogo.String
		}
		if devAbout.Valid {
			p.Developer.About = devAbout.String
		}
		if devOverview.Valid {
			p.Developer.Overview = devOverview.String
		}
		if devDisclaimer.Valid {
			p.Developer.Disclaimer = devDisclaimer.String
		}
		if devPhone.Valid {
			p.Developer.Phone = devPhone.String
		}
		if devEstYear.Valid {
			p.Developer.EstablishedYear = devEstYear.Int64
		}
		if devProjectDone.Valid {
			p.Developer.ProjectDoneNo = devProjectDone.Int64
		}
		if devCityName.Valid {
			p.Developer.CityName = devCityName.Int64
		}
		p.Developer.IsActive = bitToBool(devIsActiveBit)
		p.Developer.IsVerified = bitToBool(devIsVerifiedBit)

		// Handle Locality fields
		if locName.Valid {
			p.Locality.LocalityName = locName.String
		}
		if locURL.Valid {
			p.Locality.LocalityURL = locURL.String
		}
		if locCityID.Valid {
			p.Locality.CityID = locCityID.Int64
		}

		// Handle Configuration fields
		if configName.Valid {
			p.Configuration.ProjectConfigurationName = configName.String
		}
		if configTypeID.Valid {
			p.Configuration.ConfigurationTypeID = configTypeID.Int64
		}

		// Parse meta_keywords JSON
		if metaKeywordsJSON.Valid && len(metaKeywordsJSON.String) > 0 {
			if err := json.Unmarshal([]byte(metaKeywordsJSON.String), &p.MetaKeywords); err != nil {
				log.Printf("Warning: error parsing meta_keywords JSON for project %d: %v", p.ID, err)
			}
		}

		// Handle Project configuration type nullable fields
		if configTypeCreatedDate.Valid {
			p.ConfigType.CreatedDate = configTypeCreatedDate.Int64
		}
		if configTypeUpdatedDate.Valid {
			p.ConfigType.UpdatedDate = configTypeUpdatedDate.Int64
		}
		if propertyType.Valid {
			p.ConfigType.PropertyType = propertyType.String
		}

		// Parse project videos
		if projectVideos != nil && len(projectVideos) > 0 {
			// Log the raw BLOB data for debugging
			// log.Printf("Project ID %d: Raw BLOB data length: %d bytes", p.ID, len(projectVideos))
			// log.Printf("Project ID %d: Raw BLOB data (hex): %x", p.ID, projectVideos)

			// Try different parsing approaches
			var videoURLs []string

			// 1. Try to parse as JSON array
			if err := json.Unmarshal(projectVideos, &videoURLs); err == nil && len(videoURLs) > 0 {
				p.ProjectVideos = make([]ProjectVideo, len(videoURLs))
				for i, url := range videoURLs {
					p.ProjectVideos[i] = ProjectVideo{
						VideoURL: url,
					}
				}
				// log.Printf("Project ID %d: Found %d videos from JSON array", p.ID, len(videoURLs))
			} else {
				// 2. Try to parse as single JSON string
				var singleURL string
				if err := json.Unmarshal(projectVideos, &singleURL); err == nil && singleURL != "" {
					p.ProjectVideos = []ProjectVideo{{
						VideoURL: singleURL,
					}}
					// log.Printf("Project ID %d: Found single video URL from JSON: %s", p.ID, singleURL)
				} else {
					// 3. Try to parse as Java serialized data
					if len(projectVideos) > 2 && projectVideos[0] == 0xac && projectVideos[1] == 0xed {
						// This is Java serialized data
						for i := 0; i < len(projectVideos)-2; i++ {
							if projectVideos[i] == 0x74 { // TC_STRING marker
								length := int(projectVideos[i+1])<<8 | int(projectVideos[i+2])
								if i+3+length <= len(projectVideos) {
									url := string(projectVideos[i+3 : i+3+length])
									p.ProjectVideos = []ProjectVideo{{
										VideoURL: url,
									}}
									// log.Printf("Project ID %d: Found video URL from Java serialization: %s", p.ID, url)
									break
								}
							}
						}
					} else {
						// 4. Try to find URLs in the raw data
						urlPattern := regexp.MustCompile(`https?://[^\s<>"]+|www\.[^\s<>"]+`)
						matches := urlPattern.FindAll(projectVideos, -1)
						if len(matches) > 0 {
							p.ProjectVideos = make([]ProjectVideo, len(matches))
							for i, match := range matches {
								p.ProjectVideos[i] = ProjectVideo{
									VideoURL: string(match),
								}
							}
							log.Printf("Project ID %d: Found %d video URLs in raw data", p.ID, len(matches))
						} else {
							// 5. Try to convert to string and clean it
							strData := string(projectVideos)
							// Remove any non-printable characters
							strData = regexp.MustCompile(`[^\x20-\x7E]`).ReplaceAllString(strData, "")
							matches = urlPattern.FindAll([]byte(strData), -1)
							if len(matches) > 0 {
								p.ProjectVideos = make([]ProjectVideo, len(matches))
								for i, match := range matches {
									p.ProjectVideos[i] = ProjectVideo{
										VideoURL: string(match),
									}
								}
								log.Printf("Project ID %d: Found %d video URLs in cleaned data", p.ID, len(matches))
							} else {
								log.Printf("Project ID %d: Could not parse video data in any format", p.ID)
							}
						}
					}
				}
			}
		} else {
			log.Printf("Project ID %d: No video data found", p.ID)
		}

		// Fetch properties for this project
		propertiesQuery := `
			SELECT 
				id, about, age_of_property, amenities_para, balcony, bathrooms, bedrooms,
				builtup_area, cover_photo, covered_parking, created_date, facing,
				floor_image2d, floor_image3d, floor_para, floors, furnishing_type,
				images, is_deleted, is_featured, latlong, listing_type, location_map,
				location_advantage, location_para, logo_image, meta_description,
				meta_keywords, meta_title, open_parking, overview_para, possession_date,
				possession_status, price, product_schema, property_address, property_name,
				property_url, property_video, rera, size, updated_by_id, updated_date,
				usp, video_para, confiuration_id, developer_id, locality_id, project_id,
				user_id, highlights
			FROM property 
			WHERE project_id = ? AND is_deleted = 0
		`
		propertyRows, err := db.QueryContext(ctx, propertiesQuery, p.ID)
		if err != nil {
			log.Printf("Warning: error fetching properties for project %d: %v", p.ID, err)
		} else {
			defer propertyRows.Close()
			for propertyRows.Next() {
				var prop PropertyLegacyData
				var isDeletedBit, isFeaturedBit []byte
				err := propertyRows.Scan(
					&prop.ID, &prop.About, &prop.AgeOfProperty, &prop.AmenitiesPara,
					&prop.Balcony, &prop.Bathrooms, &prop.Bedrooms, &prop.BuiltupArea,
					&prop.CoverPhoto, &prop.CoveredParking, &prop.CreatedDate, &prop.Facing,
					&prop.FloorImage2D, &prop.FloorImage3D, &prop.FloorPara, &prop.Floors,
					&prop.FurnishingType, &prop.Images, &isDeletedBit, &isFeaturedBit,
					&prop.Latlong, &prop.ListingType, &prop.LocationMap, &prop.LocationAdvantage,
					&prop.LocationPara, &prop.LogoImage, &prop.MetaDescription, &prop.MetaKeywords,
					&prop.MetaTitle, &prop.OpenParking, &prop.OverviewPara, &prop.PossessionDate,
					&prop.PossessionStatus, &prop.Price, &prop.ProductSchema, &prop.PropertyAddress,
					&prop.PropertyName, &prop.PropertyURL, &prop.PropertyVideo, &prop.Rera,
					&prop.Size, &prop.UpdatedByID, &prop.UpdatedDate, &prop.USP, &prop.VideoPara,
					&prop.ConfiurationID, &prop.DeveloperID, &prop.LocalityID, &prop.ProjectID,
					&prop.UserID, &prop.Highlights,
				)
				if err != nil {
					log.Printf("Warning: error scanning property row: %v", err)
					continue
				}

				prop.IsDeleted = bitToBool(isDeletedBit)
				prop.IsFeatured = bitToBool(isFeaturedBit)

				// Add BHK option price to the price list if property name and price are available
				if prop.PropertyName != nil && prop.Price != 0 {
					bhkOption := BHKOptionPrice{
						Title: *prop.PropertyName,
						Price: prop.Price,
					}
					p.PriceList.BHKOptionPrices = append(p.PriceList.BHKOptionPrices, bhkOption)
				}

				// Fetch property amenities
				amenitiesQuery := `
					SELECT property_id, amenities_id 
					FROM property_amenities 
					WHERE property_id = ?
				`
				amenityRows, err := db.QueryContext(ctx, amenitiesQuery, prop.ID)
				if err != nil {
					log.Printf("Warning: error fetching amenities for property %d: %v", prop.ID, err)
				} else {
					defer amenityRows.Close()
					for amenityRows.Next() {
						var amenity PropertyAmenitiesLegacyData
						if err := amenityRows.Scan(&amenity.PropertyID, &amenity.AmenitiesID); err != nil {
							log.Printf("Warning: error scanning property amenity: %v", err)
							continue
						}
						prop.Amenities = append(prop.Amenities, amenity)
					}
				}

				p.Properties = append(p.Properties, prop)
			}
		}

		// Fetch project amenities
		projectAmenitiesQuery := `
			SELECT project_id, amenities_id 
			FROM project_amenities 
			WHERE project_id = ?
		`
		amenityRows, err := db.QueryContext(ctx, projectAmenitiesQuery, p.ID)
		if err != nil {
			log.Printf("Warning: error fetching amenities for project %d: %v", p.ID, err)
		} else {
			defer amenityRows.Close()
			for amenityRows.Next() {
				var amenity ProjectAmenitiesLegacyData
				if err := amenityRows.Scan(&amenity.ProjectID, &amenity.AmenitiesID); err != nil {
					log.Printf("Warning: error scanning project amenity: %v", err)
					continue
				}
				p.Amenities = append(p.Amenities, amenity)
			}
		}

		// Fetch FAQs for this project
		rows, err := db.QueryContext(ctx, `
			SELECT id, question, answer
			FROM faq
			WHERE project_id = ?
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching FAQs for project %d: %v", p.ID, err)
		}
		defer rows.Close()

		var faqCount int
		for rows.Next() {
			var faq FAQLegacyData
			if err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer); err != nil {
				return nil, fmt.Errorf("error scanning FAQ row: %v", err)
			}
			p.FAQs = append(p.FAQs, faq)
			faqCount++
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating FAQ rows: %v", err)
		}

		// Fetch amenities for this project
		amenityRows, err = db.QueryContext(ctx, `
			SELECT 
				a.id, a.amenities_category, a.amenities_name, 
				a.amenities_url, a.created_date, a.updated_date
			FROM 
				amenities a
				INNER JOIN project_amenities pa ON pa.amenities_id = a.id
			WHERE 
				pa.project_id = ?
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching amenities for project %d: %v", p.ID, err)
		}
		defer amenityRows.Close()

		for amenityRows.Next() {
			var amenity AmenitiesData
			var createdDate, updatedDate sql.NullInt64

			if err := amenityRows.Scan(
				&amenity.ID,
				&amenity.AmenitiesCategory,
				&amenity.AmenitiesName,
				&amenity.AmenitiesURL,
				&createdDate,
				&updatedDate,
			); err != nil {
				return nil, fmt.Errorf("error scanning amenity row: %v", err)
			}

			if createdDate.Valid {
				amenity.CreatedDate = createdDate.Int64
			}
			if updatedDate.Valid {
				amenity.UpdatedDate = updatedDate.Int64
			}

			p.AmenitiesData = append(p.AmenitiesData, amenity)
		}
		if err := amenityRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating amenity rows: %v", err)
		}

		// Fetch payment plans for this project
		paymentRows, err := db.QueryContext(ctx, `
			SELECT id, payment_plan_name, payment_plan_value, project_id
			FROM payment_plan
			WHERE project_id = ?
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching payment plans for project %d: %v", p.ID, err)
		}
		defer paymentRows.Close()

		for paymentRows.Next() {
			var plan PaymentPlanData
			if err := paymentRows.Scan(
				&plan.ID,
				&plan.PaymentPlanName,
				&plan.PaymentPlanValue,
				&plan.ProjectID,
			); err != nil {
				return nil, fmt.Errorf("error scanning payment plan row: %v", err)
			}
			p.PaymentPlans = append(p.PaymentPlans, plan)
		}
		if err := paymentRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating payment plan rows: %v", err)
		}

		// Fetch floor plans for this project
		floorRows, err := db.QueryContext(ctx, `
			SELECT 
				id, created_date, img_url, is_sold_out, price, 
				size, title, updated_date, configuration_id, 
				project_id, user_id
			FROM floorplan
			WHERE project_id = ?
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching floor plans for project %d: %v", p.ID, err)
		}
		defer floorRows.Close()

		for floorRows.Next() {
			var plan FloorPlanData
			var createdDate, updatedDate, size sql.NullInt64
			var imgURL, title sql.NullString
			var isSoldOutBit []byte

			if err := floorRows.Scan(
				&plan.ID,
				&createdDate,
				&imgURL,
				&isSoldOutBit,
				&plan.Price,
				&size,
				&title,
				&updatedDate,
				&plan.ConfigurationID,
				&plan.ProjectID,
				&plan.UserID,
			); err != nil {
				return nil, fmt.Errorf("error scanning floor plan row: %v", err)
			}

			// Handle nullable fields
			if createdDate.Valid {
				plan.CreatedDate = createdDate.Int64
			}
			if updatedDate.Valid {
				plan.UpdatedDate = updatedDate.Int64
			}
			if size.Valid {
				plan.Size = size.Int64
			}
			if imgURL.Valid {
				plan.ImgURL = imgURL.String
			}
			if title.Valid {
				plan.Title = title.String
			}
			plan.IsSoldOut = bitToBool(isSoldOutBit)

			p.FloorPlans = append(p.FloorPlans, plan)
		}
		if err := floorRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating floor plan rows: %v", err)
		}

		// Fetch RERA info
		if err := p.fetchReraInfo(db); err != nil {
			return nil, fmt.Errorf("error fetching RERA info: %v", err)
		}

		// Fetch project images
		imageRows, err := db.QueryContext(ctx, `
			SELECT project_id, image_alt_name, image_url
			FROM project_image
			WHERE project_id = ?
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching images for project %d: %v", p.ID, err)
		}
		defer imageRows.Close()

		for imageRows.Next() {
			var image ProjectImageLegacyData
			var imageAltName sql.NullString
			if err := imageRows.Scan(&image.ProjectID, &imageAltName, &image.ImageURL); err != nil {
				return nil, fmt.Errorf("error scanning project image row: %v", err)
			}
			if imageAltName.Valid {
				image.ImageAltName = imageAltName.String
			}
			p.ProjectImages = append(p.ProjectImages, image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating project image rows: %v", err)
		}

		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %v", err)
	}

	return projects, nil
}

// CheckNullFields checks and prints any null fields in the legacy project data
func CheckNullFields(projects []ProjectLegacyData) {
	for _, project := range projects {
		nullFields := make([]string, 0)

		// Check Project fields
		if project.ID == 0 {
			nullFields = append(nullFields, "ID")
		}
		if project.ProjectName == "" {
			nullFields = append(nullFields, "ProjectName")
		}
		if project.ProjectDescription == "" {
			nullFields = append(nullFields, "ProjectDescription")
		}
		if project.ProjectArea == "" {
			nullFields = append(nullFields, "ProjectArea")
		}
		if project.ProjectUnits == "" {
			nullFields = append(nullFields, "ProjectUnits")
		}
		if project.ProjectConfigurations == "" {
			nullFields = append(nullFields, "ProjectConfigurations")
		}
		if project.TotalFloor == "" {
			nullFields = append(nullFields, "TotalFloor")
		}
		if project.TotalTowers == "" {
			nullFields = append(nullFields, "TotalTowers")
		}
		if project.Status == "" {
			nullFields = append(nullFields, "Status")
		}
		if project.ProjectLaunchDate == "" {
			nullFields = append(nullFields, "ProjectLaunchDate")
		}
		if project.ProjectPossessionDate == "" {
			nullFields = append(nullFields, "ProjectPossessionDate")
		}
		if project.MetaTitle == "" {
			nullFields = append(nullFields, "MetaTitle")
		}
		if project.MetaDescription == "" {
			nullFields = append(nullFields, "MetaDescription")
		}
		if project.ProjectSchema == "" {
			nullFields = append(nullFields, "ProjectSchema")
		}
		if project.ShortAddress == "" {
			nullFields = append(nullFields, "ShortAddress")
		}
		if project.ProjectLogo == "" {
			nullFields = append(nullFields, "ProjectLogo")
		}
		if project.CoverPhoto == "" {
			nullFields = append(nullFields, "CoverPhoto")
		}
		if project.RERALink == "" {
			nullFields = append(nullFields, "RERALink")
		}
		if project.ProjectRERA == "" {
			nullFields = append(nullFields, "ProjectRERA")
		}
		if project.ProjectBrochure == "" {
			nullFields = append(nullFields, "ProjectBrochure")
		}
		if project.OverviewPara == "" {
			nullFields = append(nullFields, "OverviewPara")
		}
		if project.PriceListPara == "" {
			nullFields = append(nullFields, "PriceListPara")
		}
		if project.AmenitiesPara == "" {
			nullFields = append(nullFields, "AmenitiesPara")
		}
		if project.VideoPara == "" {
			nullFields = append(nullFields, "VideoPara")
		}
		if project.PaymentPara == "" {
			nullFields = append(nullFields, "PaymentPara")
		}
		if project.SitePlanPara == "" {
			nullFields = append(nullFields, "SitePlanPara")
		}
		if project.SitePlanImg == "" {
			nullFields = append(nullFields, "SitePlanImg")
		}
		if project.ProjectAbout == "" {
			nullFields = append(nullFields, "ProjectAbout")
		}
		if project.ProjectAddress == "" {
			nullFields = append(nullFields, "ProjectAddress")
		}
		if project.USP == "" {
			nullFields = append(nullFields, "USP")
		}

		// Check FAQs
		if len(project.FAQs) == 0 {
			nullFields = append(nullFields, "FAQs (empty array)")
		} else {
			for i, faq := range project.FAQs {
				if faq.Question == "" {
					nullFields = append(nullFields, fmt.Sprintf("FAQ[%d].Question", i))
				}
				if faq.Answer == "" {
					nullFields = append(nullFields, fmt.Sprintf("FAQ[%d].Answer", i))
				}
			}
		}

		// Print null fields if any found
		if len(nullFields) > 0 {
			log.Printf("Project ID %d (%s) has following null/empty fields:", project.ID, project.ProjectName)
			for _, field := range nullFields {
				log.Printf("- %s", field)
			}
			log.Println("-------------------")
		}
	}
}

// fetchReraInfo fetches RERA info data for a project
func (p *ProjectLegacyData) fetchReraInfo(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT id, created_on, phase, project_rera_name, qr_images,
			   rera_number, status, updated_on, project_id, user_id
		FROM rera_info
		WHERE project_id = ?
	`, p.ID)
	if err != nil {
		return fmt.Errorf("error fetching RERA info: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reraInfo ReraInfoData
		var phaseJSON, qrImagesJSON, reraNumberJSON, statusJSON []byte
		var projectReraName sql.NullString
		var createdOn, updatedOn, userID sql.NullInt64
		err := rows.Scan(
			&reraInfo.ID, &createdOn, &phaseJSON,
			&projectReraName, &qrImagesJSON,
			&reraNumberJSON, &statusJSON,
			&updatedOn, &reraInfo.ProjectID,
			&userID,
		)
		if err != nil {
			return fmt.Errorf("error scanning RERA info: %v", err)
		}

		// Handle nullable fields
		if projectReraName.Valid {
			reraInfo.ProjectReraName = projectReraName.String
		}
		if createdOn.Valid {
			reraInfo.CreatedOn = createdOn.Int64
		}
		if updatedOn.Valid {
			reraInfo.UpdatedOn = updatedOn.Int64
		}
		if userID.Valid {
			reraInfo.UserID = userID.Int64
		}

		// Initialize empty slices for JSON fields
		reraInfo.Phase = make([]string, 0)
		reraInfo.QRImages = make([]string, 0)
		reraInfo.ReraNumber = make([]string, 0)
		reraInfo.Status = make([]string, 0)

		// Helper function to check and clean JSON data
		cleanJSON := func(data []byte) []byte {
			if len(data) == 0 || string(data) == "NULL" {
				return []byte("[]")
			}
			// If it's not already an array and not already JSON-formatted
			trimmed := string(bytes.TrimSpace(data))
			if !strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "\"") {
				// Quote the string and wrap it in array brackets
				return []byte(fmt.Sprintf("[\"%s\"]", strings.ReplaceAll(trimmed, "\"", "\\\"")))
			}
			return data
		}

		// Unmarshal JSON arrays with NULL checks and cleaning
		if len(phaseJSON) > 0 {
			if err := json.Unmarshal(cleanJSON(phaseJSON), &reraInfo.Phase); err != nil {
				log.Printf("Raw phase JSON: %s", string(phaseJSON))
				return fmt.Errorf("error unmarshaling phase: %v", err)
			}
		}
		if len(qrImagesJSON) > 0 {
			if err := json.Unmarshal(cleanJSON(qrImagesJSON), &reraInfo.QRImages); err != nil {
				log.Printf("Raw qr_images JSON: %s", string(qrImagesJSON))
				return fmt.Errorf("error unmarshaling qr_images: %v", err)
			}
		}
		if len(reraNumberJSON) > 0 {
			if err := json.Unmarshal(cleanJSON(reraNumberJSON), &reraInfo.ReraNumber); err != nil {
				log.Printf("Raw rera_number JSON: %s", string(reraNumberJSON))
				return fmt.Errorf("error unmarshaling rera_number: %v", err)
			}
		}
		if len(statusJSON) > 0 {
			if err := json.Unmarshal(cleanJSON(statusJSON), &reraInfo.Status); err != nil {
				log.Printf("Raw status JSON: %s", string(statusJSON))
				return fmt.Errorf("error unmarshaling status: %v", err)
			}
		}

		p.ReraInfo = append(p.ReraInfo, reraInfo)
	}
	return rows.Err()
}
