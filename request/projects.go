package request

import (
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type AddProjectRequest struct {
	ProjectName string `json:"project_name" validate:"required"`
	ProjectURL  string `json:"project_url" validate:"required"`
	ProjectType string `json:"project_type" validate:"required"`
	Locality    string `json:"locality" validate:"required"`
	ProjectCity string `json:"project_city" validate:"required"`
	DeveloperID string `json:"developer_id" validate:"required"`
}

type UpdateProjectRequest struct {
	ProjectID    string                 `json:"project_id" validate:"required"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Status       enums.ProjectStatus    `json:"status,omitempty"`
	PriceUnit    string                 `json:"price_unit,omitempty"`
	TimelineInfo schema.TimelineInfo    `json:"timeline_info,omitempty"`
	MetaInfo     schema.SEOMeta         `json:"meta_info,omitempty"`
	WebCards     schema.ProjectWebCards `json:"web_cards,omitempty"`
	LocationInfo schema.LocationInfo    `json:"location_info,omitempty"`
	IsFeatured   bool                   `json:"is_featured,omitempty"`
	IsPremium    bool                   `json:"is_premium,omitempty"`
	IsPriority   bool                   `json:"is_priority,omitempty"`
	IsDeleted    bool                   `json:"is_deleted,omitempty"`
}
