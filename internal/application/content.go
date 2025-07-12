package application

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/VI-IM/im_backend_go/ent"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (c *application) GetProjectByCanonicalURL(ctx context.Context, url string) (*ent.Project, *imhttp.CustomError) {
	project, err := c.repo.GetProjectByCanonicalURL(ctx, url)
	if err != nil {
		fmt.Println("Error getting project by canonical URL: ", err)
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get project by canonical URL", err.Error())
	}
	return project, nil
}

func (c *application) GetPropertyByName(ctx context.Context, url string) (*ent.Property, *imhttp.CustomError) {

	cleanUrl := strings.Replace(url, "https://investmango.com/", "", -1)
	cleanUrl = strings.Replace(cleanUrl, "https://www.investmango.com/", "", -1)
	logger.Get().Debug().Msg("cleanUrl: " + cleanUrl)
	// Extract property name
	if strings.Contains(cleanUrl, "propertyforsale/") {

		parts := strings.Split(cleanUrl, "/")
		logger.Get().Debug().Msg("parts: " + strings.Join(parts, ", "))
		if len(parts) >= 2 {
			PropertyURL := parts[1]
			logger.Get().Debug().Msg("propertyName: " + PropertyURL)
			property, err := c.repo.GetPropertyByCanonicalURL(ctx, PropertyURL)
			if err != nil {
				return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get property", err.Error())
			}
			return property, nil
		}
	}
	return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid URL format", "Invalid URL format")
}
