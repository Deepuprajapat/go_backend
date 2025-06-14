package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/developer"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/property"
	"github.com/VI-IM/im_backend_go/ent/schema"
)

// Helper function to generate URL-friendly identifier
func generateIdentifier(name string) string {
	// Convert to lowercase
	identifier := strings.ToLower(name)
	// Replace spaces with hyphens
	identifier = strings.ReplaceAll(identifier, " ", "-")
	// Remove any special characters
	identifier = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, identifier)
	return identifier
}

// parseUSPs parses a JSON array of USPs from the legacy database
func parseUSPs(uspString string) []struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
} {
	var uspArray []string
	if err := json.Unmarshal([]byte(uspString), &uspArray); err != nil {
		// If parsing fails, try treating it as a single USP
		return []struct {
			Icon        string `json:"icon"`
			Description string `json:"description"`
		}{{
			Description: strings.TrimSpace(uspString),
		}}
	}

	// Convert each USP string to the new structure
	usps := make([]struct {
		Icon        string `json:"icon"`
		Description string `json:"description"`
	}, len(uspArray))

	for i, usp := range uspArray {
		usps[i] = struct {
			Icon        string `json:"icon"`
			Description string `json:"description"`
		}{
			Description: strings.TrimSpace(usp),
		}
	}
	return usps
}

// Helper function to organize amenities by category
func organizeAmenitiesByCategory(amenities []AmenitiesData) map[string][]struct {
	Icon  string `json:"icon"`
	Value string `json:"value"`
} {
	categorized := make(map[string][]struct {
		Icon  string `json:"icon"`
		Value string `json:"value"`
	})

	for _, amenity := range amenities {
		category := amenity.AmenitiesCategory
		if category == "" {
			category = "Other"
		}

		categorized[category] = append(categorized[category], struct {
			Icon  string `json:"icon"`
			Value string `json:"value"`
		}{
			Icon:  amenity.AmenitiesURL,
			Value: amenity.AmenitiesName,
		})
	}

	return categorized
}

// Convert payment plans to the new format
func convertPaymentPlans(plans []PaymentPlanData) []struct {
	Name    string `json:"name"`
	Details string `json:"details"`
} {
	result := make([]struct {
		Name    string `json:"name"`
		Details string `json:"details"`
	}, len(plans))

	for i, plan := range plans {
		result[i] = struct {
			Name    string `json:"name"`
			Details string `json:"details"`
		}{
			Name:    plan.PaymentPlanName,
			Details: plan.PaymentPlanValue,
		}
	}

	return result
}

// Convert floor plans to the new format
func convertFloorPlans(plans []FloorPlanData, configName string) []struct {
	Title        string `json:"title"`
	FlatType     string `json:"flat_type"`
	Price        string `json:"price"`
	BuildingArea string `json:"building_area"`
	Image        string `json:"image"`
} {
	result := make([]struct {
		Title        string `json:"title"`
		FlatType     string `json:"flat_type"`
		Price        string `json:"price"`
		BuildingArea string `json:"building_area"`
		Image        string `json:"image"`
	}, len(plans))

	re := regexp.MustCompile(`(\d+)\s*BHK`)

	for i, plan := range plans {
		// Extract BHK type from title using regex
		flatType := "4BHK" // Default value
		if matches := re.FindStringSubmatch(plan.Title); len(matches) > 1 {
			flatType = matches[1] + "BHK"
		}

		result[i] = struct {
			Title        string `json:"title"`
			FlatType     string `json:"flat_type"`
			Price        string `json:"price"`
			BuildingArea string `json:"building_area"`
			Image        string `json:"image"`
		}{
			Title:        plan.Title,
			FlatType:     flatType,
			Price:        fmt.Sprintf("%.2f", plan.Price),
			BuildingArea: fmt.Sprintf("%d", plan.Size),
			Image:        plan.ImgURL,
		}
	}

	return result
}

// MigrateCommonFields migrates common fields from legacy data to new schema
func MigrateCommonFields(ctx context.Context, client *ent.Client, legacyProjects []ProjectLegacyData) error {
	// Check for null fields in legacy data
	log.Println("Checking for null fields in legacy data...")
	CheckNullFields(legacyProjects)
	log.Println("Null field check completed")

	// First migrate developers since projects depend on them
	developerMap := make(map[string]*ent.Developer) // Map by developer name instead of ID
	for _, lp := range legacyProjects {
		if lp.Developer == nil {
			continue
		}

		// Generate identifier from developer name if URL is not available
		identifier := lp.Developer.DeveloperURL
		if identifier == "" {
			identifier = generateIdentifier(lp.Developer.DeveloperName)
		}

		// Create developer with common fields
		dev, err := client.Developer.Create().
			SetName(lp.Developer.DeveloperName).
			SetLegalName(lp.Developer.DeveloperLegalName).
			SetIdentifier(identifier).
			SetEstablishedYear(int(lp.Developer.EstablishedYear)).
			SetIsVerified(lp.Developer.IsVerified).
			SetMediaContent(schema.DeveloperMediaContent{
				DeveloperAddress: lp.Developer.DeveloperAddress,
				Phone:            lp.Developer.Phone,
				DeveloperLogo:    lp.Developer.DeveloperLogo,
				AltDeveloperLogo: lp.Developer.AltDeveloperLogo,
				About:            lp.Developer.About,
				Overview:         lp.Developer.Overview,
				Disclaimer:       lp.Developer.Disclaimer,
			}).
			Save(ctx)

		if err != nil {
			if ent.IsConstraintError(err) {
				// Developer already exists, fetch it
				dev, err = client.Developer.Query().
					Where(developer.Name(lp.Developer.DeveloperName)).
					Only(ctx)
				if err != nil {
					return fmt.Errorf("error fetching existing developer: %v", err)
				}
			} else {
				return fmt.Errorf("error creating developer: %v", err)
			}
		}

		// Store the developer by name for later reference
		developerMap[lp.Developer.DeveloperName] = dev
	}

	// Now migrate projects and their properties
	for _, lp := range legacyProjects {
		// log.Println("Images", lp.ProjectImages)

		// Convert FAQs to the new format
		var faqs []struct {
			Question string `json:"question"`
			Answer   string `json:"answer"`
		}
		for _, faq := range lp.FAQs {
			if faq.Question != "" && faq.Answer != "" {
				cleanQuestion := strings.TrimSpace(faq.Question)
				cleanAnswer := strings.TrimSpace(faq.Answer)
				faqs = append(faqs, struct {
					Question string `json:"question"`
					Answer   string `json:"answer"`
				}{
					Question: cleanQuestion,
					Answer:   cleanAnswer,
				})
			}
		}

		// Create web cards
		webCards := schema.ProjectWebCards{
			Images: make([]string, 0), // Initialize Images slice
			ReraInfo: schema.ReraInfo{
				WebsiteLink: lp.RERALink,
				ReraList: make([]struct {
					Phase      string `json:"phase"`
					ReraQR     string `json:"rera_qr"`
					ReraNumber string `json:"rera_number"`
					Status     string `json:"status"`
				}, 0),
			},
			Details: schema.ProjectDetails{
				Area: struct {
					Value string `json:"value"`
				}{
					Value: lp.ProjectArea,
				},
				Configuration: struct {
					Value string `json:"value"`
				}{
					Value: lp.ProjectConfigurations,
				},
				Units: struct {
					Value string `json:"value"`
				}{
					Value: lp.ProjectUnits,
				},
				LaunchDate: struct {
					Value string `json:"value"`
				}{
					Value: lp.ProjectLaunchDate,
				},
				PossessionDate: struct {
					Value string `json:"value"`
				}{
					Value: lp.ProjectPossessionDate,
				},
				TotalTowers: struct {
					Value string `json:"value"`
				}{
					Value: lp.TotalTowers,
				},
				TotalFloors: struct {
					Value string `json:"value"`
				}{
					Value: lp.TotalFloor,
				},
				ProjectStatus: struct {
					Value string `json:"value"`
				}{
					Value: lp.Status,
				},
				Type: struct {
					Value string `json:"value"`
				}{
					Value: func() string {
						if lp.ConfigType != nil {
							return lp.ConfigType.PropertyType
						}
						return ""
					}(),
				},
			},
			WhyToChoose: schema.WhyToChoose{
				ImageUrls: func() []string {
					urls := make([]string, 0)
					for _, img := range lp.ProjectImages {
						urls = append(urls, img.ImageURL)
					}
					return urls
				}(),
				USP_List: parseUSPs(lp.USP),
			},
			KnowAbout: schema.KnowAbout{
				Description:  lp.ProjectAbout,
				DownloadLink: lp.ProjectBrochure,
			},
			FloorPlan: schema.FloorPlan{
				Title:    "Floor Plans",
				Products: convertFloorPlans(lp.FloorPlans, lp.Configuration.ProjectConfigurationName),
			},
			PriceList: schema.PriceList{
				Description: lp.PriceListPara,
			},
			Amenities: schema.Amenities{
				Description:             lp.AmenitiesPara,
				CategoriesWithAmenities: organizeAmenitiesByCategory(lp.AmenitiesData),
			},
			VideoPresentation: schema.VideoPresentation{
				Description: lp.VideoPara,
				URL: func() string {
					if len(lp.ProjectVideos) == 0 {
						log.Printf("Project ID %d: No video data found", lp.ID)
						return ""
					}

					// If we have ProjectVideo structs, use the first one's URL
					if len(lp.ProjectVideos) > 0 {
						// log.Printf("Project ID %d: Found video URL from struct: %s", lp.ID, lp.ProjectVideos[0].VideoURL)
						return lp.ProjectVideos[0].VideoURL
					}

					log.Printf("Project ID %d: No valid video URL found in data", lp.ID)
					return ""
				}(),
			},
			PaymentPlans: schema.PaymentPlans{
				Description: lp.PaymentPara,
				Plans:       convertPaymentPlans(lp.PaymentPlans),
			},
			SitePlan: struct {
				Description string `json:"description"`
				Image       string `json:"image"`
			}{
				Description: lp.SitePlanPara,
				Image:       lp.SitePlanImg,
			},
			About: schema.About{
				Description:       lp.Developer.About,
				LogoURL:           lp.Developer.DeveloperLogo,
				EstablishmentYear: fmt.Sprintf("%d", lp.Developer.EstablishedYear),
				TotalProjects:     fmt.Sprintf("%d", lp.Developer.ProjectDoneNo),
				ContactDetails: struct {
					Name           string `json:"name"`
					ProjectAddress string `json:"project_address"`
					Phone          string `json:"phone"`
					BookingLink    string `json:"booking_link"`
				}{
					Name:           lp.Developer.DeveloperName,
					ProjectAddress: lp.ProjectAddress,
					Phone:          lp.Developer.Phone,
					BookingLink:    lp.ProjectBrochure,
				},
			},
			Faqs: faqs,
		}

		// First add project logo if it exists
		if lp.ProjectLogo != "" {
			webCards.Images = append(webCards.Images, lp.ProjectLogo)
		}

		// Then add all project images
		for _, img := range lp.ProjectImages {
			if img.ImageURL != lp.ProjectLogo {
				webCards.Images = append(webCards.Images, img.ImageURL)
			}
		}

		// Debug logging for project images

		// Marshal webCards to see the exact JSON SELECT id, web_cards FROM projects LIMIT 5;that will be store

		// Add RERA info from legacy data
		for _, reraInfo := range lp.ReraInfo {
			// For each phase, create a separate RERA entry
			for i := range reraInfo.Phase {
				webCards.ReraInfo.ReraList = append(webCards.ReraInfo.ReraList, struct {
					Phase      string `json:"phase"`
					ReraQR     string `json:"rera_qr"`
					ReraNumber string `json:"rera_number"`
					Status     string `json:"status"`
				}{
					Phase:      reraInfo.Phase[i],
					ReraQR:     getArrayElement(reraInfo.QRImages, i, "N/A"),
					ReraNumber: getArrayElement(reraInfo.ReraNumber, i, ""),
					Status:     getArrayElement(reraInfo.Status, i, ""),
				})
			}
		}

		// Create project with common fields
		projectCreate := client.Project.Create().
			SetBasicInfo(schema.BasicInfo{
				ProjectName:           lp.ProjectName,
				ProjectDescription:    lp.ProjectDescription,
				ProjectArea:           lp.ProjectArea,
				ProjectUnits:          lp.ProjectUnits,
				ProjectConfigurations: lp.ProjectConfigurations,
				TotalFloor:            lp.TotalFloor,
				TotalTowers:           lp.TotalTowers,
				Status:                lp.Status,
			}).
			SetTimelineInfo(schema.TimelineInfo{
				ProjectLaunchDate:     lp.ProjectLaunchDate,
				ProjectPossessionDate: lp.ProjectPossessionDate,
			}).
			SetMetaInfo(schema.SEOMeta{
				Title:         lp.MetaTitle,
				Description:   lp.MetaDescription,
				Keywords:      strings.Join(lp.MetaKeywords, ", "),
				Canonical:     lp.ProjectURL,
				ProjectSchema: lp.ProjectSchema,
			}).
			SetWebCards(webCards).
			SetLocationInfo(schema.LocationInfo{
				ShortAddress:  lp.ShortAddress,
				GoogleMapLink: lp.ProjectLocationURL,
			})

		// Debug: Print the project create mutation

		// Set other fields
		projectCreate = projectCreate.
			SetIsFeatured(lp.IsFeatured).
			SetIsPremium(lp.IsPremium).
			SetIsPriority(lp.IsPriority).
			SetIsDeleted(lp.IsDeleted).
			SetSearchContext([]string{
				lp.ProjectName,
				lp.ProjectDescription,
				lp.ShortAddress,
				lp.ProjectArea,
				lp.ProjectUnits,
				lp.ProjectConfigurations,
				lp.Status,
			})

		// Set developer edge if exists
		if dev, ok := developerMap[lp.Developer.DeveloperName]; ok {
			projectCreate.SetDeveloper(dev)
		}

		// Save project
		proj, err := projectCreate.Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				// Project already exists, update it
				proj, err = client.Project.Query().
					Where(
						project.And(
							project.IsFeatured(lp.IsFeatured),
							project.IsPremium(lp.IsPremium),
							project.IsPriority(lp.IsPriority),
						),
					).
					Only(ctx)
				if err != nil {
					return fmt.Errorf("error fetching existing project: %v", err)
				}

				// Debug: Print the project before update
				// log.Printf("Existing project before update: %+v", proj)

				// Update the project with new data
				_, err = proj.Update().
					SetBasicInfo(schema.BasicInfo{
						ProjectName:           lp.ProjectName,
						ProjectDescription:    lp.ProjectDescription,
						ProjectArea:           lp.ProjectArea,
						ProjectUnits:          lp.ProjectUnits,
						ProjectConfigurations: lp.ProjectConfigurations,
						TotalFloor:            lp.TotalFloor,
						TotalTowers:           lp.TotalTowers,
						Status:                lp.Status,
					}).
					SetTimelineInfo(schema.TimelineInfo{
						ProjectLaunchDate:     lp.ProjectLaunchDate,
						ProjectPossessionDate: lp.ProjectPossessionDate,
					}).
					SetMetaInfo(schema.SEOMeta{
						Title:         lp.MetaTitle,
						Description:   lp.MetaDescription,
						Keywords:      strings.Join(lp.MetaKeywords, ", "),
						ProjectSchema: lp.ProjectSchema,
					}).
					SetWebCards(webCards).
					SetLocationInfo(schema.LocationInfo{
						ShortAddress:  lp.ShortAddress,
						GoogleMapLink: lp.ProjectLocationURL,
					}).
					SetIsFeatured(lp.IsFeatured).
					SetIsPremium(lp.IsPremium).
					SetIsPriority(lp.IsPriority).
					SetIsDeleted(lp.IsDeleted).
					SetSearchContext([]string{
						lp.ProjectName,
						lp.ProjectDescription,
						lp.ShortAddress,
						lp.ProjectArea,
						lp.ProjectUnits,
						lp.ProjectConfigurations,
						lp.Status,
					}).
					Save(ctx)
				if err != nil {
					return fmt.Errorf("error updating project: %v", err)
				}

				// Debug: Print the project after update
				updatedProj, _ := client.Project.Get(ctx, proj.ID)
				if updatedProj != nil {
					log.Printf("Project after update: %+v", updatedProj)
					log.Printf("Updated project web cards: %+v", updatedProj.WebCards)
				}
			} else {
				return fmt.Errorf("error creating project: %v", err)
			}
		}

		// Migrate properties for this project
		for _, lp := range lp.Properties {
			// Convert property amenities to the new format
			var amenities []struct {
				Icon string `json:"icon"`
				Name string `json:"name"`
			}
			for _, amenity := range lp.Amenities {
				amenities = append(amenities, struct {
					Icon string `json:"icon"`
					Name string `json:"name"`
				}{
					Icon: "",                                     // Icon not available in legacy data
					Name: fmt.Sprintf("%d", amenity.AmenitiesID), // Using ID as name since we don't have the actual name
				})
			}

			// Create property with common fields
			propertyCreate := client.Property.Create().
				SetName(stringOrEmpty(lp.PropertyName)).
				SetDescription(stringOrEmpty(lp.About)).
				SetPropertyImages(schema.PropertyImages{
					Images: []struct {
						Order int    `json:"order"`
						Url   string `json:"url"`
						Type  string `json:"type"`
					}{
						{
							Order: 1,
							Url:   stringOrEmpty(lp.CoverPhoto),
							Type:  "cover",
						},
					},
				}).
				SetWebCards(schema.WebCards{
					PropertyDetails: schema.PropertyDetails{
						PropertyType:      "", // Not available in legacy data
						FurnishingType:    stringOrEmpty(lp.FurnishingType),
						ListingType:       stringOrEmpty(lp.ListingType),
						PossessionStatus:  stringOrEmpty(lp.PossessionStatus),
						AgeOfProperty:     stringOrEmpty(lp.AgeOfProperty),
						FloorPara:         stringOrEmpty(lp.FloorPara),
						LocationPara:      stringOrEmpty(lp.LocationPara),
						LocationAdvantage: stringOrEmpty(lp.LocationAdvantage),
						OverviewPara:      stringOrEmpty(lp.OverviewPara),
						Floors:            stringOrEmpty(lp.Floors),
						Images:            stringOrEmpty(lp.Images),
						Latlong:           stringOrEmpty(lp.Latlong),
					},
					PropertyFloorPlan: []struct {
						Title string `json:"title"`
						Plans []struct {
							Title        string `json:"title"`
							FlatType     string `json:"flat_type"`
							Price        string `json:"price"`
							BuildingArea string `json:"building_area"`
							Image        string `json:"image"`
							ExpertLink   string `json:"expert_link"`
							BrochureLink string `json:"brochure_link"`
						} `json:"plans"`
					}{},
					KnowAbout: struct {
						HtmlText string `json:"html_text"`
					}{
						HtmlText: stringOrEmpty(lp.About),
					},
					VideoPresentation: struct {
						Title    string `json:"title"`
						VideoUrl string `json:"video_url"`
					}{
						Title:    "",
						VideoUrl: stringOrEmpty(lp.PropertyVideo),
					},
					GoogleMapLink: stringOrEmpty(lp.LocationMap),
				}).
				SetBasicInfo(schema.PropertyBasicInfo{
					PropertyType: "", // Not available in legacy data
					BHKType:      "", // Not available in legacy data
					Bedrooms:     stringToInt(lp.Bedrooms),
					Bathrooms:    stringToInt(lp.Bathrooms),
				}).
				SetLocationDetails(schema.PropertyLocationDetails{
					FloorNumber: 0, // Not available in legacy data
					Facing:      stringOrEmpty(lp.Facing),
					Tower:       "", // Not available in legacy data
					Wing:        "", // Not available in legacy data
				}).
				SetPricingInfo(schema.PropertyPricingInfo{
					StartingPrice: fmt.Sprintf("%.2f", lp.Price),
					Price:         fmt.Sprintf("%.2f", lp.Price),
				}).
				SetPropertyReraInfo(schema.PropertyReraInfo{
					Phase:      "", // Not available in legacy data
					Status:     "", // Not available in legacy data
					ReraNumber: stringOrEmpty(lp.Rera),
					ReraQR:     "", // Not available in legacy data
				}).
				SetSearchContext([]string{
					stringOrEmpty(lp.PropertyName),
					stringOrEmpty(lp.About),
					stringOrEmpty(lp.PropertyAddress),
					stringOrEmpty(lp.Facing),
					fmt.Sprintf("%.2f", lp.Price),
				})

			// Add all project images to property images
			propertyImages := schema.PropertyImages{
				Images: make([]struct {
					Order int    `json:"order"`
					Url   string `json:"url"`
					Type  string `json:"type"`
				}, 0),
			}

			// Add cover photo first
			coverPhoto := stringOrEmpty(lp.CoverPhoto)
			if coverPhoto != "" {
				propertyImages.Images = append(propertyImages.Images, struct {
					Order int    `json:"order"`
					Url   string `json:"url"`
					Type  string `json:"type"`
				}{
					Order: 1,
					Url:   coverPhoto,
					Type:  "cover",
				})
			}

			// Add all project images
			for i, img := range proj.WebCards.Images {
				propertyImages.Images = append(propertyImages.Images, struct {
					Order int    `json:"order"`
					Url   string `json:"url"`
					Type  string `json:"type"`
				}{
					Order: i + 2, // Start from 2 since cover photo is 1
					Url:   img,
					Type:  "project", // Using "project" type for all project images
				})
			}

			propertyCreate.SetPropertyImages(propertyImages)

			// Set project edge
			propertyCreate.SetProject(proj)

			// Save property
			_, err := propertyCreate.Save(ctx)
			if err != nil {
				if ent.IsConstraintError(err) {
					// Property already exists, update it
					_, err = client.Property.Query().
						Where(property.Name(stringOrEmpty(lp.PropertyName))).
						Only(ctx)
					if err != nil {
						return fmt.Errorf("error fetching existing property: %v", err)
					}
				} else {
					return fmt.Errorf("error creating property: %v", err)
				}
			}
		}

		log.Printf("Successfully migrated project %d with its properties", lp.ID)
	}

	return nil
}

// Helper function to handle nullable strings
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to convert string to int
func stringToInt(s *string) int {
	if s == nil {
		return 0
	}
	val := 0
	fmt.Sscanf(*s, "%d", &val)
	return val
}

func getArrayElement(arr []string, index int, defaultValue string) string {
	if index < len(arr) {
		return arr[index]
	}
	return defaultValue
}
