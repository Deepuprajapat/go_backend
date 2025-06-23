package migration_jobs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/rs/zerolog/log"
)

func safeStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func safeInt(i *int64) int {
	if i != nil {
		return int(*i)
	}
	return 0
}

func parsePhoneJSONToString(s *string) (*string, error) {
	if s == nil {
		var zero string = ""
		log.Info().Msgf("Phone string is nil, returning 0")
		return &zero, nil
	}

	var phones []string
	if err := json.Unmarshal([]byte(*s), &phones); err != nil {
		return nil, fmt.Errorf("failed to unmarshal phone JSON: %w", err)
	}

	if len(phones) == 0 {
		return nil, fmt.Errorf("no phone number found")
	}

	// Clean up: remove dashes and whitespace
	phoneClean := strings.ReplaceAll(phones[0], "-", "")
	phoneClean = strings.TrimSpace(phoneClean)

	// Take only first 10 digits to avoid overflow
	if len(phoneClean) > 10 {
		phoneClean = phoneClean[len(phoneClean)-10:]
	}

	return &phoneClean, nil
}

func parsePropertyImagesFromProjectImages(projectImages *[]LProjectImage) (*schema.PropertyImages, error) {
	propertyImages := schema.PropertyImages{
		Images: []struct {
			Order int    `json:"order"`
			Url   string `json:"url"`
			Type  string `json:"type"`
		}{},
	}
	for _, image := range *projectImages {
		propertyImages.Images = append(propertyImages.Images, struct {
			Order int    `json:"order"`
			Url   string `json:"url"`
			Type  string `json:"type"`
		}{
			Order: 1,
			Url:   image.ImageURL,
			Type:  "property_image",
		})
	}
	return &propertyImages, nil
}

func parseWebCardsFromProject(project *LProject, floorPlans *[]LFloorPlan, property *LProperty) (*schema.WebCards, error) {

	floorPlanItems := []schema.FloorPlanItem{}
	for _, floorPlan := range *floorPlans {
		floorPlanItems = append(floorPlanItems, schema.FloorPlanItem{
			Title: safeStr(floorPlan.Title),
			FlatType: func() string {
				if floorPlan.ConfigurationID != nil {
					return fmt.Sprintf("%d", *floorPlan.ConfigurationID)
				}
				return ""
			}(),
			Price: fmt.Sprintf("%v", floorPlan.Price),
			BuildingArea: func() string {
				if floorPlan.Size != nil {
					return fmt.Sprintf("%d", *floorPlan.Size)
				}
				return ""
			}(),
			Image: safeStr(floorPlan.ImgURL),
		})
	}
	webCards := schema.WebCards{
		PropertyDetails: schema.PropertyDetails{
			PropertyType:      safeStr(project.ProjectConfigurations),
			FurnishingType:    safeStr(property.FurnishingType), // Not available in legacy data
			ListingType:       safeStr(property.ListingType),    // Not available in legacy data
			PossessionStatus:  safeStr(property.PossessionStatus),
			AgeOfProperty:     safeStr(property.AgeOfProperty), // Not available in legacy data
			FloorPara:         safeStr(property.FloorPara),
			LocationPara:      safeStr(property.LocationPara),
			LocationAdvantage: safeStr(property.LocationAdvantage), // Not available in legacy data
			OverviewPara:      safeStr(property.OverviewPara),
			Floors:            safeStr(property.Floors),
			Images:            safeStr(property.Images), // Will be populated from project images if available
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
		}{
			{
				Title: safeStr(project.FloorPara),
				Plans: []struct {
					Title        string `json:"title"`
					FlatType     string `json:"flat_type"`
					Price        string `json:"price"`
					BuildingArea string `json:"building_area"`
					Image        string `json:"image"`
					ExpertLink   string `json:"expert_link"`
					BrochureLink string `json:"brochure_link"`
				}{
					{
						Title:        safeStr(project.ProjectConfigurations),
						FlatType:     safeStr(project.ProjectConfigurations),
						Price:        "",
						BuildingArea: safeStr(project.ProjectArea),
						Image:        "",
						ExpertLink:   "",
						BrochureLink: safeStr(project.ProjectBrochure),
					},
				},
			},
		},
		KnowAbout: struct {
			HtmlText string `json:"html_text"`
		}{
			HtmlText: safeStr(project.ProjectAbout),
		},
		VideoPresentation: struct {
			Title    string `json:"title"`
			VideoUrl string `json:"video_url"`
		}{
			Title:    safeStr(project.VideoPara),
			VideoUrl: string(project.ProjectVideos), // Will be populated from project videos if available
		},
		GoogleMapLink: safeStr(project.ProjectLocationURL),
	}
	return &webCards, nil
}

func parseFloor(floorStr *string) *int64 {
	if floorStr == nil {
		return nil
	}
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(*floorStr, -1)
	if len(matches) == 0 {
		return nil
	}
	n, err := strconv.ParseInt(matches[len(matches)-1], 10, 64)
	if err != nil {
		return nil
	}
	return &n
}
