package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type GetProjectResponse struct {
	ProjectID    string                 `json:"project_id"`
	ProjectName  string                 `json:"project_name"`
	Description  string                 `json:"description"`
	Status       enums.ProjectStatus    `json:"status"`
	MinPrice     int                    `json:"min_price"`
	MaxPrice     int                    `json:"max_price"`
	PriceUnit    string                 `json:"price_unit"`
	TimelineInfo schema.TimelineInfo    `json:"timeline_info"`
	MetaInfo     schema.SEOMeta         `json:"meta_info"`
	WebCards     schema.ProjectWebCards `json:"web_cards"`
	LocationInfo schema.LocationInfo    `json:"location_info"`
	IsFeatured   bool                   `json:"is_featured"`
	IsPremium    bool                   `json:"is_premium"`
	IsPriority   bool                   `json:"is_priority"`
}

type AddProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type UpdateProjectResponse struct {
	ProjectID string `json:"project_id"`
}

func GetProjectFromEnt(project *ent.Project) *GetProjectResponse {

	return &GetProjectResponse{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		Description: project.Description,
		Status:      project.Status,
		MinPrice:    project.MinPrice,
		MaxPrice:    project.MaxPrice,
	}
}
