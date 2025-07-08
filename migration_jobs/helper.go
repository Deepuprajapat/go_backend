package migration_jobs

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

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

func parsePropertyImagesFromPropertyImages(propertyImages *string) ([]string, error) {
	if propertyImages == nil {
		return []string{}, nil
	}

	var projectImagesList []string
	if err := json.Unmarshal([]byte(*propertyImages), &projectImagesList); err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal project images: %s", *propertyImages)
		return nil, fmt.Errorf("failed to unmarshal project images: %w", err)
	}

	if len(projectImagesList) == 0 {
		log.Info().Msgf("No project images found")
		return []string{}, nil
	}

	propertyImagesList := []string{}
	for _, image := range projectImagesList {
		log.Info().Msgf("Parsing property image %+v", image)
		propertyImagesList = append(propertyImagesList, image)
	}
	return propertyImagesList, nil
}

func decodeJavaSerialized(blob []byte) []byte {
	if len(blob) == 0 {
		return []byte{}
	}

	// Try to detect if it's actually human readable (not binary)
	if utf8.Valid(blob) {
		// Maybe it's a plain string (rare case)
		return blob
	}

	// Extract string from Java serialized ArrayList
	// Find the string content after the 't' marker
	for i := 0; i < len(blob)-1; i++ {
		if blob[i] == 0x74 { // 't' marker for string in Java serialization
			strStart := i + 3 // Skip 't' and 2-byte length
			if strStart < len(blob) {
				var content []byte
				// Read until we hit a null byte or end
				for j := strStart; j < len(blob); j++ {
					if blob[j] == 0 {
						break
					}
					content = append(content, blob[j])
				}
				if len(content) > 0 {
					// If it looks like a YouTube video ID, make it a full URL
					str := string(content)
					if len(str) == 11 { // YouTube IDs are 11 characters
						return []byte("https://www.youtube.com/watch?v=" + str)
					}
					return content
				}
			}
		}
	}

	log.Warn().Msg("Failed to extract video URL from Java serialized data")
	return []byte("VIDEO_URL_DECODE_FAILED")
}
