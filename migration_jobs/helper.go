package migration_jobs

import (
	"encoding/json"
	"fmt"
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

// func parsePhoneJSONToInt32(s *string) (int32, error) {
// 	if s == nil {
// 		return 0, fmt.Errorf("phone string is nil")
// 	}

// 	var phones []string
// 	if err := json.Unmarshal([]byte(*s), &phones); err != nil {
// 		return 0, fmt.Errorf("failed to unmarshal phone JSON: %w", err)
// 	}

// 	if len(phones) == 0 {
// 		return 0, fmt.Errorf("no phone number found")
// 	}

// 	// Clean up: remove dashes and whitespace
// 	phoneClean := strings.ReplaceAll(phones[0], "-", "")
// 	phoneClean = strings.TrimSpace(phoneClean)

// 	// Convert to int32
// 	phoneInt, err := strconv.ParseInt(phoneClean, 10, 32)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to convert phone to int32: %w", err)
// 	}

// 	return int32(phoneInt), nil
// }

//	func extractNumericPhone(phoneRaw string) int32 {
//		re := regexp.MustCompile(`\d+`)
//		allNums := re.FindAllString(phoneRaw, -1)
//		joined := strings.Join(allNums, "")
//		num, _ := strconv.Atoi(joined)
//		return int32(num)
//	}
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

func parseWebCardsFromProject(project *LProject) (*schema.WebCards, error) {

	webCards := schema.WebCards{
		PropertyDetails: schema.PropertyDetails{
			PropertyType:   *project.ProjectConfigurations,
			FurnishingType: *project.ProjectConfigurations,
			ListingType:    *project.ProjectConfigurations,
		},
	}
	return &webCards, nil
}
