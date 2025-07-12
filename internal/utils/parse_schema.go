package utils



import (
	"encoding/json"
	"fmt"
)

func GetOgImageFromSchema(schemaJson string) string {
	if schemaJson == "" {
		return ""
	}

	// Use generic map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(schemaJson), &data)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return ""
	}

	image, exists := data["image"]
	if !exists {
		return ""
	}

	switch v := image.(type) {
	case string:
		// case: "image": "https://example.com/img.jpg"
		return v

	case []interface{}:
		// case: "image": ["img1", "img2"]
		if len(v) > 0 {
			if first, ok := v[0].(string); ok {
				return first
			}
		}

	case map[string]interface{}:
		// case: "image": { "url": "https://example.com/img.jpg" }
		if url, ok := v["url"].(string); ok {
			return url
		}
	}

	return ""
}
