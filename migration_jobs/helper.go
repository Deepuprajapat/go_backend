package migration_jobs

import(
	"fmt"
	"strings"
	"encoding/json"
	"strconv"
	"regexp"
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

func parsePhoneJSONToInt32(s *string) (int32, error) {
	if s == nil {
		return 0, fmt.Errorf("phone string is nil")
	}

	var phones []string
	if err := json.Unmarshal([]byte(*s), &phones); err != nil {
		return 0, fmt.Errorf("failed to unmarshal phone JSON: %w", err)
	}

	if len(phones) == 0 {
		return 0, fmt.Errorf("no phone number found")
	}

	// Clean up: remove dashes and whitespace
	phoneClean := strings.ReplaceAll(phones[0], "-", "")
	phoneClean = strings.TrimSpace(phoneClean)

	// Convert to int32
	phoneInt, err := strconv.ParseInt(phoneClean, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert phone to int32: %w", err)
	}

	return int32(phoneInt), nil
}


func extractNumericPhone(phoneRaw string) int32 {
	re := regexp.MustCompile(`\d+`)
	allNums := re.FindAllString(phoneRaw, -1)
	joined := strings.Join(allNums, "")
	num, _ := strconv.Atoi(joined)
	return int32(num)
}
