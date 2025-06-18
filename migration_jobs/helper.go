package migration_jobs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
func parsePhoneJSONToInt32(s *string) (*int32, error) {
	if s == nil {
		var zero int32 = 0
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

	// Convert to int64 first since phone numbers can be large
	phoneInt64, err := strconv.ParseInt(phoneClean, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert phone to int64: %w", err)
	}

	// Convert to int32, capping at max int32 if needed
	var result int32
	if phoneInt64 > 2147483647 { // Max int32
		result = 2147483647
	} else {
		result = int32(phoneInt64)
	}

	return &result, nil
}
